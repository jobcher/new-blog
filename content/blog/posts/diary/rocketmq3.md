---
title: "RocketMQ k8s部署 4主4从集群"
date: 2022-08-29
draft: false
authors: "jobcher"
tags: ["RocketMQ"]
categories: ["日常"]
series: ["日常系列"]
---

# RocketMQ k8s 部署 4 主 4 从集群

## 使用 NFS 配置 StatefulSet 的动态持久化存储

### 安装 NFS 服务端

```sh
sudo apt update
sudo apt install nfs-kernel-server nfs-common
```

### 安装 NFS 客户端

所有的节点都得执行

> sudo apt install nfs-common -y

### 创建目录

```sh
mkdir -p /data/storage/k8s/rocketmq
```

使用 NFS 作为`StatefulSet`持久化存储的操作记录，分别需要创建`nfs-provisioner`的`rbac`、`storageclass`、`nfs-client-provisioner`和`statefulset`的`pod`

### 创建 nfs 的 rbac

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: nfs-provisioner
  namespace: sanjiang
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: nfs-provisioner-runner
  namespace: sanjiang
rules:
  - apiGroups: [""]
    resources: ["persistentvolumes"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["persistentvolumeclaims"]
    verbs: ["get", "list", "watch", "update"]
  - apiGroups: ["storage.k8s.io"]
    resources: ["storageclasses"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["events"]
    verbs: ["list", "watch", "create", "update", "patch"]
  - apiGroups: [""]
    resources: ["services", "endpoints"]
    verbs: ["get", "create", "list", "watch", "update"]
  - apiGroups: ["extensions"]
    resources: ["podsecuritypolicies"]
    resourceNames: ["nfs-provisioner"]
    verbs: ["use"]
  - apiGroups: [""]
    resources: ["endpoints"]
    verbs: ["create", "delete", "get", "list", "watch", "patch", "update"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: run-nfs-provisioner
subjects:
  - kind: ServiceAccount
    name: nfs-provisioner
    namespace: sanjiang
roleRef:
  kind: ClusterRole
  name: nfs-provisioner-runner
  apiGroup: rbac.authorization.k8s.io
```

> kubectl apply -f nfs-rbac.yaml

### 创建 rocketmq 集群的 storageclass

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: rocketmq-nfs-storage
  namespace: sanjiang
provisioner: rocketmq/nfs
reclaimPolicy: Retain
```

> kubectl apply -f rocketmq-nfs-class.yaml  
> 查看创建情况  
> kubectl get sc -n sanjiang

### 创建 rocketmq 集群的 nfs-client-provisioner

`PROVISIONER_NAME`的值一定要和`StorageClass`中的`provisioner`相等

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rocketmq-nfs-client-provisioner
  namespace: sanjiang
spec:
  replicas: 1
  selector:
    matchLabels:
      app: rocketmq-nfs-client-provisioner
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: rocketmq-nfs-client-provisioner
    spec:
      serviceAccount: nfs-provisioner
      containers:
        - name: rocketmq-nfs-client-provisioner
          image: registry.cn-hangzhou.aliyuncs.com/open-ali/nfs-client-provisioner
          imagePullPolicy: IfNotPresent
          volumeMounts:
            - name: nfs-client-root
              mountPath: /persistentvolumes
          env:
            - name: PROVISIONER_NAME
              value: rocketmq/nfs
            - name: NFS_SERVER
              value: 193.0.40.171 #nfs ip
            - name: NFS_PATH
              value: /data/storage/k8s/rocketmq
      volumes:
        - name: nfs-client-root
          nfs:
            server: 193.0.40.171 #nfs ip
            path: /data/storage/k8s/rocketmq
```

> kubectl apply -f rocketmq-nfs.yml

## k8s 部署

### 生成文件

broker-a-s.properties

```properties
brokerClusterName=rocketmq-cluster
brokerName=broker-a
brokerId=1
namesrvAddr=rocketmq-0.rocketmq:9876
defaultTopicQueueNums=4
autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true
listenPort=20911
deleteWhen=04
fileReservedTime=120
mapedFileSizeCommitLog=1073741824
mapedFileSizeConsumeQueue=300000
diskMaxUsedSpaceRatio=88
storePathRootDir=/data/rocketmq/store
maxMessageSize=65536
brokerRole=SLAVE
flushDiskType=SYNC_FLUSH
```

broker-a.properties

```properties
brokerClusterName=rocketmq-cluster
brokerName=broker-a
brokerId=0
namesrvAddr=rocketmq-0.rocketmq:9876
defaultTopicQueueNums=4
autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true
listenPort=20911
deleteWhen=04
fileReservedTime=120
mapedFileSizeCommitLog=1073741824
mapedFileSizeConsumeQueue=300000
diskMaxUsedSpaceRatio=88
storePathRootDir=/data/rocketmq/store
maxMessageSize=65536
brokerRole=MASTER
```

broker-b-s.properties

```properties
brokerClusterName=rocketmq-cluster
brokerName=broker-b
brokerId=1
namesrvAddr=rocketmq-0.rocketmq:9876
defaultTopicQueueNums=4
autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true
listenPort=20911
deleteWhen=04
fileReservedTime=120
mapedFileSizeCommitLog=1073741824
mapedFileSizeConsumeQueue=300000
diskMaxUsedSpaceRatio=88
storePathRootDir=/data/rocketmq/store
maxMessageSize=65536
brokerRole=SLAVE
flushDiskType=SYNC_FLUSH
```

broker-b.properties

```properties
brokerClusterName=rocketmq-cluster
brokerName=broker-b
brokerId=0
namesrvAddr=rocketmq-0.rocketmq:9876
defaultTopicQueueNums=4
autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true
listenPort=20911
deleteWhen=04
fileReservedTime=120
mapedFileSizeCommitLog=1073741824
mapedFileSizeConsumeQueue=300000
diskMaxUsedSpaceRatio=88
storePathRootDir=/data/rocketmq/store
maxMessageSize=65536
brokerRole=MASTER
flushDiskType=SYNC_FLUSH
```

broker-c-s.properties

```properties
brokerClusterName=rocketmq-cluster
brokerName=broker-c
brokerId=1
namesrvAddr=rocketmq-0.rocketmq:9876
defaultTopicQueueNums=4
autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true
listenPort=20911
deleteWhen=04
fileReservedTime=120
mapedFileSizeCommitLog=1073741824
mapedFileSizeConsumeQueue=300000
diskMaxUsedSpaceRatio=88
storePathRootDir=/data/rocketmq/store
maxMessageSize=65536
brokerRole=SLAVE
flushDiskType=SYNC_FLUSH

```

broker-c.properties

```properties
brokerClusterName=rocketmq-cluster
brokerName=broker-c
brokerId=0
namesrvAddr=rocketmq-0.rocketmq:9876
defaultTopicQueueNums=4
autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true
listenPort=20911
deleteWhen=04
fileReservedTime=120
mapedFileSizeCommitLog=1073741824
mapedFileSizeConsumeQueue=300000
diskMaxUsedSpaceRatio=88
storePathRootDir=/data/rocketmq/store
maxMessageSize=65536
brokerRole=MASTER
flushDiskType=SYNC_FLUSH
```

broker-d-s.properties

```properties
brokerClusterName=rocketmq-cluster
brokerName=broker-d
brokerId=1
namesrvAddr=rocketmq-0.rocketmq:9876
defaultTopicQueueNums=4
autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true
listenPort=20911
deleteWhen=04
fileReservedTime=120
mapedFileSizeCommitLog=1073741824
mapedFileSizeConsumeQueue=300000
diskMaxUsedSpaceRatio=88
storePathRootDir=/data/rocketmq/store
maxMessageSize=65536
brokerRole=SLAVE
flushDiskType=SYNC_FLUSH
```

broker-d.properties

```properties
brokerClusterName=rocketmq-cluster
brokerName=broker-d
brokerId=0
namesrvAddr=rocketmq-0.rocketmq:9876
defaultTopicQueueNums=4
autoCreateTopicEnable=true
autoCreateSubscriptionGroup=true
listenPort=20911
deleteWhen=04
fileReservedTime=120
mapedFileSizeCommitLog=1073741824
mapedFileSizeConsumeQueue=300000
diskMaxUsedSpaceRatio=88
storePathRootDir=/data/rocketmq/store
maxMessageSize=65536
brokerRole=MASTER
flushDiskType=SYNC_FLUSH
```

运行命令

```sh
kubectl create cm rocketmq-config --from-file=broker-a.properties --from-file=broker-a-s.properties --from-file=broker-b.properties --from-file=broker-b-s.properties --from-file=broker-c.properties --from-file=broker-c-s.properties --from-file=broker-d.properties --from-file=broker-d-s.properties -n sanjiang
```

> kubectl get cm -n sanjiang|grep rocketmq

### 创建配置文件

broker-a-s.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: broker-a-s
  name: broker-a-s
  namespace: sanjiang
spec:
  ports:
    - port: 20911
      targetPort: 20911
      name: broker-port
  selector:
    app: broker-a-s
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: broker-a-s
  namespace: sanjiang
spec:
  serviceName: broker-a-s
  replicas: 1
  selector:
    matchLabels:
      app: broker-a-s
  template:
    metadata:
      labels:
        app: broker-a-s
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                      - broker-a-s
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: broker-a-s
          image: liuyi71sinacom/rocketmq-4.8.0
          imagePullPolicy: IfNotPresent
          command:
            [
              "sh",
              "-c",
              "mqbroker  -c /usr/local/rocketmq-4.8.0/conf/broker-a-s.properties",
            ]
          env:
            - name: JAVA_OPT
              value: "-server -XX:ParallelGCThreads=1 -Xms1g -Xmx1g -Xmn512m"
              #value: "-XX:MaxRAMPercentage=80.0"
          volumeMounts:
            - mountPath: /root/logs
              name: rocketmq-data
              subPath: mq-brokeroptlogs
            - mountPath: /data/rocketmq
              name: rocketmq-data
              subPath: mq-brokeroptstore
            - name: broker-config
              mountPath: /usr/local/rocketmq-4.8.0/conf/broker-a-s.properties
              subPath: broker-a-s.properties
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "touch /tmp/health"]
          livenessProbe:
            exec:
              command: ["test", "-e", "/tmp/health"]
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 20911
            initialDelaySeconds: 15
            timeoutSeconds: 5
            periodSeconds: 20
      volumes:
        - name: broker-config
          configMap:
            name: rocketmq-config
            items:
              - key: broker-a-s.properties
                path: broker-a-s.properties
  volumeClaimTemplates:
    - metadata:
        name: rocketmq-data
        namespace: sanjiang
        annotations:
          volume.beta.kubernetes.io/storage-class: "rocketmq-nfs-storage"
      spec:
        accessModes:
          - ReadWriteMany
        resources:
          requests:
            storage: 2Gi
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: broker-a-s-pv
  namespace: sanjiang
spec:
  accessModes:
    - ReadWriteMany
  capacity:
    storage: 2Gi
  persistentVolumeReclaimPolicy: Retain
  storageClassName: rocketmq-nfs-storage
  nfs:
    path: /data/storage/k8s/rocketmq/broker-a-s
    server: 193.0.40.171
```

broker-a.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: broker-a
  name: broker-a
  namespace: sanjiang
spec:
  ports:
    - port: 20911
      targetPort: 20911
      name: broker-port
  selector:
    app: broker-a
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: broker-a
  namespace: sanjiang
spec:
  serviceName: broker-a
  replicas: 1
  selector:
    matchLabels:
      app: broker-a
  template:
    metadata:
      labels:
        app: broker-a
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                      - broker-a
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: broker-a
          image: liuyi71sinacom/rocketmq-4.8.0
          imagePullPolicy: IfNotPresent
          command:
            [
              "sh",
              "-c",
              "mqbroker  -c /usr/local/rocketmq-4.8.0/conf/broker-a.properties",
            ]
          env:
            - name: JAVA_OPT
              value: "-server -XX:ParallelGCThreads=1 -Xms1g -Xmx1g -Xmn512m"
              #value: "-XX:MaxRAMPercentage=80.0"
          volumeMounts:
            - mountPath: /root/logs
              name: rocketmq-data
              subPath: mq-brokeroptlogs
            - mountPath: /data/rocketmq
              name: rocketmq-data
              subPath: mq-brokeroptstore
            - name: broker-config
              mountPath: /usr/local/rocketmq-4.8.0/conf/broker-a.properties
              subPath: broker-a.properties
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "touch /tmp/health"]
          livenessProbe:
            exec:
              command: ["test", "-e", "/tmp/health"]
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 20911
            initialDelaySeconds: 15
            timeoutSeconds: 5
            periodSeconds: 20
      volumes:
        - name: broker-config
          configMap:
            name: rocketmq-config
  volumeClaimTemplates:
    - metadata:
        name: rocketmq-data
        namespace: sanjiang
        annotations:
          volume.beta.kubernetes.io/storage-class: "rocketmq-nfs-storage"
      spec:
        accessModes:
          - ReadWriteMany
        resources:
          requests:
            storage: 2Gi
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: broker-a-pv
  namespace: sanjiang
spec:
  accessModes:
    - ReadWriteMany
  capacity:
    storage: 2Gi
  persistentVolumeReclaimPolicy: Retain
  storageClassName: rocketmq-nfs-storage
  nfs:
    path: /data/storage/k8s/rocketmq/broker-a
    server: 193.0.40.171
```

broker-b-s.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: broker-b-s
  name: broker-b-s
  namespace: sanjiang
spec:
  ports:
    - port: 20911
      targetPort: 20911
      name: broker-port
  selector:
    app: broker-b-s
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: broker-b-s
  namespace: sanjiang
spec:
  serviceName: broker-b-s
  replicas: 1
  selector:
    matchLabels:
      app: broker-b-s
  template:
    metadata:
      labels:
        app: broker-b-s
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                      - broker-b-s
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: broker-b-s
          image: liuyi71sinacom/rocketmq-4.8.0
          imagePullPolicy: IfNotPresent
          command:
            [
              "sh",
              "-c",
              "mqbroker -c /usr/local/rocketmq-4.8.0/conf/broker-b-s.properties",
            ]
          env:
            - name: JAVA_OPT
              value: "-server -XX:ParallelGCThreads=1 -Xms1g -Xmx1g -Xmn512m"
              #value: "-XX:MaxRAMPercentage=80.0"
          volumeMounts:
            - mountPath: /root/logs
              name: rocketmq-data
              subPath: mq-brokeroptlogs
            - mountPath: /data/rocketmq
              name: rocketmq-data
              subPath: mq-brokeroptstore
            - name: broker-config
              mountPath: /usr/local/rocketmq-4.8.0/conf/broker-b-s.properties
              subPath: broker-b-s.properties
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "touch /tmp/health"]
          livenessProbe:
            exec:
              command: ["test", "-e", "/tmp/health"]
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 20911
            initialDelaySeconds: 15
            timeoutSeconds: 5
            periodSeconds: 20
      volumes:
        - name: broker-config
          configMap:
            name: rocketmq-config
            items:
              - key: broker-b-s.properties
                path: broker-b-s.properties
  volumeClaimTemplates:
    - metadata:
        name: rocketmq-data
        namespace: sanjiang
        annotations:
          volume.beta.kubernetes.io/storage-class: "rocketmq-nfs-storage"
      spec:
        accessModes:
          - ReadWriteMany
        resources:
          requests:
            storage: 2Gi
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: broker-b-s-pv
  namespace: sanjiang
spec:
  accessModes:
    - ReadWriteMany
  capacity:
    storage: 2Gi
  persistentVolumeReclaimPolicy: Retain
  storageClassName: rocketmq-nfs-storage
  nfs:
    path: /data/storage/k8s/rocketmq/broker-b-s
    server: 193.0.40.171
```

broker-b.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: broker-b
  name: broker-b
  namespace: sanjiang
spec:
  ports:
    - port: 20911
      targetPort: 20911
      name: broker-port
  selector:
    app: broker-b
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: broker-b
  namespace: sanjiang
spec:
  serviceName: broker-b
  replicas: 1
  selector:
    matchLabels:
      app: broker-b
  template:
    metadata:
      labels:
        app: broker-b
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                      - broker-b
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: broker-b
          image: liuyi71sinacom/rocketmq-4.8.0
          imagePullPolicy: IfNotPresent
          command:
            [
              "sh",
              "-c",
              "mqbroker  -c /usr/local/rocketmq-4.8.0/conf/broker-b.properties",
            ]
          env:
            - name: JAVA_OPT
              value: "-server -XX:ParallelGCThreads=1 -Xms1g -Xmx1g -Xmn512m"
              #value: "-XX:MaxRAMPercentage=80.0"
          volumeMounts:
            - mountPath: /root/logs
              name: rocketmq-data
              subPath: mq-brokeroptlogs
            - mountPath: /data/rocketmq
              name: rocketmq-data
              subPath: mq-brokeroptstore
            - name: broker-config
              mountPath: /usr/local/rocketmq-4.8.0/conf/broker-b.properties
              subPath: broker-b.properties
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "touch /tmp/health"]
          livenessProbe:
            exec:
              command: ["test", "-e", "/tmp/health"]
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 20911
            initialDelaySeconds: 15
            timeoutSeconds: 5
            periodSeconds: 20
      volumes:
        - name: broker-config
          configMap:
            name: rocketmq-config
  volumeClaimTemplates:
    - metadata:
        name: rocketmq-data
        namespace: sanjiang
        annotations:
          volume.beta.kubernetes.io/storage-class: "rocketmq-nfs-storage"
      spec:
        accessModes:
          - ReadWriteMany
        resources:
          requests:
            storage: 2Gi
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: broker-b-pv
  namespace: sanjiang
spec:
  accessModes:
    - ReadWriteMany
  capacity:
    storage: 2Gi
  persistentVolumeReclaimPolicy: Retain
  storageClassName: rocketmq-nfs-storage
  nfs:
    path: /data/storage/k8s/rocketmq/broker-b
    server: 193.0.40.171
```

broker-c-s.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: broker-c-s
  name: broker-c-s
  namespace: sanjiang
spec:
  ports:
    - port: 20911
      targetPort: 20911
      name: broker-port
  selector:
    app: broker-c-s
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: broker-c-s
  namespace: sanjiang
spec:
  serviceName: broker-c-s
  replicas: 1
  selector:
    matchLabels:
      app: broker-c-s
  template:
    metadata:
      labels:
        app: broker-c-s
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                      - broker-c-s
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: broker-c-s
          image: liuyi71sinacom/rocketmq-4.8.0
          imagePullPolicy: IfNotPresent
          command:
            [
              "sh",
              "-c",
              "mqbroker  -c /usr/local/rocketmq-4.8.0/conf/broker-c-s.properties",
            ]
          env:
            - name: JAVA_OPT
              value: "-server -XX:ParallelGCThreads=1 -Xms1g -Xmx1g -Xmn512m"
              #value: "-XX:MaxRAMPercentage=80.0"
          volumeMounts:
            - mountPath: /root/logs
              name: rocketmq-data
              subPath: mq-brokeroptlogs
            - mountPath: /data/rocketmq
              name: rocketmq-data
              subPath: mq-brokeroptstore
            - name: broker-config
              mountPath: /usr/local/rocketmq-4.8.0/conf/broker-c-s.properties
              subPath: broker-c-s.properties
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "touch /tmp/health"]
          livenessProbe:
            exec:
              command: ["test", "-e", "/tmp/health"]
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 20911
            initialDelaySeconds: 15
            timeoutSeconds: 5
            periodSeconds: 20
      volumes:
        - name: broker-config
          configMap:
            name: rocketmq-config
            items:
              - key: broker-c-s.properties
                path: broker-c-s.properties
  volumeClaimTemplates:
    - metadata:
        name: rocketmq-data
        namespace: sanjiang
        annotations:
          volume.beta.kubernetes.io/storage-class: "rocketmq-nfs-storage"
      spec:
        accessModes:
          - ReadWriteMany
        resources:
          requests:
            storage: 2Gi
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: broker-c-s-pv
  namespace: sanjiang
spec:
  accessModes:
    - ReadWriteMany
  capacity:
    storage: 2Gi
  persistentVolumeReclaimPolicy: Retain
  storageClassName: rocketmq-nfs-storage
  nfs:
    path: /data/storage/k8s/rocketmq/broker-c-s
    server: 193.0.40.171
```

broker-c.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: broker-c
  name: broker-c
  namespace: sanjiang
spec:
  ports:
    - port: 20911
      targetPort: 20911
      name: broker-port
  selector:
    app: broker-c
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: broker-c
  namespace: sanjiang
spec:
  serviceName: broker-c
  replicas: 1
  selector:
    matchLabels:
      app: broker-c
  template:
    metadata:
      labels:
        app: broker-c
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                      - broker-c
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: broker-c
          image: liuyi71sinacom/rocketmq-4.8.0
          imagePullPolicy: IfNotPresent
          command:
            [
              "sh",
              "-c",
              "mqbroker  -c /usr/local/rocketmq-4.8.0/conf/broker-c.properties",
            ]
          env:
            - name: JAVA_OPT
              value: "-server -XX:ParallelGCThreads=1 -Xms1g -Xmx1g -Xmn512m"
              #value: "-XX:MaxRAMPercentage=80.0"
          volumeMounts:
            - mountPath: /root/logs
              name: rocketmq-data
              subPath: mq-brokeroptlogs
            - mountPath: /data/rocketmq
              name: rocketmq-data
              subPath: mq-brokeroptstore
            - name: broker-config
              mountPath: /usr/local/rocketmq-4.8.0/conf/broker-c.properties
              subPath: broker-c.properties
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "touch /tmp/health"]
          livenessProbe:
            exec:
              command: ["test", "-e", "/tmp/health"]
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 20911
            initialDelaySeconds: 15
            timeoutSeconds: 5
            periodSeconds: 20
      volumes:
        - name: broker-config
          configMap:
            name: rocketmq-config
  volumeClaimTemplates:
    - metadata:
        name: rocketmq-data
        namespace: sanjiang
        annotations:
          volume.beta.kubernetes.io/storage-class: "rocketmq-nfs-storage"
      spec:
        accessModes:
          - ReadWriteMany
        resources:
          requests:
            storage: 2Gi
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: broker-c-pv
  namespace: sanjiang
