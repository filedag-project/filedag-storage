package s3api

import (
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/set"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	httpstatss "github.com/filedag-project/filedag-storage/objectservice/utils/httpstats"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

type s3ApiServer struct {
	authSys *iam.AuthSys
	store   *store.StorageSys
	bmSys   *store.BucketMetadataSys
	stats   *httpstatss.APIStatsSys
}

//registerS3Router Register APIs
func (s3a *s3ApiServer) registerS3Router(router *mux.Router, stats *httpstatss.APIStatsSys) {
	// API Router
	apiRouter := router.PathPrefix("/").Subrouter()
	// Readiness Probe
	apiRouter.Methods(http.MethodGet).Path("/status").HandlerFunc(stats.RecordAPIHandler("StatusHandler", s3a.StatusHandler))
	// NotFound
	apiRouter.NotFoundHandler = http.HandlerFunc(response.NotFoundHandler)
	var routers []*mux.Router
	routers = append(routers, apiRouter.PathPrefix("/{bucket}").Subrouter())

	for _, bucket := range routers {
		// Object operations
		//HeadObject
		bucket.Methods(http.MethodHead).Path("/{object:.+}").HandlerFunc(stats.RecordAPIHandler("HeadObject", s3a.HeadObjectHandler))

		// NewMultipartUpload
		bucket.Methods(http.MethodPost).Path("/{object:.+}").HandlerFunc(stats.RecordAPIHandler("NewMultipartUploadHandler", s3a.NewMultipartUploadHandler)).Queries("uploads", "")
		// CopyObjectPart
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(stats.RecordAPIHandler("CopyObjectPartHandler", s3a.CopyObjectPartHandler)).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
		// PutObjectPart
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(stats.RecordAPIHandler("PutObjectPartHandler", s3a.PutObjectPartHandler)).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
		// ListObjectParts
		bucket.Methods(http.MethodGet).Path("/{object:.+}").HandlerFunc(stats.RecordAPIHandler("ListObjectPartsHandler", s3a.ListObjectPartsHandler)).Queries("uploadId", "{uploadId:.*}")
		// ListMultipartUploads
		bucket.Methods(http.MethodGet).HandlerFunc(stats.RecordAPIHandler("ListMultipartUploadsHandler", s3a.ListMultipartUploadsHandler)).Queries("uploads", "")
		// CompleteMultipartUpload
		bucket.Methods(http.MethodPost).Path("/{object:.+}").HandlerFunc(stats.RecordAPIHandler("CompleteMultipartUploadHandler", s3a.CompleteMultipartUploadHandler)).Queries("uploadId", "{uploadId:.*}")
		// AbortMultipart
		bucket.Methods(http.MethodDelete).Path("/{object:.+}").HandlerFunc(stats.RecordAPIHandler("AbortMultipartUploadHandler", s3a.AbortMultipartUploadHandler)).Queries("uploadId", "{uploadId:.*}")

		// ListObjectsV2
		bucket.Methods(http.MethodGet).HandlerFunc(stats.RecordAPIHandler("ListObjectsV2Handler", s3a.ListObjectsV2Handler)).Queries("list-type", "2")
		// CopyObject
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(stats.RecordAPIHandler("CopyObjectHandler", s3a.CopyObjectHandler))
		// GetObject
		bucket.Methods(http.MethodGet).Path("/{object:.+}").HandlerFunc(stats.RecordAPIHandler("GetObjectHandler", s3a.GetObjectHandler))
		// PutObject
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(stats.RecordAPIHandler("PutObjectHandler", s3a.PutObjectHandler))
		// DeleteObject
		bucket.Methods(http.MethodDelete).Path("/{object:.+}").HandlerFunc(stats.RecordAPIHandler("DeleteObjectHandler", s3a.DeleteObjectHandler))
		// DeleteMultipleObjects
		bucket.Methods(http.MethodPost).HandlerFunc(stats.RecordAPIHandler("DeleteMultipleObjectsHandler", s3a.DeleteMultipleObjectsHandler)).Queries("delete", "")

		// Bucket operations
		// GetBucketLocation
		router.Methods(http.MethodGet).HandlerFunc(stats.RecordAPIHandler("GetBucketLocationHandler", s3a.GetBucketLocationHandler)).Queries("location", "")

		// PutBucketPolicy
		bucket.Methods(http.MethodPut).HandlerFunc(stats.RecordAPIHandler("PutBucketPolicyHandler", s3a.PutBucketPolicyHandler)).Queries("policy", "")
		// DeleteBucketPolicy
		bucket.Methods(http.MethodDelete).HandlerFunc(stats.RecordAPIHandler("DeleteBucketPolicyHandler", s3a.DeleteBucketPolicyHandler)).Queries("policy", "")
		// GetBucketPolicy
		bucket.Methods(http.MethodGet).HandlerFunc(stats.RecordAPIHandler("GetBucketPolicyHandler", s3a.GetBucketPolicyHandler)).Queries("policy", "")

		// GetBucketACL
		bucket.Methods(http.MethodGet).HandlerFunc(stats.RecordAPIHandler("GetBucketAclHandler", s3a.GetBucketAclHandler)).Queries("acl", "")
		// PutBucketACL
		bucket.Methods(http.MethodPut).HandlerFunc(stats.RecordAPIHandler("PutBucketAclHandler", s3a.PutBucketAclHandler)).Queries("acl", "")

		// GetBucketCors
		bucket.Methods(http.MethodGet).HandlerFunc(stats.RecordAPIHandler("GetBucketCorsHandler", s3a.GetBucketCorsHandler)).Queries("cors", "")
		// PutBucketCors
		bucket.Methods(http.MethodPut).HandlerFunc(stats.RecordAPIHandler("PutBucketCorsHandler", s3a.PutBucketCorsHandler)).Queries("cors", "")
		// DeleteBucketCors
		bucket.Methods(http.MethodDelete).HandlerFunc(stats.RecordAPIHandler("DeleteBucketCorsHandler", s3a.DeleteBucketCorsHandler)).Queries("cors", "")

		// PutBucketTaggingHandler
		bucket.Methods(http.MethodPut).HandlerFunc(stats.RecordAPIHandler("PutBucketTaggingHandler", s3a.PutBucketTaggingHandler)).Queries("tagging", "")
		// GetBucketTaggingHandler
		router.Methods(http.MethodGet).HandlerFunc(stats.RecordAPIHandler("GetBucketTaggingHandler", s3a.GetBucketTaggingHandler)).Queries("tagging", "")
		// DeleteBucketTaggingHandler
		router.Methods(http.MethodDelete).HandlerFunc(stats.RecordAPIHandler("DeleteBucketTaggingHandler", s3a.DeleteBucketTaggingHandler)).Queries("tagging", "")

		// PutBucket
		bucket.Methods(http.MethodPut).HandlerFunc(stats.RecordAPIHandler("PutBucketHandler", s3a.PutBucketHandler))
		// HeadBucket
		bucket.Methods(http.MethodHead).HandlerFunc(stats.RecordAPIHandler("HeadBucketHandler", s3a.HeadBucketHandler))
		// DeleteBucket
		bucket.Methods(http.MethodDelete).HandlerFunc(stats.RecordAPIHandler("DeleteBucketHandler", s3a.DeleteBucketHandler))

		// ListObjectsV1
		bucket.Methods(http.MethodGet).HandlerFunc(stats.RecordAPIHandler("ListObjectsV1Handler", s3a.ListObjectsV1Handler))
	}
	// ListBuckets
	apiRouter.Methods(http.MethodGet).Path("/").HandlerFunc(stats.RecordAPIHandler("ListBucketsHandler", s3a.ListBucketsHandler))
}

//registerSTSRouter Register AWS STS compatible APIs
func (s3a *s3ApiServer) registerSTSRouter(router *mux.Router, stats *httpstatss.APIStatsSys) {
	apiRouter := router.PathPrefix("/").Subrouter()
	apiRouter.Methods(http.MethodPost).MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		ctypeOk := set.MatchSimple("application/x-www-form-urlencoded*", r.Header.Get(consts.ContentType))
		authOk := set.MatchSimple(consts.SignV4Algorithm+"*", r.Header.Get(consts.Authorization))
		noQueries := len(r.URL.RawQuery) == 0
		return ctypeOk && authOk && noQueries
	}).HandlerFunc(stats.RecordAPIHandler("assumeRole-login", s3a.AssumeRole))
}

