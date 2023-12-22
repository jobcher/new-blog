---
title: "Cloudflare Zero Trust 内网穿透"
date: 2022-09-30
draft: false
author: "jobcher"
featuredImage: "/images/cloudflare.png"
featuredImagePreview: "/images/cloudflare.png"
tags: ["内网穿透", "nps"]
categories: ["日常"]
series: ["日常系列"]
---

# Cloudflare Zero Trust 内网穿透

![cloud](/images/cloudflare.png)  
最快的 `Zero Trust` 应用访问和互联网浏览平台  
增加可见性，消除复杂性，降低远程和办公室用户的风险。杜绝数据丢失、恶意软件和网络钓鱼，保护用户、应用程序和设备安全。  
使用 Tunnel 隧道来实现内网传统，实现内网访问各类应用

## 安装部署

[https://dash.teams.cloudflare.com/](https://dash.teams.cloudflare.com/)

### Docker 部署

在 docker 环境运行 `<token>` 是你个人令牌

```sh
docker run -d --name cloudflared cloudflare/cloudflared:latest tunnel --no-autoupdate run --token <token>
```

### Linux 部署

1. X86-64 位

```sh
curl -L --output cloudflared.rpm https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-x86_64.rpm && \
sudo yum localinstall -y cloudflared.rpm && \
sudo cloudflared service install <token>

```

2. X86-32 位

```sh
curl -L --output cloudflared.rpm https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-386.rpm &&

sudo yum localinstall -y cloudflared.rpm &&

sudo cloudflared service install <token>
```

3. arm64

```sh
curl -L --output cloudflared.rpm https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-aarch64.rpm &&

sudo yum localinstall -y cloudflared.rpm &&

sudo cloudflared service install <token>

```

4. arm32

```sh
curl -L --output cloudflared.rpm https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-arm.rpm &&

sudo yum localinstall -y cloudflared.rpm &&

sudo cloudflared service install <token>
```

### windows 部署

1. 下载 [https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.msi](https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.msi)
2. 运行安装程序
3. 以管理员身份打开命令提示符
4. 运行以下命令

```sh
cloudflared.exe service install <token>
```

## 在 cloudflare 里配置内网穿透

![image](/images/cloudflare2.png)

欢迎关注我的博客[www.jobcher.com](https://www.jobcher.com/)
