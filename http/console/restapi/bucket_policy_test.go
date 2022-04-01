package restapi

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_putBucketPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "53BEEYC7TR1YSQ9SUYLW",
		STSSecretAccessKey: "4DZJ98Txgy1TWGTdRHD1YZjFPbqelCoezzvOL3ha",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiI1M0JFRVlDN1RSMVlTUTlTVVlMVyIsImV4cCI6MTY0ODc4MjcwNiwicGFyZW50IjoidGVzdDEifQ._ULv3jwmU8sHlm30lMP_XLrjdcDw57NsOlc3jkSGC7p42IfdeRP4mnnbdZNrMdla-RXcdn1kqFjDvo3Ts9sZeg",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	policy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::test22/*"}]}`
	mClient, err := NewMinioAdminClient(session)
	client := AdminClient{Client: mClient}
	err = client.putBucketPolicy(context.Background(), "test22", policy)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_getBucketPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "53BEEYC7TR1YSQ9SUYLW",
		STSSecretAccessKey: "4DZJ98Txgy1TWGTdRHD1YZjFPbqelCoezzvOL3ha",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiI1M0JFRVlDN1RSMVlTUTlTVVlMVyIsImV4cCI6MTY0ODc4MjcwNiwicGFyZW50IjoidGVzdDEifQ._ULv3jwmU8sHlm30lMP_XLrjdcDw57NsOlc3jkSGC7p42IfdeRP4mnnbdZNrMdla-RXcdn1kqFjDvo3Ts9sZeg",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	mClient, err := NewMinioAdminClient(session)
	client := AdminClient{Client: mClient}
	err = client.getBucketPolicy(context.Background(), "test22")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_removeBucketPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "53BEEYC7TR1YSQ9SUYLW",
		STSSecretAccessKey: "4DZJ98Txgy1TWGTdRHD1YZjFPbqelCoezzvOL3ha",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiI1M0JFRVlDN1RSMVlTUTlTVVlMVyIsImV4cCI6MTY0ODc4MjcwNiwicGFyZW50IjoidGVzdDEifQ._ULv3jwmU8sHlm30lMP_XLrjdcDw57NsOlc3jkSGC7p42IfdeRP4mnnbdZNrMdla-RXcdn1kqFjDvo3Ts9sZeg",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	mClient, err := NewMinioAdminClient(session)
	client := AdminClient{Client: mClient}
	err = client.removeBucketPolicy(context.Background(), "test22")
	if err != nil {
		fmt.Println(err)
	}
}
