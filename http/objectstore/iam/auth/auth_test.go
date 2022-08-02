package auth

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateCredentials(t *testing.T) {
	accessKey, secretKey, err := generateCredentials()
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
func TestCredentials_IsValid(t *testing.T) {
	testcases := []struct {
		accessKey string
		secretKey string
		expected  bool
	}{
		{
			accessKey: "test",
			secretKey: "ahsiuhfuiasdf",
			expected:  true,
		},
		{
			accessKey: "test",
			secretKey: "test",
			expected:  false,
		},
		{
			accessKey: "",
			secretKey: "ahsiuhfuiasdf",
			expected:  false,
		},
	}
	for i, testcase := range testcases {
		_, err := CreateCredentials(testcase.accessKey, testcase.secretKey)
		if (err != nil) == testcase.expected {
			t.Errorf("testcase %v failed", i)
		}
	}
}
func TestCredentials_IsExpired(t *testing.T) {
	testcases := []struct {
		accessKey  string
		secretKey  string
		expiration int64
		expected   bool
	}{
		{
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Unix(),
			expected:   true,
		},
		{
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Add(time.Hour).Unix(),
			expected:   false,
		},
		{
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Add(-time.Hour).Unix(),
			expected:   true,
		},
	}
	for i, testcase := range testcases {
		m := make(map[string]interface{})
		m["accessKey"] = testcase.accessKey
		m["secretKey"] = testcase.secretKey
		m["exp"] = testcase.expiration
		credentials, err := GetNewCredentialsWithMetadata(m, testcase.secretKey)
		if err != nil {
			t.Errorf("testcase %v failed", i)
			return
		}
		if credentials.IsExpired() != testcase.expected {
			t.Errorf("testcase %v failed expect %v", i, testcase.expected)
		}
	}
}
func TestCredentials_IsTemp(t *testing.T) {
	testcases := []struct {
		accessKey  string
		secretKey  string
		expiration int64
		expected   bool
	}{
		{
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Unix(),
			expected:   true,
		},
		{
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Add(time.Hour).Unix(),
			expected:   true,
		},
		{
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Add(-time.Hour).Unix(),
			expected:   true,
		},
		{
			accessKey: "test",
			secretKey: "ahsiuhfuiasdf",
			expected:  false,
		},
		{
			accessKey:  "test",
			secretKey:  "",
			expiration: time.Now().Add(-time.Hour).Unix(),
			expected:   false,
		},
	}
	for i, testcase := range testcases {
		m := make(map[string]interface{})
		m["accessKey"] = testcase.accessKey
		m["secretKey"] = testcase.secretKey
		m["exp"] = testcase.expiration
		credentials, err := GetNewCredentialsWithMetadata(m, testcase.secretKey)
		if err != nil {
			t.Errorf("testcase %v failed", i)
			return
		}
		if credentials.IsTemp() != testcase.expected {
			t.Errorf("testcase %v failed expect %v", i, testcase.expected)
		}
	}
}
