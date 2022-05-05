package referencecount

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
)

type IdentityRefe struct {
	DB *uleveldb.ULevelDB
}

const dagPoolRefe = "dagPoolRefe/"

func (i *IdentityRefe) AddReference(cid string) error {
	cidCode := sha256String(cid)
	var count int
	_ = i.DB.Get(dagPoolRefe+cidCode, &count)
	count++
	err := i.DB.Put(dagPoolRefe+cidCode, count)
	if err != nil {
		return err
	}
	return nil
}

func (i *IdentityRefe) QueryReference(cid string) error {
	cidCode := sha256String(cid)
	err := i.DB.Get(dagPoolRefe+cidCode, 1)
	if err != nil {
		return err
	}
	return nil
}

func (i *IdentityRefe) RemoveReference(cid string) error {
	cidCode := sha256String(cid)
	var count int
	err := i.DB.Get(dagPoolRefe+cidCode, &count)
	if count == 1 {
		err = i.DB.Delete(dagPoolRefe + cidCode)
	} else if count > 1 {
		count--
		err = i.DB.Put(dagPoolRefe+cidCode, count)
	} else {
		return errors.New("cid does not exist")
	}
	if err != nil {
		return err
	}
	return nil
}

func sha256String(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
