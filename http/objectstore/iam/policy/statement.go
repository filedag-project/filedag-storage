package policy

import (
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
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
//            "Condition": {
//                "StringEquals": {
//                    "s3:x-amz-acl": "bucket-owner-full-control"
//                }
//            }
//        }
//    ]
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
	return true
}

// Clone clones Statement structure
func (statement Statement) Clone() Statement {
	return NewStatement(statement.SID, statement.Effect, statement.Principal.Clone(),
		statement.Actions.Clone())
}

// NewStatement - creates new statement.
func NewStatement(sid ID, effect Effect, principal Principal, actionSet s3action.ActionSet) Statement {
	return Statement{
		SID:       sid,
		Effect:    effect,
		Principal: principal,
		Actions:   actionSet,
	}
}

// IsAllowed - checks given policy args is allowed to continue the Rest API.
func (statement Statement) IsAllowed(args auth.Args) bool {
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

// IsValid - checks whether statement is valid or not.
func (statement Statement) IsValid() error {

	if len(statement.Actions) == 0 {
		return errors.New(fmt.Sprintf("Action must not be empty"))
	}

	if err := statement.Actions.Validate(); err != nil {
		return err
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