spec:
  accessModes:
    - ReadWriteMany
  capacity:
    storage: 2Gi
  persistentVolumeReclaimPolicy: Retain
  storageClassName: rocketmq-nfs-storage
  nfs:
    path: /data/storage/k8s/rocketmq/broker-c
    server: 193.0.40.171
```

broker-d-s.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: broker-d-s
  name: broker-d-s
  namespace: sanjiang
spec:
  ports:
    - port: 20911
      targetPort: 20911
      name: broker-port
  selector:
    app: broker-d-s
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: broker-d-s
  namespace: sanjiang
spec:
  serviceName: broker-d-s
  replicas: 1
  selector:
    matchLabels:
      app: broker-d-s
  template:
    metadata:
      labels:
        app: broker-d-s
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                      - broker-d-s
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: broker-d-s
          image: liuyi71sinacom/rocketmq-4.8.0
          imagePullPolicy: IfNotPresent
          command:
            [
              "sh",
              "-c",
              "mqbroker  -c /usr/local/rocketmq-4.8.0/conf/broker-d-s.properties",
            ]
          env:
            - name: JAVA_OPT
              value: "-server -XX:ParallelGCThreads=1 -Xms1g -Xmx1g -Xmn512m"
              #value: "-XX:MaxRAMPercentage=80.0"
          volumeMounts:
            - mountPath: /root/logs
              name: rocketmq-data
              subPath: mq-brokeroptlogs
            - mountPath: /data/rocketmq
              name: rocketmq-data
              subPath: mq-brokeroptstore
            - name: broker-config
              mountPath: /usr/local/rocketmq-4.8.0/conf/broker-d-s.properties
              subPath: broker-d-s.properties
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "touch /tmp/health"]
          livenessProbe:
            exec:
              command: ["test", "-e", "/tmp/health"]
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 20911
            initialDelaySeconds: 15
            timeoutSeconds: 5
            periodSeconds: 20
      volumes:
        - name: broker-config
          configMap:
            name: rocketmq-config
            items:
              - key: broker-d-s.properties
                path: broker-d-s.properties
  volumeClaimTemplates:
    - metadata:
        name: rocketmq-data
        namespace: sanjiang
        annotations:
          volume.beta.kubernetes.io/storage-class: "rocketmq-nfs-storage"
      spec:
        accessModes:
          - ReadWriteMany
        resources:
          requests:
            storage: 2Gi
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: broker-d-s-pv
  namespace: sanjiang
spec:
  accessModes:
    - ReadWriteMany
  capacity:
    storage: 2Gi
  persistentVolumeReclaimPolicy: Retain
  storageClassName: rocketmq-nfs-storage
  nfs:
    path: /data/storage/k8s/rocketmq/broker-d-s
    server: 193.0.40.171
```

