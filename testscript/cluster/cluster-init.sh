#!/bin/bash

dir=$(dirname $0)

$dir/../../dagpool cluster add $dir/../../conf/node_config.json $dir/../../conf/node_config2.json $dir/../../conf/node_config3.json
$dir/../../dagpool cluster balance