//NewS3Server Start a S3Server
func NewS3Server(router *mux.Router, authSys *iam.AuthSys, bmSys *store.BucketMetadataSys, storageSys *store.StorageSys, stats *httpstatss.APIStatsSys) {
	s3server := &s3ApiServer{
		authSys: authSys,
		store:   storageSys,
		bmSys:   bmSys,
		stats:   stats,
	}
	s3server.registerSTSRouter(router, stats)
	s3server.registerS3Router(router, stats)

	router.Use(iam.SetAuthHandler)
}

// CorsHandler handler for CORS (Cross Origin Resource Sharing)
func CorsHandler(handler http.Handler) http.Handler {
	commonS3Headers := []string{
		consts.Date,
		consts.ETag,
		consts.ServerInfo,
		consts.Connection,
		consts.AcceptRanges,
		consts.ContentRange,
		consts.ContentEncoding,
		consts.ContentLength,
		consts.ContentType,
		consts.ContentDisposition,
		consts.LastModified,
		consts.ContentLanguage,
		consts.CacheControl,
		consts.RetryAfter,
		consts.AmzBucketRegion,
		consts.Expires,
		consts.Authorization,
		consts.Action,
		consts.Range,
		"X-Amz*",
		"x-amz*",
		"*",
	}

	return cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {

			return true
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPut,
			http.MethodHead,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodPatch,
		},
		AllowedHeaders:   commonS3Headers,
		ExposedHeaders:   commonS3Headers,
		AllowCredentials: true,
	}).Handler(handler)
}
