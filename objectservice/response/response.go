package response

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/datatypes"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var log = logging.Logger("resp")

type mimeType string

const (
	mimeNone mimeType = ""
	mimeJSON mimeType = "application/json"
	//mimeXML application/xml UTF-8
	mimeXML mimeType = " application/xml"
)

const (
	MaxObjectList  = 1000  // Limit number of objects in a listObjectsResponse/listObjectsVersionsResponse.
	MaxDeleteList  = 1000  // Limit number of objects deleted in a delete call.
	MaxUploadsList = 10000 // Limit number of uploads in a listUploadsResponse.
	MaxPartsList   = 10000 // Limit number of parts in a listPartsResponse.
)

// APIErrorResponse - error response format
type APIErrorResponse struct {
	XMLName   xml.Name `xml:"Error" json:"-"`
	Code      string
	Message   string
	Resource  string
	RequestID string `xml:"RequestId" json:"RequestId"`
	HostID    string `xml:"HostId" json:"HostId"`
}

//WriteSuccessResponseXML Write Success Response XML
func WriteSuccessResponseXML(w http.ResponseWriter, r *http.Request, response interface{}) {
	WriteXMLResponse(w, r, http.StatusOK, response)
}

//WriteXMLResponse Write XMLResponse
func WriteXMLResponse(w http.ResponseWriter, r *http.Request, statusCode int, response interface{}) {
	writeResponse(w, r, statusCode, encodeXMLResponse(response), mimeXML)
}
func writeResponse(w http.ResponseWriter, r *http.Request, statusCode int, response []byte, mType mimeType) {
	setCommonHeaders(w, r)
	if response != nil {
		w.Header().Set(consts.ContentLength, strconv.Itoa(len(response)))
	}
	if mType != mimeNone {
		w.Header().Set(consts.ContentType, string(mType))
	}
	w.WriteHeader(statusCode)
	if response != nil {
		log.Infof("status %d %s: %s", statusCode, mType, string(response))
		_, err := w.Write(response)
		if err != nil {
			log.Errorf("write err: %v", err)
		}
		w.(http.Flusher).Flush()
	}
}
func setCommonHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("x-amz-request-id", fmt.Sprintf("%d", time.Now().UnixNano()))
	w.Header().Set("Accept-Ranges", "bytes")
	if r.Header.Get("Origin") != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
}

// encodeXMLResponse Encodes the response headers into XML format.
func encodeXMLResponse(response interface{}) []byte {
	var bytesBuffer bytes.Buffer
	bytesBuffer.WriteString(xml.Header)
	e := xml.NewEncoder(&bytesBuffer)
	e.Encode(response)
	return bytesBuffer.Bytes()
}

//WriteSuccessResponseEmpty  Success Response Empty
func WriteSuccessResponseEmpty(w http.ResponseWriter, r *http.Request) {
	writeEmptyResponse(w, r, http.StatusOK)
}

// WriteErrorResponseJSON - writes error response in JSON format;
// useful for admin APIs.
func WriteErrorResponseJSON(w http.ResponseWriter, err apierrors.APIError, reqURL *url.URL, host string) {
	// Generate error response.
	errorResponse := getAPIErrorResponse(err, reqURL.Path, w.Header().Get(consts.AmzRequestID), host)
	encodedErrorResponse := encodeResponseJSON(errorResponse)
	writeResponseSimple(w, err.HTTPStatusCode, encodedErrorResponse, mimeJSON)
}

// getErrorResponse gets in standard error and resource value and
// provides a encodable populated response values
func getAPIErrorResponse(err apierrors.APIError, resource, requestID, hostID string) APIErrorResponse {
	return APIErrorResponse{
		Code:      err.Code,
		Message:   err.Description,
		Resource:  resource,
		RequestID: requestID,
		HostID:    hostID,
	}
}

// Encodes the response headers into JSON format.
func encodeResponseJSON(response interface{}) []byte {
	var bytesBuffer bytes.Buffer
	e := json.NewEncoder(&bytesBuffer)
	e.Encode(response)
	return bytesBuffer.Bytes()
}

// WriteSuccessResponseJSON writes success headers and response if any,
// with content-type set to `application/json`.
func WriteSuccessResponseJSON(w http.ResponseWriter, response []byte) {
	writeResponseSimple(w, http.StatusOK, response, mimeJSON)
}
func writeResponseSimple(w http.ResponseWriter, statusCode int, response []byte, mType mimeType) {
	if mType != mimeNone {
		w.Header().Set(consts.ContentType, string(mType))
	}
	w.Header().Set(consts.ContentLength, strconv.Itoa(len(response)))
	w.WriteHeader(statusCode)
	if response != nil {
		w.Write(response)
	}
}

