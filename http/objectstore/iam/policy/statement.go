package policy

import (
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"unicode/utf8"
)

//
// Statement - policy statement.
//{
//  "Version": "2012-10-17",
//  "Statement": [
//    {
//      "Effect": "Allow",
//      "Action": [
//        "s3:ListAllMyBuckets",
//        "s3:GetBucketLocation"
//      ],
//      "Resource": "arn:aws:s3:::*"
//    },
//    {
//      "Effect": "Allow",
//      "Action": "s3:ListBucket",
//      "Resource": "arn:aws:s3:::BUCKET-NAME",
//      "Condition": {"StringLike": {"s3:prefix": [
//        "",
//        "home/",
//        "home/${aws:username}/"
//      ]}}
//    },
//    {
//      "Effect": "Allow",
//      "Action": "s3:*",
//      "Resource": [
//        "arn:aws:s3:::BUCKET-NAME/home/${aws:username}",
//        "arn:aws:s3:::BUCKET-NAME/home/${aws:username}/*"
//      ]
//    }
//  ]
//}
type Statement struct {
	SID       ID                 `json:"Sid,omitempty"`
	Effect    Effect             `json:"Effect"`
	Principal Principal          `json:"Principal"`
	Actions   s3action.ActionSet `json:"Action"`
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
