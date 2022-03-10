package main

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"os"
)

const (
	dbFILE = "/tmp/leveldb2/fds.db"
)

func main() {
	err := os.Setenv("DBPATH", dbFILE)
	if err != nil {
		return
	}
	uleveldb.DBClient = uleveldb.OpenDb(os.Getenv("DBPATH"))
	startServer()
}
