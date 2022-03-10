package iam

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/etag"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/hash"
	"io"
	"io/ioutil"
	"net/http"
)

// CheckRequestAuthType Check request auth type verifies the incoming http request
// - validates the request signature
// - validates the policy action if anonymous tests bucket policies if any,
//   for authenticated requests validates IAM policies.
// returns APIErrorCode if any to be replied to the client.
func CheckRequestAuthType(ctx context.Context, r *http.Request, action s3action.Action, bucketName, objectName string) (s3Err api_errors.ErrorCode) {
	_, _, s3Err = checkRequestAuthTypeCredential(ctx, r, action, bucketName, objectName)
	return s3Err
}

// Check request auth type verifies the incoming http request
// - validates the request signature
// - validates the policy action if anonymous tests bucket policies if any,
//   for authenticated requests validates IAM policies.
// returns APIErrorCode if any to be replied to the client.
// Additionally returns the accessKey used in the request, and if this request is by an admin.
func checkRequestAuthTypeCredential(ctx context.Context, r *http.Request, action s3action.Action, bucketName, objectName string) (cred auth.Credentials, owner bool, s3Err api_errors.ErrorCode) {
	switch getRequestAuthType(r) {
	case authTypeUnknown, authTypeStreamingSigned:
		return cred, owner, api_errors.ErrSignatureVersionNotSupported
	case authTypePresignedV2, authTypeSignedV2:
		if s3Err = isReqAuthenticatedV2(r); s3Err != api_errors.ErrNone {
			return cred, owner, s3Err
		}
		cred, owner, s3Err = getReqAccessKeyV2(r)
	case authTypeSigned, authTypePresigned:
		region := ""
		switch action {
		case policy.GetBucketLocationAction, s3action.ListAllMyBucketsAction:
			region = ""
		}
		if s3Err = IsReqAuthenticated(ctx, r, region, serviceS3); s3Err != api_errors.ErrNone {
			return cred, owner, s3Err
		}
		cred, owner, s3Err = GetReqAccessKeyV4(r, region, serviceS3)
	}
	if s3Err != api_errors.ErrNone {
		return cred, owner, s3Err
	}

	if action == s3action.CreateBucketAction {
		// To extract region from XML in request body, get copy of request body.
		payload, err := ioutil.ReadAll(io.LimitReader(r.Body, consts.MaxLocationConstraintSize))
		if err != nil {
			log.Errorf("ReadAll err:%v", err)
			return cred, owner, api_errors.ErrMalformedXML
		}

		// Populate payload to extract location constraint.
		r.Body = ioutil.NopCloser(bytes.NewReader(payload))

		// Populate payload again to handle it in HTTP handler.
		r.Body = ioutil.NopCloser(bytes.NewReader(payload))
	}

	if action != s3action.ListAllMyBucketsAction && cred.AccessKey == "" {
		// Anonymous checks are not meant for ListBuckets action
		if policy.GlobalPolicySys.IsAllowed(policy.Args{
			AccountName: cred.AccessKey,
			Action:      action,
			BucketName:  bucketName,
			IsOwner:     false,
			ObjectName:  objectName,
		}) {
			// Request is allowed return the appropriate access key.
			return cred, owner, api_errors.ErrNone
		}

		if action == s3action.ListBucketVersionsAction {
			// In AWS S3 s3:ListBucket permission is same as s3:ListBucketVersions permission
			// verify as a fallback.
			if policy.GlobalPolicySys.IsAllowed(policy.Args{
				AccountName: cred.AccessKey,
				Action:      s3action.ListBucketAction,
				BucketName:  bucketName,
				IsOwner:     false,
				ObjectName:  objectName,
			}) {
				// Request is allowed return the appropriate access key.
				return cred, owner, api_errors.ErrNone
			}
		}

		return cred, owner, api_errors.ErrAccessDenied
	}

	if GlobalIAMSys.IsAllowed(policy.Args{
		AccountName: cred.AccessKey,
		Action:      action,
		BucketName:  bucketName,
		ObjectName:  objectName,
		IsOwner:     owner,
	}) {
		// Request is allowed return the appropriate access key.
		return cred, owner, api_errors.ErrNone
	}

	if action == s3action.ListBucketVersionsAction {
		// In AWS S3 s3:ListBucket permission is same as s3:ListBucketVersions permission
		// verify as a fallback.
		if GlobalIAMSys.IsAllowed(policy.Args{
			AccountName: cred.AccessKey,
			Action:      s3action.ListBucketAction,
			BucketName:  bucketName,
			ObjectName:  objectName,
			IsOwner:     owner,
		}) {
			// Request is allowed return the appropriate access key.
			return cred, owner, api_errors.ErrNone
		}
	}

	return cred, owner, api_errors.ErrAccessDenied
}

