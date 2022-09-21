package apierrors

import (
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/lock"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/filedag-project/filedag-storage/objectservice/utils/hash"
	"golang.org/x/xerrors"
)

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

func ToApiError(ctx context.Context, err error) ErrorCode {
	if ContextCanceled(ctx) {
		if ctx.Err() == context.Canceled {
			return ErrClientDisconnected
		}
	}
	errCode := ErrInternalError
	switch err.(type) {
	case lock.OperationTimedOut:
		errCode = ErrOperationTimedOut
	case hash.SHA256Mismatch:
		errCode = ErrContentSHA256Mismatch
	case hash.BadDigest:
		errCode = ErrBadDigest
	default:
		if xerrors.Is(err, store.ErrObjectNotFound) {
			errCode = ErrNoSuchKey
		}
	}
	return errCode
}
