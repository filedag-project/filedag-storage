package policy

import (
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy/condition"
	"github.com/filedag-project/filedag-storage/objectservice/iam/s3action"
	"golang.org/x/xerrors"
	"strings"
	"unicode/utf8"
)

//https://docs.aws.amazon.com/zh_cn/zh_cn/IAM/latest/UserGuide/reference_policies_elements.html

// Statement {
//    "Version": "2012-10-17",
//    "Statement": [
//        {
//            "Sid": "Only allow writes to my bucket with bucket owner full control",
//            "Effect": "Allow",
//            "Principal": {
//                "AWS": [
//                    "arn:aws:iam::111122223333:user/ExampleUser"
//                ]
//            },
//            "Action": [
//                "s3:PutObject"
//            ],
//            "Resource": "arn:aws:s3:::DOC-EXAMPLE-BUCKET/*",
//            "Conditions": {
//                "StringEquals": {
//                    "s3:x-amz-acl": "bucket-owner-full-control"
//                }
//            }
//        }
//    ]
type Statement struct {
	SID        ID                   `json:"Sid"`
	Effect     Effect               `json:"Effect"`
	Principal  Principal            `json:"Principal"`
	Actions    s3action.ActionSet   `json:"Action"`
	Resources  ResourceSet          `json:"Resource"`
	Conditions condition.Conditions `json:"Condition,omitempty"`
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
	if !statement.Conditions.Equals(st.Conditions) {
		return false
	}
	return true
}

// Clone clones Statement structure
func (statement Statement) Clone() Statement {
	return NewStatement(statement.SID, statement.Effect, statement.Principal.Clone(),
		statement.Actions.Clone(), statement.Resources.Clone(), statement.Conditions.Clone())
}

// NewStatement - creates new statement.
func NewStatement(sid ID, effect Effect, principal Principal, actionSet s3action.ActionSet, resourceSet ResourceSet,
	conditions condition.Conditions) Statement {
	return Statement{
		SID:        sid,
		Effect:     effect,
		Principal:  principal,
		Actions:    actionSet,
		Resources:  resourceSet,
		Conditions: conditions,
	}
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (statement Statement) IsAllowed(args auth.Args) bool {
	check := func() bool {
		if !statement.Principal.Match(args.AccountName) {
			return false
		}

		if !statement.Actions.Match(args.Action) {
			return false
		}

		resource := args.BucketName
		if args.ObjectName != "" {
			if !strings.HasPrefix(args.ObjectName, "/") {
				resource += "/"
			}

			resource += args.ObjectName
		} else {
			resource += "/"
		}

		// For admin statements, resource match can be ignored.
		if !statement.Resources.Match(resource, args.Conditions) {
			return false
		}
		return statement.Conditions.Evaluate(args.Conditions)

	}
	return statement.Effect.IsAllowed(check())
}

// IsValid - checks whether statement is valid or not.
func (statement Statement) IsValid() error {
	if !statement.Effect.IsValid() {
		return xerrors.Errorf("invalid Effect %v", statement.Effect)
	}

	if !statement.Principal.IsValid() {
		return xerrors.Errorf("invalid Principal %v", statement.Principal)
	}

	if len(statement.Actions) == 0 {
		return xerrors.Errorf("Action must not be empty")
	}

	if len(statement.Resources) == 0 {
		return xerrors.Errorf("Resource must not be empty")
	}

	for action := range statement.Actions {
		if action.Match(s3action.AllActions) {
			continue
		}
		if action.IsObjectAction() {
			if !statement.Resources.ObjectResourceExists() {
				return xerrors.Errorf("unsupported Resource found %v for action %v", statement.Resources, action)
			}
		} else {
			if !statement.Resources.BucketResourceExists() {
				return xerrors.Errorf("unsupported Resource found %v for action %v", statement.Resources, action)
			}
		}

		keys := statement.Conditions.Keys()
		keyDiff := keys.Difference(s3action.ActionConditionKeyMap[action])
		if !keyDiff.IsEmpty() {
			return xerrors.Errorf("unsupported condition keys '%v' used for action '%v'", keyDiff, action)
		}
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

// IsValid - checks if Effect is valid or not
func (effect Effect) IsValid() bool {
	switch effect {
	case Allow, Deny:
		return true
	}

	return false
}

// Validate - validates Statement is for given bucket or not.
func (statement Statement) Validate(bucketName string) error {
	if err := statement.IsValid(); err != nil {
		return err
	}
	return statement.Resources.Validate(bucketName)
}
