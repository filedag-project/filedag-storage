package poolservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/config"
	"github.com/filedag-project/filedag-storage/dag/node/dagnode"
	"github.com/filedag-project/filedag-storage/dag/proto"
	"github.com/filedag-project/filedag-storage/dag/slotsmgr"
	"github.com/ipfs/go-cid"
	format "github.com/ipfs/go-ipld-format"
	"github.com/syndtr/goleveldb/leveldb"
	"sort"
	"time"
)

var ErrDagNodeNotFound = errors.New("this dag node does not exist")
var ErrDagNodeRemove = errors.New("this dag node still has slots and cannot be removed")
var ErrClusterMigrating = errors.New("the cluster is migrating")

type MigrateSlot struct {
	From      string
	To        string
	SlotPairs []slotsmgr.SlotPair
}

// InitSlots Perform the slots allocation for the first time
func (d *dagPoolService) InitSlots() error {
	cfg, err := d.loadConfig()
	if err != nil {
		return err
	}
	if cfg.Version != 0 {
		return errors.New("it's already initialized")
	}
	d.dagNodesLock.Lock()
	defer d.dagNodesLock.Unlock()

	// calculate number of slots each dagnode
	nodesNum := len(d.dagNodesMap)
	piece := slotsmgr.ClusterSlots / nodesNum
	remind := slotsmgr.ClusterSlots - piece*nodesNum

	var nameList []string
	for name := range d.dagNodesMap {
		nameList = append(nameList, name)
	}
	sort.Strings(nameList)
	curIndex := 0
	for i := 0; i < nodesNum; i++ {
		curPiece := piece
		if remind > 0 {
			curPiece += 1
			remind--
		}

		node := d.dagNodesMap[nameList[i]]
		for start := curIndex; start <= curIndex+curPiece-1; start++ {
			node.AddSlot(uint64(start))

			if d.slots[start] != nil {
				log.Fatal("the slot state is inconsistent")
			}
			d.slots[start] = node
		}

		dagNode := config.DagNodeInfo{}
		dagNode.Config = *node.GetConfig()
		dagNode.SlotPairs = node.GetSlotPairs()
		cfg.Cluster = append(cfg.Cluster, dagNode)

		curIndex += curPiece
	}
	// save config
	cfg.Version = 1
	if err = d.saveConfig(cfg); err != nil {
		// rollback
		for i := 0; i < nodesNum; i++ {
			node := d.dagNodesMap[nameList[i]]
			slotPairs := node.GetSlotPairs()
			for _, pair := range slotPairs {
				for slot := pair.Start; slot <= pair.End; slot++ {
					node.ClearSlot(slot)
				}
			}
			if node.GetNumSlots() != 0 {
				log.Fatal("the slot state is inconsistent")
			}
		}
		var empty [slotsmgr.ClusterSlots]*dagnode.DagNode
		d.slots = empty
		return err
	}

	return nil
}

func (d *dagPoolService) AddDagNode(nodeConfig *config.DagNodeConfig) error {
	d.dagNodesLock.Lock()
	defer d.dagNodesLock.Unlock()

	if _, ok := d.dagNodesMap[nodeConfig.Name]; ok {
		return ErrDagNodeAlreadyExist
	}
	dagNode, err := dagnode.NewDagNode(*nodeConfig)
	if err != nil {
		log.Errorf("new dagnode err:%v", err)
		return err
	}
	go dagNode.RunHeartbeatCheck(d.parentCtx)
	d.dagNodesMap[nodeConfig.Name] = dagNode

	// update local config
	cfg, err := d.loadConfig()
	if err != nil {
		return err
	}
	cfg.Version += 1
	for _, node := range cfg.Cluster {
		if node.Config.Name == nodeConfig.Name {
			log.Fatal("local cluster config is illegal")
		}
	}
	cfg.Cluster = append(cfg.Cluster, config.DagNodeInfo{Config: *nodeConfig})
	if err = d.saveConfig(cfg); err != nil {
		// rollback
		delete(d.dagNodesMap, nodeConfig.Name)
		dagNode.Close()
		return err
	}

	return nil
}

