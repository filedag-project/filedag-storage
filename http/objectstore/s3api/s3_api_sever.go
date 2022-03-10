package s3api

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/gorilla/mux"
	"net/http"
)

type S3ApiServer struct {
}

//registerS3Router Register S3Router
func (s3a S3ApiServer) registerS3Router(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/").Subrouter()

	// Readiness Probe
	apiRouter.Methods("GET").Path("/status").HandlerFunc(s3a.StatusHandler)
	// NotFound
	apiRouter.NotFoundHandler = http.HandlerFunc(response.NotFoundHandler)
	// ListBuckets
	apiRouter.Methods("GET").Path("/").HandlerFunc(s3a.ListBucketsHandler)
	var routers []*mux.Router
	routers = append(routers, apiRouter.PathPrefix("/{bucket}").Subrouter())

	for _, bucket := range routers {
		// PutObject
		bucket.Methods("PUT").Path("/{object:.+}").HandlerFunc(s3a.PutObjectHandler)
		// PutBucket
		bucket.Methods("PUT").HandlerFunc(s3a.PutBucketHandler)

	}
}

//NewS3Server Start a S3Server
func NewS3Server(router *mux.Router) {
	var s3server S3ApiServer
	s3server.registerS3Router(router)
}
