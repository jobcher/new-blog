---
title: "Nexus3 使用和部署"
date: 2022-10-31
draft: false
authors: "jobcher"
tags: ["work"]
categories: ["work"]
series: ["问题库系列"]
---

## Nexus3 docker-compose 安装

创建外部存储

```sh
mkdir -p /data/nexus
chmod +777 -R /data/nexus
```

运行 docker-compose

```sh
version: '3'
services:
  nexus3:
    image: sonatype/nexus3:3.42.0
    container_name: nexus3
    ports:
      - 8081:8081
      - 5000:5000
    volumes:
      - /data/nexus:/nexus-data
    environment:
      - INSTALL4J_ADD_VM_PARAMS=-Xms1024m -Xmx1024m -XX:MaxDirectMemorySize=1024m -Djava.util.prefs.userRoot=/some-other-dir
    restart: always
    # 赋予外部root权限
    privileged: true
```

> docker-compose up -d
> 运行 docker-compose

##
