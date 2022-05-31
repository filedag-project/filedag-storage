package types

import (
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"strings"
	"time"
)

var unixZero = time.Unix(0, 0)

type PinType uint64

// PinType values. See PinType documentation for further explanation.
const (
	// BadType type showing up anywhere indicates a bug
	BadType PinType = 1 << iota
	// DataType is a regular, non-sharded pin. It is pinned recursively.
	// It has no associated reference.
	DataType
)

// PinDepth indicates how deep a pin should be pinned, with
// -1 meaning "to the bottom", or "recursive".
type PinDepth int

// PinMode is a PinOption that indicates how to pin something on IPFS,
// recursively or direct.
type PinMode int

// PinMode values
const (
	PinModeRecursive PinMode = 0
	PinModeDirect    PinMode = 1
)

// Cid embeds a cid.Cid with the MarshalJSON/UnmarshalJSON methods overwritten.
type Cid struct {
	cid.Cid
}

type Pin struct {
	PinOptions
	Cid Cid `json:"cid" codec:"c"`

	// See PinType comments
	Type PinType `json:"type" codec:"t,omitempty"`

	// The peers to which this pin is allocated
	Allocations []peer.ID `json:"allocations" codec:"a,omitempty"`

	// MaxDepth associated to this pin. -1 means
	// recursive.
	MaxDepth PinDepth `json:"max_depth" codec:"d,omitempty"`

	// We carry a reference CID to this pin. For
	// ClusterDAGs, it is the MetaPin CID. For the
	// MetaPin it is the ClusterDAG CID. For Shards,
	// it is the previous shard CID.
	// When not needed the pointer is nil
	Reference *Cid `json:"reference" codec:"r,omitempty"`

	// The time that the pin was submitted to the consensus layer.
	Timestamp time.Time `json:"timestamp" codec:"i,omitempty"`
}

// PinOptions wraps user-defined options for Pins
type PinOptions struct {
	Name      string            `json:"name" codec:"n,omitempty"`
	Mode      PinMode           `json:"mode" codec:"o,omitempty"`
	ShardSize uint64            `json:"shard_size" codec:"s,omitempty"`
	Metadata  map[string]string `json:"metadata" codec:"m,omitempty"`
}

// Multiaddr is a concrete type to wrap a Multiaddress so that it knows how to
// serialize and deserialize itself.
type Multiaddr struct {
	multiaddr.Multiaddr
}

// String is a string representation of a Pin.
func (pin Pin) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "cid: %s\n", pin.Cid.String())
	fmt.Fprintf(&b, "type: %s\n", pin.Type)
	fmt.Fprintf(&b, "allocations: %v\n", pin.Allocations)
	fmt.Fprintf(&b, "maxdepth: %d\n", pin.MaxDepth)
	if pin.Reference != nil {
		fmt.Fprintf(&b, "reference: %s\n", pin.Reference)
	}
	return b.String()
}

// PinCid is a shortcut to create a Pin only with a Cid.  Default is for pin to
// be recursive and the pin to be of DataType.
func PinCid(c Cid) Pin {
	return Pin{
		Cid:         c,
		Type:        DataType,
		Allocations: []peer.ID{},
		MaxDepth:    -1, // Recursive
		Timestamp:   time.Now(),
	}
}

// PinWithOpts creates a new Pin calling PinCid(c) and then sets its
// PinOptions fields with the given options. Pin fields that are linked to
// options are set accordingly (MaxDepth from Mode).
func PinWithOpts(c Cid, opts PinOptions) Pin {
	p := PinCid(c)
	p.PinOptions = opts
	return p
}

// Defined returns true if this is not a zero-object pin (the CID must be set).
func (pin Pin) Defined() bool {
	return pin.Cid.Defined()
}