broker-d.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: broker-d
  name: broker-d
  namespace: sanjiang
spec:
  ports:
    - port: 20911
      targetPort: 20911
      name: broker-port
  selector:
    app: broker-d
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: broker-d
  namespace: sanjiang
spec:
  serviceName: broker-d
  replicas: 1
  selector:
    matchLabels:
      app: broker-d
  template:
    metadata:
      labels:
        app: broker-d
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                      - broker-d
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: broker-d
          image: liuyi71sinacom/rocketmq-4.8.0
          imagePullPolicy: IfNotPresent
          command:
            [
              "sh",
              "-c",
              "mqbroker  -c /usr/local/rocketmq-4.8.0/conf/broker-d.properties",
            ]
          env:
            - name: JAVA_OPT
              value: "-server -XX:ParallelGCThreads=1 -Xms1g -Xmx1g -Xmn512m"
              #value: "-XX:MaxRAMPercentage=80.0"
          volumeMounts:
            - mountPath: /root/logs
              name: rocketmq-data
              subPath: mq-brokeroptlogs
            - mountPath: /data/rocketmq
              name: rocketmq-data
              subPath: mq-brokeroptstore
            - name: broker-config
              mountPath: /usr/local/rocketmq-4.8.0/conf/broker-d.properties
              subPath: broker-d.properties
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "touch /tmp/health"]
          livenessProbe:
            exec:
              command: ["test", "-e", "/tmp/health"]
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 20911
            initialDelaySeconds: 15
            timeoutSeconds: 5
            periodSeconds: 20
      volumes:
        - name: broker-config
          configMap:
            name: rocketmq-config
  volumeClaimTemplates:
    - metadata:
        name: rocketmq-data
        namespace: sanjiang
        annotations:
          volume.beta.kubernetes.io/storage-class: "rocketmq-nfs-storage"
      spec:
        accessModes:
          - ReadWriteMany
        resources:
          requests:
            storage: 2Gi
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: broker-d-pv
  namespace: sanjiang
