package restapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_getListBucketsResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "VXZN7MG2GATGH3AT6VQM",
		STSSecretAccessKey: "xJBjqPNQIUE227WlMg+YEFk3lkvbsJ+VdanPkAt5",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJWWFpON01HMkdBVEdIM0FUNlZRTSIsImV4cCI6MTY0ODAzMTMzNywicGFyZW50IjoidGVzdCJ9.ulBQEnLo6KY6DnxSguSIKhekWCSmGiJoApJEjw7Dp5nHTRnvLBq3BcaNIEuYpEJjbnVSKF5tBgIMayGMsX3_PA",
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
		STSAccessKeyID:     "VXZN7MG2GATGH3AT6VQM",
		STSSecretAccessKey: "xJBjqPNQIUE227WlMg+YEFk3lkvbsJ+VdanPkAt5",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJWWFpON01HMkdBVEdIM0FUNlZRTSIsImV4cCI6MTY0ODAzMTMzNywicGFyZW50IjoidGVzdCJ9.ulBQEnLo6KY6DnxSguSIKhekWCSmGiJoApJEjw7Dp5nHTRnvLBq3BcaNIEuYpEJjbnVSKF5tBgIMayGMsX3_PA",
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
