package store

import (
	"context"
	"encoding/json"
	"fmt"
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	ufsio "github.com/ipfs/go-unixfs/io"
	"io"
	"strings"
	"time"
)

var log = logging.Logger("store")

//StorageSys store sys
type StorageSys struct {
	Db         *uleveldb.ULevelDB
	DagPool    dagpoolcli.PoolClient
	CidBuilder cid.Builder
}

const objectPrefixTemplate = "object-%s-%s-%s/"
const allObjectPrefixTemplate = "object-%s-%s-"

var (
	poolUser = "pool"
	poolPass = "pool123"
)

func getPoolUser() string {
	return poolUser + "," + poolPass
}

//StoreObject store object
func (s *StorageSys) StoreObject(ctx context.Context, user, bucket, object string, reader io.ReadCloser, size int64) (ObjectInfo, error) {
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	ctx = context.WithValue(ctx, "user", getPoolUser())
	node, err := dagpoolcli.BalanceNode(ctx, reader, s.DagPool, s.CidBuilder)
	if err != nil {
		return ObjectInfo{}, err
	}
	meta := ObjectInfo{
		Bucket:           bucket,
		Name:             object,
		ModTime:          time.Now().UTC(),
		Size:             size,
		IsDir:            false,
		ETag:             node.Cid().String(),
		VersionID:        "",
		IsLatest:         true,
		DeleteMarker:     false,
		ContentType:      "application/x-msdownload",
		ContentEncoding:  "",
		Expires:          time.Unix(0, 0).UTC(),
		Parts:            nil,
		AccTime:          time.Unix(0, 0).UTC(),
		SuccessorModTime: time.Now().UTC(),
	}
	err = s.Db.Put(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), meta)
	if err != nil {
		return ObjectInfo{}, err
	}
	return meta, nil
}

//GetObject Get object
func (s *StorageSys) GetObject(ctx context.Context, user, bucket, object string) (ObjectInfo, ufsio.DagReader, error) {
	meta := ObjectInfo{}
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	err := s.Db.Get(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), &meta)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	ctx = context.WithValue(ctx, "user", getPoolUser())
	cid, err := cid.Decode(meta.ETag)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	dagNode, err := s.DagPool.Get(ctx, cid)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	reader, err := ufsio.NewDagReader(ctx, dagNode, s.DagPool)
	if err != nil {
		return ObjectInfo{}, nil, err
	}
	return meta, reader, nil
}

// HasObject has Object ?
func (s *StorageSys) HasObject(ctx context.Context, user, bucket, object string) (ObjectInfo, bool) {
	meta := ObjectInfo{}
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	err := s.Db.Get(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), &meta)
	if err != nil {
		return ObjectInfo{}, false
	}
	return meta, true
}

//DeleteObject Get object
func (s *StorageSys) DeleteObject(ctx context.Context, user, bucket, object string) error {
	//err := s.dagPool.DelFile(bucket, object)
	ctx = context.WithValue(ctx, "user", getPoolUser())
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}
	meta := ObjectInfo{}
	err := s.Db.Get(fmt.Sprintf(objectPrefixTemplate, user, bucket, object), &meta)
	if err != nil {
		return err
	}
	cid, err := cid.Decode(meta.ETag)
	if err != nil {
		return err
	}
	err = s.DagPool.Remove(ctx, cid)
	if err != nil {
		return err
	}
	err = s.Db.Delete(fmt.Sprintf(objectPrefixTemplate, user, bucket, object))
	if err != nil {
		return err
	}
	return nil
}

//ListObject list user object
func (s *StorageSys) ListObject(user, bucket string) ([]ObjectInfo, error) {
	var objs []ObjectInfo
	objMap, err := s.Db.ReadAll(fmt.Sprintf(allObjectPrefixTemplate, user, bucket))
	if err != nil {
		return nil, err
	}
	for _, v := range objMap {
		var o ObjectInfo
		json.Unmarshal([]byte(v), &o)
		objs = append(objs, o)
	}
	return objs, nil
}

//MkBucket store object
func (s *StorageSys) MkBucket(parentDirectoryPath string, bucket string) error {
	return nil
}

// ListObjectsV2Info - container for list objects version 2.
type ListObjectsV2Info struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: This element is returned only if you have delimiter request parameter
	// specified.
	ContinuationToken     string
	NextContinuationToken string

	// List of objects info for this request.
	Objects []ObjectInfo

	// List of prefixes for this request.
	Prefixes []string
}

// ListObjectsV2 list objects
//todo use more param
func (s *StorageSys) ListObjectsV2(ctx context.Context, user, bucket string, prefix string, token string, delimiter string, keys int, owner bool, after string) (ListObjectsV2Info, error) {
	objects, err := s.ListObject(user, bucket)
	var o ListObjectsV2Info
	if err != nil {
		return o, err
	}
	count := 0
	for _, v := range objects {
		if after == "" {

		} else if v.Name != after {
			continue
		}
		if count > keys {
			break
		}
		count++
		o.ContinuationToken = token
		o.IsTruncated = true
		o.Objects = append(o.Objects, v)
	}
	return o, nil
}

//Init storage sys
func (s *StorageSys) Init(poolAddr, pu, pp string) error {
	s.Db = uleveldb.DBClient
	var err error
	poolUser = pu
	poolPass = pp
	cidBuilder, err := merkledag.PrefixForCidVersion(0)
	s.CidBuilder = cidBuilder
	s.DagPool, err = dagpoolcli.NewPoolClient(poolAddr)
	if err != nil {
		return err
	}
	return nil
}

// Close storage sys
func (s *StorageSys) Close() {
	s.DagPool.Close(context.TODO())
}
