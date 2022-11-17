package poolservice

import (
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node/dagnode"
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	"github.com/howeyc/crc16"
	"github.com/syndtr/goleveldb/leveldb"
)

const clusterSlotConfig = "cluster-slot-cfg"
const SlotPrefix = "slot/"

var ErrDagNodeAlreadyExist = errors.New("this dag node already exists")

func keyHashSlot(key string) uint16 {
	return crc16.Checksum([]byte(key), crc16.IBMTable) & 0x3FFF
}

type SlotConfig struct {
	Version  int
	SlotsMap map[string][]slotsmgr.SlotPair
}

func (d *dagPoolService) clusterInit() error {
	if err := d.loadHashSlotsConfig(); err != nil {
		return err
	}

	for _, dagNodeConfig := range d.config.Cluster {
		if err := d.AddDagNode(&dagNodeConfig); err != nil {
			return err
		}

		if pairs, ok := d.slotConfig.SlotsMap[dagNodeConfig.Name]; ok {
			dagNode := d.dagNodesMap[dagNodeConfig.Name]
			for _, pair := range pairs {
				for idx := pair.Start; idx <= pair.End; idx++ {
					if err := d.addSlot(dagNode, idx); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (d *dagPoolService) loadHashSlotsConfig() error {
	if err := d.db.Get(clusterSlotConfig, &d.slotConfig); err != nil {
		if err == leveldb.ErrNotFound {
			return errors.New("the cluster slot config can not be found and needs to be initialized")
		}
		return fmt.Errorf("load cluster slot config failed, error: %v", err)
	}
	return nil
}

func (d *dagPoolService) saveHashSlotsConfig() error {
	return d.db.Put(clusterSlotConfig, &d.slotConfig)
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
		log.Fatal(errors.New("the slot state is inconsistent"))
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

// AllocateSlotsEvenly Perform the slots allocation before starting the cluster for the first time
func AllocateSlotsEvenly(cfg config.PoolConfig) error {
	db, err := uleveldb.OpenDb(cfg.LeveldbPath)
	if err != nil {
		return err
	}
	slotConfig := SlotConfig{
		SlotsMap: make(map[string][]slotsmgr.SlotPair),
	}
	nodesNum := len(cfg.ClusterConfig.Cluster)
	piece := slotsmgr.ClusterSlots / nodesNum
	remind := slotsmgr.ClusterSlots - piece*nodesNum
	cluster := cfg.ClusterConfig.Cluster
	curIndex := 0
	for i := 0; i < nodesNum; i++ {
		curPiece := piece
		if remind > 0 {
			curPiece += 1
			remind--
		}
		slotConfig.SlotsMap[cluster[i].Name] = []slotsmgr.SlotPair{
			{
				Start: uint64(curIndex),
				End:   uint64(curIndex + curPiece - 1),
			},
		}
		curIndex += curPiece
	}
	if err = db.Put(clusterSlotConfig, &slotConfig); err != nil {
		return err
	}
	db.Close()
	return nil
}
