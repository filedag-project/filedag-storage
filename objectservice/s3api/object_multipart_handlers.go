package s3api

import (
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/datatypes"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/filedag-project/filedag-storage/objectservice/utils/etag"
	"github.com/filedag-project/filedag-storage/objectservice/utils/hash"
	"github.com/filedag-project/filedag-storage/objectservice/utils/s3utils"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

// NewMultipartUploadHandler - New multipart upload.
func (s3a *s3ApiServer) NewMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object, err := getBucketAndObject(r)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(r.Context(), err))
		return
	}
	ctx := r.Context()

	log.Infof("NewMultipartUploadHandler %s %s", bucket, object)

	if err := s3utils.CheckPutObjectArgs(ctx, bucket, object); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	// Check if put is allowed
	s3err := s3a.authSys.IsPutActionAllowed(ctx, r, s3action.PutObjectAction, bucket, object)
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}

	if !s3a.bmSys.HasBucket(ctx, bucket) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	metadata, err := extractMetadata(ctx, r)
	if err != nil {
		log.Errorf("NewMultipartUploadHandler extractMetadata err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidRequest)
		return
	}

	info, err := s3a.store.NewMultipartUpload(ctx, bucket, object, metadata)
	if err != nil {
		log.Errorf("NewMultipartUploadHandler NewMultipartUpload err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}
	resp := response.GenerateInitiateMultipartUploadResponse(bucket, object, info.UploadID)

	response.WriteSuccessResponseXML(w, r, resp)
}

// PutObjectPartHandler - Put an object part in a multipart upload.
func (s3a *s3ApiServer) PutObjectPartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bucket, object, err := getBucketAndObject(r)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(r.Context(), err))
		return
	}
	// X-Amz-Copy-Source shouldn't be set for this call.
	if _, ok := r.Header[consts.AmzCopySource]; ok {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidCopySource)
		return
	}
	clientETag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidDigest)
		return
	}
	// if Content-Length is unknown/missing, deny the request
	size := r.ContentLength
	rAuthType := iam.GetRequestAuthType(r)
	if iam.IsAuthTypeStreamingSigned(rAuthType) {
		if sizeStr, ok := r.Header[consts.AmzDecodedContentLength]; ok {
			if sizeStr[0] == "" {
				response.WriteErrorResponse(w, r, apierrors.ErrMissingContentLength)
				return
			}
			size, err = strconv.ParseInt(sizeStr[0], 10, 64)
			if err != nil {
				log.Errorf("ParseInt err:%v", err)
				response.WriteErrorResponse(w, r, apierrors.ErrBadRequest)
				return
			}
		}
	}
	if size == -1 {
		response.WriteErrorResponse(w, r, apierrors.ErrMissingContentLength)
		return
	}
	// maximum Upload size for multipart objects in a single operation
	if size > consts.MaxPartSize {
		response.WriteErrorResponse(w, r, apierrors.ErrEntityTooLarge)
		return
	}

	uploadID := r.Form.Get(consts.UploadID)
	partIDString := r.Form.Get(consts.PartNumber)

	partID, err := strconv.Atoi(partIDString)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidPart)
		return
	}

	// check partID with maximum part ID for multipart objects
	if partID > consts.MaxPartID {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidMaxParts)
		return
	}

	log.Infow("PutObjectPartHandler", "bucket", bucket, "object", object, "partID", partID)

	if err := s3utils.CheckPutObjectArgs(ctx, bucket, object); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	// Check if put is allowed
	s3err := s3a.authSys.IsPutActionAllowed(ctx, r, s3action.PutObjectAction, bucket, object)
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}

	if !s3a.bmSys.HasBucket(ctx, bucket) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	var (
		md5hex              = clientETag.String()
		sha256hex           = ""
		reader    io.Reader = r.Body
	)

	switch rAuthType {
	case iam.AuthTypeStreamingSigned:
		// Initialize stream signature verifier.
		reader, s3err = iam.NewSignV4ChunkedReader(r, s3a.authSys)
		if s3err != apierrors.ErrNone {
			response.WriteErrorResponse(w, r, s3err)
			return
		}
	case iam.AuthTypeSignedV2, iam.AuthTypePresignedV2:
		s3err = s3a.authSys.IsReqAuthenticatedV2(r)
		if s3err != apierrors.ErrNone {
			response.WriteErrorResponse(w, r, s3err)
			return
		}

	case iam.AuthTypePresigned, iam.AuthTypeSigned:
		if s3err = s3a.authSys.ReqSignatureV4Verify(r, "", iam.ServiceS3); s3err != apierrors.ErrNone {
			response.WriteErrorResponse(w, r, s3err)
			return
		}
		if !iam.SkipContentSha256Cksum(r) {
			sha256hex = iam.GetContentSha256Cksum(r, iam.ServiceS3)
		}
	}

	mi, err := s3a.store.GetMultipartInfo(ctx, bucket, object, uploadID)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	hashReader, err := hash.NewReader(reader, size, md5hex, sha256hex, size)
	if err != nil {
		log.Errorf("PutObjectHandler NewReader err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	partInfo, err := s3a.store.PutObjectPart(ctx, bucket, object, uploadID, partID, hashReader, size, mi.MetaData)
	if err != nil {
		// Verify if the underlying error is signature mismatch.
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	etag := partInfo.ETag

	// We must not use the http.Header().Set method here because some (broken)
	// clients expect the ETag header key to be literally "ETag" - not "Etag" (case-sensitive).
	// Therefore, we have to set the ETag directly as map entry.
	w.Header()[consts.ETag] = []string{"\"" + etag + "\""}

	response.WriteSuccessResponseHeadersOnly(w, r)
}

// CompleteMultipartUploadHandler - Completes multipart upload.
func (s3a *s3ApiServer) CompleteMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object, err := getBucketAndObject(r)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(r.Context(), err))
		return
	}
	ctx := r.Context()

	log.Infof("CompleteMultipartUploadHandler %s %s", bucket, object)

	if err := s3utils.CheckPutObjectArgs(ctx, bucket, object); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	// Check if put is allowed
	s3err := s3a.authSys.IsPutActionAllowed(ctx, r, s3action.PutObjectAction, bucket, object)
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}

	if !s3a.bmSys.HasBucket(ctx, bucket) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	// Content-Length is required and should be non-zero
	if r.ContentLength <= 0 {
		response.WriteErrorResponse(w, r, apierrors.ErrMissingContentLength)
		return
	}

	// Get upload id.
	uploadID, _, _, _, s3Error := getObjectResources(r.Form)
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}

	complMultipartUpload := &datatypes.CompleteMultipartUpload{}
	if err = utils.XmlDecoder(r.Body, complMultipartUpload, r.ContentLength); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrMalformedXML)
		return
	}
	if len(complMultipartUpload.Parts) == 0 {
		response.WriteErrorResponse(w, r, apierrors.ErrMalformedXML)
		return
	}
	if !sort.IsSorted(datatypes.CompletedParts(complMultipartUpload.Parts)) {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidPartOrder)
		return
	}

	if _, err = s3a.store.GetMultipartInfo(ctx, bucket, object, uploadID); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	objInfo, err := s3a.store.CompleteMultiPartUpload(ctx, bucket, object, uploadID, complMultipartUpload.Parts)
	if err != nil {
		log.Errorf("CompleteMultipartUploadHandler CompleteMultiPartUpload err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}
	bucketMetas, err := s3a.bmSys.GetBucketMeta(ctx, bucket)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}
	// Generate complete multipart response.
	resp := response.GenerateCompleteMultpartUploadResponse(bucket, object, bucketMetas.Region, objInfo)
	setPutObjHeaders(w, objInfo, false)
	response.WriteSuccessResponseXML(w, r, resp)
}

