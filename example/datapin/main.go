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
	pinUrl   = "/pin"
	unpinUrl = "/unpin"
)

//go run -tags example main.go pin/unpin --addr=127.0.0.1:9985 --access-key=test --secret-key=test --bucket=aaa --object=as1234

func main() {
	var addr, accessKey, secretKey, bucket, object string
	f := flag.NewFlagSet("datapin", flag.ExitOnError)
	f.StringVar(&addr, "addr", "", "the addr of server eg.127.0.0.1:9985")
	f.StringVar(&accessKey, "access-key", "", "the access-key which have adduser policy")
	f.StringVar(&secretKey, "secret-key", "", "the secret-key which have adduser policy")
	f.StringVar(&bucket, "bucket", "", "the bucket cannot be empty")
	f.StringVar(&object, "object", "", "the object cannot be empty")
	switch os.Args[1] {
	case "pin":
		f.Parse(os.Args[2:])
		err := pin(addr, accessKey, secretKey, pinUrl, bucket, object)
		if err != nil {
			fmt.Printf("pin data err %v", err)
			return
		}
	case "unpin":
		f.Parse(os.Args[2:])
		err := pin(addr, accessKey, secretKey, unpinUrl, bucket, object)
		if err != nil {
			fmt.Printf("unpin data err %v", err)
			return
		}
	default:
		fmt.Println("expected 'adduser' subcommands")
		os.Exit(1)
	}
}
func pin(addr, accessKey, secretKey, url, bucket, object string) error {
	fmt.Println("addr:", addr)
	fmt.Println("accessKey:", accessKey)
	fmt.Println("secretKey:", secretKey)
	fmt.Println("bucket:", bucket)
	fmt.Println("object:", object)
	if addr == "" || accessKey == "" || secretKey == "" || bucket == "" || object == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go pin/unpin --addr= --access-key= --secret-key= --bucket= --object=")
		return xerrors.Errorf("check your input")
	}
	err := exampleutils.SendSignedV4Request(http.MethodPost, "http://"+addr+url+"/"+bucket+"/"+object,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return err
	}
	return nil
}
