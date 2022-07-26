package poolservice

//
//import (
//	"bytes"
//	"context"
//	"github.com/filedag-project/filedag-storage/dag/pool/client"
//	"github.com/ipfs/go-blockservice"
//	"github.com/ipfs/go-cid"
//	offline "github.com/ipfs/go-ipfs-exchange-offline"
//	logging "github.com/ipfs/go-log/v2"
//	"github.com/ipfs/go-merkledag"
//	"io/ioutil"
//	"testing"
//	"time"
//)
//
////func TestDagPoolService_Gc(t *testing.T) {
////	service, err := NewDagPoolService(config.PoolConfig{LeveldbPath: utils.TmpDirPath(t)})
////	if err != nil {
////		return
////	}
////	ctx := context.Background()
////
////	errc := make(chan error)
////	go func() {
////		errc <- service.Gc(ctx)
////		close(errc)
////	}()
////
////	go aaa(service)
////	time.Sleep(time.Second * 30)
////	//fmt.Println(<-errc)
////
////}
////
////func aaa(service *dagPoolService) {
////	//cancelFunc()
////	for {
////		fmt.Println("aaa")
////		service.gcl.Lock()
////		time.Sleep(time.Second)
////		service.gcl.Unlock()
////		time.Sleep(time.Second)
////	}
////}
////func BenchmarkXXX(b *testing.B) {
////	b.StopTimer()
////	logging.SetAllLoggers(logging.LevelInfo)
////	poolClient, err := client.NewPoolClient("127.0.0.1:50001", "dagpool", "dagpool")
////	if err != nil {
////		log.Errorf("NewPoolClient err:%v", err)
////		return
////	}
////	_, err = poolClient.DPClient.AddUser(context.TODO(), &proto.AddUserReq{
////		Username: "aaa",
////		Password: "aaa12345",
////		Capacity: 1000,
////		Policy:   "read-write",
////		User:     poolClient.User,
////	})
////	poolClient2, err := client.NewPoolClient("127.0.0.1:50001", "aaa", "aaa12345")
////	if err != nil {
////		log.Errorf("NewPoolClient err:%v", err)
////		return
////	}
////	cidBuilder, _ := merkledag.PrefixForCidVersion(0)
////	f, err := ioutil.ReadFile("gc.go")
////	b.StartTimer()
////	for i := 0; i < b.N; i++ {
////		node, err := client.BalanceNode(bytes.NewReader(append(f, byte(i))), poolClient, cidBuilder)
////		if err != nil {
////			log.Errorf("add block err:%v", err)
////			return
////		}
////		log.Infof("add block succes cid:%v", node.Cid())
////		node2, err := client.BalanceNode(bytes.NewReader(append(f, byte(i))), poolClient2, cidBuilder)
////		if err != nil {
////			log.Errorf("add block err:%v", err)
////			return
////		}
////		log.Infof("add block succes cid:%v", node2.Cid())
////		//time.Sleep(time.Millisecond * 50)
////	}
////}
//func Test_aa(t *testing.T) {
//	logging.SetAllLoggers(logging.LevelDebug)
//	//poolClient, err := client.NewPoolClient("192.168.1.159:50001", "dagpool", "dagpool")
//	//if err != nil {
//	//	log.Errorf("NewPoolClient err:%v", err)
//	//	return
//	//}
//	//ad, err := poolClient.DPClient.AddUser(context.TODO(), &proto.AddUserReq{
//	//	Username: "aaa",
//	//	Password: "aaa12345",
//	//	Capacity: 1000,
//	//	Policy:   "read-write",
//	//	User:     poolClient.User,
//	//})
//	//if err != nil {
//	//	log.Errorf("add user err:%v", err)
//	//	return
//	//}
//	//log.Infof("add user succes %v", ad.Message)
//	//
//	//go adder("192.168.1.159:50001", "dagpool", "dagpool")
//	go adder("192.168.1.159:50001", "aaa", "aaa12345")
//	//time.Sleep(time.Minute * 2)
//	//c, _ := cid.Decode("Qmedhji2WQautBzWA6PyizUNQD35BxbPAowMpL2tqC3x3t")
//	//for i := 0; i < 10; i++ {
//	//	poolClient.Remove(context.TODO(), c)
//	//}
//	time.Sleep(time.Minute * 5)
//	//poolClient.Get(context.TODO(), c)
//}
//func adder(addr, clientuser, clientpass string) {
//	poolClient, err := client.NewPoolClient(addr, clientuser, clientpass)
//	if err != nil {
//		log.Errorf("NewPoolClient err:%v", err)
//		return
//	}
//	dagServ := merkledag.NewDAGService(blockservice.New(poolClient, offline.Exchange(poolClient)))
//
//	f, err := ioutil.ReadFile("/Users/wpg/Downloads/IPFS-Desktop-0.21.0.dmg")
//	c := cid.Cid{}
//	cidBuilder, _ := merkledag.PrefixForCidVersion(0)
//	for i := 0; i < 1000; i++ {
//		node, err := client.BalanceNode(bytes.NewReader(append(f, byte(i))), dagServ, cidBuilder)
//		if err != nil {
//			log.Errorf("add block err:%v", err)
//			return
//		}
//		log.Infof("add block succes cid:%v", node.Cid())
//		c = node.Cid()
//		time.Sleep(time.Second * 50)
//	}
//	time.Sleep(time.Minute * 2)
//	get, err := poolClient.Get(context.TODO(), c)
//	if err != nil {
//		log.Errorf("%v,err%v", clientuser, err)
//		return
//	}
//	log.Infof(get.Cid().String())
//}
