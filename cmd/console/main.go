package main

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/console/routers"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
)

// main for server.
func main() {
	err := logging.SetLogLevel("*", "INFO")
	if err != nil {
		return
	}
	router := mux.NewRouter()
	routers.NewServer(router)
	err = http.ListenAndServe(":9090", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}