spec:
  accessModes:
    - ReadWriteMany
  capacity:
    storage: 2Gi
  persistentVolumeReclaimPolicy: Retain
  storageClassName: rocketmq-nfs-storage
  nfs:
    path: /data/storage/k8s/rocketmq/broker-d
    server: 193.0.40.171
```

console.yaml

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: console
  name: console
  namespace: sanjiang
spec:
  replicas: 1
  selector:
    matchLabels:
      app: console
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: console
    spec:
      containers:
        - image: styletang/rocketmq-console-ng
          name: rocketmq-console-ng
          env:
            - name: JAVA_OPTS
              value: "-Drocketmq.namesrv.addr=rocketmq:9876 -Dcom.rocketmq.sendMessageWithVIPChannel=false"
          resources: {}
status: {}
---
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels:
    app: console
  name: console
  namespace: sanjiang
spec:
  ports:
    - port: 8080
      protocol: TCP
      targetPort: 8080
      nodePort: 32080
  selector:
    app: console
  type: NodePort
```

namesrv.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: mq-namesrv
  name: rocketmq
  namespace: sanjiang
spec:
  ports:
    - port: 9876
      targetPort: 9876
      nodePort: 32076
      name: namesrv-port
  selector:
    app: mq-namesrv
  type: NodePort
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: rocketmq
  namespace: sanjiang
