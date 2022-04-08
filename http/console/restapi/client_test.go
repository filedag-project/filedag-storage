package restapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"testing"
)

func TestNewConsoleCredentials(t *testing.T) {
	got, err := NewConsoleCredentials("test", "test", "us-east-1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(got)
	tokens, err := got.Get()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokens)
}

func Test_getLoginResponse(t *testing.T) {
	accessKey := "test"
	secretKey := "test"
	loginRequest := &models.LoginRequest{
		AccessKey: &accessKey,
		SecretKey: &secretKey,
	}
	response, err := getLoginResponse(loginRequest)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
}
