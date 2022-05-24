package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"io/ioutil"
	"testing"
)

func Test_getUploadObjectResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "KIBUMQ2R8LWCC5USEMFH",
		STSSecretAccessKey: "ryQOZEdygKk4dhQ9b8uGR6loQBHHIRbRPu9NXoeN",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJLSUJVTVEyUjhMV0NDNVVTRU1GSCIsImV4cCI6MTY1MDM0MDQ0OCwicGFyZW50IjoidGVzdCJ9.f2Pc-PUQrzx8zqeXHxvG2FXZZVGrR3uMWQCf8dSSXiSSCF_IqszaycvacKrCC1QZO-DhNB9JgK3rDlEUupHRHg",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	mClient, err := NewAdminClient(session)
	client := AdminClient{Client: mClient}
	r1, _ := ioutil.ReadFile("user_objects.go")
	err = client.putObject(context.Background(), "testN", "name.go", bytes.NewReader(r1), int64(len(r1)))
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
	mClient, err := NewAdminClient(session)
	client := AdminClient{Client: mClient}
	err = client.getObject(context.Background(), "testName", "name.go")
	if err != nil {
		fmt.Println(err)
	}
}

func Test_getListObjectResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "KIBUMQ2R8LWCC5USEMFH",
		STSSecretAccessKey: "ryQOZEdygKk4dhQ9b8uGR6loQBHHIRbRPu9NXoeN",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJLSUJVTVEyUjhMV0NDNVVTRU1GSCIsImV4cCI6MTY1MDM0MDQ0OCwicGFyZW50IjoidGVzdCJ9.f2Pc-PUQrzx8zqeXHxvG2FXZZVGrR3uMWQCf8dSSXiSSCF_IqszaycvacKrCC1QZO-DhNB9JgK3rDlEUupHRHg",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	params := models.ListObjectsParams{
		BucketName: "testN",
	}
	apiServer := ApiServer{}
	objects, err := apiServer.GetListObjectsResponse(session, params)
	if err != nil {
		fmt.Println(err)
	}
	if objects != nil {
		bytes, _ := json.Marshal(objects)
		fmt.Println("objects:", string(bytes))
	}
}
