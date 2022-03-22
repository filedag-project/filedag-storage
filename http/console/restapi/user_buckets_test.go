package restapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_getListBucketsResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "NZZYDB8776AU6XTQ6C8T",
		STSSecretAccessKey: "pkd6Amn2j1duguJPpyiUSwK0+p15hvU1d5wagfTI",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJOWlpZREI4Nzc2QVU2WFRRNkM4VCIsImV4cCI6MTY0Nzg1OTA3OCwicGFyZW50IjoidGVzdCJ9.abOZuotj_5vF2Dg5EnC_nZ5YI_BtQxfukFPTntnd0A6meke68Dn_ByErFBpQuOYtJp_vYyi47TWlXLcSUKG6bQ",
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
		STSAccessKeyID:     "NZZYDB8776AU6XTQ6C8T",
		STSSecretAccessKey: "pkd6Amn2j1duguJPpyiUSwK0+p15hvU1d5wagfTI",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJOWlpZREI4Nzc2QVU2WFRRNkM4VCIsImV4cCI6MTY0Nzg1OTA3OCwicGFyZW50IjoidGVzdCJ9.abOZuotj_5vF2Dg5EnC_nZ5YI_BtQxfukFPTntnd0A6meke68Dn_ByErFBpQuOYtJp_vYyi47TWlXLcSUKG6bQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	err := getMakeBucketResponse(session, "testName", "", false)
	if err != nil {
		fmt.Println(err)
	}
}
