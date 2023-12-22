---
title: "nginx 配置和编译"
date: 2023-04-20
draft: false
featuredImage: "/images/nginx.png"
featuredImagePreview: "/images/nginx.png"
author: "jobcher"
tags: ["nginx", "linux"]
categories: ["web 服务器"]
series: ["web服务入门系列"]
---
## nginx 自编译
[nginx官网](https://nginx.org/en/)  
>下载nginx->Configure->编译->安装
### 下载nginx
```sh
wget https://nginx.org/download/nginx-1.24.0.tar.gz
```
解压目录
```sh
tar -zxvf nginx-1.24.0.tar.gz
```

### 编译最新的 nginx,并备份旧的二进制文件
```sh
cd nginx-1.24.0
./configure && make
mv /path/to/nginx /path/to/nginx.bak
```

### 覆盖旧的二进制文件
```sh
cp /path/to/nginx /path/to/nginx
```

### 执行平滑重启
```sh
nginx -s reopen
```

### 观察nginx服务是否有中断,如果正常则删除备份文件
```sh
rm /path/to/nginx.bak
```