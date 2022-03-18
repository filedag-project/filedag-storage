package s3api

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
)

//PutObjectHandler Put ObjectHandler
func (s3a *s3ApiServer) PutObjectHandler(w http.ResponseWriter, r *http.Request) {

	// http://docs.aws.amazon.com/AmazonS3/latest/dev/UploadingObjects.html

	bucket, object := getBucketAndObject(r)
	log.Infof("PutObjectHandler %s %s", bucket, object)
	_, _, err := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.PutObjectAction, "testbuckets", "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	dataReader := r.Body
	defer dataReader.Close()
	cid := ""
	var errc error
	if cid, errc = s3a.store.PutFile(".", bucket+object, r.Body); errc != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrStorePutFail)
		return
	}
	w.Write([]byte(cid))
	response.WriteSuccessResponseEmpty(w, r)
}

// GetObjectHandler - GET Object
// ----------
// This implementation of the GET operation retrieves object. To use GET,
// you must have READ access to the object.
func (s3a *s3ApiServer) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	if _, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object); s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}

}

// HeadObjectHandler - HEAD Object
// -----------
// The HEAD operation retrieves metadata from an object without returning the object itself.
func (s3a *s3ApiServer) HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	if _, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object); s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
}

// Delete objectAPIHandlers

// DeleteObjectHandler - delete an object
func (s3a *s3ApiServer) DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {

	bucket, object := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	if _, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object); s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
}

func getBucketAndObject(r *http.Request) (bucket, object string) {
	vars := mux.Vars(r)
	bucket = vars["bucket"]
	object = vars["object"]
	if !strings.HasPrefix(object, "/") {
		object = "/" + object
	}

	return
}

func passThroughResponse(proxyResponse *http.Response, w http.ResponseWriter) (statusCode int) {
	for k, v := range proxyResponse.Header {
		w.Header()[k] = v
	}
	if proxyResponse.Header.Get("Content-Range") != "" && proxyResponse.StatusCode == 200 {
		w.WriteHeader(http.StatusPartialContent)
		statusCode = http.StatusPartialContent
	} else {
		statusCode = proxyResponse.StatusCode
	}
	w.WriteHeader(statusCode)
	if n, err := io.Copy(w, proxyResponse.Body); err != nil {
		log.Infof("passthrough response read %d bytes: %v", n, err)
	}
	return statusCode
}
