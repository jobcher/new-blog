---
title: SD-webui 批量处理图片
date: 2024-01-25
draft: false
author: 'jobcher'
featuredImage: '/images/sd-logo.jpeg'
featuredImagePreview: '/images/sd-logo.jpeg'
images: ['/images/sd-logo.jpeg']
tags: ['stable diffusion']
categories: ['stable diffusion']
series: ['stable diffusion']
---
## 背景
`Stable Diffusion` 在训练数据集之前，需要先对数据进行预处理。  
本篇文章就是介绍如何对图像进行批量预处理。
## 图片上传
上传图像到你指定目录，我的目录时`/mnt/smb/`  
![上传共享文件夹](/images/1706143682775.jpg)  
打开`SD-web`地址，进入[192.168.1.232:7861](192.168.1.232:7861)，选择`附加功能`，进行图像预处理  
![图像处理](/images/1706143565210.jpg)  
## 批量抠图
### 选择`从目录进行批量处理`  
[从目录进行批量处理](/images/1706144053842.jpg)  
### 填写`输入目录`和`输出目录`  
举例：  
- 原本我的共享文件夹 地址是`\\192.168.1.249\DB Training\ai-pre-photo\out-photo` 将所有的 `\\192.168.1.249\DB Training\` 改为 `/mnt/smb/`
- 所有的 `\` 改为 `/`
因此 需填写地址如下图：  
![输入目录](/images/1706144668484.jpg)  
### 选择抠图模型
滑倒最底部选择`背景去除算法`,选择你要使用的算法，我这边选择`silueta`算法,可以根据你自己的需求使用算法  
![选择模型](/images/1706144851240.jpg)  
### 执行生成
点击`生成`按钮，开始对图像批量处理  
![执行生成](/images/1706144948763.jpg)  
### 查看生成的图像
![执行生成](/images/1706145130721.jpg)  
>在对应的`共享文件夹`也可以查看  