package main

import (
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/dag/proto"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"io/ioutil"
)

var log = logging.Logger("pool-client")
var addBlock = &cli.Command{
	Name:  "addblock",
	Usage: "add a block to dagpool eg.dagpool-client addblock --addr=127.0.0.1:50001 --client-user=dagpool --client-pass=dagpool --filepath=file.txt",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:  "client-user",
			Usage: "the client user ",
		},
		&cli.StringFlag{
			Name:  "client-pass",
			Usage: "the client pass ",
		},
		&cli.StringFlag{
			Name:  "filepath",
			Usage: "the block path that you want add,size is usually 1m",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, clientuser, clientpass, filepath string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
		}
		if cctx.String("client-user") != "" {
			clientuser = cctx.String("client-user")
		} else {
			log.Errorf("you must give the client user")
			return xerrors.Errorf("you must give the client user")
		}
		if cctx.String("client-pass") != "" {
			clientpass = cctx.String("client-pass")
		} else {
			log.Errorf("you must give the client pass")
			return xerrors.Errorf("you must give the client pass")
		}
		if cctx.String("filepath") != "" {
			filepath = cctx.String("filepath")
		} else {
			log.Errorf("you must give the filepath")
			return xerrors.Errorf("you must give the filepath")
		}
		poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		f, err := ioutil.ReadFile(filepath)
		add, err := poolClient.DPClient.Add(cctx.Context, &proto.AddReq{
			Block: f,
			User:  poolClient.User,
		})
		if err != nil {
			log.Errorf("add block err:%v", err)
			return err
		}
		log.Infof("add block succes cid:%v", add.Cid)
		return nil
	},
}
var getBlock = &cli.Command{
	Name:  "getblock",
	Usage: "get a block from dagpool eg.dagpool-client getblock --addr=127.0.0.1:50001 --client-user=dagpool --client-pass=dagpool --cid=QmZikYuqANVBRWcbb1zHAHEXzX6CsWbPz2mqRCoy92Jcge ",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:  "client-user",
			Usage: "the client user ",
		},
		&cli.StringFlag{
			Name:  "client-pass",
			Usage: "the client pass ",
		},
		&cli.StringFlag{
			Name:  "cid",
			Usage: "the block cid that you want get",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, clientuser, clientpass, cid string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
		}
		if cctx.String("client-user") != "" {
			clientuser = cctx.String("client-user")
		} else {
			log.Errorf("you must give the client user")
			return xerrors.Errorf("you must give the client user")
		}
		if cctx.String("client-pass") != "" {
			clientpass = cctx.String("client-pass")
		} else {
			log.Errorf("you must give the client pass")
			return xerrors.Errorf("you must give the client pass")
		}
		if cctx.String("cid") != "" {
			cid = cctx.String("cid")
		} else {
			log.Errorf("you must give the cid")
			return xerrors.Errorf("you must give the cid")
		}
		poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		get, err := poolClient.DPClient.Get(cctx.Context, &proto.GetReq{
			Cid:  cid,
			User: poolClient.User,
		})
		if err != nil {
			log.Errorf("get block err:%v", err)
			return err
		}
		log.Infof("get block succes block:%v", get.Block)
		return nil
	},
}
var removeBlock = &cli.Command{
	Name:  "removeblock",
	Usage: "remove a block from dagpool eg.dagpool-client removeblock --addr=127.0.0.1:50001 --client-user=dagpool --client-pass=dagpool --cid=QmZikYuqANVBRWcbb1zHAHEXzX6CsWbPz2mqRCoy92Jcge",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:  "client-user",
			Usage: "the client user ",
		},
		&cli.StringFlag{
			Name:  "client-pass",
			Usage: "the client pass ",
		},
		&cli.StringFlag{
			Name:  "cid",
			Usage: "the block cid that you want remove",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, clientuser, clientpass, cid string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
		}
		if cctx.String("client-user") != "" {
			clientuser = cctx.String("client-user")
		} else {
			log.Errorf("you must give the client user")
			return xerrors.Errorf("you must give the client user")
		}
		if cctx.String("client-pass") != "" {
			clientpass = cctx.String("client-pass")
		} else {
			log.Errorf("you must give the client pass")
			return xerrors.Errorf("you must give the client pass")
		}
		if cctx.String("cid") != "" {
			cid = cctx.String("cid")
		} else {
			log.Errorf("you must give the cid")
			return xerrors.Errorf("you must give the cid")
		}
		poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		re, err := poolClient.DPClient.Remove(cctx.Context, &proto.RemoveReq{
			Cid:  cid,
			User: poolClient.User,
		})
		if err != nil {
			log.Errorf("remove block err:%v", err)
			return err
		}
		log.Infof("remove block succes:%v", re.Message)
		return nil
	},
}

