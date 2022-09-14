package iam

import "github.com/filedag-project/filedag-storage/objectservice/iam/auth"

// AccountStatus - account status.
type AccountStatus string

// Account status per user.
const (
	AccountEnabled  AccountStatus = "on"
	AccountDisabled AccountStatus = "off"
)

// UserInfo carries information about long term users.
type UserInfo struct {
	SecretKey  string        `json:"secretKey,omitempty"`
	PolicyName []string      `json:"policyName,omitempty"`
	Status     AccountStatus `json:"status"`
}

// UserIdentity represents a user's secret key and their status
type UserIdentity struct {
	Credentials auth.Credentials `json:"credentials"`
}

func newUserIdentity(cred auth.Credentials) UserIdentity {
	return UserIdentity{Credentials: cred}
}
