package iamapi

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/filedag-project/filedag-storage/objectservice/utils/httpstats"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

const (
	DefaultTestAccessKey  = auth.DefaultAccessKey
	DefaultTestSecretKey  = auth.DefaultAccessKey
	defaultCap            = "99999"
	normalAccessKey       = "normalUser"
	normalSecretKey       = "normalUser"
	otherUserAccessKey    = "otherUser"
	otherUserSecretKey    = "otherUser"
	userNonExistAccessKey = "userNonExist"
	userNonExistSecretKey = "userNonExist"
)

var w *httptest.ResponseRecorder
var router = mux.NewRouter()

func TestMain(m *testing.M) {
	db, err := objmetadb.OpenDb((&testing.T{}).TempDir())
	if err != nil {
		println(err)
		return
	}
	defer db.Close()
	cred, err := auth.CreateCredentials(auth.DefaultAccessKey, auth.DefaultSecretKey)
	if err != nil {
		println(err)
		return
	}
	authSys := iam.NewAuthSys(db, cred)
	poolCli, done := client.NewMockPoolClient(&testing.T{})
	defer done()
	dagServ := merkledag.NewDAGService(blockservice.New(poolCli, offline.Exchange(poolCli)))
	storageSys := store.NewStorageSys(context.TODO(), dagServ, db)
	bmSys := store.NewBucketMetadataSys(db)
	bucketInfoFunc := func(ctx context.Context, accessKey string) []store.BucketInfo {
		var bucketInfos []store.BucketInfo
		bkts, err := bmSys.GetAllBucketsOfUser(ctx, accessKey)
		if err != nil {
			fmt.Printf("GetAllBucketsOfUser error: %v\n", err)
			return bucketInfos
		}
		for _, bkt := range bkts {
			info, err := storageSys.GetBucketInfo(ctx, bkt.Name)
			if err != nil {
				return nil
			}
			bucketInfos = append(bucketInfos, info)
		}
		return bucketInfos
	}

	NewIamApiServer(router, authSys, httpstats.NewHttpStatsSys(db), func(accessKey string) {}, bucketInfoFunc)
	reqPutUserOtherUrl := addUserUrl(otherUserAccessKey, otherUserSecretKey, defaultCap)
	reqPutUserOther := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserOtherUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, &testing.T{})
	resultOther := reqTest(reqPutUserOther)
	if resultOther.Code != http.StatusOK {
		panic(fmt.Sprintf("add user fail %v,%v", resultOther.Code, resultOther.Body.String()))
	}
	reqPutUserNormalUrl := addUserUrl(normalAccessKey, normalSecretKey, defaultCap)
	reqPutUserNormal := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserNormalUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, &testing.T{})
	resultNormal := reqTest(reqPutUserNormal)
	if resultNormal.Code != http.StatusOK {
		panic(fmt.Sprintf("add user fail %v,%v", resultNormal.Code, resultNormal.Body.String()))
	}
	//s3api.NewS3Server(router)
	os.Exit(m.Run())
}

func reqTest(r *http.Request) *httptest.ResponseRecorder {
	// mock a response logger
	w = httptest.NewRecorder()
	// Let the server process the mock request and record the returned response content
	router.ServeHTTP(w, r)
	//fmt.Println(w.Body.String())
	return w
}
func addUserUrl(username, secret, cap string) string {
	addUrl := "http://127.0.0.1:9985/admin/v1/add-user?"

	urlValues := make(url.Values)
	urlValues.Set(accessKey, username)
	urlValues.Set(secretKey, secret)
	urlValues.Set(capacity, cap)
	return addUrl + urlValues.Encode()
}

