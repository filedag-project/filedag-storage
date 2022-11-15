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

func TestPinAndReference(t *testing.T) {
	t.SkipNow() //delete this to test
	//utils.SetupLogLevels()
	user, pass := "dagpool", "dagpool"
	service := startTestDagPoolServer(t)
	go service.GC(context.Background())
	defer service.Close()
	testCases := []struct {
		name                                        string
		theNumberOfTimesAFileWasAddedByPin          int
		theNumberOfTimesAFileWasAddedWithoutPin     int
		theNumberOfTimesAFileWasRemovedWithoutUnPin int
		theNumberOfTimesAFileWasRemovedByUnPin      int
		theBlockDataWeAdded                         blocks.Block
		expectReference                             int64
	}{
		{
			name:                               "adding file by pin (1)",
			theBlockDataWeAdded:                blocks.NewBlock(bytes.Repeat([]byte("1234"), 1)),
			theNumberOfTimesAFileWasAddedByPin: 1,
			expectReference:                    1,
		},
		{
			name:                                    "adding file without pin(1)",
			theBlockDataWeAdded:                     blocks.NewBlock(bytes.Repeat([]byte("12345"), 1)),
			theNumberOfTimesAFileWasAddedWithoutPin: 1,
			expectReference:                         0,
		},
		{
			name:                                    "adding 2 files by pin,adding  1 file without pin",
			theBlockDataWeAdded:                     blocks.NewBlock(bytes.Repeat([]byte("123456"), 1)),
			theNumberOfTimesAFileWasAddedByPin:      2,
			theNumberOfTimesAFileWasAddedWithoutPin: 1,
			expectReference:                         2,
		},
		{
			name:                                    "adding 2 files without pin",
			theBlockDataWeAdded:                     blocks.NewBlock(bytes.Repeat([]byte("1234567"), 1)),
			theNumberOfTimesAFileWasAddedByPin:      0,
			theNumberOfTimesAFileWasAddedWithoutPin: 2,
			expectReference:                         0,
		},
		{
			name:                                    "adding 1 files by pin,adding  2 file without pin",
			theBlockDataWeAdded:                     blocks.NewBlock(bytes.Repeat([]byte("12345678"), 1)),
			theNumberOfTimesAFileWasAddedByPin:      1,
			theNumberOfTimesAFileWasAddedWithoutPin: 2,
			expectReference:                         1,
		},
		{
			name:                                    "0 testcase",
			theBlockDataWeAdded:                     blocks.NewBlock(bytes.Repeat([]byte("123456789"), 1)),
			theNumberOfTimesAFileWasAddedByPin:      0,
			theNumberOfTimesAFileWasAddedWithoutPin: 0,
			expectReference:                         0,
		},
		{
			name:                                    "adding 1 files by pin,adding 1 file without pin,removing 1 file by unpin",
			theBlockDataWeAdded:                     blocks.NewBlock(bytes.Repeat([]byte("12345679"), 1)),
			theNumberOfTimesAFileWasAddedByPin:      1,
			theNumberOfTimesAFileWasAddedWithoutPin: 1,
			theNumberOfTimesAFileWasRemovedByUnPin:  1,
			expectReference:                         0,
		},
		{
			name:                                    "adding 2 files by pin,adding 1 file without pin,removing 1 file by unpin",
			theBlockDataWeAdded:                     blocks.NewBlock(bytes.Repeat([]byte("1234569"), 1)),
			theNumberOfTimesAFileWasAddedByPin:      2,
			theNumberOfTimesAFileWasAddedWithoutPin: 1,
			theNumberOfTimesAFileWasRemovedByUnPin:  1,
			expectReference:                         1,
		},
		{
			name:                                    "adding 2 files by pin,adding 1 file without pin,removing 2 file by unpin",
			theBlockDataWeAdded:                     blocks.NewBlock(bytes.Repeat([]byte("123459"), 1)),
			theNumberOfTimesAFileWasAddedByPin:      2,
			theNumberOfTimesAFileWasAddedWithoutPin: 1,
			theNumberOfTimesAFileWasRemovedByUnPin:  2,
			expectReference:                         0,
		},
		{
			name:                                        "adding 2 files by pin,adding 1 file without pin,removing 2 file by unpin",
			theBlockDataWeAdded:                         blocks.NewBlock(bytes.Repeat([]byte("12349"), 1)),
			theNumberOfTimesAFileWasAddedByPin:          2,
			theNumberOfTimesAFileWasAddedWithoutPin:     1,
			theNumberOfTimesAFileWasRemovedWithoutUnPin: 2,
			theNumberOfTimesAFileWasRemovedByUnPin:      0,
			expectReference:                             2,
		},
		{
			name:                                        "adding 2 files by pin,adding 1 file without pin,removing 2 file by unpin,removing 1 file without unpin",
			theBlockDataWeAdded:                         blocks.NewBlock(bytes.Repeat([]byte("12359"), 1)),
			theNumberOfTimesAFileWasAddedByPin:          2,
			theNumberOfTimesAFileWasAddedWithoutPin:     1,
			theNumberOfTimesAFileWasRemovedWithoutUnPin: 2,
			theNumberOfTimesAFileWasRemovedByUnPin:      1,
			expectReference:                             1,
		},
	}
	for _, tc := range testCases {
		ctx := context.Background()
		t.Run(tc.name, func(t *testing.T) {
			//Depending on the number of times for each operation, we perform the corresponding operation
			for i := 0; i < tc.theNumberOfTimesAFileWasAddedByPin; i++ {
				err := service.Add(ctx, tc.theBlockDataWeAdded, user, pass, true)
				if err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < tc.theNumberOfTimesAFileWasAddedWithoutPin; i++ {
				err := service.Add(ctx, tc.theBlockDataWeAdded, user, pass, false)
				if err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < tc.theNumberOfTimesAFileWasRemovedWithoutUnPin; i++ {
				err := service.Remove(ctx, tc.theBlockDataWeAdded.Cid(), user, pass, false)
				if err != nil {
					t.Fatal(err)
				}
			}
			for i := 0; i < tc.theNumberOfTimesAFileWasRemovedByUnPin; i++ {
				err := service.Remove(ctx, tc.theBlockDataWeAdded.Cid(), user, pass, true)
				if err != nil {
					t.Fatal(err)
				}
			}
			get, _ := service.refCounter.Get(tc.theBlockDataWeAdded.Cid().String())
			if get != tc.expectReference {
				t.Fatalf("expectReference should be %v,but %v", tc.expectReference, get)
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
	return service
}
