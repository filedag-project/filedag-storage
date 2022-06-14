package poolservice

import (
	"context"
	"time"
)

//Gc is a goroutine to do GC
func (d *dagPoolService) Gc(ctx context.Context, gcPeriod string) error {
	duration, err := time.ParseDuration(gcPeriod)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			log.Warnf("ctx done")
			return nil
		case <-time.After(duration):
			d.gcl.Lock()
			//time.Sleep(time.Second * 5)
			if err := d.runGC(ctx); err != nil {
				d.gcl.Unlock()
				log.Error(err)
			} else {
				d.gcl.Unlock()
			}
		}
	}
}

//UnPinGc is a goroutine to do UnPin GC
func (d *dagPoolService) UnPinGc(ctx context.Context, gcPeriod string) error {
	duration, err := time.ParseDuration(gcPeriod)
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			log.Warnf("ctx done")
			return nil
		case <-time.After(duration):
			d.gcl.Lock()
			//time.Sleep(time.Second * 5)
			if err := d.runUnpinGC(ctx); err != nil {
				d.gcl.Unlock()
				log.Error(err)
			} else {
				d.gcl.Unlock()
			}
		}
	}
}
func (d *dagPoolService) runUnpinGC(ctx context.Context) error {
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

func (d *dagPoolService) runGC(ctx context.Context) error {
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
