package iam

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"testing"
)

func TestIsAllowed(t *testing.T) {

	a := GlobalIAMSys.IsAllowed(policy.Args{
		AccountName: "test",
		Groups:      nil,
		Action:      "",
		BucketName:  "",
		IsOwner:     false,
		ObjectName:  "",
	})
	fmt.Println(a)
}
