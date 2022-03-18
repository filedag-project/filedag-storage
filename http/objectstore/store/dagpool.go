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
