package s3err

import (
	"net/http"
)

// NotFoundHandler If none of the http routes match respond with MethodNotAllowed
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
}
