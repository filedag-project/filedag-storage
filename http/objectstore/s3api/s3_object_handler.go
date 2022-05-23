package s3api

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/s3action"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/filedag-project/filedag-storage/http/objectstore/store"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/etag"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils/hash"
	"github.com/gorilla/mux"
	"io/ioutil"
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
		response.WriteErrorResponse(w, r, api_errors.ErrInvalidCopySource)
		return
	}
	log.Infof("PutObjectHandler %s %s", bucket, object)
	clientETag, err := etag.FromContentMD5(r.Header)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInvalidDigest)
		return
	}
	// if Content-Length is unknown/missing, deny the request
	size := r.ContentLength
	rAuthType := iam.GetRequestAuthType(r)
	if rAuthType == iam.AuthTypeStreamingSigned {
		if sizeStr, ok := r.Header[consts.AmzDecodedContentLength]; ok {
			if sizeStr[0] == "" {
				response.WriteErrorResponse(w, r, api_errors.ErrMissingContentLength)
				return
			}
			size, err = strconv.ParseInt(sizeStr[0], 10, 64)
			if err != nil {
				log.Errorf("ParseInt err:%v", err)
				response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
				return
			}
		}
	}
	if size == -1 {
		response.WriteErrorResponse(w, r, api_errors.ErrMissingContentLength)
		return
	}
	// maximum Upload size for objects in a single operation
	if size > consts.MaxObjectSize {
		response.WriteErrorResponse(w, r, api_errors.ErrEntityTooLarge)
		return
	}
	cred, _, s3err := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.PutObjectAction, "testbuckets", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3err)
		return
	}
	if !s3a.authSys.PolicySys.Head(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}
	dataReader := r.Body

	hashReader, err1 := hash.NewReader(dataReader, size, clientETag.String(), r.Header.Get(consts.AmzContentSha256), size)
	if err1 != nil {
		log.Errorf("PutObjectHandler NewReader err:%v", err1)
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	defer dataReader.Close()
	objInfo, err2 := s3a.store.StoreObject(r.Context(), cred.AccessKey, bucket, object, hashReader)
	if err2 != nil {
		log.Errorf("PutObjectHandler StoreObject err:%v", err1)
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
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
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.authSys.PolicySys.Head(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}
	userName := cred.AccessKey
	if cred.AccessKey == "" {
		meta, err := s3a.authSys.PolicySys.GetMeta(bucket, cred.AccessKey)
		if err != nil {
			response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
			return
		}
		userName = meta.Owner
	}
	objInfo, reader, err := s3a.store.GetObject(r.Context(), userName, bucket, object)
	if err != nil {
		log.Errorf("GetObjectHandler GetObject err:%v", err)
		response.WriteErrorResponseHeadersOnly(w, r, api_errors.ErrInternalError)
		return
	}
	w.Header().Set(consts.AmzServerSideEncryption, consts.AmzEncryptionAES)

	response.SetObjectHeaders(w, r, objInfo)
	r1, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Errorf("GetObjectHandler reader readAll err:%v", err)
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	w.Header().Set(consts.ContentLength, strconv.Itoa(len(r1)))
	response.SetHeadGetRespHeaders(w, r.Form)
	_, err = w.Write(r1)
	if err != nil {
		log.Errorf("GetObjectHandler header write err:%v", err)
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
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
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.authSys.PolicySys.Head(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}
	objInfo, ok := s3a.store.HasObject(r.Context(), cred.AccessKey, bucket, object)
	if !ok {
		response.WriteErrorResponseHeadersOnly(w, r, api_errors.ErrNoSuchKey)
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
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.authSys.PolicySys.Head(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}
	objInfo, ok := s3a.store.HasObject(r.Context(), cred.AccessKey, bucket, object)
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchKey)
		return
	}
	err := s3a.store.DeleteObject(cred.AccessKey, bucket, object)
	if err != nil {
		log.Errorf("DeleteObjectHandler DeleteObject  err:%v", err)
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	setPutObjHeaders(w, objInfo, true)
	response.WriteSuccessNoContent(w)
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
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.authSys.PolicySys.Head(dstBucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}

	// Copy source path.
	cpSrcPath, err := url.QueryUnescape(r.Header.Get("X-Amz-Copy-Source"))
	if err != nil {
		// Save unescaped string as is.
		cpSrcPath = r.Header.Get(consts.AmzCopySource)
	}

	srcBucket, srcObject := pathToBucketAndObject(cpSrcPath)
	if !s3a.authSys.PolicySys.Head(srcBucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}
	log.Infof("CopyObjectHandler %s %s => %s %s", srcBucket, srcObject, dstBucket, dstObject)
	_, i, err := s3a.store.GetObject(r.Context(), cred.AccessKey, srcBucket, srcObject)
	if err != nil {
		log.Errorf("CopyObjectHandler StoreObject err:%v", err)
		response.WriteErrorResponseHeadersOnly(w, r, api_errors.ErrNoSuchKey)
		return
	}
	if (srcBucket == dstBucket && srcObject == dstObject || cpSrcPath == "") && isReplace(r) {
		object, err := s3a.store.StoreObject(r.Context(), cred.AccessKey, dstBucket, dstObject, i)
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
		response.WriteErrorResponse(w, r, api_errors.ErrInvalidCopySource)
		return
	}

	if srcBucket == dstBucket && srcObject == dstObject {
		response.WriteErrorResponse(w, r, api_errors.ErrInvalidCopyDest)
		return
	}
	obj, err := s3a.store.StoreObject(r.Context(), cred.AccessKey, dstBucket, dstObject, i)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
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
	cred, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetObjectAction, bucket, "")
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	if !s3a.authSys.PolicySys.Head(bucket, cred.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}
	// Extract all the litsObjectsV1 query params to their native values.
	prefix, marker, delimiter, maxKeys, encodingType, s3Error := getListObjectsV1Args(r.Form)
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	objs, err := s3a.store.ListObject(cred.AccessKey, bucket)
	if err != nil {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}

	listObjectsInfo := response.ListObjectsInfo{
		IsTruncated: false,
		NextMarker:  "",
		Objects:     objs,
		Prefixes:    nil,
	}
	resp := generateListObjectsV1Response(bucket, prefix, marker, delimiter, encodingType, maxKeys, listObjectsInfo)
	// Write success response.
	response.WriteSuccessResponseXML(w, r, resp)
}
func (s3a *s3ApiServer) ListObjectsV2Handler(w http.ResponseWriter, r *http.Request) {
	bucket, object := getBucketAndObject(r)

	// Check for auth type to return S3 compatible error.
	// type to return the correct error (NoSuchKey vs AccessDenied)
	cerd, _, s3Error := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetObjectAction, bucket, object)
	if s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	urlValues := r.Form
	if !s3a.authSys.PolicySys.Head(bucket, cerd.AccessKey) {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucket)
		return
	}
	// Extract all the listObjectsV2 query params to their native values.
	prefix, token, startAfter, delimiter, fetchOwner, maxKeys, encodingType, errCode := getListObjectsV2Args(urlValues)
	if errCode != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, errCode)
		return
	}
	// Validate the query params before beginning to serve the request.
	// fetch-owner is not validated since it is a boolean
	if s3Error := validateListObjectsArgs(token, delimiter, encodingType, maxKeys); s3Error != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, s3Error)
		return
	}
	var (
		listObjectsV2Info store.ListObjectsV2Info
		err               error
	)

	// Inititate a list objects operation based on the input params.
	// On success would return back ListObjectsInfo object to be
	// marshaled into S3 compatible XML header.
	listObjectsV2Info, err = s3a.store.ListObjectsV2(r.Context(), cerd.AccessKey, bucket, prefix, token, delimiter, maxKeys, fetchOwner, startAfter)

	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
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
func getListObjectsV1Args(values url.Values) (prefix, marker, delimiter string, maxkeys int, encodingType string, errCode api_errors.ErrorCode) {
	errCode = api_errors.ErrNone

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.Atoi(values.Get("max-keys")); err != nil {
			errCode = api_errors.ErrInvalidMaxKeys
			return
		}
	} else {
		maxkeys = consts.MaxObjectList
	}

	prefix = trimLeadingSlash(values.Get("prefix"))
	marker = trimLeadingSlash(values.Get("marker"))
	delimiter = values.Get("delimiter")
	encodingType = values.Get("encoding-type")
	return
}

