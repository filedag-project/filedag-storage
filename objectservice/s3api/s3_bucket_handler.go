package s3api

import (
	"encoding/xml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/iam/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	logging "github.com/ipfs/go-log/v2"
	"io"
	"net/http"
	"path"
)

var log = logging.Logger("server")

//ListBucketsHandler ListBuckets Handler
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBuckets.html
func (s3a *s3ApiServer) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.ListAllMyBucketsAction, "", "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	bucketMetas, erro := s3a.bmSys.GetAllBucketOfUser(ctx, cred.AccessKey)
	if erro != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, erro))
		return
	}
	var buckets []*s3.Bucket
	for _, b := range bucketMetas {
		buckets = append(buckets, &s3.Bucket{
			Name:         aws.String(b.Name),
			CreationDate: aws.Time(b.Created),
		})
	}

	resp := response.ListAllMyBucketsResult{
		Owner: &s3.Owner{
			ID:          aws.String(consts.DefaultOwnerID),
			DisplayName: aws.String(consts.DisplayName),
		},
		Buckets: buckets,
	}

	response.WriteSuccessResponseXML(w, r, resp)
}

// GetBucketLocationHandler - GET Bucket location.
// -------------------------
// This operation returns bucket location.
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketLocation.html
func (s3a *s3ApiServer) GetBucketLocationHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	ctx := r.Context()
	_, _, err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.ListAllMyBucketsAction, bucket, "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	bucketMetas, erro := s3a.bmSys.GetBucketMeta(ctx, bucket)
	if erro != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, erro))
		return
	}

	// Generate response.
	encodedSuccessResponse := response.LocationResponse{
		Location: bucketMetas.Region,
	}

	// Write success response.
	response.WriteSuccessResponseXML(w, r, encodedSuccessResponse)
}

//PutBucketHandler put a bucket
func (s3a *s3ApiServer) PutBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	ctx := r.Context()
	log.Infof("PutBucketHandler %s", bucket)
	region, _ := parseLocationConstraint(r)
	// avoid duplicated buckets
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.CreateBucketAction, bucket, "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	erro := s3a.bmSys.CreateBucket(ctx, bucket, cred.AccessKey, region)
	if erro != nil {
		log.Errorf("PutBucketHandler create bucket error:%v", err)
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, erro))
		return
	}
	// Make sure to add Location information here only for bucket
	if cp := pathClean(r.URL.Path); cp != "" {
		w.Header().Set(consts.Location, cp) // Clean any trailing slashes.
	}

	response.WriteSuccessResponseEmpty(w, r)
}

// HeadBucketHandler - HEAD Bucket
// ----------
// This operation is useful to determine if a bucket exists.
// The operation returns a 200 OK if the bucket exists and you
// have permission to access it. Otherwise, the operation might
// return responses such as 404 Not Found and 403 Forbidden.
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_HeadBucket.html
func (s3a *s3ApiServer) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	log.Infof("HeadBucketHandler %s", bucket)
	// avoid duplicated buckets
	_, _, err := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.HeadBucketAction, bucket, "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	if ok := s3a.bmSys.HasBucket(r.Context(), bucket); !ok {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucket)
		return
	}

	response.WriteSuccessResponseEmpty(w, r)
}

// DeleteBucketHandler delete Bucket
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucket.html
func (s3a *s3ApiServer) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	ctx := r.Context()
	log.Infof("DeleteBucketHandler %s", bucket)
	_, _, err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.DeleteBucketAction, bucket, "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	errc := s3a.bmSys.DeleteBucket(ctx, bucket)
	if errc != nil {
		log.Errorf("DeleteBucketHandler delete bucket err: %v", err)
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, errc))
		return
	}
	response.WriteSuccessNoContent(w)
}

// GetBucketAclHandler Get Bucket ACL
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketAcl.html
func (s3a *s3ApiServer) GetBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	// collect parameters
	bucket, _ := getBucketAndObject(r)
	log.Infof("GetBucketAclHandler %s", bucket)
	cred, _, err := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetBucketPolicyAction, bucket, "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}
	resp := response.AccessControlPolicy{}
	id := cred.AccessKey
	if resp.Owner.DisplayName == "" {
		resp.Owner.DisplayName = cred.AccessKey
		resp.Owner.ID = id
	}
	resp.AccessControlList.Grant = append(resp.AccessControlList.Grant, response.Grant{
		Grantee: response.Grantee{
			ID:          id,
			DisplayName: cred.AccessKey,
			Type:        "CanonicalUser",
			XMLXSI:      "CanonicalUser",
			XMLNS:       "http://www.w3.org/2001/XMLSchema-instance"},
		Permission: "FULL_CONTROL", //todo change
	})
	response.WriteSuccessResponseXML(w, r, resp)
}

// GetBucketCorsHandler Get bucket CORS
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketCors.html
func (s3a *s3ApiServer) GetBucketCorsHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteErrorResponse(w, r, apierrors.ErrNoSuchCORSConfiguration)
}

// PutBucketCorsHandler Put bucket CORS
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketCors.html
func (s3a *s3ApiServer) PutBucketCorsHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteErrorResponse(w, r, apierrors.ErrNotImplemented)
}

