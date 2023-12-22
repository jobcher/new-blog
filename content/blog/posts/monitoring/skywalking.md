---
title: "skywalking APM 监控"
date: 2022-08-10
draft: false
authors: "jobcher"
tags: ["skywalking"]
categories: ["监控"]
series: ["运维监控系列"]
---

# skywalking

基于 OpenTracing 规范，专门为微服务架构以及云原生服务。

## APM 监控

一个基于微服务架构的电商系统  
<img src="https://test.jobcher.com/images/skywalkning.png" width="100%">  
`APM `(Application Performance Management) 即应用性能管理，属于 IT 运维管理（ITOM)范畴.  
分为一下三个方面：

- Logging  
  服务在处理某个请求时打印的错误日志，可以将这些日志信息记录到`Elasticsearch`或是其他存储中。通过 Kibana 或是其他工具来分析这些日志了解服务的行为和状态，大多数情况下。日志记录的数据很分散，并且相互独立。例如错误日志，请求处理过程中关键步骤的日志等等。
- Metrics  
  `Metric`是可以聚合的，例如为电商系统中每个 HTTP 接口添加一个计数器，计算每个接口的 QPS，可以通过简单的加和计算得到系统的总负载情况。
- Tracing  
  在微服务架构系统中一请求会经过很多服务处理，调用链路会非常长，要确定中间哪个服务出现异常是非常麻烦的事情，通过分布式链路追踪，运维人员就可以构建一个请求的视图。视图上战术了一个请求从进入系统开始到返回响应的整个流程。

> `系统交互图`  
> <img src="https://test.jobcher.com/images/skywalking2.svg" width="70%">

> `系统加载图` > <img src="https://test.jobcher.com/images/skywalking3.svg" width="100%">

## 目前流行的`APM监控`

- Zipkin
  - 对 web.xml 进行修改，代码侵入
  - twitter 开源
- Cat
  - 支持 Java、C/C++、Node.Js、Python、go
  - 代码侵入，埋点
  - 美团开源
- Pinpoint
  - 基于字节码注入技术，代码无侵入
  - 韩国公司开发，社区交流滞后
  - 只支持 hbase
  - 颗粒度更细
- Skywalking  
  观测性分析平台
  - 基于字节码注入技术，代码无侵入
  - 服务、服务实例、端点指标分析
  - 服务拓扑图分析
  - 服务、服务实例和端点（Endpont）SLA 分析
  - 支持 es，h2,mysql,TiDb,sharding-sphere

# skywalking 整体框架

<img src="https://test.jobcher.com/images/skywalking4.jpg" width="100%">  
  
- 上部分 `Agent` ：负责从应用中，收集链路信息，发送给 SkyWalking OAP 服务器。目前支持 SkyWalking、Zikpin、Jaeger 等提供的 Tracing 数据信息。而我们目前采用的是，SkyWalking Agent 收集 SkyWalking Tracing 数据，传递给服务器。
- 下部分 `SkyWalking OAP` ：负责接收 Agent 发送的 Tracing 数据信息，然后进行分析(Analysis Core) ，存储到外部存储器( Storage )，最终提供查询( Query )功能。
- 右部分 `Storage` ：Tracing 数据存储。目前支持 ES、MySQL、Sharding Sphere、TiDB、H2 多种存储器。而我们目前采用的是 ES ，主要考虑是 SkyWalking 开发团队自己的生产环境采用 ES 为主。
- 左部分 `SkyWalking UI` ：负责提供控台，查看链路等等。

## skywalking 配置

## 使用 docker-compose 安装

使用 mysql 作为存储  
下载 `mysql-connector-java-8.0.30.jar`

```sh
mkdir ./libs/
mv mysql-connector-java-8.0.30.jar ./libs/
```

创建带 mysql 驱动的基础镜像

```Dockerfile
FROM apache/skywalking-oap-server:9.1.0
LABEL maintainer="nb@nbtyfood.com"
COPY ./libs/* /skywalking/oap-libs
```

上传 dockerhub 或者自己的镜像仓库，这里我是上传到自己的仓库

1. 创建镜像
   > docker build -t skywalking-mysql-server:v1.0 .
2. 打 tag，选择上传位置
   > docker tag skywalking-mysql-server:v1.0 <仓库地址>/blog/skywalking-mysql-server:v1.0
3. 上传镜像
   > docker push <仓库地址>/blog/skywalking-mysql-server:v1.0

```yml
version: "3"
services:
  skywalking-oap-server:
    image: "hub.docker.com/jobcher/skywalking-mysql-server:v1.0" #docker iamge 地址
    container_name: "oap-server"
    restart: "always"
    environment:
      - SW_STORAGE=mysql
      - SW_JDBC_URL="jdbc:mysql://10.12.12.4:3306/sk"
      - SW_DATA_SOURCE_USER=user # mysql用户名
      - SW_DATA_SOURCE_PASSWORD=password # mysql密码
    ports:
      - "10.12.12.16:12800:12800"
      - "10.12.12.16:1234:1234"
      - "10.12.12.16:11800:11800"

  skywalking-oap-ui: #UI界面
    image: "apache/skywalking-ui:9.1.0"
    container_name: "oap-ui"
    restart: "always"
    environment:
      - SW_OAP_ADDRESS=http://10.12.12.16:12800
    ports:
      - "8180:8080"
```
## skywalking 客户端部署
下载客户端  
[点击下载](https://archive.apache.org/dist/skywalking/8.14.0/apache-skywalking-java-agent-8.14.0.tgz)  
### 安装到要监控的主机
```sh
wget https://archive.apache.org/dist/skywalking/8.14.0/apache-skywalking-java-agent-8.14.0.tgz
tar -zxvf apache-skywalking-java-agent-8.14.0.tgz
```
### 配置变量
```shell
# SkyWalking Agent 配置
export SW_AGENT_NAME=rf-consumer # 配置 Agent 名字。一般来说，我们直接使用 Spring Boot 项目的 `spring.application.name` 。
export SW_AGENT_COLLECTOR_BACKEND_SERVICES=10.12.12.16:11800 # 配置 Collector 地址。
#export SW_AGENT_SPAN_LIMIT=2000 # 配置链路的最大 Span 数量。一般情况下，不需要配置，默认为 300 。主要考虑，有些新上 SkyWalking Agent 的项目，代码可能比较糟糕。
export JAVA_AGENT=-javaagent:/home/ubuntu/skywalking-agent/skywalking-agent.jar # SkyWalking Agent jar 地址。
```
### 启动java程序
```sh
java $JAVA_AGENT -jar yourapp.jar
#或者
java -javaagent:<skywalking-agent-path> -jar yourApp.jar
```
