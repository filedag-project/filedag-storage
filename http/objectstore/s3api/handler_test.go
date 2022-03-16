package s3api

import (
	"fmt"
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
	req, err := http.NewRequest("PUT", u+"?"+urlValues.Encode(), strings.NewReader(policy))
	if err != nil {
		fmt.Println(err)
		return
	}
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
	urlValues := make(url.Values)
	policy := ""
	urlValues.Set("policy", policy)
	req, err := http.NewRequest("GET", u+"?policy", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
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
