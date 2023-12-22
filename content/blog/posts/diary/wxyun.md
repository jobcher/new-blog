---
title: "网心云挂机教程 | 轻松实现睡后收入~"
date: 2022-02-11
draft: false
authors: "jobcher"
tags: ["网心云", "外快"]
categories: ["日常"]
series: ["日常系列"]
---

# 网心云挂机教程 | 轻松实现睡后收入~

首先，本文章只是分享，造成一切的后果，博主概不负责！都是成年人了……

我采用 docker 容器魔方来挂载网心云

## docker 部署

`/mnt/money/wxedge_storage`这个路径改为自己的存储路径建议`>200G`

```sh
docker run \
--name=wxedge \
--restart=always \
--privileged \
--net=host \
--tmpfs /run \
--tmpfs /tmp \
-v /mnt/money/wxedge_storage:/storage:rw \
-d \
registry.cn-hangzhou.aliyuncs.com/onething/wxedge
```

## 设备绑定

1. 进入 dockerip 地址 （http://127.0.0.1:18888)
   ![](/images/wxyun1.png)
2. 下载 app 扫码绑定
   ![](/images/wxyun2.png)

## 成功

然后坐等第二天收益到账就可以了，记得`19:00-23:00`是收益高峰期尽量保持在线~
