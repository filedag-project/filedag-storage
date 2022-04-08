package kv

import (
	"context"

	"golang.org/x/xerrors"
)

var ErrNotFound = xerrors.New("kv: key not found")

type KVDB interface {
	Put(string, []byte) error
	Delete(string) error
	Get(string) ([]byte, error)
	Size(string) (int, error)

	AllKeysChan(context.Context) (chan string, error)
	Close() error
}
