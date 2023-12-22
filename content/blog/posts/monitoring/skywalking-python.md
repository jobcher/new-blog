---
title: "skywalking python agent 安装和配置"
date: 2023-05-31
draft: false
author: "jobcher"
featuredImage: "/images/skwalking-python.png"
featuredImagePreview: "/images/skwalking-python.png"
tags: ["skywalking"]
categories: ["监控"]
series: ["运维监控系列"]
---
## 背景
skywalking 是一个APM监控，在java和微服务领域非常流行。但是在python领域，skywalking的使用率并不高。本文将介绍如何安装和配置skywalking python agent。  
`SkyWalking Python` 代理实现了一个命令行界面，可用于在部署期间将代理附加到出色的应用程序，而无需更改任何应用程序代码，就像 SkyWalking Java 代理一样。
### 安装skywalking python agent
```shell
pip install apache-skywalking
```
>运行 `sw-python` 以查看它是否可用，您需要按环境变量传递配置。
### 配置skywalking python agent
```shell
export SW_AGENT_NAME=服务名称
# 服务器地址
export SW_AGENT_COLLECTOR_BACKEND_SERVICES=127.0.0.1:11800
```

### 执行python程序
```shell
sw-python run -p python app.py
```
#### 启用 CLI 调试模式，以便在启动应用程序时查看代理日志：
```shell
sw-python -d run python app.py
```
