package store

//func TestStorageSys_Object(t *testing.T) {
//	go server.StartTestDagPoolServer(t)
//	time.Sleep(time.Second * 5)
//	var s StorageSys
//	s.DagPool, _ = dagpoolcli.NewPoolClient("localhost:9002")
//	s.Db, _ = uleveldb.OpenDb(utils.TmpDirPath(&testing.T{}))
//	r := ioutil.NopCloser(bytes.NewReader([]byte("123456")))
//	object, err := s.StoreObject(context.TODO(), "test", "testbucket", "testobject", r, 6)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("object:%v", object)
//	getObject, i, err := s.GetObject(context.TODO(), "test", "testbucket", "testobject")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(getObject)
//	all, _ := ioutil.ReadAll(i)
//	fmt.Println(string(all))
//	s.DagPool.Close(context.TODO())
//}
