package restapi

import (
	"context"
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/madmin"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/console/restapi/operations/admin_api"
	"strings"
	"time"

	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
)

// Policy evaluated constants
const (
	Unknown = 0
	Allow   = 1
	Deny    = -1
)

//func registerUsersHandlers(api *operations.ConsoleAPI) {
//// List Users
//api.AdminAPIListUsersHandler = admin_api.ListUsersHandlerFunc(func(params admin_api.ListUsersParams, session *models.Principal) middleware.Responder {
//	listUsersResponse, err := getListUsersResponse(session)
//	if err != nil {
//		return admin_api.NewListUsersDefault(int(err.Code)).WithPayload(err)
//	}
//	return admin_api.NewListUsersOK().WithPayload(listUsersResponse)
//})
//// Add User
//api.AdminAPIAddUserHandler = admin_api.AddUserHandlerFunc(func(params admin_api.AddUserParams, session *models.Principal) middleware.Responder {
//	userResponse, err := getUserAddResponse(session, params)
//	if err != nil {
//		return admin_api.NewAddUserDefault(int(err.Code)).WithPayload(err)
//	}
//	return admin_api.NewAddUserCreated().WithPayload(userResponse)
//})
//// Remove User
//api.AdminAPIRemoveUserHandler = admin_api.RemoveUserHandlerFunc(func(params admin_api.RemoveUserParams, session *models.Principal) middleware.Responder {
//	err := getRemoveUserResponse(session, params)
//	if err != nil {
//		return admin_api.NewRemoveUserDefault(int(err.Code)).WithPayload(err)
//	}
//	return admin_api.NewRemoveUserNoContent()
//})
//// Update User-Groups
//api.AdminAPIUpdateUserGroupsHandler = admin_api.UpdateUserGroupsHandlerFunc(func(params admin_api.UpdateUserGroupsParams, session *models.Principal) middleware.Responder {
//	userUpdateResponse, err := getUpdateUserGroupsResponse(session, params)
//	if err != nil {
//		return admin_api.NewUpdateUserGroupsDefault(int(err.Code)).WithPayload(err)
//	}
//
//	return admin_api.NewUpdateUserGroupsOK().WithPayload(userUpdateResponse)
//})
//// Get User
//api.AdminAPIGetUserInfoHandler = admin_api.GetUserInfoHandlerFunc(func(params admin_api.GetUserInfoParams, session *models.Principal) middleware.Responder {
//	userInfoResponse, err := getUserInfoResponse(session, params)
//	if err != nil {
//		return admin_api.NewGetUserInfoDefault(int(err.Code)).WithPayload(err)
//	}
//
//	return admin_api.NewGetUserInfoOK().WithPayload(userInfoResponse)
//})
//// Update User
//api.AdminAPIUpdateUserInfoHandler = admin_api.UpdateUserInfoHandlerFunc(func(params admin_api.UpdateUserInfoParams, session *models.Principal) middleware.Responder {
//	userUpdateResponse, err := getUpdateUserResponse(session, params)
//	if err != nil {
//		return admin_api.NewUpdateUserInfoDefault(int(err.Code)).WithPayload(err)
//	}
//
//	return admin_api.NewUpdateUserInfoOK().WithPayload(userUpdateResponse)
//})
//}

func listUsers(ctx context.Context, client MinioAdmin) ([]*models.User, error) {
	userList, err := client.listUsers(ctx)
	if err != nil {
		return []*models.User{}, err
	}

	var users []*models.User
	for _, user := range userList {
		userElem := &models.User{
			AccessKey: *user.UserName,
			Status:    "",
			Policy:    strings.Split("", ","),
		}
		users = append(users, userElem)
	}

	return users, nil
}

// getListUsersResponse performs listUsers() and serializes it to the handler's output
func getListUsersResponse(session *models.Principal) (*models.ListUsersResponse, *models.Error) {
	ctx := context.Background()
	mAdmin, err := NewMinioAdminClient(session)
	if err != nil {
		return nil, prepareError(err)
	}
	adminClient := AdminClient{Client: mAdmin}

	users, err := listUsers(ctx, adminClient)
	if err != nil {
		return nil, prepareError(err)
	}
	// serialize output
	listUsersResponse := &models.ListUsersResponse{
		Users: users,
	}
	return listUsersResponse, nil
}

