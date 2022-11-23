package iamapi

import (
	"encoding/json"
	"errors"
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iam/auth"
	"github.com/filedag-project/filedag-storage/objectservice/iam/policy"
	"github.com/filedag-project/filedag-storage/objectservice/iam/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"strconv"
)

const (
	policyDocumentVersion = "2012-10-17"
	AccessKey             = "accessKey"
	Capacity              = "capacity"
	SecretKey             = "secretKey"
	NewSecretKey          = "newSecretKey"
	UserName              = "userName"
	PolicyName            = "policyName"
	AccountStatus         = "status"
)

var validAccessKey = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9\.\-]{1,18}[A-Za-z0-9]$`)

// CreateUser  add user
func (iamApi *iamApiServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	_, ok, _ := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if !ok {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	var resp CreateUserResponse
	vars := mux.Vars(r)
	accessKey := vars[AccessKey]
	secretKey := vars[SecretKey]
	capacity := vars[Capacity]
	capa, err := strconv.ParseUint(capacity, 10, 64)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidRequestParameter)
	}
	if !auth.IsAccessKeyValid(accessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidFormatAccessKey)
		return
	}
	if !validAccessKey.MatchString(accessKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidFormatAccessKey)
		return
	}
	if !auth.IsSecretKeyValid(secretKey) {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidQueryParams)
	}
	resp.CreateUserResult.User.UserName = &accessKey
	_, err = iamApi.authSys.Iam.GetUserInfo(r.Context(), accessKey)
	if err == nil {
		response.WriteErrorResponseJSON(w, apierrors.GetAPIError(apierrors.ErrUserAlreadyExists), r.URL, r.Host)
		return
	}
	err = iamApi.authSys.Iam.AddUser(r.Context(), accessKey, secretKey, capa)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//DeleteUser delete user
func (iamApi *iamApiServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.RemoveUserAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	var resp DeleteUserResponse
	accessKey := r.FormValue(AccessKey)
	_, err := iamApi.authSys.Iam.GetUserInfo(r.Context(), accessKey)
	if err != nil {
		response.WriteErrorResponseJSON(w, apierrors.GetAPIError(apierrors.ErrNoSuchUser), r.URL, r.Host)
		return
	}
	err = iamApi.authSys.Iam.RemoveUser(r.Context(), accessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	// clean removed user's bucket
	// TODO: If the deletion fails, try again
	go func() {
		iamApi.cleanData(accessKey)
	}()

	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// AccountInfoHandler returns usage
func (iamApi *iamApiServer) AccountInfoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetUserInfoAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	accessKey := r.FormValue(AccessKey)
	if cred.AccessKey != iamApi.authSys.AdminCred.AccessKey && cred.AccessKey != accessKey {
		response.WriteErrorResponseJSON(w, apierrors.GetAPIError(apierrors.ErrAccessDenied), r.URL, r.Host)
		return
	}
	var err error
	bucketInfos := iamApi.bucketInfoFunc(ctx, cred.AccessKey)

	accountName := accessKey
	if cred.IsTemp() {
		// For derived credentials, check the parent user's permissions.
		accountName = cred.ParentUser
	}
	polices, err := iamApi.authSys.Iam.GetUserPolices(r.Context(), cred.AccessKey)
	if err != nil {
		response.WriteErrorResponseJSON(w, apierrors.GetAPIError(apierrors.ErrInternalError), r.URL, r.Host)
		return
	}
	var info = iam.UserIdentity{Credentials: iamApi.authSys.AdminCred, TotalStorageCapacity: 999999999}
	if cred.AccessKey != iamApi.authSys.AdminCred.AccessKey {
		info, err = iamApi.authSys.Iam.GetUserInfo(ctx, cred.AccessKey)
		if err != nil {
			response.WriteErrorResponseJSON(w, apierrors.GetAPIError(apierrors.ErrNoSuchUser), r.URL, r.Host)
			return
		}
	}

	var useStorageCapacity uint64
	for _, bi := range bucketInfos {
		useStorageCapacity += bi.Size
	}
	acctInfo := iam.UserInfo{
		AccountName:          accountName,
		TotalStorageCapacity: info.TotalStorageCapacity,
		UseStorageCapacity:   useStorageCapacity,
		PolicyName:           polices,
		BucketInfos:          bucketInfos,
		Status: func() iam.AccountStatus {
			if cred.IsValid() {
				return iam.AccountEnabled
			}
			return iam.AccountDisabled
		}(),
	}

	usageInfoJSON, err := json.Marshal(acctInfo)
	if err != nil {
		response.WriteErrorResponseJSON(w, apierrors.GetAPIError(apierrors.ErrMalformedJSON), r.URL, r.Host)
		return
	}

	response.WriteSuccessResponseJSON(w, usageInfoJSON)
}

// ChangePassword change password
func (iamApi *iamApiServer) ChangePassword(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}

	secret := r.FormValue(NewSecretKey)
	userName := r.FormValue(AccessKey)
	if !auth.IsSecretKeyValid(secret) {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidQueryParams)
		return
	}
	c, ok := iamApi.authSys.Iam.GetUser(r.Context(), userName)
	if !ok {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessKeyDisabled)
		return
	}
	c.SecretKey = secret
	err := iamApi.authSys.Iam.UpdateUser(r.Context(), c)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseHeadersOnly(w, r)
}

// SetStatus set user status
func (iamApi *iamApiServer) SetStatus(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "Set-Status", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}

	user := r.FormValue(AccessKey)
	status := r.FormValue(AccountStatus)
	switch status {
	case auth.AccountOn, auth.AccountOff:
	default:
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidQueryParams)
		return
	}
	c, _ := iamApi.authSys.Iam.GetUser(r.Context(), user)
	if c.AccessKey == "" {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessKeyDisabled)
		return
	}
	c.Status = status
	err := iamApi.authSys.Iam.UpdateUser(r.Context(), c)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	response.WriteSuccessResponseHeadersOnly(w, r)
}

func (iamApi *iamApiServer) AddSubUser(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	var resp CreateUserResponse
	vars := mux.Vars(r)
	userName := vars["userName"]
	secretKey := vars["secretKey"]
	capacity := vars["capacity"]
	capa, err := strconv.ParseUint(capacity, 10, 64)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInvalidRequestParameter)
	}
	resp.CreateUserResult.User.UserName = &userName
	err = iamApi.authSys.Iam.AddSubUser(r.Context(), userName, secretKey, cred.AccessKey, capa)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

func (iamApi *iamApiServer) DeleteSubUser(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	var resp CreateUserResponse
	vars := mux.Vars(r)
	userName := vars["userName"]
	c, ok := iamApi.authSys.Iam.GetUser(r.Context(), userName)
	if !ok {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	if c.ParentUser != cred.AccessKey {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	err := iamApi.authSys.Iam.RemoveUser(r.Context(), userName)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

func (iamApi *iamApiServer) GetSubUserInfo(w http.ResponseWriter, r *http.Request) {
	// todo implement SubUserInfo
	c, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	userName := r.FormValue("userName")
	info, err := iamApi.authSys.Iam.GetUserInfo(r.Context(), userName)
	if err != nil {
		response.WriteErrorResponseJSON(w, apierrors.GetAPIError(apierrors.ErrNoSuchUser), r.URL, r.Host)
		return
	}
	if c.AccessKey != info.Credentials.ParentUser {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	polices, err := iamApi.authSys.Iam.GetUserPolices(r.Context(), userName)
	if err != nil {
		response.WriteErrorResponseJSON(w, apierrors.GetAPIError(apierrors.ErrInternalError), r.URL, r.Host)
		return
	}

	user := iam.UserInfo{
		//SecretKey:  cred.SecretKey,
		PolicyName: polices,
		Status: func() iam.AccountStatus {
			if info.Credentials.IsValid() {
				return iam.AccountEnabled
			}
			return iam.AccountDisabled
		}(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		response.WriteErrorResponseJSON(w, apierrors.GetAPIError(apierrors.ErrMalformedJSON), r.URL, r.Host)
		return
	}
	response.WriteSuccessResponseJSON(w, data)
}

//GetUserList get all user
func (iamApi *iamApiServer) GetUserList(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	var resp ListUsersResponse
	users, err := iamApi.authSys.Iam.GetUserList(r.Context(), cred.AccessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	resp.ListUsersResult.Users = users
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//PutUserPolicy Put UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_PutUserPolicy.html
func (iamApi *iamApiServer) PutUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	var resp PutUserPolicyResponse
	vars := mux.Vars(r)
	userName := vars[UserName]
	policyName := vars[PolicyName]
	policyDocumentString := vars["policyDocument"]
	policyDocument, err := GetPolicyDocument(&policyDocumentString)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	var pd policy.PolicyDocument
	_ = iamApi.authSys.Iam.GetUserPolicy(r.Context(), userName, policyName, &pd)
	//if err != nil {
	//	response.WriteErrorResponse(w, r, apierrors.ErrNoSuchUserPolicy)
	//	return
	//}
	policyMergeDocument := pd.Merge(policyDocument)
	if policyMergeDocument.Version == "" && policyMergeDocument.Statement == nil {
		log.Error(errors.New("The same user policy already exists "))
		response.WriteErrorResponse(w, r, apierrors.ErrUserPolicyAlreadyExists)
		return
	}
	err = iamApi.authSys.Iam.PutUserPolicy(r.Context(), userName, policyName, policyDocument)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//GetUserPolicy  Get UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetUserPolicy.html
func (iamApi *iamApiServer) GetUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
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
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchUserPolicy)
		return
	}
	resp.GetUserPolicyResult.PolicyDocument = policyDocument.String()
	response.WriteXMLResponse(w, r, http.StatusOK, resp)

}

//ListUserPolicies  Get User all Policy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUserPolicies.html
func (iamApi *iamApiServer) ListUserPolicies(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	var resp ListUserPoliciesResponse
	userName := r.FormValue(UserName)

	policyNames, err := iamApi.authSys.Iam.GetUserPolices(r.Context(), userName)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchUserPolicy)
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
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	var resp DeleteUserPolicyResponse
	userName := r.FormValue(UserName)
	policyName := r.FormValue(PolicyName)
	err := iamApi.authSys.Iam.RemoveUserPolicy(r.Context(), userName, policyName)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchUserPolicy)
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
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponse(w, r, apierrors.ErrAccessDenied)
		return
	}
	var resp CreatePolicyResponse
	policyName := r.FormValue("policyName")
	policyDocumentString := r.FormValue("policyDocument")
	policyDocument, err := GetPolicyDocument(&policyDocumentString)
	if err != nil {
		response.WriteErrorResponse(w, r, apierrors.ErrInternalError)
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
		response.WriteErrorResponse(w, r, apierrors.ErrNoSuchBucketPolicy)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}*/
