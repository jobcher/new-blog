---
title: "helm 安装"
date: 2021-12-16
draft: false
author: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---

# helm 安装

## 脚本安装

```sh
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh

#或者可以使用这个命令
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

helm help
```

## 二进制安装

```sh
wget https://get.helm.sh/helm-v3.7.2-linux-amd64.tar.gz
tar -zxvf helm-v3.7.2-linux-amd64.tar.gz
cd helm-v3.7.2-linux-amd64
mv linux-amd64/helm /usr/local/bin/helm
helm help
```
