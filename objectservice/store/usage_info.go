package store

import "time"

type DataUsageInfo struct {
	// LastUpdate is the timestamp of when the data usage info was last updated.
	// This does not indicate a full scan.
	LastUpdate time.Time `json:"lastUpdate"`

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
	BucketsUsage map[string]BucketInfo `json:"bucketsUsageInfo"`
	// Deprecated kept here for backward compatibility reasons.
	BucketSizes map[string]uint64 `json:"bucketsSizes"`
}

// BucketInfo represents bucket usage of a bucket, and its relevant
// access type for an account
type BucketInfo struct {
	Name    string `json:"name"`
	Size    uint64 `json:"size"`
	Objects uint64 `json:"objects"`
}
