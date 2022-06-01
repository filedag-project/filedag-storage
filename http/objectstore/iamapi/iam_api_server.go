package iamapi

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/response"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/gorilla/mux"
	"net/http"
)

//iamApiServer the IamApi Server
type iamApiServer struct {
	authSys *iam.AuthSys
}

//NewIamApiServer New iamApiServer
func NewIamApiServer(router *mux.Router, db *uleveldb.ULevelDB) {
	iamApiSer := &iamApiServer{
		authSys: iam.NewAuthSys(db),
	}
	iamApiSer.registerRouter(router)

}

func (iamApi *iamApiServer) registerRouter(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/admin/v1").Subrouter()
	//root user
	apiRouter.Methods(http.MethodPost).Path("/add-user").HandlerFunc(iamApi.CreateUser).Queries("accessKey", "{accessKey:.*}", "secretKey", "{secretKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-user").HandlerFunc(iamApi.DeleteUser).Queries("accessKey", "{accessKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/change-password").HandlerFunc(iamApi.ChangePassword).Queries("accessKey", "{accessKey:.*}", "newPassword", "{newPassword:.*}")
	apiRouter.Methods(http.MethodPost).Path("/update-accessKey_status").HandlerFunc(iamApi.SetStatus).Queries("accessKey", "{accessKey:.*}", "status", "{status:.*}")
	apiRouter.Methods(http.MethodGet).Path("/user-info").HandlerFunc(iamApi.GetUserInfo).Queries("accessKey", "{accessKey:.*}")

	//sub user
	apiRouter.Methods(http.MethodPost).Path("/add-sub-user").HandlerFunc(iamApi.AddSubUser).Queries("userName", "{userName:.*}", "secretKey", "{secretKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-sub-user").HandlerFunc(iamApi.DeleteSubUser).Queries("userName", "{userName:.*}")
	apiRouter.Methods(http.MethodGet).Path("/sub-user-info").HandlerFunc(iamApi.GetSubUserInfo).Queries("userName", "{userName:.*}")
	apiRouter.Methods(http.MethodGet).Path("/list-all-sub-users").HandlerFunc(iamApi.GetUserList)

	apiRouter.Methods(http.MethodPost).Path("/put-sub-user-policy").HandlerFunc(iamApi.PutUserPolicy).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}", "policyDocument", "{policyDocument:.*}")
	apiRouter.Methods(http.MethodGet).Path("/get-sub-user-policy").HandlerFunc(iamApi.GetUserPolicy).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}")
	apiRouter.Methods(http.MethodGet).Path("/list-sub-user-policy").HandlerFunc(iamApi.ListUserPolicies).Queries("userName", "{userName:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-sub-user-policy").HandlerFunc(iamApi.DeleteUserPolicy).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}")

	//apiRouter.Methods(http.MethodPost).Path("/creat-policy").HandlerFunc(iamApi.CreatePolicy).Queries("policyName", "{policyName:.*}", "policyDocument", "{policyDocument:.*}")

	//apiRouter.Methods(http.MethodPost).Path("/creat-group").HandlerFunc(iamApi.CreatGroup).Queries("groupName", "{groupName:.*}", "version", "{version:.*}")
	//apiRouter.Methods(http.MethodGet).Path("/get_group").HandlerFunc(iamApi.GetGroup).Queries("groupName", "{groupName:.*}", "version", "{version:.*}")
	//apiRouter.Methods(http.MethodPost).Path("/delete-group").HandlerFunc(iamApi.DeleteGroup).Queries("groupName", "{groupName:.*}", "version", "{version:.*}")
	//apiRouter.Methods(http.MethodGet).Path("/list-groups").HandlerFunc(iamApi.ListGroups).Queries("pathPrefix", "{pathPrefix:.*}")

	apiRouter.NotFoundHandler = http.HandlerFunc(response.NotFoundHandler)
}
