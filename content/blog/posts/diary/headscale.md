---
title: "headscale 部署使用"
date: 2022-09-21
draft: false
authors: "jobcher"
tags: ["daliy"]
categories: ["日常"]
series: ["日常系列"]
---

# Headscale

`Tailscale` 的控制服务器是不开源的，而且对免费用户有诸多限制，这是人家的摇钱树，可以理解。好在目前有一款开源的实现叫 `Headscale`，这也是唯一的一款，希望能发展壮大。

`Headscale` 由欧洲航天局的 `Juan Font` 使用 Go 语言开发，在 BSD 许可下发布，实现了 `Tailscale` 控制服务器的所有主要功能，可以部署在企业内部，`没有任何设备数量的限制，且所有的网络流量都由自己控制`。

# Headscale 部署

我决定使用`docker-compose`进行部署

## 创建存储

```sh
#!/bin/bash
mkdir -p /opt/headscale
mkdir -p ./config
touch ./config/db.sqlite
curl https://raw.githubusercontent.com/juanfont/headscale/main/config-example.yaml -o ./config/config.yaml

```

## 运行 docker-compose 文件

创建 docker-compose.yaml

```yaml
version: "3"
services:
  headscale:
    image: headscale/headscale:latest
    volumes:
      - ./config:/etc/headscale/
      - ./data:/var/lib/headscale
    ports:
      - 8080:8080
      - 9090:9090
      - 50443:50443
    command: headscale serve
    restart: unless-stopped
```

运行

> docker-compose up -d

# Headscale 使用

## Linux 使用

```sh
wget https://pkgs.tailscale.com/stable/tailscale_1.22.2_amd64.tgz
```

解压

```sh
tar zxvf tailscale_1.22.2_amd64.tgz
```

将二进制文件复制到官方软件包默认的路径下：

```sh
cp tailscale_1.22.2_amd64/tailscaled /usr/sbin/tailscaled
cp tailscale_1.22.2_amd64/tailscale /usr/bin/tailscale
```

将 systemD service 配置文件复制到系统路径下：

```sh
cp tailscale_1.22.2_amd64/systemd/tailscaled.service /lib/systemd/system/tailscaled.service
```

启动 tailscaled.service 并设置开机自启：

```sh
systemctl enable --now tailscaled
```

查看服务状态：

```sh
systemctl status tailscaled
```
