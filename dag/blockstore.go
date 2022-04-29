package dag

//import (
//	"github.com/filedag-project/filedag-storage/dag/pool"
//	"github.com/filedag-project/filedag-storage/dag/pool/user"
//	"strings"
//	"sync"
//
//	storagekv "github.com/filedag-project/filedag-storage/kv"
//	blocks "github.com/ipfs/go-block-format"
//	"github.com/ipfs/go-cid"
//	blockstore "github.com/ipfs/go-ipfs-blockstore"
//	"golang.org/x/xerrors"
//)
//
//type blostore struct {
//	dagpool pool.Dagpool
//	batch   int
//}
//
//var userinfo *user.User
//
//func (bs *blostore) DeleteBlock(cid cid.Cid) error {
//	return bs.dagpool.DeleteNode(cid, userinfo)
//}
//
//func (bs *blostore) Has(cid cid.Cid) (bool, error) {
//	_, err := bs.dagpool.GetSize(cid, userinfo)
//	if err != nil {
//		if err == storagekv.ErrNotFound {
//			return false, nil
//		}
//		return false, err
//	}
//
//	return true, nil
//}
//
//func (bs *blostore) Get(cid cid.Cid) (blocks.Block, error) {
//	return bs.dagpool.Get(cid, userinfo)
//}
//
//func (bs *blostore) GetSize(cid cid.Cid) (int, error) {
//	n, err := bs.dagpool.GetSize(cid, userinfo)
//	if err != nil && err == storagekv.ErrNotFound {
//		return -1, blockstore.ErrNotFound
//	}
//	return n, err
//}
//
//func (bs *blostore) Put(blo blocks.Block) error {
//	return bs.dagpool.Put(blo, userinfo)
//}
//
//func (bs *blostore) PutMany(blos []blocks.Block) error {
//	var errlist []string
//	var wg sync.WaitGroup
//	batchChan := make(chan struct{}, bs.batch)
//	wg.Add(len(blos))
//	for _, blo := range blos {
//		go func(bs *blostore, blo blocks.Block) {
//			defer func() {
//				<-batchChan
//			}()
//			batchChan <- struct{}{}
//			err := bs.dagpool.Put(blo, userinfo)
//			if err != nil {
//				errlist = append(errlist, err.Error())
//			}
//		}(bs, blo)
//	}
//	wg.Wait()
//	if len(errlist) > 0 {
//		return xerrors.New(strings.Join(errlist, "\n"))
//	}
//	return nil
//}
