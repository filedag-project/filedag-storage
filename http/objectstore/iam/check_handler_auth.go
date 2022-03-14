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
		if s3Err = isReqAuthenticated(ctx, r, region, serviceS3); s3Err != api_errors.ErrNone {
			return cred, owner, s3Err
		}
		cred, owner, s3Err = getReqAccessKeyV4(r, region, serviceS3)
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
	case isRequestSignatureV4(r):
		return doesSignatureMatch(sha256sum, r, region, stype)
	case isRequestPresignedSignatureV4(r):
		return doesPresignedSignatureMatch(sha256sum, r, region, stype)
	default:
		return api_errors.ErrAccessDenied
	}
}

// Verify if request has valid AWS Signature Version '4'.
func isReqAuthenticated(ctx context.Context, r *http.Request, region string, stype serviceType) (s3Error api_errors.ErrorCode) {
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

////CheckAuth check auth right
//func (sys *iamSys) CheckAuth(f http.HandlerFunc, action s3action.Action) http.HandlerFunc {
//
//	return func(w http.ResponseWriter, r *http.Request) {
//		identity, errCode := sys.authRequest(r, action)
//		if errCode == api_errors.ErrNone {
//			if identity != nil && identity.Name != "" {
//				r.Header.Set(consts.AmzIdentityId, identity.Name)
//				if identity.isAdmin() {
//					r.Header.Set(consts.AmzIsAdmin, "true")
//				} else if _, ok := r.Header[consts.AmzIsAdmin]; ok {
//					r.Header.Del(consts.AmzIsAdmin)
//				}
//			}
//			f(w, r)
//			return
//		}
//		response.WriteErrorResponse(w, r, errCode)
//	}
//}
//func (sys *iamSys) lookupAnonymous() (identity *Identity, found bool) {
//	for _, ident := range sys.identities {
//		if ident.Name == "anonymous" {
//			return ident, true
//		}
//	}
//	return nil, false
//}
//
//type Identity struct {
//	Name        string
//	Credentials []*Credential
//	Actions     []s3action.Action
//}
//type Credential struct {
//	AccessKey string
//	SecretKey string
//}
//
//// check whether the request has valid access keys
//func (sys *iamSys) authRequest(r *http.Request, action s3action.Action) (*Identity, api_errors.ErrorCode) {
//	var identity *Identity
//	var s3Err api_errors.ErrorCode
//	var found bool
//	var authType string
//	switch getRequestAuthType(r) {
//	case authTypeStreamingSigned:
//		return identity, api_errors.ErrNone
//	case authTypeUnknown:
//		log.Infof("unknown auth type")
//		r.Header.Set(consts.AmzAuthType, "Unknown")
//		return identity, api_errors.ErrAccessDenied
//	case authTypePresignedV2, authTypeSignedV2:
//		log.Infof("v2 auth type")
//		identity, s3Err = sys.isReqAuthenticatedV2(r)
//		authType = "SigV2"
//	case authTypeSigned, authTypePresigned:
//		log.Infof("v4 auth type")
//		identity, s3Err = sys.reqSignatureV4Verify(r)
//		authType = "SigV4"
//	case authTypePostPolicy:
//		log.Infof("post policy auth type")
//		r.Header.Set(consts.AmzAuthType, "PostPolicy")
//		return identity, api_errors.ErrNone
//	case authTypeJWT:
//		log.Infof("jwt auth type")
//		r.Header.Set(consts.AmzAuthType, "Jwt")
//		return identity, api_errors.ErrNotImplemented
//	case authTypeAnonymous:
//		authType = "Anonymous"
//		identity, found = sys.lookupAnonymous()
//		if !found {
//			r.Header.Set(consts.AmzAuthType, authType)
//			return identity, api_errors.ErrAccessDenied
//		}
//	default:
//		return identity, api_errors.ErrNotImplemented
//	}
//
//	if len(authType) > 0 {
//		r.Header.Set(consts.AmzAuthType, authType)
//	}
//	if s3Err != api_errors.ErrNone {
//		return identity, s3Err
//	}
//
//	log.Infof("user name: %v actions: %v, action: %v", identity.Name, identity.Actions, action)
//
//	bucket, object := GetBucketAndObject(r)
//
//	if !identity.canDo(action, bucket, object) {
//		return identity, api_errors.ErrAccessDenied
//	}
//
//	return identity, api_errors.ErrNone
//
//}
//func GetBucketAndObject(r *http.Request) (bucket, object string) {
//	vars := mux.Vars(r)
//	bucket = vars["bucket"]
//	object = vars["object"]
//	if !strings.HasPrefix(object, "/") {
//		object = "/" + object
//	}
//
//	return
//}
//func (identity *Identity) canDo(action s3action.Action, bucket string, objectKey string) bool {
//	if identity.isAdmin() {
//		return true
//	}
//	for _, a := range identity.Actions {
//		if a == action {
//			return true
//		}
//	}
//	if bucket == "" {
//		return false
//	}
//	target := string(action) + ":" + bucket + objectKey
//	adminTarget := consts.ACTION_ADMIN + ":" + bucket + objectKey
//	limitedByBucket := string(action) + ":" + bucket
//	adminLimitedByBucket := consts.ACTION_ADMIN + ":" + bucket
//	for _, a := range identity.Actions {
//		act := string(a)
//		if strings.HasSuffix(act, "*") {
//			if strings.HasPrefix(target, act[:len(act)-1]) {
//				return true
//			}
//			if strings.HasPrefix(adminTarget, act[:len(act)-1]) {
//				return true
//			}
//		} else {
//			if act == limitedByBucket {
//				return true
//			}
//			if act == adminLimitedByBucket {
//				return true
//			}
//		}
//	}
//	return false
//}
//func (identity *Identity) isAdmin() bool {
//	for _, a := range identity.Actions {
//		if a == "Admin" {
//			return true
//		}
//	}
//	return false
//}
//func (sys *iamSys) lookupByAccessKey(accessKey string) (identity *Identity, cred *Credential, found bool) {
//
//	for _, ident := range sys.identities {
//		for _, cred := range ident.Credentials {
//			// println("checking", ident.Name, cred.AccessKey)
//			if cred.AccessKey == accessKey {
//				return ident, cred, true
//			}
//		}
//	}
//	log.Infof("could not find accessKey %s", accessKey)
//	return nil, nil, false
//}

//// setAuthHandler to validate authorization header for the incoming request.
//func setAuthHandler(h http.Handler) http.Handler {
//	// handler for validating incoming authorization headers.
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		aType := getRequestAuthType(r)
//		if aType == authTypeSigned || aType == authTypeSignedV2 || aType == authTypeStreamingSigned {
//			// Verify if date headers are set, if not reject the request
//			amzDate, errCode := parseAmzDateHeader(r)
//			if errCode != api_errors.ErrNone {
//				// All our internal APIs are sensitive towards Date
//				// header, for all requests where Date header is not
//				// present we will reject such clients.
//				writeErrorResponse(r.Context(), w, errorCodes.ToAPIErr(errCode), r.URL)
//				atomic.AddUint64(&globalHTTPStats.rejectedRequestsTime, 1)
//				return
//			}
//			// Verify if the request date header is shifted by less than globalMaxSkewTime parameter in the past
//			// or in the future, reject request otherwise.
//			curTime := UTCNow()
//			if curTime.Sub(amzDate) > globalMaxSkewTime || amzDate.Sub(curTime) > globalMaxSkewTime {
//				writeErrorResponse(r.Context(), w, errorCodes.ToAPIErr(ErrRequestTimeTooSkewed), r.URL)
//				atomic.AddUint64(&globalHTTPStats.rejectedRequestsTime, 1)
//				return
//			}
//		}
//		if isSupportedS3AuthType(aType) || aType == authTypeJWT || aType == authTypeSTS {
//			h.ServeHTTP(w, r)
//			return
//		}
//		writeErrorResponse(r.Context(), w, errorCodes.ToAPIErr(ErrSignatureVersionNotSupported), r.URL)
//		atomic.AddUint64(&globalHTTPStats.rejectedRequestsAuth, 1)
//	})
//}
