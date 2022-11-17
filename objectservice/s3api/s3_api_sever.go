package s3api

import (
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/set"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/store"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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
		// Object operations
		//HeadObject
		bucket.Methods(http.MethodHead).Path("/{object:.+}").HandlerFunc(s3a.HeadObjectHandler)

		// NewMultipartUpload
		bucket.Methods(http.MethodPost).Path("/{object:.+}").HandlerFunc(s3a.NewMultipartUploadHandler).Queries("uploads", "")
		// CopyObjectPart
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(s3a.CopyObjectPartHandler).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
		// PutObjectPart
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(s3a.PutObjectPartHandler).Queries("partNumber", "{partNumber:[0-9]+}", "uploadId", "{uploadId:.*}")
		// ListObjectParts
		bucket.Methods(http.MethodGet).Path("/{object:.+}").HandlerFunc(s3a.ListObjectPartsHandler).Queries("uploadId", "{uploadId:.*}")
		// ListMultipartUploads
		bucket.Methods(http.MethodGet).HandlerFunc(s3a.ListMultipartUploadsHandler).Queries("uploads", "")
		// CompleteMultipartUpload
		bucket.Methods(http.MethodPost).Path("/{object:.+}").HandlerFunc(s3a.CompleteMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}")
		// AbortMultipart
		bucket.Methods(http.MethodDelete).Path("/{object:.+}").HandlerFunc(s3a.AbortMultipartUploadHandler).Queries("uploadId", "{uploadId:.*}")

		// ListObjectsV2
		bucket.Methods(http.MethodGet).HandlerFunc(s3a.ListObjectsV2Handler).Queries("list-type", "2")
		// CopyObject
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HeadersRegexp("X-Amz-Copy-Source", ".*?(\\/|%2F).*?").HandlerFunc(s3a.CopyObjectHandler)
		// GetObject
		bucket.Methods(http.MethodGet).Path("/{object:.+}").HandlerFunc(s3a.GetObjectHandler)
		// PutObject
		bucket.Methods(http.MethodPut).Path("/{object:.+}").HandlerFunc(s3a.PutObjectHandler)
		// DeleteObject
		bucket.Methods(http.MethodDelete).Path("/{object:.+}").HandlerFunc(s3a.DeleteObjectHandler)
		// DeleteMultipleObjects
		bucket.Methods(http.MethodPost).HandlerFunc(s3a.DeleteMultipleObjectsHandler).Queries("delete", "")

		// Bucket operations
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

		// ListObjectsV1
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
func NewS3Server(router *mux.Router, authSys *iam.AuthSys, bmSys *store.BucketMetadataSys, storageSys *store.StorageSys) {
	s3server := &s3ApiServer{
		authSys: authSys,
		store:   storageSys,
		bmSys:   bmSys,
	}
	s3server.registerSTSRouter(router)
	s3server.registerS3Router(router)

	router.Use(iam.SetAuthHandler)
}

// CorsHandler handler for CORS (Cross Origin Resource Sharing)
func CorsHandler(handler http.Handler) http.Handler {
	commonS3Headers := []string{
		Date,
		ETag,
		ServerInfo,
		Connection,
		AcceptRanges,
		ContentRange,
		ContentEncoding,
		ContentLength,
		ContentType,
		ContentDisposition,
		LastModified,
		ContentLanguage,
		CacheControl,
		RetryAfter,
		AmzBucketRegion,
		Expires,
		Authorization,
		Action,
		Range,
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

const (
	LastModified       = "Last-Modified"
	Date               = "Date"
	ETag               = "ETag"
	ContentType        = "Content-Type"
	ContentMD5         = "Content-Md5"
	ContentEncoding    = "Content-Encoding"
	Expires            = "Expires"
	ContentLength      = "Content-Length"
	ContentLanguage    = "Content-Language"
	ContentRange       = "Content-Range"
	Connection         = "Connection"
	AcceptRanges       = "Accept-Ranges"
	AmzBucketRegion    = "X-Amz-Bucket-Region"
	ServerInfo         = "Server"
	RetryAfter         = "Retry-After"
	Location           = "Location"
	CacheControl       = "Cache-Control"
	ContentDisposition = "Content-Disposition"
	Authorization      = "Authorization"
	Action             = "Action"
	Range              = "Range"
)
