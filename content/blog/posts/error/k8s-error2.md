---
title: "k8s CNI 问题 连接认证失效"
date: 2022-09-23
draft: false
author: "jobcher"
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---

# k8s CNI 问题 连接认证失效

删除 `calico` 换成 `flannel` 后，容器没有正常启动  
network: error getting ClusterInformation: connection is unauthorized: Unauthorized]

## 解决问题

删除掉 `/etc/cni/net.d/` 目录下的 calico 配置文件即可。  
要删除`所有节点`的配置文件

```sh
sudo rm -rf /etc/cni/net.d/*calico*
```

## 不要重复网络插件

![img](https://pic.jitudisk.com/public/2022/09/22/a4d8cc49f1c84.jpg)
