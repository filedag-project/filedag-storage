package main

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"io/ioutil"
	"net/http"
)

const (
	addUserUrl        = "/admin/v1/add-user"
	delUserUrl        = "/admin/v1/remove-user"
	getUserUrl        = "/admin/v1/user-info"
	changePassUserUrl = "/admin/v1/change-password"
	setStatusUrl      = "/admin/v1/update-accessKey_status"
)

var log = logging.Logger("tools")
var cmd = &cli.Command{
	Name:  "run",
	Usage: "a user cmd eg.demotools --method --addr=127.0.0.1:9985 --access-key=test --secret-key=test --username=wpg --pass=wpg123456",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "method",
			Usage: "the method",
		},
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of server eg.127.0.0.1:9985",
		},
		&cli.StringFlag{
			Name:  "access-key",
			Usage: "the access-key which have adduser policy",
		},
		&cli.StringFlag{
			Name:  "secret-key",
			Usage: "the secret-key which have adduser policy",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the username that you want add",
		},
		&cli.StringFlag{
			Name:  "pass",
			Usage: "the password that you want add",
		},
		&cli.StringFlag{
			Name:  "new-pass",
			Usage: "the new pass that you want change",
		},
		&cli.StringFlag{
			Name:  "status",
			Usage: "the status that you want set",
		},
	},
	Action: func(cctx *cli.Context) error {
		var req *http.Request
		var err error
		switch cctx.String("method") {
		case "adduser":
			req, err = add(cctx)
		case "deluser":
			req, err = del(cctx)
		case "getuser":
			req, err = get(cctx)
		case "change-pass":
			req, err = cha(cctx)
		case "set-status":
			req, err = set(cctx)
		default:
			return xerrors.Errorf("method not allowed :%v", cctx.String("method"))
		}
		if err != nil || req == nil {
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
		log.Infof("response: %v\n", string(all))
		return nil
	},
}

func add(cctx *cli.Context) (*http.Request, error) {
	var addr, accessKey, secretKey, username, pass string
	if cctx.String("addr") != "" {
		addr = cctx.String("addr")
	} else {
		log.Errorf("you must give the addr")
		return nil, xerrors.Errorf("you must give the addr")
	}
	if cctx.String("access-key") != "" {
		accessKey = cctx.String("access-key")
	} else {
		log.Errorf("you must give the access-key")
		return nil, xerrors.Errorf("you must give the access-key")
	}
	if cctx.String("secret-key") != "" {
		secretKey = cctx.String("secret-key")
	} else {
		log.Errorf("you must give the secret-key")
		return nil, xerrors.Errorf("you must give the secret-key")
	}
	if cctx.String("username") != "" {
		username = cctx.String("username")
	} else {
		log.Errorf("you must give the username")
		return nil, xerrors.Errorf("you must give the username")
	}
	if cctx.String("pass") != "" {
		pass = cctx.String("pass")
	} else {
		log.Errorf("you must give the pass")
		return nil, xerrors.Errorf("you must give the pass")
	}

	req, err := mustNewSignedV4Request(http.MethodPost, "http://"+addr+addUserUrl+"?accessKey="+username+"&secretKey="+pass,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		log.Errorf("mustNewSignedV4Request err: %v", err)
		return nil, err
	}
	return req, nil
}
func del(cctx *cli.Context) (*http.Request, error) {
	var addr, accessKey, secretKey, username string
	if cctx.String("addr") != "" {
		addr = cctx.String("addr")
	} else {
		return nil, xerrors.Errorf("you must give the addr")
	}
	if cctx.String("access-key") != "" {
		accessKey = cctx.String("access-key")
	} else {
		return nil, xerrors.Errorf("you must give the access-key")
	}
	if cctx.String("secret-key") != "" {
		secretKey = cctx.String("secret-key")
	} else {
		return nil, xerrors.Errorf("you must give the secret-key")
	}
	if cctx.String("username") != "" {
		username = cctx.String("username")
	} else {
		return nil, xerrors.Errorf("you must give the username")
	}
	req, err := mustNewSignedV4Request(http.MethodPost, "http://"+addr+delUserUrl+"?accessKey="+username,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func get(cctx *cli.Context) (*http.Request, error) {
	var addr, accessKey, secretKey, username string
	if cctx.String("addr") != "" {
		addr = cctx.String("addr")
	} else {
		return nil, xerrors.Errorf("you must give the addr")
	}
	if cctx.String("access-key") != "" {
		accessKey = cctx.String("access-key")
	} else {
		return nil, xerrors.Errorf("you must give the access-key")
	}
	if cctx.String("secret-key") != "" {
		secretKey = cctx.String("secret-key")
	} else {
		return nil, xerrors.Errorf("you must give the secret-key")
	}
	if cctx.String("username") != "" {
		username = cctx.String("username")
	} else {
		return nil, xerrors.Errorf("you must give the username")
	}

	req, err := mustNewSignedV4Request(http.MethodGet, "http://"+addr+getUserUrl+"?accessKey="+username,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func cha(cctx *cli.Context) (*http.Request, error) {
	var addr, accessKey, secretKey, username, pass string
	if cctx.String("addr") != "" {
		addr = cctx.String("addr")
	} else {
		return nil, xerrors.Errorf("you must give the addr")
	}
	if cctx.String("access-key") != "" {
		accessKey = cctx.String("access-key")
	} else {
		return nil, xerrors.Errorf("you must give the access-key")
	}
	if cctx.String("secret-key") != "" {
		secretKey = cctx.String("secret-key")
	} else {
		return nil, xerrors.Errorf("you must give the secret-key")
	}
	if cctx.String("username") != "" {
		username = cctx.String("username")
	} else {
		return nil, xerrors.Errorf("you must give the username")
	}
	if cctx.String("new-pass") != "" {
		pass = cctx.String("new-pass")
	} else {
		return nil, xerrors.Errorf("you must give the new pass")
	}

	req, err := mustNewSignedV4Request(http.MethodPost, "http://"+addr+changePassUserUrl+"?newPassword="+pass+"&accessKey="+username,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func set(cctx *cli.Context) (*http.Request, error) {
	var addr, accessKey, secretKey, username, status string
	if cctx.String("addr") != "" {
		addr = cctx.String("addr")
	} else {
		return nil, xerrors.Errorf("you must give the addr")
	}
	if cctx.String("access-key") != "" {
		accessKey = cctx.String("access-key")
	} else {
		return nil, xerrors.Errorf("you must give the access-key")
	}
	if cctx.String("secret-key") != "" {
		secretKey = cctx.String("secret-key")
	} else {
		return nil, xerrors.Errorf("you must give the secret-key")
	}
	if cctx.String("username") != "" {
		username = cctx.String("username")
	} else {
		return nil, xerrors.Errorf("you must give the username")
	}
	if cctx.String("status") != "" {
		status = cctx.String("status")
	} else {
		return nil, xerrors.Errorf("you must give the status")
	}

	req, err := mustNewSignedV4Request(http.MethodPost, "http://"+addr+setStatusUrl+"?accessKey="+username+"&status="+status,
		0, nil, "s3", accessKey, secretKey)
	if err != nil {
		return nil, err
	}
	return req, nil
}
