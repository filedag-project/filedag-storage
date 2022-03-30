package iamapi

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/testsign"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

const (
	DefaultTestAccessKey = "test"
	DefaultTestSecretKey = "test"
)

func TestIamApiServer_GetUserList(t *testing.T) {
	u := "http://127.0.0.1:9985/admin/v1/list-user"
	req := testsign.MustNewSignedV4Request(http.MethodGet, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)

	//req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
func TestIamApiServer_AddUser(t *testing.T) {
	urlValues := make(url.Values)
	urlValues.Set("accessKey", "test2")
	urlValues.Set("secretKey", "test12345")
	u := "http://127.0.0.1:9985/admin/v1/add-user"
	req := testsign.MustNewSignedV4Request(http.MethodPost, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)

	//req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
func TestIamApiServer_GetUserInfo(t *testing.T) {
	urlValues := make(url.Values)
	user := "test1"
	urlValues.Set("userName", user)
	u := "http://127.0.0.1:9985/admin/v1/user-info"
	req := testsign.MustNewSignedV4Request(http.MethodGet, u+"?"+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)

	//req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}
