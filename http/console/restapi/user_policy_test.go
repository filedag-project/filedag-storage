package restapi

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_putUserPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "UP1GTELGU4EBM81DI9B3",
		STSSecretAccessKey: "gF5YWrT9G6kQyReY7lakz2L8dkhOl1aFIv8bRfj9",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJVUDFHVEVMR1U0RUJNODFESTlCMyIsImV4cCI6MTY0ODgwNDg1MywicGFyZW50IjoidGVzdCJ9.a4Ay0dSzqP76AxoFl10E9MhyR3Vd2wjZfThDw1fSSokMVcx1_KmGN2J4pPnxhpvFs1Fw1zXw5J-jUOHqD5A55w",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	policy := `{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::test22/*"}]}`
	mClient, err := NewMinioAdminClient(session)
	client := AdminClient{Client: mClient}
	err = client.putUserPolicy(context.Background(), "test2", policy)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_getUserPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "UP1GTELGU4EBM81DI9B3",
		STSSecretAccessKey: "gF5YWrT9G6kQyReY7lakz2L8dkhOl1aFIv8bRfj9",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJVUDFHVEVMR1U0RUJNODFESTlCMyIsImV4cCI6MTY0ODgwNDg1MywicGFyZW50IjoidGVzdCJ9.a4Ay0dSzqP76AxoFl10E9MhyR3Vd2wjZfThDw1fSSokMVcx1_KmGN2J4pPnxhpvFs1Fw1zXw5J-jUOHqD5A55w",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	mClient, err := NewMinioAdminClient(session)
	client := AdminClient{Client: mClient}
	err = client.getUserPolicy(context.Background(), "test2")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_removeUserPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "UP1GTELGU4EBM81DI9B3",
		STSSecretAccessKey: "gF5YWrT9G6kQyReY7lakz2L8dkhOl1aFIv8bRfj9",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJVUDFHVEVMR1U0RUJNODFESTlCMyIsImV4cCI6MTY0ODgwNDg1MywicGFyZW50IjoidGVzdCJ9.a4Ay0dSzqP76AxoFl10E9MhyR3Vd2wjZfThDw1fSSokMVcx1_KmGN2J4pPnxhpvFs1Fw1zXw5J-jUOHqD5A55w",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	mClient, err := NewMinioAdminClient(session)
	client := AdminClient{Client: mClient}
	err = client.removeUserPolicy(context.Background(), "test2")
	if err != nil {
		fmt.Println(err)
	}
}
