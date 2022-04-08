package restapi

import (
	"context"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/http/console/credentials"
	"github.com/filedag-project/filedag-storage/http/console/madmin"
	"github.com/filedag-project/filedag-storage/http/console/madmin/object"
	"github.com/filedag-project/filedag-storage/http/console/madmin/policy"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"io"
	"net/http"
)

const globalAppName = "Console"

// Admin interface with all functions to be implemented
type Admin interface {
	listUsers(ctx context.Context) ([]*iam.User, error)
	addUser(ctx context.Context, accessKey, SecretKey string) (*iam.User, error)
	removeUser(ctx context.Context, accessKey string) error
	getUserInfo(ctx context.Context, accessKey string) (*madmin.UserInfo, error)
	setUserStatus(ctx context.Context, accessKey string, status madmin.AccountStatus) error
	listBucketsInfo(ctx context.Context) ([]*s3.Bucket, error)
	putBucket(ctx context.Context, bucketName, location string, objectLocking bool) error
	removeBucket(ctx context.Context, bucketName, location string, objectLocking bool) error
	putObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64) error
	getObject(ctx context.Context, bucketName, objectName string) error
	removeObject(ctx context.Context, bucketName, objectName string) error
	copyObject(ctx context.Context, bucketName, objectName string) error
	listObject(ctx context.Context, bucketName string) ([]object.Object, error)
	putBucketPolicy(ctx context.Context, bucketName, policyStr string) error
	getBucketPolicy(ctx context.Context, bucketName string) (*policy.Policy, error)
	removeBucketPolicy(ctx context.Context, bucketName string) error
	putUserPolicy(ctx context.Context, userName, policyName, policyStr string) error
	getUserPolicy(ctx context.Context, userName string) (*madmin.UserPolicy, error)
	removeUserPolicy(ctx context.Context, userName, policyName string) error
	listUserPolicy(ctx context.Context, userName string) (*madmin.UserPolicies, error)
}

type AdminClient struct {
	Client *madmin.AdminClient
}

func (ac AdminClient) changePassword(ctx context.Context, accessKey, secretKey string) (*iam.User, error) {
	return ac.Client.SetUser(ctx, accessKey, secretKey, madmin.AccountEnabled)
}

// implements madmin.ListUsers()
func (ac AdminClient) listUsers(ctx context.Context) ([]*iam.User, error) {
	return ac.Client.ListUsers(ctx)
}

// implements madmin.AddUser()
func (ac AdminClient) addUser(ctx context.Context, accessKey, secretKey string) (*iam.User, error) {
	return ac.Client.AddUser(ctx, accessKey, secretKey)
}

// implements madmin.RemoveUser()
func (ac AdminClient) removeUser(ctx context.Context, accessKey string) error {
	return ac.Client.RemoveUser(ctx, accessKey)
}

//implements madmin.GetUserInfo()
func (ac AdminClient) getUserInfo(ctx context.Context, accessKey string) (*madmin.UserInfo, error) {
	return ac.Client.GetUserInfo(ctx, accessKey)
}

// implements madmin.SetUserStatus()
func (ac AdminClient) setUserStatus(ctx context.Context, accessKey string, status madmin.AccountStatus) error {
	return ac.Client.SetUserStatus(ctx, accessKey, status)
}

// implements madmin.ListServiceAccounts()
func (ac AdminClient) listServiceAccounts(ctx context.Context, user string) (madmin.ListServiceAccountsResp, error) {
	// TODO: Fix this
	return ac.Client.ListServiceAccounts(ctx, user)
}

// implements madmin.DeleteServiceAccount()
func (ac AdminClient) deleteServiceAccount(ctx context.Context, serviceAccount string) error {
	return ac.Client.DeleteServiceAccount(ctx, serviceAccount)
}

// implements madmin.InfoServiceAccount()
func (ac AdminClient) infoServiceAccount(ctx context.Context, serviceAccount string) (madmin.InfoServiceAccountResp, error) {
	return ac.Client.InfoServiceAccount(ctx, serviceAccount)
}

// implements madmin.UpdateServiceAccount()
func (ac AdminClient) updateServiceAccount(ctx context.Context, serviceAccount string, opts madmin.UpdateServiceAccountReq) error {
	return ac.Client.UpdateServiceAccount(ctx, serviceAccount, opts)
}

// AccountInfo implements madmin.AccountInfo()
func (ac AdminClient) listBucketsInfo(ctx context.Context) ([]*s3.Bucket, error) {
	return ac.Client.ListBucketsInfo(ctx, madmin.AccountOpts{})
}

