package s3api

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api/s3err"
	"github.com/gorilla/mux"
	"net/http"
)

type S3ApiServer struct {
}

//RegisterS3Router Register S3Router
func (s3a S3ApiServer) RegisterS3Router(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/").Subrouter()

	// Readiness Probe
	apiRouter.Methods("GET").Path("/status").HandlerFunc(s3a.StatusHandler)
	// NotFound
	apiRouter.NotFoundHandler = http.HandlerFunc(s3err.NotFoundHandler)
	// ListBuckets
	apiRouter.Methods("GET").Path("/").HandlerFunc(s3a.ListBucketsHandler)

}