func TestIamApiServer_AddUser(t *testing.T) {
	// test cases with inputs and expected result for AddUser.
	testCases := []struct {
		name          string
		credAccessKey string
		credSecretKey string
		accessKey     string
		secretKey     string
		cap           string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "add normal user",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          "test1",
			secretKey:          "test1234",
			cap:                defaultCap,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "The same user name already exists",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          "test1",
			secretKey:          "test1234",
			cap:                defaultCap,
			expectedRespStatus: http.StatusConflict,
		},
		{
			name:               "access key length should be between 3 and 20.",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          "1",
			secretKey:          "test1234",
			cap:                defaultCap,
			expectedRespStatus: http.StatusBadRequest,
		},
		{
			name:               "secret key length should be between 3 and 20.",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          "test2",
			secretKey:          "1",
			cap:                defaultCap,
			expectedRespStatus: http.StatusBadRequest,
		},
		{
			name:               "use normal user add user",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalAccessKey,
			accessKey:          "test2",
			secretKey:          "12345647",
			cap:                defaultCap,
			expectedRespStatus: http.StatusForbidden,
		},
	}
	// Iterating over the cases, fetching the result validating the response.
	for _, testCase := range testCases {
		// mock an HTTP request
		// add user
		t.Run(testCase.accessKey, func(t *testing.T) {
			ur := addUserUrl(testCase.accessKey, testCase.secretKey, testCase.cap)
			reqPutUser := utils.MustNewSignedV4Request(http.MethodPost, ur, 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result := reqTest(reqPutUser)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})
	}
}

//func TestIamApiServer_GetUserList(t *testing.T) {
//	// test cases with inputs and expected result for Bucket.
//	testCases := []struct {
//		isPut     bool
//		accessKey string
//		secretKey string
//		cap       string
//		// expected output.
//		expectedRespStatus int // expected response status body.
//	}{
//
//		{
//			isPut:              true,
//			accessKey:          "adminTest1",
//			secretKey:          "adminTest1",
//			cap:                defaultCap,
//			expectedRespStatus: http.StatusOK,
//		},
//		// Test case - 1.
//		// Fetching the entire Bucket and validating its contents.
//		{
//			isPut:              true,
//			accessKey:          "adminTest2",
//			secretKey:          "adminTest2",
//			cap:                defaultCap,
//			expectedRespStatus: http.StatusOK,
//		},
//		{
//			isPut:              false,
//			accessKey:          "adminTest3",
//			secretKey:          "adminTest3",
//			cap:                defaultCap,
//			expectedRespStatus: http.StatusOK,
//		},
//	}
//	addUrl := "http://127.0.0.1:9985/admin/v1/add-user"
//	queryUrl := "http://127.0.0.1:9985/admin/v1/list-all-sub-users"
//	// Iterating over the cases, fetching the object validating the response.
//	for i, testCase := range testCases {
//		// mock an HTTP request
//		if testCase.isPut {
//			// mock an HTTP request
//			// add user
//			reqPutBucket := utils.MustNewSignedV4Request(http.MethodPost, addUrl+"?accessKey="+testCase.accessKey+"&secretKey="+testCase.secretKey+"&capacity="+testCase.cap, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//			result := reqTest(reqPutBucket)
//			if result.Code != testCase.expectedRespStatus {
//				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
//			}
//		}
//		// list user
//		reqListUser := utils.MustNewSignedV4Request(http.MethodGet, queryUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//		result := reqTest(reqListUser)
//		if result.Code != testCase.expectedRespStatus {
//			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
//		}
//		var resp ListUsersResponse
//		utils.XmlDecoder(result.Body, &resp, reqListUser.ContentLength)
//		fmt.Printf("case:%v  list:%v\n", i+1, resp)
//	}
//
//}

func TestIamApiServer_AccountInfo(t *testing.T) {
	// test cases with inputs and expected result for UserInfo.
	testCases := []struct {
		name               string
		credAccessKey      string
		credSecretKey      string
		accessKey          string
		secretKey          string
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user get himself info",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "root user get normal user info",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          normalAccessKey,
			secretKey:          normalSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "root user get non-exist user info",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          userNonExistAccessKey,
			secretKey:          userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
		},
		{
			name:               "normal user get himself info",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          normalAccessKey,
			secretKey:          normalSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "normal user get other user info",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          otherUserAccessKey,
			secretKey:          otherUserSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
		{
			name:               "normal user get a non-exist user info",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          userNonExistAccessKey,
			secretKey:          userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
		},
		{
			name:               "normal user get root user info",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
	}
	u := "http://127.0.0.1:9985/console/v1/user-info?"
	// Iterating over the cases, fetching the result validating the response.
	for _, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.name, func(t *testing.T) {
			//user info
			urlValues := make(url.Values)
			urlValues.Set(accessKey, testCase.accessKey)
			userinfoReq := utils.MustNewSignedV4Request(http.MethodGet, u+urlValues.Encode(), 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result1 := reqTest(userinfoReq)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result1.Code)
			}
		})

	}
}

