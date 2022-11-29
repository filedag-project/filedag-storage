package client

import (
	"bytes"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"testing"
)

func TestBalanceNode(t *testing.T) {
	cl, _ := NewMockPoolClient(t)
	ds := merkledag.NewDAGService(blockservice.New(cl, offline.Exchange(cl)))
	testcases := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "empty",
			data:     []byte{},
			expected: true,
		},
		{
			name:     "non-empty",
			data:     []byte("\b\u0002\u0012\a1234567\u0018\a"),
			expected: true,
		},
		{
			name: "big",
			data: bytes.Repeat([]byte("1234567890"), 1000000),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			node := merkledag.NodeWithData(tc.data)
			cidBuilder, _ := merkledag.PrefixForCidVersion(0)
			_, err := BalanceNode(bytes.NewReader(tc.data), ds, cidBuilder)
			if err != nil {
				t.Errorf("BalanceNode(%v) = false, expected true", node)
			}
		})
	}

}
