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
