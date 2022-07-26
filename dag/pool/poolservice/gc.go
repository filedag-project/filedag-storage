package poolservice

import (
	"context"
	"time"
)

type gc struct {
	stopCh   chan struct{}
	gcPeriod time.Duration
}

//Stop the gc
func (g *gc) Stop() {
	g.stopCh <- struct{}{}
}

//Gc is a goroutine to do GC
func (d *dagPoolService) Gc(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Warnf("ctx done")
			return nil
		case <-d.CheckStorage():
			//time.Sleep(time.Second * 5)
			if err := d.runGC(ctx); err != nil {
				log.Error(err)
			} else {
				return err
			}
		case <-time.After(d.gc.gcPeriod):
			log.Debugf("start do gc")
			ct, cancel := context.WithCancel(ctx)
			go d.runGC(ct)
			<-d.gc.stopCh
			cancel()
		}
	}
}

func (d *dagPoolService) CheckStorage() <-chan int {
	//todo check storage if reaches the maximum value
	return nil
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
			<-d.gc.stopCh
			cancel()
		}
	}
}
func (d *dagPoolService) runStoreGC(ctx context.Context) error {
	//log.Warnf("RunUnpinGC")
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
	d.gc.stopCh <- struct{}{}
	return nil
}

func (d *dagPoolService) runGC(ctx context.Context) error {

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
			log.Errorf("DeleteManyBlock err:%v", err)
			continue
		}
	}
	d.gc.stopCh <- struct{}{}
	return nil
}
