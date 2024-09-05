FROM golang:1.21 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    netbase

RUN rm -rf /var/lib/apt/lists/ && \
    apt-get autoremove -y && \
    apt-get autoclean -y

RUN apt-get update && \
    apt-get install -y ffmpeg


COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8882
EXPOSE 8002
VOLUME /data/conf

#CMD ["./server", "-conf", "/data/conf"]
#./* 会报错
CMD ["./video"]
