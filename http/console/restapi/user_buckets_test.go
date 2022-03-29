package restapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_getListBucketsResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "XTPCOTBMXP890HLPNG3J",
		STSSecretAccessKey: "pkYPrtqeIUxfctauTduI1Wy8jSMEVvfvoe7MgBbV",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJYVFBDT1RCTVhQODkwSExQTkczSiIsImV4cCI6MTY0ODQ2ODA4NywicGFyZW50IjoidGVzdDEifQ.vIx5NS4xeNhzH_4HTfR-IZUvJIEpjGxT07aWXKVjb3BOdXUlmxuUR7cgvEPUm0YwoodggAcVugw5pZnJbXAu-w",
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
		STSAccessKeyID:     "XTPCOTBMXP890HLPNG3J",
		STSSecretAccessKey: "pkYPrtqeIUxfctauTduI1Wy8jSMEVvfvoe7MgBbV",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJYVFBDT1RCTVhQODkwSExQTkczSiIsImV4cCI6MTY0ODQ2ODA4NywicGFyZW50IjoidGVzdDEifQ.vIx5NS4xeNhzH_4HTfR-IZUvJIEpjGxT07aWXKVjb3BOdXUlmxuUR7cgvEPUm0YwoodggAcVugw5pZnJbXAu-w",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	err := getMakeBucketResponse(session, "testName", "", false)
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
