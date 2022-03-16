package s3api

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/testsign"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestS3ApiServer_PutBucketPolicyHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test"
	urlValues := make(url.Values)
	policy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::bucket/*"}]}`
	urlValues.Set("policy", policy)
	req := testsign.MustNewSignedV4Request(http.MethodPut, u+"?"+urlValues.Encode(), int64(len(policy)), strings.NewReader(policy), t)
	//req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	req.ContentLength = int64(len(policy))
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
func TestS3ApiServer_GetBucketPolicyHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test"
	req := testsign.MustNewSignedV4Request(http.MethodGet, u+"?policy", 0, nil, t)

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
func TestS3ApiServer_DeleteBucketPolicyHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test"
	req := testsign.MustNewSignedV4Request(http.MethodDelete, u+"?policy", 0, nil, t)

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
func TestS3ApiServer_HeadBucketHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test"
	//req.Header.Set("Content-Type", "text/plain")
	req := testsign.MustNewSignedV4Request(http.MethodHead, u, 0, nil, t)
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
func TestS3ApiServer_DeleteBucketHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test"
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

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
