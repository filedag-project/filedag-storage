package s3api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iamapi"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/filedag-project/filedag-storage/objectservice/utils/httpstats"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var w *httptest.ResponseRecorder
var router = mux.NewRouter()
var normalUser, normalSecret = "testA", "testA12345"
var (
	nonExistBucket, wrongAccessKey, wrongSecretKey = "/nonexist", "wrongAccessKey", "wrongSecretKey"
	nonExistObject                                 = "/nonexist"
)

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
	err = authSys.Iam.AddUser(context.TODO(), normalUser, normalSecret, 999999)
	if err != nil {
		println(err)
		return
	}
	poolCli, done := client.NewMockPoolClient(&testing.T{})
	defer done()
	dagServ := merkledag.NewDAGService(blockservice.New(poolCli, offline.Exchange(poolCli)))
	storageSys := store.NewStorageSys(context.TODO(), dagServ, db)
	bmSys := store.NewBucketMetadataSys(db)
	storageSys.SetNewBucketNSLock(bmSys.NewNSLock)
	storageSys.SetHasBucket(bmSys.HasBucket)
	bmSys.SetEmptyBucket(storageSys.EmptyBucket)
	cleanData := func(accessKey string) {
		ctx := context.Background()
		bkts, err := bmSys.GetAllBucketsOfUser(ctx, accessKey)
		if err != nil {
			log.Errorf("GetAllBucketsOfUser error: %v", err)
		}
		for _, bkt := range bkts {
			if err = storageSys.CleanObjectsInBucket(ctx, bkt.Name); err != nil {
				log.Errorf("CleanObjectsInBucket error: %v", err)
				continue
			}
			if err = bmSys.DeleteBucket(ctx, bkt.Name); err != nil {
				log.Errorf("DeleteBucket error: %v", err)
			}
		}
	}

	iamapi.NewIamApiServer(router, authSys, httpstats.NewHttpStatsSys(db), cleanData, func(ctx context.Context, accessKey string) []store.BucketInfo { return nil })
	NewS3Server(router, authSys, bmSys, storageSys, httpstats.NewHttpStatsSys(db))
	os.Exit(m.Run())
}
func reqTest(r *http.Request) *httptest.ResponseRecorder {
	// mock a response logger
	w = httptest.NewRecorder()
	// Let the server process the mock request and record the returned response content
	router.ServeHTTP(w, r)
	return w
}

