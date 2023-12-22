---
title: "yaml 语法"
date: 2021-12-30
draft: false
author: "jobcher"
tags: ["运维"]
categories: ["基础"]
series: ["基础知识系列"]
---

# yaml 语法

我们使用 YAML 是因为它像 XML 或 JSON 是一种利于人们读写的数据格式. 此外在大多数变成语言中有使用 YAML 的库.YAML 语法的基本概述, 它被用来描述一个 playbooks(我们的配置管理语言).

# 基本的 YAML

对于 Ansible, 每一个 YAML 文件都是从一个列表开始. 列表中的每一项都是一个键值对, 通常它们被称为一个 “哈希” 或 “字典”. 所以, 我们需要知道如何在 YAML 中编写列表和字典.  
YAML 还有一个小的怪癖. 所有的 YAML 文件(无论和 Ansible 有没有关系)开始行都应该是 ---. 这是 YAML 格式的一部分, 表明一个文件的开始.

```yaml
---
# 一个美味水果的列表
- Apple
- Orange
- Strawberry
- Mango
```
