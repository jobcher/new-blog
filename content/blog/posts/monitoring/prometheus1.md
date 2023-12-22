---
title: "prometheus grafana alertmanager 安装配置"
date: 2022-01-13
draft: false
authors: "jobcher"
tags: ["prometheus"]
categories: ["监控"]
series: ["运维监控系列"]
---

# prometheus+grafana+alertmanager 安装配置

服务器监控告警系统搭建，通过 exporter 获取节点信息到 prometheus。prometheus 配置规则，使 garfana 和 alertmanager 能够接受到数据，分别展示数据和发送告警

## 参数

VM :192.168.99.78

| 端口 | 服务              |
| :--- | :---------------- |
| 9100 | node_exporter     |
| 3000 | grafana           |
| 9090 | prometheus        |
| 9115 | blackbox_exporter |

## 安装

### grafa 安装

1. docker 安装

```sh
docker run -d -p 3000:3000 \
--name=grafana \
-v grafana-storage:/var/lib/grafana \
grafana/grafana:8.3.3
```

### prometheus 安装

1. 下载

```sh
wget https://github.com/prometheus/prometheus/releases/download/v2.32.1/prometheus-2.32.1.linux-amd64.tar.gz
tar -zxvf prometheus-2.32.1.linux-amd64.tar.gz
cd prometheus-2.32.1.linux-amd64
mkdir -p file_sd
mkdir -p rules
```

2. 运行 prometheus

```sh
killall prometheus
nohup ./prometheus --config.file=prometheus.yml &
# 查看运行状况
tail -f nohup.out
```

### node_exporter 安装

1. docker-compose 安装

```yml
version: "3"
services:
  node-exporter:
    image: prom/node-exporter:v1.3.1
    container_name: node-exporter
    restart: always
    ports:
      - "9100:9100"
```

```sh
docker-compose up -d
```

2. 二进制安装

```sh
wget https://github.com/prometheus/node_exporter/releases/download/v1.3.1/node_exporter-1.3.1.linux-amd64.tar.gz
tar -zxvf node_exporter-1.3.1.linux-amd64.tar.gz
cd node_exporter-1.3.1.linux-amd64
nohup ./node_exporter &
```

### blackbox_exporter

1. 二进制安装

```sh
wget https://github.com/prometheus/blackbox_exporter/releases/download/v0.19.0/blackbox_exporter-0.19.0.linux-amd64.tar.gz
tar -zxvf blackbox_exporter-0.19.0.linux-amd64.tar.gz
cd blackbox_exporter-0.19.0.linux-amd64
nohup ./blackbox_exporter &
```

### Alertmanager

1. 二进制安装

```sh
wget https://github.com/prometheus/alertmanager/releases/download/v0.23.0/alertmanager-0.23.0.linux-amd64.tar.gz
tar -zxvf alertmanager-0.23.0.linux-amd64.tar.gz
cd alertmanager-0.23.0.linux-amd64
nohup ./alertmanager &
```

2. docker-compose 安装

```yaml
version: "3"
services:
  alertmanager:
    image: "prom/alertmanager:v0.22.2"
    volumes:
      - "/etc/localtime:/etc/localtime"
      - "./alertmanager.yml:/etc/alertmanager/alertmanager.yml"
    ports:
      - "9093:9093"
    restart: "always"
    container_name: "alertmanager"
```

## prometheus.yml 配置

```yaml
global:
  scrape_interval: 15s # By default, scrape targets every 15 seconds.

  # Attach these labels to any time series or alerts when communicating with
  # external systems (federation, remote storage, Alertmanager).
  external_labels:
    monitor: "codelab-monitor"

# Alertmanager configuration
# alerting:
#   alertmanagers:
#   - static_configs:
#     - targets:
#       - 192.168.99.78:9093

# Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
# rule_files:
#   - "./rules/blackbox.yaml"
#   - "./rules/node-exporter.yaml"

# A scrape configuration containing exactly one endpoint to scrape:
# Here it's Prometheus itself.
scrape_configs:
  # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
  - job_name: "prometheus"
    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s
    static_configs:
      - targets: ["localhost:9090"]

  - job_name: "node-exporter"
    file_sd_configs:
      - files:
          - "./file_sd/node-exporter.yaml"
        refresh_interval: 5s

  - job_name: "blackbox"
    metrics_path: /probe
    scrape_interval: 30s
    scrape_timeout: 30s
    params:
      module: [http_2xx] # Look for a HTTP 200 response.
    file_sd_configs:
      - files:
          - "./file_sd/blackbox.yaml"
        refresh_interval: 5s
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: 192.168.99.78:9115
```

### node-exporter.yaml

```yaml
- targets: ["192.168.99.78:9100"]
  labels:
    instance: <实例名称>
- targets: ["<IP>:9100"]
  labels:
    instance: 实例名称
```

### blackbox.yaml

```yaml
- targets:
    - https://www.jobcher.com
    - https://<域名>
```
