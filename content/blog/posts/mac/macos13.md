---
title: "macOS 13 升级 软件失效"
date: 2022-10-31
draft: false
featuredImage: "/images/macOS_13.jpg"
featuredImagePreview: "/images/macOS_13.jpg"
author: "jobcher"
tags: ["brew"]
categories: ["Mac"]
series: ["Mac系列"]
---

# 关于 macOS 13 软件失效

Warning: You are using macOS 13.  
We do not provide support for this pre-release version.  
You will encounter build failures with some formulae.  
Please create pull requests instead of asking for help on Homebrew's GitHub,
Twitter or any other official channels. You are responsible for resolving
any issues you experience while you are running this
pre-release version.

> 简单来说就是 macOS13 版本 暂时不提供技术支持

## 解决方法

升级完 macos13 之后发现了比较麻烦的问题，很多软件出现了不兼容，这真的很无奈，对于我们这些做 IT 的人来说，这是致命的。我以 git 软件举例，有以下几个方法。

### 1. 使用时间机器恢复备份

> 使用前提：你之前备份了系统，并且系统正常
> 这个方法更加一劳永逸，因为我们并不确定还有什么软件不支持 macos13

### 2. 重新安装 xcode-select

```sh
xcode-select --install
```
