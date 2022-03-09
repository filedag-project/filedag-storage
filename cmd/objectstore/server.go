package main

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iamapi"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
	"github.com/filedag-project/filedag-storage/http/objectstore/sts_api"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
)

var log = logging.Logger("sever")

//startServer Start a IamServer
func startServer() {
	err := logging.SetLogLevel("*", "INFO")
	if err != nil {
		return
	}
	router := mux.NewRouter()
	sts_api.RegisterSTSRouter(router)
	s3api.NewS3Server(router)
	iamapi.NewIamApiServer(router)
	for _, ip := range utils.MustGetLocalIP4().ToSlice() {
		log.Infof("start sever at http://%v:%v", ip, 9985)
	}
	err = http.ListenAndServe(":9985", router)
	if err != nil {
		log.Errorf("ListenAndServe err%v", err)

		return
	}

}
