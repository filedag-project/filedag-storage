package apierrors

import (
	"context"
	"errors"
)

var (
	//ErrConfigNotFound config file not found
	ErrConfigNotFound = errors.New("config file not found")
)

// GenericBucketError - generic object layer error.
type GenericBucketError struct {
	Bucket string
	Err    error
}

// NotImplemented If a feature is not implemented
type NotImplemented struct {
	Message string
}

// ContextCanceled returns whether a context is canceled.
func ContextCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
