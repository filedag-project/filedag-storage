package iam

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"testing"
)

func TestLoadUsers(t *testing.T) {
	db, _ := uleveldb.OpenDb(t.TempDir())
	iamSys := NewIdentityAMSys(db)

	//add user
	//err := iamSys.AddUser(context.Background(), "test1", "test12345")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//list user
	mc, err := iamSys.store.loadUsers(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mc)
	//err := iamSys.store.loadUser(context.Background(), "test", m)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(m)
	//a, err := iamSys.loadUsers(r.Context())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(a)
	//err := iamSys.removeUserIdentity(r.Context(), "s")
	//if err != nil {
	//	return
	//}
}