func (d *dagPoolService) GetDagNode(dagNodeName string) (*config.DagNodeConfig, error) {
	d.dagNodesLock.RLock()
	defer d.dagNodesLock.RUnlock()

	if nd, ok := d.dagNodesMap[dagNodeName]; ok {
		return nd.GetConfig(), nil
	}
	return nil, ErrDagNodeNotFound
}

func (d *dagPoolService) RemoveDagNode(dagNodeName string) (*config.DagNodeConfig, error) {
	d.dagNodesLock.Lock()
	defer d.dagNodesLock.Unlock()

	if d.state != StateOk {
		return nil, ErrDagNodeRemove
	}

	if nd, ok := d.dagNodesMap[dagNodeName]; ok {
		// make sure the dagnode has no slot
		if nd.GetNumSlots() != 0 {
			return nil, ErrDagNodeRemove
		}

		// update local config
		cfg, err := d.loadConfig()
		if err != nil {
			return nil, err
		}
		cfg.Version += 1
		for i, node := range cfg.Cluster {
			if node.Config.Name == dagNodeName {
				cfg.Cluster = append(cfg.Cluster[:i], cfg.Cluster[(i+1):]...)
				break
			}
		}
		if err = d.saveConfig(cfg); err != nil {
			return nil, err
		}

		dagNodeCfg := nd.GetConfig()
		delete(d.dagNodesMap, dagNodeName)
		// close dagnode
		nd.Close()
		return dagNodeCfg, nil
	}
	return nil, ErrDagNodeNotFound
}

func (d *dagPoolService) MigrateSlots(fromDagNodeName, toDagNodeName string, pairs []slotsmgr.SlotPair) error {
	d.dagNodesLock.Lock()
	defer d.dagNodesLock.Unlock()

	if d.state == StateMigrating {
		return ErrClusterMigrating
	}

	err := d.migrateSlotsByName(fromDagNodeName, toDagNodeName, pairs)
	if err == nil {
		// start to migrate data
		select {
		case d.migratingCh <- struct{}{}:
		default:
		}
	}
	return err
}

func (d *dagPoolService) migrateSlotsByName(fromDagNodeName, toDagNodeName string, pairs []slotsmgr.SlotPair) error {
	fromNode, fromOk := d.dagNodesMap[fromDagNodeName]
	toNode, toOk := d.dagNodesMap[toDagNodeName]
	if !fromOk || !toOk {
		return ErrDagNodeNotFound
	}
	for _, pair := range pairs {
		for slot := pair.Start; slot <= pair.End; slot++ {
			if d.slots[slot] != fromNode {
				return fmt.Errorf("dagnode[%v] does not own the slot %d", fromDagNodeName, slot)
			}
		}
	}
	cfg, err := d.loadConfig()
	if err != nil {
		return err
	}

	d.migrateSlotsByNode(fromNode, toNode, pairs)

	var updateSlots []uint16
	rollback := func() {
		d.migrateSlotsByNode(toNode, fromNode, pairs)
		for _, idx := range updateSlots {
			if errR := d.slotMigrateRepo.Remove(idx); errR != nil {
				log.Warnw("slotMigrateRepo.Remove error", "slot", idx)
			}
		}
	}
	for _, pair := range pairs {
		for slot := pair.Start; slot <= pair.End; slot++ {
			if err = d.slotMigrateRepo.Set(uint16(slot), fromDagNodeName); err != nil {
				// rollback
				rollback()
				return err
			}
			updateSlots = append(updateSlots, uint16(slot))
		}
	}

	// update local config
	cfg.Version += 1
	cfg.Cluster = nil
	for _, node := range d.dagNodesMap {
		dagNode := config.DagNodeInfo{}
		dagNode.Config = *node.GetConfig()
		dagNode.SlotPairs = node.GetSlotPairs()
		cfg.Cluster = append(cfg.Cluster, dagNode)
	}
	if err = d.saveConfig(cfg); err != nil {
		// rollback
		rollback()
		return err
	}
	d.state = StateMigrating

	return nil
}

