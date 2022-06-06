package s3api

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/gorilla/mux"
	"net/http"
)

func (s3a *s3ApiServer) PinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getBucketAndObject a")
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	object := vars["object"]
	//bucket, object := getBucketAndObject(r)
	fmt.Println("getBucketAndObject last")
	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetObjectAction, bucket, object)
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.authSys.PolicySys.Head(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}
	objInfo, ok := s3a.store.HasObject(r.Context(), cred.AccessKey, bucket, object)
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchKey)
		return
	}
	err := s3a.store.PinObject(r.Context(), cred.AccessKey, bucket, object)
	if err != nil {
		log.Errorf("PinObjectHandler PinObject  err:%v", err)
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	setPutObjHeaders(w, objInfo, true)
	response.WriteSuccessNoContent(w)
}

func (s3a *s3ApiServer) UnPinHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)

	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetObjectAction, bucket, object)
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.authSys.PolicySys.Head(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}
	objInfo, ok := s3a.store.HasObject(r.Context(), cred.AccessKey, bucket, object)
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchKey)
		return
	}
	err := s3a.store.UnPinObject(r.Context(), cred.AccessKey, bucket, object)
	if err != nil {
		log.Errorf("UnPinObjectHandler UnPinObject  err:%v", err)
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	setPutObjHeaders(w, objInfo, true)
	response.WriteSuccessNoContent(w)
}
