package response

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
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
		w.Header().Set("Content-Length", strconv.Itoa(len(response)))
	}
	if mType != mimeNone {
		w.Header().Set("Content-Type", string(mType))
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
func WriteErrorResponseJSON(ctx context.Context, w http.ResponseWriter, err api_errors.APIError, reqURL *url.URL, host string) {
	// Generate error response.
	errorResponse := getAPIErrorResponse(ctx, err, reqURL.Path, w.Header().Get(consts.AmzRequestID), host)
	encodedErrorResponse := encodeResponseJSON(errorResponse)
	writeResponseSimple(w, err.HTTPStatusCode, encodedErrorResponse, mimeJSON)
}

// getErrorResponse gets in standard error and resource value and
// provides a encodable populated response values
func getAPIErrorResponse(ctx context.Context, err api_errors.APIError, resource, requestID, hostID string) APIErrorResponse {
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
