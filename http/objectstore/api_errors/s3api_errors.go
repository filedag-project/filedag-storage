package api_errors

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// APIError structure
type APIError struct {
	Code           string
	Description    string
	HTTPStatusCode int
}

// RESTErrorResponse - error response format
type RESTErrorResponse struct {
	XMLName    xml.Name `xml:"Error" json:"-"`
	Code       string   `xml:"Code" json:"Code"`
	Message    string   `xml:"Message" json:"Message"`
	Resource   string   `xml:"Resource" json:"Resource"`
	RequestID  string   `xml:"RequestId" json:"RequestId"`
	Key        string   `xml:"Key,omitempty" json:"Key,omitempty"`
	BucketName string   `xml:"BucketName,omitempty" json:"BucketName,omitempty"`

	// Underlying HTTP status code for the returned error
	StatusCode int `xml:"-" json:"-"`
}

// Error - Returns S3 error string.
func (e RESTErrorResponse) Error() string {
	if e.Message == "" {
		msg, ok := s3ErrorResponseMap[e.Code]
		if !ok {
			msg = fmt.Sprintf("Error response code %s.", e.Code)
		}
		return msg
	}
	return e.Message
}

// ErrorCode type of error status.
type ErrorCode int

// Error codes, see full list at http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
const (
	ErrNone ErrorCode = iota
	ErrAccessDenied
	ErrAuthorizationHeaderMalformed
	ErrMethodNotAllowed
	ErrBucketNotEmpty
	ErrBucketAlreadyExists
	ErrBucketAlreadyOwnedByYou
	ErrNoSuchBucket
	ErrNoSuchBucketPolicy
	ErrNoSuchCORSConfiguration
	ErrNoSuchLifecycleConfiguration
	ErrNoSuchKey
	ErrNoSuchUpload
	ErrInvalidBucketName
	ErrInvalidDigest
	ErrInvalidMaxKeys
	ErrInvalidMaxUploads
	ErrInvalidMaxParts
	ErrInvalidPartNumberMarker
	ErrInvalidPart
	ErrInvalidServiceSTS
	ErrInvalidServiceS3
	ErrInvalidRequestVersion
	ErrPolicyTooLarge
	ErrInternalError
	ErrInvalidCopyDest
	ErrInvalidCopySource
	ErrInvalidTag
	ErrMissingSecurityHeader
	ErrInvalidToken
	ErrInvalidAccessKeyID
	ErrAccessKeyDisabled
	ErrJsonMarshal
	ErrAuthHeaderEmpty
	ErrSetHeader
	ErrSignatureVersionNotSupported
	ErrMalformedPOSTRequest
	ErrPOSTFileRequired
	ErrPostPolicyConditionInvalidFormat
	ErrEntityTooSmall
	ErrEntityTooLarge
	ErrMissingFields
	ErrMissingCredTag
	ErrCredMalformed
	ErrMalformedXML
	ErrMalformedDate
	ErrMalformedPresignedDate
	ErrMalformedCredentialDate
	ErrMissingSignHeadersTag
	ErrMissingSignTag
	ErrUnsignedHeaders
	ErrInvalidQueryParams
	ErrInvalidQuerySignatureAlgo
	ErrExpiredPresignRequest
	ErrMalformedExpires
	ErrNegativeExpires
	ErrMaximumExpires
	ErrSignatureDoesNotMatch
	ErrStoreMkdirFail
	ErrStorePutFail
	ErrContentSHA256Mismatch
	ErrRequestNotReadyYet
	ErrMalformedPolicy
	ErrMissingDateHeader
	ErrInvalidRequest
	ErrAuthNotSetup
	ErrNotImplemented
	ErrPreconditionFailed

	ErrExistingObjectIsDirectory
	ErrExistingObjectIsFile

	ErrSetBucketPolicyFail

	ErrReader
)

