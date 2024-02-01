package reference

import (
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"github.com/syndtr/goleveldb/leveldb"
	"strings"
)

const CachePrefix = "cache/"

type CacheSet struct {
	db objmetadb.ObjStoreMetaDBAPI
}

func NewCacheSet(db objmetadb.ObjStoreMetaDBAPI) *CacheSet {
	return &CacheSet{db: db}
}

func (s *CacheSet) Add(key string) error {
	exist := true
	return s.db.Put(CachePrefix+key, &exist)
}

func (s *CacheSet) Has(key string) (bool, error) {
	var exist bool
	err := s.db.Get(CachePrefix+key, &exist)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *CacheSet) AllKeysChan(ctx context.Context) (<-chan string, error) {
	all, err := s.db.ReadAllChan(ctx, CachePrefix, "")
	if err != nil {
		return nil, err
	}
	kc := make(chan string)
	go func() {
		defer close(kc)
		for entry := range all {
			strs := strings.Split(entry.GetKey(), "/")
			if len(strs) < 2 {
				return
			}
			select {
			case <-ctx.Done():
				return
			case kc <- strs[1]:
			}
		}
	}()

	return kc, nil
}

func (s *CacheSet) Remove(key string) error {
	return s.db.Delete(CachePrefix + key)
}
