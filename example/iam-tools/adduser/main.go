//go:build example
// +build example

package main

import (
	"flag"
	"fmt"
	tools "github.com/filedag-project/filedag-storage/example/iam-tools"
	"io/ioutil"
	"net/http"
)

const (
	addUserUrl = "/admin/v1/add-user"
)

func main() {
	var addr, accessKey, secretKey, username, pass string

	flag.StringVar(&addr, "addr", "", "the addr of server eg.127.0.0.1:9985")
	flag.StringVar(&accessKey, "access-key", "", "the access-key which have adduser policy")
	flag.StringVar(&secretKey, "secret-key", "", "the secret-key which have adduser policy")
	flag.StringVar(&username, "username", "", "the username that you want add")
	flag.StringVar(&pass, "pass", "", "the password that you want add")

	flag.Parse()
	if addr == "" || accessKey == "" || secretKey == "" || username == "" || pass == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go --addr= --access-key= secret-key= --username= --pass=")
		return
	}
	run(addr, accessKey, secretKey, username, pass)
}
func run(addr, accessKey, secretKey, username, pass string) {
	req, err := tools.MustNewSignedV4Request(http.MethodPost, "http://"+addr+addUserUrl+"?accessKey="+username+"&secretKey="+pass,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		fmt.Printf("mustNewSignedV4Request err: %v", err)
		return
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Do req err: %v\n", err)
		return
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read resp err: %v\n", err)
		return
	}
	fmt.Printf("response: %v\n", string(all))
}
