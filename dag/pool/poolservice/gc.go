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
			}
			d.gcl.Unlock()
		}
	}
}

func (d *dagPoolService) RunGC(ctx context.Context) error {
	referMap, err := d.refer.QueryAllCacheReference()
	if err != nil {
		return err
	}
	if len(referMap) == 0 {
		log.Warnf("no need to gc")
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
