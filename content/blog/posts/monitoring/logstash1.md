---
title: "Logstash 自动重载配置文件"
date: 2022-10-31
draft: false
featuredImage: "/images/lostash.png"
featuredImagePreview: "/images/lostash.png"
author: "jobcher"
tags: ["logstash"]
categories: ["监控"]
series: ["运维监控系列"]
---

## 工作原理

1. 检测到配置文件变化
2. 通过停止所有输入停止当前`pipline`
3. 用新的配置创建一个新的管道
4. 检查配置文件语法是否正确
5. 检查所有的输入和输出是否可以初始化
6. 检查成功使用新的 pipeline 替换当前的`pipeline`
7. 检查失败,使用旧的继续工作.
8. 在重载过程中,jvm 没有重启.

## Logstash 自动重新加载配置

为了可以`自动检测`配置文件的变动和自动重新加载配置文件,需要在启动的时候使用以下命令:

```sh
./bin/lagstash -f configfile.conf --config.reload.automatic
```

> 启动 Logstash 的时候使用`--config.reload.automatic`或`-r`选项来开启自动重载配置。

## 修改检测间隔时间

默认检测配置文件的间隔时间是`3秒`,可以通过以下命令改变

```sh
--config.reload.interval <second>
```

如果 Logstash 已经运行并且没有开启自动重载，你可以强制 Logstash 重新载入配置文件并且重启管道通过发送一个 SIGHUP 信号。比如：

```sh
kill -1 <pid>
```

其中<pid>是正在运行的 Logstash 的进程号。

## 注意！！！

> `stdin`输入插件不支持自动重启.  
> `syslog`作为输入源,当重载配置文件时,会崩溃.  
> [解决方法](https://github.com/logstash-plugins/logstash-input-syslog/issues/40)
