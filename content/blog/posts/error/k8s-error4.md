---
title: "k8s master 节点 Unauthorized"
date: 2025-12-03
draft: false
authors: "jobcher"
featuredImage: '/images/error.jpeg'
featuredImagePreview: '/images/error.jpeg'
images: ['/images/error.jpeg']
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---
## 故障
多master集群的k8s其中一台master 出现 `You must be logged in to the server (Unauthorized)`，说明当前节点上的 kubectl 无法认证到 apiserver。  
共有三台master：
- master1 正常
- master2 异常
- master3 异常  
## 排查问题
### 用 admin 证书直接用 curl 测试认证
```sh
curl --cert /etc/kubernetes/pki/admin.crt --key /etc/kubernetes/pki/admin.key -k https://10.10.10.68:6443/healthz
```
- 结果 HTTP/1.1 200 OK → TLS + 认证成功（client cert 被接受）。
- 若返回 401/403 → TLS 成功但授权失败（RBAC 或 client cert 虽被接受但没有权限）。把完整返回贴来。
- 若 curl 报 TLS 错误 → client cert/CA 有问题（贴错误信息）。  

### 把 kubeconfig 中的 client cert/key 解码到临时文件并检查它们（确认 kubectl 用的是这个证书、并检查证书 Subject/CN/到期日）  
```sh
# 解码证书与私钥
awk '/client-certificate-data/ {print $2}' ~/.kube/config | base64 -d > /tmp/kube_client.crt
awk '/client-key-data/ {print $2}' ~/.kube/config | base64 -d > /tmp/kube_client.key
chmod 600 /tmp/kube_client.key

# 检查证书主体与到期日
openssl x509 -in /tmp/kube_client.crt -noout -text | egrep "Subject:|Not Before|Not After"
```
如果看到 Not After 在过去 → 证书过期（那就说明证书过期了）。
Subject 的 CN 通常是 kubernetes-admin 或 admin / system:admin，我们要确认 CN 是否合理。  
  
> 我这边发现证书过期了，所以导致认证失败了。  

## 解决
### 使用 kubeadm 重新生成 admin kubeconfig.  
在任意 master 上执行：  
```sh
# 备份旧配置
mv ~/.kube/config ~/.kube/config.bak

# 重新生成 admin kubeconfig
kubeadm init phase kubeconfig admin

# 拷贝到用户目录
mkdir -p ~/.kube
cp /etc/kubernetes/admin.conf ~/.kube/config
chmod 600 ~/.kube/config

```
测试
```sh
kubectl get nodes
```