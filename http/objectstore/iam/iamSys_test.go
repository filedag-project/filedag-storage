package iam

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"testing"
)

func TestIsAllowed(t *testing.T) {
	db, _ := uleveldb.OpenDb(utils.TmpDirPath(&testing.T{}))
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
