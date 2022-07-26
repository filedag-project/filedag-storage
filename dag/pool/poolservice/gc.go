package poolservice

import (
	"context"
	"time"
)

type gc struct {
	// gc state
	runningCache bool
	runningStore bool

	//gc stop channel
	stopCacheCh chan struct{}
	stopStoreCh chan struct{}
	normalCh    chan struct{}
	gcPeriod    time.Duration
}

//Stop the gc
func (g *gc) Stop() {
	if g.runningCache {
		g.stopCacheCh <- struct{}{}
	}
	if g.runningStore {
		g.stopStoreCh <- struct{}{}
	}
}

//Gc is a goroutine to do GC
func (d *dagPoolService) Gc(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Warnf("ctx done")
			return nil
		//case <-d.CheckStorage():
		//	//time.Sleep(time.Second * 5)
		//	if err := d.runGC(ctx); err != nil {
		//		log.Error(err)
		//	} else {
		//		return err
		//	}
		case <-time.After(d.gc.gcPeriod):
			log.Debugf("start do gc")
			ct, cancel := context.WithCancel(ctx)
			go d.runGC(ct)
			select {
			case <-d.gc.stopCacheCh:
				log.Debugf(" cache gc inter stop ")
				cancel()
			case <-d.gc.normalCh:
				log.Debugf(" cache gc normal stop ")
				cancel()
			}
		}
	}
}

//StoreGc is a goroutine to do UnPin GC
func (d *dagPoolService) StoreGc(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Warnf("ctx done")
			return nil
		case <-time.After(d.gc.gcPeriod):
			//time.Sleep(time.Second * 5)
			log.Debugf("start do store gc")
			ct, cancel := context.WithCancel(ctx)
			go d.runStoreGC(ct)
			select {
			case <-d.gc.stopStoreCh:
				log.Debugf("store gc inter stop ")
				cancel()
			case <-d.gc.normalCh:
				log.Debugf("store gc normal stop ")
				cancel()
			}
		}
	}
}
func (d *dagPoolService) runStoreGC(ctx context.Context) error {
	d.gc.runningStore = true
	defer func() {
		d.gc.normalCh <- struct{}{}
		d.gc.runningStore = false
	}()
	needGCCids, err := d.refer.QueryAllStoreNonRef()
	if err != nil {
		return err
	}
	if len(needGCCids) == 0 {
		//log.Warnf("no need for unpin gc")
		return nil
	}
	for _, c := range needGCCids {
		err = d.refer.RemoveRecord(c.String(), true)
		node, err := d.GetNode(ctx, c)
		if err != nil {
			return err
		}
		node.DeleteBlock(ctx, c)
	}
	return nil
}

func (d *dagPoolService) runGC(ctx context.Context) error {
	d.gc.runningCache = true
	defer func() {
		d.gc.normalCh <- struct{}{}
		d.gc.runningCache = false
	}()
	needGCCids, err := d.refer.QueryAllCacheRef()
	if err != nil {
		return err
	}
	if len(needGCCids) == 0 {
		//log.Warnf("no need for gc")
		return nil
	}
	for _, ci := range needGCCids {
		node, err := d.GetNode(ctx, ci)
		if err != nil {
			return err
		}
		d.refer.RemoveRecord(ci.String(), false)

		err = node.DeleteBlock(ctx, ci)
		if err != nil {
			log.Errorf("DeleteBlock err:%v", err)
			continue
		}
	}
	return nil
}

func (d *dagPoolService) CheckStorage() <-chan int {
	//todo check storage if reaches the maximum value
	return nil
}
