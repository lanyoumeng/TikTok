#!/bin/bash

# 启动 Flink 集群
bin/flink start-cluster

# 等待集群启动
sleep 15

# 提交 Flink CDC 作业
bin/flink-cdc.sh job/user-mysql-to-kafka.yaml
bin/flink-cdc.sh job/video-mysql-to-kafka.yaml
bin/flink-cdc.sh job/favorite-mysql-to-kafka.yaml
bin/flink-cdc.sh job/comment-mysql-to-kafka.yaml
bin/flink-cdc.sh job/relation-mysql-to-kafka.yaml
bin/flink-cdc.sh job/message-mysql-to-kafka.yaml

# 保持容器运行
tail -f $FLINK_HOME/log/flink-*.log
