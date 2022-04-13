package iam

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"testing"
)

func TestIdentityAMSys_UserApi(t *testing.T) {
	uleveldb.DBClient, _ = uleveldb.OpenDb(path)
	var iamSys IdentityAMSys
	iamSys.Init()
	var accessKey = "test1"
	var secretKey = "test12345"

	//add user
	err := iamSys.AddUser(context.Background(), accessKey, secretKey)
	if err != nil {
		fmt.Println(err)
	}

	//list user
	users, err := iamSys.GetUserList(context.Background(), "")
	if err != nil {
		fmt.Println(users)
	}

	//get user
	cred, bool := iamSys.GetUser(context.Background(), accessKey)
	fmt.Println(bool, cred)

	//remove user
	err = iamSys.RemoveUser(context.Background(), accessKey)
	if err != nil {
		fmt.Println(err)
	}
}
