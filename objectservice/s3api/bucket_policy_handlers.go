package s3api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy"
	"github.com/filedag-project/filedag-storage/objectservice/iam/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const maxBucketPolicySize = 20 * humanize.KiByte

//PutBucketPolicyHandler Put BucketPolicy
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketPolicy.html
func (s3a *s3ApiServer) PutBucketPolicyHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _, _ := getBucketAndObject(r)

	log.Infof("PutBucketPolicyHandler %s", bucket)
	_, _, s3err := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.PutBucketPolicyAction, bucket, "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}

	if r.ContentLength <= 0 {
		response.WriteErrorResponse(w, r, apierrors.ErrMissingContentLength)
		return
	}

	// Error out if Content-Length is beyond allowed size.
	if r.ContentLength > maxBucketPolicySize {
		response.WriteErrorResponse(w, r, apierrors.ErrIncompleteBody)
		return
	}
	bucketPolicyBytes, err := ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidPolicyDocument)
		return
	}
	bucketPolicy, err := policy.ParseConfig(bytes.NewReader(bucketPolicyBytes), bucket)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrMalformedPolicy)
		return
	}
	//// Version in policy must not be empty
	//if bucketPolicy.Version == "" {
	//	response.WriteErrorResponse(w, r, apierrors.ErrMalformedPolicy)
	//	return
	//}

	if err = s3a.bmSys.UpdateBucketPolicy(r.Context(), bucket, bucketPolicy); err != nil {
		log.Errorf("PutBucketPolicyHandler UpdateBucketPolicy err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ToApiError(r.Context(), err))
		return
	}
	response.WriteSuccessNoContent(w)
}

//DeleteBucketPolicyHandler Delete BucketPolicy
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_control_DeleteBucketPolicy.html
func (s3a *s3ApiServer) DeleteBucketPolicyHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _, _ := getBucketAndObject(r)
	ctx := r.Context()

	log.Infof("DeleteBucketPolicyHandler %s", bucket)
	_, _, s3err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.DeleteBucketPolicyAction, bucket, "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}
	if err := s3a.bmSys.DeleteBucketPolicy(ctx, bucket); err != nil {
		log.Errorf("DeleteBucketPolicyHandler DeleteBucketPolicy err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}
	// Success.
	response.WriteSuccessResponseHeadersOnly(w, r)
}

//GetBucketPolicyHandler Get BucketPolicy
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketPolicy.html
func (s3a *s3ApiServer) GetBucketPolicyHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _, _ := getBucketAndObject(r)
	ctx := r.Context()
	log.Infof("GetBucketPolicyHandler %s", bucket)
	_, _, s3err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.GetBucketPolicyAction, bucket, "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}

	// Read bucket access policy.
	config, err := s3a.bmSys.GetPolicyConfig(ctx, bucket)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	configData, err := json.Marshal(config)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrMalformedJSON)
		return
	}

	w.Header().Set(consts.ServerInfo, "FDS")
	w.Header().Set(consts.AmzRequestID, fmt.Sprintf("%d", time.Now().UnixNano()))
	w.Header().Set(consts.AcceptRanges, "bytes")
	if r.Header.Get("Origin") != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	w.Header().Set(consts.ContentType, "application/json")
	w.Header().Set(consts.ContentLength, strconv.Itoa(len(configData)))
	w.WriteHeader(http.StatusOK)
	if configData != nil {
		w.Write(configData)
	}
	// Write to client.
	//response.WriteSuccessResponseJSON(w, r, config)
}
