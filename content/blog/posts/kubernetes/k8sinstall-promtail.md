---
title: "Kubernetes — promtail + loki + grafana 日志系统部署"
date: 2026-01-08
draft: false
featuredImage: "/images/promtail.jpg"
featuredImagePreview: "/images/promtail.jpg"
images: ["/images/promtail.jpg"]
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---
## 背景
在k8s 部署一套 promtail + loki + grafana 日志系统。日志由 **Promtail** 从 Kubernetes 集群中收集并发送到 Loki。Promtail 会提取以下标签：  
- `namespace`: Pod 所在的命名空间
- `pod_name`: Pod 名称
- `deployment_name`: Deployment 名称（从 Pod 的 `app` 标签或 Pod 名称中提取）
- `container`: 容器名称

## 部署 loki
- **Deployment**: Loki 主服务
- **Service**: 提供集群内部访问
- **ConfigMap**: Loki 配置文件
- **PVC**: 数据持久化存储（10Gi）  
- **PV**: 数据存储类

### 创建 YAML 文件

首先创建 namespace:
```bash
kubectl create namespace logging
```

创建 loki 目录并创建 pv.yaml
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: loki-data-pv
spec:
  capacity:
    storage: 10Gi # 存储大小
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: host-loki
  hostPath:
    path: /data/loki # 存储位置
    type: DirectoryOrCreate
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
                - worker-1 # 改为你要存储的k8s节点
```
创建pvc.yaml
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: loki-data
  namespace: logging
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: host-loki
```
创建configmap.yaml
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: loki-config
  namespace: logging
data:
  loki.yaml: |
    auth_enabled: false

    server:
      http_listen_port: 3100
      grpc_listen_port: 9096

    common:
      path_prefix: /loki
      storage:
        filesystem:
          chunks_directory: /loki/chunks
          rules_directory: /loki/rules
      replication_factor: 1
      ring:
        instance_addr: 127.0.0.1
        kvstore:
          store: inmemory

    schema_config:
      configs:
        - from: 2020-10-24
          store: boltdb-shipper
          object_store: filesystem
          schema: v11
          index:
            prefix: index_
            period: 24h

    ruler:
      alertmanager_url: http://localhost:9093

    analytics:
      reporting_enabled: false

    limits_config:
      ingestion_rate_mb: 16
      ingestion_burst_size_mb: 32
      max_query_length: 721h
      max_query_parallelism: 32
      max_streams_per_user: 10000
      max_line_size: 0
      max_query_series: 500
      reject_old_samples: true
      reject_old_samples_max_age: 168h

    chunk_store_config:
      max_look_back_period: 0s

    table_manager:
      retention_deletes_enabled: true
      retention_period: 720h

    compactor:
      working_directory: /loki/compactor
      shared_store: filesystem
      compaction_interval: 10m
      retention_enabled: true
      retention_delete_delay: 2h
      retention_delete_worker_count: 150
```
创建 deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: loki
  namespace: logging
spec:
  replicas: 1
  selector:
    matchLabels:
      app: loki
  template:
    metadata:
      labels:
        app: loki
    spec:
      securityContext:
        fsGroup: 10001
      initContainers:
        - name: init-storage
          image: busybox:1.36
          securityContext:
            runAsUser: 0
          command:
            - sh
            - -c
            - |
              mkdir -p /loki/chunks /loki/rules /loki/compactor
              chown -R 10001:10001 /loki
              chmod -R 755 /loki
          volumeMounts:
            - name: storage
              mountPath: /loki
      containers:
        - name: loki
          image: grafana/loki:2.9.2
          securityContext:
            runAsUser: 10001
            runAsNonRoot: true
            readOnlyRootFilesystem: false
          ports:
            - containerPort: 3100
              name: http
            - containerPort: 9096
              name: grpc
          args:
            - -config.file=/etc/loki/loki.yaml
          volumeMounts:
            - name: config
              mountPath: /etc/loki
            - name: storage
              mountPath: /loki
          livenessProbe:
            httpGet:
              path: /ready
              port: 3100
            initialDelaySeconds: 45
            periodSeconds: 30
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /ready
              port: 3100
            initialDelaySeconds: 15
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "2Gi"
              cpu: "1000m"
      volumes:
        - name: config
          configMap:
            name: loki-config
        - name: storage
          persistentVolumeClaim:
            claimName: loki-data
```
创建service.yaml
```yaml
apiVersion: v1
kind: Service
metadata:
  name: loki
  namespace: logging
spec:
  type: ClusterIP
  selector:
    app: loki
  ports:
    - port: 3100
      targetPort: 3100
      protocol: TCP
      name: http
    - port: 9096
      targetPort: 9096
      protocol: TCP
      name: grpc
```

