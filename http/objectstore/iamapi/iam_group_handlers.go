package iamapi

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"net/http"
	"strconv"
)

// CreatGroup
//Creates a new group.
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateGroup.html
func (iamApi *iamApiServer) CreatGroup(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	groupName := r.FormValue("groupName")
	version := r.FormValue("version")
	atoi, _ := strconv.Atoi(version)
	err := iamApi.authSys.Iam.CreateGroup(r.Context(), groupName, atoi)
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
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	groupName := r.FormValue("groupName")
	_, err := iamApi.authSys.Iam.GetGroup(r.Context(), groupName)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	var resp GetGroupResponse
	resp.GroupResult = GroupResult{
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
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	groupName := r.FormValue("groupName")
	err := iamApi.authSys.Iam.DeleteGroup(r.Context(), groupName)
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
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	p := r.FormValue("pathPrefix")
	gi, err := iamApi.authSys.Iam.ListGroups(r.Context(), p)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	var resp ListGroupsResponse
	for _, g := range gi {
		resp.GroupResult.Groups = append(resp.GroupResult.Groups, GroupMember{GM: Group{
			Path:      p,
			GroupName: g.Name,
			GroupId:   strconv.Itoa(g.Version),
			Arn:       "",
		}})
	}

	response.WriteSuccessResponseXML(w, r, resp)
}
