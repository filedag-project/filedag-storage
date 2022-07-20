package main

import (
	"context"
	"errors"
	"fmt"
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/iamapi"
	"github.com/filedag-project/filedag-storage/http/objectstore/s3api"
	"github.com/filedag-project/filedag-storage/http/objectstore/uleveldb"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	logging "github.com/ipfs/go-log/v2"
	"github.com/ipfs/go-merkledag"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	EnvRootUser     = "FILEDAG_ROOT_USER"
	EnvRootPassword = "FILEDAG_ROOT_PASSWORD"
)

var log = logging.Logger("sever")

func missingCredentialError(user, pwd string) error {
	return errors.New(fmt.Sprintf("Missing credential environment variable, user is \"%s\" and password is\"%s\"."+
		" Root user and password are expected to be specified via environment variables "+
		"FILEDAG_ROOT_USER and FILEDAG_ROOT_PASSWORD respectively", user, pwd))
}

//startServer Start a IamServer
func startServer(cctx *cli.Context) {
	listen := cctx.String("listen")
	datadir := cctx.String("datadir")
	poolAddr := cctx.String("pool-addr")
	poolUser := cctx.String("pool-user")
	poolPassword := cctx.String("pool-password")

	user := cctx.String("root-user")
	password := cctx.String("root-password")
	if user == "" || password == "" {
		log.Fatal(missingCredentialError(user, password))
	}
	cred, err := auth.CreateCredentials(user, password)
	if err != nil {
		log.Fatal("Invalid credentials. Please provide correct credentials. " +
			"Root user length should be at least 3, and password length at least 8 characters")
	}

	db, err := uleveldb.OpenDb(datadir)
	if err != nil {
		return
	}
	defer db.Close()
	router := mux.NewRouter()
	authSys := iam.NewAuthSys(db, cred)
	iamapi.NewIamApiServer(router, authSys)
	poolClient, err := dagpoolcli.NewPoolClient(poolAddr, poolUser, poolPassword)
	if err != nil {
		log.Fatalf("connect dagpool server err: %v", err)
	}
	defer poolClient.Close(context.TODO())
	dagServ := merkledag.NewDAGService(blockservice.New(poolClient, offline.Exchange(poolClient)))
	s3api.NewS3Server(router, dagServ, poolClient, authSys, db)
	if strings.HasPrefix(listen, ":") {
		for _, ip := range utils.MustGetLocalIP4().ToSlice() {
			log.Infof("start sever at http://%v%v", ip, listen)
		}
	} else {
		log.Infof("start sever at http://%v", listen)
	}
	go func() {
		if err = http.ListenAndServe(listen, router); err != nil {
			log.Errorf("Listen And Serve err%v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")
	log.Info("Server exit")
}

var startCmd = &cli.Command{
	Name:  "daemon",
	Usage: "Start a filedag storage process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "listen",
			Usage: "set server listen",
			Value: ":9985",
		},
		&cli.StringFlag{
			Name:  "datadir",
			Usage: "directory to store data in",
			Value: "./store-data",
		},
		&cli.StringFlag{
			Name:  "pool-addr",
			Usage: "set the pool rpc address you want connect",
		},
		&cli.StringFlag{
			Name:  "pool-user",
			Usage: "set pool user",
		},
		&cli.StringFlag{
			Name:  "pool-password",
			Usage: "set pool password",
		},
		&cli.StringFlag{
			Name:    "root-user",
			Usage:   "set root filedag root user",
			EnvVars: []string{EnvRootUser},
			Value:   auth.DefaultAccessKey,
		},
		&cli.StringFlag{
			Name:    "root-password",
			Usage:   "set root filedag root password",
			EnvVars: []string{EnvRootPassword},
			Value:   auth.DefaultSecretKey,
		},
	},
	Action: func(cctx *cli.Context) error {
		startServer(cctx)
		return nil
	},
}
