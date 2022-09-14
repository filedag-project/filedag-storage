package s3api

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/datatypes"
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

const (
	DefaultTestAccessKey = auth.DefaultAccessKey
	DefaultTestSecretKey = auth.DefaultSecretKey
)

func TestS3ApiServer_PutObjectHandler(t *testing.T) {
	bucketName := "/testbucketputo"
	objectName := "/testobjectputo"
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
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	result := reqTest(reqPutBucket)
	if result.Code != http.StatusOK {
		t.Fatalf("the response status of putbucket: %d", result.Code)
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+testCase.objectName, int64(len(r1)), bytes.NewReader(testCase.data), "s3", testCase.accessKey, testCase.secretKey, t)
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
	bucketName := "/testbucketfeto"
	objectName := "/testobjectgeto"

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
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	a := reqTest(reqPutBucket)
	fmt.Println("putbucket:", a.Body.String())
	r1 := "1234567"
	reqputObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := utils.MustNewSignedV4Request(http.MethodGet, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		// Add test case specific headers to the request.
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: get:%v\n", i+1, result.Body.String())
	}

}
func TestS3ApiServer_HeadObjectHandler(t *testing.T) {
	bucketName := "/testbucketheado"
	objectName := "/testobjectheado"

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
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := utils.MustNewSignedV4Request(http.MethodHead, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		// Add test case specific headers to the request.
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: head:%v\n", i+1, result.Body.String())
	}

}
func TestS3ApiServer_DeleteObjectHandler(t *testing.T) {
	bucketName := "/testbucketdelo"
	objectName := "/testobjectdelo"

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
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := utils.MustNewSignedV4Request(http.MethodDelete, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		// Add test case specific headers to the request.
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: delete:%v\n", i+1, result.Body.String())
	}

}
func TestS3ApiServer_DeleteMultipleObjectsHandler(t *testing.T) {
	bucketName := "/testbucketdelobjs"

	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())

	delObjReq := datatypes.DeleteObjectsRequest{
		Quiet: false,
	}
	for i := 0; i < 3; i++ {
		data := fmt.Sprintf("1234567%d", i)
		objName := fmt.Sprintf("obj%d", i)
		reqputObject := utils.MustNewSignedV4Request(http.MethodPut, fmt.Sprintf("%s/%s", bucketName, objName), int64(len(data)), bytes.NewReader([]byte(data)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
		// Add test case specific headers to the request.
		reqTest(reqputObject)

		delObjReq.Objects = append(delObjReq.Objects, datatypes.ObjectToDelete{
			ObjectV: datatypes.ObjectV{
				ObjectName: objName,
			},
		})
	}

	// Marshal delete request.
	deleteReqBytes, err := xml.Marshal(delObjReq)
	require.NoError(t, err)
	req := utils.MustNewSignedV4Request(http.MethodPost, bucketName+"?delete=", int64(len(deleteReqBytes)),
		bytes.NewReader(deleteReqBytes), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	result := reqTest(req)
	if result.Code != http.StatusOK {
		t.Fatalf("Expected the response status to be `%d`, but instead found `%d`", http.StatusOK, result.Code)
	}
	deleteResp := response.DeleteObjectsResponse{}
	delRespBytes, err := ioutil.ReadAll(result.Body)
	require.NoError(t, err)
	err = xml.Unmarshal(delRespBytes, &deleteResp)
	require.NoError(t, err)
	t.Log(deleteResp.DeletedObjects)
	delObjMap := make(map[string]datatypes.DeletedObject)
	for _, obj := range delObjReq.Objects {
		delObjMap[obj.ObjectName] = datatypes.DeletedObject{
			ObjectName: obj.ObjectName,
			VersionID:  obj.VersionID,
		}
	}
	for i := 0; i < 3; i++ {
		// All the objects should be under deleted list (including non-existent object)
		obj, ok := delObjMap[deleteResp.DeletedObjects[i].ObjectName]
		require.True(t, ok)
		require.Equal(t, deleteResp.DeletedObjects[i], obj)
	}
	require.Zero(t, len(deleteResp.Errors))
}

func TestS3ApiServer_CopyObjectHandler(t *testing.T) {
	bucketName := "/testbucketcopy"
	objectName := "/testobjectcopy"

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
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := utils.MustNewSignedV4Request(http.MethodPut, testCase.dstbucketName+testCase.dstobjectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		req.Header.Set("X-Amz-Copy-Source", url.QueryEscape(testCase.bucketName+testCase.objectName)) // Add test case specific headers to the request.
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: copy:%v\n", i+1, result.Body.String())
	}

}
func TestS3ApiServer_ListObjectsV2Handler(t *testing.T) {
	bucketName := "/testbucketlist"
	objectName := "/testobjectlist"

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
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		req := utils.MustNewSignedV4Request(http.MethodGet, testCase.bucketName+"?list-type=2", 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
		result := reqTest(req)
		if result.Code != testCase.expectedRespStatus {
			t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
		}
		fmt.Printf("Case %d: list:%v\n", i+1, result.Body.String())
	}

}

func TestWholeNoUserAPI(t *testing.T) {
	bucketName := "/testbucketwhole"
	objectName := "/testobjectwhole"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
	reqPutUser := utils.MustNewSignedV4Request(http.MethodPost, "/admin/v1/add-user"+"?accessKey="+userName+"&secretKey="+secret, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println(reqTest(reqPutUser).Body.String())
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", userName, secret, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Body.String())
	r1 := "1234567"
	reqputObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", userName, secret, t)
	// Add test case specific headers to the request.
	reqTest(reqputObject)
	req := utils.MustNewSignedV4Request(http.MethodGet, bucketName+objectName, 0, nil, "s3", userName, secret, t)
	fmt.Println("getobject", reqTest(req).Body.String())
}
func addCustomHeaders(req *http.Request, customHeaders http.Header) {
	for k, values := range customHeaders {
		for _, value := range values {
			req.Header.Set(k, value)
		}
	}
}
