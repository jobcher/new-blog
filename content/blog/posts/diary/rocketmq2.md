---
title: "RocketMQ docker-compose部署 4主4从集群"
date: 2022-08-15
draft: true
author: "jobcher"
tags: ["RocketMQ"]
categories: ["日常"]
series: ["日常系列"]
---

# RocketMQ docker-compose 部署 4 主 4 从集群

V 4.8.0  
采用`4主4从`，`同步模式`。HA 实现上采用`Master/Slave+Failover`组件方式
每台主机运行`三个容器`，分别为`NameServer`、`BrokerMaster`、`SlaveMaster`，每个 Master 和 Slave 分别存放在不同的机器上

## 架构


| IP           | 角色         | 服务       |
| :----------- | :----------- | :--------- |
| 193.0.40.172 | NameServer   | -          |
| 193.0.40.172 | BrokerMaster | broker-a   |
| 193.0.40.172 | SlaveMaster  | broker-d-s |
| 193.0.40.172 | BrokerMaster | broker-b   |
| 193.0.40.172 | SlaveMaster  | broker-a-s |
| 193.0.40.172 | BrokerMaster | broker-c   |
| 193.0.40.172 | SlaveMaster  | broker-b-s |
| 193.0.40.172 | BrokerMaster | broker-d   |
| 193.0.40.172 | SlaveMaster  | broker-c-s |

## 部署

### 安装 docker-compose

```sh
#!/bin/bash
# 下载安装 v2.4.1 docker-compose
curl -L https://get.daocloud.io/docker/compose/releases/download/v2.4.1/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
docker-compose --version
```

执行 `docker-compose --version` 查看是否安装成功

### 生成配置文件

