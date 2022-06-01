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
	setStatusUrl = "/admin/v1/update-accessKey_status"
)

//go run -tags example main.go run set-status --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg --status=off

func main() {
	var addr, accessKey, secretKey, username, status string
	f := flag.NewFlagSet("set-status", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of server eg.127.0.0.1:9985")
	f.StringVar(&accessKey, "access-key", "", "the access-key which have set-status policy")
	f.StringVar(&secretKey, "secret-key", "", "the secret-key which have set-status policy")
	f.StringVar(&username, "username", "", "the username that you want set-status")
	f.StringVar(&status, "status", "", "the status that you want set")
	switch os.Args[1] {
	case "set-status":
		f.Parse(os.Args[2:])
		err := set(addr, accessKey, secretKey, username, status)
		if err != nil {
			fmt.Printf("set status err %v", err)
			return
		}
	default:
		fmt.Println("expected 'set-status' subcommands")
		os.Exit(1)
	}
}
func set(addr, accessKey, secretKey, username, status string) error {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" || status == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go set-status --addr= --access-key= --secret-key= --username= --status=")
		return xerrors.Errorf("check your input")
	}
	err := exampleutils.SendSignedV4Request(http.MethodPost, "http://"+addr+setStatusUrl+"?accessKey="+username+"&status="+status,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return err
	}
	return nil
}
