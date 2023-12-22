---
title: "rsync 文件同步"
date: 2021-12-20
draft: false
authors: "jobcher"
tags: ["daliy", "rsync"]
categories: ["日常"]
series: ["日常系列"]
---

# rsync 文件同步

rsync 是一个常用的 Linux 应用程序，用于文件同步

## 安装

```sh
# Debian or Ubuntu
$ sudo apt-get install rsync

# Red Hat
$ sudo yum install rsync

# Arch Linux
$ sudo pacman -S rsync
```

## 基本用法

使用 rsync 命令时，可以作为 cp 和 mv 命令的替代方法，将源目录同步到目标目录。  
-r 表示递归，即包含子目录。注意，-r 是必须的，否则 rsync 运行不会成功。source 目录表示源目录，destination 表示目标目录。  
-a 参数可以替代-r，除了可以递归同步以外，还可以同步元信息（比如修改时间、权限等）。由于 rsync 默认使用文件大小和修改时间决定文件是否需要更新

```sh
rsync -r source destination
```

远程同步

```sh
rsync -av <源地址>/ <用户名>@<ip地址>:/<目标地址>
```

1. 友情地址：[mysql 迁移](https://www.jobcher.com/mysqldump/)
