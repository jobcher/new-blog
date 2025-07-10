---
title: "kindle 邮件发送失败，错误代码E999"
date: 2025-07-10
draft: false
authors: "jobcher"
featuredImage: '/images/kindle.png'
featuredImagePreview: '/images/kindle.png'
images: ['/images/kindle.png']
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---
## 背景
最近看到了一本电子书《just for fun》 Linux之父：林纳斯的自传。就想上传到kindle上去看这本书，但是之前epub的上传都没有问题，唯独这本一直上传失败。E999 故障。我直接说我对于这个问题的解决办法。  
![e999](/images/e999.png)  
## 解决方法
1. 将epub转换为mobi格式  
2. 将装换好的mobi文件再转换为新的epub文件  
3. 上传新的epub文件到kindle  
## 总结
感觉是原本的epub文件格式有损坏导致，amzon无法正确识别epub格式，所有按照这个方式转换一遍，基本上上传就没有问题了。  