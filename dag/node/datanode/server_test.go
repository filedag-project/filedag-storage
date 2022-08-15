package datanode

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/kv/badger"
	"github.com/filedag-project/filedag-storage/kv/mutcask"
	"testing"
)

func TestServer_Put(t *testing.T) {
	testcases := []struct {
		name           string
		kvType         KVType
		key            string
		data           []byte
		expectResponse string
	}{
		{"badge", KVBadge, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a"), "success"},
		{"mutcask", KVMutcask, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a"), "success"},
		{"badge", KVBadge, "", []byte("\b\u0002\u0012\a1\u0018\a"), "failed"},
		{"mutcask", KVMutcask, "", []byte("123"), "success"},
		{"badge", KVBadge, "", []byte(""), "failed"},
		{"mutcask", KVMutcask, "", []byte(""), "success"},
		{"badge", KVBadge, "@#", []byte("@#$$&*^@*"), "success"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var ser server
			switch tc.kvType {
			case KVBadge:
				newBadger, err := badger.NewBadger(t.TempDir())
				if err != nil {
					return
				}
				ser = server{kvdb: newBadger}
			case KVMutcask:
				newMutcask, err := mutcask.NewMutcask(mutcask.PathConf(t.TempDir()), mutcask.CaskNumConf(6))
				if err != nil {
					return
				}
				ser = server{kvdb: newMutcask}
			}

			res, _ := ser.Put(context.Background(), &proto.AddRequest{Key: tc.key, DataBlock: tc.data})
			if res.Message != tc.expectResponse {
				t.Errorf("expected response %s, got %s", tc.expectResponse, res.Message)
			}
		})
	}

}
func TestServer_Get(t *testing.T) {
	testcases := []struct {
		name           string
		kvType         KVType
		key            string
		expectResponse []byte
	}{
		{"badge", KVBadge, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a")},
		{"mutcask", KVMutcask, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a")},
		{"badge", KVBadge, "122", []byte("\b\u0002\u0012\a1\u0018\a")},
		{"mutcask", KVMutcask, "", []byte("123")},
		{"badge", KVBadge, "122", []byte("")},
		{"mutcask", KVMutcask, "122", []byte("#")},
		{"badge", KVBadge, "@#", []byte("@#$$&*^@*")},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var ser server
			switch tc.kvType {
			case KVBadge:
				newBadger, err := badger.NewBadger(t.TempDir())
				if err != nil {
					return
				}
				ser = server{kvdb: newBadger}
			case KVMutcask:
				newMutcask, err := mutcask.NewMutcask(mutcask.PathConf(t.TempDir()), mutcask.CaskNumConf(6))
				if err != nil {
					return
				}
				ser = server{kvdb: newMutcask}
			}
			ser.Put(context.Background(), &proto.AddRequest{Key: tc.key, DataBlock: tc.expectResponse})
			res, _ := ser.Get(context.Background(), &proto.GetRequest{Key: tc.key})
			if string(res.DataBlock) != string(tc.expectResponse) {
				t.Errorf("expected response %s, got %s", tc.expectResponse, res.DataBlock)
			}
		})
	}
}
func TestServer_Delete(t *testing.T) {
	testcases := []struct {
		name           string
		kvType         KVType
		key            string
		set            bool
		expectResponse string
	}{
		{"badge", KVBadge, "1234567", true, "success"},
		{"mutcask", KVMutcask, "1234567", true, "success"},
		{"badge", KVBadge, "122", false, "success"},
		{"mutcask", KVMutcask, "", false, "success"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var ser server
			switch tc.kvType {
			case KVBadge:
				newBadger, err := badger.NewBadger(t.TempDir())
				if err != nil {
					return
				}
				ser = server{kvdb: newBadger}
			case KVMutcask:
				newMutcask, err := mutcask.NewMutcask(mutcask.PathConf(t.TempDir()), mutcask.CaskNumConf(6))
				if err != nil {
					return
				}
				ser = server{kvdb: newMutcask}
			}
			if tc.set {
				ser.Put(context.Background(), &proto.AddRequest{Key: tc.key, DataBlock: []byte("\b\u0002\u0012\a1234567\u0018\a")})

			}
			res, _ := ser.Delete(context.Background(), &proto.DeleteRequest{Key: tc.key})
			if res.Message != tc.expectResponse {
				t.Errorf("expected response %s, got %s", tc.expectResponse, res.Message)
			}
		})
	}
}
func TestServer_Size(t *testing.T) {
	testcases := []struct {
		name           string
		kvType         KVType
		key            string
		dataBlock      []byte
		expectResponse int64
	}{
		{"badge", KVBadge, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a"), 13},
		{"mutcask", KVMutcask, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a"), 13},
		{"badge", KVBadge, "122", []byte("1"), 1},
		{"mutcask", KVMutcask, "122", []byte("1"), 1},
		{"badge", KVBadge, "122", []byte(""), 0},
		{"mutcask", KVMutcask, "122", []byte(""), 0},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var ser server
			switch tc.kvType {
			case KVBadge:
				newBadger, err := badger.NewBadger(t.TempDir())
				if err != nil {
					return
				}
				ser = server{kvdb: newBadger}
			case KVMutcask:
				newMutcask, err := mutcask.NewMutcask(mutcask.PathConf(t.TempDir()), mutcask.CaskNumConf(6))
				if err != nil {
					return
				}
				ser = server{kvdb: newMutcask}
			}

			ser.Put(context.Background(), &proto.AddRequest{Key: tc.key, DataBlock: tc.dataBlock})

			res, _ := ser.Size(context.Background(), &proto.SizeRequest{Key: tc.key})
			if res.Size != tc.expectResponse {
				t.Errorf("expected response %d, got %d", tc.expectResponse, res.Size)
			}
		})
	}
}