// AbortMultipartUploadHandler - Aborts multipart upload.
func (s3a *s3ApiServer) AbortMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//bucket, object := xhttp.GetBucketAndObject(r)
	//
	//// Get upload id.
	//uploadID, _, _, _ := getObjectResources(r.URL.Query())
	//
	//response, errCode := s3a.abortMultipartUpload(&s3.AbortMultipartUploadInput{
	//	Bucket:   aws.String(bucket),
	//	Key:      objectKey(aws.String(object)),
	//	UploadId: aws.String(uploadID),
	//})
	//
	//if errCode != s3err.ErrNone {
	//	s3err.WriteErrorResponse(w, r, errCode)
	//	return
	//}
	//
	//writeSuccessResponseXML(w, r, response)

}

// ListMultipartUploadsHandler - Lists multipart uploads.
func (s3a *s3ApiServer) ListMultipartUploadsHandler(w http.ResponseWriter, r *http.Request) {
	//bucket, _ := xhttp.GetBucketAndObject(r)
	//
	//prefix, keyMarker, uploadIDMarker, delimiter, maxUploads, encodingType := getBucketMultipartResources(r.URL.Query())
	//if maxUploads < 0 {
	//	s3err.WriteErrorResponse(w, r, s3err.ErrInvalidMaxUploads)
	//	return
	//}
	//if keyMarker != "" {
	//	// Marker not common with prefix is not implemented.
	//	if !strings.HasPrefix(keyMarker, prefix) {
	//		s3err.WriteErrorResponse(w, r, s3err.ErrNotImplemented)
	//		return
	//	}
	//}
	//
	//response, errCode := s3a.listMultipartUploads(&s3.ListMultipartUploadsInput{
	//	Bucket:         aws.String(bucket),
	//	Delimiter:      aws.String(delimiter),
	//	EncodingType:   aws.String(encodingType),
	//	KeyMarker:      aws.String(keyMarker),
	//	MaxUploads:     aws.Int64(int64(maxUploads)),
	//	Prefix:         aws.String(prefix),
	//	UploadIdMarker: aws.String(uploadIDMarker),
	//})
	//
	//if errCode != s3err.ErrNone {
	//	s3err.WriteErrorResponse(w, r, errCode)
	//	return
	//}
	//
	//// TODO handle encodingType
	//
	//writeSuccessResponseXML(w, r, response)
}

