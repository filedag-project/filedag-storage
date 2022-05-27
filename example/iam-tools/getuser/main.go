package main

import (
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/example/iam-tools/exampleutils"
	"golang.org/x/xerrors"
	"net/http"
	"os"
)

const (
	getUserUrl = "/admin/v1/user-info"
)

//go run -tags example main.go getuser --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg

func main() {
	var addr, accessKey, secretKey, username string
	f := flag.NewFlagSet("getuser", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of server eg.127.0.0.1:9985")
	f.StringVar(&accessKey, "access-key", "", "the access-key which have getuser policy")
	f.StringVar(&secretKey, "secret-key", "", "the secret-key which have getuser policy")
	f.StringVar(&username, "username", "", "the username that you want get")

	switch os.Args[1] {
	case "getuser":
		f.Parse(os.Args[2:])
		err := get(addr, accessKey, secretKey, username)
		if err != nil {
			fmt.Printf("get user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'getuser' subcommands")
		os.Exit(1)
	}
}
func get(addr, accessKey, secretKey, username string) error {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go getuser --addr= --access-key= --secret-key= --username= ")
		return xerrors.Errorf("check your input")
	}
	err := exampleutils.SendSignedV4Request(http.MethodGet, "http://"+addr+getUserUrl+"?accessKey="+username,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return err
	}
	return nil
}