func (d *dagPoolService) migrateSlotsByNode(fromNode, toNode *dagnode.DagNode, pairs []slotsmgr.SlotPair) {
	for _, pair := range pairs {
		for slot := pair.Start; slot <= pair.End; slot++ {
			fromNode.ClearSlot(slot)
			toNode.AddSlot(slot)

			d.slots[slot] = toNode
			d.importingSlotsFrom[slot] = fromNode
		}
	}
}

func (d *dagPoolService) migrateSlotsDataTask(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case _, ok := <-d.migratingCh:
		if !ok {
			return
		}

		numSlotOk := 0
		for slot, from := range d.importingSlotsFrom {
			if from == nil {
				numSlotOk++
				continue
			}
			to := d.slots[slot]
			toName := to.GetConfig().Name

			// slot data migrate from 'from' to 'to'
			ch, err := d.slotKeyRepo.AllKeysChan(ctx, uint16(slot), "")
			if err != nil {
				log.Fatal(err)
			}
			toMigrateSlots := 0
			successMigrateSlots := 0
			for entry := range ch {
				if entry.Value == toName {
					continue
				}
				toMigrateSlots++
				blkCid, err := cid.Parse(entry.Key)
				if err != nil {
					log.Errorf("slotKeyRepo parse cid error, slot: %v cid: %v, error: %v", slot, entry.Key, err)
					continue
				}
				bk, err := from.Get(ctx, blkCid)
				if err != nil {
					if format.IsNotFound(err) {
						successMigrateSlots++
						continue
					}
					log.Errorw("migrating get block error", "from_node", from.GetConfig().Name, "slot", slot, "cid", blkCid, "err", err)
					continue
				}
				if err = to.Put(ctx, bk); err != nil {
					log.Errorw("migrating put block error", "to_node", toName, "slot", slot, "cid", entry.Key, "err", err)
					continue
				}

				if err = d.slotKeyRepo.Set(uint16(slot), entry.Key, toName); err != nil {
					log.Errorw("slotKeyRepo set key error", "to_node", toName, "slot", slot, "cid", entry.Key, "err", err)
					continue
				}
				if err = from.DeleteBlock(ctx, blkCid); err != nil {
					log.Warnw("migrating delete block error", "from_node", from.GetConfig().Name, "slot", slot, "cid", blkCid, "err", err)
				}
				successMigrateSlots++
			}
			if toMigrateSlots == successMigrateSlots {
				// all migrated
				if err = d.slotMigrateRepo.Remove(uint16(slot)); err == nil {
					d.importingSlotsFrom[slot] = nil
					numSlotOk++
				} else {
					log.Errorw("slotMigrateRepo.Remove failed", "slot", slot)
				}
			}
		}
		// is migration done?
		if numSlotOk == slotsmgr.ClusterSlots {
			if d.checkAllSlots() {
				d.state = StateOk
			} else {
				d.state = StateFail
			}
		} else {
			// try again
			time.AfterFunc(time.Minute, func() {
				d.migratingCh <- struct{}{}
			})
		}
	}
}

