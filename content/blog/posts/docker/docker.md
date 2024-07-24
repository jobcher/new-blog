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

## 通过 docker 脚本安装

```sh
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
```
```sh
curl -sSL https://get.daocloud.io/docker | sh
```

### CentOS 手动安装
卸载之前相关的依赖
```sh
sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine
```
下载依赖包
```sh
sudo yum install -y yum-utils
```
选择国内`阿里云`镜像
```sh
sudo yum-config-manager \
    --add-repo \
    https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
```
安装 Docker Engine-Community
```sh
sudo yum install docker-ce docker-ce-cli containerd.io docker-compose-plugin -y
```
### Ubuntu 手动安装
卸载之前相关的依赖
```sh
sudo apt-get remove docker docker-engine docker.io containerd runc
```
更新 apt 包索引
```sh
apt-get update
```
```sh
sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
```
```sh
curl -fsSL https://mirrors.ustc.edu.cn/docker-ce/linux/ubuntu/gpg | sudo apt-key add -
```
选择国内镜像
```sh
sudo add-apt-repository \
   "deb [arch=amd64] https://mirrors.ustc.edu.cn/docker-ce/linux/ubuntu/ \
  $(lsb_release -cs) \
  stable"
```
```sh
sudo apt-get update
```
安装 Docker Engine-Community
```sh
sudo apt-get install docker-ce docker-ce-cli containerd.io
```

## docker-compose 安装

```sh
#下载安装
sudo curl -L "https://github.jobcher.com/gh/https://github.com/docker/compose/releases/download/v2.29.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
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
