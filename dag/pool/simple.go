package pool

import (
	"context"
	"io"
)

var _ DAGPool = (*simplePool)(nil)

type simplePool struct {
}

func (p *simplePool) Add(ctx context.Context, r io.ReadCloser) (cidstr string, err error) {
	return "", nil
}

func (p *simplePool) AddWithSize(ctx context.Context, r io.ReadCloser, fsize int64) (cidstr string, err error) {
	return "", nil
}

func (p *simplePool) Get(ctx context.Context, cidstr string) (r io.ReadCloser, err error) {
	return nil, nil
}
