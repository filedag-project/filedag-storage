package s3api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dustin/go-humanize"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"io"
	"io/ioutil"
	"net/http"
)

const maxBucketPolicySize = 20 * humanize.KiByte

//PutBucketPolicyHandler Put BucketPolicy
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketPolicy.html
func (s3a *s3ApiServer) PutBucketPolicyHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	var ctx = context.Background()
	log.Infof("PutBucketPolicyHandler %s", bucket)
	cred, _, errc := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.PutBucketPolicyAction, bucket, "")
	if errc != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, errc)
		return
	}
	// Error out if Content-Length is beyond allowed size.
	if r.ContentLength > maxBucketPolicySize {
		response.WriteErrorResponse(w, r, api_errors.ErrPolicyTooLarge)
		return
	}
	bucketPolicyBytes, err := ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNewReaderFail)
		return
	}
	bucketPolicy, err := policy.ParseConfig(bytes.NewReader(bucketPolicyBytes), bucket)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrJsonMarshal)
		return
	}
	// Version in policy must not be empty
	if bucketPolicy.Version == "" {
		response.WriteErrorResponse(w, r, api_errors.ErrMalformedPolicy)
		return
	}

	if err = s3a.authSys.PolicySys.UpdatePolicy(ctx, cred.AccessKey, bucket, bucketPolicy); err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}

//DeleteBucketPolicyHandler Delete BucketPolicy
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_control_DeleteBucketPolicy.html
func (s3a *s3ApiServer) DeleteBucketPolicyHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	var ctx = context.Background()
	log.Infof("DeleteBucketPolicyHandler %s", bucket)
	cred, _, errc := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.DeleteBucketPolicyAction, bucket, "")
	if errc != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, errc)
		return
	}
	if err := s3a.authSys.PolicySys.DeletePolicy(ctx, cred.AccessKey, bucket, nil); err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrSetBucketPolicyFail)
		return
	}
	// Success.
	response.WriteSuccessResponseEmpty(w, r)
}

//GetBucketPolicyHandler Get BucketPolicy
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketPolicy.html
func (s3a *s3ApiServer) GetBucketPolicyHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	log.Infof("PutBucketPolicyHandler %s", bucket)
	cred, _, errc := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.GetBucketPolicyAction, bucket, "")
	if errc != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, errc)
		return
	}

	// Read bucket access policy.
	config, err := s3a.authSys.PolicySys.Get(bucket, cred.AccessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucketPolicy)
		return
	}

	configData, err := json.Marshal(config)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrJsonMarshal)
		return
	}

	// Write to client.
	response.WriteSuccessResponseJSON(w, configData)
}
