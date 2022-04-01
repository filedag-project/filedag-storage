package restapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/console/restapi/operations/admin_api"
	"testing"
)

func Test_getListUsersResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "JHWCTOHEC4EURPLMKFNJ",
		STSSecretAccessKey: "GHbQ9uP0qIaripnoujoxicjHK5x5Z45LpaYvclR+",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJKSFdDVE9IRUM0RVVSUExNS0ZOSiIsImV4cCI6MTY0ODc5ODY5MiwicGFyZW50IjoidGVzdCJ9.0xXpeEetsWDo_29KpoMC4BzLj7N16Acm0zugW6npl3ZcjJ8YMFTcPXW8zhXOYFCRkorqhlbcPl9mIx9gcjiuxw",
		AccountAccessKey:   "test",
		Hm:                 false,
	}
	got, err := getListUsersResponse(session)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(got)
}

func Test_getUserAddResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "JHWCTOHEC4EURPLMKFNJ",
		STSSecretAccessKey: "GHbQ9uP0qIaripnoujoxicjHK5x5Z45LpaYvclR+",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJKSFdDVE9IRUM0RVVSUExNS0ZOSiIsImV4cCI6MTY0ODc5ODY5MiwicGFyZW50IjoidGVzdCJ9.0xXpeEetsWDo_29KpoMC4BzLj7N16Acm0zugW6npl3ZcjJ8YMFTcPXW8zhXOYFCRkorqhlbcPl9mIx9gcjiuxw",
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
	fmt.Println(user)
}

func Test_getRemoveUserResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:     "JHWCTOHEC4EURPLMKFNJ",
		STSSecretAccessKey: "GHbQ9uP0qIaripnoujoxicjHK5x5Z45LpaYvclR+",
		STSSessionToken:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJKSFdDVE9IRUM0RVVSUExNS0ZOSiIsImV4cCI6MTY0ODc5ODY5MiwicGFyZW50IjoidGVzdCJ9.0xXpeEetsWDo_29KpoMC4BzLj7N16Acm0zugW6npl3ZcjJ8YMFTcPXW8zhXOYFCRkorqhlbcPl9mIx9gcjiuxw",
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
