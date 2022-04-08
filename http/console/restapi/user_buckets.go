package restapi

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/madmin/policy"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/go-openapi/swag"
	"strings"
	"time"
)

// getListBucketsResponse
func getListBucketsResponse(session *models.Principal) (*models.ListBucketsResponse, *models.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	mAdmin, err := NewAdminClient(session)
	if err != nil {
		return nil, prepareError(err)
	}
	adminClient := AdminClient{Client: mAdmin}
	buckets, err := getlistBuckets(ctx, adminClient)
	if err != nil {
		return nil, prepareError(err)
	}
	// serialize output
	listBucketsResponse := &models.ListBucketsResponse{
		Buckets: buckets,
		Total:   int64(len(buckets)),
	}
	return listBucketsResponse, nil
}

// getlistBuckets
func getlistBuckets(ctx context.Context, client Admin) ([]*models.Bucket, error) {
	info, err := client.listBucketsInfo(ctx)
	if err != nil {
		return []*models.Bucket{}, err
	}
	var bucketInfos []*models.Bucket
	for _, bucket := range info {
		bucketElem := &models.Bucket{
			CreationDate: bucket.CreationDate.Format(time.RFC3339),
			Details: &models.BucketDetails{
				Quota: nil,
			},
			Name: swag.String(*bucket.Name),
		}
		bucketInfos = append(bucketInfos, bucketElem)
	}
	return bucketInfos, nil
}

// getCreateBucketResponse performs putBucket() to create a bucket with its access policy
func getCreateBucketResponse(session *models.Principal, buchetName, location string, opts bool) *models.Error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	mAdmin, err := NewAdminClient(session)
	if err != nil {
		return nil
	}
	adminClient := AdminClient{Client: mAdmin}
	err = putBucket(ctx, adminClient, buchetName, location, opts)
	if err != nil {
		return prepareError(err)
	}
	return nil
}

// putBucket
func putBucket(ctx context.Context, client Admin, buchetName, location string, bool bool) error {
	err := client.putBucket(ctx, buchetName, location, bool)
	if err != nil {
		return err
	}
	return err
}

// getDeleteBucketResponse performs removeBucket() to delete a bucket
func getDeleteBucketResponse(session *models.Principal, buchetName string) *models.Error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	mAdmin, err := NewAdminClient(session)
	if err != nil {
		return nil
	}
	adminClient := AdminClient{Client: mAdmin}
	err = removeBucket(ctx, adminClient, buchetName)
	if err != nil {
		return prepareError(err)
	}
	return nil
}

// removeBucket deletes a bucket
func removeBucket(ctx context.Context, client Admin, bucketName string) error {
	return client.removeBucket(ctx, bucketName, "", false)
}

// getBucketSetPolicyResponse calls setBucketAccessPolicy() to set a access policy to a bucket
//   and returns the serialized output.
func getBucketSetPolicyResponse(session *models.Principal, bucketName string, req *models.SetBucketPolicyRequest) (*models.Bucket, *models.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	mAdmin, err := NewAdminClient(session)
	if err != nil {
		return nil, prepareError(err)
	}
	adminClient := AdminClient{Client: mAdmin}

	if err := setBucketAccessPolicy(ctx, adminClient, bucketName, *req.Access, req.Definition); err != nil {
		return nil, prepareError(err)
	}
	// set bucket access policy
	//bucket, err := getBucketInfo(ctx, adminClient, adminClient, bucketName)
	if err != nil {
		return nil, prepareError(err)
	}
	return nil, nil
}

// setBucketAccessPolicy set the access permissions on an existing bucket.
func setBucketAccessPolicy(ctx context.Context, client Admin, bucketName string, access models.BucketAccess, policyDefinition string) error {
	if strings.TrimSpace(bucketName) == "" {
		return fmt.Errorf("error: bucket name not present")
	}
	if strings.TrimSpace(string(access)) == "" {
		return fmt.Errorf("error: bucket access not present")
	}
	// Prepare policyJSON corresponding to the access type
	if access != models.BucketAccessPRIVATE && access != models.BucketAccessPUBLIC && access != models.BucketAccessCUSTOM {
		return fmt.Errorf("access: `%s` not supported", access)
	}
	if access == models.BucketAccessCUSTOM {
		err := client.putBucketPolicy(ctx, bucketName, policyDefinition)
		if err != nil {
			return err
		}
	}
	return nil
}

// getBucketPolicyResponse
func getBucketPolicyResponse(session *models.Principal, bucketName string) (*policy.Policy, *models.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	mAdmin, err := NewAdminClient(session)
	if err != nil {
		return nil, nil
	}
	adminClient := AdminClient{Client: mAdmin}
	policy, err := getBucketAccessPolicy(ctx, adminClient, bucketName)
	if err != nil {
		return nil, prepareError(err)
	}
	return policy, nil
}

// getBucketAccessPolicy
func getBucketAccessPolicy(ctx context.Context, client Admin, bucketName string) (*policy.Policy, error) {
	if strings.TrimSpace(bucketName) == "" {
		return nil, fmt.Errorf("error: bucket name not present")
	}
	policy, err := client.getBucketPolicy(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

// removeBucketPolicyResponse
func removeBucketPolicyResponse(session *models.Principal, bucketName string) *models.Error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	mAdmin, err := NewAdminClient(session)
	if err != nil {
		return nil
	}
	adminClient := AdminClient{Client: mAdmin}
	err = removeBucketAccessPolicy(ctx, adminClient, bucketName)
	if err != nil {
		return prepareError(err)
	}
	return nil
}

// removeBucketAccessPolicy
func removeBucketAccessPolicy(ctx context.Context, client Admin, bucketName string) error {
	if strings.TrimSpace(bucketName) == "" {
		return fmt.Errorf("error: bucket name not present")
	}
	err := client.removeBucketPolicy(ctx, bucketName)
	if err != nil {
		return err
	}
	return nil
}
