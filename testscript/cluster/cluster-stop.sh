#!/bin/bash

dir=$(dirname $0)

$dir/datanodes-stop.sh

ps -ef | grep dagpool | awk '{print $2}' | xargs kill

ps -ef | grep objectstore | awk '{print $2}' | xargs kill