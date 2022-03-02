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

	// For owner, it allowed by default.
	if args.IsOwner {
		return true
	}

	// Check all allow statements. If anyone statement allows, return true.
	for _, statement := range policy.Statements {
		if statement.Effect == Allow {
			if statement.IsAllowed(args) {
				return true
			}
		}
	}

	return false
}

// IsAllowed - returns if given check is allowed or not.
func (effect Effect) IsAllowed(b bool) bool {
	if effect == Allow {
		return b
	}
	return !b
}

// Merge merges two policies documents and drop
// duplicate statements if any.
func (policy Policy) Merge(input Policy) Policy {
	var mergedPolicy Policy
	for _, st := range policy.Statements {
		mergedPolicy.Statements = append(mergedPolicy.Statements, st.Clone())
	}
	for _, st := range input.Statements {
		mergedPolicy.Statements = append(mergedPolicy.Statements, st.Clone())
	}
	mergedPolicy.dropDuplicateStatements()
	return mergedPolicy
}
func (policy *Policy) dropDuplicateStatements() {
redo:
	for i := range policy.Statements {
		for _, statement := range policy.Statements[i+1:] {
			if !policy.Statements[i].Equals(statement) {
				continue
			}
			policy.Statements = append(policy.Statements[:i], policy.Statements[i+1:]...)
			goto redo
		}
	}
}

// Equals returns true if the two policies are identical
func (policy *Policy) Equals(p Policy) bool {
	if policy.ID != p.ID {
		return false
	}
	if len(policy.Statements) != len(p.Statements) {
		return false
	}
	for i, st := range policy.Statements {
		if !p.Statements[i].Equals(st) {
			return false
		}
	}
	return true
}

// IsEmpty - returns whether policy is empty or not.
func (policy Policy) IsEmpty() bool {
	return len(policy.Statements) == 0
}

// Validate - validates all statements are for given bucket or not.
func (policy Policy) Validate() error {
	return policy.isValid()
}

// isValid - checks if Policy is valid or not.
func (policy Policy) isValid() error {

	for _, statement := range policy.Statements {
		if err := statement.isValid(); err != nil {
			return err
		}
	}
	return nil
}
