package uleveldb

import (
	"fmt"
	"testing"
)

func TestULeveldb(t *testing.T) {
	err := GlobalUserLevelDB.Put("a", 10)
	if err != nil {
		return
	}
	var a int
	err = GlobalUserLevelDB.Get("a", &a)
	if err != nil {
		return
	}
	fmt.Println(a)
}
