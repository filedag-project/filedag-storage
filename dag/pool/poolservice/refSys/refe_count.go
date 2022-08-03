package refSys

import (
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/syndtr/goleveldb/leveldb"
	"sync"
)

// ReferSys Reference count system
type ReferSys struct {
	mu sync.RWMutex
	DB *uleveldb.ULevelDB
}

const dagPoolRefe = "dagPoolRefe/"

//AddReference add refer for block
func (i *ReferSys) AddReference(cid string) error {
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

//HasReference check if cid has reference
func (i *ReferSys) HasReference(cid string) bool {
	ti := 0
	i.mu.RLock()
	err := i.DB.Get(dagPoolRefe+cid, &ti)
	i.mu.RUnlock()
	if err == nil && ti != 0 {
		return true
	}
	return false
}

//queryReference query block refer
func (i *ReferSys) queryReference(cid string) (int, error) {
	var count int
	i.mu.RLock()
	err := i.DB.Get(dagPoolRefe+cid, &count)
	i.mu.RUnlock()
	if err != nil && err != leveldb.ErrNotFound {
		return 0, err
	}
	return count, nil
}

//RemoveReference reduce refer
func (i *ReferSys) RemoveReference(cid string) error {
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

//NewReferSys new a reference sys
func NewReferSys(db *uleveldb.ULevelDB) (*ReferSys, error) {
	return &ReferSys{DB: db}, nil
}
