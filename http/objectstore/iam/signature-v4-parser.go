/*
 * The following code tries to reverse engineer the Amazon S3 APIs,
 * and is mostly copied from minio implementation.
 */

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package iam

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// credentialHeader data type represents structured form of Credential
// string from authorization header.
type credentialHeader struct {
	accessKey string
	scope     struct {
		date    time.Time
		region  string
		service string
		request string
	}
}

// Return scope string.
func (c credentialHeader) getScope() string {
	return strings.Join([]string{
		c.scope.date.Format(yyyymmdd),
		c.scope.region,
		c.scope.service,
		c.scope.request,
	}, consts.SlashSeparator)
}

func getReqAccessKeyV4(r *http.Request, region string, stype serviceType) (auth.Credentials, bool, api_errors.ErrorCode) {
	ch, s3Err := parseCredentialHeader("Credential="+r.Form.Get(consts.AmzCredential), region, stype)
	if s3Err != api_errors.ErrNone {
		// Strip off the Algorithm prefix.
		v4Auth := strings.TrimPrefix(r.Header.Get("Authorization"), signV4Algorithm)
		authFields := strings.Split(strings.TrimSpace(v4Auth), ",")
		if len(authFields) != 3 {
			return auth.Credentials{}, false, api_errors.ErrMissingFields
		}
		ch, s3Err = parseCredentialHeader(authFields[0], region, stype)
		if s3Err != api_errors.ErrNone {
			return auth.Credentials{}, false, s3Err
		}
	}
	return checkKeyValid(r, ch.accessKey)
}

// parse credentialHeader string into its structured form.
func parseCredentialHeader(credElement string, region string, stype serviceType) (ch credentialHeader, aec api_errors.ErrorCode) {
	creds := strings.SplitN(strings.TrimSpace(credElement), "=", 2)
	if len(creds) != 2 {
		return ch, api_errors.ErrMissingFields
	}
	if creds[0] != "Credential" {
		return ch, api_errors.ErrMissingCredTag
	}
	credElements := strings.Split(strings.TrimSpace(creds[1]), consts.SlashSeparator)
	if len(credElements) < 5 {
		return ch, api_errors.ErrCredMalformed
	}
	accessKey := strings.Join(credElements[:len(credElements)-4], consts.SlashSeparator) // The access key may contain one or more `/`
	if !auth.IsAccessKeyValid(accessKey) {
		return ch, api_errors.ErrInvalidAccessKeyID
	}
	// Save access key id.
	cred := credentialHeader{
		accessKey: accessKey,
	}
	credElements = credElements[len(credElements)-4:]
	var e error
	cred.scope.date, e = time.Parse(yyyymmdd, credElements[0])
	if e != nil {
		return ch, api_errors.ErrMalformedCredentialDate
	}

	cred.scope.region = credElements[1]
	// Verify if region is valid.
	sRegion := cred.scope.region
	// Region is set to be empty, we use whatever was sent by the
	// request and proceed further. This is a work-around to address
	// an important problem for ListBuckets() getting signed with
	// different regions.
	if region == "" {
		region = sRegion
	}
	// Should validate region, only if region is set.
	if !isValidRegion(sRegion, region) {
		return ch, api_errors.ErrAuthorizationHeaderMalformed
	}
	if credElements[2] != string(stype) {
		switch stype {
		case serviceSTS:
			return ch, api_errors.ErrInvalidServiceSTS
		}
		return ch, api_errors.ErrInvalidServiceS3
	}
	cred.scope.service = credElements[2]
	if credElements[3] != "aws4_request" {
		return ch, api_errors.ErrInvalidRequestVersion
	}
	cred.scope.request = credElements[3]
	return cred, api_errors.ErrNone
}

// Parse signature from signature tag.
func parseSignature(signElement string) (string, api_errors.ErrorCode) {
	signFields := strings.Split(strings.TrimSpace(signElement), "=")
	if len(signFields) != 2 {
		return "", api_errors.ErrMissingFields
	}
	if signFields[0] != "Signature" {
		return "", api_errors.ErrMissingSignTag
	}
	if signFields[1] == "" {
		return "", api_errors.ErrMissingFields
	}
	signature := signFields[1]
	return signature, api_errors.ErrNone
}

// Parse slice of signed headers from signed headers tag.
func parseSignedHeader(signedHdrElement string) ([]string, api_errors.ErrorCode) {
	signedHdrFields := strings.Split(strings.TrimSpace(signedHdrElement), "=")
	if len(signedHdrFields) != 2 {
		return nil, api_errors.ErrMissingFields
	}
	if signedHdrFields[0] != "SignedHeaders" {
		return nil, api_errors.ErrMissingSignHeadersTag
	}
	if signedHdrFields[1] == "" {
		return nil, api_errors.ErrMissingFields
	}
	signedHeaders := strings.Split(signedHdrFields[1], ";")
	return signedHeaders, api_errors.ErrNone
}

// signValues data type represents structured form of AWS Signature V4 header.
type signValues struct {
	Credential    credentialHeader
	SignedHeaders []string
	Signature     string
}

// preSignValues data type represents structued form of AWS Signature V4 query string.
type preSignValues struct {
	signValues
	Date    time.Time
	Expires time.Duration
}

