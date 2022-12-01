package iamapi

import (
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/response"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	httpstatss "github.com/filedag-project/filedag-storage/objectservice/utils/httpstats"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
)

var log = logging.Logger("iamsever")

//iamApiServer the IamApi Server
type iamApiServer struct {
	authSys        *iam.AuthSys
	stats          *httpstatss.APIStatsSys
	cleanData      func(accessKey string)
	bucketInfoFunc func(ctx context.Context, accessKey string) []store.BucketInfo //todo use cache
}

//NewIamApiServer New iamApiServer
func NewIamApiServer(router *mux.Router, authSys *iam.AuthSys, stats *httpstatss.APIStatsSys, cleanData func(accessKey string), bucketInfoFunc func(ctx context.Context, accessKey string) []store.BucketInfo) {
	iamApiSer := &iamApiServer{
		authSys:        authSys,
		stats:          stats,
		cleanData:      cleanData,
		bucketInfoFunc: bucketInfoFunc,
	}
	iamApiSer.registerConsoleRouter(router, stats)
	iamApiSer.registerAdminsRouter(router, stats)

}

func (iamApi *iamApiServer) registerConsoleRouter(router *mux.Router, stats *httpstatss.APIStatsSys) {
	// API Router
	apiRouter := router.PathPrefix("/console/v1").Subrouter()
	//root user
	apiRouter.Methods(http.MethodPost).Path("/change-password").HandlerFunc(stats.RecordAPIHandler("change-password", iamApi.ChangePassword)).Queries("accessKey", "{accessKey:.*}", "newSecretKey", "{newSecretKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/update-accessKey_status").HandlerFunc(stats.RecordAPIHandler("update-accessKey_status", iamApi.SetStatus)).Queries("accessKey", "{accessKey:.*}", "status", "{status:.*}")
	apiRouter.Methods(http.MethodGet).Path("/user-info").HandlerFunc(stats.RecordAPIHandler("user-info", iamApi.AccountInfo)).Queries("accessKey", "{accessKey:.*}")

	//sub user
	//todo
	apiRouter.Methods(http.MethodPost).Path("/add-sub-user").HandlerFunc(stats.RecordAPIHandler("add-sub-user", iamApi.AddSubUser)).Queries("userName", "{userName:.*}", "secretKey", "{secretKey:.*}", "capacity", "{capacity:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-sub-user").HandlerFunc(stats.RecordAPIHandler("remove-sub-user", iamApi.DeleteSubUser)).Queries("userName", "{userName:.*}")
	apiRouter.Methods(http.MethodGet).Path("/sub-user-info").HandlerFunc(stats.RecordAPIHandler("sub-user-info", iamApi.GetSubUserInfo)).Queries("userName", "{userName:.*}")
	apiRouter.Methods(http.MethodGet).Path("/list-all-sub-users").HandlerFunc(stats.RecordAPIHandler("list-all-sub-users", iamApi.GetUserList))

	apiRouter.Methods(http.MethodPost).Path("/put-sub-user-policy").HandlerFunc(stats.RecordAPIHandler("put-sub-user-policy", iamApi.PutUserPolicy)).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}", "policyDocument", "{policyDocument:.*}")
	apiRouter.Methods(http.MethodGet).Path("/get-sub-user-policy").HandlerFunc(stats.RecordAPIHandler("get-sub-user-policy", iamApi.GetUserPolicy)).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}")
	apiRouter.Methods(http.MethodGet).Path("/list-sub-user-policy").HandlerFunc(stats.RecordAPIHandler("list-sub-user-policy", iamApi.ListUserPolicies)).Queries("userName", "{userName:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-sub-user-policy").HandlerFunc(stats.RecordAPIHandler("remove-sub-user-policy", iamApi.DeleteUserPolicy)).Queries("userName", "{userName:.*}", "policyName", "{policyName:.*}")

	//apiRouter.Methods(http.MethodPost).Path("/creat-policy").HandlerFunc(iamApi.CreatePolicy).Queries("policyName", "{policyName:.*}", "policyDocument", "{policyDocument:.*}")

	//apiRouter.Methods(http.MethodPost).Path("/creat-group").HandlerFunc(iamApi.CreatGroup).Queries("groupName", "{groupName:.*}", "version", "{version:.*}")
	//apiRouter.Methods(http.MethodGet).Path("/get_group").HandlerFunc(iamApi.GetGroup).Queries("groupName", "{groupName:.*}", "version", "{version:.*}")
	//apiRouter.Methods(http.MethodPost).Path("/delete-group").HandlerFunc(iamApi.DeleteGroup).Queries("groupName", "{groupName:.*}", "version", "{version:.*}")
	apiRouter.NotFoundHandler = http.HandlerFunc(response.NotFoundHandler)
}
func (iamApi *iamApiServer) registerAdminsRouter(router *mux.Router, stats *httpstatss.APIStatsSys) {
	// API Router
	apiRouter := router.PathPrefix("/admin/v1").Subrouter()
	//root user
	apiRouter.Methods(http.MethodPost).Path("/add-user").HandlerFunc(stats.RecordAPIHandler("add-user", iamApi.CreateUser)).Queries("accessKey", "{accessKey:.*}", "secretKey", "{secretKey:.*}", "capacity", "{capacity:.*}")
	apiRouter.Methods(http.MethodPost).Path("/remove-user").HandlerFunc(stats.RecordAPIHandler("remove-user", iamApi.DeleteUser)).Queries("accessKey", "{accessKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/change-password").HandlerFunc(stats.RecordAPIHandler("change-password", iamApi.ChangePassword)).Queries("accessKey", "{accessKey:.*}", "newSecretKey", "{newSecretKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/update-accessKey_status").HandlerFunc(stats.RecordAPIHandler("update-accessKey_status", iamApi.SetStatus)).Queries("accessKey", "{accessKey:.*}", "status", "{status:.*}")
	apiRouter.Methods(http.MethodGet).Path("/user-infos").HandlerFunc(stats.RecordAPIHandler("user-infos", iamApi.AccountInfos))
	apiRouter.Methods(http.MethodGet).Path("/request-overview").HandlerFunc(stats.RecordAPIHandler("request-overview", iamApi.RequestOverview))
	apiRouter.NotFoundHandler = http.HandlerFunc(response.NotFoundHandler)
}
