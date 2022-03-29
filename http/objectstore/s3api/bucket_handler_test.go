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
	u := "http://127.0.0.1:9985/test22"
	urlValues := make(url.Values)
	policy := `{"Version":"2008-10-17","Id":"aaaa-bbbb-cccc-dddd","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::test22/*"}]}`
	urlValues.Set("policy", policy)
	req := testsign.MustNewSignedV4Request(http.MethodPut, u+"?"+urlValues.Encode(), int64(len(policy)), strings.NewReader(policy),
		"s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
	u := "http://127.0.0.1:9985/test22"
	req := testsign.MustNewSignedV4Request(http.MethodGet, u+"?policy", 0, nil, "s3",
		DefaultTestAccessKey, DefaultTestSecretKey, t)

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
	u := "http://127.0.0.1:9985/test22"
	req := testsign.MustNewSignedV4Request(http.MethodDelete, u+"?policy", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)

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
	u := "http://127.0.0.1:9985/test22"
	//req.Header.Set("Content-Type", "text/plain")
	req := testsign.MustNewSignedV4Request(http.MethodHead, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
func TestS3ApiServer_GetBucketLocationHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test22"
	//req.Header.Set("Content-Type", "text/plain")
	req := testsign.MustNewSignedV4Request(http.MethodHead, u+"?location", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
	u := "http://127.0.0.1:9985/test22"
	req := testsign.MustNewSignedV4Request(http.MethodDelete, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)

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
func TestS3ApiServer_PutBucketHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test22"

	req := testsign.MustNewSignedV4Request(http.MethodPut, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
func TestS3ApiServer_ListBucketHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/"

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

func TestS3ApiServer_GetBucketAclHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test"
	req := testsign.MustNewSignedV4Request(http.MethodGet, u+"?acl=", 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
func TestS3ApiServer_PutBucketAclHandler(t *testing.T) {
	u := "http://127.0.0.1:9985/test"
	a := `<?xml version="1.0" encoding="UTF-8"?>
<AccessControlPolicy xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
  <Owner>
    <ID>*** Owner-Canonical-User-ID ***</ID>
    <DisplayName>owner-display-name</DisplayName>
  </Owner>
  <AccessControlList>
    <Grant>
      <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
               xsi:type="Canonical User">
        <ID>*** Owner-Canonical-User-ID ***</ID>
        <DisplayName>display-name</DisplayName>
      </Grantee>
      <Permission>FULL_CONTROL</Permission>
    </Grant>
  </AccessControlList>
</AccessControlPolicy>`
	req := testsign.MustNewSignedV4Request(http.MethodPut, u+"?acl=", int64(len(a)), strings.NewReader(a), "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
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
