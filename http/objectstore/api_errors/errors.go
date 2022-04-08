package api_errors

import "errors"

var (
	//ErrConfigNotFound config file not found
	ErrConfigNotFound = errors.New("config file not found")
)

// GenericBucketError - generic object layer error.
type GenericBucketError struct {
	Bucket string
	Err    error
}
