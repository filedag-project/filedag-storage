package madmin

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/filedag-project/filedag-storage/http/console/madmin/tags"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// ListUsers - list all users.
func (adm *AdminClient) ListUsers(ctx context.Context) (users []*iam.User, err error) {
	reqData := requestData{
		relPath: admin + adminAPIPrefix + AdminAPIVersionV1 + "/list-user",
	}
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)
	defer closeResponse(resp)
	if err != nil {
		return users, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp)
	fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		return users, httpRespToErrorResponse(resp)
	}
	var response ListUsersResponse
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return users, err
	}
	fmt.Println(response)
	users = response.ListUsersResult.Users
	return users, nil
}

// AddUser - adds a user.
func (adm *AdminClient) AddUser(ctx context.Context, accessKey, secretKey string) (*iam.User, error) {
	return adm.SetUser(ctx, accessKey, secretKey, AccountEnabled)
}

// SetUser - update user secret key or account status.
func (adm *AdminClient) SetUser(ctx context.Context, accessKey, secretKey string, status AccountStatus) (user *iam.User, err error) {
	queryValues := url.Values{}
	queryValues.Set("accessKey", accessKey)
	queryValues.Set("secretKey", secretKey)
	reqData := requestData{
		relPath:     admin + adminAPIPrefix + AdminAPIVersionV1 + "/add-user",
		queryValues: queryValues,
	}
	resp, err := adm.executeMethod(ctx, http.MethodPost, reqData)
	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(resp)
	fmt.Println(string(body))
	var response CreateUserResponse
	err = xml.Unmarshal(body, &response)
	user = &response.CreateUserResult.User
	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp)
	}
	return user, nil
}

// RemoveUser - remove a user.
func (adm *AdminClient) RemoveUser(ctx context.Context, accessKey string) error {
	queryValues := url.Values{}
	queryValues.Set("accessKey", accessKey)
	reqData := requestData{
		relPath:     admin + adminAPIPrefix + AdminAPIVersionV1 + "/remove-user",
		queryValues: queryValues,
	}
	resp, err := adm.executeMethod(ctx, http.MethodPost, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}

// AccountAccess contains information about
type AccountAccess struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
}

// BucketDetails provides information about features currently
// turned-on per bucket.
type BucketDetails struct {
	Versioning          bool         `json:"versioning"`
	VersioningSuspended bool         `json:"versioningSuspended"`
	Locking             bool         `json:"locking"`
	Replication         bool         `json:"replication"`
	Tagging             *tags.Tags   `json:"tags"`
	Quota               *BucketQuota `json:"quota"`
}

// BucketAccessInfo represents bucket usage of a bucket, and its relevant
// access type for an account
type BucketAccessInfo struct {
	Name                 string            `json:"name"`
	Size                 uint64            `json:"size"`
	Objects              uint64            `json:"objects"`
	ObjectSizesHistogram map[string]uint64 `json:"objectHistogram"`
	Details              *BucketDetails    `json:"details"`
	PrefixUsage          map[string]uint64 `json:"prefixUsage"`
	Created              time.Time         `json:"created"`
	Access               AccountAccess     `json:"access"`
}

// AccountInfo represents the account usage info of an
// account across buckets.
type AccountInfo struct {
	AccountName string
	Server      BackendInfo
	Policy      json.RawMessage // Use iam/policy.Parse to parse the result, to be done by the caller.
	Buckets     []BucketAccessInfo
}

// AccountOpts allows for configurable behavior with "prefix-usage"
type AccountOpts struct {
	PrefixUsage bool
}

// AccountStatus - account status.
type AccountStatus string

// Account status per user.
const (
	AccountEnabled  AccountStatus = "enabled"
	AccountDisabled AccountStatus = "disabled"
)

// UserInfo carries information about long term users.
type UserInfo struct {
	SecretKey  string        `json:"secretKey,omitempty"`
	PolicyName []string      `json:"policyName,omitempty"`
	Status     AccountStatus `json:"status"`
}

// GetUserInfo - get info on a user
func (adm *AdminClient) GetUserInfo(ctx context.Context, name string) (*UserInfo, error) {
	var userInfo UserInfo
	queryValues := url.Values{}
	queryValues.Set("userName", name)
	reqData := requestData{
		relPath:     admin + adminAPIPrefix + AdminAPIVersionV1 + "/user-info",
		queryValues: queryValues,
	}
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)

	defer closeResponse(resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, httpRespToErrorResponse(resp)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &userInfo, err
	}
	fmt.Println(resp.Body)

	if err = json.Unmarshal(b, &userInfo); err != nil {
		return &userInfo, err
	}

	return &userInfo, nil
}

// AddOrUpdateUserReq allows to update user details such as secret key and
// account status.
type AddOrUpdateUserReq struct {
	SecretKey string        `json:"secretKey,omitempty"`
	Status    AccountStatus `json:"status"`
}

