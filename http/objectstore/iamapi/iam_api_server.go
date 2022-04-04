package iamapi

// https://docs.aws.amazon.com/cli/latest/reference/iam/list-roles.html

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/gorilla/mux"
	"net/http"
)

//iamApiServer the IamApi Server
type iamApiServer struct {
	authSys iam.AuthSys
}

//NewIamApiServer New iamApiServer
func NewIamApiServer(router *mux.Router) {
	iamApiSer := &iamApiServer{}
	iamApiSer.authSys.Init()
	iamApiSer.registerRouter(router)

}

func (iamApi *iamApiServer) registerRouter(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/admin/v1").Subrouter()
	apiRouter.Methods(http.MethodGet).Path("/list-user").HandlerFunc(iamApi.GetUserList)
	apiRouter.Methods(http.MethodPost).Path("/add-user").HandlerFunc(iamApi.CreateUser).Queries("accessKey", "{accessKey:.*}", "secretKey", "{secretKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-user").HandlerFunc(iamApi.DeleteUser).Queries("accessKey", "{accessKey:.*}")

	apiRouter.Methods(http.MethodPost).Path("/put-user-policy").HandlerFunc(iamApi.PutUserPolicy).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}", "policyDocument", "{policyDocument:.*}")
	apiRouter.Methods(http.MethodGet).Path("/get-user-policy").HandlerFunc(iamApi.GetUserPolicy).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}")
	apiRouter.Methods(http.MethodGet).Path("/list-user-policy").HandlerFunc(iamApi.ListUserPolicies).Queries("userName", "{userName:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-user-policy").HandlerFunc(iamApi.DeleteUserPolicy).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}")
	apiRouter.Methods(http.MethodPost).Path("/creat-policy").HandlerFunc(iamApi.CreatePolicy).Queries("policyName", "{policyName:.*}", "policyDocument", "{policyDocument:.*}")
	apiRouter.Methods(http.MethodGet).Path("/user-info").HandlerFunc(iamApi.GetUserInfo).Queries("userName", "{userName:.*}")
	apiRouter.Methods(http.MethodGet).Path("/change-password").HandlerFunc(iamApi.ChangePassword).Queries("OldPassword", "{OldPassword:.*}", "NewPassword", "{NewPassword:.*}", "Version", "{Version:.*}")
	apiRouter.Methods(http.MethodGet).Path("/update-accessKey_status").HandlerFunc(iamApi.SetStatus).Queries("userName", "{userName:.*}")

	apiRouter.NotFoundHandler = http.HandlerFunc(response.NotFoundHandler)
}
