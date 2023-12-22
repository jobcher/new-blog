---
title: "Kubernetes 创建nfs存储类"
date: 2021-12-13
draft: false
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---

# Kubernetes 创建 nfs 存储类

首先你需要在别的终端上创建 nfs 服务并能提供 nfs 访问  
Kubernetes 不包含内部 NFS 驱动。你需要使用外部驱动为 NFS 创建 StorageClass。  
https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner  
安装 nfs 驱动

## 安装 nfs 驱动

```sh
#安装nfs客户端
apt-get install nfs-common
git clone https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner.git
cd nfs-subdir-external-provisioner/deploy
k3s kubectl create -f rbac.yaml
vim deployment.yaml
```

1. 编辑 deployment.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nfs-client-provisioner
  labels:
    app: nfs-client-provisioner
  # replace with namespace where provisioner is deployed
  namespace: default
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app: nfs-client-provisioner
  template:
    metadata:
      labels:
        app: nfs-client-provisioner
    spec:
      serviceAccountName: nfs-client-provisioner
      containers:
        - name: nfs-client-provisioner
          image: k8s.gcr.io/sig-storage/nfs-subdir-external-provisioner:v4.0.2
          volumeMounts:
            - name: nfs-client-root
              mountPath: /persistentvolumes
          env:
            - name: PROVISIONER_NAME
              value: k8s-sigs.io/nfs-subdir-external-provisioner
            - name: NFS_SERVER
              value: 192.168.99.235
            - name: NFS_PATH
              value: /volume2/nfs-k8s
      volumes:
        - name: nfs-client-root
          nfs:
            server: 192.168.99.235
            path: /volume2/nfs-k8s
```

2. 定义存储类

```sh
k3s kubectl create -f deployment.yaml
k3s kubectl create -f class.yaml
```

## 测试存储是否正常

```sh
k3s kubectl create -f test-claim.yaml -f test-pod.yaml
k3s kubectl delete -f test-claim.yaml -f test-pod.yaml
```

## 创建有状态 pods（mysql）

创建 mysql-deployment.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql
spec:
  ports:
    - port: 3306
      targetPort: 3306
  selector:
    app: mysql
  clusterIP: None
  type: LoadBalancer
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mysql
spec:
  selector:
    matchLabels:
      app: mysql
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
        - image: mysql:5.7
          name: mysql
          env:
            # Use secret in real usage
            - name: MYSQL_ROOT_PASSWORD
              value: password
          ports:
            - containerPort: 3306
              name: mysql
          volumeMounts:
            - name: mysql-persistent-storage
              mountPath: /var/lib/mysql
      volumes:
        - name: mysql-persistent-storage
          persistentVolumeClaim:
            claimName: mysql-pv-claim
```

创建 mysql-pv.yaml

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mysql-pv-volume
spec:
  storageClassName: managed-nfs-storage
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteOnce
  nfs:
    server: 192.168.99.235
    path: "/volume2/nfs-k8s"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pv-claim
spec:
  storageClassName: managed-nfs-storage
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 20Gi
```

## 部署 mysql

```yaml
k3s kubectl apply -f mysql-pv.yaml
k3s kubectl apply -f mysql-deployment.yaml
k3s kubectl describe deployment mysql
```
