package iam

import (
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	"github.com/filedag-project/filedag-storage/objectservice/store"
)

// AccountStatus - account status.
type AccountStatus string

// Account status per user.
const (
	AccountEnabled  AccountStatus = "on"
	AccountDisabled AccountStatus = "off"
)

// UserInfo carries information about long term users.
type UserInfo struct {
	AccountName          string             `json:"account_name"`
	TotalStorageCapacity uint64             `json:"total_storage_capacity"`
	BucketInfos          []store.BucketInfo `json:"bucket_infos"`
	UseStorageCapacity   uint64             `json:"use_storage_capacity"`
	PolicyName           []string           `json:"policy_name"`
	Status               AccountStatus      `json:"status"`
}

// UserIdentity represents a user's secret key and their status
type UserIdentity struct {
	Credentials          auth.Credentials `json:"credentials"`
	TotalStorageCapacity uint64           `json:"total_storage_capacity"`
}

func newUserIdentity(cred auth.Credentials) UserIdentity {
	return UserIdentity{Credentials: cred}
}
