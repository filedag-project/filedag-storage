package iam

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"testing"
)

func TestLoadUser(t *testing.T) {
	var iamSys IdentityAMSys
	iamSys.Init()
	//err := iamSys.saveUserIdentity(context.Background(), "test", UserIdentity{Credentials: auth.Credentials{
	//	AccessKey:    "test",
	//	SecretKey:    "test secret",
	//	Expiration:   time.Now(),
	//	SessionToken: "SessionToken",
	//	Status:       "on",
	//}})
	//if err != nil {
	//	return
	//}

	m := &auth.Credentials{}
	err := iamSys.store.loadUser(context.Background(), "test", m)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(m)
	//a, err := iamSys.loadUsers(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(a)
	//err := iamSys.removeUserIdentity(context.Background(), "s")
	//if err != nil {
	//	return
	//}
}
