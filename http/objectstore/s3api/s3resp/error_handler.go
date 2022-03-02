package s3resp

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func WriteEmptyResponse(w http.ResponseWriter, r *http.Request, statusCode int) {
	WriteResponse(w, r, statusCode, []byte{}, mimeNone)
}
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, errorCode ErrorCode) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	object := vars["object"]
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}

	apiError := GetAPIError(errorCode)
	errorResponse := getRESTErrorResponse(apiError, r.URL.Path, bucket, object)
	encodedErrorResponse := EncodeXMLResponse(errorResponse)
	WriteResponse(w, r, apiError.HTTPStatusCode, encodedErrorResponse, mimeXML)
}
func getRESTErrorResponse(err APIError, resource string, bucket, object string) RESTErrorResponse {
	return RESTErrorResponse{
		Code:       err.Code,
		BucketName: bucket,
		Key:        object,
		Message:    err.Description,
		Resource:   resource,
		RequestID:  fmt.Sprintf("%d", time.Now().UnixNano()),
	}
}

// NotFoundHandler If none of the http routes match respond with MethodNotAllowed
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
}
