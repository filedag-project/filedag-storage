#!/bin/bash

dir=$(dirname $0)

$dir/datanodes-start.sh

nohup $dir/../../dagpool daemon --datadir=/tmp/dp-db > /tmp/dp.log 2>&1 &

nohup $dir/../../objectstore daemon --datadir=/tmp/store-data --pool-addr=127.0.0.1:50001 --pool-user=dagpool --pool-password=dagpool > /tmp/objstore.log 2>&1 &