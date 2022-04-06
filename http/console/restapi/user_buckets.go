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

//func registerBucketsHandlers(api *operations.ConsoleAPI) {
//	// list buckets
//	api.UserAPIListBucketsHandler = user_api.ListBucketsHandlerFunc(func(params user_api.ListBucketsParams, session *models.Principal) middleware.Responder {
//		listBucketsResponse, err := getListBucketsResponse(session)
//		if err != nil {
//			return user_api.NewListBucketsDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewListBucketsOK().WithPayload(listBucketsResponse)
//	})
//	// create bucket
//	api.UserAPIMakeBucketHandler = user_api.MakeBucketHandlerFunc(func(params user_api.MakeBucketParams, session *models.Principal) middleware.Responder {
//		if err := getMakeBucketResponse(session, params.Body); err != nil {
//			return user_api.NewMakeBucketDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewMakeBucketCreated()
//	})
//	// delete bucket
//	api.UserAPIDeleteBucketHandler = user_api.DeleteBucketHandlerFunc(func(params user_api.DeleteBucketParams, session *models.Principal) middleware.Responder {
//		if err := getDeleteBucketResponse(session, params); err != nil {
//			return user_api.NewMakeBucketDefault(int(err.Code)).WithPayload(err)
//
//		}
//		return user_api.NewDeleteBucketNoContent()
//	})
//	// get bucket info
//	api.UserAPIBucketInfoHandler = user_api.BucketInfoHandlerFunc(func(params user_api.BucketInfoParams, session *models.Principal) middleware.Responder {
//		bucketInfoResp, err := getBucketInfoResponse(session, params)
//		if err != nil {
//			return user_api.NewBucketInfoDefault(int(err.Code)).WithPayload(err)
//		}
//
//		return user_api.NewBucketInfoOK().WithPayload(bucketInfoResp)
//	})
//	// set bucket policy
//	api.UserAPIBucketSetPolicyHandler = user_api.BucketSetPolicyHandlerFunc(func(params user_api.BucketSetPolicyParams, session *models.Principal) middleware.Responder {
//		bucketSetPolicyResp, err := getBucketSetPolicyResponse(session, params.Name, params.Body)
//		if err != nil {
//			return user_api.NewBucketSetPolicyDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewBucketSetPolicyOK().WithPayload(bucketSetPolicyResp)
//	})
//	// set bucket tags
//	api.UserAPIPutBucketTagsHandler = user_api.PutBucketTagsHandlerFunc(func(params user_api.PutBucketTagsParams, session *models.Principal) middleware.Responder {
//		err := getPutBucketTagsResponse(session, params.BucketName, params.Body)
//		if err != nil {
//			return user_api.NewPutBucketTagsDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewPutBucketTagsOK()
//	})
//	// get bucket versioning
//	api.UserAPIGetBucketVersioningHandler = user_api.GetBucketVersioningHandlerFunc(func(params user_api.GetBucketVersioningParams, session *models.Principal) middleware.Responder {
//		getBucketVersioning, err := getBucketVersionedResponse(session, params.BucketName)
//		if err != nil {
//			return user_api.NewGetBucketVersioningDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
//		}
//		return user_api.NewGetBucketVersioningOK().WithPayload(getBucketVersioning)
//	})
//	// update bucket versioning
//	api.UserAPISetBucketVersioningHandler = user_api.SetBucketVersioningHandlerFunc(func(params user_api.SetBucketVersioningParams, session *models.Principal) middleware.Responder {
//		err := setBucketVersioningResponse(session, params.BucketName, &params)
//		if err != nil {
//			return user_api.NewSetBucketVersioningDefault(500).WithPayload(err)
//		}
//		return user_api.NewSetBucketVersioningCreated()
//	})
//	// get bucket replication
//	api.UserAPIGetBucketReplicationHandler = user_api.GetBucketReplicationHandlerFunc(func(params user_api.GetBucketReplicationParams, session *models.Principal) middleware.Responder {
//		getBucketReplication, err := getBucketReplicationResponse(session, params.BucketName)
//		if err != nil {
//			return user_api.NewGetBucketReplicationDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
//		}
//		return user_api.NewGetBucketReplicationOK().WithPayload(getBucketReplication)
//	})
//	// get single bucket replication rule
//	api.UserAPIGetBucketReplicationRuleHandler = user_api.GetBucketReplicationRuleHandlerFunc(func(params user_api.GetBucketReplicationRuleParams, session *models.Principal) middleware.Responder {
//		getBucketReplicationRule, err := getBucketReplicationRuleResponse(session, params.BucketName, params.RuleID)
//		if err != nil {
//			return user_api.NewGetBucketReplicationRuleDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
//		}
//		return user_api.NewGetBucketReplicationRuleOK().WithPayload(getBucketReplicationRule)
//	})
//
//	// enable bucket encryption
//	api.UserAPIEnableBucketEncryptionHandler = user_api.EnableBucketEncryptionHandlerFunc(func(params user_api.EnableBucketEncryptionParams, session *models.Principal) middleware.Responder {
//		if err := enableBucketEncryptionResponse(session, params); err != nil {
//			return user_api.NewEnableBucketEncryptionDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewEnableBucketEncryptionOK()
//	})
//}

// getListBucketsResponse
func getListBucketsResponse(session *models.Principal) (*models.ListBucketsResponse, *models.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	mAdmin, err := NewMinioAdminClient(session)
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
func getlistBuckets(ctx context.Context, client MinioAdmin) ([]*models.Bucket, error) {
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

	mAdmin, err := NewMinioAdminClient(session)
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
func putBucket(ctx context.Context, client MinioAdmin, buchetName, location string, bool bool) error {
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

	mAdmin, err := NewMinioAdminClient(session)
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
func removeBucket(ctx context.Context, client MinioAdmin, bucketName string) error {
	return client.removeBucket(ctx, bucketName, "", false)
}

// getBucketSetPolicyResponse calls setBucketAccessPolicy() to set a access policy to a bucket
//   and returns the serialized output.
func getBucketSetPolicyResponse(session *models.Principal, bucketName string, req *models.SetBucketPolicyRequest) (*models.Bucket, *models.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	mAdmin, err := NewMinioAdminClient(session)
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
func setBucketAccessPolicy(ctx context.Context, client MinioAdmin, bucketName string, access models.BucketAccess, policyDefinition string) error {
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
	mAdmin, err := NewMinioAdminClient(session)
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
func getBucketAccessPolicy(ctx context.Context, client MinioAdmin, bucketName string) (*policy.Policy, error) {
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
	mAdmin, err := NewMinioAdminClient(session)
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
func removeBucketAccessPolicy(ctx context.Context, client MinioAdmin, bucketName string) error {
	if strings.TrimSpace(bucketName) == "" {
		return fmt.Errorf("error: bucket name not present")
	}
	err := client.removeBucketPolicy(ctx, bucketName)
	if err != nil {
		return err
	}
	return nil
}