// Verify if request has valid AWS Signature Version '2'.
func isReqAuthenticatedV2(r *http.Request) (s3Error api_errors.ErrorCode) {
	if isRequestSignatureV2(r) {
		return doesSignV2Match(r)
	}
	return doesPresignV2SignatureMatch(r)
}

func reqSignatureV4Verify(r *http.Request, region string, stype serviceType) (s3Error api_errors.ErrorCode) {
	sha256sum := getContentSha256Cksum(r, stype)
	switch {
	case IsRequestSignatureV4(r):
		return doesSignatureMatch(sha256sum, r, region, stype)
	case isRequestPresignedSignatureV4(r):
		return doesPresignedSignatureMatch(sha256sum, r, region, stype)
	default:
		return api_errors.ErrAccessDenied
	}
}

// Verify if request has valid AWS Signature Version '4'.
func IsReqAuthenticated(ctx context.Context, r *http.Request, region string, stype serviceType) (s3Error api_errors.ErrorCode) {
	if errCode := reqSignatureV4Verify(r, region, stype); errCode != api_errors.ErrNone {
		return errCode
	}
	clientETag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		return api_errors.ErrInvalidDigest
	}

	// Extract either 'X-Amz-Content-Sha256' header or 'X-Amz-Content-Sha256' query parameter (if V4 presigned)
	// Do not verify 'X-Amz-Content-Sha256' if skipSHA256.
	var contentSHA256 []byte
	if skipSHA256 := skipContentSha256Cksum(r); !skipSHA256 && isRequestPresignedSignatureV4(r) {
		if sha256Sum, ok := r.Form[consts.AmzContentSha256]; ok && len(sha256Sum) > 0 {
			contentSHA256, err = hex.DecodeString(sha256Sum[0])
			if err != nil {
				return api_errors.ErrContentSHA256Mismatch
			}
		}
	} else if _, ok := r.Header[consts.AmzContentSha256]; !skipSHA256 && ok {
		contentSHA256, err = hex.DecodeString(r.Header.Get(consts.AmzContentSha256))
		if err != nil || len(contentSHA256) == 0 {
			return api_errors.ErrContentSHA256Mismatch
		}
	}

	// Verify 'Content-Md5' and/or 'X-Amz-Content-Sha256' if present.
	// The verification happens implicit during reading.
	reader, err := hash.NewReader(r.Body, -1, clientETag.String(), hex.EncodeToString(contentSHA256), -1)
	if err != nil {
		return api_errors.ErrReader
	}
	r.Body = reader
	return api_errors.ErrNone
}

//ValidateAdminSignature validate admin Signature
func ValidateAdminSignature(ctx context.Context, r *http.Request, region string) (auth.Credentials, map[string]interface{}, bool, api_errors.ErrorCode) {
	var cred auth.Credentials
	var owner bool
	s3Err := api_errors.ErrAccessDenied
	if _, ok := r.Header[consts.AmzContentSha256]; ok &&
		getRequestAuthType(r) == authTypeSigned {
		// We only support admin credentials to access admin APIs.
		cred, owner, s3Err = GetReqAccessKeyV4(r, region, serviceS3)
		if s3Err != api_errors.ErrNone {
			return cred, nil, owner, s3Err
		}

		// we only support V4 (no presign) with auth body
		s3Err = IsReqAuthenticated(ctx, r, region, serviceS3)
	}
	if s3Err != api_errors.ErrNone {
		return cred, nil, owner, s3Err
	}

	return cred, nil, owner, api_errors.ErrNone
}