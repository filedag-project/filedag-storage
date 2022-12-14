package main

import (
	"context"
	"flag"
	"fmt"
	dagpoolcli "github.com/filedag-project/filedag-storage/dag/pool/client"
	"github.com/filedag-project/filedag-storage/objectservice/iam"
	"github.com/filedag-project/filedag-storage/objectservice/iamapi"
	"github.com/filedag-project/filedag-storage/objectservice/objmetadb"
	"github.com/filedag-project/filedag-storage/objectservice/pkg/auth"
	"github.com/filedag-project/filedag-storage/objectservice/s3api"
	"github.com/filedag-project/filedag-storage/objectservice/store"
	"github.com/filedag-project/filedag-storage/objectservice/utils"
	"github.com/filedag-project/filedag-storage/objectservice/utils/httpstats"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-blockservice"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	"github.com/ipfs/go-merkledag"
	"log"
	"net/http"
	"os"
)

//go run -tags example main.go daemon --datadir=/tmp/leveldb2/fds.db --listen=:9985 --pool-addr=127.0.0.1:50001 --pool-user=dagpool  --pool-password=dagpool --root-user=filedagadmin root-password=filedagadmin
func main() {
	var datadir, listen, poolAddr, poolUser, poolPass, rootUser, rootPass string
	f := flag.NewFlagSet("daemon", flag.ExitOnError)
	f.StringVar(&datadir, "datadir", "/tmp/leveldb2/fds.db", "directory to store data in")
	f.StringVar(&listen, "listen", ":9985", "set server listen")
	f.StringVar(&poolAddr, "pool-addr", "localhost:50001", "set the pool rpc address you want connect")
	f.StringVar(&poolUser, "pool-user", "", "set pool user ")
	f.StringVar(&poolPass, "pool-password", "", "set pool password")
	f.StringVar(&rootUser, "root-user", "filedagadmin", "set root filedag root user")
	f.StringVar(&rootPass, "root-password", "filedagadmin", "set root filedag root password")

	switch os.Args[1] {
	case "daemon":
		f.Parse(os.Args[2:])
		if poolUser == "" || poolPass == "" {
			fmt.Printf("db-path:%v, port:%v, pool-addr:%v, pool-user:%v, pool-user-pass:%v", datadir, listen, poolAddr, poolUser, poolPass)
			fmt.Println("please check your input\n " +
				"USAGE ERROR: go daemon -tags example main.go run daemon --datadir=/tmp/leveldb2/fds.db --listen=:9985 --pool-addr=127.0.0.1:50001 --pool-user=dagpool  --pool-password=dagpool --root-user=filedagadmin root-password=filedagadmin")
		} else {
			run(datadir, listen, poolAddr, poolUser, poolPass)
		}
	default:
		fmt.Println("expected 'daemon' subcommands")
		os.Exit(1)
	}
}
func run(leveldbPath, port, poolAddr, poolUser, poolPass string) {
	db, err := objmetadb.OpenDb(leveldbPath)
	if err != nil {
		fmt.Printf("OpenDb err:%v", err)
		return
	}
	defer db.Close()
	cred, err := auth.CreateCredentials(auth.DefaultAccessKey, auth.DefaultSecretKey)
	if err != nil {
		println(err)
		return
	}
	authSys := iam.NewAuthSys(db, cred)
	router := mux.NewRouter()
	poolClient, err := dagpoolcli.NewPoolClient(poolAddr, poolUser, poolPass, true)
	if err != nil {
		log.Fatalf("connect dagpool server err: %v", err)
	}
	defer poolClient.Close(context.TODO())
	dagServ := merkledag.NewDAGService(blockservice.New(poolClient, offline.Exchange(poolClient)))
	storageSys := store.NewStorageSys(context.TODO(), dagServ, db)
	bmSys := store.NewBucketMetadataSys(db)
	storageSys.SetNewBucketNSLock(bmSys.NewNSLock)
	storageSys.SetHasBucket(bmSys.HasBucket)
	bmSys.SetEmptyBucket(storageSys.EmptyBucket)
	cleanData := func(accessKey string) {
		ctx := context.Background()
		bkts, err := bmSys.GetAllBucketsOfUser(ctx, accessKey)
		if err != nil {
			log.Printf("GetAllBucketsOfUser error: %v", err)
		}
		for _, bkt := range bkts {
			if err = storageSys.CleanObjectsInBucket(ctx, bkt.Name); err != nil {
				log.Printf("CleanObjectsInBucket error: %v", err)
				continue
			}
			if err = bmSys.DeleteBucket(ctx, bkt.Name); err != nil {
				log.Printf("DeleteBucket error: %v", err)
			}
		}
	}
	bucketInfoFunc := func(ctx context.Context, accessKey string) []store.BucketInfo {
		var bucketInfos []store.BucketInfo
		bkts, err := bmSys.GetAllBucketsOfUser(ctx, accessKey)
		if err != nil {
			fmt.Printf("GetAllBucketsOfUser error: %v\n", err)
			return bucketInfos
		}
		for _, bkt := range bkts {
			info, err := storageSys.GetBucketInfo(ctx, bkt.Name)
			if err != nil {
				return nil
			}
			bucketInfos = append(bucketInfos, info)
		}
		return bucketInfos
	}
	storePoolStatsFunc := func(ctx context.Context) (store.DataUsageInfo, error) {
		return storageSys.StoreStats(ctx)
	}
	iamapi.NewIamApiServer(router, authSys, httpstats.NewHttpStatsSys(db), cleanData, bucketInfoFunc, storePoolStatsFunc)
	s3api.NewS3Server(router, authSys, bmSys, storageSys, httpstats.NewHttpStatsSys(db))

	for _, ip := range utils.MustGetLocalIP4().ToSlice() {
		fmt.Printf("start sever at http://%v%v", ip, port)
	}
	err = http.ListenAndServe(port, router)
	if err != nil {
		fmt.Printf("Listen And Serve err%v", err)
		return
	}
}
