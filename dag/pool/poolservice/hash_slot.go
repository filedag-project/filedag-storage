package poolservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/node/dagnode"
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
	"github.com/howeyc/crc16"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
)

const clusterConfig = "cluster-cfg"

var ErrDagNodeAlreadyExist = errors.New("this dag node already exists")

func keyHashSlot(key string) uint16 {
	return crc16.Checksum([]byte(key), crc16.IBMTable) & 0x3FFF
}

func (d *dagPoolService) clusterInit() error {
	cfg, err := d.loadConfig()
	if err != nil {
		return err
	}

	for _, dagNodeConfig := range cfg.Cluster {
		if err := d.AddDagNode(&dagNodeConfig.Config); err != nil {
			return err
		}

		dagNode := d.dagNodesMap[dagNodeConfig.Config.Name]
		for _, pair := range dagNodeConfig.SlotPairs {
			for idx := pair.Start; idx <= pair.End; idx++ {
				if err := d.addSlot(dagNode, idx); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *dagPoolService) checkAllSlots() bool {
	for _, node := range d.slots {
		if node == nil {
			return false
		}
	}
	return true
}

func (d *dagPoolService) addSlot(node *dagnode.DagNode, slot uint64) error {
	if d.slots[slot] != nil {
		return errors.New("slot already exists")
	}
	node.AddSlot(slot)
	d.slots[slot] = node
	return nil
}

func (d *dagPoolService) delSlot(slot uint64) error {
	if d.slots[slot] == nil {
		return errors.New("slot does not exist")
	}

	if !d.slots[slot].ClearSlot(slot) {
		log.Fatal("the slot state is inconsistent")
	}
	d.slots[slot] = nil
	return nil
}

// Delete all the slots associated with the specified node.
// The number of deleted slots is returned.
func (d *dagPoolService) delNodeSlots(node *dagnode.DagNode) int {
	deleted := 0

	for j := 0; j < slotsmgr.ClusterSlots; j++ {
		if node.GetSlot(uint64(j)) {
			d.delSlot(uint64(j))
			deleted++
		}
	}
	return deleted
}

// readBlock read block from dagnode
func (d *dagPoolService) readBlock(ctx context.Context, c cid.Cid) (blocks.Block, error) {
	slot := keyHashSlot(c.String())
	if node := d.importingSlotsFrom[slot]; node != nil {
		b, err := node.Get(ctx, c)
		if err == nil {
			return b, nil
		}
	}
	b, err := d.slots[slot].Get(ctx, c)
	if err != nil {
		if format.IsNotFound(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get block for %s: %v", c, err)
	}
	return b, nil
}

// readBlockSize read block size from dagnode
func (d *dagPoolService) readBlockSize(ctx context.Context, c cid.Cid) (int, error) {
	slot := keyHashSlot(c.String())
	if node := d.importingSlotsFrom[slot]; node != nil {
		size, err := node.GetSize(ctx, c)
		if err == nil {
			return size, nil
		}
	}
	return d.slots[slot].GetSize(ctx, c)
}

// putBlock put block to dagnode
func (d *dagPoolService) putBlock(ctx context.Context, block blocks.Block) error {
	blkCid := block.Cid()
	slot := keyHashSlot(blkCid.String())
	selNode := d.slots[slot]
	if err := selNode.Put(ctx, block); err != nil {
		return err
	}

	if err := d.slotKeyRepo.Set(slot, blkCid.String(), selNode.GetConfig().Name); err != nil {
		// rollback
		if delerr := selNode.DeleteBlock(ctx, blkCid); delerr != nil {
			log.Errorw("rollback block error", "slot", slot, "cid", blkCid, "error", delerr)
		}

		return err
	}
	return nil
}

// deleteBlock delete block from dagnode
func (d *dagPoolService) deleteBlock(ctx context.Context, c cid.Cid) error {
	slot := keyHashSlot(c.String())
	if err := d.slotKeyRepo.Remove(slot, c.String()); err != nil {
		return err
	}

	selNode := d.slots[slot]
	if err := selNode.DeleteBlock(ctx, c); err != nil {
		// rollback
		if dberr := d.slotKeyRepo.Set(slot, c.String(), selNode.GetConfig().Name); dberr != nil {
			log.Errorw("rollback slot entry error", "slot", slot, "cid", c, "error", dberr)
			return nil
		}
		return err
	}
	return nil
}
