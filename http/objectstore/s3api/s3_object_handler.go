package s3api

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

//putObjectHandler Put ObjectHandler
func (s3a *s3ApiServer) putObjectHandler(w http.ResponseWriter, r *http.Request) {

	// http://docs.aws.amazon.com/AmazonS3/latest/dev/UploadingObjects.html

	bucket, object := getBucketAndObject(r)
	log.Infof("PutObjectHandler %s %s", bucket, object)
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.PutObjectAction, "testbuckets", "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	dataReader := r.Body
	defer dataReader.Close()
	objInfo, err2 := s3a.store.StoreObject(cred.AccessKey, bucket, object, r.Body)
	if err2 != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrStorePutFail)
		return
	}
	setPutObjHeaders(w, objInfo, false)
	response.WriteSuccessResponseHeadersOnly(w, r)
}

// getObjectHandler - GET Object
// ----------
// This implementation of the GET operation retrieves object. To use GET,
// you must have READ access to the object.
func (s3a *s3ApiServer) getObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	if _, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object); s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}

}

// headObjectHandler - HEAD Object
// -----------
// The HEAD operation retrieves metadata from an object without returning the object itself.
func (s3a *s3ApiServer) headObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	if _, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object); s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
}

// DeleteObjectHandler - delete an object
// Delete objectAPIHandlers
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

// setPutObjHeaders sets all the necessary headers returned back
// upon a success Put/Copy/CompleteMultipart/Delete requests
// to activate delete only headers set delete as true
func setPutObjHeaders(w http.ResponseWriter, objInfo store.ObjectInfo, delete bool) {
	// We must not use the http.Header().Set method here because some (broken)
	// clients expect the ETag header key to be literally "ETag" - not "Etag" (case-sensitive).
	// Therefore, we have to set the ETag directly as map entry.
	if objInfo.ETag != "" && !delete {
		w.Header()[consts.ETag] = []string{`"` + objInfo.ETag + `"`}
	}

	// Set the relevant version ID as part of the response header.
	if objInfo.VersionID != "" {
		w.Header()[consts.AmzVersionID] = []string{objInfo.VersionID}
		// If version is a deleted marker, set this header as well
		if objInfo.DeleteMarker && delete { // only returned during delete object
			w.Header()[consts.AmzDeleteMarker] = []string{strconv.FormatBool(objInfo.DeleteMarker)}
		}
	}

	if objInfo.Bucket != "" && objInfo.Name != "" {

	}
}
