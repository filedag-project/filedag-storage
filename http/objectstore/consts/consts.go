package consts

import (
	"github.com/dustin/go-humanize"
	"time"
)

//some const
const (
	StreamingContentSHA256 = "STREAMING-AWS4-HMAC-SHA256-PAYLOAD"

	Authorization = "Authorization"
	ETag          = "ETag"
	// S3 object version ID
	AmzVersionID    = "x-amz-version-id"
	AmzDeleteMarker = "x-amz-delete-marker"

	AmzCredential = "X-Amz-Credential"
	ContentType   = "Content-Type"

	// AmzSignatureV2 Signature v2 related constants
	AmzSignatureV2 = "Signature"
	AmzAccessKeyID = "AWSAccessKeyId"
	// Signature V4 related contants.
	AmzSignature     = "X-Amz-Signature"
	AmzDate          = "X-Amz-Date"
	AmzExpires       = "X-Amz-Expires"
	AmzSignedHeaders = "X-Amz-SignedHeaders"
	AmzSecurityToken = "X-Amz-Security-Token"
	AmzAlgorithm     = "X-Amz-Algorithm"

	// AmzContentSha256 Signature V4 related contants.
	AmzContentSha256 = "X-Amz-Content-Sha256"

	// MaxLocationConstraintSize Limit of location constraint XML for unauthenticated PUT bucket operations.
	MaxLocationConstraintSize = 3 * humanize.MiByte
	EmptySHA256               = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	StsRequestBodyLimit       = 10 * (1 << 20) // 10 MiB
	DefaultRegion             = ""
	Expires                   = "Expires"
	ContentMD5                = "Content-Md5"
	Date                      = "Date"
	SlashSeparator            = "/"

	MaxSkewTime = 15 * time.Minute // 15 minutes skew allowed.

	ContentLength = "Content-Length"
	// Response request id.
	AmzRequestID    = "x-amz-request-id"
	SignV4Algorithm = "AWS4-HMAC-SHA256"

	// STS API version.
	StsAPIVersion = "2011-06-15"
	StsVersion    = "Version"
	StsAction     = "Action"
	AssumeRole    = "AssumeRole"

	// Dummy putBucketACL
	AmzACL = "x-amz-acl"
)

//object const
const (
	LastModified            = "Last-Modified"
	ContentEncoding         = "Content-Encoding"
	AmzTagCount             = "x-amz-tagging-count"
	AmzServerSideEncryption = "X-Amz-Server-Side-Encryption"
	AmzEncryptionAES        = "AES256"
	ContentLanguage         = "Content-Language"
	ContentDisposition      = "Content-Disposition"
)
