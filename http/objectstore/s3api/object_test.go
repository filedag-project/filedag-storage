package s3api

import (
	"bytes"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/testsign"
	"net/http"
	"net/url"
	"testing"
)

const (
	DefaultTestAccessKey = "test"
	DefaultTestSecretKey = "test"
)

func TestS3ApiServer_PutObjectHandler(t *testing.T) {
	bucketName := "/testbucket"
	objectName := "/testobject"
	r1 := "1234567"
	copySourceHeader := http.Header{}
	copySourceHeader.Set("X-Amz-Copy-Source", "somewhere")
	invalidMD5Header := http.Header{}
	invalidMD5Header.Set("Content-Md5", "42")
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		bucketName string
		objectName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			data:               []byte(r1),
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			data:               []byte(r1),
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
		{
			bucketName:         "/11",
			objectName:         objectName,
			data:               []byte(r1),
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
		{
			bucketName:         bucketName,
			objectName:         objectName,
			data:               []byte(r1),
			header:             copySourceHeader,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusBadRequest,
		},
		{
			bucketName:         bucketName,
			objectName:         objectName,
			data:               []byte(r1),
			header:             invalidMD5Header,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusBadRequest,
		},
	}
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := testsign.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+testCase.objectName, int64(len(r1)), bytes.NewReader(testCase.data), "s3", testCase.accessKey, testCase.secretKey, t)
		// Add test case specific headers to the request.
		addCustomHeaders(req, testCase.header)
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: put:%v\n", i+1, result.Body.String())
	}

}
func TestS3ApiServer_GetObjectHandler(t *testing.T) {
	bucketName := "/testbucket"
	objectName := "/testobject"

	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		bucketName string
		objectName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
		{
			bucketName:         "/11",
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
	}
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	a := reqTest(reqPutBucket)
	fmt.Println("putbucket:", a.Body.String())
	r1 := "1234567"
	reqputObject := testsign.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := testsign.MustNewSignedV4Request(http.MethodGet, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		// Add test case specific headers to the request.
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: get:%v\n", i+1, result.Body.String())
	}

}
func TestS3ApiServer_HeadObjectHandler(t *testing.T) {
	bucketName := "/testbucket"
	objectName := "/testobject"

	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		bucketName string
		objectName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
		{
			bucketName:         "/11",
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
	}
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := testsign.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := testsign.MustNewSignedV4Request(http.MethodHead, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		// Add test case specific headers to the request.
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: head:%v\n", i+1, result.Body.String())
	}

}
func TestS3ApiServer_DeleteObjectHandler(t *testing.T) {
	bucketName := "/testbucket"
	objectName := "/testobject"

	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		bucketName string
		objectName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNoContent,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
		{
			bucketName:         "/11",
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
	}
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := testsign.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := testsign.MustNewSignedV4Request(http.MethodDelete, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		// Add test case specific headers to the request.
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: delete:%v\n", i+1, result.Body.String())
	}

}
func TestS3ApiServer_CopyObjectHandler(t *testing.T) {
	bucketName := "/testbucket"
	objectName := "/testobject"

	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		bucketName    string
		objectName    string
		dstbucketName string
		dstobjectName string
		data          []byte
		header        http.Header
		accessKey     string
		secretKey     string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			dstbucketName:      bucketName,
			dstobjectName:      "/1.txt",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			bucketName:         bucketName,
			objectName:         objectName,
			dstbucketName:      bucketName,
			dstobjectName:      "/1.txt",
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
		{
			bucketName:         "/11",
			objectName:         objectName,
			dstbucketName:      bucketName,
			dstobjectName:      "/1.txt",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
		{
			bucketName:         "",
			objectName:         "",
			dstbucketName:      bucketName,
			dstobjectName:      "/1.txt",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusBadRequest,
		},
		{
			bucketName:         bucketName,
			objectName:         objectName,
			dstbucketName:      bucketName,
			dstobjectName:      objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusBadRequest,
		},
	}
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := testsign.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := testsign.MustNewSignedV4Request(http.MethodPut, testCase.dstbucketName+testCase.dstobjectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		req.Header.Set("X-Amz-Copy-Source", url.QueryEscape(testCase.bucketName+testCase.objectName)) // Add test case specific headers to the request.
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: copy:%v\n", i+1, result.Body.String())
	}

}
func TestS3ApiServer_ListObjectsV2Handler(t *testing.T) {
	bucketName := "/testbucket"
	objectName := "/testobject"

	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		bucketName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1.
		// Fetching the entire Bucket and validating its contents.
		{
			bucketName: bucketName,

			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			bucketName:         bucketName,
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
		{
			bucketName:         "/11",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
	}
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := testsign.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := testsign.MustNewSignedV4Request(http.MethodGet, testCase.bucketName+"?list-type=2", 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: copy:%v\n", i+1, result.Body.String())
	}

}

func TestWholeNoUserAPI(t *testing.T) {
	bucketName := "/testbucket"
	objectName := "/testobject"
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := testsign.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	reqGetObject, _ := http.NewRequest(http.MethodGet, bucketName+objectName, nil)
	fmt.Println(reqTest(reqGetObject).Body.String())
}

func TestWholeProcess(t *testing.T) {
	userName := "dean"
	secret := "dean123456"
	bucketName := "/testbucket"
	objectName := "/testobject"
	reqPutUser := testsign.MustNewSignedV4Request(http.MethodPost, "/admin/v1/add-user"+"?accessKey="+userName+"&secretKey="+secret, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println(reqTest(reqPutUser).Body.String())
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", userName, secret, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := testsign.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", userName, secret, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	req := testsign.MustNewSignedV4Request(http.MethodGet, bucketName+objectName, 0, nil, "s3", userName, secret, t)
	fmt.Println("getobject", reqTest(req).Body.String())
}
func TestS3ApiServer_PutObjectHandler2(t *testing.T) {
	bucketName := "/testbucket"
	objectName := "/testobject"
	reqPutBucket := testsign.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := testsign.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	r2 := "123456"
	reqputObject2 := testsign.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r2)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject2)
	req := testsign.MustNewSignedV4Request(http.MethodGet, bucketName+objectName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println(reqTest(req).Body.String())
}
func addCustomHeaders(req *http.Request, customHeaders http.Header) {
	for k, values := range customHeaders {
		for _, value := range values {
			req.Header.Set(k, value)
		}
	}
}

/*func TestS3_PutObjectHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test22/1.txt"
	r1, _ := ioutil.ReadFile("./object_test.go")

	req := testsign.MustNewSignedV4Request(http.MethodPut, u, int64(len(r1)), bytes.NewReader(r1), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)

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
//func TestS3ApiServer_ListObjectsV1Handler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test22"
//	req := testsign.MustNewSignedV4Request(http.MethodGet, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
//
//}
