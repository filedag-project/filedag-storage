package dagnode

import (
	"context"
	"errors"
)

// errNodeNotFound - cannot find the underlying configured node anymore.
var errNodeNotFound = errors.New("node not found")

// errErasureReadQuorum - did not meet read quorum.
var errErasureReadQuorum = errors.New("Read failed. Insufficient number of nodes online")

// errErasureWriteQuorum - did not meet write quorum.
var errErasureWriteQuorum = errors.New("Write failed. Insufficient number of nodes online")

// errNodeAccessDenied - we don't have write permissions on node.
var errNodeAccessDenied = errors.New("node access denied")

// Collection of basic errors.
var baseErrs = []error{
	errNodeNotFound,
}

var baseIgnoredErrs = baseErrs

// list all errors which can be ignored in entry operations.
var entryOpIgnoredErrs = append(baseIgnoredErrs, errNodeAccessDenied)

func reduceErrs(errs []error, ignoredErrs []error) (maxCount int, maxErr error) {
	errorCounts := make(map[error]int)
	for _, err := range errs {
		if IsErr(err, ignoredErrs...) {
			continue
		}
		// Errors due to context cancelation may be wrapped - group them by context.Canceled.
		if errors.Is(err, context.Canceled) {
			errorCounts[context.Canceled]++
			continue
		}
		errorCounts[err]++
	}

	max := 0
	for err, count := range errorCounts {
		switch {
		case max < count:
			max = count
			maxErr = err

		// Prefer `nil` over other error values with the same
		// number of occurrences.
		case max == count && err == nil:
			maxErr = err
		}
	}
	return max, maxErr
}

// IsErr returns whether given error is exact error.
func IsErr(err error, errs ...error) bool {
	for _, exactErr := range errs {
		if errors.Is(err, exactErr) {
			return true
		}
	}
	return false
}

// reduceQuorumErrs behaves like reduceErrs by only for returning
// values of maximally occurring errors validated against a generic
// quorum number that can be read or write quorum depending on usage.
func reduceQuorumErrs(ctx context.Context, errs []error, ignoredErrs []error, quorum int, quorumErr error) error {
	if contextCanceled(ctx) {
		return context.Canceled
	}
	maxCount, maxErr := reduceErrs(errs, ignoredErrs)
	if maxCount >= quorum {
		return maxErr
	}
	return quorumErr
}

// contextCanceled returns whether a context is canceled.
func contextCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
