package poolservice

import (
	"context"
	"sync"
	"time"
)

type gc struct {
	operateMap map[string]uint64
	lock       sync.RWMutex
}

//var m = make(map[string]string)

//Gc is a goroutine to do GC
func (d *dagPoolService) Gc(ctx context.Context, gcPeriod string) error {
	duration, err := time.ParseDuration(gcPeriod)
	if err != nil {
		return err
	}
	count := 0
	for {
		count++
		log.Warnf("Gc count:%d", count)
		select {
		case <-ctx.Done():
			log.Warnf("ctx done")
			return nil
		case <-d.CheckStorage():
			//time.Sleep(time.Second * 5)
			if err := d.runUnpinGC(ctx); err != nil {
				log.Error(err)
			} else {
				return err
			}
		case <-time.After(duration):
			log.Warnf("start do gc")
			err := d.runGC(ctx, count)
			if err != nil {
				return err
			}

		case <-time.After(duration * 2):
			log.Warnf("start do gc")
			err := d.runUnpinGC(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (d *dagPoolService) CheckStorage() <-chan int {
	//todo check storage if reaches the maximum value
	return nil
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
			//time.Sleep(time.Second * 5)
			if err := d.runUnpinGC(ctx); err != nil {
				log.Error(err)
			} else {
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
	for _, c := range needGCCids {
		if d.gc.operateMap[c.String()] != 0 {
			continue
		}
		err = d.refer.RemoveRecord(c.String(), true)
		node, err := d.GetNode(ctx, c)
		if err != nil {
			return err
		}
		node.DeleteBlock(ctx, c)
	}
	return nil
}

func (d *dagPoolService) runGC(ctx context.Context, c int) error {

	needGCCids, err := d.refer.QueryAllCacheReference()
	if err != nil {
		return err
	}
	if len(needGCCids) == 0 {
		log.Warnf("no need for gc %v", c)
		return nil
	}
	for _, ci := range needGCCids {
		d.gc.lock.RLock()

		if d.gc.operateMap[ci.String()] != 0 {
			d.gc.lock.RUnlock()
			continue
		}
		d.gc.lock.RUnlock()
		node, err := d.GetNode(ctx, ci)
		if err != nil {
			return err
		}
		d.refer.RemoveRecord(ci.String(), false)

		err = node.DeleteBlock(ctx, ci)
		if err != nil {
			log.Errorf("DeleteManyBlock err:%v", err)
			continue
		}
	}
	return nil
}
