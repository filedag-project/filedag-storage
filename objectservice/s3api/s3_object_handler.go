package s3api

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/datatypes"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/filedag-project/filedag-storage/objectservice/utils/etag"
	"github.com/filedag-project/filedag-storage/objectservice/utils/hash"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

//PutObjectHandler Put ObjectHandler
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutObject.html
func (s3a *s3ApiServer) PutObjectHandler(w http.ResponseWriter, r *http.Request) {

	// http://docs.aws.amazon.com/AmazonS3/latest/dev/UploadingObjects.html

	bucket, object := getBucketAndObject(r)
	// X-Amz-Copy-Source shouldn't be set for this call.
	if _, ok := r.Header[consts.AmzCopySource]; ok {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidCopySource)
		return
	}
	log.Infof("PutObjectHandler %s %s", bucket, object)
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
	if size == 0 {
		response.WriteErrorResponse(w, r, apierrors.ErrPutBucketInBucket)
		return
	}
	// maximum Upload size for objects in a single operation
	if size > consts.MaxObjectSize {
		response.WriteErrorResponse(w, r, apierrors.ErrEntityTooLarge)
		return
	}

	ctx := r.Context()
	// Check if put is allowed
	s3err := s3a.authSys.IsPutActionAllowed(ctx, r, s3action.PutObjectAction, bucket, object)
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}

	cred, _, s3err := s3a.authSys.GetCredential(r)
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}
	if !s3a.bmSys.HasBucket(bucket, cred.AccessKey) {
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

	hashReader, err := hash.NewReader(reader, size, md5hex, sha256hex, size)
	if err != nil {
		apiErr := apierrors.ErrInternalError
		switch err.(type) {
		case hash.SHA256Mismatch:
			apiErr = apierrors.ErrContentSHA256Mismatch
		case hash.BadDigest:
			apiErr = apierrors.ErrBadDigest
		}
		log.Errorf("PutObjectHandler NewReader err:%v", err)
		response.WriteErrorResponse(w, r, apiErr)
		return
	}
	metadata, err := extractMetadata(ctx, r)
	if err != nil {
		log.Errorf("PutObjectHandler extractMetadata err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidRequest)
		return
	}
	// TODO: if bucket is unique, there is no need to store user in Object
	objInfo, err := s3a.store.StoreObject(ctx, cred.AccessKey, bucket, object, hashReader, size, metadata)
	if err != nil {
		if apierrors.ContextCanceled(ctx) {
			if ctx.Err() == context.Canceled {
				response.WriteErrorResponse(w, r, apierrors.ErrClientDisconnected)
				return
			}
		}
		log.Errorf("PutObjectHandler StoreObject err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	setPutObjHeaders(w, objInfo, false)
	response.WriteSuccessResponseHeadersOnly(w, r)
}

// GetObjectHandler - GET Object
// ----------
// This implementation of the GET operation retrieves object. To use GET,
// you must have READ access to the object.
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetObject.html
func (s3a *s3ApiServer) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)

	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetObjectAction, bucket, object)
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.bmSys.HasBucket(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}
	userName := cred.AccessKey
	if cred.AccessKey == "" {
		meta, err := s3a.bmSys.GetBucketMeta(bucket, cred.AccessKey)
		if err != nil {
			response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
			return
		}
		userName = meta.Owner
	}
	objInfo, reader, err := s3a.store.GetObject(r.Context(), userName, bucket, object)
	if err != nil {
		log.Errorf("GetObjectHandler GetObject err:%v", err)
		if xerrors.Is(err, store.ErrObjectNotFound) {
			response.WriteErrorResponseHeadersOnly(w, r, apierrors.ErrNoSuchKey)
		} else {
			response.WriteErrorResponseHeadersOnly(w, r, apierrors.ErrInternalError)
		}
		return
	}
	w.Header().Set(consts.AmzServerSideEncryption, consts.AmzEncryptionAES)

	response.SetObjectHeaders(w, r, objInfo)
	w.Header().Set(consts.ContentLength, strconv.FormatInt(objInfo.Size, 10))
	response.SetHeadGetRespHeaders(w, r.Form)
	_, err = io.Copy(w, reader)
	if err != nil {
		log.Errorf("GetObjectHandler reader readAll err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
}

// HeadObjectHandler - HEAD Object
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadObject.html
// The HEAD operation retrieves metadata from an object without returning the object itself.
func (s3a *s3ApiServer) HeadObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)

	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetObjectAction, bucket, object)
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.bmSys.HasBucket(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}
	objInfo, ok := s3a.store.HasObject(r.Context(), cred.AccessKey, bucket, object)
	if !ok {
		response.WriteErrorResponseHeadersOnly(w, r, apierrors.ErrNoSuchKey)
		return
	}
	w.Header().Set(consts.AmzServerSideEncryption, consts.AmzEncryptionAES)

	// Set standard object headers.
	response.SetObjectHeaders(w, r, objInfo)
	// Set any additional requested response headers.
	response.SetHeadGetRespHeaders(w, r.Form)

	// Successful response.
	w.WriteHeader(http.StatusOK)

}

