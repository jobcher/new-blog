---
title: "docker 安装kong 网关"
date: 2021-12-22
draft: false
authors: "jobcher"
tags: ["docker"]
categories: ["docker"]
series: ["docker入门系列"]
---

# docker 安装 kong 网关

## 建立数据库

1. 创建网络

```sh
docker network create kong-net
```

2. 建立数据库

```sh
docker run -d --name kong-database \
  --network=kong-net \
  -p 5432:5432 \
  -e "POSTGRES_USER=kong" \
  -e "POSTGRES_DB=kong" \
  -e "POSTGRES_PASSWORD=kong123" \
  postgres:9.6
```

3. 创建 kong 数据

```sh
docker run --rm --network=kong-net \
  -e "KONG_DATABASE=postgres" \
  -e "KONG_PG_HOST=kong-database" \
  -e "KONG_PG_PASSWORD=kong123" \
  -e "KONG_PASSWORD=kong123" \
  kong:latest kong migrations bootstrap
```

## 创建 kong

1. 创建 kong gateway

```sh
  docker run -d --name kong \
  --network=kong-net \
  -e "KONG_DATABASE=postgres" \
  -e "KONG_PG_HOST=kong-database" \
  -e "KONG_PG_USER=kong" \
  -e "KONG_PG_PASSWORD=kong123" \
  -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" \
  -e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
  -e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
  -e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
  -e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
  -e "KONG_ADMIN_LISTEN=0.0.0.0:8001, 0.0.0.0:8444 ssl" \
  -p 8000:8000 \
  -p 8443:8443 \
  -p 127.0.0.1:8001:8001 \
  -p 127.0.0.1:8444:8444 \
  kong:latest
```

## 安装 konga

```sh
docker pull pantsel/konga:latest
```

```sh
docker run --rm pantsel/konga:latest \
        -c prepare \
        -a postgres \
        -u postgresql://kong:kong123@172.18.0.1:5432/konga
```

```sh
docker run -d -p 1337:1337 \
        --network kong-net \
        --name konga \
        -e "NODE_ENV=production"  \
        -e "DB_ADAPTER=postgres" \
        -e "DB_URI=postgresql://kong:kong123@172.18.0.1:5432/konga" \
        pantsel/konga
```