// SetUserStatus - adds a status for a user.
func (adm *AdminClient) SetUserStatus(ctx context.Context, accessKey string, status AccountStatus) error {
	queryValues := url.Values{}
	queryValues.Set("accessKey", accessKey)
	queryValues.Set("status", string(status))

	reqData := requestData{
		relPath:     adminAPIPrefix + "/set-user-status",
		queryValues: queryValues,
	}

	// Execute PUT on /minio/admin/v3/set-user-status to set status.
	resp, err := adm.executeMethod(ctx, http.MethodPut, reqData)

	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpRespToErrorResponse(resp)
	}

	return nil
}

// AddServiceAccountReq is the request options of the add service account admin call
type AddServiceAccountReq struct {
	Policy     json.RawMessage `json:"policy,omitempty"` // Parsed value from iam/policy.Parse()
	TargetUser string          `json:"targetUser,omitempty"`
	AccessKey  string          `json:"accessKey,omitempty"`
	SecretKey  string          `json:"secretKey,omitempty"`
}

// UpdateServiceAccountReq is the request options of the edit service account admin call
type UpdateServiceAccountReq struct {
	NewPolicy    json.RawMessage `json:"newPolicy,omitempty"` // Parsed policy from iam/policy.Parse
	NewSecretKey string          `json:"newSecretKey,omitempty"`
	NewStatus    string          `json:"newStatus,omityempty"`
}

// UpdateServiceAccount - edit an existing service account
func (adm *AdminClient) UpdateServiceAccount(ctx context.Context, accessKey string, opts UpdateServiceAccountReq) error {
	data, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	econfigBytes, err := EncryptData(adm.getSecretKey(), data)
	if err != nil {
		return err
	}

	queryValues := url.Values{}
	queryValues.Set("accessKey", accessKey)

	reqData := requestData{
		relPath:     adminAPIPrefix + "/update-service-account",
		content:     econfigBytes,
		queryValues: queryValues,
	}

	// Execute POST on /minio/admin/v3/update-service-account to edit a service account
	resp, err := adm.executeMethod(ctx, http.MethodPost, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return httpRespToErrorResponse(resp)
	}

	return nil
}

// ListServiceAccountsResp is the response body of the list service accounts call
type ListServiceAccountsResp struct {
	Accounts []string `json:"accounts"`
}

// ListServiceAccounts - list service accounts belonging to the specified user
func (adm *AdminClient) ListServiceAccounts(ctx context.Context, user string) (ListServiceAccountsResp, error) {
	queryValues := url.Values{}
	queryValues.Set("user", user)

	reqData := requestData{
		relPath:     adminAPIPrefix + "/list-service-accounts",
		queryValues: queryValues,
	}

	// Execute GET on /minio/admin/v3/list-service-accounts
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)
	defer closeResponse(resp)
	if err != nil {
		return ListServiceAccountsResp{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return ListServiceAccountsResp{}, httpRespToErrorResponse(resp)
	}

	data, err := DecryptData(adm.getSecretKey(), resp.Body)
	if err != nil {
		return ListServiceAccountsResp{}, err
	}

	var listResp ListServiceAccountsResp
	if err = json.Unmarshal(data, &listResp); err != nil {
		return ListServiceAccountsResp{}, err
	}
	return listResp, nil
}

// InfoServiceAccountResp is the response body of the info service account call
type InfoServiceAccountResp struct {
	ParentUser    string `json:"parentUser"`
	AccountStatus string `json:"accountStatus"`
	ImpliedPolicy bool   `json:"impliedPolicy"`
	Policy        string `json:"policy"`
}

// InfoServiceAccount - returns the info of service account belonging to the specified user
func (adm *AdminClient) InfoServiceAccount(ctx context.Context, accessKey string) (InfoServiceAccountResp, error) {
	queryValues := url.Values{}
	queryValues.Set("accessKey", accessKey)

	reqData := requestData{
		relPath:     adminAPIPrefix + "/info-service-account",
		queryValues: queryValues,
	}

	// Execute GET on /minio/admin/v3/info-service-account
	resp, err := adm.executeMethod(ctx, http.MethodGet, reqData)
	defer closeResponse(resp)
	if err != nil {
		return InfoServiceAccountResp{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return InfoServiceAccountResp{}, httpRespToErrorResponse(resp)
	}

	data, err := DecryptData(adm.getSecretKey(), resp.Body)
	if err != nil {
		return InfoServiceAccountResp{}, err
	}

	var infoResp InfoServiceAccountResp
	if err = json.Unmarshal(data, &infoResp); err != nil {
		return InfoServiceAccountResp{}, err
	}
	return infoResp, nil
}

// DeleteServiceAccount - delete a specified service account. The server will reject
// the request if the service account does not belong to the user initiating the request
func (adm *AdminClient) DeleteServiceAccount(ctx context.Context, serviceAccount string) error {
	queryValues := url.Values{}
	queryValues.Set("accessKey", serviceAccount)

	reqData := requestData{
		relPath:     adminAPIPrefix + "/delete-service-account",
		queryValues: queryValues,
	}

	// Execute DELETE on /minio/admin/v3/delete-service-account
	resp, err := adm.executeMethod(ctx, http.MethodDelete, reqData)
	defer closeResponse(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return httpRespToErrorResponse(resp)
	}

	return nil
}
