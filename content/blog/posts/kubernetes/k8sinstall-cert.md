---
title: "Kubernetes — SSL 证书自动更新"
date: 2025-09-29
draft: false
featuredImage: "/images/cert-manager.png"
featuredImagePreview: "/images/cert-manager.png"
images: ["/images/cert-manager.png"]
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---
## 介绍
提供一个在 Kubernetes 中使用 cert-manager + Cloudflare 自动签发并自动更新 Let’s Encrypt 证书的完整思路与示例（DNS-01 验证），方便你在集群内自动化 TLS 证书更新。

## 前置条件
- Kubernetes 集群：可正常访问外网。不做网络环境配置的教程，具体可以去看其他文章
- Cloudflare 账号：已将你的域名托管到 Cloudflare。使用 Cloudflare 做 `dns-01 挑战`
- kubectl：已连接到集群。最基本的条件，保证k8s能正常访问
- helm：推荐用 Helm 安装 cert-manager。使用helm安装，方便干净

## 安装
官方推荐用 Helm，这里我使用 1.18.2 的版本，在我这个时间点这个版本还是比较新的
### 安装 cert-manager
```sh
# 安装 cert-manager CRDs
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.18.2/cert-manager.crds.yaml
```
```sh
## Add the Jetstack Helm repository
helm repo add jetstack https://charts.jetstack.io --force-update
```
```sh
## Install the cert-manager helm chart
helm install cert-manager --namespace cert-manager --version v1.18.2 jetstack/cert-manager
```
验证：
```sh
kubectl get pods -n cert-manager
```
### 准备 Cloudflare API Token
1. 登陆 [Cloudflare Dashboard](https://dash.cloudflare.com)
- My Profile → API Tokens → Create Token。

2. 模板选：`Edit zone DNS`

3. 权限：
- Zone → DNS → Edit
- Zone → Zone → Read

4. Zone Resources: 选择需要签发证书的域名。

5. 保存 token，例如：`CF_API_TOKEN=xxxxxxxxxx`

### 创建 Secret 存放 Token
```sh
kubectl create secret generic cloudflare-api-token-secret \
  --from-literal=api-token=CF_API_TOKEN \
  -n cert-manager
```

### 配置 ClusterIssuer
`cluster-issuer.yaml`
```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-dns
spec:
  acme:
    # 生产环境地址
    server: https://acme-v02.api.letsencrypt.org/directory
    email: your-email@example.com        # 接收过期提醒
    privateKeySecretRef:
      name: letsencrypt-dns-key
    solvers:
    - dns01:
        cloudflare:
          email: your-email@example.com  # 或留空（若使用 API token 可不填 email）
          apiTokenSecretRef:
            name: cloudflare-api-token-secret
            key: api-token
```
部署yaml
```sh
kubectl apply -f cluster-issuer.yaml
```
### 申请证书
在需要证书的命名空间下创建 Certificate 对象，例如 Nginx Ingress 的域名 `example.com`：  
`certificate.yaml`
```sh
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: example-com
  namespace: default
spec:
  secretName: example-com-tls        # 生成的 secret 名称
  issuerRef:
    name: letsencrypt-dns
    kind: ClusterIssuer
  commonName: example.com
  dnsNames:
  - example.com
  - "*.example.com"                   # 可选：通配符
```
部署yaml
```sh
kubectl apply -f certificate.yaml
```

### 在 Ingress 中引用
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: example-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-dns # 配置证书自动生成
spec:
  tls:
  - hosts:
    - example.com
    secretName: example-com-tls #修改成自己要使用的tls名称会自动生成
  rules:
  - host: example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: web
            port:
              number: 80
```

>cert-manager 会在证书到期前约 30 天 自动续期。  
  
查看状态  
```sh
kubectl describe certificate example-com -n default
```

### 复用证书
因为证书只能部署在制定的 namespace下，因此我写了一个脚本，把`default`命名空间下的 `example-com-tls` Secret复制到 `argocd`、`longhorn-system`、`cattle-system` 等多个命名空间。  
```bash
#!/bin/bash
# 文件名：copy-example-cert.sh
# 作用：将 default 命名空间下的 example-com-tls Secret
#       复制到 argocd、longhorn-system、cattle-system

SRC_NS="default"
SECRET_NAME="example-com-tls"
TARGET_NAMESPACES=("argocd" "longhorn-system" "cattle-system")

for ns in "${TARGET_NAMESPACES[@]}"; do
  echo ">>> 正在复制到命名空间: $ns"
  kubectl get secret "$SECRET_NAME" -n "$SRC_NS" -o yaml \
    | sed "s/namespace: $SRC_NS/namespace: $ns/" \
    | kubectl apply -f -
done

echo "✅ 复制完成"

```
执行
```sh
./copy-example-cert.sh
```
脚本会依次输出
```console
>>> 正在复制到命名空间: argocd
secret/example-com-tls created
>>> 正在复制到命名空间: longhorn-system
secret/example-com-tls created
>>> 正在复制到命名空间: cattle-system
secret/example-com-tls created
✅ 复制完成
```