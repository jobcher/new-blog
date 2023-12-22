---
title: "nginx 汇总"
date: 2021-12-08
draft: false
author: "jobcher"
tags: ["nginx", "linux"]
categories: ["web 服务器"]
series: ["web服务入门系列"]
---

# nginx 汇总

各类 nginx 问题汇总

## 安装 nginx

```sh
    #centos
    yum install nginx
    #ubuntu
    apt install nginx
```

## http 代理

### 正向代理

```sh
    server {
        listen       80;
        server_name  www.nbtyfood.com;

        location / {
            proxy_pass http://127.0.0.1:8080;
        }
    }
```

### 反向代理

## 负载均衡

```sh
    upstream mysvr {
    server 192.168.10.121:3333;
    server 192.168.10.122:3333;
    }
    server {
        ....
        location  ~*^.+$ {
            proxy_pass  http://mysvr;  #请求转向mysvr 定义的服务器列表
        }
    }
```

1. 热备  
   如果你有 2 台服务器，当一台服务器发生事故时，才启用第二台服务器给提供服务。服务器处理请求的顺序：AAAAAA 突然 A 挂啦，BBBBBBBBBBBBBB.....

```sh
    upstream mysvr {
    server 127.0.0.1:7878;
    server 192.168.10.121:3333 backup;  #热备
    }
```

2. 轮询  
   nginx 默认就是轮询其权重都默认为 1，服务器处理请求的顺序：ABABABABAB....

```sh
    upstream mysvr {
    server 127.0.0.1:7878;
    server 192.168.10.121:3333;
    }
```

3. 加权轮询  
   跟据配置的权重的大小而分发给不同服务器不同数量的请求。如果不设置，则默认为 1。下面服务器的请求顺序为：ABBABBABBABBABB....

```sh
    upstream mysvr {
    server 127.0.0.1:7878 weight=1;w
    server 192.168.10.121:3333 weight=2;
    }
```

4. ip_hash
   nginx 会让相同的客户端 ip 请求相同的服务器。

```sh
    upstream mysvr {
    server 127.0.0.1:7878;
    server 192.168.10.121:3333;
    ip_hash;
    }
```

## web 缓存

```sh
    location /images/ {
    proxy_cache my_cache;
    proxy_ignore_headers Cache-Control;
    proxy_cache_valid any 30m;
    # ...
}
```

## 重定向

```sh
    rewrite ^/(.*) http://www.nbtyfood.com/$1 permanent;
```