// Parses signature version '4' query string of the following form.
//
//   querystring = X-Amz-Algorithm=algorithm
//   querystring += &X-Amz-Credential= urlencode(accessKey + '/' + credential_scope)
//   querystring += &X-Amz-Date=date
//   querystring += &X-Amz-Expires=timeout interval
//   querystring += &X-Amz-SignedHeaders=signed_headers
//   querystring += &X-Amz-Signature=signature
//
// verifies if any of the necessary query params are missing in the presigned request.
func doesV4PresignParamsExist(query url.Values) api_errors.ErrorCode {
	v4PresignQueryParams := []string{consts.AmzAlgorithm, consts.AmzCredential, consts.AmzSignature, consts.AmzDate, consts.AmzSignedHeaders, consts.AmzExpires}
	for _, v4PresignQueryParam := range v4PresignQueryParams {
		if _, ok := query[v4PresignQueryParam]; !ok {
			return api_errors.ErrInvalidQueryParams
		}
	}
	return api_errors.ErrNone
}

// Parses all the presigned signature values into separate elements.
func parsePreSignV4(query url.Values, region string, stype serviceType) (psv preSignValues, aec api_errors.ErrorCode) {
	// verify whether the required query params exist.
	aec = doesV4PresignParamsExist(query)
	if aec != api_errors.ErrNone {
		return psv, aec
	}

	// Verify if the query algorithm is supported or not.
	if query.Get(consts.AmzAlgorithm) != signV4Algorithm {
		return psv, api_errors.ErrInvalidQuerySignatureAlgo
	}

	// Initialize signature version '4' structured header.
	preSignV4Values := preSignValues{}

	// Save credential.
	preSignV4Values.Credential, aec = parseCredentialHeader("Credential="+query.Get(consts.AmzCredential), region, stype)
	if aec != api_errors.ErrNone {
		return psv, aec
	}

	var e error
	// Save date in native time.Time.
	preSignV4Values.Date, e = time.Parse(iso8601Format, query.Get(consts.AmzDate))
	if e != nil {
		return psv, api_errors.ErrMalformedPresignedDate
	}

	// Save expires in native time.Duration.
	preSignV4Values.Expires, e = time.ParseDuration(query.Get(consts.AmzExpires) + "s")
	if e != nil {
		return psv, api_errors.ErrMalformedExpires
	}

	if preSignV4Values.Expires < 0 {
		return psv, api_errors.ErrNegativeExpires
	}

	// Check if Expiry time is less than 7 days (value in seconds).
	if preSignV4Values.Expires.Seconds() > 604800 {
		return psv, api_errors.ErrMaximumExpires
	}

	// Save signed headers.
	preSignV4Values.SignedHeaders, aec = parseSignedHeader("SignedHeaders=" + query.Get(consts.AmzSignedHeaders))
	if aec != api_errors.ErrNone {
		return psv, aec
	}

	// Save signature.
	preSignV4Values.Signature, aec = parseSignature("Signature=" + query.Get(consts.AmzSignature))
	if aec != api_errors.ErrNone {
		return psv, aec
	}

	// Return structed form of signature query string.
	return preSignV4Values, api_errors.ErrNone
}

// Parses signature version '4' header of the following form.
//
//    Authorization: algorithm Credential=accessKeyID/credScope, \
//            SignedHeaders=signedHeaders, Signature=signature
//
func parseSignV4(v4Auth string, region string, stype serviceType) (sv signValues, aec api_errors.ErrorCode) {
	// credElement is fetched first to skip replacing the space in access key.
	credElement := strings.TrimPrefix(strings.Split(strings.TrimSpace(v4Auth), ",")[0], signV4Algorithm)
	// Replace all spaced strings, some clients can send spaced
	// parameters and some won't. So we pro-actively remove any spaces
	// to make parsing easier.
	v4Auth = strings.ReplaceAll(v4Auth, " ", "")
	if v4Auth == "" {
		return sv, api_errors.ErrAuthHeaderEmpty
	}

	// Verify if the header algorithm is supported or not.
	if !strings.HasPrefix(v4Auth, signV4Algorithm) {
		return sv, api_errors.ErrSignatureVersionNotSupported
	}

	// Strip off the Algorithm prefix.
	v4Auth = strings.TrimPrefix(v4Auth, signV4Algorithm)
	authFields := strings.Split(strings.TrimSpace(v4Auth), ",")
	if len(authFields) != 3 {
		return sv, api_errors.ErrMissingFields
	}

	// Initialize signature version '4' structured header.
	signV4Values := signValues{}

	var s3Err api_errors.ErrorCode
	// Save credentail values.
	signV4Values.Credential, s3Err = parseCredentialHeader(strings.TrimSpace(credElement), region, stype)
	if s3Err != api_errors.ErrNone {
		return sv, s3Err
	}

	// Save signed headers.
	signV4Values.SignedHeaders, s3Err = parseSignedHeader(authFields[1])
	if s3Err != api_errors.ErrNone {
		return sv, s3Err
	}

	// Save signature.
	signV4Values.Signature, s3Err = parseSignature(authFields[2])
	if s3Err != api_errors.ErrNone {
		return sv, s3Err
	}

	// Return the structure here.
	return signV4Values, api_errors.ErrNone
}
