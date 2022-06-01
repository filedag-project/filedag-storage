package client

//func TestPoolClient_Add_Get(t *testing.T) {
//	//go server.StartTestDagPoolServer(t)
//	time.Sleep(time.Second * 1)
//	logging.SetLogLevel("*", "INFO")
//	r := bytes.NewReader([]byte("123456"))
//	cidBuilder, err := merkledag.PrefixForCidVersion(0)
//
//	addr := flag.String("addr", "localhost:50001", "the address to connect to")
//	Conn, err := grpc.Dial(*addr, grpc.WithInsecure())
//	if err != nil {
//		log.Errorf("did not connect: %v", err)
//	}
//	defer Conn.Close()
//	c := proto.NewDagPoolClient(Conn)
//	pc := PoolClient{c, cidBuilder, Conn}
//	var ctx = context.Background()
//	node, err := BalanceNode(ctx, r, pc, cidBuilder)
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return
//	}
//	fmt.Println("aaaaa", node.Cid().String())
//	get, err := pc.Get(ctx, node.Cid())
//	if err != nil {
//		log.Errorf("err:%v", err)
//		return
//	}
//	fmt.Println(get.String())
//}