```sh
#!/bin/bash
#docker-compose 生成配置文件

mkdir -p /rocketmq/data/namesv1
mkdir -p /rocketmq/logs/namesv1
mkdir -p /rocketmq/data/namesv2
mkdir -p /rocketmq/logs/namesv2
mkdir -p /rocketmq/config/broker-a
mkdir -p /rocketmq/data/broker-a
mkdir -p /rocketmq/logs/broker-a
mkdir -p /rocketmq/config/broker-a-s
mkdir -p /rocketmq/data/broker-a-s
mkdir -p /rocketmq/logs/broker-a-s
mkdir -p /rocketmq/config/broker-b
mkdir -p /rocketmq/data/broker-b
mkdir -p /rocketmq/logs/broker-b
mkdir -p /rocketmq/config/broker-b-s
mkdir -p /rocketmq/data/broker-b-s
mkdir -p /rocketmq/logs/broker-b-s
mkdir -p /rocketmq/config/broker-c
mkdir -p /rocketmq/data/broker-c
mkdir -p /rocketmq/logs/broker-c
mkdir -p /rocketmq/config/broker-c-s
mkdir -p /rocketmq/data/broker-c-s
mkdir -p /rocketmq/logs/broker-c-s
mkdir -p /rocketmq/config/broker-d
mkdir -p /rocketmq/data/broker-d
mkdir -p /rocketmq/logs/broker-d
mkdir -p /rocketmq/config/broker-d-s
mkdir -p /rocketmq/data/broker-d-s
mkdir -p /rocketmq/logs/broker-d-s
cd /rocketmq/config/broker-a
cat > broker-a.conf <<EOF
#集群名称
brokerClusterName=DefaultCluster
#broker名称
brokerName=broker-a
#brokerId master用0 slave用其他
brokerId=0
#清理时机
deleteWhen=4
#文件保留时长 48小时
fileReservedTime=48
#broker角色 -ASYNC_MASTER异步复制 -SYNC_MASTER同步双写 -SLAVE
brokerRole=SYNC_MASTER
#刷盘策略 - ASYNC_FLUSH 异步刷盘 - SYNC_FLUSH 同步刷盘
flushDiskType=SYNC_FLUSH
#主机ip
brokerIP1=193.0.40.172
#对外服务的监听接口，同一台机器上部署多个broker,端口号要不相同
listenPort=10911
#namesvr
namesrvAddr=193.0.40.172:9876;193.0.40.172:9877
#是否能够自动创建topic
autoCreateTopicEnable=true
EOF
# 生成配置文件
cd /rocketmq/config/broker-a-s
cat > broker-a-s.conf <<EOF
#集群名称
brokerClusterName=DefaultCluster
#broker名称
brokerName=broker-a
#brokerId master用0 slave用其他
brokerId=1
#清理时机
deleteWhen=4
#文件保留时长 48小时
fileReservedTime=48
#broker角色 -ASYNC_MASTER异步复制 -SYNC_MASTER同步双写 -SLAVE
brokerRole=SLAVE
#刷盘策略 - ASYNC_FLUSH 异步刷盘 - SYNC_FLUSH 同步刷盘
flushDiskType=SYNC_FLUSH
#主机ip
brokerIP1=193.0.40.172
#对外服务的监听接口，同一台机器上部署多个broker,端口号要不相同
listenPort=11911
#namesrv
namesrvAddr=193.0.40.172:9876;193.0.40.172:9877
#是否能够自动创建topic
autoCreateTopicEnable=true
EOF
cd /rocketmq/config/broker-b
cat > broker-b.conf <<EOF
#集群名称
brokerClusterName=DefaultCluster
#broker名称
brokerName=broker-b
#brokerId master用0 slave用其他
brokerId=0
#清理时机
deleteWhen=4
#文件保留时长 48小时
fileReservedTime=48
#broker角色 -ASYNC_MASTER异步复制 -SYNC_MASTER同步双写 -SLAVE
brokerRole=SYNC_MASTER
#刷盘策略 - ASYNC_FLUSH 异步刷盘 - SYNC_FLUSH 同步刷盘
flushDiskType=SYNC_FLUSH
#主机ip
brokerIP1=193.0.40.172
#对外服务的监听接口，同一台机器上部署多个broker,端口号要不相同
listenPort=12911
#namesrv
namesrvAddr=193.0.40.172:9876;193.0.40.172:9877
#是否能够自动创建topic
autoCreateTopicEnable=true
EOF
cd /rocketmq/config/broker-b-s
cat > broker-b-s.conf <<EOF
#集群名称
brokerClusterName=DefaultCluster
#broker名称
brokerName=broker-b
#brokerId master用0 slave用其他
brokerId=1
#清理时机
deleteWhen=4
#文件保留时长 48小时
fileReservedTime=48
#broker角色 -ASYNC_MASTER异步复制 -SYNC_MASTER同步双写 -SLAVE
brokerRole=SLAVE
#刷盘策略 - ASYNC_FLUSH 异步刷盘 - SYNC_FLUSH 同步刷盘
flushDiskType=SYNC_FLUSH
#主机ip
brokerIP1=193.0.40.172
#对外服务的监听接口，同一台机器上部署多个broker,端口号要不相同
listenPort=13911
#namesrv
namesrvAddr=193.0.40.172:9876;193.0.40.172:9877
#是否能够自动创建topic
autoCreateTopicEnable=true
EOF
cd /rocketmq/config/broker-c
cat > broker-c.conf <<EOF
#集群名称
brokerClusterName=DefaultCluster
#broker名称
brokerName=broker-c
#brokerId master用0 slave用其他
brokerId=0
#清理时机
deleteWhen=4
#文件保留时长 48小时
fileReservedTime=48
#broker角色 -ASYNC_MASTER异步复制 -SYNC_MASTER同步双写 -SLAVE
brokerRole=SYNC_MASTER
#刷盘策略 - ASYNC_FLUSH 异步刷盘 - SYNC_FLUSH 同步刷盘
flushDiskType=SYNC_FLUSH
#主机ip
brokerIP1=193.0.40.172
#对外服务的监听接口，同一台机器上部署多个broker,端口号要不相同
listenPort=14911
#namesrv
namesrvAddr=193.0.40.172:9876;193.0.40.172:9877
#是否能够自动创建topic
autoCreateTopicEnable=true
EOF
cd /rocketmq/config/broker-c-s
cat > broker-c-s.conf <<EOF
#集群名称
brokerClusterName=DefaultCluster
#broker名称
brokerName=broker-c
#brokerId master用0 slave用其他
brokerId=1
#清理时机
deleteWhen=4
#文件保留时长 48小时
fileReservedTime=48
#broker角色 -ASYNC_MASTER异步复制 -SYNC_MASTER同步双写 -SLAVE
brokerRole=SLAVE
#刷盘策略 - ASYNC_FLUSH 异步刷盘 - SYNC_FLUSH 同步刷盘
flushDiskType=SYNC_FLUSH
#主机ip
brokerIP1=193.0.40.172
#对外服务的监听接口，同一台机器上部署多个broker,端口号要不相同
listenPort=15911
#namesrv
namesrvAddr=193.0.40.172:9876;193.0.40.172:9877
#是否能够自动创建topic
autoCreateTopicEnable=true
EOF
cd /rocketmq/config/broker-d
cat > broker-d.conf <<EOF
#集群名称
brokerClusterName=DefaultCluster
#broker名称
brokerName=broker-d
#brokerId master用0 slave用其他
brokerId=0
#清理时机
deleteWhen=4
#文件保留时长 48小时
fileReservedTime=48
#broker角色 -ASYNC_MASTER异步复制 -SYNC_MASTER同步双写 -SLAVE
brokerRole=SYNC_MASTER
#刷盘策略 - ASYNC_FLUSH 异步刷盘 - SYNC_FLUSH 同步刷盘
flushDiskType=SYNC_FLUSH
#主机ip
brokerIP1=193.0.40.172
#对外服务的监听接口，同一台机器上部署多个broker,端口号要不相同
listenPort=16911
#namesrv
namesrvAddr=193.0.40.172:9876;193.0.40.172:9877
#是否能够自动创建topic
autoCreateTopicEnable=true
EOF
cd /rocketmq/config/broker-d-s
cat > broker-d-s.conf <<EOF
#集群名称
brokerClusterName=DefaultCluster
#broker名称
brokerName=broker-d
#brokerId master用0 slave用其他
brokerId=1
#清理时机
deleteWhen=4
#文件保留时长 48小时
fileReservedTime=48
#broker角色 -ASYNC_MASTER异步复制 -SYNC_MASTER同步双写 -SLAVE
brokerRole=SLAVE
#刷盘策略 - ASYNC_FLUSH 异步刷盘 - SYNC_FLUSH 同步刷盘
flushDiskType=SYNC_FLUSH
#主机ip
brokerIP1=193.0.40.172
#对外服务的监听接口，同一台机器上部署多个broker,端口号要不相同
listenPort=17911
#namesrv
namesrvAddr=193.0.40.172:9876;193.0.40.172:9877
#是否能够自动创建topic
autoCreateTopicEnable=true
EOF
```

