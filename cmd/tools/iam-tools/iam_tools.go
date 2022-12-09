package main

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	addUserUrl        = "/admin/v1/add-user"
	delUserUrl        = "/admin/v1/remove-user"
	getUserUrl        = "/admin/v1/user-info"
	changePassUserUrl = "/admin/v1/change-password"
	setStatusUrl      = "/admin/v1/update-accessKey_status"
)

const (
	ServerApi      = "server-api"
	AdminAccessKey = "admin-access-key"
	AdminSecretKey = "admin-secret-key"
)
const (
	// Minimum length for  access key.
	accessKeyMinLen = 3

	// Maximum length for  access key.
	// There is no max length enforcement for access keys
	accessKeyMaxLen = 20

	// Minimum length for  secret key for both server and gateway mode.
	secretKeyMinLen = 8

	// Maximum secret key length , this
	// is used when autogenerating new credentials.
	// There is no max length enforcement for secret keys
	secretKeyMaxLen = 40
)

var (
	errInvalidAccessKeyLength = fmt.Errorf("username(access key) length should be between %d and %d", accessKeyMinLen, accessKeyMaxLen)
	errInvalidSecretKeyLength = fmt.Errorf("password(secret key) length should be between %d and %d", secretKeyMinLen, secretKeyMaxLen)
)

var log = logging.Logger("tools")
var addUserCmd = &cli.Command{
	Name:  "add-user",
	Usage: "Add a new user",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  ServerApi,
			Usage: "the api of objectservice server",
			Value: "http://127.0.0.1:9985",
		},
		&cli.StringFlag{
			Name:  AdminAccessKey,
			Usage: "the access-key of user which have add user policy",
			Value: auth.DefaultAccessKey,
		},
		&cli.StringFlag{
			Name:  AdminSecretKey,
			Usage: "the secret-key of user which have add user policy",
			Value: auth.DefaultSecretKey,
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the username that you want to add",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "the password that you want to add",
		},
		&cli.StringFlag{
			Name:  "capacity",
			Usage: "the capacity of  the user you want to add",
		},
	},
	Action: func(cctx *cli.Context) error {
		apiAddr := cctx.String(ServerApi)
		if !strings.HasPrefix(apiAddr, "http") {
			return xerrors.Errorf("you should set the api of objectservice server")
		}
		accessKey := cctx.String(AdminAccessKey)
		if accessKey == "" {
			return xerrors.Errorf("you should give the admin-access-key")
		}
		secretKey := cctx.String(AdminSecretKey)
		if secretKey == "" {
			return xerrors.Errorf("you should give the admin-secret-key")
		}

		username := cctx.String("username")
		if !auth.IsAccessKeyValid(username) {
			return errInvalidAccessKeyLength
		}
		password := cctx.String("password")
		if !auth.IsSecretKeyValid(password) {
			return errInvalidSecretKeyLength
		}
		capacity := cctx.String("capacity")

		req, err := mustNewSignedV4Request(http.MethodPost, apiAddr+addUserUrl+"?accessKey="+username+"&secretKey="+password+"&capacity="+capacity,
			0, nil, "s3", accessKey, secretKey)
		if err != nil {
			log.Errorf("mustNewSignedV4Request err: %v", err)
			return err
		}
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("Do req err: %v", err)
			return err
		}
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("read resp err: %v", err)
			return err
		}
		fmt.Printf("response: %v\n", string(all))
		return nil
	},
}

var delUserCmd = &cli.Command{
	Name:  "remove-user",
	Usage: "Remove a user",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  ServerApi,
			Usage: "the api of objectservice server",
			Value: "http://127.0.0.1:9985",
		},
		&cli.StringFlag{
			Name:  AdminAccessKey,
			Usage: "the access-key of user which have remove user policy",
			Value: auth.DefaultAccessKey,
		},
		&cli.StringFlag{
			Name:  AdminSecretKey,
			Usage: "the secret-key of user which have remove user policy",
			Value: auth.DefaultSecretKey,
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the user that you want to remove",
		},
	},
	Action: func(cctx *cli.Context) error {
		apiAddr := cctx.String(ServerApi)
		if !strings.HasPrefix(apiAddr, "http") {
			return xerrors.Errorf("you should set the api of objectservice server")
		}
		accessKey := cctx.String(AdminAccessKey)
		if accessKey == "" {
			return xerrors.Errorf("you should give the admin-access-key")
		}
		secretKey := cctx.String(AdminSecretKey)
		if secretKey == "" {
			return xerrors.Errorf("you should give the admin-secret-key")
		}

		username := cctx.String("username")
		if !auth.IsAccessKeyValid(username) {
			return errInvalidAccessKeyLength
		}
		req, err := mustNewSignedV4Request(http.MethodPost, apiAddr+delUserUrl+"?accessKey="+username,
			0, nil, "s3", accessKey, secretKey)
		if err != nil {
			return err
		}
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("response：%v\n", string(all))
		return nil
	},
}

