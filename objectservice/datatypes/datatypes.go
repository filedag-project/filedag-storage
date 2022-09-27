package datatypes

// DeletedObject objects deleted
type DeletedObject struct {
	DeleteMarker          bool   `xml:"DeleteMarker,omitempty"`
	DeleteMarkerVersionID string `xml:"DeleteMarkerVersionId,omitempty"`
	ObjectName            string `xml:"Key,omitempty"`
	VersionID             string `xml:"VersionId,omitempty"`
}

// ObjectV object version key/versionId
type ObjectV struct {
	ObjectName string `xml:"Key"`
	VersionID  string `xml:"VersionId"`
}

// ObjectToDelete carries key name for the object to delete.
type ObjectToDelete struct {
	ObjectV
}

// DeleteObjectsRequest - xml carrying the object key names which needs to be deleted.
type DeleteObjectsRequest struct {
	// Element to enable quiet mode for the request
	Quiet bool
	// List of objects to be deleted
	Objects []ObjectToDelete `xml:"Object"`
}

// CompletePart - represents the part that was completed, this is sent by the client
// during CompleteMultipartUpload request.
type CompletePart struct {
	// Part number identifying the part. This is a positive integer between 1 and
	// 10,000
	PartNumber int

	// Entity tag returned when the part was uploaded.
	ETag string

	// Checksum values. Optional.
	ChecksumCRC32  string
	ChecksumCRC32C string
	ChecksumSHA1   string
	ChecksumSHA256 string
}

// CompletedParts - is a collection satisfying sort.Interface.
type CompletedParts []CompletePart

func (a CompletedParts) Len() int           { return len(a) }
func (a CompletedParts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CompletedParts) Less(i, j int) bool { return a[i].PartNumber < a[j].PartNumber }

// CompleteMultipartUpload - represents list of parts which are completed, this is sent by the
// client during CompleteMultipartUpload request.
type CompleteMultipartUpload struct {
	Parts []CompletePart `xml:"Part"`
}
