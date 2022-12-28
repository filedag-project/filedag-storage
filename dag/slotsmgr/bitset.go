package slotsmgr

import "errors"

const (
	shift = 6    // 2^6 = 64
	mask  = 0x3f // 63
)

var IndexOutOfBoundError = errors.New("index should be less than max size")

type BitSet struct {
	data []uint64 // 64 bits
	max  uint64
}

func NewBitSet(size uint64) *BitSet {
	return &BitSet{
		data: make([]uint64, size>>shift+1),
		max:  size,
	}
}

func (bs *BitSet) Get(index uint64) (bool, error) {
	if index >= bs.max {
		return false, IndexOutOfBoundError
	}
	idx := index >> shift
	return bs.data[idx]&(1<<uint(index&mask)) != 0, nil
}

func (bs *BitSet) Set(index uint64, value bool) (oldValue bool, err error) {
	if index >= bs.max {
		return false, IndexOutOfBoundError
	}
	idx := index >> shift
	oldValue, err = bs.Get(index)
	if err != nil {
		return false, err
	}
	if !oldValue && value {
		bs.data[idx] |= 1 << uint(index&mask) // The corresponding bit is set to 1
	} else if oldValue && !value {
		bs.data[idx] &^= 1 << uint(index&mask) // The corresponding bit is set to 0
	}
	return oldValue, nil
}

func (bs *BitSet) Count() uint64 {
	var count uint64
	for _, b := range bs.data {
		count += swar(b)
	}
	return count
}

func swar(i uint64) uint64 {
	i = (i & 0x5555555555555555) + ((i >> 1) & 0x5555555555555555)
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
	i = (i & 0x0F0F0F0F0F0F0F0F) + ((i >> 4) & 0x0F0F0F0F0F0F0F0F)
	i = (i * 0x0101010101010101) >> 56
	return i
}
