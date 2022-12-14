package store

import (
	"encoding/xml"
	"github.com/dustin/go-humanize"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/utils/hash"
	"math"
	"time"
)

// BucketPolicy - Bucket level policy.
type BucketPolicy string

// Different types of Policies currently supported for buckets.
const (
	bucketPolicyNone      BucketPolicy = "private"
	bucketPolicyReadOnly               = "download"
	bucketPolicyReadWrite              = "public"
	bucketPolicyWriteOnly              = "upload"
)

// todo get real captivity
const defaultTotalCaptivity = math.MaxUint64
const (
	// bigFileThreshold is the point where we add readahead to put operations.
	bigFileThreshold = 64 * humanize.MiByte
	// equals unixfsChunkSize
	chunkSize int = 1 << 20

	objectKeyFormat        = "obj/%s/%s"
	allObjectPrefixFormat  = "obj/%s/%s"
	allObjectSeekKeyFormat = "obj/%s/%s"

	uploadKeyFormat        = "uploadObj/%s/%s/%s"
	allUploadPrefixFormat  = "uploadObj/%s/%s"
	allUploadSeekKeyFormat = "uploadObj/%s/%s/%s"

	deleteKeyFormat       = "delObj/%s"
	allDeletePrefixFormat = "delObj/"

	globalOperationTimeout = 5 * time.Minute
	deleteOperationTimeout = 1 * time.Minute

	maxCpuPercent        = 60
	maxUsedMemoryPercent = 80
)

// Tags is list of tags of XML request/response as per
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketTagging.html#API_GetBucketTagging_RequestBody
type Tags tagging
type tagging struct {
	XMLName xml.Name `xml:"Tagging"`
	TagSet  *TagSet  `xml:"TagSet"`
}

// TagSet represents list of unique tags.
type TagSet struct {
	TagMap   map[string]string
	IsObject bool
}

// BucketMetadata contains bucket metadata.
type BucketMetadata struct {
	Name    string
	Region  string
	Owner   string
	Created time.Time

	PolicyConfig  *policy.Policy
	TaggingConfig *Tags
}

// Read only object actions.
var readOnlyObjectActions = s3action.Action("s3:GetObject")

// Write only object actions.
var writeOnlyObjectActions = s3action.NewActionSet("s3:AbortMultipartUpload", "s3:DeleteObject", "s3:ListMultipartUploadParts", "s3:PutObject")

// Common bucket actions for both read and write policies.
var commonBucketActions = s3action.Action("s3:GetBucketLocation")

// Read only bucket actions.
var readOnlyBucketActions = s3action.Action("s3:ListBucket")

// Write only bucket actions.
var writeOnlyBucketActions = s3action.Action("s3:ListBucketMultipartUploads")

// ObjectInfo - represents object metadata.
//{
// 	Bucket = {string} "test"
// 	Name = {string} "default.exe"
// 	ModTime = {time.Time} 2022-03-18 10:54:43.308685163 +0800
// 	Size = {int64} 11604147
// 	IsDir = {bool} false
// 	ETag = {string} "a6b0b7ddb4630832ed47821af59aa125"
// 	VersionID = {string} ""
// 	IsLatest = {bool} false
// 	DeleteMarker = {bool} false
// 	ContentType = {string} "application/x-msdownload"
// 	ContentEncoding = {string} ""
// 	Expires = {time.Time} 0001-01-01 00:00:00 +0000
// 	Parts = {[]ObjectPartInfo} nil
// 	AccTime = {time.Time} 0001-01-01 00:00:00 +0000
// 	SuccessorModTime = {time.Time} 0001-01-01 00:00:00 +0000
//}
type ObjectInfo struct {
	// Name of the bucket.
	Bucket string

	// Name of the object.
	Name string

	// Date and time when the object was last modified.
	ModTime time.Time

	// Total object size.
	Size int64

	// IsDir indicates if the object is prefix.
	IsDir bool

	// Hex encoded unique entity tag of the object.
	ETag string

	// Version ID of this object.
	VersionID string

	// IsLatest indicates if this is the latest current version
	// latest can be true for delete marker or a version.
	IsLatest bool

	// DeleteMarker indicates if the versionId corresponds
	// to a delete marker on an object.
	DeleteMarker bool

	// A standard MIME type describing the format of the object.
	ContentType string

	// Specifies what content encodings have been applied to the object and thus
	// what decoding mechanisms must be applied to obtain the object referenced
	// by the Content-Type header field.
	ContentEncoding string

	// Date and time at which the object is no longer able to be cached
	Expires time.Time

	// Date and time when the object was last accessed.
	AccTime time.Time

	//  The mod time of the successor object version if any
	SuccessorModTime time.Time
}

