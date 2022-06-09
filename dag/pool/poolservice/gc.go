package poolservice

import (
	"context"
	"github.com/ipfs/go-cid"
	"strings"
	"time"
)

//todo as commandline param
const gcExpiredTime = time.Minute
const gcTime = time.Second * 10
const gcUnPinTime = time.Second * 15

func (d *dagPoolService) Gc(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Warnf("ctx done")
			return nil
		case <-time.After(gcTime):
			d.gcl.Lock()
			//time.Sleep(time.Second * 5)
			if err := d.RunGC(ctx); err != nil {
				d.gcl.Unlock()
				log.Error(err)
			} else {
				d.gcl.Unlock()
			}
		case <-time.After(gcUnPinTime):
			d.gcl.Lock()
			//time.Sleep(time.Second * 5)
			if err := d.RunUnpinGC(ctx); err != nil {
				d.gcl.Unlock()
				log.Error(err)
			} else {

			}
		}
	}
}
func (d *dagPoolService) RunUnpinGC(ctx context.Context) error {
	s, err := d.refer.QueryAllStoreNonRefer()
	if err != nil {
		return err
	}
	if len(s) == 0 {
		log.Warnf("no need for gc")
	}
	for _, v := range s {
		c, _ := cid.Decode(strings.Split(v, "/")[1])
		log.Warnf("gc del block:%v", c.String())
		node, err := d.GetNode(ctx, c)
		if err != nil {
			continue
		}
		err = node.DeleteBlock(ctx, c)
		if err != nil {
			continue
		}
		err = d.refer.RemoveRecord(v)
	}
	return nil
}

func (d *dagPoolService) RunGC(ctx context.Context) error {
	referMap, err := d.refer.QueryAllCacheReference()
	if err != nil {
		return err
	}
	if len(referMap) == 0 {
		log.Warnf("no need for gc")
	}
	for k, v := range referMap {
		if time.Now().After(time.Unix(v, 0).Add(gcExpiredTime)) {
			c, _ := cid.Decode(strings.Split(k, "/")[1])
			reference, err := d.refer.QueryReference(c.String(), true)
			if err != nil {
				return err
			}
			if reference != 0 {
				err = d.refer.RemoveRecord(k)
				log.Warnf("gc del block:%v,but this was pin", c)
				continue
			}
			log.Warnf("gc del block:%v", c)
			node, err := d.GetNode(ctx, c)
			if err != nil {
				continue
			}
			err = node.DeleteBlock(ctx, c)
			if err != nil {
				continue
			}
			err = d.refer.RemoveRecord(k)
		}
	}
	return nil
}