// error code to APIError structure, these fields carry respective
// descriptions for all the error responses.
var errorCodeResponse = map[ErrorCode]APIError{
	ErrAccessDenied: {
		Code:           "AccessDenied",
		Description:    "Access Denied.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrAuthorizationHeaderMalformed: {
		Code:           "AuthorizationHeaderMalformed",
		Description:    "The authorization header is malformed; the region is wrong; expecting 'us-east-1'.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMethodNotAllowed: {
		Code:           "MethodNotAllowed",
		Description:    "The specified method is not allowed against this resource.",
		HTTPStatusCode: http.StatusMethodNotAllowed,
	},
	ErrBucketNotEmpty: {
		Code:           "BucketNotEmpty",
		Description:    "The bucket you tried to delete is not empty",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrBucketAlreadyExists: {
		Code:           "BucketAlreadyExists",
		Description:    "The requested bucket name is not available. The bucket name can not be an existing collection, and the bucket namespace is shared by all users of the system. Please select a different name and try again.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrBucketAlreadyOwnedByYou: {
		Code:           "BucketAlreadyOwnedByYou",
		Description:    "Your previous request to create the named bucket succeeded and you already own it.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrInvalidBucketName: {
		Code:           "InvalidBucketName",
		Description:    "The specified bucket is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidDigest: {
		Code:           "InvalidDigest",
		Description:    "The Content-Md5 you specified is not valid.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMaxUploads: {
		Code:           "InvalidArgument",
		Description:    "Argument max-uploads must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMaxKeys: {
		Code:           "InvalidArgument",
		Description:    "Argument maxKeys must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidMaxParts: {
		Code:           "InvalidArgument",
		Description:    "Argument max-parts must be an integer between 0 and 2147483647",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidPartNumberMarker: {
		Code:           "InvalidArgument",
		Description:    "Argument partNumberMarker must be an integer.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidServiceSTS: {
		Code:           "AuthorizationParametersError",
		Description:    "Error parsing the Credential parameter; incorrect service. This endpoint belongs to \"sts\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidServiceS3: {
		Code:           "AuthorizationParametersError",
		Description:    "Error parsing the Credential/X-Amz-Credential parameter; incorrect service. This endpoint belongs to \"s3\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRequestVersion: {
		Code:           "AuthorizatioErrInvalidRequestVersionnQueryParametersError",
		Description:    "Error parsing the X-Amz-Credential parameter; incorrect terminal. This endpoint uses \"aws4_request\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrPolicyTooLarge: {
		Code:           "PolicyTooLarge",
		Description:    "Policy exceeds the maximum allowed document size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNoSuchBucket: {
		Code:           "NoSuchBucket",
		Description:    "The specified bucket does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchBucketPolicy: {
		Code:           "NoSuchBucketPolicy",
		Description:    "The bucket policy does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchCORSConfiguration: {
		Code:           "NoSuchCORSConfiguration",
		Description:    "The CORS configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchLifecycleConfiguration: {
		Code:           "NoSuchLifecycleConfiguration",
		Description:    "The lifecycle configuration does not exist",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchKey: {
		Code:           "NoSuchKey",
		Description:    "The specified key does not exist.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrNoSuchUpload: {
		Code:           "NoSuchUpload",
		Description:    "The specified multipart upload does not exist. The upload ID may be invalid, or the upload may have been aborted or completed.",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrInternalError: {
		Code:           "InternalError",
		Description:    "We encountered an internal error, please try again.",
		HTTPStatusCode: http.StatusInternalServerError,
	},

	ErrInvalidPart: {
		Code:           "InvalidPart",
		Description:    "One or more of the specified parts could not be found.  The part may not have been uploaded, or the specified entity tag may not match the part's entity tag.",
		HTTPStatusCode: http.StatusBadRequest,
	},

	ErrInvalidCopyDest: {
		Code:           "InvalidRequest",
		Description:    "This copy request is illegal because it is trying to copy an object to itself without changing the object's metadata, storage class, website redirect location or encryption attributes.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidCopySource: {
		Code:           "InvalidArgument",
		Description:    "Copy Source must mention the source bucket and key: sourcebucket/sourcekey.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidTag: {
		Code:           "InvalidTag",
		Description:    "The Tag value you have provided is invalid",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSecurityHeader: {
		Code:           "MissingSecurityHeader",
		Description:    "Your request was missing a required header",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidToken: {
		Code:           "Invalid token",
		Description:    "The token value you have provided is invalid",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedXML: {
		Code:           "MalformedXML",
		Description:    "The XML you provided was not well-formed or did not validate against our published schema.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAuthHeaderEmpty: {
		Code:           "InvalidArgument",
		Description:    "Authorization header is invalid -- one and only one ' ' (space) required.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSetHeader: {
		Code:           "ErrSetHeader",
		Description:    "set response header err",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrSignatureVersionNotSupported: {
		Code:           "InvalidRequest",
		Description:    "The authorization mechanism you have provided is not supported. Please use AWS4-HMAC-SHA256.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedPOSTRequest: {
		Code:           "MalformedPOSTRequest",
		Description:    "The body of your POST request is not well-formed multipart/form-data.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrPOSTFileRequired: {
		Code:           "InvalidArgument",
		Description:    "POST requires exactly one file upload per request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrPostPolicyConditionInvalidFormat: {
		Code:           "PostPolicyInvalidKeyName",
		Description:    "Invalid according to Policy: Policy Condition failed",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrEntityTooSmall: {
		Code:           "EntityTooSmall",
		Description:    "Your proposed upload is smaller than the minimum allowed object size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrEntityTooLarge: {
		Code:           "EntityTooLarge",
		Description:    "Your proposed upload exceeds the maximum allowed object size.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingFields: {
		Code:           "MissingFields",
		Description:    "Missing fields in request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingCredTag: {
		Code:           "InvalidRequest",
		Description:    "Missing Credential field for this request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrCredMalformed: {
		Code:           "AuthorizationQueryParametersError",
		Description:    "Error parsing the X-Amz-Credential parameter; the Credential is mal-formed; expecting \"<YOUR-AKID>/YYYYMMDD/REGION/SERVICE/aws4_request\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedDate: {
		Code:           "MalformedDate",
		Description:    "Invalid date format header, expected to be in ISO8601, RFC1123 or RFC1123Z time format.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMalformedPresignedDate: {
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Date must be in the ISO8601 Long Format \"yyyyMMdd'T'HHmmss'Z'\"",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSignHeadersTag: {
		Code:           "InvalidArgument",
		Description:    "Signature header missing SignedHeaders field.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingSignTag: {
		Code:           "AccessDenied",
		Description:    "Signature header missing Signature field.",
		HTTPStatusCode: http.StatusBadRequest,
	},

	ErrUnsignedHeaders: {
		Code:           "AccessDenied",
		Description:    "There were headers present in the request which were not signed",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidQueryParams: {
		Code:           "AuthorizationQueryParametersError",
		Description:    "Query-string authentication version 4 requires the X-Amz-Algorithm, X-Amz-Credential, X-Amz-Signature, X-Amz-Date, X-Amz-SignedHeaders, and X-Amz-Expires parameters.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidQuerySignatureAlgo: {
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Algorithm only supports \"AWS4-HMAC-SHA256\".",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrExpiredPresignRequest: {
		Code:           "AccessDenied",
		Description:    "Request has expired",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrMalformedExpires: {
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Expires should be a number",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNegativeExpires: {
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Expires must be non-negative",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMaximumExpires: {
		Code:           "AuthorizationQueryParametersError",
		Description:    "X-Amz-Expires must be less than a week (in seconds); that is, the given X-Amz-Expires must be less than 604800 seconds",
		HTTPStatusCode: http.StatusBadRequest,
	},

	ErrInvalidAccessKeyID: {
		Code:           "InvalidAccessKeyId",
		Description:    "The access key ID you provided does not exist in our records.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrAccessKeyDisabled: {
		Code:           "InvalidAccessKeyId",
		Description:    "Your account is disabled; please contact your administrator.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrJsonMarshal: {
		Code:           "JsonMarshalErr",
		Description:    "json Marshal err.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrRequestNotReadyYet: {
		Code:           "AccessDenied",
		Description:    "Request is not valid yet",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrMalformedPolicy: {
		Code:           "MalformedPolicy",
		Description:    "Policy has invalid resource.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSignatureDoesNotMatch: {
		Code:           "SignatureDoesNotMatch",
		Description:    "The request signature we calculated does not match the signature you provided. Check your key and signing method.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrStoreMkdirFail: {
		Code:           "StoreMkdirFail",
		Description:    "Fail to make a bucket ",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrStorePutFail: {
		Code:           "StoreMkdirFail",
		Description:    "Fail to put a object ",
		HTTPStatusCode: http.StatusForbidden,
	},

	ErrContentSHA256Mismatch: {
		Code:           "XAmzContentSHA256Mismatch",
		Description:    "The provided 'x-amz-content-sha256' header does not match what was computed.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrMissingDateHeader: {
		Code:           "AccessDenied",
		Description:    "AWS authentication requires a valid Date or x-amz-date header",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInvalidRequest: {
		Code:           "InvalidRequest",
		Description:    "Invalid Request",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrAuthNotSetup: {
		Code:           "InvalidRequest",
		Description:    "Signed request requires setting up SeaweedFS S3 authentication",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrNotImplemented: {
		Code:           "NotImplemented",
		Description:    "A header you provided implies functionality that is not implemented",
		HTTPStatusCode: http.StatusNotImplemented,
	},
	ErrPreconditionFailed: {
		Code:           "PreconditionFailed",
		Description:    "At least one of the pre-conditions you specified did not hold",
		HTTPStatusCode: http.StatusPreconditionFailed,
	},
	ErrExistingObjectIsDirectory: {
		Code:           "ExistingObjectIsDirectory",
		Description:    "Existing Object is a directory.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrExistingObjectIsFile: {
		Code:           "ExistingObjectIsFile",
		Description:    "Existing Object is a file.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrSetBucketPolicyFail: {
		Code:           "SetBucketPolicyFail",
		Description:    "Set bucketPolicy fail.",
		HTTPStatusCode: http.StatusConflict,
	},
	ErrReader: {
		Code:           "ErrReader",
		Description:    "Can not New Reader",
		HTTPStatusCode: http.StatusConflict,
	},
}

// GetAPIError provides API Error for input API error code.
func GetAPIError(code ErrorCode) APIError {
	return errorCodeResponse[code]
}

// STSErrorCode type of error status.
type STSErrorCode int

// STSError structure
type STSError struct {
	Code           string
	Description    string
	HTTPStatusCode int
}

// Error codes, non exhaustive list - http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRoleWithSAML.html
const (
	ErrSTSNone STSErrorCode = iota
	ErrSTSAccessDenied
	ErrSTSMissingParameter
	ErrSTSInvalidParameterValue
	ErrSTSWebIdentityExpiredToken
	ErrSTSClientGrantsExpiredToken
	ErrSTSInvalidClientGrantsToken
	ErrSTSMalformedPolicyDocument
	ErrSTSInsecureConnection
	ErrSTSInvalidClientCertificate
	ErrSTSNotInitialized
	ErrSTSInternalError
)

type stsErrorCodeMap map[STSErrorCode]STSError

//ToSTSErr code to err
func (e stsErrorCodeMap) ToSTSErr(errCode STSErrorCode) STSError {
	apiErr, ok := e[errCode]
	if !ok {
		return e[ErrSTSInternalError]
	}
	return apiErr
}

// StsErrCodes error code to STSError structure, these fields carry respective
// descriptions for all the error responses.
var StsErrCodes = stsErrorCodeMap{
	ErrSTSAccessDenied: {
		Code:           "AccessDenied",
		Description:    "Generating temporary credentials not allowed for this request.",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrSTSMissingParameter: {
		Code:           "MissingParameter",
		Description:    "A required parameter for the specified action is not supplied.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSInvalidParameterValue: {
		Code:           "InvalidParameterValue",
		Description:    "An invalid or out-of-range value was supplied for the input parameter.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSWebIdentityExpiredToken: {
		Code:           "ExpiredToken",
		Description:    "The web identity token that was passed is expired or is not valid. Get a new identity token from the identity provider and then retry the request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSClientGrantsExpiredToken: {
		Code:           "ExpiredToken",
		Description:    "The client grants that was passed is expired or is not valid. Get a new client grants token from the identity provider and then retry the request.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSInvalidClientGrantsToken: {
		Code:           "InvalidClientGrantsToken",
		Description:    "The client grants token that was passed could not be validated.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSMalformedPolicyDocument: {
		Code:           "MalformedPolicyDocument",
		Description:    "The request was rejected because the policy document was malformed.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSInsecureConnection: {
		Code:           "InsecureConnection",
		Description:    "The request was made over a plain HTTP connection. A TLS connection is required.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSInvalidClientCertificate: {
		Code:           "InvalidClientCertificate",
		Description:    "The provided client certificate is invalid. Retry with a different certificate.",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrSTSNotInitialized: {
		Code:           "STSNotInitialized",
		Description:    "STS API not initialized, please try again.",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrSTSInternalError: {
		Code:           "InternalError",
		Description:    "We encountered an internal error generating credentials, please try again.",
		HTTPStatusCode: http.StatusInternalServerError,
	},
}
