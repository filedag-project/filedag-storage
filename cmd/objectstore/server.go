package main

import (
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/set"
	"github.com/filedag-project/filedag-storage/http/objectstore/iamapi"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
	"github.com/gorilla/mux"
	logging "github.com/ipfs/go-log/v2"
	"net"
	"net/http"
	"runtime"
)

var log = logging.Logger("sever")

//StartServer Start a IamServer
func StartServer() {
	err := logging.SetLogLevel("*", "INFO")
	if err != nil {
		return
	}
	router := mux.NewRouter().SkipClean(true)
	s3api.NewS3Server(router)
	iamapi.NewIamApiServer(router)
	for _, ip := range mustGetLocalIP4().ToSlice() {
		log.Infof("start sever at http://%v:%v", ip, 9985)
	}
	http.ListenAndServe(":9985", router)

}

// mustGetLocalIP4 returns IPv4 addresses of localhost.  It panics on error.
func mustGetLocalIP4() (ipList set.StringSet) {
	ipList = set.NewStringSet()
	ifs, err := net.Interfaces()
	if err != nil {
		log.Errorf("Unable to get IP addresses of this host %v", err)

	}

	for _, interf := range ifs {
		addrs, err := interf.Addrs()
		if err != nil {
			continue
		}
		if runtime.GOOS == "windows" && interf.Flags&net.FlagUp == 0 {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.To4() != nil {
				ipList.Add(ip.String())
			}
		}
	}

	return ipList
}
