package iamapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/testsign"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

const (
	DefaultTestAccessKey = "test"
	DefaultTestSecretKey = "test"
)

var w *httptest.ResponseRecorder
var router = mux.NewRouter()

func TestMain(m *testing.M) {
	var err error
	uleveldb.DBClient, err = uleveldb.OpenDb("./test")
	if err != nil {
		return
	}
	defer uleveldb.DBClient.Close()
	NewIamApiServer(router)
	s3api.NewS3Server(router)
	os.Exit(m.Run())
}

func reqTest(r *http.Request) *httptest.ResponseRecorder {
	// mock a response logger
	w = httptest.NewRecorder()
	// Let the server process the mock request and record the returned response content
	router.ServeHTTP(w, r)
	fmt.Println(w.Body.String())
	return w
}
func TestIamApiServer_AddUser(t *testing.T) {
	// test cases with inputs and expected result for User.
	testCases := []struct {
		isRemove  bool
		accessKey string
		secretKey string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		// Fetching the entire User and validating its contents.
		{
			isRemove:           true,
			accessKey:          "test1",
			secretKey:          "test1234",
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong The same user name already exists ..
		{
			isRemove:           false,
			accessKey:          "test1",
			secretKey:          "test1234",
			expectedRespStatus: http.StatusConflict,
		},
		// Test case - 3.
		// error  access key length should be between 3 and 20.
		{
			isRemove:           false,
			accessKey:          "1",
			secretKey:          "test1234",
			expectedRespStatus: http.StatusInternalServerError,
		},
		// Test case - 4.
		// error  secret key length should be between 3 and 20.
		{
			isRemove:           false,
			accessKey:          "test2",
			secretKey:          "1",
			expectedRespStatus: http.StatusInternalServerError,
		},
	}
	addUrl := "http://127.0.0.1:9985/admin/v1/add-user"
	removeUrl := "http://127.0.0.1:9985/admin/v1/remove-user"
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		if testCase.isRemove {
			reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPost, removeUrl+"?accessKey="+testCase.accessKey, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
			result1 := reqTest(reqPutBucket)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result1.Code)
			}
		}
		// mock an HTTP request
		reqPutUser := testsign.MustNewSignedV4Request(http.MethodPost, addUrl+"?accessKey="+testCase.accessKey+"&secretKey="+testCase.secretKey, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
		result := reqTest(reqPutUser)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
	}
}

//func TestIamApiServer_GetUserList2(t *testing.T) {
//	u := "http://127.0.0.1:9985/admin/v1/list-all-sub-users"
//	req := testsign.MustNewSignedV4Request(http.MethodGet, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//	w = httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//	assert.Equal(t, http.StatusOK, w.Code)
//	fmt.Println(w.Body.String())
//}

func TestIamApiServer_GetUserList(t *testing.T) {
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		isPut     bool
		accessKey string
		secretKey string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{

		{
			isPut:              true,
			accessKey:          "adminTest1",
			secretKey:          "adminTest1",
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			isPut:              true,
			accessKey:          "adminTest2",
			secretKey:          "adminTest2",
			expectedRespStatus: http.StatusOK,
		},
		{
			isPut:              false,
			accessKey:          "adminTest3",
			secretKey:          "adminTest3",
			expectedRespStatus: http.StatusOK,
		},
	}
	addUrl := "http://127.0.0.1:9985/admin/v1/add-user"
	queryUrl := "http://127.0.0.1:9985/admin/v1/list-all-sub-users"
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		// mock an HTTP request
		if testCase.isPut {
			// mock an HTTP request
			reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPost, addUrl+"?accessKey="+testCase.accessKey+"&secretKey="+testCase.secretKey, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
			result := reqTest(reqPutBucket)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
			}
		}
		reqListUser := testsign.MustNewSignedV4Request(http.MethodGet, queryUrl, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
		result := reqTest(reqListUser)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		var resp ListUsersResponse
		utils.XmlDecoder(result.Body, &resp, reqListUser.ContentLength)
		fmt.Printf("case:%v  list:%v\n", i+1, resp)
	}

}

