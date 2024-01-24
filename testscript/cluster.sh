function init() {
    nohup ./datanode daemon --listen=127.0.0.1:9011 --datadir=./dn-data1 >dn1log.log 2>&1 &
    nohup ./datanode daemon --listen=127.0.0.1:9012 --datadir=./dn-data2 >dn2log.log 2>&1 &
    nohup ./datanode daemon --listen=127.0.0.1:9013 --datadir=./dn-data3 >dn3log.log 2>&1 &

    nohup ./datanode daemon --listen=127.0.0.1:9014 --datadir=./dn-data4 > dn4.log 2>&1 &
    nohup ./datanode daemon --listen=127.0.0.1:9015 --datadir=./dn-data5 > dn5.log 2>&1 &
    nohup ./datanode daemon --listen=127.0.0.1:9016 --datadir=./dn-data6 > dn6.log 2>&1 &

    nohup ./datanode daemon --listen=127.0.0.1:9017 --datadir=./dn-data7 > dn7.log 2>&1 &
    nohup ./datanode daemon --listen=127.0.0.1:9018 --datadir=./dn-data8 > dn8.log 2>&1 &
    nohup ./datanode daemon --listen=127.0.0.1:9019 --datadir=./dn-data9 > dn9.log 2>&1 &
    echo "start datanode success"
    sleep 3
    nohup ./dagpool daemon --datadir=./dp-data > dp.log 2>&1 &
    echo "start dagpool success"
    sleep 3
    nohup ./objectstore daemon --pool-addr=127.0.0.1:50001 --pool-user=dagpool --pool-password=dagpool --datadir=./store-data >obj.log 2>&1 &
    echo "start objectstore success"
    sleep 3
  ./dagpool cluster add conf/node_config.json conf/node_config2.json conf/node_config3.json
   echo "add  datanode to dagpool cluster success"
}

#health report of storage nodes and status
function statusAndBalance() {
    ./dagpool cluster status
    echo "*****Balance*****"
    ./dagpool cluster balance
    ./dagpool cluster status
}

#expansion of storage nodes
function expansion() {
   echo "2.expansion"
   nohup ./datanode daemon --listen=127.0.0.1:9020 --datadir=./dn-data10 > dn10.log 2>&1 &
   nohup ./datanode daemon --listen=127.0.0.1:9021 --datadir=./dn-data11 > dn11.log 2>&1 &
   nohup ./datanode daemon --listen=127.0.0.1:9022 --datadir=./dn-data12 > dn12.log 2>&1 &
   ./dagpool cluster add conf/node_config4.json
   ./dagpool cluster status
}
#scaling of storage nodes
function scaling() {
  echo "3.scaling"
   ./dagpool cluster migrate dag_node3 dag_node2 10923-16383
   sleep 3
   ./dagpool cluster remove dag_node3
   ./dagpool cluster status
}
function clean() {
    rm -rf ./dn-data*
    rm -rf ./dp-data*
    rm -f ./dn*.log
    rm -rf ./dp-db
    rm -f ./dp.log
    rm -rf ./store-data
    rm -f ./obj.log
    ps -ef | grep datanode | awk '{print $2}' | xargs kill
    ps -ef | grep dagpool | awk '{print $2}' | xargs kill

    ps -ef | grep objectstore | awk '{print $2}' | xargs kill
    }
# 根据参数选择执行对应的函数
case "$1" in
    "init")
        init
        ;;
      "clean")
        clean
        ;;
    "statusAndBalance")
        statusAndBalance
        ;;
    "expansion")
        expansion
        ;;
    "scaling")
        scaling
        ;;
    *)
        echo "Invalid argument. Usage: $0 {init|statusAndBalance|expansion|scaling}"
        exit 1
        ;;
esac

exit 0
