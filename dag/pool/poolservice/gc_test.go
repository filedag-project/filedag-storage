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
		name                           string
		theBlockDataAdded              blocks.Block
		theBlockDataAddedToInterruptGC blocks.Block
		isPinWhenAddFile               bool
		addFileByPinToInterruptGc      bool
		addFileWithoutPinToInterruptGc bool
	}{
		{
			name:                           "add file by pin and don't interrupt gc",
			theBlockDataAdded:              blocks.NewBlock(bytes.Repeat([]byte("1234"), 1)),
			theBlockDataAddedToInterruptGC: blocks.NewBlock(bytes.Repeat([]byte("1234"), 1)),
			isPinWhenAddFile:               true,
		},
		{
			name:                           "add file without pin and don't interrupt gc",
			theBlockDataAdded:              blocks.NewBlock(bytes.Repeat([]byte("12345"), 1)),
			theBlockDataAddedToInterruptGC: blocks.NewBlock(bytes.Repeat([]byte("123456"), 1)),
			isPinWhenAddFile:               false,
		},
		{
			name:                           "add file without pin and adding another file without pin to interrupt gc ",
			theBlockDataAdded:              blocks.NewBlock(bytes.Repeat([]byte("123457"), 1)),
			theBlockDataAddedToInterruptGC: blocks.NewBlock(bytes.Repeat([]byte("1234578"), 1)),
			isPinWhenAddFile:               false,
			addFileByPinToInterruptGc:      false,
			addFileWithoutPinToInterruptGc: true,
		},
		{
			name:                           "add file without pin and adding another file by pin to interrupt gc ",
			theBlockDataAdded:              blocks.NewBlock(bytes.Repeat([]byte("12345789"), 1)),
			theBlockDataAddedToInterruptGC: blocks.NewBlock(bytes.Repeat([]byte("123457890"), 1)),
			isPinWhenAddFile:               false,
			addFileByPinToInterruptGc:      true,
			addFileWithoutPinToInterruptGc: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//these are operations
			interruptGC = false
			err := service.Add(context.TODO(), tc.theBlockDataAdded, user, pass, tc.isPinWhenAddFile)
			if err != nil {
				t.Fatalf("add block err:%v", err)
			}
			err = service.Add(context.TODO(), tc.theBlockDataAddedToInterruptGC, user, pass, tc.isPinWhenAddFile)
			if err != nil {
				t.Fatalf("add block err:%v", err)
			}

			if tc.addFileByPinToInterruptGc {
				interruptGC = true
				<-startgc
				service.InterruptGC()
			}
			// sleep to see gc status
			time.Sleep(time.Second * 7)
		})
	}
}

var startgc = make(chan int)
var interruptGC bool

func (d dagPoolService) GCTest(ctx context.Context) {
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
		node, err := d.getDagNodeInfo(ctx, blkCid)
		if err != nil {
			return err
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
