package uleveldb

import (
	"fmt"
	"testing"
)

func TestULeveldb(t *testing.T) {
	db, err := OpenDb(t.TempDir())
	err = db.Put("a", 10)
	if err != nil {
		return
	}
	var a int
	err = db.Get("a", &a)
	db.Close()
	if err != nil {
		return
	}
	fmt.Println(a)
}
