package sts_api

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/consts"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/set"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	parentClaim = "parent"
	expClaim    = "exp"
)

// stsAPIHandlers implements and provides http handlers for AWS STS API.
type stsAPIHandlers struct{}

// registerSTSRouter - registers AWS STS compatible APIs.
func registerSTSRouter(router *mux.Router) {
	// Initialize STS.
	sts := &stsAPIHandlers{}

	// STS Router
	stsRouter := router.NewRoute().PathPrefix(consts.SlashSeparator).Subrouter()

	// Assume roles with no JWT, handles AssumeRole.
	stsRouter.Methods(http.MethodPost).MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
		ctypeOk := set.MatchSimple("application/x-www-form-urlencoded*", r.Header.Get(consts.ContentType))
		authOk := set.MatchSimple(consts.SignV4Algorithm+"*", r.Header.Get(consts.Authorization))
		noQueries := len(r.URL.RawQuery) == 0
		return ctypeOk && authOk && noQueries
	}).HandlerFunc(sts.AssumeRole)
}
func parseForm(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	for k, v := range r.PostForm {
		if _, ok := r.Form[k]; !ok {
			r.Form[k] = v
		}
	}
	return nil
}

// AssumeRole - implementation of AWS STS API AssumeRole to get temporary
// credentials for regular users on Minio.
// https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html
func (sts *stsAPIHandlers) AssumeRole(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	// Check auth here (otherwise r.Form will have unexpected values from
	// the call to `parseForm` below), but return failure only after we are
	// able to validate that it is a valid STS request, so that we are able
	// to send an appropriate audit log.
	user, isErrCodeSTS, stsErr := checkAssumeRoleAuth(ctx, r)

	if err := parseForm(r); err != nil {
		response.WriteSTSErrorResponse(ctx, w, true, api_errors.ErrSTSInvalidParameterValue, err)
		return
	}

	if r.Form.Get(consts.StsVersion) != consts.StsAPIVersion {
		response.WriteSTSErrorResponse(ctx, w, true, api_errors.ErrSTSMissingParameter, fmt.Errorf("invalid STS API version %s, expecting %s", r.Form.Get(consts.StsAPIVersion), consts.StsAPIVersion))
		return
	}

	action := r.Form.Get(consts.StsAction)
	switch action {
	case consts.AssumeRole:
	default:
		response.WriteSTSErrorResponse(ctx, w, true, api_errors.ErrSTSInvalidParameterValue, fmt.Errorf("unsupported action %s", action))
		return
	}

	// Validate the authentication result here so that failures will be
	// audit-logged.
	if stsErr != api_errors.ErrSTSNone {
		response.WriteSTSErrorResponse(ctx, w, isErrCodeSTS, stsErr, nil)
		return
	}
	defaultExpiryDuration := time.Duration(60) * time.Minute // Defaults to 1hr.

	m := map[string]interface{}{
		expClaim:    time.Now().UTC().Add(defaultExpiryDuration).Unix(),
		parentClaim: user.AccessKey,
	}

	secret := auth.GetDefaultActiveCred().SecretKey
	cred, err := auth.GetNewCredentialsWithMetadata(m, secret)
	if err != nil {
		response.WriteSTSErrorResponse(ctx, w, true, api_errors.ErrSTSInternalError, err)
		return
	}
	assumeRoleResponse := &response.AssumeRoleResponse{
		Result: response.AssumeRoleResult{
			Credentials: cred,
		},
	}
	assumeRoleResponse.ResponseMetadata.RequestID = w.Header().Get(consts.AmzRequestID)
	response.WriteSuccessResponseXML(w, r, response.EncodeResponse(assumeRoleResponse))
}
func checkAssumeRoleAuth(ctx context.Context, r *http.Request) (user auth.Credentials, isErrCodeSTS bool, stsErr api_errors.STSErrorCode) {
	if !iam.IsRequestSignatureV4(r) {
		return user, true, api_errors.ErrSTSAccessDenied
	}

	s3Err := iam.IsReqAuthenticated(ctx, r, consts.DefaultRegion, iam.ServiceSTS)
	if s3Err != api_errors.ErrNone {
		return user, false, api_errors.STSErrorCode(s3Err)
	}

	user, _, s3Err = iam.GetReqAccessKeyV4(r, consts.DefaultRegion, iam.ServiceSTS)
	if s3Err != api_errors.ErrNone {
		return user, false, api_errors.STSErrorCode(s3Err)
	}

	// Session tokens are not allowed in STS AssumeRole requests.
	if getSessionToken(r) != "" {
		return user, true, api_errors.ErrSTSAccessDenied
	}

	return user, true, api_errors.ErrSTSNone
}

// Fetch the security token set by the client.
func getSessionToken(r *http.Request) (token string) {
	token = r.Header.Get(consts.AmzSecurityToken)
	if token != "" {
		return token
	}
	return r.Form.Get(consts.AmzSecurityToken)
}
