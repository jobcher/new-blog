---
title: "iStoreOS(旁路由)使用openclash实现dns劫持"
date: 2025-04-07
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
我主要介绍通过openwrt中的openclash覆写hosts实现dns劫持的方法。  
  
## 操作
1. 进入openwrt后台，进入openclash，进入覆写设置，进入dns设置  
![openclash](/images/openclash-dns01.png)  

2. 勾选hosts，并写入hosts信息,`10.10.10.6`是内网nas的ip地址，你可以改成任意的IP  
```sh
'nas.com': 10.10.10.6
'*.nas.com': 10.10.10.6
```  
![openclash](/images/openclash-dns02.png)  
3. 保存并应用
![openclash](/images/openclash-dns03.png)  

4. 测试
```sh
ping nas.com
ping www.nas.com
```
![openclash](/images/openclash-dns04.png)  

## 总结
openclash实现dns劫持的方法非常简单，只需要在openclash中配置hosts即可。