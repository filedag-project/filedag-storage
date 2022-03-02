package iamerrors

import "errors"

var (
	ErrConfigNotFound = errors.New("config file not found")
)

// GenericBucketError - generic object layer error.
type GenericBucketError struct {
	Bucket string
	Err    error
}
