package referencecount

import (
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"strconv"
	"sync"
	"time"
)

//ReferSys reference sys
type ReferSys struct {
	cacheMu sync.RWMutex
	storeMu sync.RWMutex
	DB      *uleveldb.ULevelDB
}

const gcTime = time.Minute
const dagPoolReferCache = "dagPoolReferCache/"
const dagPoolReferPin = "dagPoolReferPin/"

//AddReference add refer for block
func (i *ReferSys) AddReference(cid string, isPin bool) error {
	//cidCode := sha256String(cid)
	if !isPin {
		i.cacheMu.Lock()
		ti := time.Now().Unix()
		err := i.DB.Put(dagPoolReferCache+cid, ti)
		i.cacheMu.Unlock()
		if err != nil {
			return err
		}
	} else {
		var count int64
		i.storeMu.Lock()
		err := i.DB.Get(dagPoolReferPin+cid, &count)
		count++
		err = i.DB.Put(dagPoolReferPin+cid, count)
		i.storeMu.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

//QueryReference query block refer
func (i *ReferSys) QueryReference(cid string, isPin bool) (uint64, error) {
	//cidCode := sha256String(cid)
	if !isPin {
		ti := 0
		i.cacheMu.RLock()
		err := i.DB.Get(dagPoolReferCache+cid, &ti)
		i.cacheMu.RUnlock()
		if err != nil {
			return 0, err
		}
		if ti != 0 {
			return 2, nil
		}
		return 0, errors.New("no record")

	} else {
		var count uint64
		i.storeMu.RLock()
		err := i.DB.Get(dagPoolReferPin+cid, &count)
		i.storeMu.RUnlock()
		if err != nil {
			return 0, err
		}
		return count, nil
	}
}

//RemoveReference reduce refer
func (i *ReferSys) RemoveReference(cid string, isPin bool) error {
	//cidCode := sha256String(cid)

	if isPin {
		var count int
		i.storeMu.Lock()
		err := i.DB.Get(dagPoolReferPin+cid, &count)
		if count == 0 {
			return errors.New("cid does not exist")
		} else if count >= 1 {
			count--
			err = i.DB.Put(dagPoolReferPin+cid, count)
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

//QueryAllCacheReference query all cache refer record
func (i *ReferSys) QueryAllCacheReference() (map[string]int64, error) {
	i.cacheMu.RLock()
	defer i.cacheMu.RUnlock()
	all, err := i.DB.ReadAll(dagPoolReferCache)
	if err != nil {
		return nil, err
	}
	var m = make(map[string]int64)
	for k, v := range all {
		m[k], _ = strconv.ParseInt(v, 10, 64)
	}
	return m, nil
}

//RemoveRecord remove record in db
func (i *ReferSys) RemoveRecord(c string) error {
	i.cacheMu.Lock()
	defer i.cacheMu.Unlock()

	err := i.DB.Delete(c)
	if err != nil {
		return err
	}
	return nil
}

//NewIdentityRefe new a reference sys
func NewIdentityRefe(db *uleveldb.ULevelDB) *ReferSys {
	return &ReferSys{DB: db}
}

//
//func sha256String(s string) string {
//	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
//}
