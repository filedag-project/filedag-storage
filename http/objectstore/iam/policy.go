package iam

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"unicode/utf8"
)

// ID - policy ID.
type ID string

// IsValid - checks if ID is valid or not.
func (id ID) IsValid() bool {
	return utf8.ValidString(string(id))
}

// Policy - iam bucket iamp.
type Policy struct {
	ID         ID `json:"ID,omitempty"`
	Version    string
	Statements []policy.Statement `json:"Statement"`
}
type PolicyDocument struct {
	Version   string              `json:"Version"`
	Statement []*policy.Statement `json:"Statement"`
}
type Policies struct {
	Policies map[string]PolicyDocument `json:"policies"`
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (p Policy) IsAllowed(args auth.Args) bool {
	// Check all deny statements. If any one statement denies, return false.
	for _, statement := range p.Statements {
		if statement.Effect == policy.Deny {
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
	for _, statement := range p.Statements {
		if statement.Effect == policy.Allow {
			if statement.IsAllowed(args) {
				return true
			}
		}
	}

	return false
}

// Merge merges two policies documents and drop
// duplicate statements if any.
func (p Policy) Merge(input Policy) Policy {
	var mergedPolicy Policy
	for _, st := range p.Statements {
		mergedPolicy.Statements = append(mergedPolicy.Statements, st.Clone())
	}
	for _, st := range input.Statements {
		mergedPolicy.Statements = append(mergedPolicy.Statements, st.Clone())
	}
	mergedPolicy.dropDuplicateStatements()
	return mergedPolicy
}
func (p *Policy) dropDuplicateStatements() {
redo:
	for i := range p.Statements {
		for _, statement := range p.Statements[i+1:] {
			if !p.Statements[i].Equals(statement) {
				continue
			}
			p.Statements = append(p.Statements[:i], p.Statements[i+1:]...)
			goto redo
		}
	}
}

// Equals returns true if the two policies are identical
func (p *Policy) Equals(policy Policy) bool {
	if p.ID != policy.ID {
		return false
	}
	if len(p.Statements) != len(policy.Statements) {
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
func (p Policy) IsEmpty() bool {
	return len(p.Statements) == 0
}

// Validate - validates all statements are for given bucket or not.
func (p Policy) Validate() error {
	return p.isValid()
}

// isValid - checks if Policy is valid or not.
func (p Policy) isValid() error {

	for _, statement := range p.Statements {
		if err := statement.IsValid(); err != nil {
			return err
		}
	}
	return nil
}

// DefaultPolicies - list of canned policies available in FileDagStorage.
var DefaultPolicies = []struct {
	Name       string
	Definition Policy
}{
	// ReadWrite - provides full access to all buckets and all objects.
	{
		Name: "readwrite",
		Definition: Policy{
			Statements: []policy.Statement{
				{
					SID:     "",
					Effect:  policy.Allow,
					Actions: s3action.NewActionSet(s3action.AllActions),
				},
			},
		},
	},

	// ReadOnly - read only.
	{
		Name: "readonly",
		Definition: Policy{
			Statements: []policy.Statement{
				{
					SID:     "",
					Effect:  policy.Allow,
					Actions: s3action.NewActionSet(s3action.GetBucketLocationAction, s3action.GetObjectAction),
				},
			},
		},
	},

	// WriteOnly - provides write access.
	{
		Name: "writeonly",
		Definition: Policy{

			Statements: []policy.Statement{
				{
					SID:     "",
					Effect:  policy.Allow,
					Actions: s3action.NewActionSet(s3action.PutObjectAction),
				},
			},
		},
	},
}
