---
title: "ProXmoX VE升级 apt-get update 报错"
date: 2022-01-30
draft: false
authors: "jobcher"
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---

# ProXmoX VE 升级 apt-get update 报错

## 解决方法

```sh
vim /etc/apt/sources.list.d/pve-enterprise.list
#注释掉
#deb https://enterprise.proxmox.com/debian/pve stretch pve-enterprise
```

## 添加内容

```sh
echo "deb http://download.proxmox.com/debian/pve stretch pve-no-subscription" > /etc/apt/sources.list.d/pve-install-repo.list
wget http://download.proxmox.com/debian/proxmox-ve-release-5.x.gpg -O /etc/apt/trusted.gpg.d/proxmox-ve-release-5.x.gpg
```

## 更新系统

```sh
apt update && apt dist-upgrade
```

## 结尾

升级完成后，可以执行`pveversion -v`查看下最新的软件版本。然后执行`reboot`重启物理服务器
