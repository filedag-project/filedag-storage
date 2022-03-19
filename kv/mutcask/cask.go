package mutcask

import (
	"encoding/binary"
	"hash/crc32"
	"os"
)

const MaxKeySize = 128

// max key size 128 byte +  1 byte which record the key size
const HintKeySize = MaxKeySize + 1

// HintKeySize + 8 bytes value offset + 4 bytes value size
const HintEncodeSize = HintKeySize + 8 + 4

type Hint struct {
	Key     string
	KOffset uint64
	VOffset uint64
	VSize   uint32
}

/**
		key		:	value offset	:	value size
		128+1   :   8   			:   4
**/
func (h *Hint) Encode() (ret []byte, err error) {
	kl := len(h.Key)
	if kl > 128 {
		return nil, ErrKeySizeTooLong
	}
	ret = make([]byte, HintEncodeSize)
	ret[0] = uint8(kl)
	copy(ret[1:HintKeySize], []byte(h.Key))
	binary.LittleEndian.PutUint64(ret[HintKeySize:HintKeySize+8], h.VOffset)
	binary.LittleEndian.PutUint32(ret[HintKeySize+8:], h.VSize)
	return
}

func (h *Hint) From(buf []byte) (err error) {
	if len(buf) != HintEncodeSize {
		return ErrHintFormat
	}
	keylen := uint8(buf[0])
	key := make([]byte, keylen)
	copy(key, buf[1:1+keylen])
	h.Key = string(key)
	h.VOffset = binary.LittleEndian.Uint64(buf[HintKeySize : HintKeySize+8])
	h.VSize = binary.LittleEndian.Uint32(buf[HintKeySize+8:])
	return
}

/**
		crc32	:	value
		4 		: 	xxxx
**/
func EncodeValue(v []byte) (ret []byte) {
	ret = make([]byte, 4+len(v))
	c32 := crc32.ChecksumIEEE(v)
	binary.LittleEndian.PutUint32(ret[0:4], c32)
	copy(ret[4:], v)
	return
}

func DecodeValue(buf []byte, verify bool) (v []byte, err error) {
	if len(buf) <= 4 {
		return nil, ErrValueFormat
	}
	if verify {
		vcheck := binary.LittleEndian.Uint32(buf[:4])

		c32 := crc32.ChecksumIEEE(buf[4:])
		// make sure data not rotted
		if vcheck != c32 {
			return nil, ErrDataRotted
		}
	}
	v = make([]byte, len(buf)-4)
	copy(v, buf[4:])
	return
}

type Cask struct {
	lock    chan struct{}
	vLog    *os.File
	hintLog *os.File
}