// DeleteObjectHandler - delete an object
// Delete objectAPIHandlers
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteObject.html
func (s3a *s3ApiServer) DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)

	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetObjectAction, bucket, object)
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.bmSys.HasBucket(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}
	objInfo, ok := s3a.store.HasObject(r.Context(), cred.AccessKey, bucket, object)
	if !ok {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchKey)
		return
	}
	err := s3a.store.DeleteObject(r.Context(), cred.AccessKey, bucket, object)
	if err != nil {
		log.Errorf("DeleteObjectHandler DeleteObject  err:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	setPutObjHeaders(w, objInfo, true)
	response.WriteSuccessNoContent(w)
}

// DeleteMultipleObjectsHandler - Delete multiple objects
func (s3a *s3ApiServer) DeleteMultipleObjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	bucket, _ := getBucketAndObject(r)

	// Content-Md5 is requied should be set
	// http://docs.aws.amazon.com/AmazonS3/latest/API/multiobjectdeleteapi.html
	if _, ok := r.Header[consts.ContentMD5]; !ok {
		response.WriteErrorResponse(w, r, apierrors.ErrMissingContentMD5)
		return
	}

	// Content-Length is required and should be non-zero
	// http://docs.aws.amazon.com/AmazonS3/latest/API/multiobjectdeleteapi.html
	if r.ContentLength <= 0 {
		response.WriteErrorResponse(w, r, apierrors.ErrMissingContentLength)
		return
	}

	// The max. XML contains 100000 object names (each at most 1024 bytes long) + XML overhead
	const maxBodySize = 2 * 100000 * 1024

	// Unmarshal list of keys to be deleted.
	deleteObjectsReq := &datatypes.DeleteObjectsRequest{}
	if err := utils.XmlDecoder(r.Body, deleteObjectsReq, maxBodySize); err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrMalformedXML)
		return
	}

	objects := make([]datatypes.ObjectV, len(deleteObjectsReq.Objects))
	// Convert object name delete objects if it has `/` in the beginning.
	for i := range deleteObjectsReq.Objects {
		deleteObjectsReq.Objects[i].ObjectName = trimLeadingSlash(deleteObjectsReq.Objects[i].ObjectName)
		objects[i] = deleteObjectsReq.Objects[i].ObjectV
	}

	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.DeleteObjectAction, bucket, "")
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}

	// Before proceeding validate if bucket exists.
	if !s3a.bmSys.HasBucket(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	// Return Malformed XML as S3 spec if the number of objects is empty
	if len(deleteObjectsReq.Objects) == 0 || len(deleteObjectsReq.Objects) > response.MaxDeleteList {
		response.WriteErrorResponse(w, r, apierrors.ErrMalformedXML)
		return
	}

	objectsToDelete := map[datatypes.ObjectToDelete]int{}

	type deleteResult struct {
		delInfo datatypes.DeletedObject
		errInfo response.DeleteError
	}

	deleteResults := make([]deleteResult, len(deleteObjectsReq.Objects))

	for index, object := range deleteObjectsReq.Objects {
		_, _, s3Error = s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.DeleteObjectAction, bucket, object.ObjectName)
		if s3Error != apierrors.ErrNone {
			if s3Error == apierrors.ErrSignatureDoesNotMatch || s3Error == apierrors.ErrInvalidAccessKeyID {
				response.WriteErrorResponse(w, r, s3Error)
				return
			}
			apiErr := apierrors.GetAPIError(s3Error)
			deleteResults[index].errInfo = response.DeleteError{
				Code:      apiErr.Code,
				Message:   apiErr.Description,
				Key:       object.ObjectName,
				VersionID: object.VersionID,
			}
			continue
		}

		// Avoid duplicate objects, we use map to filter them out.
		if _, ok := objectsToDelete[object]; !ok {
			objectsToDelete[object] = index
		}
	}

	toNames := func(input map[datatypes.ObjectToDelete]int) (output []datatypes.ObjectToDelete) {
		output = make([]datatypes.ObjectToDelete, len(input))
		idx := 0
		for obj := range input {
			output[idx] = obj
			idx++
		}
		return
	}

	// Disable timeouts and cancellation
	ctx = utils.BgContext(ctx)

	deleteList := toNames(objectsToDelete)
	dObjects := make([]datatypes.DeletedObject, len(deleteList))
	errs := make([]error, len(deleteList))
	for i, obj := range deleteList {
		errs[i] = s3a.store.DeleteObject(ctx, cred.AccessKey, bucket, obj.ObjectName)
		if errs[i] == nil || xerrors.Is(errs[i], store.ErrObjectNotFound) {
			dObjects[i] = datatypes.DeletedObject{
				ObjectName: obj.ObjectName,
			}
			errs[i] = nil
		}
	}

	for i := range errs {
		objToDel := datatypes.ObjectToDelete{
			ObjectV: datatypes.ObjectV{
				ObjectName: dObjects[i].ObjectName,
				VersionID:  dObjects[i].VersionID,
			},
		}
		dindex := objectsToDelete[objToDel]
		if errs[i] == nil {
			deleteResults[dindex].delInfo = dObjects[i]
			continue
		}
		apiErr := apierrors.GetAPIError(apierrors.ErrInternalError)
		deleteResults[dindex].errInfo = response.DeleteError{
			Code:      apiErr.Code,
			Message:   apiErr.Description,
			Key:       deleteList[i].ObjectName,
			VersionID: deleteList[i].VersionID,
		}
	}

	// Generate response
	deleteErrors := make([]response.DeleteError, 0, len(deleteObjectsReq.Objects))
	deletedObjects := make([]datatypes.DeletedObject, 0, len(deleteObjectsReq.Objects))
	for _, deleteResult := range deleteResults {
		if deleteResult.errInfo.Code != "" {
			deleteErrors = append(deleteErrors, deleteResult.errInfo)
		} else {
			deletedObjects = append(deletedObjects, deleteResult.delInfo)
		}
	}

	resp := generateMultiDeleteResponse(deleteObjectsReq.Quiet, dObjects, deleteErrors)

	// Write success response.
	response.WriteSuccessResponseXML(w, r, resp)
}

