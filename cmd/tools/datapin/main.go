package main

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/http/objectstore/iam/auth"
	"github.com/filedag-project/filedag-storage/http/objectstore/utils"
	"github.com/google/martian/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	utils.SetupLogLevels()
	local := []*cli.Command{
		pin,
		unpin,
	}
	app := &cli.App{
		Name:                 "object-client",
		Usage:                "test pin/unpin interface",
		Version:              "0.0.1",
		EnableBashCompletion: true,
		Commands:             local,
	}
	app.Setup()
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

var pin = &cli.Command{
	Name:  "pin",
	Usage: "pin a block eg.go run main.go pin --addr= --access-key= --secret-key= --bucket= --object=",
	Flags: []cli.Flag{
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
			Name:  "bucket",
			Usage: "the bucket cannot be empty",
		},
		&cli.StringFlag{
			Name:  "object",
			Usage: "the object cannot be empty",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, accessKey, secretKey, bucket, object string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			return xerrors.Errorf("you must enter a addr")
		}
		if cctx.String("access-key") != "" {
			accessKey = cctx.String("access-key")
		} else {
			return xerrors.Errorf("you must enter a access-key")
		}
		if cctx.String("secret-key") != "" {
			secretKey = cctx.String("secret-key")
		} else {
			return xerrors.Errorf("you must enter a secret-key")
		}
		if cctx.String("bucket") != "" {
			bucket = cctx.String("bucket")
		} else {
			return xerrors.Errorf("you must enter a bucket")
		}
		if cctx.String("object") != "" {
			object = cctx.String("object")
		} else {
			return xerrors.Errorf("you must enter a object")
		}
		fmt.Println("++++++++")
		req, err := mustNewSignedV4Request(http.MethodPost, "http://"+addr+"/pin?bucket="+bucket+"&object="+object, 0, nil, "s3", accessKey, secretKey)
		client := http.DefaultClient
		resp, err := client.Do(req)
		fmt.Println(resp)
		if err != nil {
			fmt.Println("12")
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

var unpin = &cli.Command{
	Name:  "unpin",
	Usage: "unpin a block eg.go run main.go unpin --addr= --access-key= --secret-key= --bucket= --object=",
	Flags: []cli.Flag{
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
			Name:  "bucket",
			Usage: "the bucket cannot be empty",
		},
		&cli.StringFlag{
			Name:  "object",
			Usage: "the object cannot be empty",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, accessKey, secretKey, bucket, object string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			return xerrors.Errorf("you must enter a addr")
		}
		if cctx.String("access-key") != "" {
			accessKey = cctx.String("access-key")
		} else {
			return xerrors.Errorf("you must enter a access-key")
		}
		if cctx.String("secret-key") != "" {
			secretKey = cctx.String("secret-key")
		} else {
			return xerrors.Errorf("you must enter a secret-key")
		}
		if cctx.String("bucket") != "" {
			bucket = cctx.String("bucket")
		} else {
			return xerrors.Errorf("you must enter a bucket")
		}
		if cctx.String("object") != "" {
			object = cctx.String("object")
		} else {
			return xerrors.Errorf("you must enter a object")
		}
		req, err := mustNewSignedV4Request(http.MethodPost, "http://"+addr+"/unpin?bucket="+bucket+"&object="+object, 0, nil, "s3", accessKey, secretKey)
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

var ispin = &cli.Command{
	Name:  "ispin",
	Usage: "whether the block is pinned\n\n eg.go run main.go ispin --addr= --access-key= --secret-key= --bucket= --object=",
	Flags: []cli.Flag{
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
			Name:  "bucket",
			Usage: "the bucket cannot be empty",
		},
		&cli.StringFlag{
			Name:  "object",
			Usage: "the object cannot be empty",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, accessKey, secretKey, bucket, object string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			return xerrors.Errorf("you must enter a addr")
		}
		if cctx.String("access-key") != "" {
			accessKey = cctx.String("access-key")
		} else {
			return xerrors.Errorf("you must enter a access-key")
		}
		if cctx.String("secret-key") != "" {
			secretKey = cctx.String("secret-key")
		} else {
			return xerrors.Errorf("you must enter a secret-key")
		}
		if cctx.String("bucket") != "" {
			bucket = cctx.String("bucket")
		} else {
			return xerrors.Errorf("you must enter a bucket")
		}
		if cctx.String("object") != "" {
			object = cctx.String("object")
		} else {
			return xerrors.Errorf("you must enter a object")
		}
		req, err := mustNewSignedV4Request(http.MethodPost, "http://"+addr+"/ispin?bucket="+bucket+"&object="+object, 0, nil, "s3", accessKey, secretKey)
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

//mustNewSignedV4Request  NewSignedV4Request
func mustNewSignedV4Request(method string, urlStr string, contentLength int64, body io.ReadSeeker, st string, accessKey, secretKey string) (*http.Request, error) {
	log.Infof(urlStr)
	req, err := utils.NewRequest(method, urlStr, contentLength, body)
	if err != nil {
		return nil, err
	}
	cred := &auth.Credentials{AccessKey: accessKey, SecretKey: secretKey}
	if err = utils.SignRequestV4(req, cred.AccessKey, cred.SecretKey, utils.ServiceType(st)); err != nil {
		return nil, err
	}
	return req, nil
}
