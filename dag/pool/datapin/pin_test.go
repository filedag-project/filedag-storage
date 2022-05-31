package datapin

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/datapin/types"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/ipfs/go-cid"
	"testing"
)

func TestPin_test(t *testing.T) {
	c, err := cid.Decode("QmP63DkAFEnDYNjDYBpyNDfttu1fvUw99x1brscPzpqmmq")
	pin := Pin{
		Cid:  types.Cid{Cid: c},
		Name: "testname",
		Meta: map[string]string{
			"meta": "data",
		},
	}
	db, _ := uleveldb.OpenDb(utils.TmpDirPath(&testing.T{}))
	blockPin, err := NewBlockPin(db)
	pinSer := PinService{
		blockPin: blockPin,
	}
	err = pinSer.AddPin(pin)
	if err != nil {
		fmt.Println(err)
	}
	err = pinSer.RemovePin(types.Cid{Cid: c})
	if err != nil {
		fmt.Println(err)
	}
}
