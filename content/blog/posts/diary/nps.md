---
title: "自建服务器内网穿透"
date: 2021-12-08
draft: false
author: "jobcher"
tags: ["内网穿透", "nps"]
categories: ["日常"]
series: ["日常系列"]
---

# 内网穿透

文章中使用的内网穿透前提是必须具有公网 IP 的云服务器，不符合条件的同学可以跳过了。

## nps 内网穿透

nps 是一款轻量级、高性能、功能强大的内网穿透代理服务器。

### 在公网服务器上安装 nps sever 端

```sh
    wget https://github.com/ehang-io/nps/releases/download/v0.26.10/linux_amd64_server.tar.gz
    tar -zxvf linux_amd64_server.tar.gz
    sudo ./nps install
    sudo nps start
```

### 在控制端安装 npc client 端

```sh
    wget https://github.com/ehang-io/nps/releases/download/v0.26.10/linux_amd64_client.tar.gz
    tar -zxvf linux_amd64_client.tar.gz
    sudo ./npc -server=ip:port -vkey=web界面中显示的密钥
    sudo npc start
```

npc 安装完成可以进入 web 页面穿透端口和域名  
http://localhost:8080

## frps 内网穿透

frps 相对于 nps 的劣势是有断流的风险  
frps 相对于 nps 的优势是对于高流量的媒体服务能够提供更可靠的支持

### 安装 frps

```sh
    wget https://code.aliyun.com/MvsCode/frps-onekey/raw/master/install-frps.sh -O ./install-frps.sh
    chmod 700 ./install-frps.sh
    ./install-frps.sh install
```

卸载 frps 服务

```sh
    ./install-frps.sh uninstall
```

更新 frps 服务

```sh
    ./install-frps.sh update
```

Server management（服务管理器）

```sh
    Usage: /etc/init.d/frps {start|stop|restart|status|config|version}
```
