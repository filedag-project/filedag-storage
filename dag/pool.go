package dag

import (
	"context"
	"io"
)

type DAGPool interface {
	Add(ctx context.Context, r io.ReadCloser, user, pass string) (cidstr string, err error)
	AddWithSize(ctx context.Context, r io.ReadCloser, fsize int64, user, pass string) (cidstr string, err error)
	Get(ctx context.Context, cidstr string, user, pass string) (r io.ReadSeekCloser, err error)
}
