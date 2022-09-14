package s3api

import (
	"bytes"
	"encoding/json"
	"github.com/dustin/go-humanize"
	"github.com/filedag-project/filedag-storage/http/objectstore/apierrors"
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

	log.Infof("PutBucketPolicyHandler %s", bucket)
	cred, _, errc := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.PutBucketPolicyAction, bucket, "")
	if errc != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, errc)
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

	if err = s3a.bmSys.UpdateBucketPolicy(r.Context(), cred.AccessKey, bucket, bucketPolicy); err != nil {
		log.Errorf("PutBucketPolicyHandler UpdateBucketPolicy err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}

//DeleteBucketPolicyHandler Delete BucketPolicy
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_control_DeleteBucketPolicy.html
func (s3a *s3ApiServer) DeleteBucketPolicyHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)

	log.Infof("DeleteBucketPolicyHandler %s", bucket)
	cred, _, errc := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.DeleteBucketPolicyAction, bucket, "")
	if errc != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, errc)
		return
	}
	if err := s3a.bmSys.DeleteBucketPolicy(r.Context(), cred.AccessKey, bucket, nil); err != nil {
		log.Errorf("DeleteBucketPolicyHandler DeleteBucketPolicy err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	// Success.
	response.WriteSuccessResponseEmpty(w, r)
}

//GetBucketPolicyHandler Get BucketPolicy
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketPolicy.html
func (s3a *s3ApiServer) GetBucketPolicyHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	log.Infof("GetBucketPolicyHandler %s", bucket)
	cred, _, errc := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetBucketPolicyAction, bucket, "")
	if errc != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, errc)
		return
	}

	// Read bucket access policy.
	config, err := s3a.bmSys.GetPolicyConfig(bucket, cred.AccessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucketPolicy)
		return
	}

	configData, err := json.Marshal(config)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrMalformedJSON)
		return
	}

	// Write to client.
	response.WriteSuccessResponseJSON(w, configData)
}
