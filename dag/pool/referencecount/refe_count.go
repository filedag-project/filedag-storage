package referencecount

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"sync"
)

//IdentityRefe reference sys
type IdentityRefe struct {
	cacheMu sync.RWMutex
	storeMu sync.RWMutex
	DB      *uleveldb.ULevelDB
}

const dagPoolRefeCache = "dagPoolRefeCache/"
const dagPoolRefeStore = "dagPoolRefeStore/"

//AddReference add refer for block
func (i *IdentityRefe) AddReference(cid string, isCache bool) error {
	cidCode := sha256String(cid)
	var count int
	if isCache {
		i.cacheMu.Lock()
		err := i.DB.Get(dagPoolRefeCache+cidCode, &count)
		count++
		err = i.DB.Put(dagPoolRefeCache+cidCode, count)
		i.cacheMu.Unlock()
		if err != nil {
			return err
		}
	} else {
		i.storeMu.Lock()
		err := i.DB.Get(dagPoolRefeStore+cidCode, &count)
		count++
		err = i.DB.Put(dagPoolRefeStore+cidCode, count)
		i.storeMu.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

//QueryReference query block refer
func (i *IdentityRefe) QueryReference(cid string, isCache bool) (int, error) {
	cidCode := sha256String(cid)
	var count int
	if isCache {
		i.cacheMu.RLock()
		err := i.DB.Get(dagPoolRefeCache+cidCode, &count)
		i.cacheMu.RUnlock()
		if err != nil {
			return 0, err
		}
	} else {
		i.storeMu.RLock()
		err := i.DB.Get(dagPoolRefeStore+cidCode, &count)
		i.storeMu.RUnlock()
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

//RemoveReference reduce refer
func (i *IdentityRefe) RemoveReference(cid string, isCache bool) error {
	cidCode := sha256String(cid)
	var count int
	if isCache {
		i.cacheMu.Lock()
		err := i.DB.Get(dagPoolRefeCache+cidCode, &count)
		if count == 1 {
			err = i.DB.Delete(dagPoolRefeCache + cidCode)
		} else if count > 1 {
			count--
			err = i.DB.Put(dagPoolRefeCache+cidCode, count)
		} else {
			return errors.New("cid does not exist")
		}
		i.cacheMu.Unlock()
		if err != nil {
			return err
		}
	} else {
		i.storeMu.Lock()
		err := i.DB.Get(dagPoolRefeStore+cidCode, &count)
		if count == 1 {
			err = i.DB.Delete(dagPoolRefeStore + cidCode)
		} else if count > 1 {
			count--
			err = i.DB.Put(dagPoolRefeStore+cidCode, count)
		} else {
			return errors.New("cid does not exist")
		}
		i.storeMu.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

//NewIdentityRefe new a reference sys
func NewIdentityRefe(db *uleveldb.ULevelDB) (IdentityRefe, error) {
	return IdentityRefe{DB: db}, nil
}

func sha256String(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
