package datapin

import (
	"github.com/filedag-project/filedag-storage/dag/pool/datapin/types"
)

// Pin contains basic information about a Pin and pinning options.
type Pin struct {
	Cid     types.Cid         `json:"cid"`
	Name    PinName           `json:"name"`
	Origins []types.Multiaddr `json:"origins"`
	Meta    map[string]string `json:"meta"`
}

// PinName is a string limited to 255 chars when serializing JSON.
type PinName string
