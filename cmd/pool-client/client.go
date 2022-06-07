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
	Usage: "add a block to dagpool eg.dagpool-client addblock --addr=127.0.0.1:50001 --username=wpg --pass=wpg123456 --filepath='file.txt' ",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
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
			Name:  "filepath",
			Usage: "the block path that you want add,size is usually 1m",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, username, pass, filepath string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
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
		if cctx.String("filepath") != "" {
			filepath = cctx.String("filepath")
		} else {
			log.Errorf("you must give the filepath")
			return xerrors.Errorf("you must give the filepath")
		}
		poolClient, err := client.NewPoolClient(addr, username, pass)
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
	Usage: "get a block from dagpool eg.dagpool-client getblock --addr=127.0.0.1:50001 --username=wpg --pass=wpg123456 --cid='QmZikYuqANVBRWcbb1zHAHEXzX6CsWbPz2mqRCoy92Jcge' ",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
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
			Name:  "cid",
			Usage: "the block cid that you want get",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, username, pass, cid string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
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
		if cctx.String("cid") != "" {
			cid = cctx.String("cid")
		} else {
			log.Errorf("you must give the cid")
			return xerrors.Errorf("you must give the cid")
		}
		poolClient, err := client.NewPoolClient(addr, username, pass)
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
	Usage: "remove a block from dagpool eg.dagpool-client removeblock --addr=127.0.0.1:50001 --username=wpg --pass=wpg123456 --cid='QmZikYuqANVBRWcbb1zHAHEXzX6CsWbPz2mqRCoy92Jcge' ",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "the addr of dagpool server eg.127.0.0.1:50001",
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
			Name:  "cid",
			Usage: "the block cid that you want remove",
		},
	},
	Action: func(cctx *cli.Context) error {
		var addr, username, pass, cid string
		if cctx.String("addr") != "" {
			addr = cctx.String("addr")
		} else {
			log.Errorf("you must give the addr")
			return xerrors.Errorf("you must give the addr")
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
		if cctx.String("cid") != "" {
			cid = cctx.String("cid")
		} else {
			log.Errorf("you must give the cid")
			return xerrors.Errorf("you must give the cid")
		}
		poolClient, err := client.NewPoolClient(addr, username, pass)
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
