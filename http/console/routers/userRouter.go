package routers

import (
	"github.com/filedag-project/filedag-storage/http/console/controllers"
	"github.com/gorilla/mux"
)

func registerRouter(router *mux.Router) {
	control := &controllers.Control{}
	// API Router
	apiRouter := router.PathPrefix("/").Subrouter()

	// Readiness Probe
	apiRouter.Methods("POST").Path("/login").HandlerFunc(control.UserLogin)
	apiRouter.Methods("GET").Path("/user/list").HandlerFunc(control.ListUsers)
	apiRouter.Methods("POST").Path("/user/add").HandlerFunc(control.AddUser)
	apiRouter.Methods("POST").Path("/user/remove").HandlerFunc(control.RemoveUser)
	apiRouter.Methods("GET").Path("/user/info").HandlerFunc(control.UserInfo)
	apiRouter.Methods("POST").Path("/user/policy").HandlerFunc(control.SetUserPolicy)
	apiRouter.Methods("GET").Path("/user/policy").HandlerFunc(control.GetUserPolicy)
	apiRouter.Methods("GET").Path("/user/policy/list").HandlerFunc(control.ListUserPolicy)
	apiRouter.Methods("POST").Path("/user/policy/remove").HandlerFunc(control.RemoveUserPolicy)
}

//NewServer Start a Server
func NewServer(router *mux.Router) {
	registerRouter(router)
}
