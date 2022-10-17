#!/bin/bash

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
    nohup ../datanode daemon --listen=127.0.0.1:9011 --datadir=./dn-data1 >dn1log.log 2>&1 &

    nohup ../datanode daemon --listen=127.0.0.1:9012 --datadir=./dn-data2 >dn2log.log 2>&1 &

    nohup ../datanode daemon --listen=127.0.0.1:9013 --datadir=./dn-data3 >dn3log.log 2>&1 &

    nohup ../dagpool daemon --datadir=./dp-data --config=../conf/node_config.json  >dplog.log 2>&1 &

    nohup ../objectstore daemon --pool-addr=127.0.0.1:50001 --pool-user=dagpool --pool-password=dagpool --datadir=./store-data >objlog.log 2>&1 &

    sleep 3

    ../iam-tools add-user --admin-access-key=filedagadmin --admin-secret-key=filedagadmin --username=testA --password=testAtestA >/dev/null 2>&1
    ../iam-tools add-user --admin-access-key=filedagadmin --admin-secret-key=filedagadmin --username=testB --password=testBtestB >/dev/null 2>&1

    mc --version
    if [ $? -ne 0 ]; then
      # install mc
      uNames=`uname -s`
      osName=${uNames: 0: 4}
      if [ "$osName" = "Darw" ] # Darwin
      then
        echo "Mac OS X"
        wget https://dl.min.io/client/mc/release/darwin-amd64/mc
        chmod +x mc
        mv mc /usr/local/bin
      elif [ "$osName" = "Linu" ] # Linux
      then
        echo "GNU/Linux"
        wget https://dl.min.io/client/mc/release/linux-amd64/mc
        chmod +x mc
        mv mc /usr/local/bin
      else
        echo "unknown os"
      fi

      mc --version
      if [ $? -ne 0 ]; then
        echo "install mc failed"
        echo "manually download and install the mc, https://dl.min.io/client/mc/release/"
      fi
    fi

    mc alias set loc http://127.0.0.1:9985 filedagadmin filedagadmin
    mc alias set loc1 http://127.0.0.1:9985 testA testAtestA
    mc alias set loc2 http://127.0.0.1:9985 testB testBtestB
}
function close() {
    ps aux | grep objectstore | grep -v grep | awk '{print $2}'| xargs kill -9
    ps aux | grep dagpool | grep -v grep | awk '{print $2}'| xargs kill -9
    ps aux | grep datanode | grep -v grep | awk '{print $2}'| xargs kill -9

    rm dn1log.log dn2log.log dn3log.log dplog.log objlog.log
    rm -rf ./store-data
    rm -rf ./dp-data
    rm -rf ./dn-data*

    mc alias remove loc
    mc alias remove loc1
    mc alias remove loc2
}

function test_except() {
    if [ "$2" -eq "$3" ]; then
      echo -e "\033[1m\033[;32m[OK] $1, expected the result to be $2, and found $3 \033[0m"
    else
      echo -e "\033[1m\033[;31m[FAILED] $1, expected the result to be $2, but instead found $3 \033[0m"
    fi
}
