package iam

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"testing"
)

func TestIsAllowed(t *testing.T) {
	var globalIAMSys *IAMSys
	globalIAMSys.IsAllowed(policy.Args{
		AccountName: "",
		Groups:      nil,
		Action:      "",
		BucketName:  "",
		IsOwner:     false,
		ObjectName:  "",
	})
}
