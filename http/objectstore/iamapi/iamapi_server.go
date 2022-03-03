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
	apiRouter.Methods(http.MethodGet).Path("/user-list").HandlerFunc(iama.GetUserList)
	//
	// NotFound
	apiRouter.NotFoundHandler = http.HandlerFunc(s3resp.NotFoundHandler)
}
