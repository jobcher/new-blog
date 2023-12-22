---
title: "k8s本地联调神器kt-connect"
date: 2022-04-14
draft: false
author: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---

# k8s 本地联调神器 kt-connect

[转载自 Bboysoul'sBlog](https://www.bboy.app/2022/04/11/k8s%E6%9C%AC%E5%9C%B0%E8%81%94%E8%B0%83%E7%A5%9E%E5%99%A8kt-connect/)  
k8s 集群内部的服务网络怎么和我们本地网络打通。kt-connect 就是用来解决这个问题的

## 使用方法

下载安装什么的都很简单，一个二进制而已

```sh
https://github.com/alibaba/kt-connect
```

如果你安装好了，那么直接使用下面的命令使用就好了

```sh
sudo ktctl connect
```

当然也可以指定配置文件

```sh
sudo ktctl --kubeconfig ~/.kube/local connect
```

执行完成之后，这个集群的所有`svc`都可以直接在本地解析，当然直接 ping pod 的 ip 也是可以的
