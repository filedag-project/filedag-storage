package restapi

import (
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_getListBucketsResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "KIBUMQ2R8LWCC5USEMFH",
		STSSecretAccessKey: "ryQOZEdygKk4dhQ9b8uGR6loQBHHIRbRPu9NXoeN",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJLSUJVTVEyUjhMV0NDNVVTRU1GSCIsImV4cCI6MTY1MDM0MDQ0OCwicGFyZW50IjoidGVzdCJ9.f2Pc-PUQrzx8zqeXHxvG2FXZZVGrR3uMWQCf8dSSXiSSCF_IqszaycvacKrCC1QZO-DhNB9JgK3rDlEUupHRHg",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	apiServer := ApiServer{}
	info, err := apiServer.GetListBucketsResponse(session)
	if err != nil {
		fmt.Println(err)
	}
	if info != nil {
		bytes, _ := json.Marshal(info)
		fmt.Println("listBuckets", string(bytes))
	}
}

func Test_getCreateBucketResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "KIBUMQ2R8LWCC5USEMFH",
		STSSecretAccessKey: "ryQOZEdygKk4dhQ9b8uGR6loQBHHIRbRPu9NXoeN",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJLSUJVTVEyUjhMV0NDNVVTRU1GSCIsImV4cCI6MTY1MDM0MDQ0OCwicGFyZW50IjoidGVzdCJ9.f2Pc-PUQrzx8zqeXHxvG2FXZZVGrR3uMWQCf8dSSXiSSCF_IqszaycvacKrCC1QZO-DhNB9JgK3rDlEUupHRHg",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	apiServer := ApiServer{}
	err := apiServer.GetCreateBucketResponse(session, "testN", "", false)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_getDeleteBucketResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "CZJLK4KJUG02NY2UHQQA",
		STSSecretAccessKey: "Ql3ppQPiF7eB+Y4at+TUIbWafiuEt0Wst77SWLF1",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJDWkpMSzRLSlVHMDJOWTJVSFFRQSIsImV4cCI6MTY0OTIyODkxMiwicGFyZW50IjoidGVzdCJ9.T-lutD97PK5IsuiNRtKBCiMZYg5wI0o1SjbKSvBmdYItkUptF1x3s91RXFNtZxhRrtbOxGqHtE3lAlVZXxaaoQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	apiServer := ApiServer{}
	err := apiServer.GetDeleteBucketResponse(session, "testN")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_putBucketPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "KIBUMQ2R8LWCC5USEMFH",
		STSSecretAccessKey: "ryQOZEdygKk4dhQ9b8uGR6loQBHHIRbRPu9NXoeN",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJLSUJVTVEyUjhMV0NDNVVTRU1GSCIsImV4cCI6MTY1MDM0MDQ0OCwicGFyZW50IjoidGVzdCJ9.f2Pc-PUQrzx8zqeXHxvG2FXZZVGrR3uMWQCf8dSSXiSSCF_IqszaycvacKrCC1QZO-DhNB9JgK3rDlEUupHRHg",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	policy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testN/*"}]}`
	var str = "CUSTOM"
	custom := models.BucketAccess(str)
	request := &models.SetBucketPolicyParams{
		BucketName: "testN",
		Access:     &custom,
		Definition: policy,
	}
	apiServer := ApiServer{}
	bucket, err := apiServer.GetBucketSetPolicyResponse(session, request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(bucket)
}

func Test_getBucketPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "CZJLK4KJUG02NY2UHQQA",
		STSSecretAccessKey: "Ql3ppQPiF7eB+Y4at+TUIbWafiuEt0Wst77SWLF1",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJDWkpMSzRLSlVHMDJOWTJVSFFRQSIsImV4cCI6MTY0OTIyODkxMiwicGFyZW50IjoidGVzdCJ9.T-lutD97PK5IsuiNRtKBCiMZYg5wI0o1SjbKSvBmdYItkUptF1x3s91RXFNtZxhRrtbOxGqHtE3lAlVZXxaaoQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	apiServer := ApiServer{}
	policy, err := apiServer.GetBucketPolicyResponse(session, "testN")
	if err != nil {
		fmt.Println(err)
	}
	if policy != nil {
		bytes, _ := json.Marshal(policy)
		fmt.Println("policy:", string(bytes))
	}
}

func Test_removeBucketPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "CZJLK4KJUG02NY2UHQQA",
		STSSecretAccessKey: "Ql3ppQPiF7eB+Y4at+TUIbWafiuEt0Wst77SWLF1",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJDWkpMSzRLSlVHMDJOWTJVSFFRQSIsImV4cCI6MTY0OTIyODkxMiwicGFyZW50IjoidGVzdCJ9.T-lutD97PK5IsuiNRtKBCiMZYg5wI0o1SjbKSvBmdYItkUptF1x3s91RXFNtZxhRrtbOxGqHtE3lAlVZXxaaoQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	apiServer := ApiServer{}
	err := apiServer.RemoveBucketPolicyResponse(session, "testN")
	if err != nil {
		fmt.Println(err)
	}
}
