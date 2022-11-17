package poolservice

import (
	"context"
	"github.com/ipfs/go-cid"
	"time"
)

//GC is a goroutine to do GC
func (d *dagPoolService) GC(ctx context.Context) {
	timer := time.NewTimer(d.gcPeriod)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			log.Info("starting GC...")
			func() {
				var finish chan<- struct{}
				taskCtx, cancel := context.WithCancel(ctx)
				defer func() {
					cancel()
					select {
					case finish <- struct{}{}:
					default:
						if finish != nil {
							// never reach here
							log.Fatal("GC error")
						}
					}
				}()
				go func() {
					select {
					case finish = <-d.gcControl.Interrupt():
						cancel()
					case <-taskCtx.Done():
					}
				}()
				if err := d.runGC(taskCtx); err != nil {
					log.Errorf("GC err: %v", err)
				}
			}()
			log.Info("GC completed")
			timer.Reset(d.gcPeriod)
		case finish := <-d.gcControl.Interrupt():
			finish <- struct{}{}
		}
	}
}

func (d *dagPoolService) runGC(ctx context.Context) error {
	keys, err := d.cacheSet.AllKeysChan(ctx)
	if err != nil {
		return err
	}

	for key := range keys {
		// is pinned?
		if has, err := d.refCounter.Has(key); err != nil {
			return err
		} else if has {
			continue
		}

		blkCid, err := cid.Decode(key)
		if err != nil {
			log.Warnw("decode cid error", "cid", key, "error", err)
			continue
		}
		node := d.slots[keyHashSlot(blkCid.String())]
		if err = d.cacheSet.Remove(key); err != nil {
			log.Warnw("remove cache key error", "cid", key, "error", err)
			continue
		}
		log.Infow("delete block", "cid", key)
		if err = node.DeleteBlock(ctx, blkCid); err != nil {
			if err := d.cacheSet.Add(key); err != nil {
				log.Errorw("rollback cache key error", "cid", key, "error", err)
			}

			log.Warnw("delete block data error", "cid", key, "error", err)
			continue
		}
	}
	return nil
}

func (d *dagPoolService) InterruptGC() {
	d.gcControl.WaitInterrupt()
}

type GcControl struct {
	interruptCh chan chan<- struct{}
}

func NewGcControl() *GcControl {
	return &GcControl{
		interruptCh: make(chan chan<- struct{}),
	}
}

func (c *GcControl) WaitInterrupt() {
	finished := make(chan struct{})
	defer close(finished)
	c.interruptCh <- finished
	<-finished
}

func (c *GcControl) Interrupt() chan chan<- struct{} {
	return c.interruptCh
}
