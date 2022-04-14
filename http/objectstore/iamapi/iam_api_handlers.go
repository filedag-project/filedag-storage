package iamapi

import (
	"encoding/json"
	"errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

const (
	policyDocumentVersion = "2012-10-17"
	AccessKey             = "accessKey"
	SecretKey             = "secretKey"
	UserName              = "userName"
	PolicyName            = "policyName"
)

// CreateUser  add user
func (iamApi *iamApiServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	_, ok, _ := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp CreateUserResponse
	vars := mux.Vars(r)
	accessKey := vars[AccessKey]
	secretKey := vars[SecretKey]
	resp.CreateUserResult.User.UserName = &accessKey
	err := iamApi.authSys.Iam.AddUser(r.Context(), accessKey, secretKey)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//DeleteUser delete user
func (iamApi *iamApiServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp DeleteUserResponse
	accessKey := r.FormValue(AccessKey)
	err := iamApi.authSys.Iam.RemoveUser(r.Context(), accessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//GetUserInfo get user info
func (iamApi *iamApiServer) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	userName := r.FormValue(AccessKey)

	cred, ok := iamApi.authSys.Iam.GetUserInfo(r.Context(), userName)
	if !ok {
		response.WriteErrorResponseJSON(w, api_errors.GetAPIError(api_errors.ErrAccessKeyDisabled), r.URL, r.Host)
		return
	}
	polices, err := iamApi.authSys.Iam.GetUserPolices(r.Context(), userName)
	if err != nil {
		response.WriteErrorResponseJSON(w, api_errors.GetAPIError(api_errors.ErrInternalError), r.URL, r.Host)
		return
	}

	user := iam.UserInfo{
		SecretKey:  cred.SecretKey,
		PolicyName: polices,
		Status:     iam.AccountStatus(cred.Status),
	}

	data, err := json.Marshal(user)
	if err != nil {
		response.WriteErrorResponseJSON(w, api_errors.GetAPIError(api_errors.ErrJsonMarshal), r.URL, r.Host)
		return
	}
	response.WriteSuccessResponseJSON(w, data)
}

// ChangePassword change password
func (iamApi *iamApiServer) ChangePassword(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}

	secret := r.FormValue("newPassword")
	if !auth.IsSecretKeyValid(secret) {
		response.WriteErrorResponse(w, r, api_errors.ErrInvalidQueryParams)
	}
	c, ok := iamApi.authSys.Iam.GetUser(r.Context(), cred.AccessKey)
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessKeyDisabled)
		return
	}
	c.SecretKey = secret
	err := iamApi.authSys.Iam.UpdateUser(r.Context(), c)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}

// SetStatus set user status
func (iamApi *iamApiServer) SetStatus(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}

	user := r.FormValue(AccessKey)
	status := r.FormValue("status")
	c, ok := iamApi.authSys.Iam.GetUser(r.Context(), user)
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessKeyDisabled)
		return
	}
	c.Status = status
	err := iamApi.authSys.Iam.UpdateUser(r.Context(), c)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseEmpty(w, r)
}

func (iamApi *iamApiServer) AddSubUser(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp CreateUserResponse
	vars := mux.Vars(r)
	userName := vars["userName"]
	secretKey := vars["secretKey"]
	resp.CreateUserResult.User.UserName = &userName
	err := iamApi.authSys.Iam.AddSubUser(r.Context(), userName, secretKey, cred.AccessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

func (iamApi *iamApiServer) DeleteSubUser(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp CreateUserResponse
	vars := mux.Vars(r)
	userName := vars["userName"]
	c, ok := iamApi.authSys.Iam.GetUser(r.Context(), userName)
	if !ok {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	if c.ParentUser != cred.AccessKey {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	err := iamApi.authSys.Iam.RemoveUser(r.Context(), userName)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

func (iamApi *iamApiServer) GetSubUserInfo(w http.ResponseWriter, r *http.Request) {
	c, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	userName := r.FormValue("userName")
	cred, ok := iamApi.authSys.Iam.GetUserInfo(r.Context(), userName)
	if !ok {
		response.WriteErrorResponseJSON(w, api_errors.GetAPIError(api_errors.ErrAccessKeyDisabled), r.URL, r.Host)
		return
	}
	if c.AccessKey != cred.ParentUser {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	polices, err := iamApi.authSys.Iam.GetUserPolices(r.Context(), userName)
	if err != nil {
		response.WriteErrorResponseJSON(w, api_errors.GetAPIError(api_errors.ErrInternalError), r.URL, r.Host)
		return
	}

	user := iam.UserInfo{
		SecretKey:  cred.SecretKey,
		PolicyName: polices,
		Status:     iam.AccountStatus(cred.Status),
	}

	data, err := json.Marshal(user)
	if err != nil {
		response.WriteErrorResponseJSON(w, api_errors.GetAPIError(api_errors.ErrJsonMarshal), r.URL, r.Host)
		return
	}
	response.WriteSuccessResponseJSON(w, data)
}

//GetUserList get all user
func (iamApi *iamApiServer) GetUserList(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp ListUsersResponse
	users, err := iamApi.authSys.Iam.GetUserList(r.Context(), cred.AccessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	resp.ListUsersResult.Users = users
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//PutUserPolicy Put UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_PutUserPolicy.html
func (iamApi *iamApiServer) PutUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp PutUserPolicyResponse
	vars := mux.Vars(r)
	userName := vars[UserName]
	policyName := vars[PolicyName]
	policyDocumentString := vars["policyDocument"]
	policyDocument, err := GetPolicyDocument(&policyDocumentString)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	var pd policy.PolicyDocument
	err = iamApi.authSys.Iam.GetUserPolicy(r.Context(), userName, policyName, &pd)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchUserPolicy)
		return
	}
	policyMergeDocument := pd.Merge(policyDocument)
	if policyMergeDocument.Version == "" && policyMergeDocument.Statement == nil {
		log.Error(errors.New("The same user policy already exists "))
		response.WriteErrorResponse(w, r, api_errors.ErrUserPolicyAlreadyExists)
		return
	}
	err = iamApi.authSys.Iam.PutUserPolicy(r.Context(), userName, policyName, policyDocument)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//GetUserPolicy  Get UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetUserPolicy.html
func (iamApi *iamApiServer) GetUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp GetUserPolicyResponse
	userName := r.FormValue(UserName)
	policyName := r.FormValue(PolicyName)

	resp.GetUserPolicyResult.UserName = userName
	resp.GetUserPolicyResult.PolicyName = policyName
	policyDocument := policy.PolicyDocument{Version: policyDocumentVersion}
	err := iamApi.authSys.Iam.GetUserPolicy(r.Context(), userName, policyName, &policyDocument)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchUserPolicy)
		return
	}
	resp.GetUserPolicyResult.PolicyDocument = policyDocument.String()
	response.WriteXMLResponse(w, r, http.StatusOK, resp)

}

//ListUserPolicies  Get User all Policy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUserPolicies.html
func (iamApi *iamApiServer) ListUserPolicies(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp ListUserPoliciesResponse
	userName := r.FormValue(UserName)

	policyNames, err := iamApi.authSys.Iam.GetUserPolices(r.Context(), userName)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchUserPolicy)
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
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp DeleteUserPolicyResponse
	userName := r.FormValue(UserName)
	policyName := r.FormValue(PolicyName)
	err := iamApi.authSys.Iam.RemoveUserPolicy(r.Context(), userName, policyName)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchUserPolicy)
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
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
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
	err = iamApi.authSys.Iam.CreatePolicy(r.Context(), policyName, policyDocument)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrNoSuchBucketPolicy)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}*/
