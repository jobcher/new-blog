---
title: "docker进阶使用"
date: 2021-12-09
draft: false
author: "jobcher"
tags: ["docker", "dockerfile", "DockerCompose"]
categories: ["docker"]
series: ["docker进阶系列"]
---

# docker 进阶使用

dockerfile 和 docker compose 的配置

# Dockerfile 使用

Dockerfile 是一个用来构建镜像的文本文件，文本内容包含了一条条构建镜像所需的指令和说明。

例子：

```dockerfile
    FROM nginx
    RUN echo '这是一个本地构建的nginx镜像' > /usr/share/nginx/html/index.html
```

保存 Dockerfile 文件并在本地路径执行

```sh
    docker build -t nginx:v1-test .
    docker run -name docker run --name nginx-test -d  -p 8080:80 nginx:v1-test
```

浏览 nginx 页面确认更新内容

    curl 127.0.0.1:8080
    输出：
    这是一个本地构建的nginx镜像

## Docker 命令详解

### COPY

复制指令，从上下文目录中复制文件或者目录到容器里指定路径。

```dockerfile
    COPY [--chown=<user>:<group>] <源路径1>...  <目标路径>
    COPY [--chown=<user>:<group>] ["<源路径1>",...  "<目标路径>"]
```

<源路径>：源文件或者源目录，这里可以是通配符表达式，其通配符规则要满足 Go 的 filepath.Match 规则。例如：

```dockerfile
    COPY hom* /mydir/
    COPY hom?.txt /mydir/
```

### FROM

FROM：定制的镜像都是基于 FROM 的镜像

```dockerfile
    FROM nginx
```

### RUN

RUN：用于执行后面跟着的命令行命令

shell：

```dockerfile
    RUN <命令行命令>
    # <命令行命令> 等同于，在终端操作的 shell 命令。
```

exec：

```dockerfile
    RUN ["可执行文件", "参数1", "参数2"]
    # 例如：
    # RUN ["./test.php", "dev", "offline"] 等价于 RUN ./test.php dev offline
```

### ADD

ADD 指令和 COPY 的使用格类似

ADD 的优点：在执行 <源文件> 为 tar 压缩文件的话，压缩格式为 gzip, bzip2 以及 xz 的情况下，会自动复制并解压到 <目标路径>。

ADD 的缺点：在不解压的前提下，无法复制 tar 压缩文件。会令镜像构建缓存失效，从而可能会令镜像构建变得比较缓慢。具体是否使用，可以根据是否需要自动解压来决定。

### CMD

类似于 RUN 指令，用于运行程序，但二者运行的时间点不同:

1. CMD 在 docker run 时运行。
2. RUN 是在 docker build。

Dockerfile 中如果存在多个 CMD 指令，仅最后一个生效。

```dockerfile
    CMD <shell 命令>
    CMD ["<可执行文件或命令>","<param1>","<param2>",...]
    CMD ["<param1>","<param2>",...]  # 该写法是为 ENTRYPOINT 指令指定的程序提供默认参数
```

## 通过 dockerfile 文件封装 hugo

dokcerfile 文件

```dockerfile
    FROM nginx:1.21
    COPY public/ /usr/share/nginx/html
```

docker.sh 文件

```sh
    #/!bin/bash
    echo "删除旧的docker"
    docker ps
    docker stop nginx-hugo
    docker rm nginx-hugo
    docker rmi nginx:hugo
    echo "生成新的docker"
    hugo -t LoveIt -D
    docker build -t nginx:hugo .
    docker run --name nginx-hugo -d -p 8080:80 nginx:hugo
    echo "显示端口"
    netstat -lntp
```

执行脚本：

```sh
    sh update.sh
```