// CopyObjectHandler - Copy Object
// ----------
// This implementation of the PUT operation adds an object to a bucket
// while reading the object from another source.
// Notice: The S3 client can send secret keys in headers for encryption related jobs,
// the handler should ensure to remove these keys before sending them to the object layer.
// Currently these keys are:
//   - X-Amz-Server-Side-Encryption-Customer-Key
//   - X-Amz-Copy-Source-Server-Side-Encryption-Customer-Key
func (s3a *s3ApiServer) CopyObjectHandler(w http.ResponseWriter, r *http.Request) {
	dstBucket, dstObject := getBucketAndObject(r)

	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetObjectAction, dstBucket, dstObject)
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.bmSys.HasBucket(dstBucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	// Copy source path.
	cpSrcPath, err := url.QueryUnescape(r.Header.Get("X-Amz-Copy-Source"))
	if err != nil {
		// Save unescaped string as is.
		cpSrcPath = r.Header.Get(consts.AmzCopySource)
	}

	srcBucket, srcObject := pathToBucketAndObject(cpSrcPath)
	if !s3a.bmSys.HasBucket(srcBucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}
	log.Debugf("CopyObjectHandler %s %s => %s %s", srcBucket, srcObject, dstBucket, dstObject)
	srcObjInfo, srcReader, err := s3a.store.GetObject(r.Context(), cred.AccessKey, srcBucket, srcObject)
	if err != nil {
		log.Errorf("CopyObjectHandler StoreObject err:%v", err)
		response.WriteErrorResponseHeadersOnly(w, r, apierrors.ErrNoSuchKey)
		return
	}
	metadata := make(map[string]string)
	metadata[strings.ToLower(consts.ContentType)] = srcObjInfo.ContentType
	metadata[strings.ToLower(consts.ContentEncoding)] = srcObjInfo.ContentEncoding
	if (srcBucket == dstBucket && srcObject == dstObject || cpSrcPath == "") && isReplace(r) {
		object, err := s3a.store.StoreObject(r.Context(), cred.AccessKey, dstBucket, dstObject, srcReader, srcObjInfo.Size, metadata)
		if err != nil {
			return
		}
		response.WriteSuccessResponseXML(w, r, response.CopyObjectResult{
			ETag:         fmt.Sprintf("%x", object.ETag),
			LastModified: time.Now().UTC(),
		})
		return
	}

	// If source object is empty or bucket is empty, reply back invalid copy source.
	if srcObject == "" || srcBucket == "" {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidCopySource)
		return
	}

	if srcBucket == dstBucket && srcObject == dstObject {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidCopyDest)
		return
	}
	obj, err := s3a.store.StoreObject(r.Context(), cred.AccessKey, dstBucket, dstObject, srcReader, srcObjInfo.Size, metadata)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}

	setEtag(w, obj.ETag)

	resp2 := response.CopyObjectResult{
		ETag:         obj.ETag,
		LastModified: time.Now().UTC(),
	}

	response.WriteSuccessResponseXML(w, r, resp2)
}

