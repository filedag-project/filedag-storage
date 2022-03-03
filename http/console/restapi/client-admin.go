// This file is part of MinIO Console Server
// Copyright (c) 2021 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package restapi

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/console/credentials"
	"github.com/filedag-project/filedag-storage/http/console/madmin"
	models "github.com/filedag-project/filedag-storage/http/console/model"
	"net/http"
)

const globalAppName = "MinIO Console"

// MinioAdmin interface with all functions to be implemented
// by mock when testing, it should include all MinioAdmin respective api calls
// that are used within this project.
type MinioAdmin interface {
	listUsers(ctx context.Context) (map[string]madmin.UserInfo, error)
	addUser(ctx context.Context, acessKey, SecretKey string) error
	removeUser(ctx context.Context, accessKey string) error
	getUserInfo(ctx context.Context, accessKey string) (madmin.UserInfo, error)
	setUserStatus(ctx context.Context, accessKey string, status madmin.AccountStatus) error
	listGroups(ctx context.Context) ([]string, error)
	updateGroupMembers(ctx context.Context, greq madmin.GroupAddRemove) error
	getGroupDescription(ctx context.Context, group string) (*madmin.GroupDesc, error)
	setGroupStatus(ctx context.Context, group string, status madmin.GroupStatus) error
}

// Interface implementation
//
// Define the structure of a minIO Client and define the functions that are actually used
// from minIO api.
type AdminClient struct {
	Client *madmin.AdminClient
}

func (ac AdminClient) changePassword(ctx context.Context, accessKey, secretKey string) error {
	return ac.Client.SetUser(ctx, accessKey, secretKey, madmin.AccountEnabled)
}

// implements madmin.ListUsers()
func (ac AdminClient) listUsers(ctx context.Context) (map[string]madmin.UserInfo, error) {
	return ac.Client.ListUsers(ctx)
}

// implements madmin.AddUser()
func (ac AdminClient) addUser(ctx context.Context, accessKey, secretKey string) error {
	return ac.Client.AddUser(ctx, accessKey, secretKey)
}

// implements madmin.RemoveUser()
func (ac AdminClient) removeUser(ctx context.Context, accessKey string) error {
	return ac.Client.RemoveUser(ctx, accessKey)
}

//implements madmin.GetUserInfo()
func (ac AdminClient) getUserInfo(ctx context.Context, accessKey string) (madmin.UserInfo, error) {
	return ac.Client.GetUserInfo(ctx, accessKey)
}

// implements madmin.SetUserStatus()
func (ac AdminClient) setUserStatus(ctx context.Context, accessKey string, status madmin.AccountStatus) error {
	return ac.Client.SetUserStatus(ctx, accessKey, status)
}

// implements madmin.ListGroups()
func (ac AdminClient) listGroups(ctx context.Context) ([]string, error) {
	return ac.Client.ListGroups(ctx)
}

// implements madmin.UpdateGroupMembers()
func (ac AdminClient) updateGroupMembers(ctx context.Context, greq madmin.GroupAddRemove) error {
	return ac.Client.UpdateGroupMembers(ctx, greq)
}

// implements madmin.GetGroupDescription(group)
func (ac AdminClient) getGroupDescription(ctx context.Context, group string) (*madmin.GroupDesc, error) {
	return ac.Client.GetGroupDescription(ctx, group)
}

// implements madmin.SetGroupStatus(group, status)
func (ac AdminClient) setGroupStatus(ctx context.Context, group string, status madmin.GroupStatus) error {
	return ac.Client.SetGroupStatus(ctx, group, status)
}

// implements madmin.ListServiceAccounts()
func (ac AdminClient) listServiceAccounts(ctx context.Context, user string) (madmin.ListServiceAccountsResp, error) {
	// TODO: Fix this
	return ac.Client.ListServiceAccounts(ctx, user)
}

// implements madmin.DeleteServiceAccount()
func (ac AdminClient) deleteServiceAccount(ctx context.Context, serviceAccount string) error {
	return ac.Client.DeleteServiceAccount(ctx, serviceAccount)
}

// implements madmin.InfoServiceAccount()
func (ac AdminClient) infoServiceAccount(ctx context.Context, serviceAccount string) (madmin.InfoServiceAccountResp, error) {
	return ac.Client.InfoServiceAccount(ctx, serviceAccount)
}

// implements madmin.UpdateServiceAccount()
func (ac AdminClient) updateServiceAccount(ctx context.Context, serviceAccount string, opts madmin.UpdateServiceAccountReq) error {
	return ac.Client.UpdateServiceAccount(ctx, serviceAccount, opts)
}

// AccountInfo implements madmin.AccountInfo()
func (ac AdminClient) AccountInfo(ctx context.Context) (madmin.AccountInfo, error) {
	return ac.Client.AccountInfo(ctx, madmin.AccountOpts{})
}

func NewMinioAdminClient(sessionClaims *models.Principal) (*madmin.AdminClient, error) {
	adminClient, err := newAdminFromClaims(sessionClaims)
	if err != nil {
		return nil, err
	}
	return adminClient, nil
}

// newAdminFromClaims creates a minio admin from Decrypted claims using Assume role credentials
func newAdminFromClaims(claims *models.Principal) (*madmin.AdminClient, error) {
	tlsEnabled := getMinIOEndpointIsSecure()
	endpoint := getMinIOEndpoint()

	adminClient, err := madmin.NewWithOptions(endpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(claims.STSAccessKeyID, claims.STSSecretAccessKey, claims.STSSessionToken),
		Secure: tlsEnabled,
	})
	if err != nil {
		return nil, err
	}
	adminClient.SetCustomTransport(GetConsoleHTTPClient().Transport)
	return adminClient, nil
}

// newAdminFromCreds Creates a minio client using custom credentials for connecting to a remote host
func newAdminFromCreds(accessKey, secretKey, endpoint string, tlsEnabled bool) (*madmin.AdminClient, error) {
	minioClient, err := madmin.NewWithOptions(endpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: tlsEnabled,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

// httpClient is a custom http client, this client should not be called directly and instead be
// called using GetConsoleHTTPClient() to ensure is initialized and the certificates are loaded correctly
var httpClient *http.Client

// GetConsoleHTTPClient will initialize the console HTTP Client with fully populated custom TLS
// Transport that with loads certs at
// - ${HOME}/.console/certs/CAs
// - ${HOME}/.minio/certs/CAs
func GetConsoleHTTPClient() *http.Client {
	if httpClient == nil {
		httpClient = PrepareConsoleHTTPClient(false)
	}
	return httpClient
}

