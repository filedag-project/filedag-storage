package response

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
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

func WriteSuccessResponseXML(w http.ResponseWriter, r *http.Request, response interface{}) {
	WriteXMLResponse(w, r, http.StatusOK, response)
}
func WriteXMLResponse(w http.ResponseWriter, r *http.Request, statusCode int, response interface{}) {
	WriteResponse(w, r, statusCode, EncodeXMLResponse(response), mimeXML)
}
func WriteResponse(w http.ResponseWriter, r *http.Request, statusCode int, response []byte, mType mimeType) {
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

// EncodeXMLResponse Encodes the response headers into XML format.
func EncodeXMLResponse(response interface{}) []byte {
	var bytesBuffer bytes.Buffer
	bytesBuffer.WriteString(xml.Header)
	e := xml.NewEncoder(&bytesBuffer)
	e.Encode(response)
	return bytesBuffer.Bytes()
}

//WriteSuccessResponseEmpty  Success Response Empty
func WriteSuccessResponseEmpty(w http.ResponseWriter, r *http.Request) {
	WriteEmptyResponse(w, r, http.StatusOK)
}

// WriteErrorResponseJSON - writes error response in JSON format;
// useful for admin APIs.
func WriteErrorResponseJSON(ctx context.Context, w http.ResponseWriter, err api_errors.APIError, reqURL *url.URL, host string) {
	// Generate error response.
	errorResponse := getAPIErrorResponse(ctx, err, reqURL.Path, w.Header().Get(consts.AmzRequestID), host)
	encodedErrorResponse := encodeResponseJSON(errorResponse)
	writeResponse(w, err.HTTPStatusCode, encodedErrorResponse, mimeJSON)
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
	writeResponse(w, http.StatusOK, response, mimeJSON)
}
func writeResponse(w http.ResponseWriter, statusCode int, response []byte, mType mimeType) {
	if mType != mimeNone {
		w.Header().Set(consts.ContentType, string(mType))
	}
	w.Header().Set(consts.ContentLength, strconv.Itoa(len(response)))
	w.WriteHeader(statusCode)
	if response != nil {
		w.Write(response)
	}
}