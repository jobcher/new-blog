---
title: "解决Elasticsearch索引只读（read-only）"
date: 2023-09-20
draft: false
author: "jobcher"
featuredImage: '/images/elasticsearch.png'
featuredImagePreview: '/images/elasticsearch.png'
images: ['/images/elasticsearch.png']
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---

# 背景
这两天有开发向我反馈说elasticsearch有报错，嘿，我定睛一看，这不是进入只读状态了，看来是存储达到额度，我马上加个新的数据节点，平衡一下存储压力  
报错信息：  
```sh
Elasticsearch Error {type:cluster_block_exception,reason:”blocked by: [FORBIDDEN/12/index read-only / allow delete (api)];}
```
## 新建服务器，安装elasticsearch
为了和之前的服务器一样，我简单写一下我elasticsearch版本和服务器系统版本  
|软件|版本|
|:----|:----|
|centos|7.9|
|elasticsearch|6.7.2|
|JDK|1.8.61|
|内存|32G|
  
### 安装和配置elasticsearch
使用rpm 安装  
```bash
wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-6.7.2.rpm
```
```bash
rpm --install elasticsearch-6.7.2.rpm
```
配置参数，进入`/etc/elasticsearch`目录  
修改配置`vim elasticsearch.yml`  
```yaml
# ======================== Elasticsearch Configuration ========================= 

cluster.name: cluster-prod-es # 集群名称

node.name: node-x # 节点名称

path.data: /var/lib/elasticsearch # 数据存储

path.logs: /var/log/elasticsearch  # 日志存储

network.host: 192.168.0.170 # 主机IP地址 

http.port: 9200 # 端口号

discovery.zen.ping.unicast.hosts: ["192.168.0.171", "192.168.0.172", "192.168.0.173"] # 集群节点

discovery.zen.minimum_master_nodes: 2 #防止脑裂
```
修改配置`vim jvm.options`  
```yaml
-Xms16g
-Xmx16g
```
启动elasticsearch 服务  
```sh
systemctl start elasticsearch.service
```
### 等待数据节点自动调节
这里要等待一会儿数据节点自动调节，这个调节时间取决于你数据的大小，一般来说几个小时也好了，如果数据重要性不太高的话，和领导沟通一下，就别傻坐着等他迁移完，正常下班就行了。
### 关闭索引只读
对了，千万不要忘记关闭只读状态，虽然你新增了节点，但是当前的只读状态并没有关闭，所以要执行一下命令关闭只读状态  
```sh
curl -XPUT -H "Content-Type: application/json" http://192.168.0.170:9200/_all/_settings -d '{"index.blocks.read_only_allow_delete": null}'
```
或者你在head-elasticsearch上操作  
![head-es](/images/head-es.jpg)  
## 感谢你的阅读