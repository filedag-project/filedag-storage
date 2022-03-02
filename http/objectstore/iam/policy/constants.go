package policy

import "github.com/filedag-project/filedag-storage/http/objectstore/iam/action"

const (
	PutObjectAction         = "s3:PutObject"
	GetBucketLocationAction = "s3:GetBucketLocation"
	GetObjectAction         = "s3:GetObject"
)

// DefaultPolicies - list of canned policies available in FileDagStorage.
var DefaultPolicies = []struct {
	Name       string
	Definition Policy
}{
	// ReadWrite - provides full access to all buckets and all objects.
	{
		Name: "readwrite",
		Definition: Policy{
			Statements: []Statement{
				{
					SID:     "",
					Effect:  Allow,
					Actions: action.NewActionSet(action.AllActions),
				},
			},
		},
	},

	// ReadOnly - read only.
	{
		Name: "readonly",
		Definition: Policy{
			Statements: []Statement{
				{
					SID:     "",
					Effect:  Allow,
					Actions: action.NewActionSet(GetBucketLocationAction, GetObjectAction),
				},
			},
		},
	},

	// WriteOnly - provides write access.
	{
		Name: "writeonly",
		Definition: Policy{

			Statements: []Statement{
				{
					SID:     "",
					Effect:  Allow,
					Actions: action.NewActionSet(PutObjectAction),
				},
			},
		},
	},
}
