package s3api

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/datatypes"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
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
	folder := "/testfolder/"
	r1 := "1234567"
	copySourceHeader := http.Header{}
	copySourceHeader.Set("X-Amz-Copy-Source", "somewhere")
	invalidMD5Header := http.Header{}
	invalidMD5Header.Set("Content-Md5", "42")
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	result := reqTest(reqPutBucket)
	if result.Code != http.StatusOK {
		t.Fatalf("the response status of putbucket: %d", result.Code)
	}
	reqPutBucketNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal", 0, nil, "s3", normalUser, normalSecret, t)
	resultNormal := reqTest(reqPutBucketNormal)
	if result.Code != http.StatusOK {
		t.Fatalf("the response status of reqPutBucketNormal: %d", resultNormal.Code)
	}

	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name       string
		bucketName string
		objectName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1
		{
			name:               "root user put obj",
			bucketName:         bucketName,
			objectName:         objectName,
			data:               []byte(r1),
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "normal user put obj",
			bucketName:         bucketName + "normal",
			objectName:         objectName,
			data:               []byte(r1),
			accessKey:          normalUser,
			secretKey:          normalSecret,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "wrong accessKey",
			bucketName:         bucketName,
			objectName:         objectName,
			data:               []byte(r1),
			accessKey:          wrongAccessKey,
			secretKey:          wrongSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
		{
			name:               "non-exist bucket",
			bucketName:         nonExistBucket,
			objectName:         objectName,
			data:               []byte(r1),
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
		{
			name:               "add copySourceHeader",
			bucketName:         bucketName,
			objectName:         objectName,
			data:               []byte(r1),
			header:             copySourceHeader,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusBadRequest,
		},
		{
			name:               "invalidMD5Header",
			bucketName:         bucketName,
			objectName:         objectName,
			data:               []byte(r1),
			header:             invalidMD5Header,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusBadRequest,
		},
		{
			name:               "fileFolder",
			bucketName:         bucketName,
			objectName:         folder,
			data:               nil,
			header:             nil,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "fileFolder-bad-StatusMovedPermanently",
			bucketName:         bucketName,
			objectName:         objectName + "///",
			data:               nil,
			header:             nil,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusMovedPermanently,
		},
		{
			name:               "fileFolder-bad-StatusBadRequest",
			bucketName:         bucketName,
			objectName:         objectName + "/" + "image/",
			data:               nil,
			header:             nil,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusBadRequest,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+testCase.objectName, int64(len(testCase.data)), bytes.NewReader(testCase.data), "s3", testCase.accessKey, testCase.secretKey, t)
			// Add test case specific headers to the request.
			addCustomHeaders(req, testCase.header)
			result := reqTest(req)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
			}
		})
	}

}
func TestS3ApiServer_GetObjectHandler(t *testing.T) {
	bucketName := "/testbucketfeto"
	objectName := "/testobjectgeto"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutBucketResult := reqTest(reqPutBucket)
	if reqPutBucketResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketResult.Code)
	}
	r1 := "1234567"
	reqPutObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutObjectResult := reqTest(reqPutObject)
	if reqPutObjectResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectResult.Code)
	}

	//normal user bucket and object
	reqPutBucketNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal", 0, nil, "s3", normalUser, normalSecret, t)
	reqPutBucketNormalResult := reqTest(reqPutBucketNormal)
	if reqPutBucketNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketNormalResult.Code)
	}
	reqPutObjectNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal"+objectName+"normal", int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", normalUser, normalSecret, t)
	reqPutObjectNormalResult := reqTest(reqPutObjectNormal)
	if reqPutObjectNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectNormalResult.Code)
	}

	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name       string
		bucketName string
		objectName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1
		{
			name:               "root user get object",
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "normal user get object",
			bucketName:         bucketName + "normal",
			objectName:         objectName + "normal",
			accessKey:          normalUser,
			secretKey:          normalSecret,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			name:               "wrong accessKey.",
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          wrongAccessKey,
			secretKey:          wrongSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
		{
			name:               "non-exist bucket",
			bucketName:         nonExistBucket,
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
		{
			name:               "non-exist object",
			bucketName:         bucketName,
			objectName:         nonExistObject,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := utils.MustNewSignedV4Request(http.MethodGet, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			// Add test case specific headers to the request.
			result := reqTest(req)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})
	}

}
func TestS3ApiServer_HeadObjectHandler(t *testing.T) {
	bucketName := "/testbucketheado"
	objectName := "/testobjectheado"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutBucketResult := reqTest(reqPutBucket)
	if reqPutBucketResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketResult.Code)
	}
	r1 := "1234567"
	reqPutObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutObjectResult := reqTest(reqPutObject)
	if reqPutObjectResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectResult.Code)
	}

	//normal user bucket and object
	reqPutBucketNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal", 0, nil, "s3", normalUser, normalSecret, t)
	reqPutBucketNormalResult := reqTest(reqPutBucketNormal)
	if reqPutBucketNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketNormalResult.Code)
	}
	reqPutObjectNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal"+objectName+"normal", int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", normalUser, normalSecret, t)
	reqPutObjectNormalResult := reqTest(reqPutObjectNormal)
	if reqPutObjectNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectNormalResult.Code)
	}
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name       string
		bucketName string
		objectName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1
		{
			name:               "root user head object",
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		{
			name:               "normal user head object",
			bucketName:         bucketName + "normal",
			objectName:         objectName + "normal",
			accessKey:          normalUser,
			secretKey:          normalSecret,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "wrong accessKey.",
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          wrongAccessKey,
			secretKey:          wrongSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
		{
			name:               "non-exist bucket",
			bucketName:         nonExistBucket,
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
		{
			name:               "non-exist object",
			bucketName:         bucketName,
			objectName:         nonExistObject,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := utils.MustNewSignedV4Request(http.MethodHead, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			// Add test case specific headers to the request.
			result := reqTest(req)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})
	}

}
func TestS3ApiServer_DeleteObjectHandler(t *testing.T) {
	bucketName := "/testbucketdelo"
	objectName := "/testobjectdelo"
	folderName := "/testfolder/"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutBucketResult := reqTest(reqPutBucket)
	if reqPutBucketResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketResult.Code)
	}
	r1 := "1234567"
	reqPutObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutObjectResult := reqTest(reqPutObject)
	if reqPutObjectResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectResult.Code)
	}

	//normal user bucket and object
	reqPutBucketNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal", 0, nil, "s3", normalUser, normalSecret, t)
	reqPutBucketNormalResult := reqTest(reqPutBucketNormal)
	if reqPutBucketNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketNormalResult.Code)
	}
	reqPutObjectNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal"+objectName+"normal", int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", normalUser, normalSecret, t)
	reqPutObjectNormalResult := reqTest(reqPutObjectNormal)
	if reqPutObjectNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectNormalResult.Code)
	}
	reqPutFolder := utils.MustNewSignedV4Request(http.MethodPut, bucketName+folderName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutFolderResult := reqTest(reqPutFolder)
	if reqPutFolderResult.Code != http.StatusOK {
		t.Fatalf("reqPutFolderResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutFolderResult.Code)
	}
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name       string
		bucketName string
		objectName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1
		{
			name:               "root user del object",
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNoContent,
		},
		{
			name:               "normal user del object",
			bucketName:         bucketName + "normal",
			objectName:         objectName + "normal",
			accessKey:          normalUser,
			secretKey:          normalSecret,
			expectedRespStatus: http.StatusNoContent,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			name:               "wrong accessKey.",
			bucketName:         bucketName,
			objectName:         objectName,
			accessKey:          wrongAccessKey,
			secretKey:          wrongSecretKey,
			expectedRespStatus: http.StatusForbidden,
		},
		{
			name:               "non-exist bucket",
			bucketName:         nonExistBucket,
			objectName:         objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
		{
			name:               "non-exist object",
			bucketName:         bucketName,
			objectName:         nonExistObject,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
		{
			name:               "folder",
			bucketName:         bucketName,
			objectName:         folderName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNoContent,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := utils.MustNewSignedV4Request(http.MethodDelete, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			// Add test case specific headers to the request.
			result := reqTest(req)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})

	}

}
func TestS3ApiServer_DeleteMultipleObjectsHandler(t *testing.T) {
	bucketName := "testbucketdelobjs"

	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, "/"+bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	fmt.Println("putbucket:", reqTest(reqPutBucket).Code)

	delObjReq := datatypes.DeleteObjectsRequest{
		Quiet: false,
	}
	for i := 0; i < 3; i++ {
		data := fmt.Sprintf("1234567%d", i)
		objName := fmt.Sprintf("obj%d", i)
		reqputObject := utils.MustNewSignedV4Request(http.MethodPut, fmt.Sprintf("/%s/%s", bucketName, objName), int64(len(data)), bytes.NewReader([]byte(data)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
		// Add test case specific headers to the request.
		result := reqTest(reqputObject)
		if result.Code != http.StatusOK {
			t.Fatalf("Expected the response status to be `%d`, but instead found `%d`", http.StatusOK, result.Code)
		}

		delObjReq.Objects = append(delObjReq.Objects, datatypes.ObjectToDelete{
			ObjectV: datatypes.ObjectV{
				ObjectName: objName,
			},
		})
	}

	// Marshal delete request.
	deleteReqBytes, err := xml.Marshal(delObjReq)
	require.NoError(t, err)
	req := utils.MustNewSignedV4Request(http.MethodPost, "/"+bucketName+"?delete=", int64(len(deleteReqBytes)),
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

	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutBucketResult := reqTest(reqPutBucket)
	if reqPutBucketResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketResult.Code)
	}
	r1 := "1234567"
	reqPutObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutObjectResult := reqTest(reqPutObject)
	if reqPutObjectResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectResult.Code)
	}

	//normal user bucket and object
	reqPutBucketNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal", 0, nil, "s3", normalUser, normalSecret, t)
	reqPutBucketNormalResult := reqTest(reqPutBucketNormal)
	if reqPutBucketNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketNormalResult.Code)
	}
	reqPutObjectNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal"+objectName+"normal", int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", normalUser, normalSecret, t)
	reqPutObjectNormalResult := reqTest(reqPutObjectNormal)
	if reqPutObjectNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectNormalResult.Code)
	}

	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name          string
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
		// Test case - 1
		{
			name:               "root user copy object",
			bucketName:         bucketName,
			objectName:         objectName,
			dstbucketName:      bucketName,
			dstobjectName:      "1.txt",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "normal user copy object",
			bucketName:         bucketName + "normal",
			objectName:         objectName + "normal",
			dstbucketName:      bucketName + "normal",
			dstobjectName:      "1.txt",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		// Test case - 2.
		// wrong accessKey.
		{
			name:               "wrong accessKey.",
			bucketName:         bucketName,
			objectName:         objectName,
			dstbucketName:      bucketName,
			dstobjectName:      "1.txt",
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
		{
			name:               "source non-exist bucket",
			bucketName:         nonExistBucket,
			objectName:         objectName,
			dstbucketName:      bucketName,
			dstobjectName:      "1.txt",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusConflict,
		},
		{
			name:               "source non-exist object",
			bucketName:         bucketName,
			objectName:         nonExistObject,
			dstbucketName:      bucketName,
			dstobjectName:      "1.txt",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusConflict,
		},
		{
			name:               "dst non-exist bucket",
			bucketName:         nonExistBucket,
			objectName:         objectName,
			dstbucketName:      bucketName,
			dstobjectName:      "1.txt",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusConflict,
		},
		{
			name:               "dst non-exist object",
			bucketName:         bucketName,
			objectName:         nonExistObject,
			dstbucketName:      bucketName,
			dstobjectName:      "1.txt",
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusConflict,
		},
		{
			name:               "copy to same place",
			bucketName:         bucketName,
			objectName:         objectName,
			dstbucketName:      bucketName,
			dstobjectName:      objectName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusBadRequest,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := utils.MustNewSignedV4Request(http.MethodPut, testCase.dstbucketName+testCase.dstobjectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			req.Header.Set("X-Amz-Copy-Source", url.QueryEscape(testCase.bucketName+testCase.objectName)) // Add test case specific headers to the request.
			result := reqTest(req)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
			}
		})
	}

}
func TestS3ApiServer_ListObjectsV2Handler(t *testing.T) {
	bucketName := "/testbucketlistv2"
	objectName := "/testobjectlist"
	folderName := "/testfolder/"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutBucketResult := reqTest(reqPutBucket)
	if reqPutBucketResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketResult.Code)
	}
	r1 := "1234567"
	reqPutObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutObjectResult := reqTest(reqPutObject)
	if reqPutObjectResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectResult.Code)
	}

	//normal user bucket and object
	reqPutBucketNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal", 0, nil, "s3", normalUser, normalSecret, t)
	reqPutBucketNormalResult := reqTest(reqPutBucketNormal)
	if reqPutBucketNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutBucketNormalResult.Code)
	}
	reqPutObjectNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal"+objectName+"normal", int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", normalUser, normalSecret, t)
	reqPutObjectNormalResult := reqTest(reqPutObjectNormal)
	if reqPutObjectNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectNormalResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutObjectNormalResult.Code)
	}
	reqPutFolder := utils.MustNewSignedV4Request(http.MethodPut, bucketName+folderName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutFolderResult := reqTest(reqPutFolder)
	if reqPutFolderResult.Code != http.StatusOK {
		t.Fatalf("reqPutFolderResult : Expected the response status to be `%d`, but instead found `%d`", 200, reqPutFolderResult.Code)
	}
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name       string
		bucketName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectedRespStatus int // expected response status body.
	}{
		// Test case - 1
		{
			name:               "root user list objects",
			bucketName:         bucketName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
		{
			name:               "normal user list objects",
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
			accessKey:          "1",
			secretKey:          "1",
			expectedRespStatus: http.StatusForbidden,
		},
		{
			name:               "non-exist bucket",
			bucketName:         nonExistBucket,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusNotFound,
		},
		{
			name:               "list folder",
			bucketName:         bucketName,
			accessKey:          DefaultTestAccessKey,
			secretKey:          DefaultTestSecretKey,
			expectedRespStatus: http.StatusOK,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := utils.MustNewSignedV4Request(http.MethodGet, testCase.bucketName+"?list-type=2"+"&&prefix="+folderName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			result := reqTest(req)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespStatus, result.Code)
			}
		})
	}

}