func (s3a *s3ApiServer) ListObjectsV1Handler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)

	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.ListBucketAction, bucket, "")
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.bmSys.HasBucket(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}
	// Extract all the litsObjectsV1 query params to their native values.
	prefix, marker, delimiter, maxKeys, encodingType, s3Error := getListObjectsV1Args(r.Form)
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	objs, err := s3a.store.ListObjects(r.Context(), cred.AccessKey, bucket, prefix, marker, delimiter, maxKeys)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	resp := generateListObjectsV1Response(bucket, prefix, marker, delimiter, encodingType, maxKeys, objs)
	// Write success response.
	response.WriteSuccessResponseXML(w, r, resp)
}
func (s3a *s3ApiServer) ListObjectsV2Handler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)

	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cerd, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.ListBucketAction, bucket, object)
	if s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.bmSys.HasBucket(bucket, cerd.AccessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	urlValues := r.Form
	// Extract all the listObjectsV2 query params to their native values.
	prefix, token, startAfter, delimiter, fetchOwner, maxKeys, encodingType, errCode := getListObjectsV2Args(urlValues)
	if errCode != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, errCode)
		return
	}
	// Validate the query params before beginning to serve the request.
	// fetch-owner is not validated since it is a boolean
	if s3Error := validateListObjectsArgs(token, delimiter, encodingType, maxKeys); s3Error != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}

	// Initiate a list objects operation based on the input params.
	// On success would return back ListObjectsInfo object to be
	// marshaled into S3 compatible XML header.
	listObjectsV2Info, err := s3a.store.ListObjectsV2(r.Context(), cerd.AccessKey, bucket, prefix, token, delimiter,
		maxKeys, fetchOwner, startAfter)

	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}

	resp := GenerateListObjectsV2Response(bucket, prefix, token, listObjectsV2Info.NextContinuationToken, startAfter,
		delimiter, encodingType, listObjectsV2Info.IsTruncated,
		maxKeys, listObjectsV2Info.Objects, listObjectsV2Info.Prefixes)

	// Write success response.
	response.WriteSuccessResponseXML(w, r, resp)
}

func getBucketAndObject(r *http.Request) (bucket, object string) {
	vars := mux.Vars(r)
	bucket = vars["bucket"]
	object = vars["object"]
	if !strings.HasPrefix(object, "/") {
		object = "/" + object
	}
	return
}

