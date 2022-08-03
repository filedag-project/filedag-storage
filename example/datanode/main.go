package main

import (
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/node/datanode"
	"os"
)

//go run -tags example main.go daemon --host=127.0.0.1 --port=9010 --path=/tmp/dag/data
func main() {
	var host, port, path string
	f := flag.NewFlagSet("daemon", flag.ExitOnError)
	f.StringVar(&host, "host", "127.0.0.1", "set host eg.:127.0.0.1")
	f.StringVar(&port, "port", "9010", "set port eg.:9010")
	f.StringVar(&path, "path", "/tmp/dag/data", "set data node path")

	switch os.Args[1] {
	case "daemon":
		f.Parse(os.Args[2:])
		if host == "" || port == "" || path == "" {
			fmt.Printf("host:%v, port:%v, path:%v", host, port, path)
			fmt.Println("please check your input\n " +
				"USAGE ERROR: go run -tags example main.go daemon --host= --port= --path= ")
		} else {
			run(host, port, path)
		}
	default:
		fmt.Println("expected 'daemon' subcommands")
		os.Exit(1)
	}
}
func run(host, port, path string) {
	datanode.MutDataNodeServer(fmt.Sprintf("%s:%s", host, port), datanode.KVBadge, path)
}
