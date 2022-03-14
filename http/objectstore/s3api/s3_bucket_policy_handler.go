package s3api

import (
	"bytes"
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"io"
	"io/ioutil"
	"net/http"
)

//PutBucketPolicyHandler Put BucketPolicy
func (s3a *s3ApiServer) PutBucketPolicyHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := GetBucketAndObject(r)
	var ctx = context.Background()
	log.Infof("PutBucketPolicyHandler %s", bucket)
	cred, _, errc := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.CreateBucketAction, bucket, "")
	if errc != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, errc)
		return
	}
	bucketPolicyBytes, err := ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrReader)
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

	if err = s3a.authSys.PolicySys.Update(ctx, cred.AccessKey, bucket, bucketPolicy); err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}
