package datanode

import (
	"context"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/kv/badger"
	"github.com/filedag-project/filedag-storage/kv/mutcask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestServer_Put(t *testing.T) {
	KeyError := status.Error(codes.Unknown, "Key cannot be empty")
	testcases := []struct {
		name        string
		kvType      KVType
		key         string
		data        []byte
		expectError error
	}{
		{"badge", KVBadge, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a"), nil},
		{"mutcask", KVMutcask, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a"), nil},
		{"badge", KVBadge, "", []byte("\b\u0002\u0012\a1\u0018\a"), KeyError},
		{"mutcask", KVMutcask, "", []byte("123"), nil},
		{"badge", KVBadge, "", []byte(""), KeyError},
		{"mutcask", KVMutcask, "", []byte(""), nil},
		{"badge", KVBadge, "@#", []byte("@#$$&*^@*"), nil},
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

			_, err := ser.Put(context.Background(), &proto.AddRequest{Key: tc.key, Data: tc.data})
			if err != nil {
				if err.Error() != tc.expectError.Error() {
					t.Errorf("expected response error %v, got %v", tc.expectError, err)
				}
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
			ser.Put(context.Background(), &proto.AddRequest{Key: tc.key, Data: tc.expectResponse})
			res, _ := ser.Get(context.Background(), &proto.GetRequest{Key: tc.key})
			if string(res.Data) != string(tc.expectResponse) {
				t.Errorf("expected response %s, got %s", tc.expectResponse, res.Data)
			}
		})
	}
}
func TestServer_Delete(t *testing.T) {
	testcases := []struct {
		name        string
		kvType      KVType
		key         string
		set         bool
		expectError error
	}{
		{"badge", KVBadge, "1234567", true, nil},
		{"mutcask", KVMutcask, "1234567", true, nil},
		{"badge", KVBadge, "122", false, nil},
		{"mutcask", KVMutcask, "", false, nil},
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
				ser.Put(context.Background(), &proto.AddRequest{Key: tc.key, Data: []byte("\b\u0002\u0012\a1234567\u0018\a")})

			}
			_, err := ser.Delete(context.Background(), &proto.DeleteRequest{Key: tc.key})
			if err != tc.expectError {
				t.Errorf("expected response error %s, got %s", tc.expectError, err)
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
		{"badge", KVBadge, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a"), HeaderSize + 13},
		{"mutcask", KVMutcask, "1234567", []byte("\b\u0002\u0012\a1234567\u0018\a"), HeaderSize + 13},
		{"badge", KVBadge, "122", []byte("1"), HeaderSize + 1},
		{"mutcask", KVMutcask, "122", []byte("1"), HeaderSize + 1},
		{"badge", KVBadge, "122", []byte(""), HeaderSize},
		{"mutcask", KVMutcask, "122", []byte(""), HeaderSize},
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

			ser.Put(context.Background(), &proto.AddRequest{Key: tc.key, Data: tc.dataBlock})

			res, _ := ser.Size(context.Background(), &proto.SizeRequest{Key: tc.key})
			if res.Size != tc.expectResponse {
				t.Errorf("expected response %d, got %d", tc.expectResponse, res.Size)
			}
		})
	}
}