func TestS3ApiServer_PutBucketHandler(t *testing.T) {
	bucketName := "/testbucketput"
	// test cases with inputs and expected result for PutBucket.
	testCases := []struct {
		name       string
		bucketName string
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		{
			name:               "root user list bucket",
			bucketName:         bucketName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			name:               "wrong accessKey",
			bucketName:         bucketName,
			accessKey:          wrongAccessKey,
			secretKey:          wrongSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
		{
			name:               "normal accessKey",
			bucketName:         bucketName + "normal",
			accessKey:          normalUser,
			secretKey:          normalSecret,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "wrong bucket name",
			bucketName:         "/Aq",
			accessKey:          normalUser,
			secretKey:          normalSecret,
			expectedRespStatus: http.StatusBadRequest,
		},
	}
	// Iterating over the cases,
	for _, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.name, func(t *testing.T) {
			reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			result := reqTest(reqPutBucket)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})
	}

}
func TestS3ApiServer_HeadBucketHandler(t *testing.T) {
	bucketName := "/testbuckethead"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultBucket := reqTest(reqPutBucket)
	if resultBucket.Code != http.StatusOK {
		t.Fatalf("reqPutNormalBucket : Expected the response status to be `%d`, but instead found `%d`", 200, resultBucket.Code)
	}
	reqPutNormalBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal", 0, nil, "s3", normalUser, normalSecret, t)
	resultNormalBucket := reqTest(reqPutNormalBucket)
	if resultNormalBucket.Code != http.StatusOK {
		t.Fatalf("reqPutNormalBucket : Expected the response status to be `%d`, but instead found `%d`", 200, resultNormalBucket.Code)
	}
	// test cases with inputs and expected result for HeadBucket.
	testCases := []struct {
		name       string
		bucketName string
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.

		{
			name:               "root user head bucket",
			bucketName:         bucketName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "normal user head bucket",
			bucketName:         bucketName + "normal",
			accessKey:          normalUser,
			secretKey:          normalSecret,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			name:               "wrong accessKey.",
			bucketName:         bucketName,
			accessKey:          wrongAccessKey,
			secretKey:          wrongSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
		// Test case - 3.
		{
			name:               "non-exist bucket",
			bucketName:         nonExistBucket,
			accessKey:          normalUser,
			secretKey:          normalSecret,
			expectedRespStatus: http.StatusNotFound,
		},
	}
	// Iterating over the cases,
	for _, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.name, func(t *testing.T) {
			reqHeadBucket := utils.MustNewSignedV4Request(http.MethodHead, testCase.bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			result2 := reqTest(reqHeadBucket)
			if result2.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result2.Code)
			}
		})

	}

}
func TestS3ApiServer_ListBucketHandler(t *testing.T) {
	bucketName := "/testbucketlist"
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name       string
		bucketName string
		has        bool
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{

		// Test case - 1.
		{
			name:               "root user list bucket(0)",
			has:                false,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		{
			name:               "root user list bucket(1) ",
			bucketName:         bucketName + "has",
			has:                true,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 3.
		// wrong accessKey.
		{
			name:               "wrong accessKey",
			has:                false,
			bucketName:         bucketName + "wrong",
			accessKey:          wrongAccessKey,
			secretKey:          wrongSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},

		{
			name:               "normal user list  bucket",
			bucketName:         bucketName + "normal",
			has:                true,
			accessKey:          normalUser,
			secretKey:          normalSecret,
			expectedRespStatus: http.StatusOK,
		},
	}
	// Iterating over the cases,
	for i, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.has {
				reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
				result1 := reqTest(reqPutBucket)
				if result1.Code != testCase.expectedRespStatus {
					t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result1.Code)
				}
			}
			reqListBucket := utils.MustNewSignedV4Request(http.MethodGet, "/", 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			result := reqTest(reqListBucket)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
			}
			var resp response.ListAllMyBucketsResult
			utils.XmlDecoder(result.Body, &resp, reqListBucket.ContentLength)
			fmt.Printf("case:%v  list:%v\n", i+1, resp)
		})
	}

}
func TestS3ApiServer_DeleteBucketHandler(t *testing.T) {
	bucketName := "/testbucketdel"
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name       string
		bucketName string
		has        bool
		accessKey  string
		secretKey  string
		// expected output.
		expectedPutRespStatus int // expected response status body.
		expectedDelRespStatus int
	}{

		// Test case - 1.
		{
			name:                  "root user del empty bucket",
			has:                   true,
			bucketName:            bucketName + "empty",
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusOK,
			expectedDelRespStatus: http.StatusNoContent,
		},
		{
			name:                  "root user del not-empty bucket",
			bucketName:            bucketName + "notempty",
			has:                   true,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusOK,
			expectedDelRespStatus: http.StatusConflict,
		},
		// Test case - 2.
		// non exist bucket.
		{
			name:                  "root user del non-exist bucket",
			bucketName:            bucketName + "nonexist",
			has:                   false,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusOK,
			expectedDelRespStatus: http.StatusNotFound,
		},

		// wrong accessKey.
		{
			name:                  "wrong accessKey.",
			bucketName:            bucketName + "wrong",
			has:                   true,
			accessKey:             "1",
			secretKey:             "1",
			expectedPutRespStatus: http.StatusForbidden,
			expectedDelRespStatus: http.StatusForbidden,
		},
		// Test case - 3.
		{
			name:                  "normal user del empty bucket",
			has:                   true,
			bucketName:            bucketName + "emptynormal",
			accessKey:             normalUser,
			secretKey:             normalSecret,
			expectedPutRespStatus: http.StatusOK,
			expectedDelRespStatus: http.StatusNoContent,
		},
		{
			name:                  "normal user del not-empty bucket",
			bucketName:            bucketName + "notemptynormal",
			has:                   true,
			accessKey:             normalUser,
			secretKey:             normalSecret,
			expectedPutRespStatus: http.StatusOK,
			expectedDelRespStatus: http.StatusConflict,
		},
		// Test case - 2.
		// non exist bucket.
		{
			name:                  "normal user del non-exist bucket",
			bucketName:            bucketName + "nonexistnormal",
			has:                   false,
			accessKey:             normalUser,
			secretKey:             normalSecret,
			expectedPutRespStatus: http.StatusOK,
			expectedDelRespStatus: http.StatusNotFound,
		},
	}

	// Iterating over the cases,
	for _, testCase := range testCases {
		// mock an HTTP request
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.has {
				reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
				result1 := reqTest(reqPutBucket)
				if result1.Code != testCase.expectedPutRespStatus {
					t.Fatalf("Case %s:reqPutBucket Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedPutRespStatus, result1.Code)
				}
				if strings.Contains(testCase.bucketName, "notempty") {
					r1 := "1234567"
					reqPutNotEmpty := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+"/"+"object1", int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", testCase.accessKey, testCase.secretKey, t)
					resultPutNotEmpty := reqTest(reqPutNotEmpty)
					if resultPutNotEmpty.Code != testCase.expectedPutRespStatus {
						t.Fatalf("Case %s: reqPutNotEmpty Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedPutRespStatus, resultPutNotEmpty.Code)
					}
				}
			}
			reqDeleteBucket := utils.MustNewSignedV4Request(http.MethodDelete, testCase.bucketName, 0,
				nil, "s3", testCase.accessKey, testCase.secretKey, t)
			result := reqTest(reqDeleteBucket)
			if result.Code != testCase.expectedDelRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedDelRespStatus, result.Code)
			}
		})
	}

}

//func TestS3ApiServer_GetBucketLocationHandler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test22"
//	//req.Header.Set("Content-Type", "text/plain")
//	req := testsign.MustNewSignedV4Request(http.MethodHead, u+"?location", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
//
//func TestS3ApiServer_GetBucketAclHandler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test"
//	req := testsign.MustNewSignedV4Request(http.MethodGet, u+"?acl=", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
//func TestS3ApiServer_PutBucketAclHandler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test"
//	a := `<?xml version="1.0" encoding="UTF-8"?>
//<AccessControlPolicy xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
//  <Owner>
//    <ID>*** Owner-Canonical-User-ID ***</ID>
//    <DisplayName>owner-display-name</DisplayName>
//  </Owner>
//  <AccessControlList>
//    <Grant>
//      <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
//               xsi:type="Canonical User">
//        <ID>*** Owner-Canonical-User-ID ***</ID>
//        <DisplayName>display-name</DisplayName>
//      </Grantee>
//      <Permission>FULL_CONTROL</Permission>
//    </Grant>
//  </AccessControlList>
//</AccessControlPolicy>`
//	req := testsign.MustNewSignedV4Request(http.MethodPut, u+"?acl=", int64(len(a)), strings.NewReader(a), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
/*func TestS3_PutBucketHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test22"

	req := testsign.MustNewSignedV4Request(http.MethodPut, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	//req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}*/
