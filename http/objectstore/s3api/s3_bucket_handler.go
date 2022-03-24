package s3api

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	logging "github.com/ipfs/go-log/v2"
	"io"
	"net/http"
)

var log = logging.Logger("server")

//ListBucketsHandler ListBuckets Handler
func (s3a *s3ApiServer) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.ListAllMyBucketsAction, "testbuckets", "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	bucketMetas, erro := s3a.authSys.PolicySys.GetAllBucketOfUser(cred.AccessKey)
	if erro != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	var buckets []*s3.Bucket
	for _, b := range bucketMetas {
		buckets = append(buckets, &s3.Bucket{
			Name:         aws.String(b.Name),
			CreationDate: aws.Time(b.Created),
		})
	}

	resp := response.ListAllMyBucketsResult{
		Owner: &s3.Owner{
			ID:          aws.String(consts.DefaultOwnerID),
			DisplayName: aws.String(cred.AccessKey),
		},
		Buckets: buckets,
	}

	response.WriteSuccessResponseXML(w, r, resp)
}

//PutBucketHandler put a bucket
func (s3a *s3ApiServer) PutBucketHandler(w http.ResponseWriter, r *http.Request) {

	bucket, _ := getBucketAndObject(r)
	log.Infof("PutBucketHandler %s", bucket)

	// avoid duplicated buckets
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.CreateBucketAction, bucket, "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	get, _ := s3a.authSys.PolicySys.Get(cred.AccessKey, bucket)
	if get != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrBucketAlreadyExists)
		return
	}
	// create the folder for bucket, but lazily create actual collection
	if err := s3a.store.MkBucket("", bucket); err != nil {
		log.Errorf("PutBucketHandler mkdir: %v", err)
		response.WriteErrorResponse(w, r, api_errors.ErrStoreMkdirFail)
		return
	}
	erro := s3a.authSys.PolicySys.Set(bucket, cred.AccessKey)
	if erro != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrSetBucketPolicyFail)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}

// HeadBucketHandler - HEAD Bucket
// ----------
// This operation is useful to determine if a bucket exists.
// The operation returns a 200 OK if the bucket exists and you
// have permission to access it. Otherwise, the operation might
// return responses such as 404 Not Found and 403 Forbidden.
func (s3a *s3ApiServer) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {

	bucket, _ := getBucketAndObject(r)
	log.Infof("HeadBucketHandler %s", bucket)
	// avoid duplicated buckets
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.HeadBucketAction, bucket, "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	if ok := s3a.authSys.PolicySys.Head(bucket, cred.AccessKey); !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}

	response.WriteSuccessResponseEmpty(w, r)
}

// DeleteBucketHandler delete Bucket
func (s3a *s3ApiServer) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {

	bucket, _ := getBucketAndObject(r)
	log.Infof("DeleteBucketHandler %s", bucket)

	// avoid duplicated buckets
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.DeleteBucketAction, bucket, "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	errc := s3a.authSys.PolicySys.Delete(context.Background(), cred.AccessKey, bucket)
	if errc != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucketPolicy)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}

// GetBucketAclHandler Get Bucket ACL
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketAcl.html
func (s3a *s3ApiServer) GetBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	// collect parameters
	bucket, _ := getBucketAndObject(r)
	log.Infof("GetBucketAclHandler %s", bucket)
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.GetBucketPolicyAction, bucket, "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	resp := response.AccessControlPolicy{}
	id := cred.AccessKey
	if resp.Owner.DisplayName == "" {
		resp.Owner.DisplayName = cred.AccessKey
		resp.Owner.ID = id
	}
	resp.AccessControlList.Grant = append(resp.AccessControlList.Grant, response.Grant{
		Grantee: response.Grantee{
			ID:          id,
			DisplayName: cred.AccessKey,
			Type:        "CanonicalUser",
			XMLXSI:      "CanonicalUser",
			XMLNS:       "http://www.w3.org/2001/XMLSchema-instance"},
		Permission: "FULL_CONTROL", //todo change
	})
	response.WriteSuccessResponseXML(w, r, resp)
}

// GetBucketCorsHandler Get bucket CORS
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketCors.html
func (s3a *s3ApiServer) GetBucketCorsHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteErrorResponse(w, r, api_errors.ErrNoSuchCORSConfiguration)
}

// PutBucketCorsHandler Put bucket CORS
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketCors.html
func (s3a *s3ApiServer) PutBucketCorsHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteErrorResponse(w, r, api_errors.ErrNotImplemented)
}

// DeleteBucketCorsHandler Delete bucket CORS
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketCors.html
func (s3a *s3ApiServer) DeleteBucketCorsHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteErrorResponse(w, r, http.StatusNoContent)
}

// PutBucketAclHandler Put bucket ACL
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketAcl.html
func (s3a *s3ApiServer) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)

	// Allow putBucketACL if policy action is set, since this is a dummy call
	// we are simply re-purposing the bucketPolicyAction.
	_, _, err := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.PutBucketPolicyAction, bucket, "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	aclHeader := r.Header.Get(consts.AmzACL)
	if aclHeader == "" {
		acl := &response.AccessControlPolicy{}
		if errc := utils.XmlDecoder(r.Body, acl, r.ContentLength); errc != nil {
			if errc == io.EOF {
				response.WriteErrorResponse(w, r, api_errors.ErrMissingSecurityHeader)
				return
			}
			response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
			return
		}

		if len(acl.AccessControlList.Grant) == 0 {
			response.WriteErrorResponse(w, r, api_errors.ErrNotImplemented)
			return
		}

		if acl.AccessControlList.Grant[0].Permission != "FULL_CONTROL" {
			response.WriteErrorResponse(w, r, api_errors.ErrNotImplemented)
			return
		}
	}

	if aclHeader != "" && aclHeader != "private" {
		response.WriteErrorResponse(w, r, api_errors.ErrNotImplemented)
		return
	}
}
