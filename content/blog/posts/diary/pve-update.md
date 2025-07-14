---
title: "Proxmox VE（PVE） 更新到最新版本"
date: 2025-07-14
draft: false
authors: "jobcher"
featuredImage: "/images/Proxmox-logo-860.png"
featuredImagePreview: "/images/Proxmox-logo-860.png"
images: ["/images/Proxmox-logo-860.png"]
tags: ["pve"]
categories: ["日常"]
series: ["日常系列"]
---
## 背景
由于增加了新的pve机器组了pve集群，这过程中发生很多事情，打算记录一下过程中发生的问题。  
## 如何不使用企业源完成pve更新  
### 修改订阅源
![pve 订阅源](/images/pve-update-1.png)  
禁用企业源  
![pve 禁用企业源](/images/pve-update-2.png)  
添加 no-subscription 源  
![pve 添加 no-subscription 源](/images/pve-update-3.png)  
![pve 确认 no-subscription 源](/images/pve-update-4.png)  
### 执行命令
到具体的pve机器上执行命令
```sh
apt update
apt dist-upgrade
# 查看版本
pveversion
```

## HA 集群无法正常运行
出现`lrm pve (old timestamp - dead?)`获取不到pve主机状态  
查看 lrm 服务状态
```sh
systemctl status pve-ha-lrm
systemctl status pve-ha-crm
```
重启 lrm 服务
```sh
# 重启 lrm 服务
systemctl restart pve-ha-lrm
# 重启 crm 服务
systemctl restart pve-ha-crm
```