执行`rocker-config.sh`,确认配置文件是否生成

> tree /rocketmq

### 执行 docker-compose.yaml 文件

```yaml
version: "3"
services:
  rocketmq-namesv1:
    image: apache/rocketmq:4.8.0
    container_name: rocketmq-namesv1
    restart: always
    ports:
      - 9876:9876
    volumes:
      - /rocketmq/logs/namesv1:/home/rocketmq/logs
    environment:
      JAVA_OPT_EXT: -server -Xms256M -Xmx256M -Xmn128m
    command: sh mqnamesrv
    networks:
      rocketmq:
        aliases:
          - rocketmq-namesv1
  rocketmq-namesv2:
    image: apache/rocketmq:4.8.0
    container_name: rocketmq-namesv2
    restart: always
    ports:
      - 9877:9876
    volumes:
      - /rocketmq/logs/namesv2:/home/rocketmq/logs
    environment:
      JAVA_OPT_EXT: -server -Xms256M -Xmx256M -Xmn128m
    command: sh mqnamesrv
    networks:
      rocketmq:
        aliases:
          - rocketmq-namesv2
  broker-a:
    image: apache/rocketmq:4.8.0
    container_name: broker-a
    links:
      - rocketmq-namesv1:rocketmq-namesv1
      - rocketmq-namesv2:rocketmq-namesv2
    ports:
      - 10909:10909
      - 10911:10911
      - 10912:10912
    environment:
      TZ: Asia/Shanghai
      NAMESRV_ADDR: "rocketmq-namesv1:9876"
      JAVA_OPT_EXT: "-server -Xms256M -Xmx256M -Xmn128m"
    volumes:
      - /rocketmq/logs/broker-a:/home/rocketmq/logs
      - /rocketmq/config/broker-a/broker-a.conf:/home/rocketmq/rocketmq-4.8.0/conf/broker.conf
    command: sh mqbroker  -c /home/rocketmq/rocketmq-4.8.0/conf/broker.conf autoCreateTopicEnable=true &
    networks:
      rocketmq:
        aliases:
          - broker-a
  broker-a-s:
    image: apache/rocketmq:4.8.0
    container_name: broker-a-s
    links:
      - rocketmq-namesv1:rocketmq-namesv1
      - rocketmq-namesv2:rocketmq-namesv2
    ports:
      - 11909:10909
      - 11911:11911
      - 11912:10912
    environment:
      TZ: Asia/Shanghai
      NAMESRV_ADDR: "rocketmq-namesv1:9876"
      JAVA_OPT_EXT: "-server -Xms256M -Xmx256M -Xmn128m"
    volumes:
      - /rocketmq/logs/broker-a-s:/home/rocketmq/logs
      - /rocketmq/config/broker-a-s/broker-a-s.conf:/home/rocketmq/rocketmq-4.8.0/conf/broker.conf
    command: sh mqbroker  -c /home/rocketmq/rocketmq-4.8.0/conf/broker.conf autoCreateTopicEnable=true &
    networks:
      rocketmq:
        aliases:
          - broker-a-s
  broker-b:
    image: apache/rocketmq:4.8.0
    container_name: broker-b
    links:
      - rocketmq-namesv1:rocketmq-namesv1
      - rocketmq-namesv2:rocketmq-namesv2
    ports:
      - 12909:10909
      - 12911:12911
      - 12912:10912
    environment:
      TZ: Asia/Shanghai
      NAMESRV_ADDR: "rocketmq-namesv1:9876"
      JAVA_OPT_EXT: "-server -Xms256M -Xmx256M -Xmn128m"
    volumes:
      - /rocketmq/logs/broker-b:/home/rocketmq/logs
      - /rocketmq/config/broker-b/broker-b.conf:/home/rocketmq/rocketmq-4.8.0/conf/broker.conf
    command: sh mqbroker  -c /home/rocketmq/rocketmq-4.8.0/conf/broker.conf autoCreateTopicEnable=true &
    networks:
      rocketmq:
        aliases:
          - broker-b
  broker-b-s:
    image: apache/rocketmq:4.8.0
    container_name: broker-b-s
    links:
      - rocketmq-namesv1:rocketmq-namesv1
      - rocketmq-namesv2:rocketmq-namesv2
    ports:
      - 13909:10909
      - 13911:13911
      - 13912:10912
    environment:
      TZ: Asia/Shanghai
      NAMESRV_ADDR: "rocketmq-namesv1:9876"
      JAVA_OPT_EXT: "-server -Xms256M -Xmx256M -Xmn128m"
    volumes:
      - /rocketmq/logs/broker-b-s:/home/rocketmq/logs
      - /rocketmq/config/broker-b-s/broker-b-s.conf:/home/rocketmq/rocketmq-4.8.0/conf/broker.conf
    command: sh mqbroker  -c /home/rocketmq/rocketmq-4.8.0/conf/broker.conf autoCreateTopicEnable=true &
    networks:
      rocketmq:
        aliases:
          - broker-b-s
  broker-c:
    image: apache/rocketmq:4.8.0
    container_name: broker-c
    links:
      - rocketmq-namesv1:rocketmq-namesv1
      - rocketmq-namesv2:rocketmq-namesv2
    ports:
      - 14909:10909
      - 14911:14911
      - 14912:10912
    environment:
      TZ: Asia/Shanghai
      NAMESRV_ADDR: "rocketmq-namesv1:9876"
      JAVA_OPT_EXT: "-server -Xms256M -Xmx256M -Xmn128m"
    volumes:
      - /rocketmq/logs/broker-c:/home/rocketmq/logs
      - /rocketmq/config/broker-c/broker-c.conf:/home/rocketmq/rocketmq-4.8.0/conf/broker.conf
    command: sh mqbroker  -c /home/rocketmq/rocketmq-4.8.0/conf/broker.conf autoCreateTopicEnable=true &
    networks:
      rocketmq:
        aliases:
          - broker-c
  broker-c-s:
    image: apache/rocketmq:4.8.0
    container_name: broker-c-s
    links:
      - rocketmq-namesv1:rocketmq-namesv1
      - rocketmq-namesv2:rocketmq-namesv2
    ports:
      - 15909:10909
      - 15911:15911
      - 15912:10912
    environment:
      TZ: Asia/Shanghai
      NAMESRV_ADDR: "rocketmq-namesv1:9876"
      JAVA_OPT_EXT: "-server -Xms256M -Xmx256M -Xmn128m"
    volumes:
      - /rocketmq/logs/broker-c-s:/home/rocketmq/logs
      - /rocketmq/config/broker-c-s/broker-c-s.conf:/home/rocketmq/rocketmq-4.8.0/conf/broker.conf
    command: sh mqbroker  -c /home/rocketmq/rocketmq-4.8.0/conf/broker.conf autoCreateTopicEnable=true &
    networks:
      rocketmq:
        aliases:
          - broker-c-s
  broker-d:
    image: apache/rocketmq:4.8.0
    container_name: broker-d
    links:
      - rocketmq-namesv1:rocketmq-namesv1
      - rocketmq-namesv2:rocketmq-namesv2
    ports:
      - 16909:10909
      - 16911:16911
      - 16912:10912
    environment:
      TZ: Asia/Shanghai
      NAMESRV_ADDR: "rocketmq-namesv1:9876"
      JAVA_OPT_EXT: "-server -Xms256M -Xmx256M -Xmn128m"
    volumes:
      - /rocketmq/logs/broker-d:/home/rocketmq/logs
      - /rocketmq/config/broker-d/broker-d.conf:/home/rocketmq/rocketmq-4.8.0/conf/broker.conf
    command: sh mqbroker  -c /home/rocketmq/rocketmq-4.8.0/conf/broker.conf autoCreateTopicEnable=true &
    networks:
      rocketmq:
        aliases:
          - broker-d
  broker-d-s:
    image: apache/rocketmq:4.8.0
    container_name: broker-d-s
    links:
      - rocketmq-namesv1:rocketmq-namesv1
      - rocketmq-namesv2:rocketmq-namesv2
    ports:
      - 17909:10909
      - 17911:17911
      - 17912:10912
    environment:
      TZ: Asia/Shanghai
      NAMESRV_ADDR: "rocketmq-namesv1:9876"
      JAVA_OPT_EXT: "-server -Xms256M -Xmx256M -Xmn128m"
    volumes:
      - /rocketmq/logs/broker-d-s:/home/rocketmq/logs
      - /rocketmq/config/broker-d-s/broker-d-s.conf:/home/rocketmq/rocketmq-4.8.0/conf/broker.conf
    command: sh mqbroker  -c /home/rocketmq/rocketmq-4.8.0/conf/broker.conf autoCreateTopicEnable=true &
    networks:
      rocketmq:
        aliases:
          - broker-d-s
  rocketmq-console:
    image: styletang/rocketmq-console-ng
    container_name: rocketmq-console
    ports:
      - 8090:8080
    environment:
      JAVA_OPTS: -Drocketmq.namesrv.addr=rocketmq-namesv1:9876;rocketmq-namesv2:9877 -Dcom.rocketmq.sendMessageWithVIPChannel=false
    networks:
      rocketmq:
        aliases:
          - rocketmq-console
networks:
  rocketmq:
    name: rocketmq
    driver: bridge
```

执行 docker-compose 部署

> docker-compose up -d
