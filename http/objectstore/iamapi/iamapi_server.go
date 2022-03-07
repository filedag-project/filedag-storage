package iamapi

// https://docs.aws.amazon.com/cli/latest/reference/iam/list-roles.html

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api/s3resp"
	"github.com/gorilla/mux"
	"net/http"
)

//iamApiServer the IamApi Server
type iamApiServer struct {
}

//NewIamApiServer New iamApiServer
func NewIamApiServer(router *mux.Router) {
	iamApiSer := &iamApiServer{}
	iam.GlobalIAMSys.Init(context.Background())
	iamApiSer.registerRouter(router)

}

func (iamApi *iamApiServer) registerRouter(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/admin").Subrouter()
	apiRouter.Methods(http.MethodGet).Path("/list-user").HandlerFunc(iamApi.GetUserList)
	apiRouter.Methods(http.MethodPost).Path("/add-user").HandlerFunc(iamApi.AddUser).Queries("accessKey", "{accessKey:.*}", "secretKey", "{secretKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-user").HandlerFunc(iamApi.RemoveUser).Queries("accessKey", "{accessKey:.*}")

	apiRouter.Methods(http.MethodPost).Path("/put-user-policy").HandlerFunc(iamApi.PutUserPolicy).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}", "policyDocument", "{policyDocument:.*}")
	apiRouter.Methods(http.MethodPost).Path("/get-user-policy").HandlerFunc(iamApi.GetUserPolicy).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-user-policy").HandlerFunc(iamApi.RemoveUserPolicy).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}")
	apiRouter.Methods(http.MethodPost).Path("/creat-policy").HandlerFunc(iamApi.CreatePolicy).Queries("policyName", "{policyName:.*}", "policyDocument", "{policyDocument:.*}")

	apiRouter.NotFoundHandler = http.HandlerFunc(s3resp.NotFoundHandler)
}
