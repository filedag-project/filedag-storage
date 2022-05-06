package dag

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"io"
)

type DAGPool interface {
	Add(ctx context.Context, r io.ReadCloser) (cidstr string, err error)
	AddWithSize(ctx context.Context, r io.ReadCloser, fsize int64) (cidstr string, err error)
	Get(ctx context.Context, cidstr string) (r io.ReadSeekCloser, err error)
	AddUser(ctx context.Context, user, pass string, policy userpolicy.DagPoolPolicy, cap uint64) error
}
