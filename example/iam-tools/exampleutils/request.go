package exampleutils

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"io"
	"io/ioutil"
	"net/http"
)

//SendSignedV4Request  NewSignedV4Request
func SendSignedV4Request(method string, urlStr string, contentLength int64, body io.ReadSeeker, st string, accessKey, secretKey string) error {
	req, err := utils.NewRequest(method, urlStr, contentLength, body)
	if err != nil {
		return err
	}
	cred := &auth.Credentials{AccessKey: accessKey, SecretKey: secretKey}
	if err = utils.SignRequestV4(req, cred.AccessKey, cred.SecretKey, utils.ServiceType(st)); err != nil {
		return err
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("response: %v\n", string(all))
	return nil
}
