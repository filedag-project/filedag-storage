package s3api

import (
	"context"
	"encoding/xml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
	"time"
)

var log = logging.Logger("sever")

//ListAllMyBucketsResult  List All Buckets Result
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListAllMyBucketsResult"`
	Owner   *s3.Owner
	Buckets []*s3.Bucket `xml:"Buckets>Bucket"`
}

//ListBucketsHandler ListBuckets Handler
func (s3a *S3ApiServer) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	err := iam.CheckRequestAuthType(context.Background(), r, s3action.ListAllMyBucketsAction, "testbuckets", "")
	if err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	var buckets []*s3.Bucket
	buckets = append(buckets, &s3.Bucket{
		Name:         aws.String("testbuckets"),
		CreationDate: aws.Time(time.Unix(0, 0).UTC()),
	})
	resp := ListAllMyBucketsResult{
		Owner: &s3.Owner{
			ID:          aws.String("fds"),
			DisplayName: aws.String("fds admin"),
		},
		Buckets: buckets,
	}

	response.WriteSuccessResponseXML(w, r, resp)
}

//CreateBucketHandler Create Bucket
func (s3a *S3ApiServer) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := GetBucketAndObject(r)
	log.Infof("CreateBucketHandler %s", bucket)

	// create the folder for bucket, but lazily create actual collection
	if err := store.Mkdir(".", bucket); err != nil {
		log.Errorf("CreateBucketHandler mkdir: %v", err)
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}
