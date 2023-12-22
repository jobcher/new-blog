---
title: "Vue3 + vite + nginx项目部署后404问题"
date: 2023-08-10
draft: false
author: "jobcher"
featuredImage: '/images/nginx.png'
featuredImagePreview: '/images/nginx.png'
images: ['/images/nginx.png']
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---

## Vue3 + vite + nginx项目部署后404问题
`vue3 + vite + nginx`  
在服务器上部署后打开首页都没问题，打开其他路径全部 404。  
`nginx` 报错日志：`No such file or directory`  
  
> 其实查看 `build` 后的`dist文件夹`可以发现，只有一个`index.html`，当你访问别的路径时`nignx`查找不到所以就报错了
  
## 解决方案
在 nginx.conf 中添加: `try_files $uri $uri/ /index.html;`  
```conf
server {
          listen       80;
          server_name  localhost;
          location / {
          root   /dist;
          index  index.html index.htm;
          # 在配置文件的此处加上这句话
          try_files $uri $uri/ /index.html;
        }
    }
```

## 总结
其实上述改动就是告诉 nignx 找不到文件的时候就访问 `index.html` 就可以了。

究其原因其实就是是 vue3 的 router 使用了`history`模式，该模式与之前`hash`模式的具体区别可以自行百度一下，不在此赘述。