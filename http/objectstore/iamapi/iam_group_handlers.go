package iamapi

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"net/http"
	"strconv"
)

// CreatGroup
//Creates a new group.
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateGroup.html
func (iamApi *iamApiServer) CreatGroup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	groupName := r.FormValue("groupName")
	version := r.FormValue("version")
	atoi, _ := strconv.Atoi(version)
	err := iamApi.authSys.Iam.CreateGroup(ctx, groupName, atoi)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	var resp CreateGroupResponse
	response.WriteSuccessResponseXML(w, r, resp)
}

// GetGroup
//Returns a list of IAM users that are in the specified IAM group. You can paginate the results using the MaxItems and Marker parameters.
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetGroup.html
func (iamApi *iamApiServer) GetGroup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	groupName := r.FormValue("groupName")
	_, err := iamApi.authSys.Iam.GetGroup(ctx, groupName)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	var resp GetGroupResponse
	resp.GroupResult = GetGroupResult{
		G: Group{
			Path:      "",
			GroupName: groupName,
			GroupId:   "",
			Arn:       "",
		},
	}

	response.WriteSuccessResponseXML(w, r, resp)
}

// DeleteGroup
//Deletes the specified IAM group. The group must not contain any users or have any attached policies.
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeleteGroup.html
func (iamApi *iamApiServer) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	groupName := r.FormValue("groupName")
	err := iamApi.authSys.Iam.DeleteGroup(ctx, groupName)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}

	response.WriteSuccessResponseEmpty(w, r)
}

// ListGroups
//Lists the IAM groups that have the specified path prefix.
//You can paginate the results using the MaxItems and Marker parameters.
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListGroups.html
func (iamApi *iamApiServer) ListGroups(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	//_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	//if s3err != api_errors.ErrNone {
	//	response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
	//	return
	//}
	p := r.FormValue("pathPrefix")
	_, err := iamApi.authSys.Iam.ListGroups(ctx, p)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	var resp ListGroupsResponse
	response.WriteSuccessResponseXML(w, r, resp)
}
