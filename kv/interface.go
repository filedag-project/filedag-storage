package kv

import "context"

type KVDB interface {
	Put(string, []byte) error
	Delete(string) error
	Get(string) ([]byte, error)
	Size(string) (int, error)

	AllKeysChan(context.Context) (chan string, error)
	Close() error
}
