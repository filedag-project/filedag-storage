package policy

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"golang.org/x/xerrors"
	"unicode/utf8"
)

type Statement struct {
	SID       ID                 `json:"Sid,omitempty"`
	Effect    Effect             `json:"Effect"`
	Principal Principal          `json:"Principal"`
	Actions   s3action.ActionSet `json:"Action"`
	Resources ResourceSet        `json:"Resource"`
}

// ID - policy ID.
type ID string

// IsValid - checks if ID is valid or not.
func (id ID) IsValid() bool {
	return utf8.ValidString(string(id))
}

// Effect - policy statement effect Allow or Deny.
type Effect string

const (
	// Allow - allow effect.
	Allow Effect = "Allow"

	// Deny - deny effect.
	Deny = "Deny"
)

// Equals checks if two statements are equal
func (statement Statement) Equals(st Statement) bool {
	if statement.Effect != st.Effect {
		return false
	}
	if !statement.Principal.Equals(st.Principal) {
		return false
	}
	if !statement.Actions.Equals(st.Actions) {
		return false
	}
	if !statement.Resources.Equals(st.Resources) {
		return false
	}
	if !statement.Resources.Equals(st.Resources) {
		return false
	}
	return true
}

// Clone clones Statement structure
func (statement Statement) Clone() Statement {
	return NewStatement(statement.SID, statement.Effect, statement.Principal.Clone(),
		statement.Actions.Clone(), statement.Resources.Clone())
}

// NewStatement - creates new statement.
func NewStatement(sid ID, effect Effect, principal Principal, actionSet s3action.ActionSet, resourceSet ResourceSet) Statement {
	return Statement{
		SID:       sid,
		Effect:    effect,
		Principal: principal,
		Actions:   actionSet,
		Resources: resourceSet,
	}
}

// IsValid - checks whether statement is valid or not.
func (statement Statement) IsValid() error {
	if !statement.Principal.IsValid() {
		return xerrors.Errorf("invalid Principal %v", statement.Principal)
	}

	if len(statement.Actions) == 0 {
		return xerrors.Errorf("Action must not be empty")
	}

	if len(statement.Resources) == 0 {
		return xerrors.Errorf("Resource must not be empty")
	}

	return nil
}

// IsAllowed - returns if given check is allowed or not.
func (effect Effect) IsAllowed(b bool) bool {
	if effect == Allow {
		return b
	}
	return !b
}

// Validate - validates Statement is for given bucket or not.
func (statement Statement) Validate(bucketName string) error {
	if err := statement.IsValid(); err != nil {
		return err
	}
	return statement.Resources.Validate(bucketName)
}
