package main

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"io"
	"net/http"
)

//mustNewSignedV4Request  NewSignedV4Request
func mustNewSignedV4Request(method string, urlStr string, contentLength int64, body io.ReadSeeker, st string, accessKey, secretKey string) (*http.Request, error) {
	log.Infof(urlStr)
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
