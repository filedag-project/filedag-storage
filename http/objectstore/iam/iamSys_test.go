package iam

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"testing"
)

func TestIsAllowed(t *testing.T) {
	db, _ := uleveldb.OpenDb(t.TempDir())
	iamSys := NewIdentityAMSys(db)
	a := iamSys.IsAllowed(context.Background(), auth.Args{
		AccountName: "test",
		Groups:      nil,
		Action:      "",
		BucketName:  "",
		IsOwner:     false,
		ObjectName:  "",
	})
	fmt.Println(a)
}
