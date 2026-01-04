---
title: "openwrt 硬盘扩容"
date: 2025-01-04
draft: false
featuredImage: "/images/openclash.jpg"
featuredImagePreview: "/images/openclash.jpg"
images: ["/images/openclash.jpg"]
authors: "jobcher"
tags: ["daliy"]
categories: ["日常"]
series: ["日常系列"]
---
## 背景
本教程基于 `ext4` 文件系统的 OpenWrt
  
## 操作
1. OpenWrt 镜像默认的磁盘大小是比较小的，没安装几个软件就不够用了，所以需要我们手动来扩容一下：  
![openclash](/images/openwrt-kr-1.png)  
```sh
opkg update
opkg install cfdisk
```

2. 首先 PVE 下给之前的 SATA 硬盘增加 2GB 空间：
![openclash](/images/openwrt-kr-2.png)  
硬件更改后，记得重启一下这个 OpenWrt 的 VM  
然后 OpenWrt 安装 cfdisk 工具：  
![openclash](/images/openwrt-kr-3.png)  
后面扩容分区需要用到这个工具。  

3. 使用 cfdisk 来扩容分区，可以看到末尾有 2GB 空闲分区：
![openclash](/images/openwrt-kr-4.png)  
选中第二个分区，选择下面的「Resize」 调整磁盘分区大小：  
![openclash](/images/openwrt-kr-5.png)  
最后选择第二个分区，选择`「Write」` 保存我们上面的操作：  
![openclash](/images/openwrt-kr-6.png)  
此时记住我们当前的第 2 个分区路径为：`/dev/sda2` 下一步操作需要用到这个路径信息.  

4. 设置循环
OpenWrt 安装 losetup 工具：  
```sh
opkg update
opke install losetup resize2fs
```
接着设置循环设备并挂载，操作完重启一下：  
```sh
losetup /dev/loop0 /dev/sda2
resize2fs -f /dev/loop0
reboot
```

## 总结
经过上述几步操作，最终扩容成功