---
title: "Golang go build 编译不同版本"
date: 2022-04-17
draft: false
author: "jobcher"
tags: ["golang"]
categories: ["go"]
series: ["go系列"]
---

# Golang go build 编译不同系统下的可执行文件

## Mac 系统编译

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build test.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build test.go
```

## Linux 系统编译

```sh
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build test.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build test.go
```

## windows 系统编译

```bat
SET CGO_ENABLED=0 SET GOOS=darwin3 SET GOARCH=amd64 go build test.go
SET CGO_ENABLED=0 SET GOOS=linux SET GOARCH=amd64 go build  test.go
```

- `GOOS`：目标可执行程序运行操作系统，支持 `darwin，freebsd，linux，windows`
- `GOARCH`：目标可执行程序操作系统构架，包括 `386，amd64，arm`
