package restapi

import (
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_getListBucketsResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "CZJLK4KJUG02NY2UHQQA",
		STSSecretAccessKey: "Ql3ppQPiF7eB+Y4at+TUIbWafiuEt0Wst77SWLF1",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJDWkpMSzRLSlVHMDJOWTJVSFFRQSIsImV4cCI6MTY0OTIyODkxMiwicGFyZW50IjoidGVzdCJ9.T-lutD97PK5IsuiNRtKBCiMZYg5wI0o1SjbKSvBmdYItkUptF1x3s91RXFNtZxhRrtbOxGqHtE3lAlVZXxaaoQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	info, err := getListBucketsResponse(session)
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
		STSAccessKeyID:     "CZJLK4KJUG02NY2UHQQA",
		STSSecretAccessKey: "Ql3ppQPiF7eB+Y4at+TUIbWafiuEt0Wst77SWLF1",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJDWkpMSzRLSlVHMDJOWTJVSFFRQSIsImV4cCI6MTY0OTIyODkxMiwicGFyZW50IjoidGVzdCJ9.T-lutD97PK5IsuiNRtKBCiMZYg5wI0o1SjbKSvBmdYItkUptF1x3s91RXFNtZxhRrtbOxGqHtE3lAlVZXxaaoQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	err := getCreateBucketResponse(session, "testN", "", false)
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
	err := getDeleteBucketResponse(session, "testN")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_putBucketPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "CZJLK4KJUG02NY2UHQQA",
		STSSecretAccessKey: "Ql3ppQPiF7eB+Y4at+TUIbWafiuEt0Wst77SWLF1",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJDWkpMSzRLSlVHMDJOWTJVSFFRQSIsImV4cCI6MTY0OTIyODkxMiwicGFyZW50IjoidGVzdCJ9.T-lutD97PK5IsuiNRtKBCiMZYg5wI0o1SjbKSvBmdYItkUptF1x3s91RXFNtZxhRrtbOxGqHtE3lAlVZXxaaoQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	policy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::testN/*"}]}`
	var str = "CUSTOM"
	custom := models.BucketAccess(str)
	request := &models.SetBucketPolicyRequest{
		Access:     &custom,
		Definition: policy,
	}
	bucket, err := getBucketSetPolicyResponse(session, "testN", request)
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

	policy, err := getBucketPolicyResponse(session, "testN")
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
	err := removeBucketPolicyResponse(session, "testN")
	if err != nil {
		fmt.Println(err)
	}
}
