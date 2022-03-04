package store

import (
	"io"
	"io/ioutil"
	"os"
)

//PutFile store object
func PutFile(parentDirectoryPath string, dirName string, reader io.Reader) (string, error) {

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
func Mkdir(parentDirectoryPath string, bucket string) error {
	err := os.Mkdir(parentDirectoryPath+bucket, 777)
	if err != nil {
		return err
	}
	return nil
}