// implements minio.MakeBucketWithContext(ctx, bucketName, location, objectLocking)
func (ac AdminClient) putBucket(ctx context.Context, bucketName, location string, objectLocking bool) error {
	return ac.Client.PutBucket(ctx, bucketName, madmin.AccountOpts{})
}

// implements minio.RemoveBucket(ctx, bucketName, location, objectLocking)
func (ac AdminClient) removeBucket(ctx context.Context, bucketName, location string, objectLocking bool) error {
	return ac.Client.RemoveBucket(ctx, bucketName, madmin.AccountOpts{})
}

func (ac AdminClient) putObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64) error {
	return ac.Client.PutObject(ctx, bucketName, objectName, reader, objectSize)
}

func (ac AdminClient) getObject(ctx context.Context, bucketName, objectName string) error {
	return ac.Client.GetObject(ctx, bucketName, objectName)
}

func (ac AdminClient) removeObject(ctx context.Context, bucketName, objectName string) error {
	return ac.Client.RemoveObject(ctx, bucketName, objectName)
}

func (ac AdminClient) copyObject(ctx context.Context, bucketName, objectName string) error {
	return ac.Client.CopyObject(ctx, bucketName, objectName)
}

func (ac AdminClient) headObject(ctx context.Context, bucketName, objectName string) error {
	return ac.Client.HeadObject(ctx, bucketName, objectName)
}

func (ac AdminClient) listObject(ctx context.Context, bucketName string) ([]object.Object, error) {
	return ac.Client.ListObject(ctx, bucketName)
}

func (ac AdminClient) putBucketPolicy(ctx context.Context, bucketName, policyStr string) error {
	return ac.Client.PutBucketPolicy(ctx, bucketName, policyStr)
}

func (ac AdminClient) getBucketPolicy(ctx context.Context, bucketName string) (*policy.Policy, error) {
	return ac.Client.GetBucketPolicy(ctx, bucketName)
}

func (ac AdminClient) removeBucketPolicy(ctx context.Context, bucketName string) error {
	return ac.Client.RemoveBucketPolicy(ctx, bucketName)
}

func (ac AdminClient) putUserPolicy(ctx context.Context, bucketName, policyName, policyStr string) error {
	return ac.Client.PutUserPolicy(ctx, bucketName, policyName, policyStr)
}

func (ac AdminClient) getUserPolicy(ctx context.Context, bucketName string) (*madmin.UserPolicy, error) {
	return ac.Client.GetUserPolicy(ctx, bucketName)
}

func (ac AdminClient) removeUserPolicy(ctx context.Context, bucketName, policyName string) error {
	return ac.Client.RemoveUserPolicy(ctx, bucketName, policyName)
}

func (ac AdminClient) listUserPolicy(ctx context.Context, bucketName string) (*madmin.UserPolicies, error) {
	return ac.Client.ListUserPolicy(ctx, bucketName)
}

func NewAdminClient(sessionClaims *models.Principal) (*madmin.AdminClient, error) {
	adminClient, err := newAdminFromClaims(sessionClaims)
	if err != nil {
		return nil, err
	}
	return adminClient, nil
}

// newAdminFromClaims creates a minio admin from Decrypted claims using Assume role credentials
func newAdminFromClaims(claims *models.Principal) (*madmin.AdminClient, error) {
	tlsEnabled := getEndpointIsSecure()
	endpoint := getEndpoint()

	adminClient, err := madmin.NewWithOptions(endpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(claims.STSAccessKeyID, claims.STSSecretAccessKey, claims.STSSessionToken),
		Secure: tlsEnabled,
	})
	if err != nil {
		return nil, err
	}
	adminClient.SetCustomTransport(GetConsoleHTTPClient().Transport)
	return adminClient, nil
}

// newAdminFromCreds Creates a minio client using custom credentials for connecting to a remote host
func newAdminFromCreds(accessKey, secretKey, endpoint string, tlsEnabled bool) (*madmin.AdminClient, error) {
	minioClient, err := madmin.NewWithOptions(endpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: tlsEnabled,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

// httpClient is a custom http client, this client should not be called directly and instead be
// called using GetConsoleHTTPClient() to ensure is initialized and the certificates are loaded correctly
var httpClient *http.Client

// GetConsoleHTTPClient will initialize the console HTTP Client with fully populated custom TLS
// Transport that with loads certs at
// - ${HOME}/.console/certs/CAs
func GetConsoleHTTPClient() *http.Client {
	if httpClient == nil {
		httpClient = PrepareConsoleHTTPClient(false)
	}
	return httpClient
}
