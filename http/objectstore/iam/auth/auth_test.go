package auth

import (
	"fmt"
	"testing"
)

func TestGenerateCredentials(t *testing.T) {
	accessKey, secretKey, err := GenerateCredentials()
	if err != nil {
		return
	}
	fmt.Println(accessKey, secretKey)
	credentials, err := CreateCredentials(accessKey, secretKey)
	if err != nil {
		return
	}
	fmt.Printf("credentials %+v", credentials)
}

func TestGetNewCredentialsWithMetadata(t *testing.T) {
	accessKey, secretKey, err := GenerateCredentials()
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
