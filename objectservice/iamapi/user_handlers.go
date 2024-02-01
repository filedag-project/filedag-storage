package iamapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/filedag-project/filedag-storage/objectservice/apierrors"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/policy"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/s3action"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
)

const (
	defaultPolicyDocumentVersion = "2012-10-17"
	accessKey                    = "accessKey"
	capacity                     = "capacity"
	secretKey                    = "secretKey"
	newSecretKey                 = "newSecretKey"
	oldSecretKey                 = "oldSecretKey"
	userName                     = "userName"
	policyName                   = "policyName"
	accountStatus                = "status"
)

var validAccessKey = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9\.\-]{1,18}[A-Za-z0-9]$`)

// CreateUser  add user
func (iamApi *iamApiServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	_, ok, _ := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if !ok {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	username := r.FormValue(accessKey)
	userSecret := r.FormValue(secretKey)
	userCapacity := r.FormValue(capacity)
	capa, err := strconv.ParseUint(userCapacity, 10, 64)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidRequestParameter))
		return
	}
	if capa > 1<<50*10 {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidRequestParameter))
		return
	}
	if !auth.IsAccessKeyValid(username) {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidFormatAccessKey))
		return
	}
	if !validAccessKey.MatchString(username) {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidFormatAccessKey))
		return
	}
	if !auth.IsSecretKeyValid(userSecret) {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidRequestParameter))
		return
	}
	_, err = iamApi.authSys.Iam.GetUserInfo(r.Context(), username)
	if err == nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrUserAlreadyExists))
		return
	}
	err = iamApi.authSys.Iam.AddUser(r.Context(), username, userSecret, capa)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	response.WriteSuccessResponseJSON(w, r, nil)
}

//DeleteUser delete user
func (iamApi *iamApiServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.RemoveUserAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	username := r.FormValue(accessKey)
	_, err := iamApi.authSys.Iam.GetUserInfo(r.Context(), username)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrNoSuchUser))
		return
	}
	if username != cred.AccessKey && cred.AccessKey != iamApi.authSys.AdminCred.AccessKey {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	err = iamApi.authSys.Iam.RemoveUser(r.Context(), username)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	// clean removed user's bucket
	// TODO: If the deletion fails, try again
	go func() {
		iamApi.cleanData(username)
	}()
	response.WriteSuccessResponseJSON(w, r, nil)
}

// AccountInfo returns usage
func (iamApi *iamApiServer) AccountInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetUserInfoAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(s3err))
		return
	}
	accountName := r.FormValue(accessKey)
	if cred.AccessKey != accountName {
		if accountName == iamApi.authSys.AdminCred.AccessKey {
			response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
			return
		}
		ucred, ok := iamApi.authSys.Iam.GetUser(ctx, accountName)
		if ok {
			// user exist:
			//1) tmp user
			//2) other user
			if ucred.IsTemp() {
				// For derived credentials, check the parent user's permissions.
				accountName = ucred.ParentUser
			} else {
				// only root user can get other user info
				if cred != iamApi.authSys.AdminCred {
					response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
					return
				}
			}
		} else {
			response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrNoSuchUser))
			return
		}
	}
	var err error
	bucketInfos := iamApi.bucketInfoFunc(ctx, accountName)
	var info = iam.UserIdentity{Credentials: iamApi.authSys.AdminCred, TotalStorageCapacity: math.MaxUint64}
	if accountName != iamApi.authSys.AdminCred.AccessKey {
		info, err = iamApi.authSys.Iam.GetUserInfo(ctx, accountName)
		if err != nil {
			response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrNoSuchUser))
			return
		}
	}

	var useStorageCapacity, bucketsCount, objectsCount uint64

	for _, bi := range bucketInfos {
		useStorageCapacity += bi.Size
		bucketsCount++
		objectsCount += bi.Objects
	}
	acctInfo := iam.UserOverView{
		AccountName:          accountName,
		TotalStorageCapacity: info.TotalStorageCapacity,
		UseStorageCapacity:   useStorageCapacity,
		BucketsCount:         bucketsCount,
		ObjectsCount:         objectsCount,
	}
	response.WriteSuccessResponseJSON(w, r, acctInfo)
}

// GetPolicyName get PolicyName
func (iamApi *iamApiServer) GetPolicyName(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetUserInfoAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(s3err))
		return
	}
	bucket := r.FormValue("bucketName")
	bucketPolicyBytes, err := ioutil.ReadAll(io.LimitReader(r.Body, r.ContentLength))
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidPolicyDocument))
		return
	}
	bucketPolicy, err := policy.ParseConfig(bytes.NewReader(bucketPolicyBytes), bucket)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrMalformedPolicy))
		return
	}
	bucketPolicyName := store.GetPolicyName(bucketPolicy.Statements, bucket, "")
	response.WriteSuccessResponseJSON(w, r, bucketPolicyName)
}

//IsAdmin check user if admin
func (iamApi *iamApiServer) IsAdmin(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, s3action.GetUserInfoAction, "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(s3err))
		return
	}
	isAdmin := false
	if cred.AccessKey == iamApi.authSys.AdminCred.AccessKey {
		isAdmin = true
	}
	response.WriteSuccessResponseJSON(w, r, isAdmin)
}

// todo review SubUser

func (iamApi *iamApiServer) AddSubUser(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	var resp CreateUserResponse
	vars := mux.Vars(r)
	userName := vars["userName"]
	secretKey := vars["secretKey"]
	capacity := vars["capacity"]
	capa, err := strconv.ParseUint(capacity, 10, 64)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInvalidRequestParameter))
	}
	resp.CreateUserResult.User.UserName = &userName
	err = iamApi.authSys.Iam.AddSubUser(r.Context(), userName, secretKey, cred.AccessKey, capa)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	response.WriteSuccessResponseJSON(w, r, nil)
}

func (iamApi *iamApiServer) DeleteSubUser(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	vars := mux.Vars(r)
	username := vars["userName"]
	c, ok := iamApi.authSys.Iam.GetUser(r.Context(), username)
	if !ok {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	if c.ParentUser != cred.AccessKey {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	err := iamApi.authSys.Iam.RemoveUser(r.Context(), username)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	response.WriteSuccessResponseJSON(w, r, nil)
}

func (iamApi *iamApiServer) GetSubUserInfo(w http.ResponseWriter, r *http.Request) {
	// todo implement SubUserInfo
	c, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	username := r.FormValue("userName")
	info, err := iamApi.authSys.Iam.GetUserInfo(r.Context(), username)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrNoSuchUser))
		return
	}
	if c.AccessKey != info.Credentials.ParentUser {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	polices, err := iamApi.authSys.Iam.GetUserPolices(r.Context(), username)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
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
	response.WriteSuccessResponseJSON(w, r, user)
}

//GetUserList get all user
func (iamApi *iamApiServer) GetUserList(w http.ResponseWriter, r *http.Request) {
	cred, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	var resp ListUsersResponse
	users, err := iamApi.authSys.Iam.GetUserList(r.Context(), cred.AccessKey)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	resp.ListUsersResult.Users = users
	response.WriteSuccessResponseJSON(w, r, resp)
}

//PutUserPolicy Put UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_PutUserPolicy.html
func (iamApi *iamApiServer) PutUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	vars := mux.Vars(r)
	username := vars[userName]
	policyName := vars[policyName]
	policyDocumentString := vars["policyDocument"]
	policyDocument, err := GetPolicyDocument(&policyDocumentString)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	var pd policy.PolicyDocument
	_ = iamApi.authSys.Iam.GetUserPolicy(r.Context(), username, policyName, &pd)
	//if err != nil {
	//	response.WriteErrorResponseJSON(w,r, apierrors.GetAPIError(apierrors.ErrNoSuchUserPolicy)
	//	return
	//}
	policyMergeDocument := pd.Merge(policyDocument)
	if policyMergeDocument.Version == "" && policyMergeDocument.Statement == nil {
		log.Error(errors.New("The same user policy already exists "))
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrUserPolicyAlreadyExists))
		return
	}
	err = iamApi.authSys.Iam.PutUserPolicy(r.Context(), username, policyName, policyDocument)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrInternalError))
		return
	}
	response.WriteSuccessResponseJSON(w, r, nil)
}

//GetUserPolicy  Get UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetUserPolicy.html
func (iamApi *iamApiServer) GetUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	var resp GetUserPolicyResponse
	username := r.FormValue(userName)
	policiesName := r.FormValue(policyName)

	resp.GetUserPolicyResult.UserName = username
	resp.GetUserPolicyResult.PolicyName = policiesName
	policyDocument := policy.PolicyDocument{Version: defaultPolicyDocumentVersion}
	err := iamApi.authSys.Iam.GetUserPolicy(r.Context(), username, policiesName, &policyDocument)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrNoSuchUserPolicy))
		return
	}
	resp.GetUserPolicyResult.PolicyDocument = policyDocument.String()
	response.WriteSuccessResponseJSON(w, r, resp)

}

//ListUserPolicies  Get User all Policy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_ListUserPolicies.html
func (iamApi *iamApiServer) ListUserPolicies(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	var resp ListUserPoliciesResponse
	username := r.FormValue(userName)

	policyNames, err := iamApi.authSys.Iam.GetUserPolices(r.Context(), username)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrNoSuchUserPolicy))
		return
	}
	var members []string
	for _, v := range policyNames {
		members = append(members, v)
	}
	resp.ListUserPoliciesResult.PolicyNames.Member = members
	response.WriteSuccessResponseJSON(w, r, resp)

}

//DeleteUserPolicy Remove UserPolicy
//https://docs.aws.amazon.com/IAM/latest/APIReference/API_DeleteUserPolicy.html
func (iamApi *iamApiServer) DeleteUserPolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrAccessDenied))
		return
	}
	var resp DeleteUserPolicyResponse
	username := r.FormValue(userName)
	policiesName := r.FormValue(policyName)
	err := iamApi.authSys.Iam.RemoveUserPolicy(r.Context(), username, policiesName)
	if err != nil {
		response.WriteErrorResponseJSON(w, r, apierrors.GetAPIError(apierrors.ErrNoSuchUserPolicy))
		return
	}
	response.WriteSuccessResponseJSON(w, r, resp)
}

//GetPolicyDocument Get PolicyDocument
func GetPolicyDocument(policyD *string) (policyDocument policy.PolicyDocument, err error) {
	if err = json.Unmarshal([]byte(*policyD), &policyDocument); err != nil {
		return policy.PolicyDocument{}, err
	}
	return policyDocument, err
}
func (iamApi *iamApiServer) getAllUserInfos(ctx context.Context) ([]iam.UserInfo, error) {
	userIdentities, err := iamApi.authSys.Iam.GetAllUser(ctx)
	if err != nil {
		return nil, err
	}
	var allUserInfo []iam.UserInfo
	for _, userIdentity := range userIdentities {
		bucketInfos := iamApi.bucketInfoFunc(ctx, userIdentity.Credentials.AccessKey)

		polices, err := iamApi.authSys.Iam.GetUserPolices(ctx, userIdentity.Credentials.AccessKey)
		if err != nil {
			return nil, err
		}
		var useStorageCapacity uint64
		for _, bi := range bucketInfos {
			useStorageCapacity += bi.Size
		}
		acctInfo := iam.UserInfo{
			AccountName:          userIdentity.Credentials.AccessKey,
			TotalStorageCapacity: userIdentity.TotalStorageCapacity,
			UseStorageCapacity:   useStorageCapacity,
			PolicyName:           polices,
			BucketInfos:          bucketInfos,
			Status: func() iam.AccountStatus {
				if userIdentity.Credentials.IsValid() {
					return iam.AccountEnabled
				}
				return iam.AccountDisabled
			}(),
		}
		allUserInfo = append(allUserInfo, acctInfo)
	}
	return allUserInfo, nil
}

/*//CreatePolicy Create Policy
func (iamApi *iamApiServer) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	_, _, s3err := iamApi.authSys.CheckRequestAuthTypeCredential(r.Context(), r, "", "", "")
	if s3err != apierrors.ErrNone {
		response.WriteErrorResponseJSON(w,r, apierrors.GetAPIError(apierrors.ErrAccessDenied)
		return
	}
	var resp CreatePolicyResponse
	policyName := r.FormValue("policyName")
	policyDocumentString := r.FormValue("policyDocument")
	policyDocument, err := GetPolicyDocument(&policyDocumentString)
	if err != nil {
		response.WriteErrorResponseJSON(w,r, apierrors.GetAPIError(apierrors.ErrInternalError)
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
		response.WriteErrorResponseJSON(w,r, apierrors.GetAPIError(apierrors.ErrNoSuchBucketPolicy)
		return
	}
	response.WriteXMLResponse(w, r, http.StatusOK, resp)
}*/
