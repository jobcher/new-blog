---
title: "linux系统开启root权限"
date: 2021-12-12
draft: false
author: "jobcher"
tags: ["linux"]
categories: ["日常"]
series: ["日常系列"]
---

# linux 系统开启 root 权限

1. 修改 ssh 服务配置文件

```sh
sudo su -
sudo vim /etc/ssh/sshd_config
```

2. 增加权限  
   在# Authentication: 下输入

```bash
PermitRootLogin yes
```

3. 更改 root 密码，重启服务

```sh
sudo passwd root
service sshd restart
```
