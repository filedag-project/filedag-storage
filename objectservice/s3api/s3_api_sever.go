package s3api

import (
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/set"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"github.com/gorilla/mux"
	ipld "github.com/ipfs/go-ipld-format"
	"net/http"
)

type s3ApiServer struct {
	authSys *iam.AuthSys
	store   *store.StorageSys
	bmSys   *store.BucketMetadataSys
}

//registerS3Router Register APIs
func (s3a *s3ApiServer) registerS3Router(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/").Subrouter()
	// Readiness Probe
	apiRouter.Methods(http.MethodGet).Path("/status").HandlerFunc(s3a.StatusHandler)
	// NotFound
	apiRouter.NotFoundHandler = http.HandlerFunc(response.NotFoundHandler)
	var routers []*mux.Router
	routers = append(routers, apiRouter.PathPrefix("/{bucket}").Subrouter())

	for _, bucket := range routers {
		// ListObjectsV2
		bucket.Methods(http.MethodGet).HandlerFunc(s3a.ListObjectsV2Handler).Queries("list-type", "2")
		// CopyObject
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(s3a.CopyObjectHandler)

		// GetObject
		bucket.Methods(http.MethodGet).Path("/{object:.+}").HandlerFunc(s3a.GetObjectHandler)

		// PutObject
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(s3a.PutObjectHandler)
		//HeadObject
		bucket.Methods(http.MethodHead).Path("/{object:.+}").HandlerFunc(s3a.HeadObjectHandler)

		// DeleteObject
		bucket.Methods(http.MethodDelete).Path("/{object:.+}").HandlerFunc(s3a.DeleteObjectHandler)
		// DeleteMultipleObjects
		bucket.Methods(http.MethodPost).HandlerFunc(s3a.DeleteMultipleObjectsHandler).Queries("delete", "")
		// GetBucketLocation
		router.Methods(http.MethodGet).HandlerFunc(s3a.GetBucketLocationHandler).Queries("location", "")

		// PutBucketPolicy
		bucket.Methods(http.MethodPut).HandlerFunc(s3a.PutBucketPolicyHandler).Queries("policy", "")
		// DeleteBucketPolicy
		bucket.Methods(http.MethodDelete).HandlerFunc(s3a.DeleteBucketPolicyHandler).Queries("policy", "")
		// GetBucketPolicy
		bucket.Methods(http.MethodGet).HandlerFunc(s3a.GetBucketPolicyHandler).Queries("policy", "")

		// GetBucketACL
		bucket.Methods(http.MethodGet).HandlerFunc(s3a.GetBucketAclHandler).Queries("acl", "")
		// PutBucketACL
		bucket.Methods(http.MethodPut).HandlerFunc(s3a.PutBucketAclHandler).Queries("acl", "")

		// GetBucketCors
		bucket.Methods(http.MethodGet).HandlerFunc(s3a.GetBucketCorsHandler).Queries("cors", "")
		// PutBucketCors
		bucket.Methods(http.MethodPut).HandlerFunc(s3a.PutBucketCorsHandler).Queries("cors", "")
		// DeleteBucketCors
		bucket.Methods(http.MethodDelete).HandlerFunc(s3a.DeleteBucketCorsHandler).Queries("cors", "")

		// PutBucketTaggingHandler
		bucket.Methods(http.MethodPut).HandlerFunc(s3a.PutBucketTaggingHandler).Queries("tagging", "")
		// GetBucketTaggingHandler
		router.Methods(http.MethodGet).HandlerFunc(s3a.GetBucketTaggingHandler).Queries("tagging", "")
		// DeleteBucketTaggingHandler
		router.Methods(http.MethodDelete).HandlerFunc(s3a.DeleteBucketTaggingHandler).Queries("tagging", "")

		// PutBucket
		bucket.Methods(http.MethodPut).HandlerFunc(s3a.PutBucketHandler)
		// HeadBucket
		bucket.Methods(http.MethodHead).HandlerFunc(s3a.HeadBucketHandler)
		// DeleteBucket
		bucket.Methods(http.MethodDelete).HandlerFunc(s3a.DeleteBucketHandler)
		bucket.Methods(http.MethodGet).HandlerFunc(s3a.ListObjectsV1Handler)
	}
	// ListBuckets
	apiRouter.Methods(http.MethodGet).Path("/").HandlerFunc(s3a.ListBucketsHandler)
}

//registerSTSRouter Register AWS STS compatible APIs
func (s3a *s3ApiServer) registerSTSRouter(router *mux.Router) {
	apiRouter := router.PathPrefix("/").Subrouter()
	apiRouter.Methods(http.MethodPost).MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		ctypeOk := set.MatchSimple("application/x-www-form-urlencoded*", r.Header.Get(consts.ContentType))
		authOk := set.MatchSimple(consts.SignV4Algorithm+"*", r.Header.Get(consts.Authorization))
		noQueries := len(r.URL.RawQuery) == 0
		return ctypeOk && authOk && noQueries
	}).HandlerFunc(s3a.AssumeRole)
}

//NewS3Server Start a S3Server
func NewS3Server(router *mux.Router, dagService ipld.DAGService, authSys *iam.AuthSys, db *uleveldb.ULevelDB) {
	s3server := &s3ApiServer{
		authSys: authSys,
		store:   store.NewStorageSys(dagService, db),
		bmSys:   store.NewBucketMetadataSys(db),
	}
	s3server.store.SetNewBucketNSLock(s3server.bmSys.NewNSLock)
	s3server.store.SetHasBucket(s3server.bmSys.HasBucket)
	s3server.bmSys.SetEmptyBucket(s3server.store.EmptyBucket)
	s3server.registerSTSRouter(router)
	s3server.registerS3Router(router)

	router.Use(iam.SetAuthHandler)
}
