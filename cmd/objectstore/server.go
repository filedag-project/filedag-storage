package main

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iamapi"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
	"github.com/gorilla/mux"
	"net/http"
)

//StartServer Start a IamServer
func StartServer() {
	router := mux.NewRouter().SkipClean(true)
	s3api.NewS3Server(router)
	iamapi.NewIamApiServer(router)
	http.ListenAndServe(":9985", router)
}