spec:
  serviceName: rocketmq
  replicas: 1
  selector:
    matchLabels:
      app: mq-namesrv
  template:
    metadata:
      labels:
        app: mq-namesrv
    spec:
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: "app"
                    operator: In
                    values:
                      - mq-namesrv
              topologyKey: "kubernetes.io/hostname"
      containers:
        - name: mq-namesrv
          image: liuyi71sinacom/rocketmq-4.8.0
          imagePullPolicy: IfNotPresent
          command: ["sh", "/usr/local/rocketmq-4.8.0/bin/mqnamesrv"]
          ports:
            - containerPort: 9876
              protocol: TCP
          env:
            - name: JAVA_OPT
              value: "-server -XX:ParallelGCThreads=1 -Xms1g -Xmx1g -Xmn512m"
              #value: "-XX:MaxRAMPercentage=80.0"
          lifecycle:
            postStart:
              exec:
                command: ["/bin/sh", "-c", "touch /tmp/health"]
          livenessProbe:
            exec:
              command: ["test", "-e", "/tmp/health"]
            initialDelaySeconds: 5
            timeoutSeconds: 5
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 9876
            initialDelaySeconds: 15
            timeoutSeconds: 5
            periodSeconds: 20
```

ingress.yaml

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: rocketmq
  namespace: sanjiang
spec:
  ingressClassName: nginx
  rules:
    - host: mq.sj.com
      http:
        paths:
          - path: /
            backend:
              serviceName: console
              servicePort: 8080
```

### 部署服务

```sh
kubectl apply -f .
```
