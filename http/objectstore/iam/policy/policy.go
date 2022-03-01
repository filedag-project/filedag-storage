package policy

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/action"
)

// ID - policy ID.
type ID string

// Args - arguments to policy to check whether it is allowed
type Args struct {
	AccountName string        `json:"account"`
	Groups      []string      `json:"groups"`
	Action      action.Action `json:"action"`
	BucketName  string        `json:"bucket"`
	IsOwner     bool          `json:"owner"`
	ObjectName  string        `json:"object"`
}

// Policy - bucket policy.
type Policy struct {
	ID         ID          `json:"ID,omitempty"`
	Statements []Statement `json:"Statement"`
}

// Statement - policy statement.
type Statement struct {
	SID       ID               `json:"Sid,omitempty"`
	Effect    Effect           `json:"Effect"`
	Principal Principal        `json:"Principal"`
	Actions   action.ActionSet `json:"Action"`
}

// Effect - policy statement effect Allow or Deny.
type Effect string

const (
	// Allow - allow effect.
	Allow Effect = "Allow"

	// Deny - deny effect.
	Deny = "Deny"
)

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (policy Policy) IsAllowed(args Args) bool {
	// Check all deny statements. If any one statement denies, return false.
	for _, statement := range policy.Statements {
		if statement.Effect == Deny {
			if !statement.IsAllowed(args) {
				return false
			}
		}
	}

	// For owner, its allowed by default.
	if args.IsOwner {
		return true
	}

	// Check all allow statements. If any one statement allows, return true.
	for _, statement := range policy.Statements {
		if statement.Effect == Allow {
			if statement.IsAllowed(args) {
				return true
			}
		}
	}

	return false
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (statement Statement) IsAllowed(args Args) bool {
	check := func() bool {
		if !statement.Principal.Match(args.AccountName) {
			return false
		}

		if !statement.Actions.Contains(args.Action) {
			return false
		}
		return true

	}
	return statement.Effect.IsAllowed(check())
}

// IsAllowed - returns if given check is allowed or not.
func (effect Effect) IsAllowed(b bool) bool {
	if effect == Allow {
		return b
	}
	return !b
}

// NewStatement - creates new statement.
func NewStatement(sid ID, effect Effect, principal Principal, actionSet action.ActionSet) Statement {
	return Statement{
		SID:       sid,
		Effect:    effect,
		Principal: principal,
		Actions:   actionSet,
	}
}
