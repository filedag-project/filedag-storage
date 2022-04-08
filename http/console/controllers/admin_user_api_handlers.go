package controllers

import (
	"encoding/json"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"io/ioutil"
	"net/http"
)

// ListUsers user list
func (control *Control) ListUsers(w http.ResponseWriter, r *http.Request) {
	var resp *models.ListUsersResponse
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	resp, error := control.apiServer.GetListUsersResponse(principal)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// AddUser user add
func (control *Control) AddUser(w http.ResponseWriter, r *http.Request) {
	var addUserRequest *models.AddUserRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	json.Unmarshal(body, &addUserRequest)
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	param := models.AddUserParams{
		Body: addUserRequest,
	}
	resp, error := control.apiServer.GetUserAddResponse(principal, param)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// RemoveUser user remove
func (control *Control) RemoveUser(w http.ResponseWriter, r *http.Request) {
	var removeUserParams *models.RemoveUserParams
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	json.Unmarshal(body, &removeUserParams)
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	error := control.apiServer.RemoveUserResponse(principal, *removeUserParams)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, "")
}

// UserInfo user info
func (control *Control) UserInfo(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	name := r.FormValue("name")
	userInfoParams := &models.GetUserInfoParams{
		Name: name,
	}
	resp, error := control.apiServer.GetUserInfoResponse(principal, *userInfoParams)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// SetUserPolicy set user policy
func (control *Control) SetUserPolicy(w http.ResponseWriter, r *http.Request) {
	var setUserPolicyParams *models.SetUserPolicyParams
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	json.Unmarshal(body, &setUserPolicyParams)
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	policy := `{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::test22/*"}]}`
	setUserPolicyParams.Definition = policy
	error := control.apiServer.GetUserSetPolicyResponse(principal, principal.AccountAccessKey, setUserPolicyParams)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, "")
}

// GetUserPolicy get user policy
func (control *Control) GetUserPolicy(w http.ResponseWriter, r *http.Request) {
	//name := r.FormValue("name")
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	resp, error := control.apiServer.GetUserPolicyResponse(principal, principal.AccountAccessKey)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// ListUserPolicy get user list policy
func (control *Control) ListUserPolicy(w http.ResponseWriter, r *http.Request) {
	//name := r.FormValue("name")
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	resp, error := control.apiServer.ListUserPolicyResponse(principal, principal.AccountAccessKey)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// RemoveUserPolicy remove user policy
func (control *Control) RemoveUserPolicy(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		PolicyName string `json:"policy_name"`
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	var params Params
	json.Unmarshal(body, &params)
	//name := r.FormValue("policy_name")
	token := r.Header.Get("Authorization")
	principal, err := VerifyToken(token)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	error := control.apiServer.RemoveUserPolicyResponse(principal, principal.AccountAccessKey, params.PolicyName)
	if error != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, "")
}
