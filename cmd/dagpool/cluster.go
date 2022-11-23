package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
	"github.com/filedag-project/filedag-storage/dag/utils"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"strconv"
	"strings"
)

var clusterCmd = &cli.Command{
	Name:  "cluster",
	Usage: "Manage dagpool cluster nodes",
	Subcommands: []*cli.Command{
		status,
		addDagNode,
		getDagNode,
		removeDagNode,
		initSlots,
		balanceSlots,
		migrateSlots,
	},
}

var status = &cli.Command{
	Name:  "status",
	Usage: "Displays the current status of the cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
		&cli.BoolFlag{
			Name:  "detail",
			Usage: "displays the detail of cluster",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")

		cli, err := client.NewPoolClusterClient(addr)
		if err != nil {
			return err
		}
		defer cli.Close(cctx.Context)
		reply, err := cli.Status(cctx.Context)
		if err != nil {
			return err
		}

		fmt.Printf("cluster_state: %s\ncluster_dagnodes: %d\ncluster_dagnodes_info:\n",
			reply.State, len(reply.Statuses))

		for _, status := range reply.Statuses {
			fmt.Printf("  dagnode_name: %s\n  dagnode_slots: %s\n",
				status.Node.Name, utils.ToSlotPairs(status.Pairs))
			if cctx.Bool("detail") {
				fmt.Printf("  dagnode_set:\n    nodes:\n")
				for idx, nd := range status.Node.Nodes {
					st := "fail"
					if nd.State != nil && *nd.State {
						st = "ok"
					}
					fmt.Printf("      set_index: %d, rpc_address: %s, state: %s\n", idx, nd.RpcAddress, st)
				}
				fmt.Printf("    data_blocks: %d\n    parity_blocks: %d\n",
					status.Node.DataBlocks, status.Node.ParityBlocks)
			}
			fmt.Println()
		}
		return nil
	},
}

var addDagNode = &cli.Command{
	Name:      "add",
	Usage:     "Add a dagnode to the dag pool cluster",
	ArgsUsage: "dagnode_config_path [...]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
		&cli.StringFlag{
			Name:  "format",
			Usage: "the format type of dagnode config file",
			Value: "json",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")
		if cctx.String("format") != "json" {
			return errors.New("not support this format")
		}
		if cctx.NArg() == 0 {
			return errors.New("at least one dagnode configuration file path is required")
		}

		cli, err := client.NewPoolClusterClient(addr)
		if err != nil {
			return err
		}
		defer cli.Close(cctx.Context)
		var dagNodes []*config.DagNodeConfig
		for i := 0; i < cctx.NArg(); i++ {
			path := cctx.Args().Get(i)
			var nc config.DagNodeConfig
			cfgBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			if err = json.Unmarshal(cfgBytes, &nc); err != nil {
				return err
			}
			dagNodes = append(dagNodes, &nc)
		}
		successCount := 0
		for _, node := range dagNodes {
			if err = cli.AddDagNode(cctx.Context, node); err != nil {
				fmt.Printf("Error: add dagnode failed, name=%s, error=%v\n", node.Name, err)
			}
			successCount++
		}
		fmt.Printf("add dagnodes:\ntotal: %v success: %v failed: %v\n",
			len(dagNodes), successCount, len(dagNodes)-successCount)
		return nil
	},
}

var getDagNode = &cli.Command{
	Name:      "get",
	Usage:     "Get a dagnode from the dag pool cluster",
	ArgsUsage: "dagnode_name",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")
		if cctx.NArg() != 1 {
			return errors.New("a dagnode name is required")
		}

		cli, err := client.NewPoolClusterClient(addr)
		if err != nil {
			return err
		}
		defer cli.Close(cctx.Context)

		dagnode, err := cli.GetDagNode(cctx.Context, cctx.Args().First())
		if err != nil {
			return err
		}

		fmt.Printf("dagnode_name: %s\n", dagnode.Name)
		fmt.Printf("dagnode_set:\n  nodes:\n")
		for idx, nd := range dagnode.Nodes {
			fmt.Printf("    set_index: %d, rpc_address: %s\n", idx, nd)
		}
		fmt.Printf("  data_blocks: %d\n  parity_blocks: %d\n",
			dagnode.DataBlocks, dagnode.ParityBlocks)
		return nil
	},
}

