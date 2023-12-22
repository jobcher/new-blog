---
title: "SSH 通过 443 端口连接 GitHub"
date: 2023-09-27
draft: false
featuredImage: "/images/github.png"
featuredImagePreview: "/images/github.png"
images: ["/images/github.png"]
authors: "jobcher"
tags: ["gitlab"]
categories: ["gitlab"]
series: ["git入门系列"]
---

# 背景
`GitHub` 提供了两种协议供用户使用 Git 连接—— SSH 和 HTTPS。理论上我可以随意选择两者之一连接到我在 GitHub 上的代码仓库，无论是将云端的仓库 clone 到本地，还是将本地的修改 push 到云端。然而，出于一些奇奇怪怪的原因，我所在的办公网络环境禁止了 `22 端口`，而 `22 端口`正是 GitHub 提供 SSH 访问的端口号。尽管可以换用 HTTPS 协议，但无论如何将我电脑上的所有代码仓库的上游都从 `git@github.com:...` 修改称 `https://github.com/...` 仍然是一个繁重的体力活。

## 解决
为了让Git也能通过上述端口用 SSH 访问 GitHub，我们为上述 SSH 连接方式设置一个别名。首先找到SSH的配置文件，它的路径一般是`~/.ssh/config`，如果这个文件不存在的话也可以创建一个。然后，在其中增加以下内容：  
```sh
Host github.com
  HostName ssh.github.com
  User git
  Port 443
```
其中，Host 是别名，HostName 是实际的域名地址，Port 是端口号。因为我希望当我在用 SSH 连接 github.com 时，实际访问的是 ssh.github.com，所以 Host 和 HostName 分别设置成这两个域名（注意不要颠倒顺序）。

## 测试连接
```sh 
ssh -vT git@github.com
```