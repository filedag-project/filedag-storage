#!/bin/bash

ps -ef | grep datanode | awk '{print $2}' | xargs kill
