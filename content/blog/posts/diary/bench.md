---
title: "linux 网络测速"
date: 2022-02-22
draft: false
author: "jobcher"
tags: ["javascript"]
categories: ["日常"]
series: ["日常系列"]
---

# linux 网络测速

## 一键测试脚本`bench.sh`

适用于各种 Linux 发行版的网络（下行）和 IO 测试：

1. 显示当前测试的各种系统信息
2. 取自世界多处的知名数据中心的测试点，下载测试比较全面
3. 支持 IPv6 下载测速
4. IO 测试三次，并显示平均值

```sh
wget -qO- bench.sh | bash
#或者下面这命令下载执行
curl -Lso- bench.sh | bash
```

欢迎关注我的博客[www.jobcher.com](https://www.jobcher.com/)