func TestIamApiServer_RemoveUser(t *testing.T) {
	// test cases with inputs and expected result for RemoveUser.
	himselfUserAccessKey := "himselfUser"
	himselfUserSecretKey := "himselfUser"
	removeUserAccessKey := "removeUser"
	removeUserSecretKey := "removeUser"
	reqPutUserHimselfUrl := addUserUrl(himselfUserAccessKey, himselfUserSecretKey, defaultCap)
	reqPutUserHimself := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserHimselfUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultHimself := reqTest(reqPutUserHimself)
	if resultHimself.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultHimself.Code, resultHimself.Body.String())
	}
	reqPutUserRemoveUrl := addUserUrl(removeUserAccessKey, removeUserSecretKey, defaultCap)
	reqPutUserRemove := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserRemoveUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultRemove := reqTest(reqPutUserRemove)
	if resultHimself.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultRemove.Code, resultRemove.Body.String())
	}
	testCases := []struct {
		name               string
		credAccessKey      string
		credSecretKey      string
		accessKey          string
		secretKey          string
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		{
			name:               "root user remove a user",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          removeUserAccessKey,
			secretKey:          removeUserSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		{
			name:               "root user remove a non-exist user",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          userNonExistAccessKey,
			secretKey:          userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
		},
		{
			name:               "user remove himself",
			credAccessKey:      himselfUserAccessKey,
			credSecretKey:      himselfUserSecretKey,
			accessKey:          himselfUserAccessKey,
			secretKey:          himselfUserSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "user remove other user",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          otherUserAccessKey,
			secretKey:          otherUserSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
		{
			name:               "user remove non-exist user",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          userNonExistAccessKey,
			secretKey:          userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
		},
	}

	removeUrl := "http://127.0.0.1:9985/admin/v1/remove-user?"
	// Iterating over the cases, fetching the result validating the response.
	for _, testCase := range testCases {
		// remove user
		t.Run(testCase.name, func(t *testing.T) {
			urlValues := make(url.Values)
			urlValues.Set(accessKey, testCase.accessKey)
			reqPutBucket := utils.MustNewSignedV4Request(http.MethodPost, removeUrl+urlValues.Encode(), 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result1 := reqTest(reqPutBucket)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result1.Code)
			}
		})

	}

}

