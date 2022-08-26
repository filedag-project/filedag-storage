package reference

import (
	"context"
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/xerrors"
	"strings"
	"sync"
)

const RefPrefix = "ref/"

type RefCounter struct {
	mut      sync.Mutex
	db       *uleveldb.ULevelDB
	cacheSet *CacheSet
}

func NewRefCounter(db *uleveldb.ULevelDB, cacheSet *CacheSet) *RefCounter {
	return &RefCounter{
		db:       db,
		cacheSet: cacheSet,
	}
}

func (rc *RefCounter) Incr(key string) error {
	return rc.IncrOrCreate(key, nil)
}

func (rc *RefCounter) IncrOrCreate(key string, createFunc func() error) error {
	rc.mut.Lock()
	defer rc.mut.Unlock()
	var count int64
	err := rc.db.Get(RefPrefix+key, &count)
	if xerrors.Is(err, leveldb.ErrNotFound) {
		if createFunc != nil {
			if err = createFunc(); err != nil {
				return err
			}
		}
	} else if err != nil {
		return err
	}
	count++
	return rc.db.Put(RefPrefix+key, count)
}

func (rc *RefCounter) Get(key string) (int64, error) {
	var count int64
	err := rc.db.Get(RefPrefix+key, &count)
	return count, err
}

func (rc *RefCounter) Has(key string) (bool, error) {
	var count int64
	err := rc.db.Get(RefPrefix+key, &count)
	if err != nil {
		if xerrors.Is(err, leveldb.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return count > 0, err
}

func (rc *RefCounter) Decr(key string) error {
	rc.mut.Lock()
	defer rc.mut.Unlock()
	var count int64
	err := rc.db.Get(RefPrefix+key, &count)
	if err != nil && !xerrors.Is(err, leveldb.ErrNotFound) {
		return err
	}
	if xerrors.Is(err, leveldb.ErrNotFound) || count == 0 {
		return errors.New("reference count of key is zero")
	}
	count--
	if count == 0 {
		// move to cache
		if err = rc.moveToCache(key); err != nil {
			return err
		}
		return rc.db.Delete(RefPrefix + key)
	}
	return rc.db.Put(RefPrefix+key, count)
}

// AllKeysChan query keys which reference count value equals count param
func (rc *RefCounter) AllKeysChan(ctx context.Context, count int64) (<-chan string, error) {
	all, err := rc.db.ReadAll(RefPrefix)
	if err != nil {
		return nil, err
	}
	kc := make(chan string)
	go func() {
		defer close(kc)
		for k, v := range all {
			if v == "0" {
				strs := strings.Split(k, "/")
				if len(strs) < 2 {
					return
				}
				select {
				case <-ctx.Done():
					return
				case kc <- strs[1]:
				}
			}
		}
	}()

	return kc, nil
}

//Remove remove zero reference counter
func (rc *RefCounter) Remove(key string, force bool) error {
	rc.mut.Lock()
	defer rc.mut.Unlock()
	var count int64
	err := rc.db.Get(RefPrefix+key, &count)
	if err != nil {
		if xerrors.Is(err, leveldb.ErrNotFound) {
			return nil
		}
		return err
	}
	if !force && count > 0 {
		return errors.New("reference count of key is greater than zero")
	}
	// move to cache
	if err = rc.moveToCache(key); err != nil {
		return err
	}
	return rc.db.Delete(RefPrefix + key)
}

func (rc *RefCounter) moveToCache(key string) error {
	return rc.cacheSet.Add(key)
}
