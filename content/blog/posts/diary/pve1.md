---
title: "Proxmox VE 在线扩容磁盘分区"
date: 2022-01-19
draft: false
author: "jobcher"
tags: ["扩容磁盘", "pve"]
categories: ["日常"]
series: ["日常系列"]
---

# Proxmox VE 在线扩容磁盘分区

## 添加磁盘大小

![pve1-1](/images/pve1.png)

## 在 VM 上做扩容操作

1. 安装 growpart

```sh
yum install -y epel-release
yum install -y cloud-utils
```

2. 查看系统盘 路径

```sh
fdisk -l
df -h
```

![pve1-2](/images/pve2.png)

3. 扩容设备并重启

```sh
growpart /dev/sda 2 #2代表是第二块系统分区，不是sda2,中间有空格
reboot
```

4. 重启执行命令

```sh
xfs_growfs /dev/sda2 #(xfs 文件系统)
resize2fs /dev/sda2 #(ext4 文件系统)

```

5. 更新完成

```sh
df -h
```

## 逻辑卷没有正常扩容的情况
1. 检查当前逻辑卷属于哪个卷组:
```sh
vgdisplay
```
![pve1-1.png](/images/pve1-1.png)  
检查卷组中是有足够的空间可以扩容,还有99g  

2. 扩展逻辑卷大小到200G:
```sh
lvextend -L +99G /dev/mapper/ubuntu--vg-ubuntu--lv
```
![pve1-2.png](/images/pve1-2.png)  

3. 调整文件系统大小到逻辑卷大小:
```sh
resize2fs /dev/mapper/ubuntu--vg-ubuntu--lv
```
![pve1-3.png](/images/pve1-3.png)  
4. 检查
```sh
df -h /dev/mapper/ubuntu--vg-ubuntu--lv
```
成功扩容