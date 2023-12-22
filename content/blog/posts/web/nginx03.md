---
title: "nginx.conf 配置文件详解"
date: 2021-12-24
draft: false
authors: "jobcher"
tags: ["nginx", "linux"]
categories: ["web 服务器"]
series: ["web服务入门系列"]
---

# nginx.conf 配置文件详解

```sh
# vim nginx.conf
user nobody nobody; # 运行 nginx 的所属组和所有者
worker_processes 2; # 开启两个 nginx 工作进程,一般几个 CPU 核心就写几
error_log logs/error.log notice; # 错误日志路径
pid logs/nginx.pid; # pid 路径

events
{
    worker_connections 1024; # 一个进程能同时处理 1024 个请求
}

http
{
    include mime.types;
    default_type application/octet-stream;

    log_format main ‘$remote_addr – $remote_user [$time_local] “$request” ‘
    ‘$status $body_bytes_sent “$http_referer” ‘
    ‘”$http_user_agent” “$http_x_forwarded_for”‘;
    access_log logs/access.log main; # 默认访问日志路径
    sendfile on;
    keepalive_timeout 65; # keepalive 超市时间
    # 开始配置一个域名,一个 server 配置段一般对应一个域名
    server
    {
        listen 80; #
        # 在本机所有 ip 上监听 80,也可以写为 192.168.1.202:80,这样的话,就只监听 192.168.1.202 上的 80 口
        server_name www.nbtyfood.com; # 域名
        root /www/html/www.nbtyfood.com; # 站点根目录（程序目录）
        index index.html index.htm; # 索引文件
        location /
        { # 可以有多个 location
            root /www/html/www.nbtyfood.com; # 站点根目录（程序目录）
        }
        error_page 500 502 503 504 /50x.html;
        # 定义错误页面,如果是 500 错误,则把站点根目录下的 50x.html 返回给用户
        location = /50x.html {
        root /www/html/www.nbtyfood.com;
        }
    }
# 开始配置站点 bbs.nbtyfood.com
    server
    {
        listen 80;
        server_name bbs.nbtyfood.com;
        root /www/html/bbs.nbtyfood.com;
        index index.html index.htm; # 索引文件
        location /
        {
            root /www/html/bbs.nbtyfood.com;
        }
        error_page 500 502 503 504 /50x.html;
        location = /50x.html {
            root /www/html/bbs.nbtyfood.com;
        }
    }
}

```