func TestWholeProcess(t *testing.T) {
	userName := "dean"
	secret := "dean123456"
	capacity := "9999999"
	bucketName := "/testbucket"
	objectName := "/testobject"
	addUrl := "http://127.0.0.1:9985/admin/v1/add-user?"

	urlValues := make(url.Values)
	urlValues.Set("accessKey", userName)
	urlValues.Set("secretKey", secret)
	urlValues.Set("capacity", capacity)
	reqPutUser := utils.MustNewSignedV4Request(http.MethodPost, addUrl+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutUserResult := reqTest(reqPutUser)
	if reqPutUserResult.Code != http.StatusOK {
		t.Fatalf("reqPutUserResult expect 200 but found %v", reqPutUserResult.Code)
	}
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", userName, secret, t)
	reqPutBucketResult := reqTest(reqPutBucket)
	if reqPutBucketResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketResult expect 200 but found %v", reqPutBucketResult.Code)
	}
	r1 := "1234567"
	reqPutObject := utils.MustNewSignedV4Request(http.MethodPut, bucketName+objectName, int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", userName, secret, t)
	// Add test case specific headers to the request.
	reqPutObjectResult := reqTest(reqPutObject)
	if reqPutObjectResult.Code != http.StatusOK {
		t.Fatalf("reqPutObjectResult expect 200 but found %v", reqPutObjectResult.Code)
	}
	reqGetObj := utils.MustNewSignedV4Request(http.MethodGet, bucketName+objectName, 0, nil, "s3", userName, secret, t)
	reqGetObjResult := reqTest(reqGetObj)
	if reqGetObjResult.Code != http.StatusOK {
		t.Fatalf("reqGetObjResult expect 200 but found %v", reqGetObjResult.Code)
	}
}
func addCustomHeaders(req *http.Request, customHeaders http.Header) {
	for k, values := range customHeaders {
		for _, value := range values {
			req.Header.Set(k, value)
		}
	}
}
func TestS3ApiServer_PutObjectFolderHandler(t *testing.T) {
	bucketName := "/testbucketputf"
	keyName := "/testobjectputf"
	r1 := "1234567"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultPutBucket := reqTest(reqPutBucket)
	if resultPutBucket.Code != http.StatusOK {
		t.Fatalf("the response status of putbucket: %d", resultPutBucket.Code)
	}
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name       string
		bucketName string
		objectName string
		data       []byte
		header     http.Header
		accessKey  string
		secretKey  string
		// expected output.
		expectBody            string
		expectedRespStatus    int // expected response status body.
		expectedRespGetStatus int // expected response status body.
	}{
		{
			name:                  "create folder",
			bucketName:            bucketName,
			objectName:            keyName + "/aaa/bbb/ccc/",
			data:                  nil,
			header:                nil,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectBody:            "",
			expectedRespStatus:    http.StatusOK,
			expectedRespGetStatus: http.StatusOK,
		},
		{
			name:                  "create folder,but file in body",
			bucketName:            bucketName,
			objectName:            keyName + "/aaa/bbb/ccc/",
			data:                  []byte(r1),
			header:                nil,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectBody:            r1,
			expectedRespStatus:    http.StatusOK,
			expectedRespGetStatus: http.StatusOK,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+testCase.objectName, int64(len(testCase.data)), bytes.NewReader(testCase.data), "s3", testCase.accessKey, testCase.secretKey, t)
			// Add test case specific headers to the request.
			addCustomHeaders(req, testCase.header)
			result := reqTest(req)
			if result.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result.Code)
			}

			reqGet := utils.MustNewSignedV4Request(http.MethodGet, testCase.bucketName+testCase.objectName, 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			// Add test case specific headers to the request.
			resultGet := reqTest(reqGet)
			fmt.Println(resultGet.Body.String())
			if resultGet.Body.String() != testCase.expectBody {
				t.Fatalf("Case %s: Expected the response status to be `%s`, but instead found `%s`", testCase.name, testCase.expectBody, resultGet.Body.String())
			}
			if resultGet.Code != testCase.expectedRespGetStatus {
				t.Fatalf("Case %s: Expected the response status to be `%d`, but instead found `%d`", testCase.name, testCase.expectedRespGetStatus, resultGet.Code)
			}
		})
	}

}

