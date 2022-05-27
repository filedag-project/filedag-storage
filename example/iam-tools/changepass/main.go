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
	changePassUserUrl = "/admin/v1/change-password"
)

//go run -tags example main.go run change-pass --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg --new-pass=qwe123456

func main() {
	var addr, accessKey, secretKey, username, newPass string
	f := flag.NewFlagSet("change-pass", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of server eg.127.0.0.1:9985")
	f.StringVar(&accessKey, "access-key", "", "the access-key which have change-pass policy")
	f.StringVar(&secretKey, "secret-key", "", "the secret-key which have change-pass policy")
	f.StringVar(&username, "username", "", "the username that you want change pass")
	f.StringVar(&newPass, "new-pass", "", "the pass that you want change")
	switch os.Args[1] {
	case "change-pass":
		f.Parse(os.Args[2:])
		err := cha(addr, accessKey, secretKey, username, newPass)
		if err != nil {
			fmt.Printf("change pass err %v", err)
			return
		}
	default:
		fmt.Println("expected 'change-pass' subcommands")
		os.Exit(1)
	}
}
func cha(addr, accessKey, secretKey, username, pass string) error {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" || pass == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go change-pass  --addr= --access-key= --secret-key= --username= --new-pass=")
		return xerrors.Errorf("check your input")
	}
	err := exampleutils.SendSignedV4Request(http.MethodPost, "http://"+addr+changePassUserUrl+"?newPassword="+pass+"&accessKey="+username,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return err
	}
	return nil
}
