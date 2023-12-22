---
title: "docker 命令(2)"
date: 2022-02-18
draft: false
authors: "jobcher"
tags: ["docker"]
categories: ["docker"]
series: ["docker进阶系列"]
---

# docker 命令(2)

## `docker ps` 命令

`docker ps` 能查看所有运行中的容器  
`docker ps -a` 能查看所有的容器  
`docker rm -f $(docker ps -aq)` 强制删除所有容器

## `docker run`和`docker create`有什么区别

`docker create`命令能够基于镜像创建容器。  
该命令执行的效果类似于`docker run -d`，即创建一个将在系统后台运行的容器。  
但是与`docker run -d`不同的是，`docker create`创建的容器`并未实际启动`，还需要执行`docker start`命令或`docker run`命令以启动容器。  
事实上，`docker creat`e 命令常用于`在启动容器之前进行必要的设置`。
