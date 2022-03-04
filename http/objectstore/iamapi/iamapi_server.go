package iamapi

// https://docs.aws.amazon.com/cli/latest/reference/iam/list-roles.html

import (
	"context"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api/s3resp"
	"github.com/gorilla/mux"
	"net/http"
)

type IamApiServer struct {
}

func NewIamApiServer(router *mux.Router) {
	iamApiServer := &IamApiServer{}
	iam.GlobalIAMSys.Init(context.Background())
	iamApiServer.registerRouter(router)

}

func (iama *IamApiServer) registerRouter(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/admin").Subrouter()
	apiRouter.Methods(http.MethodGet).Path("/list-user").HandlerFunc(iama.GetUserList)
	apiRouter.Methods(http.MethodPost).Path("/add-user").HandlerFunc(iama.AddUser).Queries("accessKey", "{accessKey:.*}", "secretKey", "{secretKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-user").HandlerFunc(iama.RemoveUser).Queries("accessKey", "{accessKey:.*}")

	apiRouter.Methods(http.MethodPost).Path("/put-user-policy").HandlerFunc(iama.PutUserPolicy).Queries("accessKey", "{accessKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/get-user-policy").HandlerFunc(iama.GetUserPolicy).Queries("accessKey", "{accessKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-user-policy").HandlerFunc(iama.DeleteUserPolicy).Queries("accessKey", "{accessKey:.*}")

	apiRouter.NotFoundHandler = http.HandlerFunc(s3resp.NotFoundHandler)
}
