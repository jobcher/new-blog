---
title: "Golang 初识（安装、使用）"
date: 2022-09-21
draft: false
authors: "jobcher"
tags: ["golang"]
categories: ["go"]
series: ["go系列"]
---

# Golang 初识（安装、使用）

## Go 导学

go 语言由 google 公司推出。

- 运行速度快，简单易学
- 适合区块链开发
- 拥有丰富指令
- 可以直接包含 C 语言
- 语言层面支持并发

### Go 方向

- 网络编程
- 服务器编程
- 区块链开发

## 环境安装

### 安装环境

安装包下载  
https://golang.google.cn/dl/

### windows 部署

```sh
wget https://golang.google.cn/dl/go1.19.1.windows-amd64.msi
# 直接安装
```

### GOPATH 设置

在环境变量 `PATH` 上直接配置安装地址

## 编写第一个程序

```go
package main

import "fmt"

func main() {
	fmt.Println("Hello World!")
}
```
