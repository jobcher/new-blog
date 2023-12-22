---
title: "k8s.gcr.io国内无法连接解决方法"
date: 2022-09-22
draft: false
authors: "jobcher"
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---

# k8s.gcr.io 国内无法连接解决方法

Get https://k8s.gcr.io/v2/: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)  
这个一看知道什么原因了，应该 GFW！那好吧，只能给 docker 加个代理了。

## 解决问题

添加 mirror 站点

```sh
registry.cn-hangzhou.aliyuncs.com/google_containers
```
