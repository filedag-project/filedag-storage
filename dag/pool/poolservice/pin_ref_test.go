package poolservice

import (
	"bytes"
	"context"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node/datanode"
	blocks "github.com/ipfs/go-block-format"
	"testing"
	"time"
)

func TestPinRef(t *testing.T) {
	t.SkipNow() //delete this to test
	//utils.SetupLogLevels()
	user, pass := "dagpool", "dagpool"
	service := startTestDagPoolServer(t)
	go service.GC(context.Background())
	defer service.Close()
	testCases := []struct {
		name          string
		pinAddTimes   int
		cacheAddTimes int
		rmCacheTimes  int
		rmUnPinTimes  int
		bl1           blocks.Block
		ref           int64
	}{
		{
			name:          "pin-1-0",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("1234"), 1)),
			pinAddTimes:   1,
			cacheAddTimes: 0,
			ref:           1,
		},
		{
			name:          "pin-0-1",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("12345"), 1)),
			pinAddTimes:   0,
			cacheAddTimes: 1,
			ref:           0,
		},
		{
			name:          "pin-2-1",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("123456"), 1)),
			pinAddTimes:   2,
			cacheAddTimes: 1,
			ref:           2,
		},
		{
			name:          "pin-0-2",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("1234567"), 1)),
			pinAddTimes:   0,
			cacheAddTimes: 2,
			ref:           0,
		},
		{
			name:          "pin-1-2",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("12345678"), 1)),
			pinAddTimes:   1,
			cacheAddTimes: 2,
			ref:           1,
		},
		{
			name:          "pin-0-0",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("123456789"), 1)),
			pinAddTimes:   0,
			cacheAddTimes: 0,
			ref:           0,
		},
		{
			name:          "pin-1-1-1",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("12345679"), 1)),
			pinAddTimes:   1,
			cacheAddTimes: 1,
			rmUnPinTimes:  1,
			ref:           0,
		},
		{
			name:          "pin-2-1-1",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("1234569"), 1)),
			pinAddTimes:   2,
			cacheAddTimes: 1,
			rmUnPinTimes:  1,
			ref:           1,
		},
		{
			name:          "pin-2-1-2",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("123459"), 1)),
			pinAddTimes:   2,
			cacheAddTimes: 1,
			rmUnPinTimes:  2,
			ref:           0,
		},
		{
			name:          "pin-2-1-2-0",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("12349"), 1)),
			pinAddTimes:   2,
			cacheAddTimes: 1,
			rmCacheTimes:  2,
			rmUnPinTimes:  0,
			ref:           2,
		},
		{
			name:          "pin-2-1-2-1",
			bl1:           blocks.NewBlock(bytes.Repeat([]byte("12359"), 1)),
			pinAddTimes:   2,
			cacheAddTimes: 1,
			rmCacheTimes:  2,
			rmUnPinTimes:  1,
			ref:           1,
		},
	}
	for _, tc := range testCases {
		ctx := context.Background()
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i < tc.pinAddTimes; i++ {
				err := service.Add(ctx, tc.bl1, user, pass, true)
				if err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < tc.cacheAddTimes; i++ {
				err := service.Add(ctx, tc.bl1, user, pass, false)
				if err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < tc.rmCacheTimes; i++ {
				err := service.Remove(ctx, tc.bl1.Cid(), user, pass, false)
				if err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < tc.rmUnPinTimes; i++ {
				err := service.Remove(ctx, tc.bl1.Cid(), user, pass, true)
				if err != nil {
					t.Fatal(err)
				}
			}
			get, _ := service.refCounter.Get(tc.bl1.Cid().String())
			if get != tc.ref {
				t.Fatalf("ref should be %v,but %v", tc.ref, get)
			}
		})
	}
}
func startTestDagPoolServer(t *testing.T) *dagPoolService {
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
				SetIndex:   0,
				RpcAddress: "127.0.0.1:9021",
			},
			{
				SetIndex:   1,
				RpcAddress: "127.0.0.1:9022",
			},
			{
				SetIndex:   2,
				RpcAddress: "127.0.0.1:9023",
			},
		}
		dagc = config.ClusterConfig{
			Cluster: []config.DagNodeConfig{
				{
					Nodes:        dagdc,
					DataBlocks:   2,
					ParityBlocks: 1,
				},
			},
		}
		cfg = config.PoolConfig{
			Listen:        "127.0.0.1:50002",
			ClusterConfig: dagc,
			LeveldbPath:   t.TempDir(),
			RootUser:      user,
			RootPassword:  pass,
			GcPeriod:      time.Second * 5,
		}
	)
	service, err := NewDagPoolService(context.TODO(), cfg)
	if err != nil {
		t.Fatalf("NewDagPoolService err:%v", err)
	}
	return service
}
