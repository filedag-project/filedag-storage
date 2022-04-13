package blockstore

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
)

var testdata = [][]byte{
	[]byte("Variety is the spice of life"),
	[]byte("Bad times make a good man"),
	[]byte("There is no royal road to learning"),
	[]byte("Doubt is the key to knowledge"),
	[]byte("The greatest test of courage on earth is to bear defeat without losing heart"),
	[]byte("A man's best friends are his ten fingers"),
	[]byte("Only they who fulfill their duties in everyday matters will fulfill them on great occasions"),
	[]byte("The shortest way to do many things is to only one thing at a time"),
}
var blockdatas []blocks.Block

func init() {
	//cb := cid.V1Builder{Codec: cid.DagCBOR, MhType: mh.SHA2_256}
	cb := cid.V0Builder{}
	for _, d := range testdata {
		id, _ := cb.Sum(d)
		fmt.Println(id)
		b, _ := blocks.NewBlockWithCid(d, id)
		blockdatas = append(blockdatas, b)
	}
}
func TestMutcaskbs(t *testing.T) {
	bstore, err := NewMutcaskbs(&Config{
		Path: tmpdirpath(t),
	})
	if err != nil {
		t.Fatal("failed to init mutcaskbs")
	}
	//ctx := context.TODO()
	// test put block data
	for _, blo := range blockdatas {
		if err := bstore.Put(blo); err != nil {
			t.Fatal(fmt.Sprintf("Put failed: %s", err))
		}
	}
	// test get block data
	for _, blo := range blockdatas {
		b, err := bstore.Get(blo.Cid())
		if err != nil {
			t.Fatal(fmt.Sprintf("Get failed: %s", err))
		}
		if !bytes.Equal(b.RawData(), blo.RawData()) {
			t.Fatal("data not match")
		}
	}
}

func tmpdirpath(t *testing.T) string {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to make temp dir: %s", err)
	}
	return tmpdir
}
