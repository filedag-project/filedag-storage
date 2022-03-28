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
	DefaultRegion             = "US"
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
	AmzACL         = "x-amz-acl"
	Location       = "Location"
	DefaultOwnerID = "02d6176db174dc93cb1b899f7c6078f08654445fe8cf1b6ce98d8855f66bdbf4"
)

//object const
const (
	AmzCopySource           = "X-Amz-Copy-Source"
	AmzDecodedContentLength = "X-Amz-Decoded-Content-Length"
	MaxObjectSize           = 5 * humanize.TiByte
	LastModified            = "Last-Modified"
	ContentEncoding         = "Content-Encoding"
	AmzTagCount             = "x-amz-tagging-count"
	AmzServerSideEncryption = "X-Amz-Server-Side-Encryption"
	AmzEncryptionAES        = "AES256"
	ContentLanguage         = "Content-Language"
	ContentDisposition      = "Content-Disposition"
)
