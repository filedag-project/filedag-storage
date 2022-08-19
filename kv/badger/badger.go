package badger

import (
	"context"
	"github.com/dgraph-io/badger"
	"github.com/filedag-project/filedag-storage/kv"
)

var _ kv.KVDB = (*badgerDb)(nil)

func NewBadger(path string) (bg *badgerDb, err error) {
	opts := badger.DefaultOptions(path)
	opts.SyncWrites = false
	opts.ValueThreshold = 256
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &badgerDb{db: db}, nil
}

type badgerDb struct {
	db *badger.DB
}

func (b *badgerDb) Put(key string, value []byte) error {
	wb := b.db.NewWriteBatch()
	defer wb.Cancel()
	err := wb.SetEntry(badger.NewEntry([]byte(key), value).WithMeta(0))
	if err != nil {
		return err
	}
	return wb.Flush()
}

func (b *badgerDb) Delete(key string) error {
	wb := b.db.NewWriteBatch()
	defer wb.Cancel()
	return wb.Delete([]byte(key))
}

func (b *badgerDb) Get(key string) ([]byte, error) {
	var ival []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		ival, err = item.ValueCopy(nil)
		return err
	})
	return ival, err
}

func (b *badgerDb) Size(key string) (int, error) {
	size := 0
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		size = int(item.ValueSize())
		return err
	})
	return size, err
}

func (b *badgerDb) AllKeysChan(ctx context.Context) (chan string, error) {
	kc := make(chan string)
	go func(ctx context.Context, b *badgerDb) {
		defer close(kc)
		_ = b.db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false
			it := txn.NewIterator(opts)
			defer it.Close()
			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				k := item.Key()
				select {
				case <-ctx.Done():
					return nil
				default:
					kc <- string(k)
				}
			}
			return nil
		})
	}(ctx, b)
	return kc, nil
}

func (b *badgerDb) Close() error {
	return b.db.Close()
}
