package restapi

import (
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/console/restapi/operations/admin_api"
	"testing"
)

func Test_getListUsersResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "G4H0E4D5X3C9T5V2WTR8",
		STSSecretAccessKey: "kI+8s47yq783LjgGlQafGMLUj4koAJEEGnBcNjez",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJHNEgwRTRENVgzQzlUNVYyV1RSOCIsImV4cCI6MTY0OTIzMjg1MywicGFyZW50IjoidGVzdCJ9.QgvVQT5JorikftAF_D0ZTb6ofA_lulieM5YFhXlXBlPdHaKNsJAafvt74tQU7og6nrZBPIocniIF4mIigKNvsQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	users, err := getListUsersResponse(session)
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
		STSAccessKeyID:     "CZJLK4KJUG02NY2UHQQA",
		STSSecretAccessKey: "Ql3ppQPiF7eB+Y4at+TUIbWafiuEt0Wst77SWLF1",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJDWkpMSzRLSlVHMDJOWTJVSFFRQSIsImV4cCI6MTY0OTIyODkxMiwicGFyZW50IjoidGVzdCJ9.T-lutD97PK5IsuiNRtKBCiMZYg5wI0o1SjbKSvBmdYItkUptF1x3s91RXFNtZxhRrtbOxGqHtE3lAlVZXxaaoQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	accessKey := "admin"
	secretKey := "admin1234"
	param := admin_api.AddUserParams{
		Body: &models.AddUserRequest{
			AccessKey: &accessKey,
			SecretKey: &secretKey,
		},
	}

	user, err := getUserAddResponse(session, param)
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
		STSAccessKeyID:     "CZJLK4KJUG02NY2UHQQA",
		STSSecretAccessKey: "Ql3ppQPiF7eB+Y4at+TUIbWafiuEt0Wst77SWLF1",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJDWkpMSzRLSlVHMDJOWTJVSFFRQSIsImV4cCI6MTY0OTIyODkxMiwicGFyZW50IjoidGVzdCJ9.T-lutD97PK5IsuiNRtKBCiMZYg5wI0o1SjbKSvBmdYItkUptF1x3s91RXFNtZxhRrtbOxGqHtE3lAlVZXxaaoQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	name := "admin"
	param := admin_api.RemoveUserParams{
		Name: name,
	}

	err := getRemoveUserResponse(session, param)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err)
}

func Test_getUserInfoResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "ND04P21NLS472ILQ7O4A",
		STSSecretAccessKey: "mrRX9dPPvKfJK+1xNkpDLRCOZ+Dnjz02wIBmL+9X",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJORDA0UDIxTkxTNDcySUxRN080QSIsImV4cCI6MTY0OTIzNzM2OSwicGFyZW50IjoidGVzdCJ9.KWSRJujOlwdVjVBqfvInP_umvupeAkQc5r4MCvkqEb_Q_rI0G4yGe04FIif1uVb_pDSEGHrXVYV-XFBxA9LDRQ",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	name := "admin"
	param := admin_api.GetUserInfoParams{
		Name: name,
	}

	user, err := getUserInfoResponse(session, param)
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
	request := &models.SetUserPolicyRequest{
		Access:     &custom,
		Name:       "read",
		Definition: policy,
	}
	err := getUserSetPolicyResponse(session, "test2", request)
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
	userPolicys, err := listUserPolicyResponse(session, "test2")
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
	userPolicy, err := getUserPolicyResponse(session, "test2")
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
	err := removeUserPolicyResponse(session, "test2", "read")
	if err != nil {
		fmt.Println(err)
	}
}
