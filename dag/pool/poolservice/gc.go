package poolservice

import (
	"context"
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
		}
	}
}
func (d *dagPoolService) UnPinGc(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Warnf("ctx done")
			return nil
		case <-time.After(gcUnPinTime):
			d.gcl.Lock()
			//time.Sleep(time.Second * 5)
			if err := d.RunUnpinGC(ctx); err != nil {
				d.gcl.Unlock()
				log.Error(err)
			} else {
				d.gcl.Unlock()
			}
		}
	}
}
func (d *dagPoolService) RunUnpinGC(ctx context.Context) error {
	log.Warnf("RunUnpinGC")
	needGCCids, err := d.refer.QueryAllStoreNonRefer()
	if err != nil {
		return err
	}
	if len(needGCCids) == 0 {
		log.Warnf("no need for unpin gc")
	}
	m := make(map[string][]string)
	for _, v := range needGCCids {
		name, _ := d.nrSys.Get(v)
		m[name] = append(m[name], v)
	}
	for n, cids := range m {
		err := d.dagNodes[n].DeleteManyBlock(ctx, cids)
		if err != nil {
			log.Errorf("DeleteManyBlock err:%v", err)
			continue
		}
		err = d.refer.RemoveRecord(cids, true)
	}
	return nil
}

func (d *dagPoolService) RunGC(ctx context.Context) error {
	needGCCids, err := d.refer.QueryAllCacheReference()
	if err != nil {
		return err
	}
	if len(needGCCids) == 0 {
		log.Warnf("no need for gc")
	}
	m := make(map[string][]string)
	for _, v := range needGCCids {
		name, _ := d.nrSys.Get(v)
		m[name] = append(m[name], v)
	}
	for n, cids := range m {
		err := d.dagNodes[n].DeleteManyBlock(ctx, cids)
		if err != nil {
			log.Errorf("DeleteManyBlock err:%v", err)
			continue
		}
		d.refer.RemoveRecord(cids, false)
	}
	return nil
}
