#!/bin/bash

# 启动 Flink 集群
./bin/start-cluster.sh

export FLINK_HOME=/home/lanmengyou/code/go_code/TikTok/flink/flink-1.18.0

# 等待集群启动
sleep 5

# 提交 Flink CDC 作业
bin/flink-cdc.sh job/user-mysql-to-kafka.yaml
bin/flink-cdc.sh job/video-mysql-to-kafka.yaml
bin/flink-cdc.sh job/favorite-mysql-to-kafka.yaml
bin/flink-cdc.sh job/comment-mysql-to-kafka.yaml
bin/flink-cdc.sh job/relation-mysql-to-kafka.yaml
bin/flink-cdc.sh job/message-mysql-to-kafka.yaml

# 保持容器运行
tail -f $FLINK_HOME/log/flink-*.log
