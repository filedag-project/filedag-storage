package main

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
	"github.com/gorilla/mux"
	"net/http"
)

//StartS3Server Start a S3Server
func StartS3Server() {
	var s3server s3api.S3ApiServer
	router := mux.NewRouter().SkipClean(true)
	s3server.RegisterS3Router(router)
	http.ListenAndServe("127.0.0.1:9985", router)
}
