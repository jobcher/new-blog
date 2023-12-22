---
title: "prometheus 配置"
date: 2022-01-13
draft: false
author: "jobcher"
tags: ["prometheus"]
categories: ["监控"]
series: ["运维监控系列"]
---

# prometheus 配置

![prometheus](/images/prometheus.svg)  
Prometheus 是由 SoundCloud 开源监控告警解决方案

## 组件

1. Prometheus Server， 主要用于抓取数据和存储时序数据，另外还提供查询和 Alert Rule 配置管理。
2. client libraries，用于对接 Prometheus Server, 可以查询和上报数据。
3. push gateway ，用于批量，短期的监控数据的汇总节点，主要用于业务数据汇报等。
   各种汇报数据的 exporters ，例如汇报机器数据的 node_exporter, 汇报 MongoDB 信息的 MongoDB exporter 等等。
4. 用于告警通知管理的 alertmanager 。

## 运行逻辑

1. Prometheus server 定期从静态配置的 targets 或者服务发现的 targets 拉取数据。
2. 当新拉取的数据大于配置内存缓存区的时候，Prometheus 会将数据持久化到磁盘（如果使用 remote storage 将持久化到云端）。
3. Prometheus 可以配置 rules，然后定时查询数据，当条件触发的时候，会将 alert 推送到配置的 Alertmanager。
4. Alertmanager 收到警告的时候，可以根据配置，聚合，去重，降噪，最后发送警告。
   可以使用 API， Prometheus Console 或者 Grafana 查询和聚合数据。

## 安装 prometheus

1. 使用预编译的二进制文件安装

```sh
wget https://github.com/prometheus/prometheus/releases/download/v2.32.1/prometheus-2.32.1.linux-amd64.tar.gz
tar -zxvf prometheus-2.32.1.linux-amd64.tar.gz
cd prometheus-2.32.1.linux-amd64
```

2. 使用 docker 安装

```sh
mkdir -p opt/prometheus
vim prometheus.yml
docker run \
    -p 9090:9090 \
    -v /path/to/prometheus.yml:/opt/prometheus/prometheus.yml \
    prom/prometheus
```
