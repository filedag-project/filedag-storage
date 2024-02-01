package iamapi

import (
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"math"
	"net/http"
)

// AccountInfos returns all user usage
func (iamApi *iamApiServer) AccountInfos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetUserInfoAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(s3err))
		return
	}
	if cred.AccessKey != iamApi.authSys.AdminCred.AccessKey {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	infos, err := iamApi.getAllUserInfos(ctx)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	bucketInfos := iamApi.bucketInfoFunc(ctx, iamApi.authSys.AdminCred.AccessKey)
	var useStorageCapacity uint64
	for _, bi := range bucketInfos {
		useStorageCapacity += bi.Size
	}
	polices, err := iamApi.authSys.Iam.GetUserPolices(r.Context(), iamApi.authSys.AdminCred.AccessKey)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	infos = append(infos, iam.UserInfo{
		AccountName:          iamApi.authSys.AdminCred.AccessKey,
		TotalStorageCapacity: math.MaxUint64,
		BucketInfos:          bucketInfos,
		UseStorageCapacity:   useStorageCapacity,
		PolicyName:           polices,
		Status:               iam.AccountEnabled,
	})
	response.WriteSuccessResponseJSON(w, r, infos)
}

// RequestOverview returns all user request
func (iamApi *iamApiServer) RequestOverview(w http.ResponseWriter, r *http.Request) {
	//to implement
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetUserInfoAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(s3err))
		return
	}
	if cred.AccessKey != iamApi.authSys.AdminCred.AccessKey {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	stats, err := iamApi.stats.GetCurrentStats(r.Context())
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	response.WriteSuccessResponseJSON(w, r, stats)
}

// StorePoolStats Store Pool Stats
func (iamApi *iamApiServer) StorePoolStats(w http.ResponseWriter, r *http.Request) {
	//to implement
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetPoolStatsAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(s3err))
		return
	}
	if cred.AccessKey != iamApi.authSys.AdminCred.AccessKey {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	storePoolStats, err := iamApi.storePoolStatsFunc(r.Context())
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ToApiError(r.Context(), err)))
		return
	}
	response.WriteSuccessResponseJSON(w, r, storePoolStats)
}

// ChangePassword change password
func (iamApi *iamApiServer) ChangePassword(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.ChangePassWordAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}

	secret := r.FormValue(newSecretKey)
	username := r.FormValue(accessKey)
	oldSecret := r.FormValue(oldSecretKey)
	credChange, ok := iamApi.authSys.Iam.GetUser(r.Context(), username)
	if !ok {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrNoSuchUser))
		return
	}
	if cred.AccessKey != iamApi.authSys.AdminCred.AccessKey {
		if credChange.SecretKey != oldSecret || username != cred.AccessKey {
			response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
			return
		}
	}
	if !auth.IsSecretKeyValid(secret) {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidRequestParameter))
		return
	}

	credChange.SecretKey = secret
	m := make(map[string]interface{})
	var err error
	credChange.SessionToken, err = auth.JWTSignWithAccessKey(username, m, auth.DefaultSecretKey)
	err = iamApi.authSys.Iam.UpdateUser(r.Context(), credChange)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	response.WriteSuccessResponseJSON(w, r, nil)
}

// SetStatus set user status
func (iamApi *iamApiServer) SetStatus(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.SetStatusAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}

	username := r.FormValue(accessKey)
	status := r.FormValue(accountStatus)
	c, _ := iamApi.authSys.Iam.GetUser(r.Context(), username)
	if c.AccessKey == "" {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrNoSuchUser))
		return
	}
	if username != cred.AccessKey && cred.AccessKey != iamApi.authSys.AdminCred.AccessKey {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	switch status {
	case auth.AccountOn, auth.AccountOff:
	default:
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidRequestParameter))
		return
	}
	if c.Status == status {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidRequestParameter))
		return
	}
	c.Status = status
	err := iamApi.authSys.Iam.UpdateUser(r.Context(), c)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	response.WriteSuccessResponseJSON(w, r, nil)
}