// set status
func TestIamApiServer_SetStatus(t *testing.T) {
	offUserAccessKey := "offUser"
	offUserSecretKey := "offUser1234"
	otherUserOffAccessKey := "otherUserOffUser"
	otherUserOffSecretKey := "otherUserOffUser1234"
	reqPutUserOffUrl := addUserUrl(offUserAccessKey, offUserSecretKey, defaultCap)
	reqPutUserOff := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserOffUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultOff := reqTest(reqPutUserOff)
	if resultOff.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultOff.Code, resultOff.Body.String())
	}
	reqPutUserOtherOffUrl := addUserUrl(otherUserOffAccessKey, otherUserOffSecretKey, defaultCap)
	reqPutUserOtherOff := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserOtherOffUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultOtherOff := reqTest(reqPutUserOtherOff)
	if resultOtherOff.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultOtherOff.Code, resultOtherOff.Body.String())
	}
	testCases := []struct {
		name          string
		credAccessKey string
		credSecretKey string
		accessKey     string
		secretKey     string
		status        string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user set a user off",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          offUserAccessKey,
			secretKey:          offUserSecretKey,
			expectedRespStatus: http.StatusOK,
			status:             "off",
		},
		{
			name:               "root user set a user on",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          offUserAccessKey,
			secretKey:          offUserSecretKey,
			expectedRespStatus: http.StatusOK,
			status:             "on",
		},
		{
			name:               "root user set a non exist user off",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          userNonExistAccessKey,
			secretKey:          userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
			status:             "off",
		},
		{
			name:               "user set himself off",
			credAccessKey:      offUserAccessKey,
			credSecretKey:      offUserSecretKey,
			accessKey:          offUserAccessKey,
			secretKey:          offUserSecretKey,
			expectedRespStatus: http.StatusOK,
			status:             "off",
		},
		{
			// after user set status off only root user can set status on
			name:               "user set himself on",
			credAccessKey:      offUserAccessKey,
			credSecretKey:      offUserSecretKey,
			accessKey:          offUserAccessKey,
			secretKey:          offUserSecretKey,
			expectedRespStatus: http.StatusForbidden,
			status:             "on",
		},
		{
			name:               "user set a non exist user on",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          userNonExistAccessKey,
			secretKey:          userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
			status:             "on",
		},
		{
			name:               "user set other user off",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          otherUserOffAccessKey,
			secretKey:          otherUserOffSecretKey,
			expectedRespStatus: http.StatusForbidden,
			status:             "off",
		},
		{
			name:               "user set other user on",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			accessKey:          otherUserOffAccessKey,
			secretKey:          otherUserOffSecretKey,
			expectedRespStatus: http.StatusForbidden,
			status:             "on",
		},
	}
	setStatusUrl := "http://127.0.0.1:9985/admin/v1/update-accessKey_status?"
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			urlValues := make(url.Values)
			urlValues.Set(accessKey, testCase.accessKey)
			urlValues.Set(accountStatus, testCase.status)
			reqSetStatus := utils.MustNewSignedV4Request(http.MethodPost, setStatusUrl+urlValues.Encode(), 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result := reqTest(reqSetStatus)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})

	}
}
func TestIamApiServer_AccountInfos(t *testing.T) {
	// test cases with inputs and expected result for UserInfo.
	testCases := []struct {
		name               string
		credAccessKey      string
		credSecretKey      string
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user get himself info",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "root user get himself info",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
	}
	u := "http://127.0.0.1:9985/admin/v1/user-infos"
	// Iterating over the cases, fetching the result validating the response.
	for _, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.credAccessKey, func(t *testing.T) {
			//user info
			userinfoReq := utils.MustNewSignedV4Request(http.MethodGet, u, 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result1 := reqTest(userinfoReq)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result1.Code)
			}
			fmt.Println(result1.Body.String())
		})

	}
}
func TestIamApiServer_request_overview(t *testing.T) {
	// test cases with inputs and expected result for UserInfo.
	testCases := []struct {
		name               string
		credAccessKey      string
		credSecretKey      string
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user get himself info",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "root user get himself info",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
	}
	u := "http://127.0.0.1:9985/admin/v1/request-overview"
	// Iterating over the cases, fetching the result validating the response.
	for _, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.credAccessKey, func(t *testing.T) {
			//user info
			userinfoReq := utils.MustNewSignedV4Request(http.MethodGet, u, 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result1 := reqTest(userinfoReq)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result1.Code)
			}
			fmt.Println(result1.Body.String())
		})

	}
}

