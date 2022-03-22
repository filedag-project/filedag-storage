package store

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"io"
	"io/ioutil"
	"os"
	"time"
)

//StorageSys store sys
type StorageSys struct {
	db      *uleveldb.ULevelDB
	dagPool dagPoolClient
}

const objectPrefixTemplate = "object-%s-%s-%s/"

//StoreObject store object
func (s *StorageSys) StoreObject(user, bucket, object string, reader io.Reader) (ObjectInfo, error) {
	cid, err := s.dagPool.PutFile(bucket, object, reader)
	if err != nil {
		return ObjectInfo{}, err
	}
	all, err := ioutil.ReadAll(reader)
	meta := ObjectInfo{
		Bucket:                     bucket,
		Name:                       object,
		ModTime:                    time.Now().UTC(),
		Size:                       int64(len(all)),
		IsDir:                      false,
		ETag:                       cid,
		VersionID:                  "",
		IsLatest:                   false,
		DeleteMarker:               false,
		RestoreExpires:             time.Unix(0, 0).UTC(),
		RestoreOngoing:             false,
		ContentType:                "application/x-msdownload",
		ContentEncoding:            "",
		Expires:                    time.Unix(0, 0).UTC(),
		StorageClass:               "STANDARD",
		UserDefined:                nil,
		UserTags:                   "",
		Parts:                      nil,
		Writer:                     nil,
		Reader:                     nil,
		PutObjReader:               nil,
		AccTime:                    time.Unix(0, 0).UTC(),
		Legacy:                     false,
		VersionPurgeStatusInternal: "",
		VersionPurgeStatus:         "",
		NumVersions:                0,
		SuccessorModTime:           time.Now().UTC(),
	}
	err = s.db.Put(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), meta)
	if err != nil {
		return ObjectInfo{}, err
	}
	return meta, nil
}

//GetObject Get object
func (s *StorageSys) GetObject(user, bucket, object string) (ObjectInfo, io.Reader, error) {
	reader, err := s.dagPool.GetFile(bucket, object)
	meta := ObjectInfo{}
	err = s.db.Get(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), &meta)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	return meta, reader, nil
}

//DeleteObject Get object
func (s *StorageSys) DeleteObject(user, bucket, object string) error {
	err := s.dagPool.DelFile(bucket, object)
	err = s.db.Delete(fmt.Sprintf(objectPrefixTemplate, user, bucket, object))
	if err != nil {
		return err
	}
	return nil
}

//MkBucket store object
func (s *StorageSys) MkBucket(parentDirectoryPath string, bucket string) error {
	err := os.Mkdir(parentDirectoryPath+bucket, 0777)
	if err != nil {
		return err
	}
	return nil
}

//Init storage sys
func (s *StorageSys) Init() {
	s.db = uleveldb.DBClient
}
