package uleveldb

import (
	"fmt"
	"testing"
)

func TestULeveldb(t *testing.T) {
	err := GlobalLevelDB.Put("a", 10)
	if err != nil {
		return
	}
	var a int
	err = GlobalLevelDB.Get("a", &a)
	if err != nil {
		return
	}
	fmt.Println(a)
}
