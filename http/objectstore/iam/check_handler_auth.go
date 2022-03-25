package iam

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/etag"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// AuthSys auth and sign system
type AuthSys struct {
	Iam       IdentityAMSys
	PolicySys IPolicySys
}

//Init AuthSys
func (s *AuthSys) Init() {
	s.Iam.Init()
	s.PolicySys.Init()
}

// CheckRequestAuthTypeCredential Check request auth type verifies the incoming http request
// - validates the request signature
// - validates the policy action if anonymous tests bucket policies if any,
//   for authenticated requests validates IAM policies.
// returns APIErrorCode if any to be replied to the client.
// Additionally, returns the accessKey used in the request, and if this request is by an admin.
func (s *AuthSys) CheckRequestAuthTypeCredential(ctx context.Context, r *http.Request, action s3action.Action, bucketName, objectName string) (cred auth.Credentials, owner bool, s3Err api_errors.ErrorCode) {
	switch GetRequestAuthType(r) {
	case authTypeUnknown, AuthTypeStreamingSigned:
		return cred, owner, api_errors.ErrSignatureVersionNotSupported
	case authTypePresignedV2, authTypeSignedV2:
		if s3Err = s.isReqAuthenticatedV2(r); s3Err != api_errors.ErrNone {
			return cred, owner, s3Err
		}
		cred, owner, s3Err = s.getReqAccessKeyV2(r)
	case authTypeSigned, authTypePresigned:
		region := ""
		switch action {
		case s3action.GetBucketLocationAction, s3action.ListAllMyBucketsAction:
			region = ""
		}
		if s3Err = s.IsReqAuthenticated(ctx, r, region, serviceS3); s3Err != api_errors.ErrNone {
			return cred, owner, s3Err
		}
		cred, owner, s3Err = s.GetReqAccessKeyV4(r, region, serviceS3)
	}
	if s3Err != api_errors.ErrNone {
		return cred, owner, s3Err
	}
	cred, _ = s.Iam.GetUser(ctx, cred.ParentUser)
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
		pol, err := s.PolicySys.bmSys.GetConfig(bucketName, cred.AccessKey)
		if pol != (bucketMetadata{}) {
			return cred, owner, api_errors.ErrBucketAlreadyExists
		}
	}

	if action != s3action.ListAllMyBucketsAction && cred.AccessKey == "" {
		// Anonymous checks are not meant for ListBuckets action
		if s.PolicySys.IsAllowed(auth.Args{
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
			if s.PolicySys.IsAllowed(auth.Args{
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

	if s.Iam.IsAllowed(auth.Args{
		AccountName: cred.AccessKey,
		Action:      action,
		BucketName:  bucketName,
		Conditions:  getConditions(r, cred.AccessKey),
		ObjectName:  objectName,
		IsOwner:     owner,
	}) {
		// Request is allowed return the appropriate access key.
		return cred, owner, api_errors.ErrNone
	}

	if action == s3action.ListBucketVersionsAction {
		// In AWS S3 s3:ListBucket permission is same as s3:ListBucketVersions permission
		// verify as a fallback.
		if s.Iam.IsAllowed(auth.Args{
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
func (s *AuthSys) isReqAuthenticatedV2(r *http.Request) (s3Error api_errors.ErrorCode) {
	if isRequestSignatureV2(r) {
		return s.doesSignV2Match(r)
	}
	return s.doesPresignV2SignatureMatch(r)
}

func (s *AuthSys) reqSignatureV4Verify(r *http.Request, region string, stype serviceType) (s3Error api_errors.ErrorCode) {
	sha256sum := getContentSha256Cksum(r, stype)
	switch {
	case IsRequestSignatureV4(r):
		return s.doesSignatureMatch(sha256sum, r, region, stype)
	case isRequestPresignedSignatureV4(r):
		return s.doesPresignedSignatureMatch(sha256sum, r, region, stype)
	default:
		return api_errors.ErrAccessDenied
	}
}

// IsReqAuthenticated Verify if request has valid AWS Signature Version '4'.
func (s *AuthSys) IsReqAuthenticated(ctx context.Context, r *http.Request, region string, stype serviceType) (s3Error api_errors.ErrorCode) {
	if errCode := s.reqSignatureV4Verify(r, region, stype); errCode != api_errors.ErrNone {
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
func (s *AuthSys) ValidateAdminSignature(ctx context.Context, r *http.Request, region string) (auth.Credentials, map[string]interface{}, bool, api_errors.ErrorCode) {
	var cred auth.Credentials
	var owner bool
	s3Err := api_errors.ErrAccessDenied
	if _, ok := r.Header[consts.AmzContentSha256]; ok &&
		GetRequestAuthType(r) == authTypeSigned {
		// We only support admin credentials to access admin APIs.
		cred, owner, s3Err = s.GetReqAccessKeyV4(r, region, serviceS3)
		if s3Err != api_errors.ErrNone {
			return cred, nil, owner, s3Err
		}

		// we only support V4 (no presign) with auth body
		s3Err = s.IsReqAuthenticated(ctx, r, region, serviceS3)
	}
	if s3Err != api_errors.ErrNone {
		return cred, nil, owner, s3Err
	}

	return cred, nil, owner, api_errors.ErrNone
}

func getConditions(r *http.Request, username string) map[string][]string {
	currTime := time.Now().UTC()

	principalType := "Anonymous"
	if username != "" {
		principalType = "User"
	}

	at := GetRequestAuthType(r)
	var signatureVersion string
	switch at {
	case authTypeSignedV2, authTypePresignedV2:
		signatureVersion = signV2Algorithm
	case authTypeSigned, authTypePresigned, AuthTypeStreamingSigned, authTypePostPolicy:
		signatureVersion = signV4Algorithm
	}

	var authtype string
	switch at {
	case authTypePresignedV2, authTypePresigned:
		authtype = "REST-QUERY-STRING"
	case authTypeSignedV2, authTypeSigned, AuthTypeStreamingSigned:
		authtype = "REST-HEADER"
	case authTypePostPolicy:
		authtype = "POST"
	}

	args := map[string][]string{
		"CurrentTime":      {currTime.Format(time.RFC3339)},
		"EpochTime":        {strconv.FormatInt(currTime.Unix(), 10)},
		"SecureTransport":  {strconv.FormatBool(r.TLS != nil)},
		"UserAgent":        {r.UserAgent()},
		"Referer":          {r.Referer()},
		"principaltype":    {principalType},
		"userid":           {username},
		"username":         {username},
		"signatureversion": {signatureVersion},
		"AuthType":         {authtype},
	}

	cloneHeader := r.Header.Clone()

	for key, values := range cloneHeader {
		if existingValues, found := args[key]; found {
			args[key] = append(existingValues, values...)
		} else {
			args[key] = values
		}
	}

	cloneURLValues := make(url.Values, len(r.Form))
	for k, v := range r.Form {
		cloneURLValues[k] = v
	}

	for key, values := range cloneURLValues {
		if existingValues, found := args[key]; found {
			args[key] = append(existingValues, values...)
		} else {
			args[key] = values
		}
	}

	return args
}
