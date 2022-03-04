package iamapi

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api/s3resp"
	"net/http"
)

//GetUserList get all user
func (iama *IamApiServer) GetUserList(w http.ResponseWriter, r *http.Request) {
	var resp ListUsersResponse
	resp.ListUsersResult.Users = iam.GlobalIAMSys.GetUserList(context.Background())
	s3resp.WriteXMLResponse(w, r, http.StatusOK, resp)
}

// AddUser  add user
func (iama *IamApiServer) AddUser(w http.ResponseWriter, r *http.Request) {
	var resp CreateUserResponse
	values := r.URL.Query()
	accessKey := values.Get("accessKey")
	secretKey := values.Get("secretKey")
	resp.CreateUserResult.User.UserName = &accessKey
	err := iam.GlobalIAMSys.AddUser(context.Background(), accessKey, secretKey)
	if err != nil {
		s3resp.WriteErrorResponse(w, r, s3resp.ErrInternalError)
		return
	}
	s3resp.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//RemoveUser delete user
func (iama *IamApiServer) RemoveUser(w http.ResponseWriter, r *http.Request) {
	var resp CreateUserResponse
	accessKey := r.FormValue("accessKey")
	resp.CreateUserResult.User.UserName = &accessKey
	err := iam.GlobalIAMSys.RemoveUser(context.Background(), accessKey)
	if err != nil {
		s3resp.WriteErrorResponse(w, r, s3resp.ErrInternalError)
		return
	}
	s3resp.WriteXMLResponse(w, r, http.StatusOK, resp)
}

//PutUserPolicy Put UserPolicy
func (iama *IamApiServer) PutUserPolicy(w http.ResponseWriter, r *http.Request) {

}

//GetUserPolicy  Get UserPolicy
func (iama *IamApiServer) GetUserPolicy(w http.ResponseWriter, r *http.Request) {

}

//DeleteUserPolicy Delete eUserPolicy
func (iama *IamApiServer) DeleteUserPolicy(w http.ResponseWriter, r *http.Request) {

}
