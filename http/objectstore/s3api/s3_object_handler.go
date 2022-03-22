package s3api

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//PutObjectHandler Put ObjectHandler
func (s3a *s3ApiServer) PutObjectHandler(w http.ResponseWriter, r *http.Request) {

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

// GetObjectHandler - GET Object
// ----------
// This implementation of the GET operation retrieves object. To use GET,
// you must have READ access to the object.
func (s3a *s3ApiServer) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object)
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}

	objInfo, _, err := s3a.store.GetObject(cred.AccessKey, bucket, object)
	if err != nil {
		response.WriteErrorResponseHeadersOnly(w, r, api_errors.ErrInternalError)
		return
	}
	w.Header().Set(consts.AmzServerSideEncryption, consts.AmzEncryptionAES)

	if err = response.SetObjectHeaders(w, r, objInfo); err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrSetHeader)
		return
	}
	//todo use reader
	r1, _ := ioutil.ReadFile("./go.mod")
	w.Header().Set(consts.ContentLength, strconv.Itoa(len(r1)))
	response.SetHeadGetRespHeaders(w, r.Form)
	_, err = w.Write(r1)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrReader)
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
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object)
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}

	objInfo, _, err := s3a.store.GetObject(cred.AccessKey, bucket, object)
	if err != nil {
		response.WriteErrorResponseHeadersOnly(w, r, api_errors.ErrInternalError)
		return
	}
	w.Header().Set(consts.AmzServerSideEncryption, consts.AmzEncryptionAES)

	// Set standard object headers.
	if err = response.SetObjectHeaders(w, r, objInfo); err != nil {
		response.WriteErrorResponseHeadersOnly(w, r, api_errors.ErrSetHeader)
		return
	}
	// Set any additional requested response headers.
	response.SetHeadGetRespHeaders(w, r.Form)

	// Successful response.
	w.WriteHeader(http.StatusOK)

}

// DeleteObjectHandler - delete an object
// Delete objectAPIHandlers
func (s3a *s3ApiServer) DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object)
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	objInfo, _, err := s3a.store.GetObject(cred.AccessKey, bucket, object)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	err = s3a.store.DeleteObject(cred.AccessKey, bucket, object)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	setPutObjHeaders(w, objInfo, true)
	response.WriteSuccessNoContent(w)
}

// CopyObjectHandler - Copy Object
// ----------
// This implementation of the PUT operation adds an object to a bucket
// while reading the object from another source.
// Notice: The S3 client can send secret keys in headers for encryption related jobs,
// the handler should ensure to remove these keys before sending them to the object layer.
// Currently these keys are:
//   - X-Amz-Server-Side-Encryption-Customer-Key
//   - X-Amz-Copy-Source-Server-Side-Encryption-Customer-Key
func (s3a *s3ApiServer) CopyObjectHandler(w http.ResponseWriter, r *http.Request) {
	dstBucket, dstObject := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, dstBucket, dstObject)
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	// Copy source path.
	cpSrcPath, err := url.QueryUnescape(r.Header.Get("X-Amz-Copy-Source"))
	if err != nil {
		// Save unescaped string as is.
		cpSrcPath = r.Header.Get("X-Amz-Copy-Source")
	}

	srcBucket, srcObject := pathToBucketAndObject(cpSrcPath)

	log.Infof("CopyObjectHandler %s %s => %s %s", srcBucket, srcObject, dstBucket, dstObject)

	_, i, err := s3a.store.GetObject(cred.AccessKey, srcBucket, srcObject)
	if err != nil {
		response.WriteErrorResponseHeadersOnly(w, r, api_errors.ErrInternalError)
		return
	}
	if (srcBucket == dstBucket && srcObject == dstObject || cpSrcPath == "") && isReplace(r) {
		object, err := s3a.store.StoreObject(cred.AccessKey, dstBucket, dstObject, i)
		if err != nil {
			return
		}
		response.WriteSuccessResponseXML(w, r, response.CopyObjectResult{
			ETag:         fmt.Sprintf("%x", object.ETag),
			LastModified: time.Now().UTC(),
		})
		return
	}

	// If source object is empty or bucket is empty, reply back invalid copy source.
	if srcObject == "" || srcBucket == "" {
		response.WriteErrorResponse(w, r, api_errors.ErrInvalidCopySource)
		return
	}

	if srcBucket == dstBucket && srcObject == dstObject {
		response.WriteErrorResponse(w, r, api_errors.ErrInvalidCopyDest)
		return
	}
	obj, err := s3a.store.StoreObject(cred.AccessKey, dstBucket, dstObject, i)

	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrStorePutFail)
		return
	}

	setEtag(w, obj.ETag)

	resp2 := response.CopyObjectResult{
		ETag:         obj.ETag,
		LastModified: time.Now().UTC(),
	}

	response.WriteSuccessResponseXML(w, r, resp2)
}

func (s3a *s3ApiServer) ListObjectsV1Handler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	_, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object)
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
}

func (s3a *s3ApiServer) ListObjectsV2Handler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)
	var ctx = context.Background()
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	_, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetObjectAction, bucket, object)
	if s3Error != api_errors.ErrNone {
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
func pathToBucketAndObject(path string) (bucket, object string) {
	path = strings.TrimPrefix(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 2 {
		return parts[0], "/" + parts[1]
	}
	return parts[0], "/"
}
func setEtag(w http.ResponseWriter, etag string) {
	if etag != "" {
		if strings.HasPrefix(etag, "\"") {
			w.Header().Set("ETag", etag)
		} else {
			w.Header().Set("ETag", "\""+etag+"\"")
		}
	}
}
func isReplace(r *http.Request) bool {
	return r.Header.Get("X-Amz-Metadata-Directive") == "REPLACE"
}
