package iam

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"net/http"
	"time"
)

// SetAuthHandler to validate authorization header for the incoming request.
func SetAuthHandler(h http.Handler) http.Handler {
	// handler for validating incoming authorization headers.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aType := GetRequestAuthType(r)
		if aType == AuthTypeSigned || aType == AuthTypeSignedV2 || aType == AuthTypeStreamingSigned {
			// Verify if date headers are set, if not reject the request
			amzDate, errCode := parseAmzDateHeader(r)
			if errCode != api_errors.ErrNone {
				// All our internal APIs are sensitive towards Date
				// header, for all requests where Date header is not
				// present we will reject such clients.
				response.WriteErrorResponse(w, r, errCode)
				return
			}
			// Verify if the request date header is shifted by less than globalMaxSkewTime parameter in the past
			// or in the future, reject request otherwise.
			curTime := time.Now().UTC()
			if curTime.Sub(amzDate) > consts.GlobalMaxSkewTime || amzDate.Sub(curTime) > consts.GlobalMaxSkewTime {
				response.WriteErrorResponse(w, r, api_errors.ErrRequestTimeTooSkewed)
				return
			}
		}
		if isSupportedS3AuthType(aType) || aType == AuthTypeJWT || aType == AuthTypeSTS {
			h.ServeHTTP(w, r)
			return
		}
		response.WriteErrorResponse(w, r, api_errors.ErrSignatureVersionNotSupported)
	})
}

// Supported amz date formats.
var amzDateFormats = []string{
	// Do not change this order, x-amz-date format is usually in
	// iso8601Format rest are meant for relaxed handling of other
	// odd SDKs that might be out there.
	iso8601Format,
	time.RFC1123,
	time.RFC1123Z,
	// Add new AMZ date formats here.
}

// Supported Amz date headers.
var amzDateHeaders = []string{
	// Do not chane this order, x-amz-date value should be
	// validated first.
	"x-amz-date",
	"date",
}

// parseAmzDate - parses date string into supported amz date formats.
func parseAmzDate(amzDateStr string) (amzDate time.Time, apiErr api_errors.ErrorCode) {
	for _, dateFormat := range amzDateFormats {
		amzDate, err := time.Parse(dateFormat, amzDateStr)
		if err == nil {
			return amzDate, api_errors.ErrNone
		}
	}
	return time.Time{}, api_errors.ErrMalformedDate
}

// parseAmzDateHeader - parses supported amz date headers, in
// supported amz date formats.
func parseAmzDateHeader(req *http.Request) (time.Time, api_errors.ErrorCode) {
	for _, amzDateHeader := range amzDateHeaders {
		amzDateStr := req.Header.Get(amzDateHeader)
		if amzDateStr != "" {
			return parseAmzDate(amzDateStr)
		}
	}
	// Date header missing.
	return time.Time{}, api_errors.ErrMissingDateHeader
}

// List of all support S3 auth types.
var supportedS3AuthTypes = map[AuthType]struct{}{
	AuthTypeAnonymous:       {},
	AuthTypePresigned:       {},
	AuthTypePresignedV2:     {},
	AuthTypeSigned:          {},
	AuthTypeSignedV2:        {},
	AuthTypePostPolicy:      {},
	AuthTypeStreamingSigned: {},
}

// Validate if the authType is valid and supported.
func isSupportedS3AuthType(aType AuthType) bool {
	_, ok := supportedS3AuthTypes[aType]
	return ok
}
