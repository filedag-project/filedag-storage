package node

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/filedag-project/filedag-storage/dag/proto"
	logging "github.com/ipfs/go-log/v2"
	"google.golang.org/grpc"
	"strconv"
)

var log = logging.Logger("dag-node")

// RepairDisk prepare disk repair
func (d *DagNode) RepairDisk(ip, port string) error {
	ctx := context.TODO()
	keyCodeMap, err := d.db.ReadAll("")
	if err != nil {
		return err
	}
	index := -1
	dataNode := new(DataNode)
	for i, node := range d.Nodes {
		if node.Ip == ip && node.Port == port {
			dataNode = &d.Nodes[i]
			index = i
		}
	}
	if index == -1 {
		return errors.New("the host does not exist")
	}
	for key, value := range keyCodeMap {
		keyCode := sha256String(key)
		size, err := dataNode.Client.Size(ctx, &proto.SizeRequest{
			Key: keyCode,
		})
		if err != nil {
			return err
		}
		if size.Size > 0 {
			continue
		}
		merged := make([][]byte, 0)
		for i, node := range d.Nodes {
			if i == index {
				merged = append(merged, nil)
				continue
			}
			res, err := node.Client.Get(ctx, &proto.GetRequest{Key: keyCode})
			if err != nil {
				log.Infof("this node err : %s,: %s", node.Ip, node.Port)
				return err
			}
			if len(res.DataBlock) == 0 {
				log.Infof("There is no data in this node")
				merged = append(merged, nil)
				continue
			}
			merged = append(merged, res.DataBlock)
		}
		i64, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			log.Errorf("strconv fail :%v", err)
		}
		enc, err := NewErasure(d.dataBlocks, d.parityBlocks, i64)
		if err != nil {
			log.Errorf("new erasure fail :%v", err)
			return err
		}
		err = enc.DecodeDataBlocks(merged)
		if err != nil {
			log.Errorf("decode date blocks fail :%v", err)
			return err
		}
		dataByte := merged[index]
		_, err = dataNode.Client.Put(ctx, &proto.AddRequest{Key: keyCode, DataBlock: dataByte})
		if err != nil {
			log.Errorf("data node put fail :%v", err)
			return err
		}
	}
	return err
}

// RepairHost prepare host repair
func (d *DagNode) RepairHost(oldIp, newIp, oldPort, newPort string) error {
	ctx := context.TODO()
	index, err := d.modifyConfig(oldIp, newIp, oldPort, newPort)
	if err != nil {
		log.Errorf("modify node config fail")
		return err
	}
	keyCodeMap, err := d.db.ReadAll("")
	if err != nil {
		return err
	}
	newDataNode := d.Nodes[index]
	for key, value := range keyCodeMap {
		keyCode := sha256String(key)
		merged := make([][]byte, len(d.Nodes))
		for i, node := range d.Nodes {
			if i == index {
				merged[i] = nil
				continue
			}
			res, err := node.Client.Get(ctx, &proto.GetRequest{Key: keyCode})
			if err != nil {
				log.Infof("this node err : %s,: %s", node.Ip, node.Port)
				return err
			}
			if len(res.DataBlock) == 0 {
				log.Infof("There is no data in this node")
				merged[i] = nil
				continue
			}
			merged[i] = res.DataBlock
		}
		i64, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Errorf("strconv fail :%v", err)
		}
		enc, err := NewErasure(d.dataBlocks, d.parityBlocks, i64)
		if err != nil {
			log.Errorf("new erasure fail :%v", err)
			return err
		}
		err = enc.DecodeDataBlocks(merged)
		if err != nil {
			log.Errorf("decode date blocks fail :%v", err)
			return err
		}
		dataByte := merged[index]
		_, err = newDataNode.Client.Put(ctx, &proto.AddRequest{Key: keyCode, DataBlock: dataByte})
		if err != nil {
			log.Errorf("data node put fail :%v", err)
			return err
		}
		if err != nil {
			break
		}
	}
	return err
}

//modify node config
func (d *DagNode) modifyConfig(oldIp, newIp, oldPort, newPort string) (int, error) {
	index := -1
	for i, node := range d.Nodes {
		if node.Ip == oldIp && node.Port == oldPort {
			index = i
		}
	}
	if index == -1 {
		return index, errors.New("the old ip does not exist")
	}
	addr := flag.String("addr"+fmt.Sprint(len(d.Nodes)+1), fmt.Sprintf("%s:%s", newIp, newPort), "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		conn.Close()
		return index, err
	}
	client := proto.NewMutCaskClient(conn)
	d.Nodes[index].Client = client
	d.Nodes[index].Ip = newIp
	d.Nodes[index].Port = newPort
	return index, nil
}
