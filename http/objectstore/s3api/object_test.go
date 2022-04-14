package s3api

import (
	"bytes"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/testsign"
	"net/http"
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
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		bucketName string
		objectName string
		data       []byte
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
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := testsign.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+testCase.objectName, int64(len(r1)), bytes.NewReader(testCase.data), "s3", testCase.accessKey, testCase.secretKey, t)
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: put:%v\n", i+1, result.Body.String())
	}

}

//func TestS3ApiServer_GetObjectHandler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test/1.txt"
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
//func TestS3ApiServer_CopyObjectHandler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test1/11.txt"
//	req := testsign.MustNewSignedV4Request(http.MethodPut, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//	req.Header.Set("X-Amz-Copy-Source", url.QueryEscape("/test/1.txt"))
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
//func TestS3ApiServer_HeadObjectHandler(t *testing.T) {
//	u := "http://127.0.0.1:9985/test/1.txt"
//	req := testsign.MustNewSignedV4Request(http.MethodHead, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
