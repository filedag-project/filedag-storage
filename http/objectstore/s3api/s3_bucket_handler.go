package s3api

import (
	"context"
	"encoding/xml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
)

var log = logging.Logger("server")

//ListAllMyBucketsResult  List All Buckets Result
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListAllMyBucketsResult"`
	Owner   *s3.Owner
	Buckets []*s3.Bucket `xml:"Buckets>Bucket"`
}

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

	resp := ListAllMyBucketsResult{
		Owner: &s3.Owner{
			ID:          aws.String(cred.AccessKey),
			DisplayName: aws.String(cred.AccessKey),
		},
		Buckets: buckets,
	}

	response.WriteSuccessResponseXML(w, r, resp)
}

//PutBucketHandler put a bucket
func (s3a *s3ApiServer) PutBucketHandler(w http.ResponseWriter, r *http.Request) {

	bucket, _ := GetBucketAndObject(r)
	log.Infof("PutBucketHandler %s", bucket)

	// avoid duplicated buckets
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(context.Background(), r, s3action.CreateBucketAction, bucket, "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	// create the folder for bucket, but lazily create actual collection
	if err := store.Mkdir("", bucket); err != nil {
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
