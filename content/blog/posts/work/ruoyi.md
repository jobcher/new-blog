---
title: "ruoyi-cloud docker部署"
date: 2022-03-03
draft: false
author: "jobcher"
tags: ["work"]
categories: ["work"]
series: ["问题库系列"]
---

## 基础环境安装

```sh
# docker 脚本安装
curl -sSL https://get.daocloud.io/docker | sh

#docker compose 脚本安装
curl -L https://get.daocloud.io/docker/compose/releases/download/v2.4.1/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose


#可执行权限
sudo chmod +x /usr/local/bin/docker-compose
#创建软链：
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
#测试是否安装成功
docker-compose --version

```

## 下载安装

```sh
git clone https://gitlab.sanjiang.com/it-group/ruoyi-cloud.git
```

## 编译

```sh
cd ruoyi-cloud
mvn clean install -DskipTests
```

## 复制 jar 包

```sh
cd ./docker
./copy.sh
```

## 部署 docker

```sh
./deploy.sh base
./deploy.sh modules
```

## 检查 docker

```sh
docker ps -a | grep ruoyi
docker logs -f ruoyi-auth
docker logs -f ruoyi-gateway
docker logs -f ruoyi-modules-system
```