// DeleteBucketCorsHandler Delete bucket CORS
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketCors.html
func (s3a *s3ApiServer) DeleteBucketCorsHandler(w http.ResponseWriter, r *http.Request) {
	response.WriteErrorResponse(w, r, http.StatusNoContent)
}

// PutBucketAclHandler Put bucket ACL
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketAcl.html
func (s3a *s3ApiServer) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)

	// Allow putBucketACL if policy action is set, since this is a dummy call
	// we are simply re-purposing the bucketPolicyAction.
	_, _, err := s3a.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.PutBucketPolicyAction, bucket, "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	aclHeader := r.Header.Get(consts.AmzACL)
	if aclHeader == "" {
		acl := &response.AccessControlPolicy{}
		if errc := utils.XmlDecoder(r.Body, acl, r.ContentLength); errc != nil {
			if errc == io.EOF {
				response.WriteErrorResponse(w, r, apierrors.ErrMissingSecurityHeader)
				return
			}
			response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
			return
		}

		if len(acl.AccessControlList.Grant) == 0 {
			response.WriteErrorResponse(w, r, apierrors.ErrNotImplemented)
			return
		}

		if acl.AccessControlList.Grant[0].Permission != "FULL_CONTROL" {
			response.WriteErrorResponse(w, r, apierrors.ErrNotImplemented)
			return
		}
	}

	if aclHeader != "" && aclHeader != "private" {
		response.WriteErrorResponse(w, r, apierrors.ErrNotImplemented)
		return
	}
}

// PutBucketTaggingHandler
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketTagging.html
func (s3a *s3ApiServer) PutBucketTaggingHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	ctx := r.Context()
	log.Infof("DeleteBucketHandler %s", bucket)
	_, _, err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.DeleteBucketAction, bucket, "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	tags, err1 := unmarshalXML(io.LimitReader(r.Body, r.ContentLength), false)
	if err1 != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrMalformedXML)
		return
	}

	if err1 = s3a.bmSys.UpdateBucketTagging(ctx, bucket, tags); err1 != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err1))
		return
	}

	// Write success response.
	response.WriteSuccessResponseHeadersOnly(w, r)
}

// GetBucketTaggingHandler
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketTagging.html
func (s3a *s3ApiServer) GetBucketTaggingHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	ctx := r.Context()
	log.Infof("DeleteBucketHandler %s", bucket)
	_, _, err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.DeleteBucketAction, bucket, "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	tags, err1 := s3a.bmSys.GetTaggingConfig(ctx, bucket)
	if err1 != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err1))
		return
	}
	configData, err2 := xml.Marshal(tags)
	if err2 != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrMalformedXML)
		return
	}

	// Write success response.
	response.WriteSuccessResponseXML(w, r, configData)
}

// DeleteBucketTaggingHandler
//https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketTagging.html
func (s3a *s3ApiServer) DeleteBucketTaggingHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := getBucketAndObject(r)
	ctx := r.Context()
	log.Infof("DeleteBucketHandler %s", bucket)
	_, _, err := s3a.authSys.CheckRequestAuthTypeCredential(ctx, r, s3action.DeleteBucketAction, bucket, "")
	if err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, err)
		return
	}

	if err1 := s3a.bmSys.DeleteBucketTagging(ctx, bucket); err1 != nil {
		response.WriteErrorResponse(w, r, apierrors.ToApiError(ctx, err1))
		return
	}

	// Write success response.
	response.WriteSuccessResponseHeadersOnly(w, r)
}

// Parses location constraint from the incoming reader.
func parseLocationConstraint(r *http.Request) (location string, s3Error apierrors.ErrorCode) {
	// If the request has no body with content-length set to 0,
	// we do not have to validate location constraint. Bucket will
	// be created at default region.
	locationConstraint := createBucketLocationConfiguration{}
	err := utils.XmlDecoder(r.Body, &locationConstraint, r.ContentLength)
	if err != nil && r.ContentLength != 0 {
		// Treat all other failures as XML parsing errors.
		return "", apierrors.ErrMalformedXML
	} // else for both err as nil or io.EOF
	location = locationConstraint.Location
	if location == "" {
		location = consts.DefaultRegion
	}
	return location, apierrors.ErrNone
}

// createBucketConfiguration container for bucket configuration request from client.
// Used for parsing the location from the request body for Makebucket.
type createBucketLocationConfiguration struct {
	XMLName  xml.Name `xml:"CreateBucketConfiguration" json:"-"`
	Location string   `xml:"LocationConstraint"`
}

// pathClean is like path.Clean but does not return "." for
// empty inputs, instead returns "empty" as is.
func pathClean(p string) string {
	cp := path.Clean(p)
	if cp == "." {
		return ""
	}
	return cp
}

func unmarshalXML(reader io.Reader, isObject bool) (*store.Tags, error) {
	tagging := &store.Tags{
		TagSet: &store.TagSet{
			TagMap:   make(map[string]string),
			IsObject: isObject,
		},
	}

	if err := xml.NewDecoder(reader).Decode(tagging); err != nil {
		return nil, err
	}

	return tagging, nil
}
