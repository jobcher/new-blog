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
yum -y install ntpdate
#更改时区
timedatectl set-timezone 'Asia/Shanghai'
#查看时间
timedatectl
datetime
```

### 安装 keepalived

```sh
#安装
yum -y install keepalived

```

### 创建 check_apiserver.sh
```sh
# 创建检测脚本
vim /etc/keepalived/check_apiserver.sh 
```
```sh
#!/bin/bash

# VIP 地址（你的 Kubernetes apiserver 将通过这个 VIP 暴露）
APISERVER_VIP="10.10.10.68" # 设定你自己的vip
APISERVER_PORT="6443"

errorExit() {
  echo "*** $*" 1>&2
  exit 1
}

# 检查 apiserver 的 /healthz
curl -sf --max-time 3 https://${APISERVER_VIP}:${APISERVER_PORT}/healthz \
  -k -o /dev/null || errorExit "API Server Unhealthy"
```

### 配置master
```sh
vim /etc/keepalived/keepalived.conf
```
```json
vrrp_script chk_apiserver {
    script "/etc/keepalived/check_apiserver.sh"
    interval 3
    weight -10
    fall 3
    rise 2
}

vrrp_instance VI_1 {
    state MASTER
    interface eth0 # 改为你实际网关
    virtual_router_id 51
    priority 120 #master改为最大值
    advert_int 1

    authentication {
        auth_type PASS
        auth_pass 123456 #改为实际的密码
    }

    virtual_ipaddress {
        10.10.10.68/24  #改为vip地址
    }

    track_script {
        chk_apiserver
    }
}
```
### 启动keepalived
```sh
systemctl enable keepalived --now
systemctl status keepalived
# 检测vip是否正常
ip a | grep eth0
ping 10.10.10.68
curl -k https://10.10.10.68:6443/healthz
```

### 配置 backup
```sh
vim /etc/keepalived/keepalived.conf
```
```sh
vrrp_script chk_apiserver {
    script "/etc/keepalived/check_apiserver.sh"
    interval 3
    weight -10
    fall 3
    rise 2
}

vrrp_instance VI_1 {
    state BACKUP                   # master2/master3 上改为 BACKUP
    interface eth0
    virtual_router_id 51
    priority 100                   # master2: 100, master3: 90
    advert_int 1

    authentication {
        auth_type PASS
        auth_pass 123456
    }

    virtual_ipaddress {
        10.10.10.68/24
    }

    track_script {
        chk_apiserver
    }
}
```
### 启动keepalived
```sh
systemctl enable keepalived --now
systemctl status keepalived
# 检测vip是否正常
ping 10.10.10.68
curl -k https://10.10.10.68:6443/healthz
```