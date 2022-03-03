package s3api

import (
	"bytes"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/google/martian/log"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestS3ApiServer_PutObjectHandler(t *testing.T) {
	var s3server S3ApiServer
	router := mux.NewRouter().SkipClean(true)
	s3server.RegisterS3Router(router)
	http.ListenAndServe("127.0.0.1:9985", router)
	url := "http://127.0.0.1:9985/test/test1.txt"
	fiveMBBytes := bytes.Repeat([]byte("a"), 5*humanize.KiByte)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(fiveMBBytes))
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	res, err := client.Do(req)
	log.Infof("err%v", err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
