package iamapi

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/api_errors"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/policy"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
)

const (
	policyDocumentVersion = "2012-10-17"
)

var policyLock = sync.RWMutex{}

//GetUserList get all user
func (iamApi *iamApiServer) GetUserList(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp ListUsersResponse
	resp.SetRequestId()
	resp.ListUsersResult.Users = iamApi.authSys.Iam.GetUserList(context.Background())
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// AddUser  add user
func (iamApi *iamApiServer) AddUser(w http.ResponseWriter, r *http.Request) {
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

//RemoveUser delete user
func (iamApi *iamApiServer) RemoveUser(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp CreateUserResponse
	accessKey := r.FormValue("accessKey")
	resp.CreateUserResult.User.UserName = &accessKey
	err := iamApi.authSys.Iam.RemoveUser(context.Background(), accessKey)
	if err != nil {
		response.WriteErrorResponse(w, r, api_errors.ErrInternalError)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//PutUserPolicy Put UserPolicy
func (iamApi *iamApiServer) PutUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(context.Background(), r, "", "", "")
	if s3err != api_errors.ErrNone {
		response.WriteErrorResponse(w, r, api_errors.ErrAccessDenied)
		return
	}
	var resp PutUserPolicyResponse
	userName := r.FormValue("userName")
	policyName := r.FormValue("policyName")
	policyDocumentString := r.FormValue("policyDocument")
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

//RemoveUserPolicy Remove UserPolicy
func (iamApi *iamApiServer) RemoveUserPolicy(w http.ResponseWriter, r *http.Request) {
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
func Hash(s *string) string {
	h := sha1.New()
	h.Write([]byte(*s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

//CreatePolicy Create Policy
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
	bucketMetas, err := iamApi.authSys.PolicySys.GetAllBucketOfUser(cred.AccessKey)
	if err != nil {
		response.WriteErrorResponseJSON(ctx, w, api_errors.GetAPIError(api_errors.ErrInternalError), r.URL, r.Host)
		return
	}
	var bucketAccessInfos []BucketAccessInfo
	for _, meta := range bucketMetas {
		bucketAccessInfos = append(bucketAccessInfos, BucketAccessInfo{
			Name:                 meta.Name,
			Size:                 0,
			Objects:              0,
			ObjectSizesHistogram: nil,
			Details:              nil,
			PrefixUsage:          nil,
			Created:              meta.Created,
			Access: AccountAccess{
				Read:  true,
				Write: true,
			},
		})
	}
	user := iam.UserInfo{
		SecretKey:  cred.SecretKey,
		PolicyName: "policy",
		Status:     iam.AccountStatus(cred.Status),
		MemberOf:   nil,
	}
	indent, err := json.MarshalIndent(user.PolicyName, "", " ")

	var accountInfo = AccountInfo{
		AccountName: userName,
		Policy:      indent,
		Buckets:     bucketAccessInfos,
	}

	data, err := json.Marshal(accountInfo)
	if err != nil {
		response.WriteErrorResponseJSON(ctx, w, api_errors.GetAPIError(api_errors.ErrJsonMarshal), r.URL, r.Host)
		return
	}
	response.WriteSuccessResponseJSON(w, data)
}