// ListObjectPartsHandler - Lists object parts in a multipart upload.
func (s3a *s3ApiServer) ListObjectPartsHandler(w http.ResponseWriter, r *http.Request) {
	//bucket, object := xhttp.GetBucketAndObject(r)
	//
	//uploadID, partNumberMarker, maxParts, _ := getObjectResources(r.URL.Query())
	//if partNumberMarker < 0 {
	//	s3err.WriteErrorResponse(w, r, s3err.ErrInvalidPartNumberMarker)
	//	return
	//}
	//if maxParts < 0 {
	//	s3err.WriteErrorResponse(w, r, s3err.ErrInvalidMaxParts)
	//	return
	//}
	//
	//response, errCode := s3a.listObjectParts(&s3.ListPartsInput{
	//	Bucket:           aws.String(bucket),
	//	Key:              objectKey(aws.String(object)),
	//	MaxParts:         aws.Int64(int64(maxParts)),
	//	PartNumberMarker: aws.Int64(int64(partNumberMarker)),
	//	UploadId:         aws.String(uploadID),
	//})
	//
	//if errCode != s3err.ErrNone {
	//	s3err.WriteErrorResponse(w, r, errCode)
	//	return
	//}
	//
	//writeSuccessResponseXML(w, r, response)

}

// CopyObjectPartHandler - uploads a part by copying data from an existing object as data source.
func (s3a *s3ApiServer) CopyObjectPartHandler(w http.ResponseWriter, r *http.Request) {

}

// Parse bucket url queries for ?uploads
func getBucketMultipartResources(values url.Values) (prefix, keyMarker, uploadIDMarker, delimiter string, maxUploads int, encodingType string, errCode apierrors.ErrorCode) {
	errCode = apierrors.ErrNone

	if values.Get("max-uploads") != "" {
		var err error
		if maxUploads, err = strconv.Atoi(values.Get("max-uploads")); err != nil {
			errCode = apierrors.ErrInvalidMaxUploads
			return
		}
	} else {
		maxUploads = response.MaxUploadsList
	}

	prefix = trimLeadingSlash(values.Get("prefix"))
	keyMarker = trimLeadingSlash(values.Get("key-marker"))
	uploadIDMarker = values.Get("upload-id-marker")
	delimiter = values.Get("delimiter")
	encodingType = values.Get("encoding-type")
	return
}

// Parse object url queries
func getObjectResources(values url.Values) (uploadID string, partNumberMarker, maxParts int, encodingType string, errCode apierrors.ErrorCode) {
	var err error
	errCode = apierrors.ErrNone

	if values.Get("max-parts") != "" {
		if maxParts, err = strconv.Atoi(values.Get("max-parts")); err != nil {
			errCode = apierrors.ErrInvalidMaxParts
			return
		}
	} else {
		maxParts = response.MaxPartsList
	}

	if values.Get("part-number-marker") != "" {
		if partNumberMarker, err = strconv.Atoi(values.Get("part-number-marker")); err != nil {
			errCode = apierrors.ErrInvalidPartNumberMarker
			return
		}
	}

	uploadID = values.Get("uploadId")
	encodingType = values.Get("encoding-type")
	return
}
