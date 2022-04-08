package routers

import (
	"github.com/filedag-project/filedag-storage/http/console/controllers"
	"github.com/gorilla/mux"
)

func registerRouter(router *mux.Router) {
	control := &controllers.Control{}
	// API Router
	userRouter := router.PathPrefix("/user").Subrouter()

	userRouter.Methods("POST").Path("/login").HandlerFunc(control.UserLogin)
	userRouter.Methods("GET").Path("/list").HandlerFunc(control.ListUsers)
	userRouter.Methods("POST").Path("/add").HandlerFunc(control.AddUser)
	userRouter.Methods("POST").Path("/remove").HandlerFunc(control.RemoveUser)
	userRouter.Methods("GET").Path("/info").HandlerFunc(control.UserInfo)
	userRouter.Methods("POST").Path("/policy").HandlerFunc(control.SetUserPolicy)
	userRouter.Methods("GET").Path("/policy").HandlerFunc(control.GetUserPolicy)
	userRouter.Methods("GET").Path("/policy/list").HandlerFunc(control.ListUserPolicy)
	userRouter.Methods("POST").Path("/policy/remove").HandlerFunc(control.RemoveUserPolicy)

	bucketRouter := router.PathPrefix("/bucket").Subrouter()
	bucketRouter.Methods("GET").Path("/list").HandlerFunc(control.ListBuckets)
	bucketRouter.Methods("POST").Path("/add").HandlerFunc(control.AddBucket)
	bucketRouter.Methods("POST").Path("/remove").HandlerFunc(control.RemoveBucket)
	bucketRouter.Methods("POST").Path("/policy").HandlerFunc(control.SetBucketPolicy)
	bucketRouter.Methods("GET").Path("/policy").HandlerFunc(control.GetBucketPolicy)
	bucketRouter.Methods("POST").Path("/policy/remove").HandlerFunc(control.RemoveBucketPolicy)

	objectRouter := router.PathPrefix("/object").Subrouter()
	objectRouter.Methods("GET").Path("/list").HandlerFunc(control.ListObjects)
}

//NewServer Start a Server
func NewServer(router *mux.Router) {
	registerRouter(router)
}
