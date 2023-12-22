---
title: "nginx 日志格式整理"
date: 2021-12-13
draft: false
author: "jobcher"
tags: ["nginx", "linux"]
categories: ["web 服务器"]
series: ["web服务入门系列"]
---

# nginx 日志配置

## 语法

```conf
access_log path [format [buffer=size] [gzip[=level]] [flush=time] [if=condition]]; # 设置访问日志
access_log off; # 关闭访问日志
```

例子：

```conf
access_log /var/logs/nginx-access.log
access_log /var/logs/nginx-access.log buffer=32k gzip flush=1m
```

## 使用 log_format 自定义日志格式

Nginx 预定义了名为 combined 日志格式，如果没有明确指定日志格式默认使用该格式：

```conf
log_format combined '$remote_addr - $remote_user [$time_local] '
                    '"$request" $status $body_bytes_sent '
                    '"$http_referer" "$http_user_agent"';
```

如果不想使用 Nginx 预定义的格式，可以通过 log_format 指令来自定义。

### 语法

```conf
log_format name [escape=default|json] string ...;
```

| 变量                  | 含义                                                                                                                       |
| :-------------------- | :------------------------------------------------------------------------------------------------------------------------- |
| $bytes_sent           | 发送给客户端的总字节数                                                                                                     |
| $body_bytes_sent      | 发送给客户端的字节数，不包括响应头的大小                                                                                   |
| $connection           | 连接序列号                                                                                                                 |
| $connection_requests  | 当前通过连接发出的请求数量                                                                                                 |
| $msec                 | 日志写入时间，单位为秒，精度是毫秒                                                                                         |
| $pipe                 | 如果请求是通过 http 流水线发送，则其值为"p"，否则为“."                                                                     |
| $request_length       | 请求长度（包括请求行，请求头和请求体）                                                                                     |
| $request_time         | 请求处理时长，单位为秒，精度为毫秒，从读入客户端的第一个字节开始，直到把最后一个字符发送张客户端进行日志写入为止           |
| $status               | 响应状态码                                                                                                                 |
| $time_iso8601         | 标准格式的本地时间,形如“2017-05-24T18:31:27+08:00”                                                                         |
| $time_local           | 通用日志格式下的本地时间，如"24/May/2017:18:31:27 +0800"                                                                   |
| $http_referer         | 请求的 referer 地址。                                                                                                      |
| $http_user_agent      | 客户端浏览器信息。                                                                                                         |
| $remote_addr          | 客户端 IP                                                                                                                  |
| $http_x_forwarded_for | 当前端有代理服务器时，设置 web 节点记录客户端地址的配置，此参数生效的前提是代理服务器也要进行相关的 x_forwarded_for 设置。 |
| $request              | 完整的原始请求行，如 "GET / HTTP/1.1"                                                                                      |
| $remote_user          | 客户端用户名称，针对启用了用户认证的请求                                                                                   |
| $request_uri          | 完整的请求地址，如 "https://daojia.com/"                                                                                   |

例子：

```conf
access_log /var/logs/nginx-access.log main
log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                  '$status $body_bytes_sent "$http_referer" '
                  '"$http_user_agent" "$http_x_forwarded_for"';
```
