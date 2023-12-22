---
title: "搭建docker registry 镜像仓库"
date: 2021-12-22
draft: false
author: "jobcher"
tags: ["docker"]
categories: ["docker"]
series: ["docker入门系列"]
---

# 搭建 docker registry 镜像仓库

## 获取镜像

```sh
docker pull registry:2.7.1
```

```sh
docker pull hyper/docker-registry-web
```

## 容器运行

```sh
mkdir -p /opt/data/registry
docker run -d -p 5000:5000 -v /opt/data/registry:/var/lib/registry --name registry registry:2.7.1
```

```sh
docker run -d -p 8080:8080 --name registry-web --link registry \
   -e REGISTRY_URL=http://192.168.99.146:5000/v2 \
   -e REGISTRY_TRUST_ANY_SSL=true \
   -e REGISTRY_BASIC_AUTH="GjhYGDGi2HhkJB" \
   -e REGISTRY_NAME=192.168.99.146:5000 \
   hyper/docker-registry-web
```

## 上传容器

```sh
vim /etc/docker/daemon.json
{
    "insecure-registries": ["192.168.99.146:5000"]
}

docker tag sjtfreaks/hogo-nginx:v1.1 192.168.99.146:5000/sjtfreaks/hogo-nginx:v1.1
docker push 192.168.99.146:5000/sjtfreaks/hogo-nginx:v1.1
```
