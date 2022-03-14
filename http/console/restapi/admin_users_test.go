package restapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/credentials"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/console/pkg/auth"
	"github.com/filedag-project/filedag-storage/http/console/restapi/operations/admin_api"
	"testing"
)

func Test_getListUsersResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "W2W2JWQUI52SVMMJK5MV",
		STSSecretAccessKey: "vSZmxbVcx+lP3iUZgqqXT0PZTYH2HAEsuQ+9hLJT",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJXMlcySldRVUk1MlNWTU1KSzVNViIsImV4cCI6MTY0NjI5Mzk4OSwicGFyZW50IjoibWluaW9hZG1pbiJ9.kmC564DCOOiDUpl4FiAWDx0839tTtuuZmvxoN_tSxpnqOBt_W8zoZgASi-ag9jD29kUThnVjR4I92qIs-TTL9g",
		AccountAccessKey:   "minioadmin",
		Hm:                 false,
	}
	got, got1 := getListUsersResponse(session)
	if got1 != nil {
		fmt.Println(got1)
	}
	fmt.Println(got)
}

func Test_getUserAddResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "W2W2JWQUI52SVMMJK5MV",
		STSSecretAccessKey: "vSZmxbVcx+lP3iUZgqqXT0PZTYH2HAEsuQ+9hLJT",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJXMlcySldRVUk1MlNWTU1KSzVNViIsImV4cCI6MTY0NjI5Mzk4OSwicGFyZW50IjoibWluaW9hZG1pbiJ9.kmC564DCOOiDUpl4FiAWDx0839tTtuuZmvxoN_tSxpnqOBt_W8zoZgASi-ag9jD29kUThnVjR4I92qIs-TTL9g",
		AccountAccessKey:   "minioadmin",
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
	fmt.Println(user)
}

func Test_getRemoveUserResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "W2W2JWQUI52SVMMJK5MV",
		STSSecretAccessKey: "vSZmxbVcx+lP3iUZgqqXT0PZTYH2HAEsuQ+9hLJT",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJXMlcySldRVUk1MlNWTU1KSzVNViIsImV4cCI6MTY0NjI5Mzk4OSwicGFyZW50IjoibWluaW9hZG1pbiJ9.kmC564DCOOiDUpl4FiAWDx0839tTtuuZmvxoN_tSxpnqOBt_W8zoZgASi-ag9jD29kUThnVjR4I92qIs-TTL9g",
		AccountAccessKey:   "minioadmin",
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

func Test_decryptClaims(t *testing.T) {
	token := ""
	claims, err := auth.SessionTokenAuthenticate(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(claims)
}

func Test_encryptClaims(t *testing.T) {
	value := &credentials.Value{
		AccessKeyID:     "",
		SecretAccessKey: "",
		SessionToken:    "",
	}
	accountAccessKey := "admin"
	features := true
	sessionFeatures := &auth.SessionFeatures{
		HideMenu: features,
	}
	claims, err := auth.NewEncryptedTokenForClient(value, accountAccessKey, sessionFeatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(claims)
}
