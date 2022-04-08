package restapi

import (
	"context"
	"encoding/base64"
	"github.com/filedag-project/filedag-storage/http/console/madmin/object"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/console/restapi/operations/user_api"
	"io"
	"path"
	"strconv"
	"time"
)

// GetListObjectsResponse returns a list of objects
func (apiServer *ApiServer) GetListObjectsResponse(session *models.Principal, params models.ListObjectsParams) (*models.ListObjectsResponse, *models.Error) {
	var prefix string
	var recursive bool
	var withVersions bool
	var withMetadata bool
	if params.Prefix != nil {
		encodedPrefix := SanitizeEncodedPrefix(*params.Prefix)
		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
		if err != nil {
			return nil, prepareError(err)
		}
		prefix = string(decodedPrefix)
	}
	if params.Recursive != nil {
		recursive = *params.Recursive
	}
	if params.WithVersions != nil {
		withVersions = *params.WithVersions
	}
	if params.WithMetadata != nil {
		withMetadata = *params.WithMetadata
	}
	// bucket request needed to proceed
	if params.BucketName == "" {
		return nil, prepareError(errBucketNameNotInRequest)
	}
	mAdmin, err := NewAdminClient(session)
	if err != nil {
		return nil, prepareError(err)
	}
	// create a minioClient interface implementation
	// defining the client to be used
	adminClient := AdminClient{Client: mAdmin}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	objs, err := listBucketObjects(ctx, adminClient, params.BucketName, prefix, recursive, withVersions, withMetadata)
	if err != nil {
		return nil, prepareError(err)
	}

	resp := &models.ListObjectsResponse{
		Objects: objs,
		Total:   int64(len(objs)),
	}
	return resp, nil
}

