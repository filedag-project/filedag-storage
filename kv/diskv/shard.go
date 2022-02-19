package diskv

import (
	"path/filepath"

	"golang.org/x/xerrors"
)

const defaultShardBytes = 4
const defaultShardLevel = 2

type ShardFun func(key string) (parent, path string, err error)

func DefaultShardFun(key string) (parent, path string, err error) {
	keyLen := len(key)
	if keyLen < defaultShardBytes*defaultShardLevel {
		return "", "", xerrors.New("key is too short")
	}
	for i := 0; i < defaultShardLevel; i++ {
		parent = filepath.Join(parent, key[keyLen-(i+1)*defaultShardBytes:keyLen-i*defaultShardBytes])
	}
	path = filepath.Join(parent, key)
	return
}
