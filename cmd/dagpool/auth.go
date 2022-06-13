package main

import (
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"
)

var authCmd = &cli.Command{
	Name:  "auth",
	Usage: "manage dagpool user permissions",
	Subcommands: []*cli.Command{
		addUser,
		removeUser,
		getUser,
		updateUser,
	},
}

var addUser = &cli.Command{
	Name:  "add",
	Usage: "add a user for dagpool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:    "root-user",
			Usage:   "set root user",
			EnvVars: []string{EnvRootUser},
			Value:   "dagpool",
		},
		&cli.StringFlag{
			Name:    "root-password",
			Usage:   "set root password",
			EnvVars: []string{EnvRootPassword},
			Value:   "dagpool",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "set the username ",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "set the password",
		},
		&cli.Uint64Flag{
			Name:  "capacity",
			Usage: "set the capacity",
		},
		&cli.StringFlag{
			Name:  "policy",
			Usage: "set the policy, enum: only-read, only-write, read-write",
			Value: string(userpolicy.ReadWrite),
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")
		rootUser := cctx.String("root-user")
		if rootUser == "" {
			return xerrors.New("root user is invalid")
		}
		rootPassword := cctx.String("root-password")

		username := cctx.String("username")
		if username == "" {
			return xerrors.Errorf("you must give the username")
		}
		password := cctx.String("password")
		if password == "" {
			return xerrors.Errorf("you must give the password")
		}
		capacity := cctx.Uint64("capacity")

		policy := cctx.String("policy")
		if !userpolicy.CheckValid(policy) {
			return xerrors.Errorf("the policy is invalid")
		}
		poolClient, err := client.NewPoolClient(addr, rootUser, rootPassword)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		if err = poolClient.AddUser(cctx.Context, username, password, capacity, policy); err != nil {
			log.Errorf("add user err:%v", err)
			return err
		}
		return nil
	},
}

var getUser = &cli.Command{
	Name:  "get",
	Usage: "get the user config from dagpool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:    "root-user",
			Usage:   "set root user",
			EnvVars: []string{EnvRootUser},
			Value:   "dagpool",
		},
		&cli.StringFlag{
			Name:    "root-password",
			Usage:   "set root password",
			EnvVars: []string{EnvRootPassword},
			Value:   "dagpool",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "set the username to query",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")
		rootUser := cctx.String("root-user")
		if rootUser == "" {
			err := xerrors.New("root user is invalid")
			fmt.Println(err)
			return err
		}
		rootPassword := cctx.String("root-password")

		username := cctx.String("username")
		if username == "" {
			err := xerrors.Errorf("you must give the username")
			fmt.Println(err)
			return err
		}

		poolClient, err := client.NewPoolClient(addr, rootUser, rootPassword)
		if err != nil {
			fmt.Println("NewPoolClient err:", err)
			return err
		}
		reply, err := poolClient.QueryUser(cctx.Context, username)
		if err != nil {
			fmt.Println("get user err:", err)
			return err
		}
		fmt.Printf("username:%v policy:%v capacity:%v\n", reply.Username, reply.Policy, reply.Capacity)
		return nil
	},
}

var updateUser = &cli.Command{
	Name:  "update",
	Usage: "update the user config",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:    "root-user",
			Usage:   "set root user",
			EnvVars: []string{EnvRootUser},
			Value:   "dagpool",
		},
		&cli.StringFlag{
			Name:    "root-password",
			Usage:   "set root password",
			EnvVars: []string{EnvRootPassword},
			Value:   "dagpool",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "set the username ",
		},
		&cli.StringFlag{
			Name:  "new-password",
			Usage: "set the new password",
		},
		&cli.Uint64Flag{
			Name:  "new-capacity",
			Usage: "set the new capacity",
		},
		&cli.StringFlag{
			Name:  "new-policy",
			Usage: "set the new policy, enum: only-read, only-write, read-write",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")
		rootUser := cctx.String("root-user")
		if rootUser == "" {
			return xerrors.New("root user is invalid")
		}
		rootPassword := cctx.String("root-password")

		username := cctx.String("username")
		if username == "" {
			return xerrors.Errorf("you must give the username")
		}
		password := cctx.String("new-password")
		if password == "" {
			return xerrors.Errorf("you must give the password")
		}
		capacity := cctx.Uint64("new-capacity")

		policy := cctx.String("new-policy")
		if !userpolicy.CheckValid(policy) {
			return xerrors.Errorf("the policy is invalid")
		}

		poolClient, err := client.NewPoolClient(addr, rootUser, rootPassword)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		if err = poolClient.UpdateUser(cctx.Context, username, password, capacity, policy); err != nil {
			log.Errorf("update user err:%v", err)
			return err
		}
		return nil
	},
}

var removeUser = &cli.Command{
	Name:  "remove",
	Usage: "remove a user from dagpool",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:    "root-user",
			Usage:   "set root user",
			EnvVars: []string{EnvRootUser},
			Value:   "dagpool",
		},
		&cli.StringFlag{
			Name:    "root-password",
			Usage:   "set root password",
			EnvVars: []string{EnvRootPassword},
			Value:   "dagpool",
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "set the username to query",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")
		rootUser := cctx.String("root-user")
		if rootUser == "" {
			return xerrors.New("root user is invalid")
		}
		rootPassword := cctx.String("root-password")

		username := cctx.String("username")
		if username == "" {
			return xerrors.Errorf("you must give the username")
		}

		poolClient, err := client.NewPoolClient(addr, rootUser, rootPassword)
		if err != nil {
			log.Errorf("NewPoolClient err:%v", err)
			return err
		}
		if err = poolClient.RemoveUser(cctx.Context, username); err != nil {
			log.Errorf("remove user err:%v", err)
			return err
		}
		return nil
	},
}
