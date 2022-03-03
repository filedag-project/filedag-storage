package iamerrors

import "errors"

var (
	//ErrConfigNotFound config file not found
	ErrConfigNotFound = errors.New("config file not found")
	// ErrInvalidArgument means that input argument is invalid.
	ErrInvalidArgument = errors.New("Invalid arguments specified")

	// ErrNoSuchUser error returned to IAM subsystem when user doesn't exist.
	ErrNoSuchUser = errors.New("Specified user does not exist")

	// ErrNoSuchGroup error returned to IAM subsystem when groups doesn't exist.
	ErrNoSuchGroup = errors.New("Specified group does not exist")

	// ErrGroupNotEmpty error returned to IAM subsystem when a non-empty group needs to be
	// deleted.
	ErrGroupNotEmpty = errors.New("Specified group is not empty - cannot remove it")

	// ErrNoSuchPolicy error returned to IAM subsystem when policy doesn't exist.
	ErrNoSuchPolicy = errors.New("Specified canned policy does not exist")

	// ErrPolicyInUse error returned when policy to be deleted is in use.
	ErrPolicyInUse = errors.New("Specified policy is in use and cannot be deleted.")
)

// GenericBucketError - generic object layer error.
type GenericBucketError struct {
	Bucket string
	Err    error
}
