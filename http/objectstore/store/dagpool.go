package store

import "io"

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
	return nil, nil
}

//DelFile del file
func (d dagPoolClient) DelFile(bucket, object string) error {
	//todo implement me
	return nil
}
