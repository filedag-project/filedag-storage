package s3api

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/set"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/gorilla/mux"
	"net/http"
)

type s3ApiServer struct {
	authSys iam.AuthSys
}

//registerS3Router Register S3Router
func (s3a *s3ApiServer) registerS3Router(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/").Subrouter()
	apiRouter.Methods(http.MethodPost).MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		ctypeOk := set.MatchSimple("application/x-www-form-urlencoded*", r.Header.Get(consts.ContentType))
		authOk := set.MatchSimple(consts.SignV4Algorithm+"*", r.Header.Get(consts.Authorization))
		noQueries := len(r.URL.RawQuery) == 0
		return ctypeOk && authOk && noQueries
	}).HandlerFunc(s3a.AssumeRole)
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
	var s3server s3ApiServer
	s3server.authSys.Init()
	s3server.registerS3Router(router)
}
