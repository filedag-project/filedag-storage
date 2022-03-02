package iam

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"sync"
)

const (
	statusEnabled  = "enabled"
	statusDisabled = "disabled"
)

// IAMSys - config system.
type IAMSys struct {
	sync.Mutex
	// Persistence layer for IAM subsystem
	store *IAMStoreSys
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (sys *IAMSys) IsAllowed(args policy.Args) bool {

	return true
}

// Init - initializes config system
func (sys *IAMSys) Init(ctx context.Context) {
	sys.Lock()
	defer sys.Unlock()
	sys.initStore()

}

// initStore initializes IAM stores
func (sys *IAMSys) initStore() {
	sys.store = &IAMStoreSys{newIAMLevelDBStore()}
}
