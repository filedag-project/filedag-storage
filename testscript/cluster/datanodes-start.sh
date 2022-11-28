#!/bin/bash

dir=$(dirname $0)

nohup $dir/../../datanode daemon --listen=127.0.0.1:9011 --datadir=/tmp/dn-data1 > /tmp/dn1.log 2>&1 &
nohup $dir/../../datanode daemon --listen=127.0.0.1:9012 --datadir=/tmp/dn-data2 > /tmp/dn2.log 2>&1 &
nohup $dir/../../datanode daemon --listen=127.0.0.1:9013 --datadir=/tmp/dn-data3 > /tmp/dn3.log 2>&1 &

nohup $dir/../../datanode daemon --listen=127.0.0.1:9014 --datadir=/tmp/dn-data4 > /tmp/dn4.log 2>&1 &
nohup $dir/../../datanode daemon --listen=127.0.0.1:9015 --datadir=/tmp/dn-data5 > /tmp/dn5.log 2>&1 &
nohup $dir/../../datanode daemon --listen=127.0.0.1:9016 --datadir=/tmp/dn-data6 > /tmp/dn6.log 2>&1 &

nohup $dir/../../datanode daemon --listen=127.0.0.1:9017 --datadir=/tmp/dn-data7 > /tmp/dn7.log 2>&1 &
nohup $dir/../../datanode daemon --listen=127.0.0.1:9018 --datadir=/tmp/dn-data8 > /tmp/dn8.log 2>&1 &
nohup $dir/../../datanode daemon --listen=127.0.0.1:9019 --datadir=/tmp/dn-data9 > /tmp/dn9.log 2>&1 &