package mutcask

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"sync"

	"github.com/filedag-project/filedag-storage/kv"
)

const MaxKeySize = 128

// max key size 128 byte +  1 byte which record the key size + 1 byte delete flag
const HintKeySize = MaxKeySize + 1 + 1

// HintKeySize + 8 bytes value offset + 4 bytes value size
const HintEncodeSize = HintKeySize + 8 + 4

const (
	HintDeletedFlag = byte(1)
)

type Hint struct {
	Key     string
	KOffset uint64
	VOffset uint64
	VSize   uint32
	Deleted bool
}

/**
		key		:	value offset	:	value size
		128+1   :   8   			:   4
**/
func (h *Hint) Encode() (ret []byte, err error) {
	kl := len(h.Key)
	if kl > MaxKeySize {
		return nil, ErrKeySizeTooLong
	}
	ret = make([]byte, HintEncodeSize)
	if h.Deleted {
		ret[0] = HintDeletedFlag
	}
	ret[1] = uint8(kl)
	copy(ret[2:HintKeySize], []byte(h.Key))
	binary.LittleEndian.PutUint64(ret[HintKeySize:HintKeySize+8], h.VOffset)
	binary.LittleEndian.PutUint32(ret[HintKeySize+8:], h.VSize)
	return
}

func (h *Hint) From(buf []byte) (err error) {
	if len(buf) != HintEncodeSize {
		return ErrHintFormat
	}
	if buf[0] == HintDeletedFlag {
		h.Deleted = true
	}
	keylen := uint8(buf[1])
	key := make([]byte, keylen)
	copy(key, buf[2:2+keylen])
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

const (
	opread = iota
	opwrite
	opdelete
)

type action struct {
	optype   int
	hint     *Hint
	key      string
	value    []byte
	retvchan chan retv
}

type retv struct {
	err  error
	data []byte
}

type Cask struct {
	id          uint32
	close       func()
	closeChan   chan struct{}
	actChan     chan *action
	vLog        *os.File
	vLogSize    uint64
	hintLog     *os.File
	hintLogSize uint64
	keyMap      *KeyMap
}

func NewCask(id uint32) *Cask {
	cc := make(chan struct{})
	cask := &Cask{
		id:        id,
		closeChan: cc,
		actChan:   make(chan *action),
	}
	cask.keyMap = &KeyMap{
		m: make(map[string]*Hint),
	}
	var once sync.Once
	cask.close = func() {
		once.Do(func() {
			close(cc)
		})
	}
	go func(cask *Cask) {
		for {
			select {
			case <-cask.closeChan:
				return
			case act := <-cask.actChan:
				switch act.optype {
				case opread:
					cask.doread(act)
				case opdelete:
					cask.dodelete(act)
				case opwrite:
					cask.dowrite(act)
				default:
					fmt.Printf("unkown op type %d\n", act.optype)
				}
			}

		}

	}(cask)
	return cask
}

func (c *Cask) Close() {
	c.close()
	if c.hintLog != nil {
		c.hintLog.Close()
	}
	if c.vLog != nil {
		c.vLog.Close()
	}
}

func (c *Cask) Put(key string, value []byte) (err error) {
	retvc := make(chan retv)
	c.actChan <- &action{
		optype:   opwrite,
		key:      key,
		value:    value,
		retvchan: retvc,
	}
	ret := <-retvc

	return ret.err
}

func (c *Cask) Delete(key string) (err error) {
	hint, has := c.keyMap.Get(key)
	if !has || hint.Deleted {
		return nil
	}
	retvc := make(chan retv)
	c.actChan <- &action{
		optype:   opdelete,
		key:      key,
		hint:     hint,
		retvchan: retvc,
	}
	ret := <-retvc

	return ret.err
}

func (c *Cask) Read(key string) (v []byte, err error) {
	hint, has := c.keyMap.Get(key)
	if !has || hint.Deleted {
		return nil, kv.ErrNotFound
	}

	retvc := make(chan retv)
	c.actChan <- &action{
		optype:   opread,
		hint:     hint,
		retvchan: retvc,
	}
	ret := <-retvc
	if ret.err != nil {
		return nil, ret.err
	}
	return ret.data, nil

}

func (c *Cask) Size(key string) (int, error) {
	hint, has := c.keyMap.Get(key)
	if !has || hint.Deleted {
		return -1, kv.ErrNotFound
	}
	return int(hint.VSize - 4), nil
}

func (c *Cask) doread(act *action) {
	var err error
	defer func() {
		if err != nil {
			act.retvchan <- retv{err: err}
		}
	}()
	buf := make([]byte, act.hint.VSize)
	_, err = c.vLog.ReadAt(buf, int64(act.hint.VOffset))
	if err != nil {
		return
	}
	v, err := DecodeValue(buf, true)
	if err != nil {
		return
	}
	act.retvchan <- retv{data: v}
}

func (c *Cask) dodelete(act *action) {
	var err error
	defer func() {
		if err != nil {
			act.retvchan <- retv{err: err}
		}
	}()
	// operations for one cask actually did in a sync style, so there is no need to use actomic
	fsize := c.hintLogSize // atomic.LoadUint64(&c.hintLogSize)
	//fmt.Printf("%d | %s hint offset: %d, %d, file size: %d\n", c.id, act.key, act.hint.KOffset, act.hint.KOffset+HintEncodeSize, fsize)
	if act.hint.KOffset+HintEncodeSize > fsize {
		err = ErrReadHintBeyondRange
		return
	}
	_, err = c.hintLog.WriteAt([]byte{HintDeletedFlag}, int64(act.hint.KOffset))
	if err != nil {
		return
	}
	// h := &Hint{
	// 	Deleted: true,
	// 	Key:     act.hint.Key,
	// 	KOffset: act.hint.KOffset,
	// 	VOffset: act.hint.VOffset,
	// 	VSize:   act.hint.VSize,
	// }
	act.hint.Deleted = true
	c.keyMap.Add(act.key, act.hint)
	// truncate the last hint
	act.retvchan <- retv{}
}

func (c *Cask) dowrite(act *action) {
	var err error
	defer func() {
		if err != nil {
			act.retvchan <- retv{err: err}
		}
	}()

	var hint = &Hint{}
	var isAddNew bool
	// check if key value already been saved
	if h, has := c.keyMap.Get(act.key); has {
		hint = h
		// the crc code
		buf := make([]byte, 4)
		_, err = c.vLog.ReadAt(buf, int64(h.VOffset))
		if err != nil {
			return
		}

		crcRecord := binary.LittleEndian.Uint32(buf)
		crcv := crc32.ChecksumIEEE(act.value)
		// value has same crc code
		if crcRecord == crcv && !h.Deleted {
			act.retvchan <- retv{}
			return
		}
		// value has same crc code but in deleted state
		// should just update hint to undeleted state
		if crcRecord == crcv && h.Deleted {
			if hint.KOffset+HintEncodeSize > c.hintLogSize {
				err = ErrReadHintBeyondRange
				return
			}
			_, err = c.hintLog.WriteAt([]byte{HintDeletedFlag}, int64(hint.KOffset))
			if err != nil {
				return
			}
			hint.Deleted = false
			c.keyMap.Add(hint.Key, hint)
			return
		}
	} else {
		isAddNew = true
		hint.KOffset = c.hintLogSize
	}

	// record file size as value offset
	voffset := c.vLogSize
	// encode value
	encbytes := EncodeValue(act.value)
	// record encoded value size
	vsize := uint32(len(encbytes))
	// write to vlog file
	_, err = c.vLog.WriteAt(encbytes, int64(voffset))
	if err != nil {
		return
	}
	// operations for one cask actually did in a sync style, so there is no need to use actomic
	// update vlog file size
	//atomic.AddUint64(&c.vLogSize, uint64(vsize))
	c.vLogSize += uint64(vsize)

	hint.Key = act.key
	hint.VOffset = voffset
	hint.VSize = vsize
	hint.Deleted = false

	encHintBytes, err := hint.Encode()
	if err != nil {
		return
	}
	_, err = c.hintLog.WriteAt(encHintBytes, int64(hint.KOffset))
	if err != nil {
		return
	}

	if isAddNew {
		// update hint log file size
		// atomic.AddUint64(&c.hintLogSize, HintEncodeSize)
		c.hintLogSize += HintEncodeSize
	}

	fmt.Printf("update key map for %d\n", c.id)
	c.keyMap.Add(hint.Key, hint)
	act.retvchan <- retv{}
	fmt.Printf("put %s = %s\n", act.key, act.value)
}
