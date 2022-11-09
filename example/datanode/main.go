package main

import (
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/node/datanode"
	"os"
)

//go run -tags example main.go daemon --listen=127.0.0.1:9011 --datadir=/tmp/dn-data1 --kvdb=badger
func main() {
	var listen, kvdb, datadir string
	f := flag.NewFlagSet("daemon", flag.ExitOnError)
	f.StringVar(&listen, "listen", "127.0.0.1:9010", "set server listen")
	f.StringVar(&kvdb, "kvdb", "badger", "choose kvdb, badger or mutcask")
	f.StringVar(&datadir, "datadir", "/tmp/dag/data", "directory to store data in")

	switch os.Args[1] {
	case "daemon":
		f.Parse(os.Args[2:])
		if listen == "" || kvdb == "" || datadir == "" {
			fmt.Printf("listen:%v, kvdb:%v, datadir:%v", listen, kvdb, datadir)
			fmt.Println("please check your input\n " +
				"USAGE ERROR: go run -tags example main.go daemon --listen=127.0.0.1:9011 --datadir=/tmp/dn-data1 --kvdb=badger")
		} else {
			run(listen, kvdb, datadir)
		}
	default:
		fmt.Println("expected 'daemon' subcommands")
		os.Exit(1)
	}
}
func run(host, port, path string) {
	datanode.StartDataNodeServer(fmt.Sprintf("%s:%s", host, port), datanode.KVBadge, path)
}
