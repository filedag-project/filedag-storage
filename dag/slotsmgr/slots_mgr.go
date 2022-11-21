package slotsmgr

import (
	"errors"
	"fmt"
)

const ClusterSlots = 16384

var IndexRangeError = errors.New("end index should be equal or greater than start index")

// SlotPair means a range of [start, end]
type SlotPair struct {
	Start uint64
	End   uint64
}

func (sp SlotPair) String() string {
	if sp.Start == sp.End {
		return fmt.Sprintf("%v", sp.Start)
	}
	return fmt.Sprintf("%v-%v", sp.Start, sp.End)
}

func (sp SlotPair) Count() uint64 {
	return sp.End - sp.Start + 1
}

type SlotsManager struct {
	bitset *BitSet
}

func NewSlotsManager() *SlotsManager {
	return &SlotsManager{
		bitset: NewBitSet(ClusterSlots),
	}
}

func (hs *SlotsManager) Get(index uint64) (bool, error) {
	return hs.bitset.Get(index)
}

func (hs *SlotsManager) Set(index uint64, value bool) (bool, error) {
	return hs.bitset.Set(index, value)
}

func (hs *SlotsManager) Count() uint64 {
	return hs.bitset.Count()
}

func (hs *SlotsManager) SetRange(pair SlotPair, value bool) error {
	if pair.Start > pair.End {
		return IndexRangeError
	}
	for i := pair.Start; i <= pair.End; i++ {
		if _, err := hs.bitset.Set(i, value); err != nil {
			return err
		}
	}
	return nil
}

func (hs *SlotsManager) ToSlotPair() []SlotPair {
	var slotPairs []SlotPair
	pair := SlotPair{}
	isSetStart := false
	for i, n := range hs.bitset.data {
		if !isSetStart && n == 0 {
			continue
		}
		if isSetStart && n == 0xffffffffffffffff {
			continue
		}
		start := i << shift
		end := (i + 1) << shift
		for ; start < end; start++ {
			if n&(1<<uint(start&mask)) != 0 {
				if !isSetStart {
					isSetStart = true
					pair.Start = uint64(start)
				}
			} else {
				if isSetStart {
					pair.End = uint64(start - 1)
					slotPairs = append(slotPairs, pair)
					pair = SlotPair{}
					isSetStart = false
				}
			}
		}
	}
	return slotPairs
}
