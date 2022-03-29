package restapi

import (
	"bytes"
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"io/ioutil"
	"testing"
)

func Test_getUploadObjectResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "RFHEVBU72KNPLRNIYR6C",
		STSSecretAccessKey: "VXyVlRZkIqR2Lmyv2xUwKaLgg3ONMvHchlHaXb0c",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJSRkhFVkJVNzJLTlBMUk5JWVI2QyIsImV4cCI6MTY0ODUyMDM0OSwicGFyZW50IjoidGVzdDEifQ.6_Q7oQ_YNlufwTFI-aTVGYQudbKa_Inp6IwxB_OuoPRFDFyfNa_tYF8DdBSLgxTtXlY5ub5Aehy8FTGvIkt8Fw",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	mClient, err := NewMinioAdminClient(session)
	client := AdminClient{Client: mClient}
	r1, _ := ioutil.ReadFile("user_objects.go")
	err = client.putObject(context.Background(), "testName", "name.go", bytes.NewReader(r1), int64(len(r1)))
	if err != nil {
		fmt.Println(err)
	}
}

func Test_getDownloadObjectResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "DP36WI5WQVCRG41QFG09",
		STSSecretAccessKey: "7avgd5HfvUpq6ZMhXZBtwA+NWD9qi1BS7lUl2c3h",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJEUDM2V0k1V1FWQ1JHNDFRRkcwOSIsImV4cCI6MTY0ODU0Nzg0OSwicGFyZW50IjoidGVzdDEifQ.r7oF_OSRHnHFXO8H1WZ_pqG7I7feQFFEkxXzCXFIXwA6JagSNY7h_h391ZZ25MzTah6EE70gFH00wNf24HSmQg",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	mClient, err := NewMinioAdminClient(session)
	client := AdminClient{Client: mClient}
	err = client.getObject(context.Background(), "testName", "name.go")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_getListObjectResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "0A1WWEB961CSSLZ6E7UT",
		STSSecretAccessKey: "++yo0SF0V5J8c+unVvnZO6MUsJlYkKY6I1UzMBd+",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiIwQTFXV0VCOTYxQ1NTTFo2RTdVVCIsImV4cCI6MTY0ODU1MzkxNiwicGFyZW50IjoidGVzdDEifQ.ifE6jEcUobnfl2M6rvlUxG_PwUzFqlMKGdvFI0cn2pBCHLlgyG71gahwOAjGvOIqhhURCnTLdCuOGlInmPnsIA",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	mClient, err := NewMinioAdminClient(session)
	client := AdminClient{Client: mClient}
	err = client.listObject(context.Background(), "testName")
	if err != nil {
		fmt.Println(err)
	}
}