func TestS3ApiServer_PutObjectAsyncHandler(t *testing.T) {
	bucketName := "/testbucketputfasync"
	keyName := "/testobjectputf"
	r1 := "1234567"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultPutBucket := reqTest(reqPutBucket)
	if resultPutBucket.Code != http.StatusOK {
		t.Fatalf("the response status of putbucket: %d", resultPutBucket.Code)
	}
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name        string
		bucketName  string
		objectName1 string
		objectName2 string
		data1       []byte
		data2       []byte
		header      http.Header
		accessKey   string
		secretKey   string
		// expected output.
		expectBody            string
		expectedRespStatus    int // expected response status body.
		expectedRespGetStatus int // expected response status body.
	}{
		{
			name:                  "Async create folder",
			bucketName:            bucketName,
			objectName1:           keyName + "/aaa/bbb/ccc/a.jpg",
			objectName2:           keyName + "/aaa/bbb/ccc/",
			data1:                 []byte(r1),
			data2:                 nil,
			header:                nil,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectBody:            "",
			expectedRespStatus:    http.StatusOK,
			expectedRespGetStatus: http.StatusOK,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			go func() {
				req1 := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+testCase.objectName1, int64(len(testCase.data1)), bytes.NewReader(testCase.data1), "s3", testCase.accessKey, testCase.secretKey, t)
				// Add test case specific headers to the request.
				addCustomHeaders(req1, testCase.header)
				result1 := reqTest(req1)
				if result1.Code != testCase.expectedRespStatus {
					t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result1.Code)
				}
			}()
			req2 := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+testCase.objectName2, int64(len(testCase.data2)), bytes.NewReader(testCase.data2), "s3", testCase.accessKey, testCase.secretKey, t)
			// Add test case specific headers to the request.
			addCustomHeaders(req2, testCase.header)
			result2 := reqTest(req2)
			if result2.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result2.Code)
			}
		})
	}
}

