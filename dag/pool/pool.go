package pool

import (
	"github.com/filedag-project/filedag-storage/dag/pool/user"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	storagekv "github.com/filedag-project/filedag-storage/kv"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"golang.org/x/xerrors"
	"strings"
	"sync"
)

type Dagpool struct {
	kv       storagekv.KVDB
	batch    int
	identity user.IdentityUser
}

func (d *Dagpool) DeleteNode(cid cid.Cid, user *user.User) error {
	if !d.identity.CheckUserPolicy(user.Username, user.Password, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
	}
	return d.kv.Delete(cid.String())
}

func (d *Dagpool) Has(cid cid.Cid) (bool, error) {
	_, err := d.kv.Size(cid.String())
	if err != nil {
		if err == storagekv.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (d *Dagpool) Get(cid cid.Cid, user *user.User) (blocks.Block, error) {
	if !d.identity.CheckUserPolicy(user.Username, user.Password, userpolicy.OnlyWrite) {
		return nil, userpolicy.AccessDenied
	}
	data, err := d.kv.Get(cid.String())
	if err != nil {
		if err == storagekv.ErrNotFound {
			return nil, blockstore.ErrNotFound
		}
		return nil, err
	}
	b, err := blocks.NewBlockWithCid(data, cid)
	if err == blocks.ErrWrongHash {
		return nil, blockstore.ErrHashMismatch
	}
	return b, err
}

func (d *Dagpool) GetSize(cid cid.Cid, user *user.User) (int, error) {
	if !d.identity.CheckUserPolicy(user.Username, user.Password, userpolicy.OnlyWrite) {
		return 0, userpolicy.AccessDenied
	}
	n, err := d.kv.Size(cid.String())
	if err != nil && err == storagekv.ErrNotFound {
		return -1, blockstore.ErrNotFound
	}
	return n, err
}

func (d *Dagpool) Put(block blocks.Block, user *user.User) error {
	if !d.identity.CheckUserPolicy(user.Username, user.Password, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
	}
	//1、验证用户名
	//2、添加引用计数
	//3、存储
	return d.kv.Put(block.Cid().String(), block.RawData())
}

func (d *Dagpool) PutMany(blos []blocks.Block, user *user.User) error {
	if !d.identity.CheckUserPolicy(user.Username, user.Password, userpolicy.OnlyWrite) {
		return userpolicy.AccessDenied
	}
	var errlist []string
	var wg sync.WaitGroup
	batchChan := make(chan struct{}, d.batch)
	wg.Add(len(blos))
	for _, blo := range blos {
		go func(d *Dagpool, block blocks.Block) {
			defer func() {
				<-batchChan
			}()
			batchChan <- struct{}{}
			err := d.kv.Put(blo.Cid().String(), blo.RawData())
			if err != nil {
				errlist = append(errlist, err.Error())
			}
		}(d, blo)
	}
	wg.Wait()
	if len(errlist) > 0 {
		return xerrors.New(strings.Join(errlist, "\n"))
	}
	return nil
}
