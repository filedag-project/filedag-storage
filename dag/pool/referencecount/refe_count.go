package referencecount

import (
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

type IdentityRefe struct {
	mu sync.RWMutex
	DB *uleveldb.ULevelDB
}

const dagPoolRefe = "dagPoolRefe/"

func (i *IdentityRefe) AddReference(cid string) error {
	var count int
	i.mu.Lock()
	err := i.DB.Get(dagPoolRefe+cid, &count)
	count++
	err = i.DB.Put(dagPoolRefe+cid, count)
	i.mu.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func (i *IdentityRefe) QueryReference(cid string) (int, error) {
	var count int
	i.mu.RLock()
	err := i.DB.Get(dagPoolRefe+cid, &count)
	i.mu.RUnlock()
	if err != nil && err != leveldb.ErrNotFound {
		return 0, err
	}
	return count, nil
}

func (i *IdentityRefe) RemoveReference(cid string) error {
	var count int
	i.mu.Lock()
	err := i.DB.Get(dagPoolRefe+cid, &count)
	if count == 1 {
		err = i.DB.Delete(dagPoolRefe + cid)
	} else if count > 1 {
		count--
		err = i.DB.Put(dagPoolRefe+cid, count)
	} else {
		return errors.New("cid does not exist")
	}
	i.mu.Unlock()
	if err != nil {
		return err
	}
	return nil
}

func NewIdentityRefe(db *uleveldb.ULevelDB) (IdentityRefe, error) {
	return IdentityRefe{DB: db}, nil
}
