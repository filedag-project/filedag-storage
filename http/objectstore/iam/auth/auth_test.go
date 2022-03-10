package auth

import (
	"fmt"
	"testing"
)

func TestGenerateCredentials(t *testing.T) {
	accessKey, secretKey, err := generateCredentials()
	if err != nil {
		return
	}
	fmt.Println(accessKey, secretKey)
	credentials, err := createCredentials(accessKey, secretKey)
	if err != nil {
		return
	}
	fmt.Printf("credentials %+v", credentials)
}

func TestGetNewCredentialsWithMetadata(t *testing.T) {
	accessKey, secretKey, err := generateCredentials()
	if err != nil {
		return
	}
	m := make(map[string]interface{})
	m["accessKey"] = accessKey
	credentials, err := GetNewCredentialsWithMetadata(m, secretKey)
	if err != nil {
		return
	}
	fmt.Printf("credentials %+v\n", credentials)
	credentials, err = GetNewCredentialsWithMetadata(nil, "")
	if err != nil {
		return
	}
	fmt.Printf("credentials %+v\n", credentials)
}
