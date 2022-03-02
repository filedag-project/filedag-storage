package iam

// AccountStatus - account status.
type AccountStatus string

// Account status per user.
const (
	AccountEnabled  AccountStatus = "enabled"
	AccountDisabled AccountStatus = "disabled"
)

// UserInfo carries information about long term users.
type UserInfo struct {
	SecretKey  string        `json:"secretKey,omitempty"`
	PolicyName string        `json:"policyName,omitempty"`
	Status     AccountStatus `json:"status"`
	MemberOf   []string      `json:"memberOf,omitempty"`
}

// AddOrUpdateUser allows to update user details such as secret key and
// account status.
type AddOrUpdateUser struct {
	SecretKey string        `json:"secretKey,omitempty"`
	Status    AccountStatus `json:"status"`
}
