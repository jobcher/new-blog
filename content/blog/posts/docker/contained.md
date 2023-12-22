---
title: "contained 安装及使用"
date: 2022-08-17
draft: false
author: "jobcher"
tags: ["docker"]
categories: ["docker"]
series: ["docker入门系列"]
---

# contained 安装及使用

![contained](/images/containerd-horizontal-color.png)  
containerd 是一个行业标准的容器运行时，强调简单性、健壮性和可移植性。它可作为 Linux 和 Windows 的守护进程使用，可以管理其主机系统的完整容器生命周期：图像传输和存储、容器执行和监督、低级存储和网络附件等。  
`containerd` is a member of `CNCF` with `graduated` status.

1. 早在 2016 年 3 月，`Docker 1.11`的`Docker Engine`里就包含了`containerd`，而现在则是把`containerd`从`Docker Engine`里彻底剥离出来，作为一个独立的开源项目独立发展，目标是提供一个更加开放、稳定的容器运行基础设施。和原先包含在 Docker Engine 里`containerd`相比，独立的`containerd`将具有更多的功能，可以涵盖整个容器运行时管理的所有需求。
2. `containerd`并不是直接面向最终用户的，而是主要用于集成到更上层的系统里，比如`Swarm`, `Kubernetes`, `Mesos`等容器编排系统。
3. `containerd`以`Daemon`的形式运行在系统上，通过暴露底层的`gRPC API`，上层系统可以通过这些`API`管理机器上的容器。
4. 每个`containerd`只负责`一台机器`，Pull 镜像，对容器的操作（启动、停止等），网络，存储都是由 containerd 完成。具体运行容器由`runC`负责，实际上只要是符合`OCI规范`的容器都可以支持。
5. 对于容器编排服务来说，运行时只需要使用`containerd+runC`，更加轻量，容易管理。 5.独立之后`containerd`的特性演进可以和`Docker Engine`分开，专注容器运行时管理，可以更稳定。

## 安装

centos

```sh
yum install -y containerd.io
```

ubuntu

```sh
apt install -y containerd.io
```

设置开机自启

```sh
systemctl enable containerd
systemctl start containerd
systemctl status containerd
```

验证

```sh
ctr version
```

## ctr 命令

| 命令                      | 作用                                 |
| :------------------------ | :----------------------------------- |
| plugins, plugin           | 提供有关容器插件的信息               |
| version                   | 打印客户端和服务器版本               |
| containers, c, container  | 管理容器                             |
| content                   | 管理内容                             |
| events, event             | 显示容器事件                         |
| images, image, i          | 管理图像                             |
| leases                    | 管理租约                             |
| namespaces, namespace, ns | 管理租命名空间                       |
| pprof                     | 为 containerd 提供 golang pprof 输出 |
| run                       | 运行一个容器                         |
| snapshots, snapshot       | 管理快照                             |
| tasks, t, task            | 管理任务                             |
| install                   | 安装一个新包                         |
| oci                       | OCI 工具                             |
| shim                      | 直接与 shim 交互                     |
| help, h                   | 显示命令列表或一个命令的帮助         |
