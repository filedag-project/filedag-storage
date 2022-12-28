package slotmigraterepo

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"github.com/syndtr/goleveldb/leveldb"
	"strconv"
	"strings"
)

const SlotMigratePrefix = "migrate/"

// SlotMigrateRepo saves information about the slot to be transferred.
type SlotMigrateRepo struct {
	db *uleveldb.ULevelDB
}

func NewSlotMigrateRepo(db *uleveldb.ULevelDB) *SlotMigrateRepo {
	return &SlotMigrateRepo{db: db}
}

func (s *SlotMigrateRepo) Set(slot uint16, value string) error {
	return s.db.Put(fmt.Sprintf("%s%v", SlotMigratePrefix, slot), value)
}

func (s *SlotMigrateRepo) Get(slot uint16) (value string, err error) {
	err = s.db.Get(fmt.Sprintf("%s%v", SlotMigratePrefix, slot), &value)
	return value, err
}

func (s *SlotMigrateRepo) Has(slot uint16) (bool, error) {
	var val string
	err := s.db.Get(fmt.Sprintf("%s%v", SlotMigratePrefix, slot), &val)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *SlotMigrateRepo) Remove(slot uint16) error {
	return s.db.Delete(fmt.Sprintf("%s%v", SlotMigratePrefix, slot))
}

type Entry struct {
	Slot  uint16
	Value string
}

func (s *SlotMigrateRepo) AllKeysChan(ctx context.Context) (<-chan *Entry, error) {
	all, err := s.db.ReadAllChan(ctx, SlotMigratePrefix, "")
	if err != nil {
		return nil, err
	}
	kc := make(chan *Entry)
	go func() {
		defer close(kc)
		for entry := range all {
			strs := strings.Split(entry.Key, "/")
			if len(strs) < 2 {
				return
			}
			slot, err := strconv.ParseUint(strs[1], 10, 32)
			if err != nil {
				return
			}
			var val string
			if err = entry.UnmarshalValue(&val); err != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case kc <- &Entry{
				Slot:  uint16(slot),
				Value: val,
			}:
			}
		}
	}()

	return kc, nil
}
