---
title: "清理Docker容器日志"
date: 2022-09-19
draft: false
authors: "jobcher"
tags: ["docker"]
categories: ["docker"]
series: ["docker入门系列"]
---

# 清理 Docker 容器日志

如果 docker 容器正在运行，那么使用`rm -rf`方式删除日志后，通过`df -h`会发现磁盘空间并没有释放。原因是在 Linux 或者 Unix 系统中，通过`rm -rf`或者文件管理器删除文件，将会从文件系统的目录结构上解除链接（unlink）。如果文件是被打开的（有一个进程正在使用），那么进程将仍然可以读取该文件，磁盘空间也一直被占用。正确姿势是`cat /dev/null > *-json.log`，当然你也可以通过`rm -rf`删除后重启 docker。

## 日志清理脚本 clean_docker_log.sh

```sh
#!/bin/sh

echo "======== start clean docker containers logs ========"

logs=$(find /var/lib/docker/containers/ -name *-json.log)

for log in $logs
        do
                echo "clean logs : $log"
                cat /dev/null > $log
        done

echo "======== end clean docker containers logs ========"

```

> chmod +x clean_docker_log.sh && ./clean_docker_log.sh

## 设置 Docker 容器日志大小

设置一个容器服务的日志大小上限  
上述方法，日志文件迟早又会涨回来。要从根本上解决问题，需要限制容器服务的日志大小上限。这个通过配置容器`docker-compose`的`max-size`选项来实现

```yaml
nginx:
  image: nginx:1.12.1
  restart: always
  logging:
    driver: “json-file”
    options:
      max-size: “5g”
```

### 全局设置

新建`/etc/docker/daemon.json`，若有就不用新建了。添加`log-dirver`和`log-opts`参数

```json
# vim /etc/docker/daemon.json

{

  "log-driver":"json-file",
  "log-opts": {"max-size":"500m", "max-file":"3"}
}

```

`max-size=500m`，意味着一个容器日志大小上限是`500M`  
`max-file=3`，意味着一个容器有三个日志，分别是`id+.json`、`id+1.json`、`id+2.json`。

> 注意：设置的日志大小，只对新建的容器有效。
