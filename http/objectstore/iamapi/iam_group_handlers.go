package iamapi

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"net/http"
)

func (iamApi *iamApiServer) CreatGroup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
}

func (iamApi *iamApiServer) GetGroup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
}

func (iamApi *iamApiServer) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
}

func (iamApi *iamApiServer) ListGroups(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
}