func TestS3ApiServer_PutObjectInsteadFolderHandler(t *testing.T) {
	bucketName := "/testbucketputfinstead"
	keyName := "/testobjectputf"
	r1 := "1234567"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	resultPutBucket := reqTest(reqPutBucket)
	if resultPutBucket.Code != http.StatusOK {
		t.Fatalf("the response status of putbucket: %d", resultPutBucket.Code)
	}
	reqPut := utils.MustNewSignedV4Request(http.MethodPut, bucketName+keyName+"/aaa/bbb/ccc/1.jpg", int64(len(r1)), bytes.NewReader([]byte(r1)), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	// Add test case specific headers to the request.
	resultPut := reqTest(reqPut)
	if resultPut.Code != http.StatusOK {
		t.Fatalf(" Expected the response status to be `%d`, but instead found `%d`", http.StatusOK, resultPut.Code)
	}
	// test cases with inputs and expected result for Bucket.
	testCases := []struct {
		name        string
		bucketName  string
		objectName1 string
		objectName2 string
		data1       []byte
		data2       []byte
		header      http.Header
		accessKey   string
		secretKey   string
		// expected output.
		expectBody            string
		expectedRespStatus    int // expected response status body.
		expectedRespGetStatus int // expected response status body.
	}{
		{
			name:                  "Async create folder",
			bucketName:            bucketName,
			objectName1:           keyName + "/aaa/bbb/ccc/",
			objectName2:           keyName + "/aaa/bbb/ccc/",
			data1:                 nil,
			data2:                 []byte(r1),
			header:                nil,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectBody:            "",
			expectedRespStatus:    http.StatusOK,
			expectedRespGetStatus: http.StatusOK,
		},
	}
	// Iterating over the cases, fetching the object validating the response.
	for i, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			req1 := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+testCase.objectName1, int64(len(testCase.data1)), bytes.NewReader(testCase.data1), "s3", testCase.accessKey, testCase.secretKey, t)
			// Add test case specific headers to the request.
			addCustomHeaders(req1, testCase.header)
			result1 := reqTest(req1)
			if result1.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result1.Code)
			}

			req2 := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+testCase.objectName1, int64(len(testCase.data2)), bytes.NewReader(testCase.data2), "s3", testCase.accessKey, testCase.secretKey, t)
			// Add test case specific headers to the request.
			addCustomHeaders(req2, testCase.header)
			result2 := reqTest(req2)
			if result2.Code != testCase.expectedRespStatus {
				t.Fatalf("Case %d: Expected the response status to be `%d`, but instead found `%d`", i+1, testCase.expectedRespStatus, result2.Code)
			}
			req := utils.MustNewSignedV4Request(http.MethodGet, testCase.bucketName+"?prefix="+keyName+"/aaa/bbb/ccc/"+"&&delimiter=/", 0, nil, "s3", testCase.accessKey, testCase.secretKey, t)
			result := reqTest(req)
			fmt.Println(result.Body.String())
		})
	}
}
