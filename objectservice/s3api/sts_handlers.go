package s3api

import (
	"context"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"net/http"
	"strconv"
	"time"
)

const (
	parentClaim = "parent"
	expClaim    = "exp"
)

// AssumeRole - implementation of AWS STS API AssumeRole to get temporary
// credentials for regular users .
// https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html
func (s3a *s3ApiServer) AssumeRole(w http.ResponseWriter, r *http.Request) {
	// Check auth here (otherwise r.Form will have unexpected values from
	// the call to `parseForm` below), but return failure only after we are
	// able to validate that it is a valid STS request, so that we are able
	// to send an appropriate audit log.
	user, isErrCodeSTS, stsErr := s3a.checkAssumeRoleAuth(r.Context(), r)

	if err := parseForm(r); err != nil {
		response.WriteSTSErrorResponse(r.Context(), w, true, apierrors.ErrSTSInvalidParameterValue, err)
		return
	}

	if r.Form.Get(consts.StsVersion) != consts.StsAPIVersion {
		response.WriteSTSErrorResponse(r.Context(), w, true, apierrors.ErrSTSMissingParameter, fmt.Errorf("invalid STS API version %s3a, expecting %s3a", r.Form.Get(consts.StsAPIVersion), consts.StsAPIVersion))
		return
	}

	action := r.Form.Get(consts.StsAction)
	switch action {
	case consts.AssumeRole:
	default:
		response.WriteSTSErrorResponse(r.Context(), w, true, apierrors.ErrSTSInvalidParameterValue, fmt.Errorf("unsupported action %s3a", action))
		return
	}

	// Validate the authentication result here so that failures will be
	// audit-logged.
	if stsErr != apierrors.ErrSTSNone {
		response.WriteSTSErrorResponse(r.Context(), w, isErrCodeSTS, stsErr, nil)
		return
	}
	expiration, err := getDefaultExpiration(r.Form.Get(consts.StsDurationSeconds))
	if err != nil {
		response.WriteSTSErrorResponse(r.Context(), w, true, apierrors.ErrSTSInvalidParameterValue, err)
		return
	}
	m := map[string]interface{}{
		expClaim:    expiration,
		parentClaim: user.AccessKey,
	}

	secret := s3a.authSys.AdminCred.SecretKey
	cred, err := auth.GetNewCredentialsWithMetadata(m, secret)
	if err != nil {
		response.WriteSTSErrorResponse(r.Context(), w, true, apierrors.ErrSTSInternalError, err)
		return
	}
	// Set the parent of the temporary access key, so that it's access
	// policy is inherited from `user.AccessKey`.
	cred.ParentUser = user.AccessKey
	// Set the newly generated credentials.
	if err = s3a.authSys.Iam.SetTempUser(r.Context(), cred.AccessKey, cred, ""); err != nil {
		response.WriteSTSErrorResponse(r.Context(), w, true, apierrors.ErrSTSInternalError, err)
		return
	}
	assumeRoleResponse := &response.AssumeRoleResponse{
		Result: response.AssumeRoleResult{
			Credentials: cred,
		},
	}
	assumeRoleResponse.ResponseMetadata.RequestID = w.Header().Get(consts.AmzRequestID)
	response.WriteSuccessResponseXML(w, r, assumeRoleResponse)
}
func (s3a *s3ApiServer) checkAssumeRoleAuth(ctx context.Context, r *http.Request) (user auth.Credentials, isErrCodeSTS bool, stsErr apierrors.STSErrorCode) {
	if !iam.IsRequestSignatureV4(r) {
		return user, true, apierrors.ErrSTSAccessDenied
	}

	s3Err := s3a.authSys.IsReqAuthenticated(ctx, r, consts.DefaultRegion, iam.ServiceSTS)
	if s3Err != apierrors.ErrNone {
		return user, false, apierrors.STSErrorCode(s3Err)
	}

	user, _, s3Err = s3a.authSys.GetReqAccessKeyV4(r, consts.DefaultRegion, iam.ServiceSTS)
	if s3Err != apierrors.ErrNone {
		return user, false, apierrors.STSErrorCode(s3Err)
	}

	// Session tokens are not allowed in STS AssumeRole requests.
	if getSessionToken(r) != "" {
		return user, true, apierrors.ErrSTSAccessDenied
	}

	return user, true, apierrors.ErrSTSNone
}

// Fetch the security token set by the client.
func getSessionToken(r *http.Request) (token string) {
	token = r.Header.Get(consts.AmzSecurityToken)
	if token != "" {
		return token
	}
	return r.Form.Get(consts.AmzSecurityToken)
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

// getDefaultExpiration - returns the expiration seconds expected.
func getDefaultExpiration(dsecs string) (time.Duration, error) {
	defaultExpiryDuration := time.Duration(60) * time.Minute // Defaults to 1hr.
	if dsecs != "" {
		expirySecs, err := strconv.ParseInt(dsecs, 10, 64)
		if err != nil {
			return 0, errors.New("invalid token expiry")
		}

		// The duration, in seconds, of the role session.
		// The value can range from 900 seconds (15 minutes)
		// up to 365 days.
		if expirySecs < 900 || expirySecs > 31536000 {
			return 0, errors.New("invalid token expiry")
		}

		defaultExpiryDuration = time.Duration(expirySecs) * time.Second
	}
	return defaultExpiryDuration, nil
}
