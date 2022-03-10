package s3api

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func (s3a *s3ApiServer) PutObjectHandler(w http.ResponseWriter, r *http.Request) {

	// http://docs.aws.amazon.com/AmazonS3/latest/dev/UploadingObjects.html

	bucket, object := GetBucketAndObject(r)
	log.Infof("PutObjectHandler %s %s", bucket, object)

	dataReader := r.Body
	defer dataReader.Close()
	cid := ""
	var err error
	if cid, err = store.PutFile(".", bucket+object, r.Body); err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	w.Write([]byte(cid))
	response.WriteSuccessResponseEmpty(w, r)
}
func GetBucketAndObject(r *http.Request) (bucket, object string) {
	vars := mux.Vars(r)
	bucket = vars["bucket"]
	object = vars["object"]
	if !strings.HasPrefix(object, "/") {
		object = "/" + object
	}

	return
}