// addUser
func addUser(ctx context.Context, client MinioAdmin, accessKey, secretKey *string, groups []string, policies []string) (*models.User, error) {
	_, err := client.addUser(ctx, *accessKey, *secretKey)
	if err != nil {
		return nil, err
	}
	// set groups for the newly created user
	var userWithGroups *models.User
	if len(policies) > 0 {
		policyString := strings.Join(policies, ",")
		fmt.Println(policyString)
	}
	memberOf := []string{}
	status := "enabled"
	if userWithGroups != nil {
		memberOf = userWithGroups.MemberOf
		status = userWithGroups.Status
	}

	userRet := &models.User{
		AccessKey: *accessKey,
		MemberOf:  memberOf,
		Policy:    policies,
		Status:    status,
	}
	return userRet, nil
}

func getUserAddResponse(session *models.Principal, params admin_api.AddUserParams) (*models.User, *models.Error) {
	ctx := context.Background()
	mAdmin, err := NewMinioAdminClient(session)
	if err != nil {
		return nil, prepareError(err)
	}
	adminClient := AdminClient{Client: mAdmin}
	var userExists bool

	_, err = adminClient.getUserInfo(ctx, *params.Body.AccessKey)
	userExists = err == nil

	if userExists {
		return nil, prepareError(errNonUniqueAccessKey)
	}
	user, err := addUser(
		ctx,
		adminClient,
		params.Body.AccessKey,
		params.Body.SecretKey,
		params.Body.Groups,
		params.Body.Policies,
	)
	if err != nil {
		return nil, prepareError(err)
	}
	return user, nil
}

//removeUser invokes removing an user on `MinioAdmin`, then we return the response from API
func removeUser(ctx context.Context, client MinioAdmin, accessKey string) error {
	return client.removeUser(ctx, accessKey)
}

func getRemoveUserResponse(session *models.Principal, params admin_api.RemoveUserParams) *models.Error {
	ctx := context.Background()

	mAdmin, err := NewMinioAdminClient(session)
	if err != nil {
		return prepareError(err)
	}

	if session.AccountAccessKey == params.Name {
		return prepareError(errAvoidSelfAccountDelete)
	}
	adminClient := AdminClient{Client: mAdmin}

	if err := removeUser(ctx, adminClient, params.Name); err != nil {
		return prepareError(err)
	}

	return nil
}

// getUserInfo calls MinIO server get the User Information
func getUserInfo(ctx context.Context, client MinioAdmin, accessKey string) (*madmin.UserInfo, error) {
	userInfo, err := client.getUserInfo(ctx, accessKey)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func getUserInfoResponse(session *models.Principal, params admin_api.GetUserInfoParams) (*models.User, *models.Error) {
	ctx := context.Background()
	mAdmin, err := NewMinioAdminClient(session)
	if err != nil {
		return nil, prepareError(err)
	}
	adminClient := AdminClient{Client: mAdmin}
	user, err := getUserInfo(ctx, adminClient, params.Name)
	if err != nil {
		if madmin.ToErrorResponse(err).Code == "XMinioAdminNoSuchUser" {
			var errorCode int32 = 404
			errorMessage := "User doesn't exist"
			return nil, &models.Error{Code: errorCode, Message: swag.String(errorMessage), DetailedMessage: swag.String(err.Error())}
		}
		return nil, prepareError(err)
	}
	hasPolicy := true

	userInformation := &models.User{
		AccessKey: params.Name,
		Policy:    user.PolicyName,
		Status:    string(user.Status),
		HasPolicy: hasPolicy,
	}

	return userInformation, nil
}

// setUserStatus invokes setUserStatus from madmin to update user status
func setUserStatus(ctx context.Context, client MinioAdmin, user string, status string) error {
	var setStatus madmin.AccountStatus
	switch status {
	case "enabled":
		setStatus = madmin.AccountEnabled
	case "disabled":
		setStatus = madmin.AccountDisabled
	default:
		return errors.New(500, "status not valid")
	}

	return client.setUserStatus(ctx, user, setStatus)
}

// getUserSetPolicyResponse calls setUserAccessPolicy() to set a access policy to a user
//   and returns the serialized output.
func getUserSetPolicyResponse(session *models.Principal, userName string, req *models.SetUserPolicyRequest) *models.Error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	mAdmin, err := NewMinioAdminClient(session)
	if err != nil {
		return nil
	}
	adminClient := AdminClient{Client: mAdmin}

	if err := setUserAccessPolicy(ctx, adminClient, userName, *req.Access, req.Name, req.Definition); err != nil {
		return prepareError(err)
	}
	if err != nil {
		return prepareError(err)
	}
	return nil
}