func TestIamApiServer_UserInfo(t *testing.T) {
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		isPut              bool
		accessKey          string
		secretKey          string
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		{
			isPut:              true,
			accessKey:          "adminTest1",
			secretKey:          "adminTest1",
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		{
			isPut:              true,
			accessKey:          "adminTest2",
			secretKey:          "adminTest2",
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 3.
		// The specified user does not exist
		{
			isPut:              false,
			accessKey:          "adminTest3",
			secretKey:          "adminTest3",
			expectedRespStatus: http.StatusConflict,
		},
	}
	u := "http://127.0.0.1:9985/admin/v1/user-info"
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		reqPutBucket := testsign.MustNewSignedV4Request(http.MethodGet, u+"?accessKey="+testCase.accessKey, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
		result1 := reqTest(reqPutBucket)
		if result1.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result1.Code)
		}

	}
}

func TestIamApiServer_RemoveUser(t *testing.T) {
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		isPut              bool
		accessKey          string
		secretKey          string
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		{
			isPut:              true,
			accessKey:          "adminTest1",
			secretKey:          "adminTest1",
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		{
			isPut:              true,
			accessKey:          "adminTest2",
			secretKey:          "adminTest2",
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 3.
		// The specified user does not exist
		{
			isPut:              false,
			accessKey:          "adminTest3",
			secretKey:          "adminTest3",
			expectedRespStatus: http.StatusConflict,
		},
	}
	removeUrl := "http://127.0.0.1:9985/admin/v1/remove-user"
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		// mock an HTTP request
		//if testCase.isPut {
		//	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		//	result1 := reqTest(reqPutBucket)
		//	if result1.Code != testCase.expectedRespStatus {
		//		t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result1.Code)
		//	}
		//}
		reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPost, removeUrl+"?accessKey="+testCase.accessKey, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
		result1 := reqTest(reqPutBucket)
		if result1.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result1.Code)
		}

	}

}

//func TestIamApiServer_RemoveUser(t *testing.T) {
//	u := "http://127.0.0.1:9985/admin/v1/remove-user"
//	req := testsign.MustNewSignedV4Request(http.MethodPost, u+"?accessKey=test1", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//	w = httptest.NewRecorder()
//	router.ServeHTTP(w, req)
//	assert.Equal(t, http.StatusOK, w.Code)
//	fmt.Println(w.Body.String())
//}

func TestIamApiServer_PutUserPolicy(t *testing.T) {
	urlValues := make(url.Values)
	policy := `{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::test1/*"}]}`
	urlValues.Set("policyDocument", policy)
	urlValues.Set("userName", "test1")
	urlValues.Set("policyName", "read2")
	u := "http://127.0.0.1:9985/admin/v1/put-sub-user-policy?"
	req := testsign.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
	req := testsign.MustNewSignedV4Request(http.MethodGet, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
	req := testsign.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body.String())
}
func TestIamApiServer_ListUserPolicy(t *testing.T) {
	urlValues := make(url.Values)
	urlValues.Set("userName", "test1")
	u := "http://127.0.0.1:9985/admin/v1/list-sub-user-policy?"
	req := testsign.MustNewSignedV4Request(http.MethodGet, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body.String())
}

//func TestIamApiServer_ChangePassword(t *testing.T) {
//	urlValues := make(url.Values)
//	urlValues.Set("newPassword", "test2222")
//	u := "http://127.0.0.1:9985/admin/v1/change-password?"
//	req := testsign.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", "test1", "test12345", t)
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
//
//func TestIamApiServer_SetStatus(t *testing.T) {
//	urlValues := make(url.Values)
//	urlValues.Set("userName", "test1")
//	urlValues.Set("status", string(iam.AccountDisabled))
//	//urlValues.Set("status", string(iam.AccountEnabled))
//	u := "http://127.0.0.1:9985/admin/v1/update-accessKey_status?"
//	req := testsign.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", "test1", "test2222", t)
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
