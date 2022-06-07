package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/google/martian/log"
	blocks "github.com/ipfs/go-block-format"
	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
)

func main() {
	logging.SetLogLevel("*", "INFO")
	local := []*cli.Command{
		dnPut,
		dnGet,
		dnSize,
		dnDelete,
	}
	app := &cli.App{
		Name:     "datanode",
		Usage:    "send rpc request to data node",
		Version:  "0.0.1",
		Commands: local,
	}
	app.Setup()
	if err := app.Run(os.Args); err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

var dnPut = &cli.Command{
	Name:  "dnput",
	Usage: "Write data to data node eg. go run main.go dnput --addr 127.0.0.1:9010 --file ./main.go",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "set server addr",
			Value: "127.0.0.1:9010",
		},
		&cli.StringFlag{
			Name:  "file",
			Usage: "input upload file path",
			Value: "",
		},
	},
	Action: func(c *cli.Context) error {
		if c.String("file") == "" {
			log.Errorf("you must enter a file path")
			return xerrors.Errorf("you must enter a file path")
		}
		conn, err := grpc.Dial(c.String("addr"), grpc.WithInsecure())
		if err != nil {
			conn.Close()
			log.Errorf("did not connect: %v", err)
			return err
		}
		defer conn.Close()
		client := proto.NewDataNodeClient(conn)
		file, err := ioutil.ReadFile(c.String("file"))
		block := blocks.NewBlock(file)
		keyCode := sha256String(block.Cid().String())
		fmt.Println("keyCode:", keyCode)
		_, err = client.Put(context.TODO(), &proto.AddRequest{Key: keyCode, DataBlock: block.RawData()})
		if err != nil {
			log.Errorf("%s,keyCode:%s,kvdb put :%v", c.String("addr"), keyCode, err)
		}
		return err
	},
}

var dnGet = &cli.Command{
	Name:  "dnget",
	Usage: "get data to data node eg. go run main.go dnget --addr 127.0.0.1:9010 --key 8acdc995f4b9856f7e2565e9d61ae4e7e342027dd8c6ab41a27742daf672822a",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "set server addr",
			Value: "127.0.0.1:9010",
		},
		&cli.StringFlag{
			Name:  "key",
			Usage: "input key code",
			Value: "",
		},
	},
	Action: func(c *cli.Context) error {
		if c.String("key") == "" {
			log.Errorf("you must enter a key")
			return xerrors.Errorf("you must enter a key")
		}
		conn, err := grpc.Dial(c.String("addr"), grpc.WithInsecure())
		if err != nil {
			conn.Close()
			log.Errorf("did not connect: %v", err)
			return err
		}
		defer conn.Close()
		client := proto.NewDataNodeClient(conn)
		res, err := client.Get(context.TODO(), &proto.GetRequest{Key: c.String("key")})
		if err != nil {
			log.Errorf("%s,keyCode:%s,kvdb get :%v", c.String("addr"), c.String("key"), err)
			return err
		}
		fmt.Println(string(res.DataBlock))
		return nil
	},
}

var dnSize = &cli.Command{
	Name:  "dnsize",
	Usage: "get data size to data node eg. go run main.go dnsize --addr 127.0.0.1:9010 --key 8acdc995f4b9856f7e2565e9d61ae4e7e342027dd8c6ab41a27742daf672822a",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "set server addr",
			Value: "127.0.0.1:9010",
		},
		&cli.StringFlag{
			Name:  "key",
			Usage: "input key code",
			Value: "",
		},
	},
	Action: func(c *cli.Context) error {
		if c.String("key") == "" {
			log.Errorf("you must enter a key")
			return xerrors.Errorf("you must enter a key")
		}
		conn, err := grpc.Dial(c.String("addr"), grpc.WithInsecure())
		if err != nil {
			conn.Close()
			log.Errorf("did not connect: %v", err)
			return err
		}
		defer conn.Close()
		client := proto.NewDataNodeClient(conn)
		res, err := client.Size(context.TODO(), &proto.SizeRequest{Key: c.String("key")})
		if err != nil {
			log.Errorf("%s,keyCode:%s,kvdb size :%v", c.String("addr"), c.String("key"), err)
			return err
		}
		fmt.Println("size:", res.Size)
		return nil
	},
}

var dnDelete = &cli.Command{
	Name:  "dndelete",
	Usage: "remove data to data node eg. go run main.go dndelete --addr 127.0.0.1:9010 --key 8acdc995f4b9856f7e2565e9d61ae4e7e342027dd8c6ab41a27742daf672822a",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Usage: "set server addr",
			Value: "127.0.0.1:9010",
		},
		&cli.StringFlag{
			Name:  "key",
			Usage: "input key code",
			Value: "",
		},
	},
	Action: func(c *cli.Context) error {
		if c.String("key") == "" {
			log.Errorf("you must enter a key")
			return xerrors.Errorf("you must enter a key")
		}
		conn, err := grpc.Dial(c.String("addr"), grpc.WithInsecure())
		if err != nil {
			conn.Close()
			log.Errorf("did not connect: %v", err)
			return err
		}
		defer conn.Close()
		client := proto.NewDataNodeClient(conn)
		res, err := client.Delete(context.TODO(), &proto.DeleteRequest{Key: c.String("key")})
		if err != nil {
			log.Errorf("%s,keyCode:%s,kvdb delete :%v", c.String("addr"), c.String("key"), err)
			return err
		}
		fmt.Println(res.Message)
		return nil
	},
}

func sha256String(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
