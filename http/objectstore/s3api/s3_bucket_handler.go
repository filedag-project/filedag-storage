package s3api

import (
	"encoding/xml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api/s3resp"
	"net/http"
	"time"
)

//ListAllMyBucketsResult  List All Buckets Result
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListAllMyBucketsResult"`
	Owner   *s3.Owner
	Buckets []*s3.Bucket `xml:"Buckets>Bucket"`
}

//ListBucketsHandler ListBuckets Handler
func (s3a *S3ApiServer) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	var buckets []*s3.Bucket
	buckets = append(buckets, &s3.Bucket{
		Name:         aws.String("testbuckets"),
		CreationDate: aws.Time(time.Unix(0, 0).UTC()),
	})
	response := ListAllMyBucketsResult{
		Owner: &s3.Owner{
			ID:          aws.String("fds"),
			DisplayName: aws.String("fds admin"),
		},
		Buckets: buckets,
	}

	s3resp.WriteSuccessResponseXML(w, r, response)
}
