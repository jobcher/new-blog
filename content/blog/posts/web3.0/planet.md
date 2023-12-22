---
title: "Planet 下载及安装"
date: 2022-08-17
draft: false
authors: "jobcher"
tags: ["Planet", "web3.0"]
categories: ["web3.0"]
series: ["web3.0"]
---

# Planet 下载及安装

## [官网下载](https://www.planetable.xyz/)

`Planet` 是一款用于发布和关注 Web 内容的免费开源软件，它不需要集中式服务器或服务。它使用 `IPFS` 来实现`点对点`的`内容分发`。此外，您可以将您的内容链接到`以太坊名称 (.eth)`，以便其他人可以通过 `Planet` 以 `.eth` 名称关注您。由于 `IPFS` 和 `ENS` 都是去中心化的，因此您可以以`去中心化`的方式构建您的网站或关注其他网站。

## 如何使用

标准是 `EIP-1577`，这个 `Content Hash` 字段可以接受一些可能的值。例如，`IPFS`——另一种去中心化的内容分发技术。而[vitalik.eth](https://ipfs.io/ipns/vitalik.eth) 网站已经在 IPFS 上运行。

![logo](/images/vitalik-eth.png)

> 通过 Planet 关注来自 vitalik.eth 的更新

使用 `Planet` 创建网站后，右键单击侧栏中的项目，然后选择`Copy IPNS`，然后您将在粘贴板中看到如下所示的内容：

```sh
k51qzi5uqu5dgv8kzl1anc0m74n6t9ffdjnypdh846ct5wgpljc7rulynxa74a
```

## 公开 ENS

然后您可以像这样将该 `IPNS` 放入您的 `ENS   ContentHash` 中：  
![ETH](/images/set-content-hash.png)

> 确保在该字符串之前添加了 `ipns://`。

## 完成！

然后您的网站将链接到您的 `ENS`。恭喜！现在你有一个在 `ENS + IPFS` 上运行的去中心化网站！
