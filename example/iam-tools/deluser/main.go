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
	delUserUrl = "/admin/v1/remove-user"
)

//go run -tags example main.go deluser --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg

func main() {
	var addr, accessKey, secretKey, username string
	f := flag.NewFlagSet("deluser", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of server eg.127.0.0.1:9985")
	f.StringVar(&accessKey, "access-key", "", "the access-key which have deluser policy")
	f.StringVar(&secretKey, "secret-key", "", "the secret-key which have deluser policy")
	f.StringVar(&username, "username", "", "the username that you want del")

	switch os.Args[1] {
	case "deluser":
		f.Parse(os.Args[2:])
		err := del(addr, accessKey, secretKey, username)
		if err != nil {
			fmt.Printf("del user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'deluser' subcommands")
		os.Exit(1)
	}
}
func del(addr, accessKey, secretKey, username string) error {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go deluser --addr= --access-key= --secret-key= --username= ")
		return xerrors.Errorf("check your input")
	}
	err := exampleutils.SendSignedV4Request(http.MethodPost, "http://"+addr+delUserUrl+"?accessKey="+username,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return err
	}
	return nil
}