//change password
func TestIamApiServer_ChangePassword(t *testing.T) {
	changePassSuccessAccessKey := "changePassSuccess"
	changePassSuccessSecretKey := "changePassSuccess"
	himselfChangeSuccessAccessKey := "himselfChangeSuccess"
	himselfChangeSuccessSecretKey := "himselfChangeSuccess"

	thePassToChange := "thePassToChange"
	reqPutUserChangePassSuccessUrl := addUserUrl(changePassSuccessAccessKey, changePassSuccessSecretKey, defaultCap)
	reqPutUserChangePassSuccess := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserChangePassSuccessUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultChangePassSucces := reqTest(reqPutUserChangePassSuccess)
	if resultChangePassSucces.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultChangePassSucces.Code, resultChangePassSucces.Body.String())
	}
	reqPutUserHimselfUrl := addUserUrl(himselfChangeSuccessAccessKey, himselfChangeSuccessSecretKey, defaultCap)
	reqPutUserHimself := utils.MustNewSignedV4Request(http.MethodPost, reqPutUserHimselfUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultHimself := reqTest(reqPutUserHimself)
	if resultHimself.Code != http.StatusOK {
		t.Fatalf("add user fail %v,%v", resultHimself.Code, resultHimself.Body.String())
	}

	testCases := []struct {
		name          string
		credAccessKey string
		credSecretKey string
		oldSecretKey  string
		accessKey     string
		pass          string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		{
			name:               "root user change a user pass",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          changePassSuccessAccessKey,
			oldSecretKey:       changePassSuccessSecretKey,
			expectedRespStatus: http.StatusOK,
			pass:               thePassToChange,
		},
		{
			name:               "root user change root pass",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			oldSecretKey:       DefaultTestSecretKey,
			accessKey:          DefaultTestAccessKey,
			expectedRespStatus: http.StatusConflict,
			pass:               thePassToChange,
		},
		{
			name:               "root user change a non-exist user pass",
			credAccessKey:      DefaultTestAccessKey,
			credSecretKey:      DefaultTestSecretKey,
			accessKey:          userNonExistAccessKey,
			oldSecretKey:       userNonExistSecretKey,
			expectedRespStatus: http.StatusConflict,
			pass:               thePassToChange,
		},
		{
			name:               "normal user change himself pass",
			credAccessKey:      himselfChangeSuccessAccessKey,
			credSecretKey:      himselfChangeSuccessSecretKey,
			accessKey:          himselfChangeSuccessAccessKey,
			oldSecretKey:       himselfChangeSuccessSecretKey,
			expectedRespStatus: http.StatusOK,
			pass:               thePassToChange,
		},
		{
			name:               "normal user change other user pass",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			oldSecretKey:       otherUserSecretKey,
			accessKey:          otherUserAccessKey,
			expectedRespStatus: http.StatusForbidden,
			pass:               thePassToChange,
		},
		{
			name:               "normal user change a non-exist user pass",
			credAccessKey:      normalAccessKey,
			credSecretKey:      normalSecretKey,
			oldSecretKey:       userNonExistSecretKey,
			accessKey:          userNonExistAccessKey,
			expectedRespStatus: http.StatusConflict,
			pass:               thePassToChange,
		},
		{
			name:               "normal user change user err pass",
			credAccessKey:      himselfChangeSuccessAccessKey,
			credSecretKey:      thePassToChange,
			oldSecretKey:       thePassToChange,
			accessKey:          himselfChangeSuccessAccessKey,
			expectedRespStatus: http.StatusBadRequest,
			pass:               "dj",
		},
	}
	changePassUrl := "http://127.0.0.1:9985/console/v1/change-password?"
	for _, testCase := range testCases {
		// mock an HTTP request
		//change password
		t.Run(testCase.name, func(t *testing.T) {
			urlValues := make(url.Values)
			urlValues.Set(newSecretKey, testCase.pass)
			urlValues.Set(accessKey, testCase.accessKey)
			urlValues.Set(oldSecretKey, testCase.oldSecretKey)
			//urlValues.Set("status", string(iam.AccountDisabled
			reqChange := utils.MustNewSignedV4Request(http.MethodPost, changePassUrl+urlValues.Encode(), 0, nil, "s3", testCase.credAccessKey, testCase.credSecretKey, t)
			result := reqTest(reqChange)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})
	}
}

// todo more test case
func TestIamApiServer_PutUserPolicy(t *testing.T) {
	urlValues := make(url.Values)
	policy := `{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::test1/*"}]}`
	urlValues.Set("policyDocument", policy)
	urlValues.Set("userName", "test1")
	urlValues.Set("policyName", "read2")
	u := "http://127.0.0.1:9985/admin/v1/put-sub-user-policy?"
	req := utils.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body.String())
}
func TestIamApiServer_GetUserPolicy(t *testing.T) {
	urlValues := make(url.Values)
	urlValues.Set("userName", "test1")
	urlValues.Set("policyName", "read2")
	u := "http://127.0.0.1:9985/admin/v1/get-sub-user-policy?"
	req := utils.MustNewSignedV4Request(http.MethodGet, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body.String())
}
func TestIamApiServer_RemoveUserPolicy(t *testing.T) {
	urlValues := make(url.Values)
	urlValues.Set("userName", "test1")
	urlValues.Set("policyName", "read2")
	u := "http://127.0.0.1:9985/admin/v1/remove-sub-user-policy?"
	req := utils.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body.String())
}
func TestIamApiServer_ListUserPolicy(t *testing.T) {
	urlValues := make(url.Values)
	urlValues.Set("userName", "test1")
	u := "http://127.0.0.1:9985/admin/v1/list-sub-user-policy?"
	req := utils.MustNewSignedV4Request(http.MethodGet, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body.String())
}

//func TestIamApiServer_GetUserInfo(t *testing.T) {
//	urlValues := make(url.Values)
//	user := "test1"
//	urlValues.Set("userName", user)
//	u := "http://127.0.0.1:9985/admin/v1/user-info"
//	req := testsign.MustNewSignedV4Request(http.MethodGet, u+"?"+urlValues.Encode(), 0, nil, "s3", "test1", "test2222", t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