var getUserCmd = &cli.Command{
	Name:  "get-user",
	Usage: "Get a user info",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  ServerApi,
			Usage: "the api of objectservice server",
			Value: "http://127.0.0.1:9985",
		},
		&cli.StringFlag{
			Name:  AdminAccessKey,
			Usage: "the access-key of user which have get user policy",
			Value: auth.DefaultAccessKey,
		},
		&cli.StringFlag{
			Name:  AdminSecretKey,
			Usage: "the secret-key of user which have get user policy",
			Value: auth.DefaultSecretKey,
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the username that you want to get",
		},
	},
	Action: func(cctx *cli.Context) error {
		apiAddr := cctx.String(ServerApi)
		if !strings.HasPrefix(apiAddr, "http") {
			return xerrors.Errorf("you should set the api of objectservice server")
		}
		accessKey := cctx.String(AdminAccessKey)
		if accessKey == "" {
			return xerrors.Errorf("you should give the admin-access-key")
		}
		secretKey := cctx.String(AdminSecretKey)
		if secretKey == "" {
			return xerrors.Errorf("you should give the admin-secret-key")
		}

		username := cctx.String("username")
		if !auth.IsAccessKeyValid(username) {
			return errInvalidAccessKeyLength
		}
		req, err := mustNewSignedV4Request(http.MethodGet, apiAddr+getUserUrl+"?accessKey="+username,
			0, nil, "s3", accessKey, secretKey)
		if err != nil {
			return err
		}
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Printf("response：%v\n", string(all))
		return nil
	},
}

var changePassCmd = &cli.Command{
	Name:  "change-password",
	Usage: "Change the user password",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  ServerApi,
			Usage: "the api of objectservice server",
			Value: "http://127.0.0.1:9985",
		},
		&cli.StringFlag{
			Name:  AdminAccessKey,
			Usage: "the access-key of user which have change password policy",
			Value: auth.DefaultAccessKey,
		},
		&cli.StringFlag{
			Name:  AdminSecretKey,
			Usage: "the secret-key of user which have change password policy",
			Value: auth.DefaultSecretKey,
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the username that you want to change",
		},
		&cli.StringFlag{
			Name:  "new-password",
			Usage: "the new password that you want to set",
		},
	},
	Action: func(cctx *cli.Context) error {
		apiAddr := cctx.String(ServerApi)
		if !strings.HasPrefix(apiAddr, "http") {
			return xerrors.Errorf("you should set the api of objectservice server")
		}
		accessKey := cctx.String(AdminAccessKey)
		if accessKey == "" {
			return xerrors.Errorf("you should give the admin-access-key")
		}
		secretKey := cctx.String(AdminSecretKey)
		if secretKey == "" {
			return xerrors.Errorf("you should give the admin-secret-key")
		}

		username := cctx.String("username")
		if !auth.IsAccessKeyValid(username) {
			return errInvalidAccessKeyLength
		}
		password := cctx.String("new-password")
		if !auth.IsSecretKeyValid(password) {
			return errInvalidSecretKeyLength
		}

		req, err := mustNewSignedV4Request(http.MethodPost, apiAddr+changePassUserUrl+"?newSecretKey="+password+"&accessKey="+username,
			0, nil, "s3", accessKey, secretKey)
		if err != nil {
			return err
		}
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if string(all) != "" {
			fmt.Printf("response：%v\n", string(all))
		}
		return nil
	},
}

var setStatusCmd = &cli.Command{
	Name:  "set-status",
	Usage: "Set the user status",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  ServerApi,
			Usage: "the api of objectservice server",
			Value: "http://127.0.0.1:9985",
		},
		&cli.StringFlag{
			Name:  AdminAccessKey,
			Usage: "the access-key of user which have set status policy",
			Value: auth.DefaultAccessKey,
		},
		&cli.StringFlag{
			Name:  AdminSecretKey,
			Usage: "the secret-key of user which have set status policy",
			Value: auth.DefaultSecretKey,
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the username that you want to set status",
		},
		&cli.StringFlag{
			Name:  "status",
			Usage: "set the status of account, enum: on,off",
		},
	},
	Action: func(cctx *cli.Context) error {
		apiAddr := cctx.String(ServerApi)
		if !strings.HasPrefix(apiAddr, "http") {
			return xerrors.Errorf("you should set the api of objectservice server")
		}
		accessKey := cctx.String(AdminAccessKey)
		if accessKey == "" {
			return xerrors.Errorf("you should give the admin-access-key")
		}
		secretKey := cctx.String(AdminSecretKey)
		if secretKey == "" {
			return xerrors.Errorf("you should give the admin-secret-key")
		}

		username := cctx.String("username")
		if !auth.IsAccessKeyValid(username) {
			return errInvalidAccessKeyLength
		}
		status := cctx.String("status")
		switch status {
		case "on", "off":
		default:
			return xerrors.Errorf("invalid status, you should give the valid status, 'on' or 'off'")
		}

		req, err := mustNewSignedV4Request(http.MethodPost, apiAddr+setStatusUrl+"?accessKey="+username+"&status="+status,
			0, nil, "s3", accessKey, secretKey)
		if err != nil {
			return err
		}
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		all, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if string(all) != "" {
			fmt.Printf("response：%v\n", string(all))
		}
		return nil
	},
}
