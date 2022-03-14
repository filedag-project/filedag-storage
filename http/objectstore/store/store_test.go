package store

import (
	"fmt"
	"os"
	"testing"
)

func TestPutFile(t *testing.T) {
	file, err := os.Open("./test.txt")
	if err != nil {
		return
	}
	cid, err := PutFile(".", "aa.txt", file)
	if err != nil {
		return
	}
	fmt.Println(cid)
}
