---
title: "gitlab CI/CD 的使用"
date: 2021-12-03
draft: false
author: "jobcher"
tags: ["gitlab"]
categories: ["gitlab"]
series: ["git入门系列"]
---

# gitlab CI/CD 的使用

我将使用 gitlab 的流水线自动实现 hugo blog 文章的自动发布。

## 一、基础知识

## 二、安装过程

### 1.安装 gitlab runner

首先需要安装 gitlab runner 进入服务器 A  
安装方法：

1. 容器部署
2. 手动二进制文件部署
3. 通过 rpm/deb 包部署

4. docker 方式安装

安装文档：https://docs.gitlab.com/runne...

```sh
    docker run -dit \
    --name gitlab-runner \
    --restart always \
    -v /srv/gitlab-runner/config:/etc/gitlab-runner \
    -v /var/run/docker.sock:/var/run/docker.sock \
    gitlab/gitlab-runner
```

1.1 设置信息

    docker exec -it gitlab-runner gitlab-runner register

2. 非 docker 方式安装

2.1 安装 GitLab Runner

安装环境：Linux

其他环境参考：https://docs.gitlab.com/runne...

下载

```sh
    curl -L --output /usr/local/bin/gitlab-runner https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64
```

添加权限

```sh
    chmod +x /usr/local/bin/gitlab-runner
```

新建 gitlab-runner 用户

```sh
    sudo useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash
```

安装

安装时需要指定我们上面新建的用户

```sh
    gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner
```

启动

```sh
    gitlab-runner start
```

```sh
# Download the binary for your system
sudo curl -L --output /usr/local/bin/gitlab-runner https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64

# Give it permissions to execute
sudo chmod +x /usr/local/bin/gitlab-runner

# Create a GitLab CI user
sudo useradd --comment 'GitLab Runner' --create-home gitlab-runner --shell /bin/bash

# Install and run as service
sudo gitlab-runner install --user=gitlab-runner --working-directory=/home/gitlab-runner
sudo gitlab-runner start
```

### 2.配置 docker shell 链接

```sh
    ssh-keygen -t rsa
    cd .ssh/
    cat id_rsa.pub >>authorized_keys
    docker cp id_rsa gitlab-runner:/root
    docker exec -it gitlab-runner /bin/bash
    chmod 600 /root/id_rsa


    vim /etc/systemd/system/gitlab-runner.service

    "--syslog" "--user" "root" #修改为root
    wq保存退出

    systemctl daemon-reload
    systemctl restart gitlab-runner
```

### 3.配置.gitlab-ci.yml 文件

```sh
    vim .gitlab-ci.yml

    stages:
    - build
    - test
    - deploy

    build-job:
    stage: build
    script:
        - echo "上传代码"
        - echo "上传完成."

    unit-test-job:
    stage: test
    script:
        - echo
        - sleep 60
        - echo "Code coverage is 90%"

    lint-test-job:
    stage: test
    script:
        - echo "Linting code... This will take about 10 seconds."
        - sleep 10
        - echo "No lint issues found."

    deploy-job:
    stage: deploy
    script:
        - echo "Deploying application..."
        - echo "Application successfully deployed."
```
