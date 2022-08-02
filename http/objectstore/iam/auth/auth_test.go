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
		name      string
		accessKey string
		secretKey string
		expected  bool
	}{
		{
			name:      "test1",
			accessKey: "test",
			secretKey: "ahsiuhfuiasdf",
			expected:  true,
		},
		{
			name:      "test2",
			accessKey: "test",
			secretKey: "test",
			expected:  false,
		},
		{
			name:      "test3",
			accessKey: "",
			secretKey: "ahsiuhfuiasdf",
			expected:  false,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			_, err := CreateCredentials(testcase.accessKey, testcase.secretKey)
			if (err != nil) == testcase.expected {
				t.Errorf("testcase %v failed", testcase.name)
			}
		})
	}
}
func TestCredentials_IsExpired(t *testing.T) {
	testcases := []struct {
		name       string
		accessKey  string
		secretKey  string
		expiration int64
		expected   bool
	}{
		{
			name:       "test1",
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Unix(),
			expected:   true,
		},
		{
			name:       "test2",
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Add(time.Hour).Unix(),
			expected:   false,
		},
		{
			name:       "test3",
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Add(-time.Hour).Unix(),
			expected:   true,
		},
	}
	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			m := make(map[string]interface{})
			m["accessKey"] = testcase.accessKey
			m["secretKey"] = testcase.secretKey
			m["exp"] = testcase.expiration
			credentials, err := GetNewCredentialsWithMetadata(m, testcase.secretKey)
			if err != nil {
				t.Errorf("testcase %v failed", testcase.name)
				return
			}
			if credentials.IsExpired() != testcase.expected {
				t.Errorf("testcase %v failed expect %v", testcase.name, testcase.expected)
			}
		})
	}
}
func TestCredentials_IsTemp(t *testing.T) {
	testcases := []struct {
		name       string
		accessKey  string
		secretKey  string
		expiration int64
		expected   bool
	}{
		{
			name:       "test1",
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Unix(),
			expected:   true,
		},
		{
			name:       "test1",
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Add(time.Hour).Unix(),
			expected:   true,
		},
		{
			name:       "test1",
			accessKey:  "test",
			secretKey:  "ahsiuhfuiasdf",
			expiration: time.Now().Add(-time.Hour).Unix(),
			expected:   true,
		},
		{
			name:      "test2",
			accessKey: "test",
			secretKey: "ahsiuhfuiasdf",
			expected:  false,
		},
		{
			name:       "test2",
			accessKey:  "test",
			secretKey:  "",
			expiration: time.Now().Add(-time.Hour).Unix(),
			expected:   false,
		},
	}
	for _, testcase := range testcases {

		t.Run(testcase.name, func(t *testing.T) {
			m := make(map[string]interface{})
			m["accessKey"] = testcase.accessKey
			m["secretKey"] = testcase.secretKey
			m["exp"] = testcase.expiration
			credentials, err := GetNewCredentialsWithMetadata(m, testcase.secretKey)
			if err != nil {
				t.Errorf("testcase %v failed", testcase.name)
				return
			}
			if credentials.IsTemp() != testcase.expected {
				t.Errorf("testcase %v failed expect %v", testcase.name, testcase.expected)
			}
		})
	}
}