// setUserAccessPolicy set the access permissions on an existing user.
func setUserAccessPolicy(ctx context.Context, client MinioAdmin, userName string, access models.BucketAccess, policyName, policyDefinition string) error {
	if strings.TrimSpace(userName) == "" {
		return fmt.Errorf("error: user name not present")
	}
	if strings.TrimSpace(string(access)) == "" {
		return fmt.Errorf("error: user access not present")
	}
	// Prepare policyJSON corresponding to the access type
	if access != models.BucketAccessPRIVATE && access != models.BucketAccessPUBLIC && access != models.BucketAccessCUSTOM {
		return fmt.Errorf("access: `%s` not supported", access)
	}
	if access == models.BucketAccessCUSTOM {
		err := client.putUserPolicy(ctx, userName, policyName, policyDefinition)
		if err != nil {
			return err
		}
	}
	return nil
}

// getUserPolicyResponse
func getUserPolicyResponse(session *models.Principal, userName string) (*madmin.UserPolicy, *models.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	mAdmin, err := NewMinioAdminClient(session)
	if err != nil {
		return nil, nil
	}
	adminClient := AdminClient{Client: mAdmin}
	policy, err := getUserAccessPolicy(ctx, adminClient, userName)
	if err != nil {
		return nil, prepareError(err)
	}
	return policy, nil
}

// getUserAccessPolicy
func getUserAccessPolicy(ctx context.Context, client MinioAdmin, userName string) (*madmin.UserPolicy, error) {
	if strings.TrimSpace(userName) == "" {
		return nil, fmt.Errorf("error: user name not present")
	}
	userPolicy, err := client.getUserPolicy(ctx, userName)
	if err != nil {
		return nil, err
	}
	return userPolicy, nil
}

// listUserPolicyResponse
func listUserPolicyResponse(session *models.Principal, userName string) (*madmin.UserPolicies, *models.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	mAdmin, err := NewMinioAdminClient(session)
	if err != nil {
		return nil, nil
	}
	adminClient := AdminClient{Client: mAdmin}
	policy, err := listUserAccessPolicy(ctx, adminClient, userName)
	if err != nil {
		return nil, prepareError(err)
	}
	return policy, nil
}

// listUserAccessPolicy
func listUserAccessPolicy(ctx context.Context, client MinioAdmin, userName string) (*madmin.UserPolicies, error) {
	if strings.TrimSpace(userName) == "" {
		return nil, fmt.Errorf("error: user name not present")
	}
	userPolicy, err := client.listUserPolicy(ctx, userName)
	if err != nil {
		return nil, err
	}
	return userPolicy, nil
}

// removeUserPolicyResponse
func removeUserPolicyResponse(session *models.Principal, userName, policyName string) *models.Error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	mAdmin, err := NewMinioAdminClient(session)
	if err != nil {
		return nil
	}
	adminClient := AdminClient{Client: mAdmin}
	err = removeUserAccessPolicy(ctx, adminClient, userName, policyName)
	if err != nil {
		return prepareError(err)
	}
	return nil
}

// removeUserAccessPolicy
func removeUserAccessPolicy(ctx context.Context, client MinioAdmin, userName, policyName string) error {
	if strings.TrimSpace(userName) == "" {
		return fmt.Errorf("error: bucket name not present")
	}
	err := client.removeUserPolicy(ctx, userName, policyName)
	if err != nil {
		return err
	}
	return nil
}
