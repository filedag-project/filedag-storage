package pool

import (
	"github.com/gorilla/mux"
	"net/http"
)

type DagPoolServer struct {
	dp *DagPool
}

//registerS3Router Register S3Router
func (d *DagPoolServer) registerS3Router(router *mux.Router) {
	// API Router
	apiRouter := router.PathPrefix("/dagpool/").Subrouter()
	apiRouter.Methods(http.MethodPost).Path("/add").HandlerFunc(d.Add).Queries("user", "{user:.*}", "secretKey", "{secretKey:.*}")
	apiRouter.Methods(http.MethodPost).Path("/get").HandlerFunc(d.Get).Queries("user", "{user:.*}", "secretKey", "{secretKey:.*}", "cid", "{cid:.*}")
	apiRouter.Methods(http.MethodPost).Path("/delete").HandlerFunc(d.Delete).Queries("user", "{user:.*}", "secretKey", "{secretKey:.*}", "cid", "{cid:.*}")
}

func NewDagPoolServer() *DagPoolServer {
	service, err := NewDagPoolService()
	if err != nil {
		return nil
	}
	dagServer := &DagPoolServer{dp: service}
	r := mux.NewRouter()
	dagServer.registerS3Router(r)
	return dagServer
}
