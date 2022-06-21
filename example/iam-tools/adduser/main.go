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
	addUserUrl = "/admin/v1/add-user"
)

//go run -tags example main.go add-user --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg --pass=wpg123456

func main() {
	var addr, accessKey, secretKey, username, pass string
	f := flag.NewFlagSet("add-user", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of server eg.127.0.0.1:9985")
	f.StringVar(&accessKey, "access-key", "", "the access-key which have adduser policy")
	f.StringVar(&secretKey, "secret-key", "", "the secret-key which have adduser policy")
	f.StringVar(&username, "username", "", "the username that you want add")
	f.StringVar(&pass, "pass", "", "the user password that you want add")
	switch os.Args[1] {
	case "add-user":
		f.Parse(os.Args[2:])
		err := add(addr, accessKey, secretKey, username, pass)
		if err != nil {
			fmt.Printf("add user err %v", err)
			return
		}
	default:
		fmt.Println("expected 'add-user' subcommands")
		os.Exit(1)
	}
}
func add(addr, accessKey, secretKey, username, pass string) error {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" || pass == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go add-user --addr= --access-key= --secret-key= --username= --pass=")
		return xerrors.Errorf("check your input")
	}
	err := exampleutils.SendSignedV4Request(http.MethodPost, "http://"+addr+addUserUrl+"?accessKey="+username+"&secretKey="+pass,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return err
	}
	return nil
}
