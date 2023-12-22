---
title: "centos7.9 网络配置"
date: 2022-01-11
draft: false
authors: "jobcher"
tags: ["linux"]
categories: ["日常"]
series: ["日常系列"]
---

# centos7.9 网络配置

解决 centos 新机器网络不通的问题，CentOS7 默认不启动网卡的。CentOS 安装成功后,进行一下 ping 的操作,验证网络是否联通.

```sh
ping 1.1.1.1
ip addr
# 查看ip网络名称
```

1. 启用网卡
   进入 /etc/sysconfig/network-scipts 文件夹下，找到 IP 网卡名称

```sh
cd /etc/sysconfig/network-scipts
vim ifcfg-eth0
```

2. 启用 ONBOOT

```sh
#vim ifcfg-eth0
#修改
ONBOOT=YES
# esc 并:wq退出保存
```

3. 重启机器

```sh
shutdown -r now
```

## 结尾

centos 用的挺别扭，不考虑性能和性价比，我还是喜欢用 ubuntu……，简单的配置，初学者我建议还是先用 ubuntu，会少踩很多坑。当然了，用 x86 不然初学者用树莓派和 arm 设备，会碰到很多兼容性的问题。