// objectPartInfo Info of each part kept in the multipart metadata
// file after CompleteMultipartUpload() is called.
type objectPartInfo struct {
	ETag    string    `json:"etag,omitempty"`
	Number  int       `json:"number"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type MultipartInfo struct {
	Bucket    string
	Object    string
	UploadID  string
	Initiated time.Time
	MetaData  map[string]string
	// List of individual parts, maximum size of upto 10,000
	Parts []objectPartInfo
}

// putObjReader is a type that wraps sio.EncryptReader and
// underlying hash.Reader in a struct
type putObjReader struct {
	*hash.Reader              // actual data stream
	rawReader    *hash.Reader // original data stream
	sealMD5Fn    sealMD5CurrFn
}

// sealMD5CurrFn seals md5sum with object encryption key and returns sealed
// md5sum
type sealMD5CurrFn func([]byte) []byte

// versionPurgeStatusType represents status of a versioned delete or permanent delete w.r.t bucket replication
type versionPurgeStatusType string

const (
	// pending - versioned delete replication is pending.
	pending versionPurgeStatusType = "PENDING"

	// complete - versioned delete replication is now complete, erase version on disk.
	complete versionPurgeStatusType = "COMPLETE"

	// failed - versioned delete replication failed.
	failed versionPurgeStatusType = "FAILED"
)

// Empty returns true if purge status was not set.
func (v versionPurgeStatusType) Empty() bool {
	return string(v) == ""
}

// Pending  returns true if the version is pending purge.
func (v versionPurgeStatusType) Pending() bool {
	return v == pending || v == failed
}

// ListPartsInfo - represents list of all parts.
type ListPartsInfo struct {
	// Name of the bucket.
	Bucket string

	// Name of the object.
	Object string

	// Upload ID identifying the multipart upload whose parts are being listed.
	UploadID string

	// Part number after which listing begins.
	PartNumberMarker int

	// When a list is truncated, this element specifies the last part in the list,
	// as well as the value to use for the part-number-marker request parameter
	// in a subsequent request.
	NextPartNumberMarker int

	// Maximum number of parts that were allowed in the response.
	MaxParts int

	// Indicates whether the returned list of parts is truncated.
	IsTruncated bool

	// List of all parts.
	Parts []objectPartInfo

	// Any metadata set during InitMultipartUpload, including encryption headers.
	Metadata map[string]string

	// ChecksumAlgorithm if set
	ChecksumAlgorithm string
}

// ListMultipartsInfo - represnets bucket resources for incomplete multipart uploads.
type ListMultipartsInfo struct {
	// Together with upload-id-marker, this parameter specifies the multipart upload
	// after which listing should begin.
	KeyMarker string

	// Together with key-marker, specifies the multipart upload after which listing
	// should begin. If key-marker is not specified, the upload-id-marker parameter
	// is ignored.
	UploadIDMarker string

	// When a list is truncated, this element specifies the value that should be
	// used for the key-marker request parameter in a subsequent request.
	NextKeyMarker string

	// When a list is truncated, this element specifies the value that should be
	// used for the upload-id-marker request parameter in a subsequent request.
	NextUploadIDMarker string

	// Maximum number of multipart uploads that could have been included in the
	// response.
	MaxUploads int

	// Indicates whether the returned list of multipart uploads is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of multipart uploads exceeds the limit allowed or specified
	// by max uploads.
	IsTruncated bool

	// List of all pending uploads.
	Uploads []MultipartInfo

	// When a prefix is provided in the request, The result contains only keys
	// starting with the specified prefix.
	Prefix string

	// A character used to truncate the object prefixes.
	// NOTE: only supported delimiter is '/'.
	Delimiter string

	// CommonPrefixes contains all (if there are any) keys between Prefix and the
	// next occurrence of the string specified by delimiter.
	CommonPrefixes []string

	EncodingType string // Not supported yet.
}

// ListObjectsV2Info - container for list objects version 2.
type ListObjectsV2Info struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: This element is returned only if you have delimiter request parameter
	// specified.
	ContinuationToken     string
	NextContinuationToken string

	// List of objects info for this request.
	Objects []ObjectInfo

	// List of prefixes for this request.
	Prefixes []string
}

// ListObjectsInfo - container for list objects.
type ListObjectsInfo struct {
	// Indicates whether the returned list objects response is truncated. A
	// value of true indicates that the list was truncated. The list can be truncated
	// if the number of objects exceeds the limit allowed or specified
	// by max keys.
	IsTruncated bool

	// When response is truncated (the IsTruncated element value in the response is true),
	// you can use the key name in this field as marker in the subsequent
	// request to get next set of objects.
	//
	// NOTE: AWS S3 returns NextMarker only if you have delimiter request parameter specified,
	NextMarker string

	// List of objects info for this request.
	Objects []ObjectInfo

	// List of prefixes for this request.
	Prefixes []string
}

type DataUsageInfo struct {
	TotalCaptivity uint64 `json:"totalCaptivity"`
	// LastUpdate is the timestamp of when the data usage info was last updated.
	// This does not indicate a full scan.
	//LastUpdate time.Time `json:"lastUpdate"`

	// Objects total count across all buckets
	ObjectsTotalCount uint64 `json:"objectsCount"`

	// Objects total size across all buckets
	ObjectsTotalSize uint64 `json:"objectsTotalSize"`
	// Total number of buckets in this cluster
	BucketsCount uint64 `json:"bucketsCount"`

	// Buckets usage info provides following information across all buckets
	// - total size of the bucket
	// - total objects in a bucket
	// - object size histogram per bucket
	BucketsUsage []BucketInfo `json:"bucketsUsageInfo"`
	// Deprecated kept here for backward compatibility reasons.
	//BucketSizes map[string]uint64 `json:"bucketsSizes"`
}

// BucketInfo represents bucket usage of a bucket, and its relevant
// access type for an account
type BucketInfo struct {
	Name    string `json:"name"`
	Size    uint64 `json:"size"`
	Objects uint64 `json:"objects"`
}
