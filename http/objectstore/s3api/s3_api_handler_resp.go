package s3api

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type mimeType string

const (
	mimeNone mimeType = ""
	//mimeXML application/xml UTF-8
	mimeXML mimeType = " application/xml"
)

func writeSuccessResponseXML(w http.ResponseWriter, r *http.Request, response interface{}) {
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
		log.Printf("status %d %s: %s", statusCode, mType, string(response))
		_, err := w.Write(response)
		if err != nil {
			log.Printf("write err: %v", err)
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
