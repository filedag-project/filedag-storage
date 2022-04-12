package iam

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"net/http"
	"strings"
)

// Verify if request has JWT.
func isRequestJWT(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Authorization"), "Bearer")
}

// IsRequestSignatureV4 Verify if request has AWS Signature Version '4'.
func IsRequestSignatureV4(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Authorization"), signV4Algorithm)
}

// Verify if request has AWS Signature Version '2'.
func isRequestSignatureV2(r *http.Request) bool {
	return !strings.HasPrefix(r.Header.Get("Authorization"), signV4Algorithm) &&
		strings.HasPrefix(r.Header.Get("Authorization"), signV2Algorithm)
}

// Verify if request has AWS PreSign Version '4'.
func isRequestPresignedSignatureV4(r *http.Request) bool {
	_, ok := r.URL.Query()["X-Amz-Credential"]
	return ok
}

// Verify request has AWS PreSign Version '2'.
func isRequestPresignedSignatureV2(r *http.Request) bool {
	_, ok := r.URL.Query()["AWSAccessKeyId"]
	return ok
}

// Verify if request has AWS Post policy Signature Version '4'.
func isRequestPostPolicySignatureV4(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") &&
		r.Method == http.MethodPost
}

// Verify if the request has AWS Streaming Signature Version '4'. This is only valid for 'PUT' operation.
func isRequestSignStreamingV4(r *http.Request) bool {
	return r.Header.Get("x-amz-content-sha256") == consts.StreamingContentSHA256 &&
		r.Method == http.MethodPut
}

// AuthType Authorization type.
type AuthType int

// List of all supported auth types.
const (
	authTypeUnknown AuthType = iota
	authTypeAnonymous
	authTypePresigned
	authTypePresignedV2
	authTypePostPolicy
	AuthTypeStreamingSigned
	authTypeSigned
	authTypeSignedV2
	authTypeJWT
)

// GetRequestAuthType Get request authentication type.
func GetRequestAuthType(r *http.Request) AuthType {
	if isRequestSignatureV2(r) {
		return authTypeSignedV2
	} else if isRequestPresignedSignatureV2(r) {
		return authTypePresignedV2
	} else if isRequestSignStreamingV4(r) {
		return AuthTypeStreamingSigned
	} else if IsRequestSignatureV4(r) {
		return authTypeSigned
	} else if isRequestPresignedSignatureV4(r) {
		return authTypePresigned
	} else if isRequestJWT(r) {
		return authTypeJWT
	} else if isRequestPostPolicySignatureV4(r) {
		return authTypePostPolicy
	} else if _, ok := r.Header["Authorization"]; !ok {
		return authTypeAnonymous
	}
	return authTypeUnknown
}
