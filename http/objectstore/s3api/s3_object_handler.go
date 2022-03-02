package s3api

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api/s3resp"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	"github.com/google/martian/log"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func (s3a *S3ApiServer) PutObjectHandler(w http.ResponseWriter, r *http.Request) {

	// http://docs.aws.amazon.com/AmazonS3/latest/dev/UploadingObjects.html

	bucket, object := GetBucketAndObject(r)
	log.Infof("PutObjectHandler %s %s", bucket, object)

	dataReader := r.Body
	defer dataReader.Close()

	if err := store.PutFile("BucketsPath", bucket+object); err != nil {
		s3resp.WriteErrorResponse(w, r, s3resp.ErrInternalError)
		return
	}

	s3resp.WriteSuccessResponseEmpty(w, r)
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
