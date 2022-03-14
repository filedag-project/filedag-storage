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
	policy := `{"Version":"2012-10-17","Statement":[{"Action":["s3:GetBucketLocation"],"Effect":"Allow","Principal":{"AWS":["*"]},"Sid":""}]}`
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
