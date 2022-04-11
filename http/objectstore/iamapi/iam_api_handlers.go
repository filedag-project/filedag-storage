package iamapi

import (
	"context"
	"encoding/json"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	policyDocumentVersion = "2012-10-17"
)

// CreateUser  add user
func (iamApi *iamApiServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	_, ok, _ := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp CreateUserResponse
	vars := mux.Vars(r)
	accessKey := vars["accessKey"]
	secretKey := vars["secretKey"]
	resp.CreateUserResult.User.UserName = &accessKey
	err := iamApi.authSys.Iam.AddUser(context.Background(), accessKey, secretKey)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//DeleteUser delete user
func (iamApi *iamApiServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp DeleteUserResponse
	accessKey := r.FormValue("accessKey")
	err := iamApi.authSys.Iam.RemoveUser(context.Background(), accessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//GetUserInfo get user info
func (iamApi *iamApiServer) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	userName := r.FormValue("userName")
	ctx := context.Background()
	cred, ok := iamApi.authSys.Iam.GetUserInfo(ctx, userName)
	if !ok {
		response.WriteErrorResponseJSON(ctx, w, api_errors.GetAPIError(api_errors.ErrAccessKeyDisabled), r.URL, r.Host)
		return
	}
	polices, err := iamApi.authSys.Iam.GetUserPolices(ctx, userName)
	if err != nil {
		response.WriteErrorResponseJSON(ctx, w, api_errors.GetAPIError(api_errors.ErrInternalError), r.URL, r.Host)
		return
	}

	user := iam.UserInfo{
		SecretKey:  cred.SecretKey,
		PolicyName: polices,
		Status:     iam.AccountStatus(cred.Status),
	}

	data, err := json.Marshal(user)
	if err != nil {
		response.WriteErrorResponseJSON(ctx, w, api_errors.GetAPIError(api_errors.ErrJsonMarshal), r.URL, r.Host)
		return
	}
	response.WriteSuccessResponseJSON(w, data)
}

// ChangePassword change password
func (iamApi *iamApiServer) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}

	secret := r.FormValue("newPassword")
	if !auth.IsSecretKeyValid(secret) {
		response.WriteErrorResponse(w, r, api_errors.ErrInvalidQueryParams)
	}
	c, ok := iamApi.authSys.Iam.GetUser(ctx, cred.AccessKey)
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessKeyDisabled)
		return
	}
	c.SecretKey = secret
	err := iamApi.authSys.Iam.UpdateUser(ctx, c)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}

// SetStatus set user status
func (iamApi *iamApiServer) SetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(ctx, r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}

	user := r.FormValue("userName")
	status := r.FormValue("status")
	c, ok := iamApi.authSys.Iam.GetUser(ctx, user)
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessKeyDisabled)
		return
	}
	c.Status = status
	err := iamApi.authSys.Iam.UpdateUser(ctx, c)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}

func (iamApi *iamApiServer) AddSubUser(w http.ResponseWriter, r *http.Request) {
	cred, ok, _ := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp CreateUserResponse
	vars := mux.Vars(r)
	userName := vars["userName"]
	secretKey := vars["secretKey"]
	resp.CreateUserResult.User.UserName = &userName
	err := iamApi.authSys.Iam.AddSubUser(context.Background(), userName, secretKey, cred.AccessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

func (iamApi *iamApiServer) DeleteSubUser(w http.ResponseWriter, r *http.Request) {

}

func (iamApi *iamApiServer) GetSubUserInfo(w http.ResponseWriter, r *http.Request) {

}

//GetUserList get all user
func (iamApi *iamApiServer) GetUserList(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp ListUsersResponse
	resp.ListUsersResult.Users = iamApi.authSys.Iam.GetUserList(context.Background())
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//PutUserPolicy Put UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_PutUserPolicy.html
func (iamApi *iamApiServer) PutUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp PutUserPolicyResponse
	vars := mux.Vars(r)
	userName := vars["userName"]
	policyName := vars["policyName"]
	policyDocumentString := vars["policyDocument"]
	policyDocument, err := GetPolicyDocument(&policyDocumentString)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	err = iamApi.authSys.Iam.PutUserPolicy(context.Background(), userName, policyName, policyDocument)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucketPolicy)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//GetUserPolicy  Get UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetUserPolicy.html
func (iamApi *iamApiServer) GetUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp GetUserPolicyResponse
	userName := r.FormValue("userName")
	policyName := r.FormValue("policyName")

	resp.GetUserPolicyResult.UserName = userName
	resp.GetUserPolicyResult.PolicyName = policyName
	policyDocument := policy.PolicyDocument{Version: policyDocumentVersion}
	err := iamApi.authSys.Iam.GetUserPolicy(context.Background(), userName, policyName, &policyDocument)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucketPolicy)
		return
	}
	resp.GetUserPolicyResult.PolicyDocument = policyDocument.String()
	response.WriteXMLResponse(w, r, http.StatusOK, resp)

}

//ListUserPolicies  Get User all Policy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUserPolicies.html
func (iamApi *iamApiServer) ListUserPolicies(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp ListUserPoliciesResponse
	userName := r.FormValue("userName")

	policyNames, err := iamApi.authSys.Iam.GetUserPolices(context.Background(), userName)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucketPolicy)
		return
	}
	var members []string
	for _, v := range policyNames {
		members = append(members, v)
	}
	resp.ListUserPoliciesResult.PolicyNames.Member = members
	response.WriteXMLResponse(w, r, http.StatusOK, resp)

}

//DeleteUserPolicy Remove UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeleteUserPolicy.html
func (iamApi *iamApiServer) DeleteUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp PutUserPolicyResponse
	userName := r.FormValue("userName")
	policyName := r.FormValue("policyName")
	err := iamApi.authSys.Iam.RemoveUserPolicy(context.Background(), userName, policyName)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucketPolicy)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//GetPolicyDocument Get PolicyDocument
func GetPolicyDocument(policyD *string) (policyDocument policy.PolicyDocument, err error) {
	if err = json.Unmarshal([]byte(*policyD), &policyDocument); err != nil {
		return policy.PolicyDocument{}, err
	}
	return policyDocument, err
}

/*//CreatePolicy Create Policy
func (iamApi *iamApiServer) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp CreatePolicyResponse
	policyName := r.FormValue("policyName")
	policyDocumentString := r.FormValue("policyDocument")
	policyDocument, err := GetPolicyDocument(&policyDocumentString)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
	}
	policyId := Hash(&policyDocumentString)
	arn := fmt.Sprintf("arn:aws:iam:::policy/%s", policyName)
	resp.CreatePolicyResult.Policy.PolicyName = &policyName
	resp.CreatePolicyResult.Policy.Arn = &arn
	resp.CreatePolicyResult.Policy.PolicyId = &policyId
	policyLock.Lock()
	defer policyLock.Unlock()
	err = iamApi.authSys.Iam.CreatePolicy(context.Background(), policyName, policyDocument)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucketPolicy)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}*/
