package restapi

import (
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func Test_getListUsersResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "PPB2H0J6RLDINHXCTICB",
		STSSecretAccessKey: "uyK1awAXLA+IgE5EOOWg9a0W9qN7Y3iHqgWXo5zR",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJQUEIySDBKNlJMRElOSFhDVElDQiIsImV4cCI6MTY1MDM1NTgyNiwicGFyZW50IjoidGVzdCJ9.-nH15EdPGtIlj3dNFtWPdzw80ZzLgMR8g7t5YhpEPnQVQ-wvRLHB86cKJ4KBN91VlgAx3pzZeseuY1Itn5Rsmw",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	apiServer := ApiServer{}
	users, err := apiServer.GetListUsersResponse(session)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(users)
	if users != nil {
		bytes, _ := json.Marshal(users)
		fmt.Println("users:", string(bytes))
	}
}

func Test_getUserAddResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "PPB2H0J6RLDINHXCTICB",
		STSSecretAccessKey: "uyK1awAXLA+IgE5EOOWg9a0W9qN7Y3iHqgWXo5zR",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJQUEIySDBKNlJMRElOSFhDVElDQiIsImV4cCI6MTY1MDM1NTgyNiwicGFyZW50IjoidGVzdCJ9.-nH15EdPGtIlj3dNFtWPdzw80ZzLgMR8g7t5YhpEPnQVQ-wvRLHB86cKJ4KBN91VlgAx3pzZeseuY1Itn5Rsmw",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	accessKey := "admin1"
	secretKey := "admin1234"
	param := models.AddUserParams{
		Body: &models.AddUserRequest{
			AccessKey: &accessKey,
			SecretKey: &secretKey,
		},
	}
	apiServer := ApiServer{}
	user, err := apiServer.GetUserAddResponse(session, param)
	if err != nil {
		fmt.Println(err)
	}
	if user != nil {
		bytes, _ := json.Marshal(user)
		fmt.Println("users:", string(bytes))
	}
}

func Test_getRemoveUserResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "PPB2H0J6RLDINHXCTICB",
		STSSecretAccessKey: "uyK1awAXLA+IgE5EOOWg9a0W9qN7Y3iHqgWXo5zR",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJQUEIySDBKNlJMRElOSFhDVElDQiIsImV4cCI6MTY1MDM1NTgyNiwicGFyZW50IjoidGVzdCJ9.-nH15EdPGtIlj3dNFtWPdzw80ZzLgMR8g7t5YhpEPnQVQ-wvRLHB86cKJ4KBN91VlgAx3pzZeseuY1Itn5Rsmw",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	name := "admin"
	param := models.RemoveUserParams{
		Name: name,
	}
	apiServer := ApiServer{}
	err := apiServer.RemoveUserResponse(session, param)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err)
}

func Test_getUserInfoResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "PPB2H0J6RLDINHXCTICB",
		STSSecretAccessKey: "uyK1awAXLA+IgE5EOOWg9a0W9qN7Y3iHqgWXo5zR",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJQUEIySDBKNlJMRElOSFhDVElDQiIsImV4cCI6MTY1MDM1NTgyNiwicGFyZW50IjoidGVzdCJ9.-nH15EdPGtIlj3dNFtWPdzw80ZzLgMR8g7t5YhpEPnQVQ-wvRLHB86cKJ4KBN91VlgAx3pzZeseuY1Itn5Rsmw",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	name := "admin"
	param := models.GetUserInfoParams{
		Name: name,
	}
	apiServer := ApiServer{}
	user, err := apiServer.GetUserInfoResponse(session, param)
	if err != nil {
		fmt.Println(err)
	}
	if user != nil {
		byte, _ := json.Marshal(user)
		fmt.Println("user:", string(byte))
	}
}

func Test_putUserPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "ND04P21NLS472ILQ7O4A",
		STSSecretAccessKey: "mrRX9dPPvKfJK+1xNkpDLRCOZ+Dnjz02wIBmL+9X",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJORDA0UDIxTkxTNDcySUxRN080QSIsImV4cCI6MTY0OTIzNzM2OSwicGFyZW50IjoidGVzdCJ9.KWSRJujOlwdVjVBqfvInP_umvupeAkQc5r4MCvkqEb_Q_rI0G4yGe04FIif1uVb_pDSEGHrXVYV-XFBxA9LDRQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	policy := `{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::test22/*"}]}`
	var str = "CUSTOM"
	custom := models.BucketAccess(str)
	request := &models.SetUserPolicyParams{
		Access:     &custom,
		Name:       "read",
		Definition: policy,
	}
	apiServer := ApiServer{}
	err := apiServer.GetUserSetPolicyResponse(session, "test2", request)
	if err != nil {
		fmt.Println(err)
	}
}

func Test_listUserPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "ND04P21NLS472ILQ7O4A",
		STSSecretAccessKey: "mrRX9dPPvKfJK+1xNkpDLRCOZ+Dnjz02wIBmL+9X",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJORDA0UDIxTkxTNDcySUxRN080QSIsImV4cCI6MTY0OTIzNzM2OSwicGFyZW50IjoidGVzdCJ9.KWSRJujOlwdVjVBqfvInP_umvupeAkQc5r4MCvkqEb_Q_rI0G4yGe04FIif1uVb_pDSEGHrXVYV-XFBxA9LDRQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	apiServer := ApiServer{}
	userPolicys, err := apiServer.ListUserPolicyResponse(session, "test2")
	if err != nil {
		fmt.Println(err)
	}
	if userPolicys != nil {
		byte, _ := json.Marshal(userPolicys)
		fmt.Println("userPolicys:", string(byte))
	}
}

func Test_getUserPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "ND04P21NLS472ILQ7O4A",
		STSSecretAccessKey: "mrRX9dPPvKfJK+1xNkpDLRCOZ+Dnjz02wIBmL+9X",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJORDA0UDIxTkxTNDcySUxRN080QSIsImV4cCI6MTY0OTIzNzM2OSwicGFyZW50IjoidGVzdCJ9.KWSRJujOlwdVjVBqfvInP_umvupeAkQc5r4MCvkqEb_Q_rI0G4yGe04FIif1uVb_pDSEGHrXVYV-XFBxA9LDRQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	apiServer := ApiServer{}
	userPolicy, err := apiServer.GetUserPolicyResponse(session, "test2")
	if err != nil {
		fmt.Println(err)
	}
	if userPolicy != nil {
		byte, _ := json.Marshal(userPolicy)
		fmt.Println("userPolicy:", string(byte))
	}
}

func Test_removeUserPolicyResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "ND04P21NLS472ILQ7O4A",
		STSSecretAccessKey: "mrRX9dPPvKfJK+1xNkpDLRCOZ+Dnjz02wIBmL+9X",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJORDA0UDIxTkxTNDcySUxRN080QSIsImV4cCI6MTY0OTIzNzM2OSwicGFyZW50IjoidGVzdCJ9.KWSRJujOlwdVjVBqfvInP_umvupeAkQc5r4MCvkqEb_Q_rI0G4yGe04FIif1uVb_pDSEGHrXVYV-XFBxA9LDRQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	apiServer := ApiServer{}
	err := apiServer.RemoveUserPolicyResponse(session, "test2", "read")
	if err != nil {
		fmt.Println(err)
	}
}
