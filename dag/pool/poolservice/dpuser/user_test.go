package dpuser

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/poolservice/dpuser/upolicy"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"testing"
)

func newTestIdentityUserSys(t *testing.T) (*IdentityUserSys, error) {
	db, _ := objmetadb.OpenDb(t.TempDir())
	sys, err := NewIdentityUserSys(db, "pool", "pool123")
	if err != nil {
		return sys, err
	}
	return sys, nil
}
func TestIdentityUserSys_AddUser(t *testing.T) {
	sys, err := newTestIdentityUserSys(t)
	if err != nil {
		t.Fatalf("newTestIdentityUserSys %v", err)
		return
	}
	err = sys.AddUser(DagPoolUser{
		Username: "test",
		Password: "test123",
		Policy:   upolicy.ReadWrite,
		Capacity: 0,
	})
	if err != nil {
		t.Fatalf("AddUser %v", err)
		return
	}
	fmt.Println("ok")
}
func TestIdentityUserSys_QueryUser(t *testing.T) {
	sys, err := newTestIdentityUserSys(t)
	if err != nil {
		t.Fatalf("newTestIdentityUserSys %v", err)
		return
	}
	err = sys.AddUser(DagPoolUser{
		Username: "test",
		Password: "test123",
		Policy:   upolicy.ReadWrite,
		Capacity: 0,
	})
	if err != nil {
		t.Fatalf("AddUser %v", err)
		return
	}
	user, err := sys.QueryUser("test")
	if err != nil {
		t.Fatalf("QueryUser %v", err)
		return
	}
	fmt.Printf("user %v\n", user)
	_, err = sys.QueryUser("test2")
	if err == nil {
		t.Fatalf("shouldn't get user")
		return
	}
}
func TestIdentityUserSys_RemoveUser(t *testing.T) {
	sys, err := newTestIdentityUserSys(t)
	if err != nil {
		t.Fatalf("newTestIdentityUserSys %v", err)
		return
	}
	err = sys.AddUser(DagPoolUser{
		Username: "test",
		Password: "test123",
		Policy:   upolicy.ReadWrite,
		Capacity: 0,
	})
	if err != nil {
		t.Fatalf("AddUser %v", err)
		return
	}
	user, err := sys.QueryUser("test")
	if err != nil {
		t.Fatalf("QueryUser %v", err)
		return
	}
	fmt.Printf("user %v\n", user)
	err = sys.RemoveUser("test")
	if err != nil {
		t.Fatalf("RemoveUser %v", err)
		return
	}
	_, err = sys.QueryUser("test")
	if err == nil {
		t.Fatalf("shouldn't get user")
		return
	}
	fmt.Println("ok")
}
func TestIdentityUserSys_UpdateUser(t *testing.T) {
	sys, err := newTestIdentityUserSys(t)
	if err != nil {
		t.Fatalf("newTestIdentityUserSys %v", err)
		return
	}
	err = sys.AddUser(DagPoolUser{
		Username: "test",
		Password: "test123",
		Policy:   upolicy.ReadWrite,
		Capacity: 0,
	})
	if err != nil {
		t.Fatalf("AddUser %v", err)
		return
	}
	user, err := sys.QueryUser("test")
	if err != nil {
		t.Fatalf("QueryUser %v", err)
		return
	}
	fmt.Printf("user %v\n", user)
	err = sys.UpdateUser(DagPoolUser{
		Username: "test",
		Password: "test456",
		Policy:   upolicy.ReadOnly,
		Capacity: 0,
	})
	if err != nil {
		t.Fatalf("UpdateUser %v", err)
		return
	}
	user2, err := sys.QueryUser("test")
	if err != nil {
		t.Fatalf("QueryUser %v", err)
		return
	}
	if user2.Password != "test456" {
		t.Fatalf("update not success")
		return
	}
	fmt.Println("ok")
}
func TestIdentityUserSys_CheckUserPolicy(t *testing.T) {
	sys, err := newTestIdentityUserSys(t)
	if err != nil {
		t.Fatalf("newTestIdentityUserSys %v", err)
		return
	}
	//ReadWrite
	err = sys.AddUser(DagPoolUser{
		Username: "test",
		Password: "test123",
		Policy:   upolicy.ReadWrite,
		Capacity: 0,
	})
	if err != nil {
		t.Fatalf("AddUser %v", err)
		return
	}
	testCases := []struct {
		user   string
		pass   string
		pol    upolicy.DagPoolPolicy
		status bool
	}{
		{"test", "test123", upolicy.ReadOnly, true},
		{"test", "test123", upolicy.ReadWrite, true},
		{"test", "test123", upolicy.WriteOnly, true},
		{"test", "test", upolicy.WriteOnly, false},
	}
	for i, c := range testCases {
		if sys.CheckUserPolicy(c.user, c.pass, c.pol) != c.status {
			t.Fatalf("ReadWrite case %v ,fail", i)
		}
	}

	//OnlyWrite
	err = sys.AddUser(DagPoolUser{
		Username: "test",
		Password: "test123",
		Policy:   upolicy.WriteOnly,
		Capacity: 0,
	})
	if err != nil {
		t.Fatalf("AddUser %v", err)
		return
	}
	testCases = []struct {
		user   string
		pass   string
		pol    upolicy.DagPoolPolicy
		status bool
	}{
		{"test", "test123", upolicy.ReadOnly, false},
		{"test", "test123", upolicy.ReadWrite, false},
		{"test", "test123", upolicy.WriteOnly, true},
		{"test", "test", upolicy.WriteOnly, false},
	}
	for i, c := range testCases {
		if sys.CheckUserPolicy(c.user, c.pass, c.pol) != c.status {
			t.Fatalf("OnlyWrite case %v ,fail", i)
		}
	}
	//OnlyRead
	err = sys.AddUser(DagPoolUser{
		Username: "test",
		Password: "test123",
		Policy:   upolicy.ReadOnly,
		Capacity: 0,
	})
	if err != nil {
		t.Fatalf("AddUser %v", err)
		return
	}
	testCases = []struct {
		user   string
		pass   string
		pol    upolicy.DagPoolPolicy
		status bool
	}{
		{"test", "test123", upolicy.ReadOnly, true},
		{"test", "test123", upolicy.ReadWrite, false},
		{"test", "test123", upolicy.WriteOnly, false},
		{"test", "test", upolicy.WriteOnly, false},
	}
	for i, c := range testCases {
		if sys.CheckUserPolicy(c.user, c.pass, c.pol) != c.status {
			t.Fatalf("OnlyRead case %v ,fail", i)
		}
	}

}
