# 使用一个合适的基础镜像，通常使用 OpenJDK 作为基础镜像
FROM  openjdk:8-jre-slim

# 设置环境变量，您可以根据需要调整这些值
ENV FLINK_HOME /opt/flink
ENV PATH $FLINK_HOME/bin:$PATH

# 创建 Flink 的安装目录
RUN mkdir -p $FLINK_HOME

# 将本地的 Flink 文件夹复制到 Docker 镜像中的指定目录
COPY flink-1.18.0/* $FLINK_HOME/

# 设置工作目录
WORKDIR $FLINK_HOME

# 暴露 Flink 的默认端口
EXPOSE 8081 6123 6124

# 启动 Flink 集群并提交作业的脚本
COPY start-flink.sh /opt/
RUN chmod +x /opt/start-flink.sh

# 启动 Flink 集群并提交作业
ENTRYPOINT ["/opt/start-flink.sh"]
