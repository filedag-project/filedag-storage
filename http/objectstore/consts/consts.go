package consts

import (
	"github.com/dustin/go-humanize"
	"time"
)

//some const
const (
	// Iso8601TimeFormat RFC3339 a subset of the ISO8601 timestamp format. e.g 2014-04-29T18:30:38Z
	Iso8601TimeFormat = "2006-01-02T15:04:05.000Z" // Reply date format with nanosecond precision.

	StreamingContentSHA256 = "STREAMING-AWS4-HMAC-SHA256-PAYLOAD"

	Authorization = "Authorization"
	ETag          = "ETag"
	ContentType   = "Content-Type"

	// MaxLocationConstraintSize Limit of location constraint XML for unauthenticated PUT bucket operations.
	MaxLocationConstraintSize = 3 * humanize.MiByte
	EmptySHA256               = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	StsRequestBodyLimit       = 10 * (1 << 20) // 10 MiB
	DefaultRegion             = "us-east-1"
	Expires                   = "Expires"
	ContentMD5                = "Content-Md5"
	Date                      = "Date"
	SlashSeparator            = "/"

	MaxSkewTime = 15 * time.Minute // 15 minutes skew allowed.

	ContentLength = "Content-Length"

	// STS API version.
	StsAPIVersion   = "2011-06-15"
	StsVersion      = "Version"
	StsAction       = "Action"
	AssumeRole      = "AssumeRole"
	SignV4Algorithm = "AWS4-HMAC-SHA256"

	Location       = "Location"
	DefaultOwnerID = "02d6176db174dc93cb1b899f7c6078f08654445fe8cf1b6ce98d8855f66bdbf4"
	DisplayName    = "FileDagStorage"
)

// Standard S3 HTTP request constants
const (
	IfModifiedSince   = "If-Modified-Since"
	IfUnmodifiedSince = "If-Unmodified-Since"
	IfMatch           = "If-Match"
	IfNoneMatch       = "If-None-Match"

	// S3 storage class
	AmzStorageClass = "x-amz-storage-class"

	// S3 object version ID
	AmzVersionID    = "x-amz-version-id"
	AmzDeleteMarker = "x-amz-delete-marker"

	// S3 object tagging
	AmzObjectTagging = "X-Amz-Tagging"
	AmzTagCount      = "x-amz-tagging-count"
	AmzTagDirective  = "X-Amz-Tagging-Directive"

	// S3 transition restore
	AmzRestore            = "x-amz-restore"
	AmzRestoreExpiryDays  = "X-Amz-Restore-Expiry-Days"
	AmzRestoreRequestDate = "X-Amz-Restore-Request-Date"
	AmzRestoreOutputPath  = "x-amz-restore-output-path"

	// S3 extensions
	AmzCopySourceIfModifiedSince   = "x-amz-copy-source-if-modified-since"
	AmzCopySourceIfUnmodifiedSince = "x-amz-copy-source-if-unmodified-since"

	AmzCopySourceIfNoneMatch = "x-amz-copy-source-if-none-match"
	AmzCopySourceIfMatch     = "x-amz-copy-source-if-match"

	AmzCopySource                 = "X-Amz-Copy-Source"
	AmzCopySourceVersionID        = "X-Amz-Copy-Source-Version-Id"
	AmzCopySourceRange            = "X-Amz-Copy-Source-Range"
	AmzMetadataDirective          = "X-Amz-Metadata-Directive"
	AmzObjectLockMode             = "X-Amz-Object-Lock-Mode"
	AmzObjectLockRetainUntilDate  = "X-Amz-Object-Lock-Retain-Until-Date"
	AmzObjectLockLegalHold        = "X-Amz-Object-Lock-Legal-Hold"
	AmzObjectLockBypassGovernance = "X-Amz-Bypass-Governance-Retention"
	AmzBucketReplicationStatus    = "X-Amz-Replication-Status"
	AmzSnowballExtract            = "X-Amz-Meta-Snowball-Auto-Extract"

	// Multipart parts count
	AmzMpPartsCount = "x-amz-mp-parts-count"

	// Object date/time of expiration
	AmzExpiration = "x-amz-expiration"

	// Dummy putBucketACL
	AmzACL = "x-amz-acl"

	// Signature V4 related contants.
	AmzContentSha256        = "X-Amz-Content-Sha256"
	AmzDate                 = "X-Amz-Date"
	AmzAlgorithm            = "X-Amz-Algorithm"
	AmzExpires              = "X-Amz-Expires"
	AmzSignedHeaders        = "X-Amz-SignedHeaders"
	AmzSignature            = "X-Amz-Signature"
	AmzCredential           = "X-Amz-Credential"
	AmzSecurityToken        = "X-Amz-Security-Token"
	AmzDecodedContentLength = "X-Amz-Decoded-Content-Length"

	AmzMetaUnencryptedContentLength = "X-Amz-Meta-X-Amz-Unencrypted-Content-Length"
	AmzMetaUnencryptedContentMD5    = "X-Amz-Meta-X-Amz-Unencrypted-Content-Md5"

	// AWS server-side encryption headers for SSE-S3, SSE-KMS and SSE-C.
	AmzServerSideEncryption                      = "X-Amz-Server-Side-Encryption"
	AmzServerSideEncryptionKmsID                 = AmzServerSideEncryption + "-Aws-Kms-Key-Id"
	AmzServerSideEncryptionKmsContext            = AmzServerSideEncryption + "-Context"
	AmzServerSideEncryptionCustomerAlgorithm     = AmzServerSideEncryption + "-Customer-Algorithm"
	AmzServerSideEncryptionCustomerKey           = AmzServerSideEncryption + "-Customer-Key"
	AmzServerSideEncryptionCustomerKeyMD5        = AmzServerSideEncryption + "-Customer-Key-Md5"
	AmzServerSideEncryptionCopyCustomerAlgorithm = "X-Amz-Copy-Source-Server-Side-Encryption-Customer-Algorithm"
	AmzServerSideEncryptionCopyCustomerKey       = "X-Amz-Copy-Source-Server-Side-Encryption-Customer-Key"
	AmzServerSideEncryptionCopyCustomerKeyMD5    = "X-Amz-Copy-Source-Server-Side-Encryption-Customer-Key-Md5"

	AmzEncryptionAES = "AES256"
	AmzEncryptionKMS = "aws:kms"

	// Signature v2 related constants
	AmzSignatureV2 = "Signature"
	AmzAccessKeyID = "AWSAccessKeyId"

	// Response request id.
	AmzRequestID = "x-amz-request-id"
)

//object const
const (
	MaxObjectSize      = 5 * humanize.TiByte
	LastModified       = "Last-Modified"
	ContentEncoding    = "Content-Encoding"
	ContentLanguage    = "Content-Language"
	ContentDisposition = "Content-Disposition"
)

// limit
const (
	// The maximum allowed time difference between the incoming request
	// date and server date during signature verification.
	GlobalMaxSkewTime = 15 * time.Minute // 15 minutes skew allowed.
)
