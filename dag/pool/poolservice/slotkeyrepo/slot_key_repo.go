package slotkeyrepo

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"github.com/syndtr/goleveldb/leveldb"
	"strings"
)

const SlotPrefix = "slot/"

// SlotKeyRepo saves information about the slot and cid key mapping.
type SlotKeyRepo struct {
	db *uleveldb.ULevelDB
}

func NewSlotKeyRepo(db *uleveldb.ULevelDB) *SlotKeyRepo {
	return &SlotKeyRepo{db: db}
}

func (s *SlotKeyRepo) Set(slot uint16, key string, value string) error {
	return s.db.Put(fmt.Sprintf("%s%v/%s", SlotPrefix, slot, key), value)
}

func (s *SlotKeyRepo) Get(slot uint16, key string) (value string, err error) {
	err = s.db.Get(fmt.Sprintf("%s%v/%s", SlotPrefix, slot, key), &value)
	return value, err
}

func (s *SlotKeyRepo) Has(slot uint16, key string) (bool, error) {
	var val string
	err := s.db.Get(fmt.Sprintf("%s%v/%s", SlotPrefix, slot, key), &val)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *SlotKeyRepo) Remove(slot uint16, key string) error {
	return s.db.Delete(fmt.Sprintf("%s%v/%s", SlotPrefix, slot, key))
}

type SlotKeyEntry struct {
	Slot  uint16
	Key   string
	Value string
}

func (s *SlotKeyRepo) AllKeysChan(ctx context.Context, slot uint16, seekSlotKey string) (<-chan *SlotKeyEntry, error) {
	prefix := fmt.Sprintf("%s%v/", SlotPrefix, slot)
	if seekSlotKey != "" {
		seekSlotKey = fmt.Sprintf("%s%v/%s", SlotPrefix, slot, seekSlotKey)
	}
	all, err := s.db.ReadAllChan(ctx, prefix, seekSlotKey)
	if err != nil {
		return nil, err
	}
	kc := make(chan *SlotKeyEntry)
	go func() {
		defer close(kc)
		for entry := range all {
			strs := strings.Split(entry.Key, "/")
			if len(strs) < 3 {
				return
			}
			var val string
			if err = entry.UnmarshalValue(&val); err != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case kc <- &SlotKeyEntry{
				Slot:  slot,
				Key:   strs[2],
				Value: val,
			}:
			}
		}
	}()

	return kc, nil
}
