package policy

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"testing"
	"time"
)

func TestPolicySys_IsAllowed(t *testing.T) {
	initSys()
	if GlobalPolicySys.IsAllowed(Args{
		AccountName: auth.DefaultAccessKey,
		Action:      "list",
		BucketName:  "test",
		ObjectName:  "test",
		IsOwner:     false,
	}) {
		// Request is allowed return the appropriate access key.
		fmt.Println(true)
	}
}

func initSys() {
	var states []Statement
	ast := s3action.NewActionSet("list")
	principal := NewPrincipal(auth.DefaultAccessKey)
	states = append(states, NewStatement("1", Allow, principal, ast))
	globalBucketMetadataSys.Set("test", BucketMetadata{
		Name:    "test",
		Created: time.Time{},
		PolicyConfig: &Policy{
			ID:         "policyConfigtest",
			Statements: states,
		},
	})
}
