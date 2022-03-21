package mutcask

import (
	"bytes"
	"testing"
)

func TestHintEncode(t *testing.T) {
	h1 := &Hint{
		Key:     "QmYs2ezGBk63nzf3vD4EHejWfN5ZkDfTVroS7rwY2JTbnQ",
		VOffset: 4 << 10,
		VSize:   512,
	}
	bs, err := h1.Encode()
	if err != nil {
		t.Fatal()
	}
	h2 := &Hint{}
	err = h2.From(bs)
	if err != nil {
		t.Fatal()
	}
	if h1.Key != h2.Key || h1.VOffset != h2.VOffset || h1.VSize != h2.VSize {
		t.Fatal()
	}
}

func TestValueEncodeDecode(t *testing.T) {
	value := []byte("mutation of bitcask")
	encoded := EncodeValue(value)
	v, err := DecodeValue(encoded, true)
	if err != nil {
		t.Fatal()
	}
	if !bytes.Equal(value, v) {
		t.Fatal()
	}
}
