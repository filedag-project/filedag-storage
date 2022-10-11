package poolservice

import (
	"bytes"
	"context"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node/datanode"
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
	go datanode.StartDataNodeServer(":9021", datanode.KVBadge, t.TempDir())
	time.Sleep(time.Second)
	go datanode.StartDataNodeServer(":9022", datanode.KVBadge, t.TempDir())
	time.Sleep(time.Second)
	go datanode.StartDataNodeServer(":9023", datanode.KVBadge, t.TempDir())
	time.Sleep(time.Second)
	var (
		dagdc = []config.DataNodeConfig{
			{
				Ip:   "127.0.0.1",
				Port: "9021",
			},
			{
				Ip:   "127.0.0.1",
				Port: "9022",
			},
			{
				Ip:   "127.0.0.1",
				Port: "9023",
			},
		}
		dagc = []config.DagNodeConfig{
			{
				Nodes:        dagdc,
				DataBlocks:   2,
				ParityBlocks: 1,
			},
		}
		cfg = config.PoolConfig{
			Listen:        "127.0.0.1:50002",
			DagNodeConfig: dagc,
			LeveldbPath:   t.TempDir(),
			RootUser:      user,
			RootPassword:  pass,
			GcPeriod:      time.Second * 5,
		}
	)
	service, err := NewDagPoolService(cfg)
	if err != nil {
		t.Fatalf("NewDagPoolService err:%v", err)
	}
	defer service.Close()
	go service.GCTest(context.Background())
	testCases := []struct {
		name           string
		bl1            blocks.Block
		bl2            blocks.Block
		pin            bool
		pinInterrupt   bool
		nopinInterrupt bool
	}{
		{
			name:           "pin-no-interrupt",
			bl1:            blocks.NewBlock(bytes.Repeat([]byte("12345"), 1)),
			bl2:            blocks.NewBlock(bytes.Repeat([]byte("12345"), 1)),
			pin:            true,
			pinInterrupt:   false,
			nopinInterrupt: false,
		},
		{
			name:           "no-pin-no-interrupt",
			bl1:            blocks.NewBlock(bytes.Repeat([]byte("123456"), 1)),
			pinInterrupt:   false,
			nopinInterrupt: false,
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

			err := service.Add(context.TODO(), tc.bl1, user, pass, tc.pin)
			if err != nil {
				t.Fatalf("add block err:%v", err)
			}
			err = service.Add(context.TODO(), tc.bl1, user, pass, tc.pin)
			if err != nil {
				t.Fatalf("add block err:%v", err)
			}

			if tc.pinInterrupt {
				<-startgc
				service.InterruptGC()
			}
			time.Sleep(time.Second * 5)
		})
	}
}

var startgc = make(chan int)

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
		startgc <- 1
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
