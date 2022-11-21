package poolservice

import (
	"bytes"
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"testing"
	"time"
)

func Test_Gc(t *testing.T) {
	t.SkipNow() //delete this to test
	utils.SetupLogLevels()
	user, pass := "dagpool", "dagpool"
	service := startTestDagPoolServer(t)
	go service.GCTest(context.Background())
	defer service.Close()
	testCases := []struct {
		name           string
		bl1            blocks.Block
		bl2            blocks.Block
		pin            bool
		pinInterrupt   bool
		nopinInterrupt bool
	}{
		{
			name: "pin-no-interrupt",
			bl1:  blocks.NewBlock(bytes.Repeat([]byte("1234"), 1)),
			bl2:  blocks.NewBlock(bytes.Repeat([]byte("1234"), 1)),
			pin:  true,
		},
		{
			name: "no-pin-no-interrupt",
			bl1:  blocks.NewBlock(bytes.Repeat([]byte("12345"), 1)),
			bl2:  blocks.NewBlock(bytes.Repeat([]byte("123456"), 1)),
			pin:  false,
		},
		{
			name:           "no-pin-no-pin-interrupt",
			bl1:            blocks.NewBlock(bytes.Repeat([]byte("123457"), 1)),
			bl2:            blocks.NewBlock(bytes.Repeat([]byte("1234578"), 1)),
			pin:            false,
			pinInterrupt:   false,
			nopinInterrupt: true,
		},
		{
			name:           "no-pin-pin-interrupt",
			bl1:            blocks.NewBlock(bytes.Repeat([]byte("12345789"), 1)),
			bl2:            blocks.NewBlock(bytes.Repeat([]byte("123457890"), 1)),
			pin:            false,
			pinInterrupt:   true,
			nopinInterrupt: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interruptGC = false
			err := service.Add(context.TODO(), tc.bl1, user, pass, tc.pin)
			if err != nil {
				t.Fatalf("add block err:%v", err)
			}
			err = service.Add(context.TODO(), tc.bl2, user, pass, tc.pin)
			if err != nil {
				t.Fatalf("add block err:%v", err)
			}

			if tc.pinInterrupt {
				interruptGC = true
				<-startgc
				service.InterruptGC()
			}
			time.Sleep(time.Second * 7)
		})
	}
}

var startgc = make(chan int)
var interruptGC bool

func (d *dagPoolService) GCTest(ctx context.Context) {
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
				if err := d.runGCTest(taskCtx); err != nil {
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

//IExactly the same logic as runGC, just increase the deletion time to test the GC interruption problem
func (d *dagPoolService) runGCTest(ctx context.Context) error {
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
		if err = d.cacheSet.Remove(key); err != nil {
			log.Warnw("remove cache key error", "cid", key, "error", err)
			continue
		}
		//Increase delete time to test for GC interruptions
		if interruptGC {
			startgc <- 1
		}
		time.Sleep(time.Second)

		log.Infow("delete block", "cid", key)
		if err = d.deleteBlock(ctx, blkCid); err != nil {
			if err := d.cacheSet.Add(key); err != nil {
				log.Errorw("rollback cache key error", "cid", key, "error", err)
			}

			log.Warnw("delete block data error", "cid", key, "error", err)
			continue
		}
	}
	return nil
}
