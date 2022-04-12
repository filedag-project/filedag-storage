package store

import (
	"bytes"
	"io"
	"io/ioutil"
)

//dagPoolClient dagPool Client
type dagPoolClient struct {
}

//PutFile put file
func (d dagPoolClient) PutFile(bucket, object string, reader io.Reader) (string, error) {
	//todo implement me
	return "cid", nil
}

//GetFile get file
func (d dagPoolClient) GetFile(bucket, object string) (io.Reader, error) {
	//todo implement me
	//todo use reader
	r1, _ := ioutil.ReadFile("./go.mod")
	return bytes.NewReader(r1), nil
}

//DelFile del file
func (d dagPoolClient) DelFile(bucket, object string) error {
	//todo implement me
	return nil
}

//MkBucket del file
func (d dagPoolClient) MkBucket(bucket string) error {
	//todo implement me
	return nil
}
