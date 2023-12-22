---
title: "linux服务器进程打开文件过多导致服务异常"
date: 2023-11-07
draft: false
authors: "jobcher"
featuredImage: '/images/elasticsearch.png'
featuredImagePreview: '/images/elasticsearch.png'
images: ['/images/elasticsearch.png']
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---
## 背景
我的logstash是多管道部署结果发现有大量日志丢失的情况查看logstash日志出现了以下报错

>[2023-11-07T09:43:38,492][ERROR][logstash.javapipeline    ][log] Pipeline worker error, the pipeline will be stopped {:pipeline_id=>"log", :error=>"/var/lib/logstash/queue/log/checkpoint.head.tmp (Too many open files)", :exception=>Java::JavaIo::FileNotFoundException, :backtrace=>["java.base/java.io.FileOutputStream.open0(Native Method)", "java.base/java.io.FileOutputStream.open(FileOutputStream.java:298)", "java.base/java.io.FileOutputStream.<init>(FileOutputStream.java:237)", "java.base/java.io.FileOutputStream.<init>(FileOutputStream.java:187)", "org.logstash.ackedqueue.io.FileCheckpointIO.write(FileCheckpointIO.java:105)", "org.logstash.ackedqueue.Page.headPageCheckpoint(Page.java:202)", 

**这个问题是 Logstash Pipeline 在处理数据时报错,原因是打开文件过多导致"Too many open files"**  
  
## 解决方法
#### 1. 检查操作系统的文件打开数量限制,使用`ulimit -n`查看。如果太低,可以提高这个限制
打开 /etc/profile 增加ulimit 值
```sh
vim /etc/profile
## 增加，保存并退出
ulimit -n 10240
# 重载配置
source /etc/profile
```

#### 2. 适当增大Logstash的heap size,如-Xms和-Xmx设置为2g。
```sh
vim /etc/logstash/jvm.option
# 修改参数
-Xms 2g
-Xmx 2g
# 重启logstash服务
```