var removeDagNode = &cli.Command{
	Name:      "remove",
	Usage:     "Remove a dagnode from the dag pool cluster",
	ArgsUsage: "dagnode_name",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")
		if cctx.NArg() != 1 {
			return errors.New("a dagnode name is required")
		}

		cli, err := client.NewPoolClusterClient(addr)
		if err != nil {
			return err
		}
		defer cli.Close(cctx.Context)

		dagnode, err := cli.RemoveDagNode(cctx.Context, cctx.Args().First())
		if err != nil {
			return err
		}

		fmt.Printf("the dagnode is removed successfully\n")
		fmt.Printf("removed_dagnode_name: %s\n", dagnode.Name)
		fmt.Printf("removed_dagnode_set:\n  nodes:\n")
		for idx, nd := range dagnode.Nodes {
			fmt.Printf("    set_index: %d, rpc_address: %s\n", idx, nd)
		}
		fmt.Printf("  data_blocks: %d\n  parity_blocks: %d\n",
			dagnode.DataBlocks, dagnode.ParityBlocks)
		return nil
	},
}

var initSlots = &cli.Command{
	Name:  "init",
	Usage: "Init slots of the dag pool cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")

		cli, err := client.NewPoolClusterClient(addr)
		if err != nil {
			return err
		}
		defer cli.Close(cctx.Context)
		return cli.InitSlots(cctx.Context)
	},
}

var balanceSlots = &cli.Command{
	Name:  "balance",
	Usage: "Balance slots of the dag pool cluster",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")

		cli, err := client.NewPoolClusterClient(addr)
		if err != nil {
			return err
		}
		defer cli.Close(cctx.Context)
		return cli.BalanceSlots(cctx.Context)
	},
}

var migrateSlots = &cli.Command{
	Name:      "migrate",
	Usage:     "Migrate slots from a dagnode to another dagnode",
	ArgsUsage: "from_dagnode_name to_dagnode_name slots_pair(start-end) [...]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "the address of dagpool server",
			Value: "127.0.0.1:50001",
		},
	},
	Action: func(cctx *cli.Context) error {
		addr := cctx.String("address")
		if cctx.NArg() < 3 {
			return errors.New("at least input three parameters")
		}
		from := cctx.Args().Get(0)
		to := cctx.Args().Get(1)
		var slotPairs []slotsmgr.SlotPair
		for i := 2; i < cctx.NArg(); i++ {
			str := cctx.Args().Get(i)
			list := strings.Split(str, "-")
			switch len(list) {
			case 1:
				start, err := strconv.ParseUint(list[0], 10, 32)
				if err != nil {
					return err
				}
				slotPairs = append(slotPairs, slotsmgr.SlotPair{
					Start: start,
					End:   start,
				})
			case 2:
				start, err := strconv.ParseUint(list[0], 10, 32)
				if err != nil {
					return err
				}
				end, err := strconv.ParseUint(list[1], 10, 32)
				if err != nil {
					return err
				}
				if start > end {
					return errors.New("slots_pair is illegal")
				}
				slotPairs = append(slotPairs, slotsmgr.SlotPair{
					Start: start,
					End:   end,
				})
			default:
				return errors.New("slots_pair is illegal")
			}
		}

		cli, err := client.NewPoolClusterClient(addr)
		if err != nil {
			return err
		}
		defer cli.Close(cctx.Context)
		return cli.MigrateSlots(cctx.Context, from, to, slotPairs)
	},
}
