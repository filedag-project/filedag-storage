package store

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/hash"
	"io"
	"time"
)

// ObjectInfo - represents object metadata.
// Bucket = {string} "test"
// Name = {string} "default.exe"
// ModTime = {time.Time} 2022-03-18 10:54:43.308685163 +0800
// Size = {int64} 11604147
// IsDir = {bool} false
// ETag = {string} "a6b0b7ddb4630832ed47821af59aa125"
// VersionID = {string} ""
// IsLatest = {bool} false
// DeleteMarker = {bool} false
// RestoreExpires = {time.Time} 0001-01-01 00:00:00 +0000
// RestoreOngoing = {bool} false
// ContentType = {string} "application/x-msdownload"
// ContentEncoding = {string} ""
// Expires = {time.Time} 0001-01-01 00:00:00 +0000
// StorageClass = {string} "STANDARD"
// UserDefined = {map[string]string}
// UserTags = {string} ""
// Parts = {[]ObjectPartInfo} nil
// Writer = {io.WriteCloser} nil
// Reader = {*hash.Reader | 0x0} nil
// putObjReader = {*putObjReader | 0x0} nil
// AccTime = {time.Time} 0001-01-01 00:00:00 +0000
// Legacy = {bool} false
// VersionPurgeStatusInternal = {string} ""
// VersionPurgeStatus = {VersionPurgeStatusType} ""
// NumVersions = {int} 0
// SuccessorModTime = {time.Time} 0001-01-01 00:00:00 +0000
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

	// RestoreExpires indicates date a restored object expires
	RestoreExpires time.Time

	// RestoreOngoing indicates if a restore is in progress
	RestoreOngoing bool

	// A standard MIME type describing the format of the object.
	ContentType string

	// Specifies what content encodings have been applied to the object and thus
	// what decoding mechanisms must be applied to obtain the object referenced
	// by the Content-Type header field.
	ContentEncoding string

	// Date and time at which the object is no longer able to be cached
	Expires time.Time

	// Specify object storage class
	StorageClass string

	// User-Defined metadata
	UserDefined map[string]string

	// User-Defined object tags
	UserTags string

	// List of individual parts, maximum size of upto 10,000
	Parts []objectPartInfo `json:"-"`

	// Implements writer and reader used by CopyObject API
	Writer       io.WriteCloser `json:"-"`
	Reader       *hash.Reader   `json:"-"`
	PutObjReader *putObjReader  `json:"-"`

	// Date and time when the object was last accessed.
	AccTime time.Time

	Legacy bool // indicates object on disk is in legacy data format

	// internal representation of version purge status
	VersionPurgeStatusInternal string
	VersionPurgeStatus         versionPurgeStatusType

	// The total count of all versions of this object
	NumVersions int
	//  The mod time of the successor object version if any
	SuccessorModTime time.Time
}

// objectPartInfo Info of each part kept in the multipart metadata
// file after CompleteMultipartUpload() is called.
type objectPartInfo struct {
	ETag       string `json:"etag,omitempty"`
	Number     int    `json:"number"`
	Size       int64  `json:"size"`
	ActualSize int64  `json:"actualSize"`
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