func (d *dagPoolService) BalanceSlots() error {
	d.dagNodesLock.Lock()
	defer d.dagNodesLock.Unlock()

	if d.state == StateMigrating {
		return ErrClusterMigrating
	}

	// calculate number of slots each dagnode
	nodesNum := len(d.dagNodesMap)
	piece := slotsmgr.ClusterSlots / nodesNum
	remind := slotsmgr.ClusterSlots - piece*nodesNum

	var nameList []string
	for name := range d.dagNodesMap {
		nameList = append(nameList, name)
	}
	sort.Strings(nameList)
	type MigrateInfo struct {
		DagNodeName string
		NumSlots    int
	}

	availableList := make([]MigrateInfo, 0)
	requireList := make([]MigrateInfo, 0)
	for i := 0; i < nodesNum; i++ {
		curPiece := piece
		if remind > 0 {
			curPiece += 1
			remind--
		}
		node := d.dagNodesMap[nameList[i]]
		numSlots := node.GetNumSlots()
		// Is it necessary to adjust?
		if curPiece == numSlots {
			continue
		}
		if curPiece > numSlots {
			availableList = append(availableList, MigrateInfo{
				DagNodeName: nameList[i],
				NumSlots:    curPiece - numSlots,
			})
		} else {
			requireList = append(requireList, MigrateInfo{
				DagNodeName: nameList[i],
				NumSlots:    numSlots - curPiece,
			})
		}
	}

	// calculate migrate slots
	var migrateSlots []*MigrateSlot
	var available MigrateInfo
	for _, require := range requireList {
		requireSlots := require.NumSlots
		for {
			if requireSlots == 0 {
				break
			}
			if available.NumSlots == 0 {
				if len(availableList) == 0 {
					log.Fatal("the slot state is inconsistent")
				}
				available = availableList[0]
				availableList = availableList[1:]
			}

			node := d.dagNodesMap[available.DagNodeName]
			pairs := node.GetSlotPairs()
			toMigrateSlots := available.NumSlots
			if toMigrateSlots > requireSlots {
				toMigrateSlots = requireSlots
			}
			remain := toMigrateSlots
			var migrateSlotPairs []slotsmgr.SlotPair
			for index := len(pairs) - 1; index >= 0; index++ {
				if remain == 0 {
					break
				}
				if remain >= int(pairs[index].Count()) {
					migrateSlotPairs = append(migrateSlotPairs, pairs[index])
					remain -= int(pairs[index].Count())
				} else {
					migrateSlotPairs = append(migrateSlotPairs, slotsmgr.SlotPair{
						Start: pairs[index].End - uint64(remain) + 1,
						End:   pairs[index].End,
					})
					remain = 0
				}
			}
			migrateSlot := MigrateSlot{
				From:      available.DagNodeName,
				To:        require.DagNodeName,
				SlotPairs: migrateSlotPairs,
			}
			migrateSlots = append(migrateSlots, &migrateSlot)
			requireSlots -= toMigrateSlots
		}
	}

	var err error
	migrated := false
	for _, migrate := range migrateSlots {
		if err = d.migrateSlotsByName(migrate.From, migrate.To, migrate.SlotPairs); err != nil {
			break
		}
		migrated = true
	}

	if migrated {
		// start to migrate data
		select {
		case d.migratingCh <- struct{}{}:
		default:
		}
	}

	return err
}

func (d *dagPoolService) Status() (*proto.StatusReply, error) {
	list := make([]*proto.DagNodeStatus, 0, len(d.dagNodesMap))
	for _, node := range d.dagNodesMap {
		pairs := node.GetSlotPairs()
		newPairs := make([]*proto.SlotPair, 0, len(pairs))
		for _, p := range pairs {
			newPairs = append(newPairs, &proto.SlotPair{Start: uint32(p.Start), End: uint32(p.End)})
		}
		cfg := node.GetConfig()
		dataNodes := make([]*proto.DataNodeInfo, 0, len(cfg.Nodes))
		for idx, nd := range cfg.Nodes {
			state := node.GetDataNodeState(idx)
			dataNodes = append(dataNodes, &proto.DataNodeInfo{
				RpcAddress: nd,
				State:      &state,
			})
		}
		st := &proto.DagNodeStatus{
			Node: &proto.DagNodeInfo{
				Name:         cfg.Name,
				Nodes:        dataNodes,
				DataBlocks:   int32(cfg.DataBlocks),
				ParityBlocks: int32(cfg.ParityBlocks),
			},
			Pairs: newPairs,
		}
		list = append(list, st)
	}
	return &proto.StatusReply{
		State:    d.state.String(),
		Statuses: list,
	}, nil
}

func (d *dagPoolService) loadConfig() (*config.ClusterConfig, error) {
	var cfg config.ClusterConfig
	if err := d.db.Get(clusterConfig, &cfg); err != nil {
		if err != leveldb.ErrNotFound {
			return nil, fmt.Errorf("load cluster config failed, error: %v", err)
		}
	}
	return &cfg, nil
}

func (d *dagPoolService) saveConfig(cfg *config.ClusterConfig) error {
	return d.db.Put(clusterConfig, cfg)
}
