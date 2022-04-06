package user

import (
	"encoding/json"
	"time"
)

// AccountInfo represents the account usage info of an
// account across buckets.
type AccountInfo struct {
	AccountName string
	Policy      json.RawMessage // Use iam/policy.Parse to parse the result, to be done by the caller.
	Buckets     []BucketAccessInfo
}

// BucketAccessInfo represents bucket usage of a bucket, and its relevant
// access type for an account
type BucketAccessInfo struct {
	Name                 string            `json:"name"`
	Size                 uint64            `json:"size"`
	Objects              uint64            `json:"objects"`
	ObjectSizesHistogram map[string]uint64 `json:"objectHistogram"`
	Details              *BucketDetails    `json:"details"`
	PrefixUsage          map[string]uint64 `json:"prefixUsage"`
	Created              time.Time         `json:"created"`
	Access               AccountAccess     `json:"access"`
}

// BucketDetails provides information about features currently
// turned-on per bucket.
type BucketDetails struct {
	Versioning          bool `json:"versioning"`
	VersioningSuspended bool `json:"versioningSuspended"`
	Locking             bool `json:"locking"`
	Replication         bool `json:"replication"`
}

// AccountAccess contains information about
type AccountAccess struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
}
