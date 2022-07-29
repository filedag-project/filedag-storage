package iam

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
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
	Iam       *IdentityAMSys
	PolicySys *iPolicySys
	AdminCred auth.Credentials
}

//NewAuthSys new an AuthSys
func NewAuthSys(db *uleveldb.ULevelDB, adminCred auth.Credentials) *AuthSys {
	return &AuthSys{
		Iam:       NewIdentityAMSys(db),
		PolicySys: newIPolicySys(db),
		AdminCred: adminCred,
	}
}

// CheckRequestAuthTypeCredential Check request auth type verifies the incoming http request
// - validates the request signature
// - validates the policy action if anonymous tests bucket policies if any,
//   for authenticated requests validates IAM policies.
// returns APIErrorCode if any to be replied to the client.
// Additionally, returns the accessKey used in the request, and if this request is by an admin.
func (s *AuthSys) CheckRequestAuthTypeCredential(ctx context.Context, r *http.Request, action s3action.Action, bucketName, objectName string) (cred auth.Credentials, owner bool, s3Err api_errors.ErrorCode) {
	switch GetRequestAuthType(r) {
	case AuthTypeUnknown, AuthTypeStreamingSigned:
		return cred, owner, api_errors.ErrSignatureVersionNotSupported
	case AuthTypePresignedV2, AuthTypeSignedV2:
		if s3Err = s.IsReqAuthenticatedV2(r); s3Err != api_errors.ErrNone {
			return cred, owner, s3Err
		}
		cred, owner, s3Err = s.getReqAccessKeyV2(r)
	case AuthTypeSigned, AuthTypePresigned:
		region := ""
		switch action {
		case s3action.GetBucketLocationAction, s3action.ListAllMyBucketsAction:
			region = ""
		}
		if s3Err = s.IsReqAuthenticated(ctx, r, region, ServiceS3); s3Err != api_errors.ErrNone {
			return cred, owner, s3Err
		}
		cred, owner, s3Err = s.GetReqAccessKeyV4(r, region, ServiceS3)
	}
	if s3Err != api_errors.ErrNone {
		return cred, owner, s3Err
	}
	if cred.IsTemp() {
		cred, _ = s.Iam.GetUser(ctx, cred.ParentUser)
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
		ok := s.PolicySys.bmSys.HasBucket(bucketName, cred.AccessKey)
		if ok {
			return cred, owner, api_errors.ErrBucketAlreadyExists
		}
	}

	if action != s3action.ListAllMyBucketsAction && cred.AccessKey == "" {
		// Anonymous checks are not meant for ListBuckets action
		if s.PolicySys.isAllowed(auth.Args{
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
			if s.PolicySys.isAllowed(auth.Args{
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

	if s.Iam.IsAllowed(r.Context(), auth.Args{
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
		if s.Iam.IsAllowed(r.Context(), auth.Args{
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
func (s *AuthSys) IsReqAuthenticatedV2(r *http.Request) (s3Error api_errors.ErrorCode) {
	if isRequestSignatureV2(r) {
		return s.doesSignV2Match(r)
	}
	return s.doesPresignV2SignatureMatch(r)
}

func (s *AuthSys) ReqSignatureV4Verify(r *http.Request, region string, stype serviceType) (s3Error api_errors.ErrorCode) {
	sha256sum := GetContentSha256Cksum(r, stype)
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
	if errCode := s.ReqSignatureV4Verify(r, region, stype); errCode != api_errors.ErrNone {
		return errCode
	}
	clientETag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		return api_errors.ErrInvalidDigest
	}

	// Extract either 'X-Amz-Content-Sha256' header or 'X-Amz-Content-Sha256' query parameter (if V4 presigned)
	// Do not verify 'X-Amz-Content-Sha256' if skipSHA256.
	var contentSHA256 []byte
	if skipSHA256 := SkipContentSha256Cksum(r); !skipSHA256 && isRequestPresignedSignatureV4(r) {
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
		return api_errors.ErrInternalError
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
		GetRequestAuthType(r) == AuthTypeSigned {
		// We only support admin credentials to access admin APIs.
		cred, owner, s3Err = s.GetReqAccessKeyV4(r, region, ServiceS3)
		if s3Err != api_errors.ErrNone {
			return cred, nil, owner, s3Err
		}

		// we only support V4 (no presign) with auth body
		s3Err = s.IsReqAuthenticated(ctx, r, region, ServiceS3)
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
	case AuthTypeSignedV2, AuthTypePresignedV2:
		signatureVersion = signV2Algorithm
	case AuthTypeSigned, AuthTypePresigned, AuthTypeStreamingSigned, AuthTypePostPolicy:
		signatureVersion = signV4Algorithm
	}

	var authtype string
	switch at {
	case AuthTypePresignedV2, AuthTypePresigned:
		authtype = "REST-QUERY-STRING"
	case AuthTypeSignedV2, AuthTypeSigned, AuthTypeStreamingSigned:
		authtype = "REST-HEADER"
	case AuthTypePostPolicy:
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

// IsPutActionAllowed - check if PUT operation is allowed on the resource, this
// call verifies bucket policies and IAM policies, supports multi user
// checks etc.
func (s *AuthSys) IsPutActionAllowed(ctx context.Context, r *http.Request, action s3action.Action, bucketName, objectName string) (s3Err api_errors.ErrorCode) {
	var cred auth.Credentials
	var owner bool
	switch GetRequestAuthType(r) {
	case AuthTypeUnknown:
		return api_errors.ErrSignatureVersionNotSupported
	case AuthTypeSignedV2, AuthTypePresignedV2:
		cred, owner, s3Err = s.getReqAccessKeyV2(r)
	case AuthTypeStreamingSigned, AuthTypePresigned, AuthTypeSigned:
		region := ""
		cred, owner, s3Err = s.GetReqAccessKeyV4(r, region, ServiceS3)
	}
	if s3Err != api_errors.ErrNone {
		return s3Err
	}

	// Do not check for PutObjectRetentionAction permission,
	// if mode and retain until date are not set.
	// Can happen when bucket has default lock config set
	if action == s3action.PutObjectRetentionAction &&
		r.Header.Get(consts.AmzObjectLockMode) == "" &&
		r.Header.Get(consts.AmzObjectLockRetainUntilDate) == "" {
		return api_errors.ErrNone
	}

	if cred.AccessKey == "" {
		if s.PolicySys.isAllowed(auth.Args{
			AccountName: cred.AccessKey,
			Action:      action,
			BucketName:  bucketName,
			Conditions:  getConditions(r, ""),
			IsOwner:     false,
			ObjectName:  objectName,
		}) {
			return api_errors.ErrNone
		}
		return api_errors.ErrAccessDenied
	}

	if s.Iam.IsAllowed(ctx, auth.Args{
		AccountName: cred.AccessKey,
		Action:      action,
		BucketName:  bucketName,
		Conditions:  getConditions(r, cred.AccessKey),
		ObjectName:  objectName,
		IsOwner:     owner,
	}) {
		return api_errors.ErrNone
	}
	return api_errors.ErrAccessDenied
}

func (s *AuthSys) GetCredential(r *http.Request) (cred auth.Credentials, owner bool, s3Err api_errors.ErrorCode) {
	switch GetRequestAuthType(r) {
	case AuthTypeUnknown:
		s3Err = api_errors.ErrSignatureVersionNotSupported
	case AuthTypeSignedV2, AuthTypePresignedV2:
		cred, owner, s3Err = s.getReqAccessKeyV2(r)
	case AuthTypeStreamingSigned, AuthTypePresigned, AuthTypeSigned:
		region := ""
		cred, owner, s3Err = s.GetReqAccessKeyV4(r, region, ServiceS3)
	}
	return
}
