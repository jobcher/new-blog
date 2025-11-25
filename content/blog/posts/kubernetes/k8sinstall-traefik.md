---
title: "Kubernetes — metalLB + Traefik 部署"
date: 2025-11-25
draft: false
featuredImage: "/images/traefik-logo.png"
featuredImagePreview: "/images/traefik-logo.png"
images: ["/images/traefik-logo.png"]
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---
## 背景
鉴于 Ingress NGINX 将在 2026 年 3 月停止积极维护（只保留 “best-effort maintenance”）考虑切换到Traefik。Traefik 官方推荐是最直接的替代，因为 Traefik 围绕 Ingress NGINX 的兼容层做了优化：它对部分常见的 nginx-ingress 注解提供了兼容支持。

## MEtalLB 安装
```sh
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.15.2/config/manifests/metallb-native.yaml
```
```sh
kubectl get pods -n metallb-system
```
创建 `metallb-config.yaml`
```sh
# metallb-config.yaml
apiVersion: metallb.io/v1beta1
kind: IPAddressPool
metadata:
  name: local-pool
  namespace: metallb-system
spec:
  addresses:
    - 10.10.10.180-10.10.10.181  # ← 修改为你的局域网可用 IP
---
apiVersion: metallb.io/v1beta1
kind: L2Advertisement
metadata:
  name: l2adv
  namespace: metallb-system
```
```sh
kubectl apply -f metallb-config.yaml
```

# 安装
```sh
helm repo add traefik https://traefik.github.io/charts
helm repo update
```
helm 安装
```sh
helm install traefik traefik/traefik \
  -n traefik --create-namespace \
  --set service.type=LoadBalancer \
  --set ingressClass.enabled=true \
  --set ingressClass.isDefaultClass=true \
  --set dashboard.enabled=true \
  --set api.dashboard=true \
  --set api.insecure=false \
  --set ports.web.expose.enabled=true \
  --set ports.websecure.expose.enabled=true \
  --set ports.websecure.tls.enabled=true \
  --set metrics.prometheus.enabled=true

```
验证
```sh
kubectl get pods -n traefik
kubectl get svc -n traefik
```

### 启用dashboard
创建`traefik-dashboard.yaml`
```yaml
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: traefik-dashboard
  namespace: traefik
  annotations:
    kubernetes.io/ingress.class: traefik
spec:
  entryPoints:
    - websecure
  routes:
    - match: Host(`traefik.jobcher.com`) # 改为你自己的域名地址
      kind: Rule
      services:
        - name: api@internal
          kind: TraefikService
  tls:
    secretName: jobcher-com-tls # 改为你自己的tls证书
```
部署
```sh
kubectl -n traefik apply -f traefik-dashboard.yaml
```
验证
```sh
kubectl -n traefik get ingressRoute
```

## 访问地址
https://traefik.jobcher.com  
![traefik-dashboard](/images/traefik-dashboard.jpg)  