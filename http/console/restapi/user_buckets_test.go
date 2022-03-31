package restapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_getListBucketsResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "14TNGTIC75AF4NTQRRGT",
		STSSecretAccessKey: "tfbPF4OJMaKA3n1DcdNN8GUHUmPKIkiaHxqyRD35",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiIxNFROR1RJQzc1QUY0TlRRUlJHVCIsImV4cCI6MTY0ODcyNzAzMywicGFyZW50IjoidGVzdCJ9.JvU8w7UGlaQB7vtXgojw8hsUk-WtB7rSnLD52l25kg_KCys6-tOLnHy-k9_XAWLG5SKlShU_riuC-Bk6SZjXaw",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	info, err := getListBucketsResponse(session)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}

func Test_getMakeBucketResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "23F021N2VZEUTF9DXRZK",
		STSSecretAccessKey: "OJp+6hFC1Z3XzjRHVjiacmFzWAFk5prBhQpjBJdb",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiIyM0YwMjFOMlZaRVVURjlEWFJaSyIsImV4cCI6MTY0ODcxOTcwNiwicGFyZW50IjoidGVzdCJ9.cQoY1yWHJo4Z77IytMGtHqkVO1-SGV4i47swTiog_WTiZ7KkTUtfcxvYKzp4dcQpd3y6OnCYtDvTqQMVpn6akA",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	err := getMakeBucketResponse(session, "test22", "", false)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_getDeleteBucketResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "4D3OPGC2ZBA1MM07NS9R",
		STSSecretAccessKey: "rXBvmgbEkb9A0v05KYN6uVprPRDbDr+CiVlNmoSK",
		STSSessionToken:    "rXBvmgbEkb9A0v05KYN6uVprPRDbDr+CiVlNmoSK eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiI0RDNPUEdDMlpCQTFNTTA3TlM5UiIsImV4cCI6MTY0ODAwNTg2NywicGFyZW50IjoidGVzdCJ9.A6AzQXStykxB6IGQ4HMgo6lDP2Amet5WDBXDNAU8M6SxWSI7z7DTeIRMtcGW-ciXAXnqya6UOSOcMtbgheUfWQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	err := getDeleteBucketResponse(session, "testName2")
	if err != nil {
		fmt.Println(err)
	}
}
