---
title: "iOS 17 「待机显示」适配普通 iPhone（非 Pro/Max），屏幕在充电时常亮"
date: 2023-09-26
draft: false
featuredImage: "/images/Appinn-feature-images-2023-09-25T152200.586-1536x669.webp"
featuredImagePreview: "/images/Appinn-feature-images-2023-09-25T152200.586-1536x669.webp"
images: ["/images/Appinn-feature-images-2023-09-25T152200.586-1536x669.webp"]
author: "jobcher"
tags: ["daliy"]
categories: ["福利"]
series: ["福利系列"]
---
## 让 iOS 17 「待机显示」适配普通 iPhone（非 Pro/Max），屏幕在充电时常亮
## 背景
让 `iPhone 14/15 Pro/Max 以下机型`也可以在屏幕`激活待机显示功能`（充电且横置）时`保持常亮`，以显示小组件内容。  
2023 年新发布的 iOS 17 有一个新功能：待机显示 StandBy，它能在 iPhone 横向放置且充电时，全屏显示小组件，比如时钟、日历等等，但所有的`非 Pro/Max 机型`，由于没有`全天候显示屏`（`显示屏能够以低至 1 Hz 的刷新率运行`），所以并不能持续使用待机显示，会在几秒钟之后熄灭屏幕，当感受到震动（轻轻拍一下桌子），或有新通知时，再次点亮。  
## 原理
![standby-yuanli](/images/photo_2023-09-25_15-52-04.jpg)  
- 如果 iPhone 接入充电器充电 > 等待 19 秒 > 开关一次低电量模式 > 再运行一次脚本
- 如果没有充电 > 关闭低电量模式 > 停止脚本
之后，使用自动化功能，当 iPhone 接入充电器时，自动运行这个脚本就行了。  
## 教程
## [获取链接](https://www.icloud.com/shortcuts/f50d508c7ec9471a9fb94a3b2b57f1af)  
  
自带19秒等待间隔，如果遇到熄屏可以修改为15秒。
{{< bilibili BV11N4y1Z7Hk >}}  
## 注意
- 小组件自动轮换失效
- 烧屏：请务必小心，长期使用可能会导致的屏幕问题
- 充电前需要手动熄屏，否则会因为不断的开关低电量模式而无法进入待机显示
