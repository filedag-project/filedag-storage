package s3api

import (
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"net/http"
	"strings"
	"testing"
)

func TestS3ApiServer_PutGetBucketPolicyHandler(t *testing.T) {
	bucketName := "/testbucketputpoliy"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutBucketResult := reqTest(reqPutBucket)
	if reqPutBucketResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketResult expect 200 ,but found %v", reqPutBucketResult.Code)
	}

	reqPutBucketWrong := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"wrong", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutBucketWrongResult := reqTest(reqPutBucketWrong)
	if reqPutBucketWrongResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketResult expect 200 ,but found %v", reqPutBucketWrongResult.Code)
	}
	reqPutBucketNameDoseNotMatch := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"bucketnotmatch", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutBucketNameDoseNotMatchResult := reqTest(reqPutBucketNameDoseNotMatch)
	if reqPutBucketNameDoseNotMatchResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketNameDoseNotMatchResult expect 200 ,but found %v", reqPutBucketNameDoseNotMatchResult.Code)
	}
	reqPutBucketNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal", 0, nil, "s3", normalUser, normalSecret, t)
	reqPutBucketNormalResult := reqTest(reqPutBucketNormal)
	if reqPutBucketNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketNormalResult expect 200 ,but found %v", reqPutBucketNormalResult.Code)
	}
	correctPolicy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666","filedagadmin"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testbucketputpoliy/*"}]}`
	normalCorrectPolicy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666","testA"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testbucketputpoliynormal/*"}]}`
	accessDeniedPolicy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testbucketputpoliy/*"}]}`
	wrongPolicy := `{"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testbucketputpoliypoliy/*"}]}`
	bucketNameDoseNotMatchPolicy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666","filedagadmin"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testbucketpoliyae/*"}]}`

	teatCases := []struct {
		name                  string
		bucketName            string
		policyJson            string
		accessKey             string
		secretKey             string
		expectedPutRespStatus int // expected response status body.
		expectedGetRespStatus int // expected response status body.
	}{
		{
			name:                  "correctPolicy",
			bucketName:            bucketName,
			policyJson:            correctPolicy,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusNoContent,
			expectedGetRespStatus: http.StatusOK,
		},
		{
			name:                  "accessDeniedPolicy",
			bucketName:            bucketName,
			policyJson:            accessDeniedPolicy,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusNoContent,
			expectedGetRespStatus: http.StatusForbidden,
		},
		{
			name:                  "wrongPolicy",
			bucketName:            bucketName + "wrong",
			policyJson:            wrongPolicy,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusBadRequest,
			expectedGetRespStatus: http.StatusOK,
		},
		{
			name:                  "normal user",
			bucketName:            bucketName + "normal",
			policyJson:            normalCorrectPolicy,
			accessKey:             normalUser,
			secretKey:             normalSecret,
			expectedPutRespStatus: http.StatusNoContent,
			expectedGetRespStatus: http.StatusOK,
		},
		{
			name:                  "bucketName dose not match policy",
			bucketName:            bucketName + "bucketnotmatch",
			policyJson:            bucketNameDoseNotMatchPolicy,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusBadRequest,
			expectedGetRespStatus: http.StatusOK,
		},
		{
			name:                  "nonExistBucket",
			bucketName:            nonExistBucket,
			policyJson:            bucketNameDoseNotMatchPolicy,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusNotFound,
			expectedGetRespStatus: http.StatusNotFound,
		},
	}
	for _, testCase := range teatCases {
		t.Run(testCase.name, func(t *testing.T) {
			reqPutPolicy := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+"?policy", int64(len(testCase.policyJson)), strings.NewReader(testCase.policyJson),
				"s3", testCase.accessKey, testCase.secretKey, t)
			reqPutPolicyResult := reqTest(reqPutPolicy)
			if reqPutPolicyResult.Code != testCase.expectedPutRespStatus {
				t.Fatalf("reqPutPolicyResult : Expected the response status to be `%d`, but instead found `%d`", testCase.expectedPutRespStatus, reqPutPolicyResult.Code)
			}
			reqGetPolicy := utils.MustNewSignedV4Request(http.MethodGet, testCase.bucketName+"?policy", 0, nil, "s3",
				testCase.accessKey, testCase.secretKey, t)
			reqGetPolicyResult := reqTest(reqGetPolicy)
			if reqGetPolicyResult.Code != testCase.expectedGetRespStatus {
				t.Fatalf("reqGetPolicyResult : Expected the response status to be `%d`, but instead found `%d`", testCase.expectedGetRespStatus, reqGetPolicyResult.Code)
			}
		})
	}
}
func TestS3ApiServer_DelBucketPolicyHandler(t *testing.T) {
	bucketName := "/testbucketdelpoliy"
	reqPutBucket := utils.MustNewSignedV4Request(http.MethodPut, bucketName, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
	reqPutBucketResult := reqTest(reqPutBucket)
	if reqPutBucketResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketResult expect 200 ,but found %v", reqPutBucketResult.Code)
	}
	reqPutBucketNormal := utils.MustNewSignedV4Request(http.MethodPut, bucketName+"normal", 0, nil, "s3", normalUser, normalSecret, t)
	reqPutBucketNormalResult := reqTest(reqPutBucketNormal)
	if reqPutBucketNormalResult.Code != http.StatusOK {
		t.Fatalf("reqPutBucketNormalResult expect 200 ,but found %v", reqPutBucketNormalResult.Code)
	}
	correctPolicy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666","filedagadmin"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testbucketdelpoliy/*"}]}`
	normalCorrectPolicy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666","testA"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testbucketdelpoliynormal/*"}]}`
	teatCases := []struct {
		name                  string
		bucketName            string
		policyJson            string
		accessKey             string
		secretKey             string
		expectedPutRespStatus int // expected response status body.
		expectedGetRespStatus int // expected response status body.
	}{
		{
			name:                  "root user del bucket policy",
			bucketName:            bucketName,
			policyJson:            correctPolicy,
			accessKey:             DefaultTestAccessKey,
			secretKey:             DefaultTestSecretKey,
			expectedPutRespStatus: http.StatusNoContent,
			expectedGetRespStatus: http.StatusOK,
		},
		{
			name:                  "normal user del bucket policy",
			bucketName:            bucketName + "normal",
			policyJson:            normalCorrectPolicy,
			accessKey:             normalUser,
			secretKey:             normalSecret,
			expectedPutRespStatus: http.StatusNoContent,
			expectedGetRespStatus: http.StatusOK,
		},
	}
	for _, testCase := range teatCases {
		t.Run(testCase.name, func(t *testing.T) {
			reqPutPolicy := utils.MustNewSignedV4Request(http.MethodPut, testCase.bucketName+"?policy", int64(len(testCase.policyJson)), strings.NewReader(testCase.policyJson),
				"s3", testCase.accessKey, testCase.secretKey, t)
			reqPutPolicyResult := reqTest(reqPutPolicy)
			if reqPutPolicyResult.Code != testCase.expectedPutRespStatus {
				t.Fatalf("reqPutPolicyResult : Expected the response status to be `%d`, but instead found `%d`", testCase.expectedPutRespStatus, reqPutPolicyResult.Code)
			}
			reqDelPolicy := utils.MustNewSignedV4Request(http.MethodDelete, testCase.bucketName+"?policy", 0, nil, "s3",
				testCase.accessKey, testCase.secretKey, t)
			reqDelPolicyResult := reqTest(reqDelPolicy)
			if reqDelPolicyResult.Code != testCase.expectedGetRespStatus {
				t.Fatalf("reqDelPolicyResult : Expected the response status to be `%d`, but instead found `%d`", testCase.expectedGetRespStatus, reqDelPolicyResult.Code)
			}
		})
	}
}
