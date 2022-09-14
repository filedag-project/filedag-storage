package response

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectservice/apierrors"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func writeEmptyResponse(w http.ResponseWriter, r *http.Request, statusCode int) {
	writeResponse(w, r, statusCode, []byte{}, mimeNone)
}
func WriteErrorResponseHeadersOnly(w http.ResponseWriter, r *http.Request, err apierrors.ErrorCode) {
	writeResponse(w, r, apierrors.GetAPIError(err).HTTPStatusCode, nil, mimeNone)
}

//WriteErrorResponse write ErrorResponse
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, errorCode apierrors.ErrorCode) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	object := vars["object"]
	if strings.HasPrefix(object, "/") {
		object = object[1:]
	}

	apiError := apierrors.GetAPIError(errorCode)
	errorResponse := getRESTErrorResponse(apiError, r.URL.Path, bucket, object)
	WriteXMLResponse(w, r, apiError.HTTPStatusCode, errorResponse)
}
func getRESTErrorResponse(err apierrors.APIError, resource string, bucket, object string) apierrors.RESTErrorResponse {
	return apierrors.RESTErrorResponse{
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
