---
title: "Keepalived高可用"
date: 2022-01-05
draft: false
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---

# Keepalived 高可用

配置文件存放位置：/usr/share/doc/keepalived/samples  
VVRP 虚拟路由冗余协议

## 组成

LB 集群：Load Balancing，负载均衡集群，平均分配给多个节点  
HA 集群：High Availability，高可用集群，保证服务可用  
HPC 集群：High Performance Computing，高性能集群

## 配置

keepalived+LVS+nginx

1. 各节点时间必须同步：ntp, chrony
2. 关闭防火墙及 SELinux

### 同步各节点时间

```sh
#安装ntpdate
apt install ntpdate
#更改时区
timedatectl set-timezone 'Asia/Shanghai'
#查看时间
timedatectl
datetime
```

### 安装 keepalived

```sh
#安装
apt install keepalived
#更改模板
cd /usr/share/doc/keepalived/samples

```
