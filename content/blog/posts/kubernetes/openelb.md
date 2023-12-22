---
title: "OpenELB：让k8s私有环境对外暴露端口"
date: 2022-04-13
draft: false
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---

# OpenELB：云原生负载均衡器插件

OpenELB 是一个开源的云原生负载均衡器实现，可以在基于裸金属服务器、边缘以及虚拟化的 Kubernetes 环境中使用 LoadBalancer 类型的 Service 对外暴露服务。

## 在 Kubernetes 中安装 OpenELB

```sh
kubectl apply -f https://raw.githubusercontent.com/openelb/openelb/master/deploy/openelb.yaml
```

- 查看状态

```sh
kubectl get po -n openelb-system
```

## 使用 kubectl 删除 OpenELB

```sh
kubectl delete -f https://raw.githubusercontent.com/openelb/openelb/master/deploy/openelb.yaml
```

```sh
kubectl get ns
```

## 配置 OpenELB

```sh
kubectl edit configmap kube-proxy -n kube-system

# 修改 网卡
ipvs:
  strictARP: true
```

### 重启组件

```sh
kubectl rollout restart daemonset kube-proxy -n kube-system
```

### 为 master1 节点添加一个 annotation 来指定网卡：

```sh
kubectl annotate nodes master1 layer2.openelb.kubesphere.io/v1alpha1="192.168.0.2"
```

### 创建地址池 `layer2-eip.yaml`

```yaml
apiVersion: network.kubesphere.io/v1alpha2
kind: Eip
metadata:
  name: layer2-eip
spec:
  address: 192.168.0.91-192.168.0.100
  interface: eth0
  protocol: layer2
```

### 创建部署 `jobcher-service.yaml`

```yaml
#暴露端口
apiVersion: v1
kind: Service
metadata:
  name: jobcher-service
  annotations:
    lb.kubesphere.io/v1alpha1: openelb
    protocol.openelb.kubesphere.io/v1alpha1: layer2
    eip.openelb.kubesphere.io/v1alpha2: layer2-eip
  labels:
    app: jobcher-blog
spec:
  selector:
    app: jobcher-blog
  ports:
    - name: jobcher-port
      protocol: TCP
      port: 80
      targetPort: 80
  type: LoadBalancer
```
