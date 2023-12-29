---
title: "docker 和 docker-compose 安装"
date: 2021-12-28
draft: false
authors: "jobcher"
tags: ["docker"]
categories: ["docker"]
series: ["docker入门系列"]
---

# 安装 docker

通过 docker 脚本安装

```sh
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
```
```sh
curl -sSL https://get.daocloud.io/docker | sh
```

## docker-compose 安装

```sh
#下载安装
sudo curl -L "https://github.jobcher.com/gh/https://github.com/docker/compose/releases/download/v2.2.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
#可执行权限
sudo chmod +x /usr/local/bin/docker-compose
#创建软链：
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
#测试是否安装成功
docker-compose --version
```

## docker 命令

常用 docker 命令

```sh
    #查看容器
    docker ps
    #查看镜像
    docker images
    #停止当前所有容器
    docker stop $(docker ps -aq)
    #删除当前停止的所有容器
    docker rm $(docker ps -aq)
    #删除镜像
    docker rmi nginx
```
