---
title: "安装 docker 出现 ERROR: Unsupported distribution 'ol' 问题"
date: 2022-01-11
draft: false
authors: "jobcher"
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---

# 安装 docker 出现 ERROR: Unsupported distribution 'ol' 问题

部署 docker 安装出现 ERROR: Unsupported distribution 'ol'

1. 确认是不是 arm 架构

```sh
uname -r
```

2. 确认使用的是不是 oracle 服务器系统,如果是请继续操作，安装依赖：

```sh
dnf install -y dnf-utils zip unzip
dnf config-manager --add-repo=https://download.docker.com/linux/centos/docker-ce.repo
```

3. 安装 docker

```sh
dnf remove -y runc
dnf install -y docker-ce --nobest
```

4. 完成 docker 安装并检查

```sh
systemctl enable docker.service
systemctl start docker.service
```

```sh
#检查
systemctl status docker.service
docker info
docker version
```

## 结尾

该问题主要是 oracle 没有支持依赖导致的~oracle 还是很不错的~