// listBucketObjects gets an array of objects in a bucket
func listBucketObjects(ctx context.Context, client AdminClient, bucketName string, prefix string, recursive, withVersions bool, withMetadata bool) ([]*object.Object, error) {
	var objects []*object.Object
	//opts := minio.ListObjectsOptions{
	//	Prefix:       prefix,
	//	Recursive:    recursive,
	//	WithVersions: withVersions,
	//	WithMetadata: withMetadata,
	//}
	//if withMetadata {
	//	opts.MaxKeys = 1
	//}
	listObjects, err := client.listObject(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	for _, lsObj := range listObjects {
		objects = append(objects, &lsObj)
	}
	return objects, nil
}

func getUploadObjectResponse(session *models.Principal, params user_api.PostBucketsBucketNameObjectsUploadParams) *models.Error {
	ctx := context.Background()
	mClient, err := NewAdminClient(session)
	if err != nil {
		return prepareError(err)
	}
	client := AdminClient{Client: mClient}
	if err := uploadFiles(ctx, client, params); err != nil {
		return prepareError(err, ErrorGeneric)
	}
	return nil
}

// uploadFiles gets files from http.Request form and uploads them to MinIO
func uploadFiles(ctx context.Context, client AdminClient, params user_api.PostBucketsBucketNameObjectsUploadParams) error {
	var prefix string
	if params.Prefix != nil {
		encodedPrefix := SanitizeEncodedPrefix(*params.Prefix)
		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
		if err != nil {
			return err
		}
		prefix = string(decodedPrefix)
	}

	// parse a request body as multipart/form-data.
	// 32 << 20 is default max memory
	mr, err := params.HTTPRequest.MultipartReader()
	if err != nil {
		return err
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		size, err := strconv.ParseInt(p.FormName(), 10, 64)
		if err != nil {
			return err
		}

		//contentType := p.Header.Get("content-type")

		err = client.putObject(ctx, params.BucketName, path.Join(prefix, p.FileName()), p, size)

		if err != nil {
			return err
		}
	}

	return nil
}

//// getShareObjectResponse returns a share object url
//func getShareObjectResponse(session *models.Principal, params user_api.ShareObjectParams) (*string, *models.Error) {
//	ctx := context.Background()
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return nil, prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	s3Client, err := newS3BucketClient(session, params.BucketName, prefix)
//	if err != nil {
//		return nil, prepareError(err)
//	}
//	// create a mc S3Client interface implementation
//	// defining the client to be used
//	mcClient := mcClient{client: s3Client}
//	var expireDuration string
//	if params.Expires != nil {
//		expireDuration = *params.Expires
//	}
//	url, err := getShareObjectURL(ctx, mcClient, params.VersionID, expireDuration)
//	if err != nil {
//		return nil, prepareError(err)
//	}
//	return url, nil
//}
//
//func getShareObjectURL(ctx context.Context, client MCClient, versionID string, duration string) (url *string, err error) {
//	// default duration 7d if not defined
//	if strings.TrimSpace(duration) == "" {
//		duration = "168h"
//	}
//
//	expiresDuration, err := time.ParseDuration(duration)
//	if err != nil {
//		return nil, err
//	}
//	objURL, pErr := client.shareDownload(ctx, versionID, expiresDuration)
//	if pErr != nil {
//		return nil, pErr.Cause
//	}
//	return &objURL, nil
//}
//
//func getSetObjectLegalHoldResponse(session *models.Principal, params user_api.PutObjectLegalHoldParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	err = setObjectLegalHold(ctx, minioClient, params.BucketName, prefix, params.VersionID, *params.Body.Status)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func setObjectLegalHold(ctx context.Context, client MinioClient, bucketName, prefix, versionID string, status models.ObjectLegalHoldStatus) error {
//	var lstatus minio.LegalHoldStatus
//	if status == models.ObjectLegalHoldStatusEnabled {
//		lstatus = minio.LegalHoldEnabled
//	} else {
//		lstatus = minio.LegalHoldDisabled
//	}
//	return client.putObjectLegalHold(ctx, bucketName, prefix, minio.PutObjectLegalHoldOptions{VersionID: versionID, Status: &lstatus})
//}
//
//func getSetObjectRetentionResponse(session *models.Principal, params user_api.PutObjectRetentionParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	err = setObjectRetention(ctx, minioClient, params.BucketName, params.VersionID, prefix, params.Body)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func setObjectRetention(ctx context.Context, client MinioClient, bucketName, versionID, prefix string, retentionOps *models.PutObjectRetentionRequest) error {
//	if retentionOps == nil {
//		return errors.New("object retention options can't be nil")
//	}
//	if retentionOps.Expires == nil {
//		return errors.New("object retention expires can't be nil")
//	}
//
//	var mode minio.RetentionMode
//	if *retentionOps.Mode == models.ObjectRetentionModeGovernance {
//		mode = minio.Governance
//	} else {
//		mode = minio.Compliance
//	}
//	retentionUntilDate, err := time.Parse(time.RFC3339, *retentionOps.Expires)
//	if err != nil {
//		return err
//	}
//	opts := minio.PutObjectRetentionOptions{
//		GovernanceBypass: retentionOps.GovernanceBypass,
//		RetainUntilDate:  &retentionUntilDate,
//		Mode:             &mode,
//		VersionID:        versionID,
//	}
//	return client.putObjectRetention(ctx, bucketName, prefix, opts)
//}
//
//func deleteObjectRetentionResponse(session *models.Principal, params user_api.DeleteObjectRetentionParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	err = deleteObjectRetention(ctx, minioClient, params.BucketName, prefix, params.VersionID)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func deleteObjectRetention(ctx context.Context, client MinioClient, bucketName, prefix, versionID string) error {
//	opts := minio.PutObjectRetentionOptions{
//		GovernanceBypass: true,
//		VersionID:        versionID,
//	}
//
//	return client.putObjectRetention(ctx, bucketName, prefix, opts)
//}
//
//func getPutObjectTagsResponse(session *models.Principal, params user_api.PutObjectTagsParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	err = putObjectTags(ctx, minioClient, params.BucketName, prefix, params.VersionID, params.Body.Tags)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func putObjectTags(ctx context.Context, client MinioClient, bucketName, prefix, versionID string, tagMap map[string]string) error {
//	opt := minio.PutObjectTaggingOptions{
//		VersionID: versionID,
//	}
//	otags, err := tags.MapToObjectTags(tagMap)
//	if err != nil {
//		return err
//	}
//	return client.putObjectTagging(ctx, bucketName, prefix, otags, opt)
//}
//
//// Restore Object Version
//func getPutObjectRestoreResponse(session *models.Principal, params user_api.PutObjectRestoreParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//
//	err = restoreObject(ctx, minioClient, params.BucketName, prefix, params.VersionID)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func restoreObject(ctx context.Context, client MinioClient, bucketName, prefix, versionID string) error {
//	// Select required version
//	srcOpts := minio.CopySrcOptions{
//		Bucket:    bucketName,
//		Object:    prefix,
//		VersionID: versionID,
//	}
//
//	// Destination object, same as current bucket
//	replaceMetadata := make(map[string]string)
//	replaceMetadata["copy-source"] = versionID
//
//	dstOpts := minio.CopyDestOptions{
//		Bucket:       bucketName,
//		Object:       prefix,
//		UserMetadata: replaceMetadata,
//	}
//
//	// Copy object call
//	_, err := client.copyObject(ctx, dstOpts, srcOpts)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// Metadata Response from minio-go API
//func getObjectMetadataResponse(session *models.Principal, params user_api.GetObjectMetadataParams) (*models.Metadata, *models.Error) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return nil, prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return nil, prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//
//	objectInfo, err := getObjectInfo(ctx, minioClient, params.BucketName, prefix)
//
//	if err != nil {
//		return nil, prepareError(err)
//	}
//
//	metadata := &models.Metadata{ObjectMetadata: objectInfo.Metadata}
//
//	return metadata, nil
//}
//
//func getObjectInfo(ctx context.Context, client MinioClient, bucketName, prefix string) (minio.ObjectInfo, error) {
//	objectData, err := client.statObject(ctx, bucketName, prefix, minio.GetObjectOptions{})
//
//	if err != nil {
//		return minio.ObjectInfo{}, err
//	}
//
//	return objectData, nil
//}
//
//// newClientURL returns an abstracted URL for filesystems and object storage.
//func newClientURL(urlStr string) *mc.ClientURL {
//	scheme, rest := getScheme(urlStr)
//	if strings.HasPrefix(rest, "//") {
//		// if rest has '//' prefix, skip them
//		var authority string
//		authority, rest = splitSpecial(rest[2:], "/", false)
//		if rest == "" {
//			rest = "/"
//		}
//		host := getHost(authority)
//		if host != "" && (scheme == "http" || scheme == "https") {
//			return &mc.ClientURL{
//				Scheme:          scheme,
//				Type:            objectStorage,
//				Host:            host,
//				Path:            rest,
//				SchemeSeparator: "://",
//				Separator:       '/',
//			}
//		}
//	}
//	return &mc.ClientURL{
//		Type:      fileSystem,
//		Path:      rest,
//		Separator: filepath.Separator,
//	}
//}
//
//// Maybe rawurl is of the form scheme:path. (Scheme must be [a-zA-Z][a-zA-Z0-9+-.]*)
//// If so, return scheme, path; else return "", rawurl.
//func getScheme(rawurl string) (scheme, path string) {
//	urlSplits := strings.Split(rawurl, "://")
//	if len(urlSplits) == 2 {
//		scheme, uri := urlSplits[0], "//"+urlSplits[1]
//		// ignore numbers in scheme
//		validScheme := regexp.MustCompile("^[a-zA-Z]+$")
//		if uri != "" {
//			if validScheme.MatchString(scheme) {
//				return scheme, uri
//			}
//		}
//	}
//	return "", rawurl
//}
//
//// Assuming s is of the form [s delimiter s].
//// If so, return s, [delimiter]s or return s, s if cutdelimiter == true
//// If no delimiter found return s, "".
//func splitSpecial(s string, delimiter string, cutdelimiter bool) (string, string) {
//	i := strings.Index(s, delimiter)
//	if i < 0 {
//		// if delimiter not found return as is.
//		return s, ""
//	}
//	// if delimiter should be removed, remove it.
//	if cutdelimiter {
//		return s[0:i], s[i+len(delimiter):]
//	}
//	// return split strings with delimiter
//	return s[0:i], s[i:]
//}
//
//// getHost - extract host from authority string, we do not support ftp style username@ yet.
//func getHost(authority string) (host string) {
//	i := strings.LastIndex(authority, "@")
//	if i >= 0 {
//		// TODO support, username@password style userinfo, useful for ftp support.
//		return
//	}
//	return authority
//}