### 部署
1. 创建 PV 和 PVC:
```bash
kubectl apply -f pv.yaml
kubectl apply -f pvc.yaml
```

2. 创建 ConfigMap:
```bash
kubectl apply -f configmap.yaml
```

3. 部署 Deployment:
```bash
kubectl apply -f deployment.yaml
```

4. 创建 Service:
```bash
kubectl apply -f service.yaml
```

5. 验证部署状态:
```bash
# 检查 Loki Pod 状态
kubectl get pods -n logging -l app=loki

# 查看 Loki 日志
kubectl logs -n logging -l app=loki --tail=50
```

## 部署 promtail
Promtail 是 Grafana Loki 的日志收集代理，用于从 Kubernetes 集群中收集容器日志并发送到 Loki。  
- **DaemonSet**: 在每个节点上运行一个 Promtail Pod，收集节点上的容器日志
- **ConfigMap**: 包含 Promtail 配置，定义日志收集规则和 Loki 输出
- **ServiceAccount & RBAC**: 提供访问 Kubernetes API 的权限，用于获取 Pod 元数据  
  
### 创建 YAML 文件

创建 promtail 目录并创建 configmap.yaml
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: promtail-config
  namespace: logging
data:
  promtail.yaml: |
    server:
      http_listen_port: 3101
      grpc_listen_port: 9096

    positions:
      filename: /tmp/positions.yaml

    clients:
      - url: http://loki.logging.svc.cluster.local:3100/loki/api/v1/push

    scrape_configs:
      - job_name: kubernetes-pods
        kubernetes_sd_configs:
          - role: pod
        pipeline_stages:
          - docker: {}
        relabel_configs:
          # 设置日志文件路径
          - source_labels:
              - __meta_kubernetes_namespace
              - __meta_kubernetes_pod_name
              - __meta_kubernetes_pod_uid
            separator: _
            target_label: __tmp_pod_path
          - source_labels:
              - __tmp_pod_path
              - __meta_kubernetes_pod_container_name
            separator: /
            target_label: __path__
            replacement: /var/log/pods/$1/$2/*.log
          # 提取 namespace 标签
          - source_labels:
              - __meta_kubernetes_namespace
            target_label: namespace
          # 提取 pod_name 标签
          - source_labels:
              - __meta_kubernetes_pod_name
            target_label: pod_name
          # 提取 deployment_name 标签（优先从 app 标签获取）
          - source_labels:
              - __meta_kubernetes_pod_label_app
            target_label: deployment_name
            regex: (.+)
          # 如果 app 标签不存在，从 pod 名称中提取（格式：deployment-name-replicaset-hash）
          - source_labels:
              - __meta_kubernetes_pod_name
            target_label: deployment_name
            regex: '^(.+?)-[0-9a-z]+-[0-9a-z]+$'
            replacement: '${1}'
            action: replace
          # 提取 container 标签
          - source_labels:
              - __meta_kubernetes_pod_container_name
            target_label: container
          # 只保留有效的 pod 日志
          - action: keep
            source_labels:
              - __meta_kubernetes_pod_name
              - __meta_kubernetes_pod_node_name
              - __meta_kubernetes_namespace
          # 移除所有 __meta_kubernetes 前缀的标签
          - action: labeldrop
            regex: '__meta_kubernetes.*'
```
创建 daemonset.yaml
```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: promtail
  namespace: logging
spec:
  selector:
    matchLabels:
      app: promtail
  template:
    metadata:
      labels:
        app: promtail
    spec:
      serviceAccountName: promtail
      tolerations:
        - effect: NoSchedule
          operator: Exists
        - effect: NoExecute
          operator: Exists
      containers:
        - name: promtail
          image: grafana/promtail:2.9.2
          args:
            - -config.file=/etc/promtail/promtail.yaml
          ports:
            - name: http-metrics
              containerPort: 3101
          env:
            - name: HOSTNAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
          volumeMounts:
            - name: config
              mountPath: /etc/promtail
            - name: varlog
              mountPath: /var/log
              readOnly: true
            - name: varlibdockercontainers
              mountPath: /var/lib/docker/containers
              readOnly: true
            - name: positions
              mountPath: /tmp
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "256Mi"
              cpu: "200m"
          securityContext:
            runAsUser: 0
            runAsGroup: 0
            runAsNonRoot: false
      volumes:
        - name: config
          configMap:
            name: promtail-config
        - name: varlog
          hostPath:
            path: /var/log
        - name: varlibdockercontainers
          hostPath:
            path: /var/lib/docker/containers
        - name: positions
          emptyDir: {}
```
创建 serviceaccount.yaml
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: promtail
  namespace: logging

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: promtail
rules:
  - apiGroups: [""]
    resources:
      - nodes
      - nodes/proxy
      - services
      - endpoints
      - pods
    verbs: ["get", "list", "watch"]
  - apiGroups:
      - ""
    resources:
      - configmaps
    verbs: ["get"]

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: promtail
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: promtail
subjects:
  - kind: ServiceAccount
    name: promtail
    namespace: logging
```

### 部署
1. 创建 ConfigMap:

```bash
kubectl apply -f configmap.yaml
```

2. 创建 ServiceAccount 和 RBAC:

```bash
kubectl apply -f serviceaccount.yaml
```

3. 部署 DaemonSet:

```bash
kubectl apply -f daemonset.yaml
```

4. 验证部署状态:

```bash
# 检查 Promtail Pod 是否在所有节点上运行
kubectl get pods -n logging -l app=promtail

# 查看 Promtail 日志
kubectl logs -n logging -l app=promtail --tail=50
```

## 部署 grafana
Grafana 用于可视化 Loki 中的日志数据，提供强大的查询和仪表板功能。  
- **Deployment**: Grafana 主服务
- **Service**: 提供集群内部访问，你也可以 nodeport 直接访问
- **Ingress**: 提供外部访问（通过 Traefik）
- **ConfigMap**: 数据源配置（自动配置 Loki 数据源）
- **PVC**: 数据持久化存储（10Gi，保存仪表板和配置）  
  
### 创建 YAML 文件

创建 grafana 目录并创建 configmap.yaml
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: grafana-datasources
  namespace: logging
data:
  datasources.yaml: |
    apiVersion: 1
    datasources:
      - name: Loki
        type: loki
        access: proxy
        url: http://loki.logging.svc.cluster.local:3100
        isDefault: false
        editable: true
        jsonData:
          maxLines: 1000
          derivedFields:
            - datasourceUid: loki
              matcherRegex: "traceID=(\\w+)"
              name: TraceID
              url: '$${__value.raw}'
```
创建 deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grafana
  namespace: logging
spec:
  replicas: 1
  selector:
    matchLabels:
      app: grafana
  template:
    metadata:
      labels:
        app: grafana
    spec:
      initContainers:
        - name: init-storage
          image: busybox:1.36
          securityContext:
            runAsUser: 0
          command:
            - sh
            - -c
            - |
              mkdir -p /var/lib/grafana/plugins /var/lib/grafana/data /var/lib/grafana/logs
              chown -R root:root /var/lib/grafana
              chmod -R 755 /var/lib/grafana
          volumeMounts:
            - name: storage
              mountPath: /var/lib/grafana
      containers:
        - name: grafana
          image: grafana/grafana:10.2.2
          securityContext:
            runAsUser: 0
            runAsNonRoot: false
            readOnlyRootFilesystem: false
          ports:
            - containerPort: 3000
              name: http
          env:
            - name: GF_SECURITY_ADMIN_USER
              value: admin
            - name: GF_SECURITY_ADMIN_PASSWORD
              value: admin  # 建议通过 Secret 管理
            - name: GF_INSTALL_PLUGINS
              value: ""
            - name: GF_SERVER_ROOT_URL
              value: "%(protocol)s://%(domain)s:%(http_port)s/"
            - name: GF_SERVER_SERVE_FROM_SUB_PATH
              value: "false"
          volumeMounts:
            - name: storage
              mountPath: /var/lib/grafana
            - name: datasources
              mountPath: /etc/grafana/provisioning/datasources
          livenessProbe:
            httpGet:
              path: /api/health
              port: 3000
            initialDelaySeconds: 60
            periodSeconds: 30
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /api/health
              port: 3000
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
      volumes:
        - name: storage
          persistentVolumeClaim:
            claimName: grafana-data
        - name: datasources
          configMap:
            name: grafana-datasources
```
创建 pv.yaml
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: grafana-data-pv
spec:
  capacity:
    storage: 10Gi # 存储大小
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: host-grafana
  hostPath:
    path: /data/grafana # 存储目录
    type: DirectoryOrCreate
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
                - worker-1 # 修改成你自己的节点
```
创建 pvc.yaml
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: grafana-data
  namespace: logging
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
  storageClassName: host-grafana
```
创建 service.yaml
```yaml
apiVersion: v1
kind: Service
metadata:
  name: grafana
  namespace: logging
spec:
  type: ClusterIP
  selector:
    app: grafana
  ports:
    - port: 3000
      targetPort: 3000
      protocol: TCP
      name: http
```
创建 ingress.yaml
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: grafana-ingress
  namespace: logging
  annotations:
    traefik.ingress.kubernetes.io/router.entrypoints: web,websecure
    traefik.ingress.kubernetes.io/router.tls: "true"
spec:
  tls:
    - hosts:
        - grafana-dev.jobcher.com           #改成你自己的域名
      secretName: jobcher-com-tls           #改成你自己的SSL证书
  rules:
    - host: grafana-dev.jobcher.com         #改成你自己的域名
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: grafana
                port:
                  number: 3000
```

### 部署
1. 创建 PV和PVC:
```bash
kubectl apply -f pv.yaml
kubectl apply -f pvc.yaml
```

2. 创建 ConfigMap（数据源配置）:
```bash
kubectl apply -f configmap.yaml
```

3. 部署 Deployment:
```bash
kubectl apply -f deployment.yaml
```

4. 创建 Service:
```bash
kubectl apply -f service.yaml
```

5. 创建 Ingress（可选，用于外部访问）:
```bash
kubectl apply -f ingress.yaml
```

6. 验证部署状态:

```bash
# 检查所有组件状态
kubectl get pods -n logging

# 检查 Loki 服务
kubectl get svc -n logging

# 查看 Loki 日志
kubectl logs -n logging -l app=loki --tail=50

# 查看 Grafana 日志
kubectl logs -n logging -l app=grafana --tail=50
```

## 验证和测试

### 验证 Loki 是否正常工作
```bash
# 检查 Loki 健康状态
kubectl exec -n logging -it deployment/loki -- wget -q -O - http://localhost:3100/ready

# 查询 Loki 中的日志流
kubectl exec -n logging -it deployment/loki -- wget -q -O - "http://localhost:3100/loki/api/v1/label/namespace/values"
```

### 在 Grafana 中查看日志
1. 访问 Grafana（通过 Ingress 或端口转发）:
```bash
# 端口转发（如果使用 ClusterIP）
kubectl port-forward -n logging svc/grafana 3000:3000
```
2. 使用默认账号登录：`admin` / `admin`
3. 进入 **Explore** 页面，选择 **Loki** 数据源
4. 使用 LogQL 查询日志，例如：
   - `{namespace="default"}` - 查看 default 命名空间的日志
   - `{pod_name="your-pod-name"}` - 查看特定 Pod 的日志
   - `{deployment_name="your-deployment"}` - 查看特定 Deployment 的日志

## 注意事项

### 存储配置
- 确保节点上的 `/var/log/pods` 路径存在（Kubernetes 1.14+ 标准日志路径）
- 根据实际需求调整 PV 的存储大小和节点选择
- 如果使用 containerd 而不是 Docker，日志路径格式相同，但可能需要在配置中调整

### 安全配置
- Promtail 需要以 root 用户运行（uid 0）才能访问节点上的日志文件
- Grafana 管理员密码建议通过 Secret 管理，而不是直接写在 Deployment 中
- 生产环境建议启用 Loki 的认证功能

### 性能优化
- 根据集群规模调整资源限制（CPU、内存）
- 根据日志量调整 Loki 的 `ingestion_rate_mb` 和 `ingestion_burst_size_mb`
- 定期清理旧日志，根据 `retention_period` 配置自动删除

### 标签提取逻辑
- `deployment_name` 标签的提取逻辑：
  - 优先使用 Pod 的 `app` 标签值
  - 如果 `app` 标签不存在，则从 Pod 名称中提取（格式：`deployment-name-replicaset-hash`）
- 确保 Pod 有正确的标签，以便 Promtail 正确提取元数据

### 故障排查
- 如果 Promtail 无法收集日志，检查 Pod 是否有权限访问 `/var/log/pods`
- 如果 Loki 无法接收日志，检查网络连接和 Service 配置
- 使用 `kubectl logs` 查看各组件的日志进行排查