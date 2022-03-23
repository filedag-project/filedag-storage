package s3api

import (
	"bytes"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/testsign"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestStsAPIHandlers_AssumeRole(t *testing.T) {
	body := bytes.NewReader([]byte("Version=2011-06-15&Action=AssumeRole"))
	req := testsign.MustNewSignedV4Request(http.MethodPost, "http://127.0.0.1:9985/", 0, body, "sts", "test1", "testsecretKey", t)
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(err)
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	header := resp.Header
	fmt.Printf("resp%v,%v", string(all), header)
}