// Parse bucket url queries for ListObjects V2.
func getListObjectsV2Args(values url.Values) (prefix, token, startAfter, delimiter string, fetchOwner bool, maxkeys int, encodingType string, errCode api_errors.ErrorCode) {
	errCode = api_errors.ErrNone

	// The continuation-token cannot be empty.
	if val, ok := values["continuation-token"]; ok {
		if len(val[0]) == 0 {
			errCode = api_errors.ErrInvalidToken
			return
		}
	}

	if values.Get("max-keys") != "" {
		var err error
		if maxkeys, err = strconv.Atoi(values.Get("max-keys")); err != nil {
			errCode = api_errors.ErrInvalidMaxKeys
			return
		}
	} else {
		maxkeys = consts.MaxObjectList
	}

	prefix = trimLeadingSlash(values.Get("prefix"))
	startAfter = trimLeadingSlash(values.Get("start-after"))
	delimiter = values.Get("delimiter")
	fetchOwner = values.Get("fetch-owner") == "true"
	encodingType = values.Get("encoding-type")

	if token = values.Get("continuation-token"); token != "" {
		decodedToken, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			errCode = api_errors.ErrIncorrectContinuationToken
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
// Special conditions required by MinIO server are as below
// - delimiter if set should be equal to '/', otherwise the request is rejected.
// - marker if set should have a common prefix with 'prefix' param, otherwise
//   the request is rejected.
func validateListObjectsArgs(marker, delimiter, encodingType string, maxKeys int) api_errors.ErrorCode {
	// Max keys cannot be negative.
	if maxKeys < 0 {
		return api_errors.ErrInvalidMaxKeys
	}

	if encodingType != "" {
		// AWS S3 spec only supports 'url' encoding type
		if !strings.EqualFold(encodingType, "url") {
			return api_errors.ErrInvalidEncodingMethod
		}
	}

	return api_errors.ErrNone
}

// GenerateListObjectsV2Response Generates an ListObjectsV2 response for the said bucket with other enumerated options.
func GenerateListObjectsV2Response(bucket, prefix, token, nextToken, startAfter, delimiter, encodingType string, isTruncated bool, maxKeys int, objects []store.ObjectInfo, prefixes []string) response.ListObjectsV2Response {
	contents := make([]response.Object, 0, len(objects))
	a := consts.DefaultOwnerID
	b := "fds"
	owner := s3.Owner{
		ID:          &a,
		DisplayName: &b,
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
func generateListObjectsV1Response(bucket, prefix, marker, delimiter, encodingType string, maxKeys int, resp response.ListObjectsInfo) response.ListObjectsResponse {
	contents := make([]response.Object, 0, len(resp.Objects))
	a := consts.DefaultOwnerID
	b := "fds"
	owner := s3.Owner{
		ID:          &a,
		DisplayName: &b,
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
