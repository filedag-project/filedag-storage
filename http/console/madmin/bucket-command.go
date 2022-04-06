package madmin

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ListBucketsInfo returns the usage info for the authenticating account.
func (adm *AdminClient) ListBucketsInfo(ctx context.Context, opts AccountOpts) ([]*s3.Bucket, error) {
	var buckets []*s3.Bucket
	q := make(url.Values)
	if opts.PrefixUsage {
		q.Set("prefix-usage", "true")
	}
	resp, err := adm.executeMethod(ctx, http.MethodGet,
		requestData{
			relPath:     "/",
			queryValues: q,
		},
	)
	//fmt.Println(resp.Body)
	defer closeResponse(resp)
	if err != nil {
		return buckets, err
	}
	// Check response http status code
	if resp.StatusCode != http.StatusOK {
		return buckets, httpRespToErrorResponse(resp)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return buckets, err
	}
	fmt.Println(string(respBytes))
	var response ListAllMyBucketsResult
	err = xml.Unmarshal(respBytes, &response)
	if err != nil {
		return buckets, err
	}
	buckets = response.Buckets
	return buckets, nil
}

type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListAllMyBucketsResult"`
	Owner   *s3.Owner
	Buckets []*s3.Bucket `xml:"Buckets>Bucket"`
}

// PutBucket returns the usage info for the authenticating account.
func (adm *AdminClient) PutBucket(ctx context.Context, bucketName string, opts AccountOpts) error {

	queryValues := url.Values{}
	//queryValues.Set("bucketName", bucketName)

	reqData := requestData{
		//relPath:     adminAPIPrefix + "/add-user",
		relPath:     "/" + bucketName,
		queryValues: queryValues,
	}

	// Execute PUT on /minio/admin/v3/add-user to set a user.
	resp, err := adm.executeMethod(ctx, http.MethodPut, reqData)

	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}

// RemoveBucket returns the usage info for the authenticating account.
func (adm *AdminClient) RemoveBucket(ctx context.Context, bucketName string, opts AccountOpts) error {

	queryValues := url.Values{}
	//queryValues.Set("bucketName", bucketName)

	reqData := requestData{
		//relPath:     adminAPIPrefix + "/add-user",
		relPath:     "/" + bucketName,
		queryValues: queryValues,
	}

	// Execute PUT on /minio/admin/v3/add-user to set a user.
	resp, err := adm.executeMethod(ctx, http.MethodDelete, reqData)

	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}
