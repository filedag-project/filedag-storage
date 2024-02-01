#!/bin/bash

function green_echo() {
    local green=$(tput setaf 2)
    local reset=$(tput sgr0)
    echo "${green}$1${reset}"
}

function start_datanode() {
    local start_port=$1
    local end_port=$2
    local datadir_prefix=$3

    for port in $(seq $start_port $end_port); do
        nohup ./datanode daemon --listen=127.0.0.1:$port --datadir=./$datadir_prefix-data$port > dn$port.log 2>&1 &
    done
    green_echo "Start datanodes success"
    sleep 3
}

function start_dagpool() {
    nohup ./dagpool daemon --datadir=./dp-data > dp.log 2>&1 &
    green_echo "Start dagpool success"
    sleep 3
}

function start_objectstore() {
    nohup ./objectstore daemon --pool-addr=127.0.0.1:50001 --pool-user=dagpool --pool-password=dagpool --datadir=./store-data > obj.log 2>&1 &
    green_echo "Start objectstore success"
    sleep 3
}

function add_datanodes_to_dagpool() {
    ./dagpool cluster add conf/node_config.json conf/node_config2.json conf/node_config3.json > /dev/null 2>&1
    green_echo "Add datanodes to dagpool cluster success"
}
function init() {
    start_datanode 9011 9019 dn
    start_dagpool
    start_objectstore
    add_datanodes_to_dagpool
}
function status_and_balance() {
    ./dagpool cluster status
    echo ""
    echo "*****Balance*****"
    ./dagpool cluster balance
    result
}

function expansion() {
    echo ""
    echo "*****Expansion*****"
    start_datanode 9020 9022 dn
    ./dagpool cluster add conf/node_config4.json > /dev/null 2>&1
    result
}

function scaling() {
    echo ""
    echo "*****Scaling*****"
    ./dagpool cluster migrate dag_node3 dag_node2 10923-16383 > /dev/null 2>&1
    sleep 3
    ./dagpool cluster remove dag_node3 > /dev/null 2>&1
    result
}

function clean() {
    rm -rf ./dn-data* ./dp-data* ./dp-db ./store-data
    rm -f ./dn*.log ./dp.log ./obj.log
    ps -ef | grep datanode | awk '{print $2}' | xargs kill
    ps -ef | grep dagpool | awk '{print $2}' | xargs kill
    ps -ef | grep objectstore | awk '{print $2}' | xargs kill
}

function result() {
    local command_output
    command_output=$(./dagpool cluster status 2>&1)
    echo "$(green_echo "$command_output")"
}

case "$1" in
    "init")
        init
        ;;
    "all")
        init
        status_and_balance
        expansion
        scaling
        ;;
    "clean")
        clean
        ;;
    "statusAndBalance")
        status_and_balance
        ;;
    "expansion")
        expansion
        ;;
    "scaling")
        scaling
        ;;
    *)
        echo "Invalid argument. Usage: $0 {init|all|clean|statusAndBalance|expansion|scaling}"
        exit 1
        ;;
esac

exit 0
