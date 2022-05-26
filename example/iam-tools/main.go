package main

import (
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	addUserUrl        = "/admin/v1/add-user"
	delUserUrl        = "/admin/v1/remove-user"
	getUserUrl        = "/admin/v1/user-info"
	changePassUserUrl = "/admin/v1/change-password"
	setStatusUrl      = "/admin/v1/update-accessKey_status"
)

//go run -tags example main.go --method=adduser --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg --pass=wpg123456
//go run -tags example main.go --method=deluser --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg
//go run -tags example main.go --method=getuser --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg
//go run -tags example main.go --method=change-pass --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg
//go run -tags example main.go --method=set-status --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg

func main() {
	var method, addr, accessKey, secretKey, username, pass, newPass, status string

	flag.StringVar(&method, "method", "", "the method you want use")
	flag.StringVar(&addr, "addr", "", "the addr of server eg.127.0.0.1:9985")
	flag.StringVar(&accessKey, "access-key", "", "the access-key which have adduser policy")
	flag.StringVar(&secretKey, "secret-key", "", "the secret-key which have adduser policy")
	flag.StringVar(&username, "username", "", "the username that you want add")
	flag.StringVar(&pass, "pass", "", "the password that you want add")
	flag.StringVar(&newPass, "new-pass", "", "the username that you want add")
	flag.StringVar(&status, "status", "", "the username that you want add")

	flag.Parse()
	run(method, addr, accessKey, secretKey, username, pass, newPass, status)
}
func run(method, addr, accessKey, secretKey, username, pass, newPass, status string) {
	var req *http.Request
	var err error
	switch method {
	case "adduser":
		req, err = add(addr, accessKey, secretKey, username, pass)
		if err != nil {
			return
		}
	case "deluser":
		req, err = del(addr, accessKey, secretKey, username)
		if err != nil {
			return
		}
	case "getuser":
		req, err = get(addr, accessKey, secretKey, username)
		if err != nil {
			return
		}
	case "change-pass":
		req, err = cha(addr, accessKey, secretKey, username, newPass)
		if err != nil {
			return
		}
	case "set-status":
		req, err = set(addr, accessKey, secretKey, username, status)
		if err != nil {
			return
		}
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
func add(addr, accessKey, secretKey, username, pass string) (*http.Request, error) {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" || pass == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go --method=adduser --addr= --access-key= --secret-key= --username= ")
		return nil, xerrors.Errorf("check your input")
	}
	return mustNewSignedV4Request(http.MethodPost, "http://"+addr+addUserUrl+"?accessKey="+username+"&secretKey="+pass,
		0, nil, "s3", accessKey, secretKey)
}
func del(addr, accessKey, secretKey, username string) (*http.Request, error) {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go --method=deluser --addr= --access-key= --secret-key= --username= ")
		return nil, xerrors.Errorf("check your input")
	}
	return mustNewSignedV4Request(http.MethodPost, "http://"+addr+delUserUrl+"?accessKey="+username,
		0, nil, "s3", accessKey, secretKey)
}
func get(addr, accessKey, secretKey, username string) (*http.Request, error) {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go --method=getuser --addr= --access-key= --secret-key= --username= ")
		return nil, xerrors.Errorf("check your input")
	}
	return mustNewSignedV4Request(http.MethodGet, "http://"+addr+getUserUrl+"?accessKey="+username,
		0, nil, "s3", accessKey, secretKey)
}
func cha(addr, accessKey, secretKey, username, pass string) (*http.Request, error) {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go --method=change-pass  --addr= --access-key= --secret-key= --username= getuser")
		return nil, xerrors.Errorf("check your input")
	}
	return mustNewSignedV4Request(http.MethodPost, "http://"+addr+changePassUserUrl+"?newPassword="+pass+"&accessKey="+username,
		0, nil, "s3", accessKey, secretKey)
}
func set(addr, accessKey, secretKey, username, status string) (*http.Request, error) {
	if addr == "" || accessKey == "" || secretKey == "" || username == "" {
		fmt.Println("please check your input\n " +
			"USAGE ERROR: go run -tags example main.go --method=set-status --addr= --access-key= --secret-key= --username= getuser")
		return nil, xerrors.Errorf("check your input")
	}
	return mustNewSignedV4Request(http.MethodPost, "http://"+addr+setStatusUrl+"?accessKey="+username+"&status="+status,
		0, nil, "s3", accessKey, secretKey)
}

//mustNewSignedV4Request  NewSignedV4Request
func mustNewSignedV4Request(method string, urlStr string, contentLength int64, body io.ReadSeeker, st string, accessKey, secretKey string) (*http.Request, error) {
	req, err := utils.NewRequest(method, urlStr, contentLength, body)
	if err != nil {
		return nil, err
	}
	cred := &auth.Credentials{AccessKey: accessKey, SecretKey: secretKey}
	if err = utils.SignRequestV4(req, cred.AccessKey, cred.SecretKey, utils.ServiceType(st)); err != nil {
		return nil, err
	}
	return req, nil
}
