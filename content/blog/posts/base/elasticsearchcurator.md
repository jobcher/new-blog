---
title: "使用 ElasticSearch Curator 7天定期删除日志"
date: 2023-11-06
draft: false
author: "jobcher"
tags: ["运维"]
featuredImage: "/images/elasticsearch.png"
featuredImagePreview: "/images/elasticsearch.png"
images: ["/images/elasticsearch.png"]
categories: ["基础"]
series: ["基础知识系列"]
---
# 使用 ElasticSearch Curator 7天定期删除日志
## 背景
Curator 是 Elastic 官方发布的一个管理 Elasticsearch 索引的工具，可以完成许多索引生命周期的管理工作。  
我使用的 elasticseraech 8.0 以上的版本，所有我直接安装最新版的curator,服务器是centos 7 的

## 二进制安装
### 下载
```sh
wget https://packages.elastic.co/curator/5/centos/7/Packages/elasticsearch-curator-5.8.4-1.x86_64.rpm
```
### 安装 curator
```sh
rpm -ivh elasticsearch-curator-5.8.4-1.x86_64.rpm
curator --version
```

### 进入安装文件，创建文件
```sh
cd /opt/elasticsearch-curator
mkdir log
cd log
touch run.log
```

### 创建`config.yml`文件在`log`目录下
config.yml样例如下： 配置说明参考官网说明：[config.yml](https://www.elastic.co/guide/en/elasticsearch/client/curator/8.0/configfile.html)
```yml
# Rmember, leave a key empty if there is no value.  None will be a string,
# not a Python "NoneType"
client:
  hosts: 
    - 192.168.10.17  # elasticsearch IP 地址
  port: 9200
  url_prefix:
  use_ssl: False
  certificate:
  client_cert:
  client_key:
  ssl_no_validate: False
  http_auth: elastic:password # elastic 密码，没有就不用写
  timeout: 30
  master_only: False

logging:
  loglevel: INFO
  logfile: /opt/elasticsearch-curator/log/run.log
  logformat: default
  blacklist: ['elasticsearch', 'urllib3']
```

### 创建 `elk-7-action.yml` 执行 7天自动删除所有日志
aelk-7-action.yml 样例如下： 配置说明参考官网说明：[action.yml](https://www.elastic.co/guide/en/elasticsearch/client/curator/8.0/actionfile.html)
```yml
# Remember, leave a key empty if there is no value.  None will be a string,
# not a Python "NoneType"
#
# Also remember that all examples have 'disable_action' set to True.  If you
# want to use this action as a template, be sure to set this to False after
# copying it.
actions:
  1:
    action: delete_indices
    description: >-
      Delete indices older than 7 days (based on index creation_date)
    options:
      timeout_override:
      continue_if_exception: False
      disable_action: False
    filters:
    - filtertype: age
      source: creation_date 
      direction: older 
      unit: days
      unit_count: 7
```
### 执行
```sh
curator --config /opt/elasticsearch-curator/log/config.yml /opt/elasticsearch-curator/log/elk-7-action.yml
```

### 定时执行
```sh
crontab -e
0 0 * * * curator --config /opt/elasticsearch-curator/log/config.yml /opt/elasticsearch-curator/log/elk-7-action.yml
```
wq 保存定时任务

## 总结
curator适用于基于时间或者template其他方式创建的索引，不适合单一索引存储N久历史数据的操作的场景。