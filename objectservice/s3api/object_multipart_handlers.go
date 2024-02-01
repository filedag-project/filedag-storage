package s3api

import (
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/datatypes"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/filedag-project/filedag-storage/objectservice/utils/etag"
	"github.com/filedag-project/filedag-storage/objectservice/utils/hash"
	"github.com/filedag-project/filedag-storage/objectservice/utils/s3utils"
	"io"
	"net/http"
	"net/textproto"
	"net/url"
	"path"
	"sort"
	"strconv"
	"strings"
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

	if err := s3utils.CheckNewMultipartArgs(ctx, bucket, object); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	_, _, s3err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.PutObjectAction, bucket, object)
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
	metadata[consts.AmzMetaFileSize] = textproto.MIMEHeader(r.Header).Get(consts.AmzMetaFileSize)
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

	if err := s3utils.CheckPutObjectPartArgs(ctx, bucket, object); err != nil {
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
	r.Header.Set("file-type", path.Ext(object))
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

	if err := s3utils.CheckCompleteMultipartArgs(ctx, bucket, object); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	_, _, s3err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.PutObjectAction, bucket, object)
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
	r.Header.Set("file-size", strconv.FormatInt(objInfo.Size, 10))
	r.Header.Set("file-type", path.Ext(object))
	response.WriteSuccessResponseXML(w, r, resp)
}

// AbortMultipartUploadHandler - Aborts multipart upload.
func (s3a *s3ApiServer) AbortMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object, err := getBucketAndObject(r)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(r.Context(), err))
		return
	}
	ctx := r.Context()

	if err := s3utils.CheckAbortMultipartArgs(ctx, bucket, object); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	_, _, s3err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.AbortMultipartUploadAction, bucket, object)
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}

	if !s3a.bmSys.HasBucket(ctx, bucket) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	// Get upload id.
	uploadID, _, _, _, s3Error := getObjectResources(r.Form)
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}

	log.Infow("AbortMultipartUploadHandler", "bucket", bucket, "object", object, "uploadID", uploadID)

	err = s3a.store.AbortMultipartUpload(ctx, bucket, object, uploadID)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}
	response.WriteSuccessNoContent(w)
}

// ListMultipartUploadsHandler - GET Bucket (List Multipart uploads)
// -------------------------
// This operation lists in-progress multipart uploads. An in-progress
// multipart upload is a multipart upload that has been initiated,
// using the Initiate Multipart Upload request, but has not yet been
// completed or aborted. This operation returns at most 1,000 multipart
// uploads in the response.
func (s3a *s3ApiServer) ListMultipartUploadsHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _, _ := getBucketAndObject(r)
	ctx := r.Context()

	log.Infof("ListMultipartUploadsHandler %s", bucket)
	_, _, s3err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.ListBucketMultipartUploadsAction, bucket, "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}

	if !s3a.bmSys.HasBucket(ctx, bucket) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	prefix, keyMarker, uploadIDMarker, delimiter, maxUploads, encodingType, errCode := getBucketMultipartResources(r.Form)
	if errCode != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, errCode)
		return
	}

	if maxUploads < 0 {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidMaxUploads)
		return
	}

	if keyMarker != "" {
		// Marker not common with prefix is not implemented.
		if !strings.HasPrefix(keyMarker, prefix) {
			response.WriteErrorResponse(w, r, apierrors.ErrNotImplemented)
			return
		}
	}

	if err := s3utils.CheckListMultipartArgs(ctx, bucket, prefix, keyMarker, uploadIDMarker, delimiter); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	listMultipartsInfo, err := s3a.store.ListMultipartUploads(ctx, bucket, prefix, keyMarker, uploadIDMarker, delimiter, maxUploads)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}
	// generate response
	resp := response.GenerateListMultipartUploadsResponse(bucket, listMultipartsInfo, encodingType)
	// write success response.
	response.WriteSuccessResponseXML(w, r, resp)
}

// ListObjectPartsHandler - Lists object parts in a multipart upload.
func (s3a *s3ApiServer) ListObjectPartsHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object, err := getBucketAndObject(r)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(r.Context(), err))
		return
	}
	ctx := r.Context()

	log.Infof("ListObjectPartsHandler %s %s", bucket, object)

	_, _, s3err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.ListMultipartUploadPartsAction, bucket, object)
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}

	if !s3a.bmSys.HasBucket(ctx, bucket) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	uploadID, partNumberMarker, maxParts, encodingType, s3Error := getObjectResources(r.Form)
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if partNumberMarker < 0 {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidPartNumberMarker)
		return
	}
	if maxParts < 0 {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidMaxParts)
		return
	}

	if err := s3utils.CheckListPartsArgs(ctx, bucket, object); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}

	listPartsInfo, err := s3a.store.ListObjectParts(ctx, bucket, object, uploadID, partNumberMarker, maxParts)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err))
		return
	}
	resp := response.GenerateListPartsResponse(listPartsInfo, encodingType)
	response.WriteSuccessResponseXML(w, r, resp)
}

// CopyObjectPartHandler - uploads a part by copying data from an existing object as data source.
func (s3a *s3ApiServer) CopyObjectPartHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteErrorResponse(w, r, apierrors.ErrNotImplemented)
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
		maxUploads = consts.MaxUploadsList
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
		maxParts = consts.MaxPartsList
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
