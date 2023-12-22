---
title: "安装 minIO Azure S3网关"
date: 2022-01-07
draft: false
authors: "jobcher"
tags: ["minIO"]
categories: ["日常"]
series: ["日常系列"]
---

# 安装 minIO

## 通过 docker 安装

```sh
docker run -p 9000:9000 -p 41863:41863 -d --name azure-s3 \
 -e "MINIO_ACCESS_KEY=azure存储账户" \
 -e "MINIO_SECRET_KEY=azure存储密码" \
 minio/minio gateway azure --console-address ":41863"
```

## 通过 docker-compose 安装

```yaml
version: "3"
services:
  minio:
    image: "minio/minio:RELEASE.2022-01-04T07-41-07Z.fips"
    container_name: "minio"
    restart: "always"
    volumes:
      - "/etc/localtime:/etc/localtime"
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      - "MINIO_ROOT_USER=azure存储账户"
      - "MINIO_ROOT_PASSWORD=azure存储密码"
    command:
      - --console-address ":41863"
```
