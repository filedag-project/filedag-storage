package pool

import (
	"context"
	"io"
)

type DAGPool interface {
	Add(ctx context.Context, r io.ReadCloser) (cidstr string, err error)
	AddWithSize(ctx context.Context, r io.ReadCloser, fsize int64) (cidstr string, err error)
	Get(ctx context.Context, cidstr string) (r io.ReadSeekCloser, err error)
}
