---
title: "linux常用命令"
date: 2022-02-12
draft: false
authors: "jobcher"
tags: ["运维"]
categories: ["基础"]
series: ["基础知识系列"]
---

# linux 常用命令

## 软件操作命令

```sh
#软件包管理器
yum
# 安装软件
yum install xxxx
# 卸载软件
yum remove xxx
# 搜索软件
yum search xxx
# 清理缓存
yum clean packages
# 列出已安装
yum list
# 软件包信息
yum info
```

## 服务器硬件资源和磁盘操作

```sh
# 内存
free -h
# 硬盘
df -h
# 负载
w/top/htop
# 查看cpu
cat /proc/cpuinfo
# 查看磁盘
fdisk -l
```

## 文件和文件夹操作命令

| 命令  | 解释             |
| :---- | :--------------- |
| ls    | 查看目录下的文件 |
| touch | 新建文件         |
| mkdir | 新建目录         |
| cd    | 进入目录         |
| rm    | 删除文件和目录   |
| cp    | 复制             |
| mv    | 移动             |
| pwd   | 显示路径         |

## 系统用户操作命令

## 防火墙相关设置

## 提权操作 sudo 和文件传输