// WriteSuccessNoContent writes success headers with http status 204
func WriteSuccessNoContent(w http.ResponseWriter) {
	writeResponseSimple(w, http.StatusNoContent, nil, mimeNone)
}

//ListAllMyBucketsResult  List All Buckets Result
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListAllMyBucketsResult"`
	Owner   *s3.Owner
	Buckets []*s3.Bucket `xml:"Buckets>Bucket"`
}

//WriteSuccessResponseHeadersOnly write SuccessResponseHeadersOnly
func WriteSuccessResponseHeadersOnly(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, r, http.StatusOK, nil, mimeNone)
}

type CopyObjectResponse struct {
	CopyObjectResult CopyObjectResult `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CopyObjectResult"`
}

type CopyObjectResult struct {
	LastModified time.Time `xml:"http://s3.amazonaws.com/doc/2006-03-01/ LastModified"`
	ETag         string    `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ETag"`
}

// LocationResponse - format for location response.
type LocationResponse struct {
	XMLName  xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ LocationConstraint" json:"-"`
	Location string   `xml:",chardata"`
}

// ListObjectsResponse - format for list objects response.
type ListObjectsResponse struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListBucketResult" json:"-"`

	Name   string
	Prefix string
	Marker string

	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects. Server lists objects in alphabetical
	// order Note: This element is returned only if you have delimiter request parameter
	// specified. If response does not include the NextMaker and it is truncated,
	// you can use the value of the last Key in the response as the marker in the
	// subsequent request to get the next set of object keys.
	NextMarker string `xml:"NextMarker,omitempty"`

	MaxKeys   int
	Delimiter string
	// A flag that indicates whether or not ListObjects returned all of the results
	// that satisfied the search criteria.
	IsTruncated bool

	Contents       []Object
	CommonPrefixes []CommonPrefix

	// Encoding type used to encode object keys in the response.
	EncodingType string `xml:"EncodingType,omitempty"`
}

// ListObjectsV2Response - format for list objects response.
type ListObjectsV2Response struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ ListBucketResult" json:"-"`

	Name       string
	Prefix     string
	StartAfter string `xml:"StartAfter,omitempty"`
	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects. Server lists objects in alphabetical
	// order Note: This element is returned only if you have delimiter request parameter
	// specified. If response does not include the NextMaker and it is truncated,
	// you can use the value of the last Key in the response as the marker in the
	// subsequent request to get the next set of object keys.
	ContinuationToken     string `xml:"ContinuationToken,omitempty"`
	NextContinuationToken string `xml:"NextContinuationToken,omitempty"`

	KeyCount  int
	MaxKeys   int
	Delimiter string
	// A flag that indicates whether or not ListObjects returned all of the results
	// that satisfied the search criteria.
	IsTruncated bool

	Contents       []Object
	CommonPrefixes []CommonPrefix

	// Encoding type used to encode object keys in the response.
	EncodingType string `xml:"EncodingType,omitempty"`
}

// Object container for object metadata
type Object struct {
	Key          string
	LastModified string // time string of format "2006-01-02T15:04:05.000Z"
	ETag         string
	Size         int64

	// Owner of the object.
	Owner s3.Owner

	// The class of storage used to store the object.
	StorageClass string

	// UserMetadata user-defined metadata
	UserMetadata StringMap `xml:"UserMetadata,omitempty"`
}

// StringMap is a map[string]string
type StringMap map[string]string

// MarshalXML - StringMap marshals into XML.
func (s StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}

	for key, value := range s {
		t := xml.StartElement{}
		t.Name = xml.Name{
			Space: "",
			Local: key,
		}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{Name: t.Name})
	}

	tokens = append(tokens, xml.EndElement{
		Name: start.Name,
	})

	for _, t := range tokens {
		if err := e.EncodeToken(t); err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	return e.Flush()
}

// CommonPrefix container for prefix response in ListObjectsResponse
type CommonPrefix struct {
	Prefix string
}

// DeleteError structure.
type DeleteError struct {
	Code      string
	Message   string
	Key       string
	VersionID string `xml:"VersionId"`
}

// DeleteObjectsResponse container for multiple object deletes.
type DeleteObjectsResponse struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ DeleteResult" json:"-"`

	// Collection of all deleted objects
	DeletedObjects []datatypes.DeletedObject `xml:"Deleted,omitempty"`

	// Collection of errors deleting certain objects.
	Errors []DeleteError `xml:"Error,omitempty"`
}
