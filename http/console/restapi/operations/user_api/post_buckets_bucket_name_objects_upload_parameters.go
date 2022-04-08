package user_api

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

// NewPostBucketsBucketNameObjectsUploadParams creates a new PostBucketsBucketNameObjectsUploadParams object
//
// There are no default values defined in the spec.
func NewPostBucketsBucketNameObjectsUploadParams() PostBucketsBucketNameObjectsUploadParams {

	return PostBucketsBucketNameObjectsUploadParams{}
}

// PostBucketsBucketNameObjectsUploadParams contains all the bound params for the post buckets bucket name objects upload operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostBucketsBucketNameObjectsUpload
type PostBucketsBucketNameObjectsUploadParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	BucketName string
	/*
	  In: query
	*/
	Prefix *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostBucketsBucketNameObjectsUploadParams() beforehand.
func (o *PostBucketsBucketNameObjectsUploadParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	rBucketName, rhkBucketName, _ := route.Params.GetOK("bucket_name")
	if err := o.bindBucketName(rBucketName, rhkBucketName, route.Formats); err != nil {
		res = append(res, err)
	}

	qPrefix, qhkPrefix, _ := qs.GetOK("prefix")
	if err := o.bindPrefix(qPrefix, qhkPrefix, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindBucketName binds and validates parameter BucketName from path.
func (o *PostBucketsBucketNameObjectsUploadParams) bindBucketName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.BucketName = raw

	return nil
}

// bindPrefix binds and validates parameter Prefix from query.
func (o *PostBucketsBucketNameObjectsUploadParams) bindPrefix(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.Prefix = &raw

	return nil
}
