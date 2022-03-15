package store

import (
	"io"
	"io/ioutil"
	"os"
)

//StorageSys store sys
type StorageSys struct {
}

//PutFile store object
func (s *StorageSys) PutFile(parentDirectoryPath string, dirName string, reader io.Reader) (string, error) {

	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(parentDirectoryPath+"/"+dirName, all, 0644)
	if err != nil {
		return "", err
	}
	return "cid", nil

}

//Mkdir store object
func (s *StorageSys) Mkdir(parentDirectoryPath string, bucket string) error {
	err := os.Mkdir(parentDirectoryPath+bucket, 0777)
	if err != nil {
		return err
	}
	return nil
}

//Init storage sys
func (s *StorageSys) Init() {

}