// setPutObjHeaders sets all the necessary headers returned back
// upon a success Put/Copy/CompleteMultipart/Delete requests
// to activate delete only headers set delete as true
func setPutObjHeaders(w http.ResponseWriter, objInfo store.ObjectInfo, delete bool) {
	// We must not use the http.Header().Set method here because some (broken)
	// clients expect the ETag header key to be literally "ETag" - not "Etag" (case-sensitive).
	// Therefore, we have to set the ETag directly as map entry.
	if objInfo.ETag != "" && !delete {
		w.Header()[consts.ETag] = []string{`"` + objInfo.ETag + `"`}
	}

	// Set the relevant version ID as part of the response header.
	if objInfo.VersionID != "" {
		w.Header()[consts.AmzVersionID] = []string{objInfo.VersionID}
		// If version is a deleted marker, set this header as well
		if objInfo.DeleteMarker && delete { // only returned during delete object
			w.Header()[consts.AmzDeleteMarker] = []string{strconv.FormatBool(objInfo.DeleteMarker)}
		}
	}

	if objInfo.Bucket != "" && objInfo.Name != "" {

	}
}
func pathToBucketAndObject(path string) (bucket, object string) {
	path = strings.TrimPrefix(path, "/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) == 2 {
		return parts[0], "/" + parts[1]
	}
	return parts[0], "/"
}
func setEtag(w http.ResponseWriter, etag string) {
	if etag != "" {
		if strings.HasPrefix(etag, "\"") {
			w.Header().Set("ETag", etag)
		} else {
			w.Header().Set("ETag", "\""+etag+"\"")
		}
	}
}
func isReplace(r *http.Request) bool {
	return r.Header.Get("X-Amz-Metadata-Directive") == "REPLACE"
}

// Parse bucket url queries
func getListObjectsV1Args(values url.Values) (prefix, marker, delimiter string, maxkeys int, encodingType string, errCode apierrors.ErrorCode) {
	errCode = apierrors.ErrNone

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.Atoi(values.Get("max-keys")); err != nil {
			errCode = apierrors.ErrInvalidMaxKeys
			return
		}
	} else {
		maxkeys = response.MaxObjectList
	}

	prefix = trimLeadingSlash(values.Get("prefix"))
	marker = trimLeadingSlash(values.Get("marker"))
	delimiter = values.Get("delimiter")
	encodingType = values.Get("encoding-type")
	return
}

// Parse bucket url queries for ListObjects V2.
func getListObjectsV2Args(values url.Values) (prefix, token, startAfter, delimiter string, fetchOwner bool, maxkeys int, encodingType string, errCode apierrors.ErrorCode) {
	errCode = apierrors.ErrNone

	// The continuation-token cannot be empty.
	if val, ok := values["continuation-token"]; ok {
		if len(val[0]) == 0 {
			errCode = apierrors.ErrInvalidToken
			return
		}
	}

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.Atoi(values.Get("max-keys")); err != nil {
			errCode = apierrors.ErrInvalidMaxKeys
			return
		}
		// Over flowing count - reset to maxObjectList.
		if maxkeys > response.MaxObjectList {
			maxkeys = response.MaxObjectList
		}
	} else {
		maxkeys = response.MaxObjectList
	}

	prefix = trimLeadingSlash(values.Get("prefix"))
	startAfter = trimLeadingSlash(values.Get("start-after"))
	delimiter = values.Get("delimiter")
	fetchOwner = values.Get("fetch-owner") == "true"
	encodingType = values.Get("encoding-type")

	if token = values.Get("continuation-token"); token != "" {
		decodedToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			errCode = apierrors.ErrIncorrectContinuationToken
			return
		}
		token = string(decodedToken)
	}
	return
}
func trimLeadingSlash(ep string) string {
	if len(ep) > 0 && ep[0] == '/' {
		// Path ends with '/' preserve it
		if ep[len(ep)-1] == '/' && len(ep) > 1 {
			ep = path.Clean(ep)
			ep += "/"
		} else {
			ep = path.Clean(ep)
		}
		ep = ep[1:]
	}
	return ep
}

// Validate all the ListObjects query arguments, returns an APIErrorCode
// if one of the args do not meet the required conditions.
// - delimiter if set should be equal to '/', otherwise the request is rejected.
// - marker if set should have a common prefix with 'prefix' param, otherwise
//   the request is rejected.
func validateListObjectsArgs(marker, delimiter, encodingType string, maxKeys int) apierrors.ErrorCode {
	// Max keys cannot be negative.
	if maxKeys < 0 {
		return apierrors.ErrInvalidMaxKeys
	}

	if encodingType != "" {
		// AWS S3 spec only supports 'url' encoding type
		if !strings.EqualFold(encodingType, "url") {
			return apierrors.ErrInvalidEncodingMethod
		}
	}

	return apierrors.ErrNone
}

