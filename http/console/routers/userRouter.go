package routers

import (
	"github.com/filedag-project/filedag-storage/http/console/controllers"
	"github.com/gorilla/mux"
)

func registerRouter(router *mux.Router) {
	control := &controllers.Control{}
	// API Router
	apiRouter := router.PathPrefix("/user").Subrouter()

	apiRouter.Methods("POST").Path("/login").HandlerFunc(control.UserLogin)
	apiRouter.Methods("GET").Path("/list").HandlerFunc(control.ListUsers)
	apiRouter.Methods("POST").Path("/add").HandlerFunc(control.AddUser)
	apiRouter.Methods("POST").Path("/remove").HandlerFunc(control.RemoveUser)
	apiRouter.Methods("GET").Path("/info").HandlerFunc(control.UserInfo)
	apiRouter.Methods("POST").Path("/policy").HandlerFunc(control.SetUserPolicy)
	apiRouter.Methods("GET").Path("/policy").HandlerFunc(control.GetUserPolicy)
	apiRouter.Methods("GET").Path("/policy/list").HandlerFunc(control.ListUserPolicy)
	apiRouter.Methods("POST").Path("/policy/remove").HandlerFunc(control.RemoveUserPolicy)

	bucketRouter := router.PathPrefix("/bucket").Subrouter()

	bucketRouter.Methods("GET").Path("/list").HandlerFunc(control.ListBuckets)
	bucketRouter.Methods("POST").Path("/add").HandlerFunc(control.AddBucket)
	bucketRouter.Methods("POST").Path("/remove").HandlerFunc(control.RemoveBucket)
	bucketRouter.Methods("POST").Path("/policy").HandlerFunc(control.SetBucketPolicy)
	bucketRouter.Methods("GET").Path("/policy").HandlerFunc(control.GetBucketPolicy)
	bucketRouter.Methods("POST").Path("/policy/remove").HandlerFunc(control.RemoveBucketPolicy)
}

//NewServer Start a Server
func NewServer(router *mux.Router) {
	registerRouter(router)
}