var addUser = &cli.Command{
	Name:  "adduser",
	Usage: "add a user to dagpool eg.dagpool-client adduser --addr=127.0.0.1:50001 --client-user=dagpool --client-pass=dagpool --username=wpg --pass=wpg12345 --capacity=1000 --policy=only-read",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:  "client-user",
			Usage: "the client user ",
		},
		&cli.StringFlag{
			Name:  "client-pass",
			Usage: "the client pass ",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the username ",
		},
		&cli.StringFlag{
			Name:  "pass",
			Usage: "the password ",
		},
		&cli.StringFlag{
			Name:  "capacity",
			Usage: "the capacity ",
		},
		&cli.StringFlag{
			Name:  "policy",
			Usage: "the policy ",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, clientuser, clientpass, username, pass, policy string
		var capacity uint64
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
		}
		if cctx.String("client-user") != "" {
			clientuser = cctx.String("client-user")
		} else {
			log.Errorf("you must give the client user")
			return xerrors.Errorf("you must give the client user")
		}
		if cctx.String("client-pass") != "" {
			clientpass = cctx.String("client-pass")
		} else {
			log.Errorf("you must give the client pass")
			return xerrors.Errorf("you must give the client pass")
		}
		if cctx.String("username") != "" {
			username = cctx.String("username")
		} else {
			log.Errorf("you must give the username")
			return xerrors.Errorf("you must give the username")
		}
		if cctx.String("pass") != "" {
			pass = cctx.String("pass")
		} else {
			log.Errorf("you must give the pass")
			return xerrors.Errorf("you must give the pass")
		}
		if cctx.String("capacity") != "" {
			capacity = cctx.Uint64("capacity")
		} else {
			log.Errorf("you must give the capacity")
			return xerrors.Errorf("you must give the capacity")
		}
		if cctx.String("policy") != "" {
			policy = cctx.String("policy")
		} else {
			log.Errorf("you must give the policy")
			return xerrors.Errorf("you must give the policy")
		}
		poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		add, err := poolClient.DPClient.AddUser(cctx.Context, &proto.AddUserReq{
			Username: username,
			Password: pass,
			Capacity: capacity,
			Policy:   policy,
			User:     poolClient.User,
		})
		if err != nil {
			log.Errorf("add user err:%v", err)
			return err
		}
		log.Infof("add user succes %v", add.Message)
		return nil
	},
}
var removeUser = &cli.Command{
	Name:  "remove-user",
	Usage: "remove a user from dagpool eg.dagpool-client remove-user --addr=127.0.0.1:50001 --client-user=dagpool --client-pass=dagpool --username=wpg --pass=wpg12345",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:  "client-user",
			Usage: "the client user ",
		},
		&cli.StringFlag{
			Name:  "client-pass",
			Usage: "the client pass ",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the username ",
		},
		&cli.StringFlag{
			Name:  "pass",
			Usage: "the password ",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, clientuser, clientpass, username, pass string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
		}
		if cctx.String("client-user") != "" {
			clientuser = cctx.String("client-user")
		} else {
			log.Errorf("you must give the client user")
			return xerrors.Errorf("you must give the client user")
		}
		if cctx.String("client-pass") != "" {
			clientpass = cctx.String("client-pass")
		} else {
			log.Errorf("you must give the client pass")
			return xerrors.Errorf("you must give the client pass")
		}
		if cctx.String("username") != "" {
			username = cctx.String("username")
		} else {
			log.Errorf("you must give the username")
			return xerrors.Errorf("you must give the username")
		}
		if cctx.String("pass") != "" {
			pass = cctx.String("pass")
		} else {
			log.Errorf("you must give the pass")
			return xerrors.Errorf("you must give the pass")
		}

		poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		re, err := poolClient.DPClient.RemoveUser(cctx.Context, &proto.RemoveUserReq{
			Username: username,
			Password: pass,
		})
		if err != nil {
			log.Errorf("remove user err:%v", err)
			return err
		}
		log.Infof("remove user succes %v", re.Message)
		return nil
	},
}
var getUser = &cli.Command{
	Name:  "get-user",
	Usage: "get a user from dagpool eg.dagpool-client get-user --addr=127.0.0.1:50001 --client-user=dagpool --client-pass=dagpool --username=wpg --pass=wpg12345",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:  "client-user",
			Usage: "the client user ",
		},
		&cli.StringFlag{
			Name:  "client-pass",
			Usage: "the client pass ",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the username ",
		},
		&cli.StringFlag{
			Name:  "pass",
			Usage: "the password ",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, clientuser, clientpass, username, pass string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
		}
		if cctx.String("client-user") != "" {
			clientuser = cctx.String("client-user")
		} else {
			log.Errorf("you must give the client user")
			return xerrors.Errorf("you must give the client user")
		}
		if cctx.String("client-pass") != "" {
			clientpass = cctx.String("client-pass")
		} else {
			log.Errorf("you must give the client pass")
			return xerrors.Errorf("you must give the client pass")
		}
		if cctx.String("username") != "" {
			username = cctx.String("username")
		} else {
			log.Errorf("you must give the username")
			return xerrors.Errorf("you must give the username")
		}
		if cctx.String("pass") != "" {
			pass = cctx.String("pass")
		} else {
			log.Errorf("you must give the pass")
			return xerrors.Errorf("you must give the pass")
		}

		poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		re, err := poolClient.DPClient.QueryUser(cctx.Context, &proto.QueryUserReq{
			Username: username,
			Password: pass,
		})
		if err != nil {
			log.Errorf("get user err:%v", err)
			return err
		}
		log.Infof("get user: %v succes,policy:%v,capacity:%v", re.Username, re.Policy, re.Capacity)
		return nil
	},
}
var updateUser = &cli.Command{
	Name:  "update-user",
	Usage: "update the user eg.dagpool-client update-user --addr=127.0.0.1:50001 --client-user=dagpool --client-pass=dagpool --username=wpg --pass=wpg12345 --capacity=1000 --policy=only-read",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:  "client-user",
			Usage: "the client user ",
		},
		&cli.StringFlag{
			Name:  "client-pass",
			Usage: "the client pass ",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "the username ",
		},
		&cli.StringFlag{
			Name:  "pass",
			Usage: "the password ",
		},
		&cli.StringFlag{
			Name:  "capacity",
			Usage: "the capacity ",
		},
		&cli.StringFlag{
			Name:  "policy",
			Usage: "the policy ",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, clientuser, clientpass, username, pass, policy string
		var capacity uint64
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
		}
		if cctx.String("client-user") != "" {
			clientuser = cctx.String("client-user")
		} else {
			log.Errorf("you must give the client user")
			return xerrors.Errorf("you must give the client user")
		}
		if cctx.String("client-pass") != "" {
			clientpass = cctx.String("client-pass")
		} else {
			log.Errorf("you must give the client pass")
			return xerrors.Errorf("you must give the client pass")
		}
		if cctx.String("username") != "" {
			username = cctx.String("username")
		} else {
			log.Errorf("you must give the username")
			return xerrors.Errorf("you must give the username")
		}
		if cctx.String("pass") != "" {
			pass = cctx.String("pass")
		} else {
			log.Errorf("you must give the pass")
			return xerrors.Errorf("you must give the pass")
		}
		if cctx.String("capacity") != "" {
			capacity = cctx.Uint64("capacity")
		} else {
			log.Errorf("you must give the capacity")
			return xerrors.Errorf("you must give the capacity")
		}
		if cctx.String("policy") != "" {
			policy = cctx.String("policy")
		} else {
			log.Errorf("you must give the policy")
			return xerrors.Errorf("you must give the policy")
		}
		poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		re, err := poolClient.DPClient.UpdateUser(cctx.Context, &proto.UpdateUserReq{
			NewUsername: username,
			NewPassword: pass,
			Capacity:    capacity,
			Policy:      policy,
			User:        poolClient.User,
		})
		if err != nil {
			log.Errorf("update user err:%v", err)
			return err
		}
		log.Infof("update user succes %v", re.Message)
		return nil
	},
}