// GenerateListObjectsV2Response Generates an ListObjectsV2 response for the said bucket with other enumerated options.
func GenerateListObjectsV2Response(bucket, prefix, token, nextToken, startAfter, delimiter, encodingType string, isTruncated bool, maxKeys int, objects []store.ObjectInfo, prefixes []string) response.ListObjectsV2Response {
	contents := make([]response.Object, 0, len(objects))
	id := consts.DefaultOwnerID
	name := consts.DisplayName
	owner := s3.Owner{
		ID:          &id,
		DisplayName: &name,
	}
	data := response.ListObjectsV2Response{}

	for _, object := range objects {
		content := response.Object{}
		if object.Name == "" {
			continue
		}
		content.Key = utils.S3EncodeName(object.Name, encodingType)
		content.LastModified = object.ModTime.UTC().Format(consts.Iso8601TimeFormat)
		if object.ETag != "" {
			content.ETag = "\"" + object.ETag + "\""
		}
		content.Size = object.Size
		content.Owner = owner
		contents = append(contents, content)
	}
	data.Name = bucket
	data.Contents = contents

	data.EncodingType = encodingType
	data.StartAfter = utils.S3EncodeName(startAfter, encodingType)
	data.Delimiter = utils.S3EncodeName(delimiter, encodingType)
	data.Prefix = utils.S3EncodeName(prefix, encodingType)
	data.MaxKeys = maxKeys
	data.ContinuationToken = base64.StdEncoding.EncodeToString([]byte(token))
	data.NextContinuationToken = base64.StdEncoding.EncodeToString([]byte(nextToken))
	data.IsTruncated = isTruncated

	commonPrefixes := make([]response.CommonPrefix, 0, len(prefixes))
	for _, prefix := range prefixes {
		prefixItem := response.CommonPrefix{}
		prefixItem.Prefix = utils.S3EncodeName(prefix, encodingType)
		commonPrefixes = append(commonPrefixes, prefixItem)
	}
	data.CommonPrefixes = commonPrefixes
	data.KeyCount = len(data.Contents) + len(data.CommonPrefixes)
	return data
}

// generates an ListObjectsV1 response for the said bucket with other enumerated options.
func generateListObjectsV1Response(bucket, prefix, marker, delimiter, encodingType string, maxKeys int, resp store.ListObjectsInfo) response.ListObjectsResponse {
	contents := make([]response.Object, 0, len(resp.Objects))
	id := consts.DefaultOwnerID
	name := consts.DisplayName
	owner := s3.Owner{
		ID:          &id,
		DisplayName: &name,
	}
	data := response.ListObjectsResponse{}

	for _, object := range resp.Objects {
		content := response.Object{}
		if object.Name == "" {
			continue
		}
		content.Key = utils.S3EncodeName(object.Name, encodingType)
		content.LastModified = object.ModTime.UTC().Format(consts.Iso8601TimeFormat)
		if object.ETag != "" {
			content.ETag = "\"" + object.ETag + "\""
		}
		content.Size = object.Size
		content.StorageClass = ""
		content.Owner = owner
		contents = append(contents, content)
	}
	data.Name = bucket
	data.Contents = contents

	data.EncodingType = encodingType
	data.Prefix = utils.S3EncodeName(prefix, encodingType)
	data.Marker = utils.S3EncodeName(marker, encodingType)
	data.Delimiter = utils.S3EncodeName(delimiter, encodingType)
	data.MaxKeys = maxKeys
	data.NextMarker = utils.S3EncodeName(resp.NextMarker, encodingType)
	data.IsTruncated = resp.IsTruncated

	prefixes := make([]response.CommonPrefix, 0, len(resp.Prefixes))
	for _, prefix := range resp.Prefixes {
		prefixItem := response.CommonPrefix{}
		prefixItem.Prefix = utils.S3EncodeName(prefix, encodingType)
		prefixes = append(prefixes, prefixItem)
	}
	data.CommonPrefixes = prefixes
	return data
}

// generate multi objects delete response.
func generateMultiDeleteResponse(quiet bool, deletedObjects []datatypes.DeletedObject, errs []response.DeleteError) response.DeleteObjectsResponse {
	deleteResp := response.DeleteObjectsResponse{}
	if !quiet {
		deleteResp.DeletedObjects = deletedObjects
	}
	deleteResp.Errors = errs
	return deleteResp
}