package iam

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"testing"
)

func TestIsAllowed(t *testing.T) {
	var iamSys IdentityAMSys
	uleveldb.DBClient = uleveldb.OpenDb("/tmp/leveldb2/test")
	iamSys.Init()
	a := iamSys.IsAllowed(auth.Args{
		AccountName: "test",
		Groups:      nil,
		Action:      "",
		BucketName:  "",
		IsOwner:     false,
		ObjectName:  "",
	})
	fmt.Println(a)
}
