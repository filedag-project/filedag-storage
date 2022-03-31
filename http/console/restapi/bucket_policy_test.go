package restapi

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_putBucketPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "GG377WW9TPE7G9AHQ0ZR",
		STSSecretAccessKey: "5vXPZAiaKdN5+c4XN4lubOuMnCgKl1cB6Sa+blEN",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJHRzM3N1dXOVRQRTdHOUFIUTBaUiIsImV4cCI6MTY0ODcyNzkyMiwicGFyZW50IjoidGVzdDEifQ.V177tBgnH8KmsUA_0Arc1hkYgTjUkADgMnivvQsTEOxqwqphTN3K2xDwxhj1Vsbc9VW0xfD7NnoYxn8HXutJKA",
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
		STSAccessKeyID:     "14TNGTIC75AF4NTQRRGT",
		STSSecretAccessKey: "tfbPF4OJMaKA3n1DcdNN8GUHUmPKIkiaHxqyRD35",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiIxNFROR1RJQzc1QUY0TlRRUlJHVCIsImV4cCI6MTY0ODcyNzAzMywicGFyZW50IjoidGVzdCJ9.JvU8w7UGlaQB7vtXgojw8hsUk-WtB7rSnLD52l25kg_KCys6-tOLnHy-k9_XAWLG5SKlShU_riuC-Bk6SZjXaw",
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
