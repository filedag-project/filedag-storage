# .mc/config.json
# {
#	"version": "10",
#	"aliases": {
#		"loc": {
#			"url": "http://127.0.0.1:9985",
#			"accessKey": "filedagadmin",
#			"secretKey": "filedagadmin",
#			"api": "s3v4",
#			"path": "auto"
#			},
#		"loc1": {
#			"url": "http://127.0.0.1:9985",
#			"accessKey": "testA",
#			"secretKey": "testAtestA",
#			"api": "s3v4",
#			"path": "auto"
#			},
#		"loc2": {
#			"url": "http://127.0.0.1:9985",
#			"accessKey": "testB",
#			"secretKey": "testBtestB",
#			"api": "s3v4",
#			"path": "auto"
#			}
#}

if [ ! -x "../datanode" ]|| [ ! -x "../dagpool" ]||[ ! -x "../objectstore" ] ;then
  cd ..
  print "make project"
  make
  cd testscript || exit
  fi

function init() {
    nohup ../datanode daemon --listen=127.0.0.1:9011 --datadir=/tmp/dn-data1 >dn1log.log 2>&1 &

    nohup ../datanode daemon --listen=127.0.0.1:9012 --datadir=/tmp/dn-data2 >dn2log.log 2>&1 &

    nohup ../datanode daemon --listen=127.0.0.1:9013 --datadir=/tmp/dn-data3 >dn3log.log 2>&1 &

    nohup ../dagpool daemon --datadir=/tmp/dagpool-db --config=../conf/node_config.json  >dplog.log 2>&1 &

    nohup ../objectstore daemon --pool-addr=127.0.0.1:50001 --pool-user=dagpool --pool-password=dagpool --datadir=../store-data >objlog.log 2>&1 &

    sleep 3s

    ../iam-tools add-user --admin-access-key=filedagadmin --admin-secret-key=filedagadmin --username=testA --password=testAtestA >/dev/null 2>&1
    ../iam-tools add-user --admin-access-key=filedagadmin --admin-secret-key=filedagadmin --username=testB --password=testBtestB >/dev/null 2>&1

}
function close() {
    ps aux | grep datanode | grep -v grep | awk '{print $2}'| xargs kill -9
    ps aux | grep dagpool | grep -v grep | awk '{print $2}'| xargs kill -9
    ps aux | grep objectstore | grep -v grep | awk '{print $2}'| xargs kill -9
    rm dn1log.log dn2log.log dn3log.log dplog.log objlog.log
    rm -rf ../store-data
    rm -rf ../dp-data
}

function test_except() {
    if [ "$2" -eq "$3" ]; then
      if [ "$2" -eq 1 ];then
        print "\033[1m\033[;32m$1\033[1m\033[;33m,expect fail get fail,\033[0m\033[1m\033[;32mthe test success \033[0m"
      else
        print "\033[1m\033[;32m$1,expect success get success,\033[1m\033[;32mthe test success \033[0m"
      fi
    else
      print "\033[1m\033[;31m$1,expect $2 get $3,failed \033[0m"
    fi
}
