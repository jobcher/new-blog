---
title: "Kubernetes — skywalking + banyanDB APM监控部署"
date: 2026-01-09
draft: false
featuredImage: "/images/skywalking.webp"
featuredImagePreview: "/images/skywalking.webp"
images: ["/images/skywalking.webp"]
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---
## 介绍

Apache SkyWalking 是一个应用性能监控（APM）系统，用于分布式系统的监控、追踪和诊断。它提供了完整的可观测性解决方案，帮助开发者和运维人员快速定位和解决分布式系统中的性能问题。

### 组件说明

- **OAP (Observability Analysis Platform)**: 核心分析平台，负责数据收集、分析和存储
- **UI**: Web 界面，用于可视化和查询监控数据
- **BanyanDB**: 高性能时序数据库，作为 SkyWalking 的存储后端
- **etcd**: 分布式键值存储，BanyanDB 使用它来存储元数据

### 功能特性

- **分布式追踪**: 自动收集和关联分布式系统的调用链，支持跨服务追踪
- **服务拓扑**: 可视化服务之间的依赖关系，实时展示服务调用图
- **性能指标**: 收集服务的性能指标（延迟、吞吐量、错误率等）
- **日志关联**: 将日志与追踪数据关联，通过 TraceID 快速定位问题
- **告警机制**: 支持基于指标的告警规则配置

## 部署

### 前置要求

在开始部署之前，请确保满足以下要求：

1. **Kubernetes 集群**: 版本 1.20+
2. **命名空间**: 已创建 `logging` 命名空间（或根据实际情况修改）
3. **存储**: 确保节点有足够的存储空间（建议至少 10Gi）
4. **网络**: 确保 Pod 之间可以正常通信

> **提示**: 本文档中的所有资源都部署在 `logging` 命名空间中，如需使用其他命名空间，请修改相应的 YAML 文件。

### 部署 BanyanDB

BanyanDB 是 SkyWalking 的存储后端，需要先部署 BanyanDB 和 etcd。

#### 1. 部署 etcd

etcd 用于存储 BanyanDB 的元数据。

**创建 etcd 持久化存储卷**

创建文件 `banyandb-etcd-pv.yaml`:
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: banyandb-etcd-data-pv
spec:
  capacity:
    storage: 5Gi # 存储大小
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: host-banyandb-etcd
  hostPath:
    path: /data/banyandb-etcd # 存储位置
    type: DirectoryOrCreate
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values: 
                - worker-1  # 修改为你要存储的节点名称
```

> **注意**: 请根据实际情况修改 `worker-1` 为你的节点名称，并确保该节点存在 `/data/banyandb-etcd` 目录或具有创建权限。

**创建 etcd 持久化存储声明**

创建文件 `banyandb-etcd-pvc.yaml`:
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: banyandb-etcd-data
  namespace: logging
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi  #存储大小
  storageClassName: host-banyandb-etcd
```

**创建 etcd 服务**

创建文件 `banyandb-etcd-service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: banyandb-etcd
  namespace: logging
spec:
  clusterIP: None
  ports:
    - name: client
      port: 2379
      targetPort: 2379
    - name: peer
      port: 2380
      targetPort: 2380
  selector:
    app: banyandb-etcd
```

**创建 etcd StatefulSet**

创建文件 `banyandb-etcd-statefulset.yaml`:
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: banyandb-etcd
  namespace: logging
spec:
  serviceName: banyandb-etcd
  replicas: 1
  selector:
    matchLabels:
      app: banyandb-etcd
  template:
    metadata:
      labels:
        app: banyandb-etcd
    spec:
      containers:
        - name: etcd
          image: quay.io/coreos/etcd:v3.5.9
          ports:
            - name: client
              containerPort: 2379
            - name: peer
              containerPort: 2380
          env:
            - name: ETCD_NAME
              value: "banyandb-etcd-0"
            - name: ETCD_DATA_DIR
              value: /etcd-data
            - name: ETCD_LISTEN_CLIENT_URLS
              value: "http://0.0.0.0:2379"
            - name: ETCD_ADVERTISE_CLIENT_URLS
              value: "http://banyandb-etcd-0.banyandb-etcd:2379"
            - name: ETCD_LISTEN_PEER_URLS
              value: "http://0.0.0.0:2380"
            - name: ETCD_INITIAL_ADVERTISE_PEER_URLS
              value: "http://banyandb-etcd-0.banyandb-etcd:2380"
            - name: ETCD_INITIAL_CLUSTER
              value: "banyandb-etcd-0=http://banyandb-etcd-0.banyandb-etcd:2380"
            - name: ETCD_INITIAL_CLUSTER_TOKEN
              value: "banyandb-etcd-cluster"
            - name: ETCD_INITIAL_CLUSTER_STATE
              value: "new"
          volumeMounts:
            - name: data
              mountPath: /etcd-data
          livenessProbe:
            httpGet:
              path: /health
              port: 2379
            initialDelaySeconds: 30
            periodSeconds: 30
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /health
              port: 2379
            initialDelaySeconds: 10
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
        - name: data
          persistentVolumeClaim:
            claimName: banyandb-etcd-data
```

#### 2. 部署 BanyanDB

**创建 BanyanDB 持久化存储卷**

创建文件 `banyandb-pv.yaml`:
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: banyandb-data-pv
spec:
  capacity:
    storage: 5Gi #存储大小
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: host-banyandb
  hostPath:
    path: /data/banyandb #存储位置
    type: DirectoryOrCreate
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: kubernetes.io/hostname
              operator: In
              values:
                - worker-1  # 修改为你要存储的节点名称
```

> **注意**: 请根据实际情况修改 `worker-1` 为你的节点名称，并确保该节点存在 `/data/banyandb` 目录或具有创建权限。

**创建 BanyanDB 持久化存储声明**

创建文件 `banyandb-pvc.yaml`:
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: banyandb-data
  namespace: logging
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi #存储大小
  storageClassName: host-banyandb
```

**创建 BanyanDB 服务**

创建文件 `banyandb-service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: banyandb
  namespace: logging
spec:
  ports:
    - name: grpc
      port: 17912
      targetPort: 17912
    - name: http
      port: 17913
      targetPort: 17913
  selector:
    app: banyandb
```

**创建 BanyanDB StatefulSet**

创建文件 `banyandb-statefulset.yaml`:
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: banyandb
  namespace: logging
spec:
  serviceName: banyandb
  replicas: 1
  selector:
    matchLabels:
      app: banyandb
  template:
    metadata:
      labels:
        app: banyandb
    spec:
      initContainers:
        - name: wait-for-etcd
          image: busybox:1.35
          command: ['sh', '-c']
          args:
            - |
              until nc -z banyandb-etcd 2379; do
                echo "Waiting for etcd to be ready..."
                sleep 2
              done
              echo "etcd is ready!"
      containers:
        - name: banyandb
          image: apache/skywalking-banyandb:0.9.0
          args: ["standalone"]
          ports:
            - name: grpc
              containerPort: 17912
            - name: http
              containerPort: 17913
          env:
            - name: BANYANDB_STANDALONE
              value: "true"
            - name: BANYANDB_ETCD_ENDPOINTS
              value: "http://banyandb-etcd:2379"
            - name: BANYANDB_GRPC_HOST
              value: "0.0.0.0"
            - name: BANYANDB_GRPC_PORT
              value: "17912"
            - name: BANYANDB_HTTP_HOST
              value: "0.0.0.0"
            - name: BANYANDB_HTTP_PORT
              value: "17913"
            - name: BANYANDB_DATA_DIR
              value: "/data"
          livenessProbe:
            tcpSocket:
              port: 17913
            initialDelaySeconds: 60
            periodSeconds: 30
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            tcpSocket:
              port: 17913
            initialDelaySeconds: 30
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
          volumeMounts:
            - name: data
              mountPath: /data
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: banyandb-data
```

#### 3. 部署 BanyanDB 组件

在包含所有 BanyanDB YAML 文件的目录下执行：

```bash
kubectl apply -f .
```

验证部署状态：

```bash
# 检查 etcd 状态
kubectl get pods -n logging | grep banyandb-etcd

# 检查 BanyanDB 状态
kubectl get pods -n logging | grep banyandb

# 查看 BanyanDB 日志
kubectl logs -n logging -l app=banyandb --tail=50
```

> **提示**: 确保 etcd 和 BanyanDB 都处于 `Running` 状态后再继续部署 SkyWalking。

---

### 部署 SkyWalking

#### 1. 部署 OAP Server

OAP (Observability Analysis Platform) 是 SkyWalking 的核心组件，负责数据收集、分析和存储。

**创建 OAP Deployment**

创建文件 `oap-deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: skywalking-oap
  namespace: logging
spec:
  replicas: 1
  selector:
    matchLabels:
      app: skywalking-oap
  template:
    metadata:
      labels:
        app: skywalking-oap
    spec:
      initContainers:
        - name: wait-for-banyandb
          image: busybox:1.35
          command: ['sh', '-c']
          args:
            - |
              echo "Waiting for BanyanDB to be ready..."
              # 等待 BanyanDB 服务可用
              until nc -z banyandb.logging.svc.cluster.local 17912; do
                echo "Waiting for BanyanDB gRPC port..."
                sleep 3
              done
              # 额外等待几秒，确保 gRPC 服务完全启动
              sleep 5
              echo "BanyanDB is ready!"
      containers:
        - name: oap
          image: apache/skywalking-oap-server:10.3.0
          ports:
            - containerPort: 11800
            - containerPort: 12800
          env:
            - name: SW_STORAGE
              value: banyandb
            - name: SW_STORAGE_BANYANDB_TARGETS
              value: banyandb:17912
            - name: SW_STORAGE_BANYANDB_GRPC_TLS_ENABLED
              value: "false"
            - name: JAVA_OPTS
              value: "-Xms1g -Xmx1g"
            # 启用 Log Receiver 模块 - 接收来自应用的日志
            - name: SW_RECEIVER_LOGGING
              value: default
            # 启用 Log Analyzer 模块 - 分析日志数据
            - name: SW_LOG_ANALYZER
              value: default
            # LAL 文件配置
            - name: SW_LOG_LAL_FILES
              value: default
            # 日志接收器 gRPC 处理器配置（与 agent 9.5.0 的 GRPCLogClientAppender 兼容）
            - name: SW_RECEIVER_LOGGING_DEFAULT_HANDLERS
              value: grpc
            # 日志接收器 gRPC 服务配置
            - name: SW_RECEIVER_LOGGING_DEFAULT_GRPC_HOST
              value: 0.0.0.0
            - name: SW_RECEIVER_LOGGING_DEFAULT_GRPC_PORT
              value: "11800"
            # 日志接收器 gRPC 最大消息大小（支持大日志）
            - name: SW_RECEIVER_LOGGING_DEFAULT_GRPC_MAX_MESSAGE_SIZE
              value: "10485760"
            # 日志接收器 gRPC 最大连接数
            - name: SW_RECEIVER_LOGGING_DEFAULT_GRPC_MAX_CONCURRENT_STREAMS
              value: "100"
          livenessProbe:
            httpGet:
              path: /healthcheck
              port: 12800
            initialDelaySeconds: 60
            periodSeconds: 30
            timeoutSeconds: 5
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /healthcheck
              port: 12800
            initialDelaySeconds: 30
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
```

> **说明**: 
> - `SW_STORAGE=banyandb`: 使用 BanyanDB 作为存储后端
> - `SW_RECEIVER_LOGGING`: 启用日志接收功能
> - `SW_LOG_ANALYZER`: 启用日志分析功能
> - `JAVA_OPTS`: 设置 JVM 内存参数，根据实际需求调整

**创建 OAP Service**

创建文件 `oap-service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: skywalking-oap
  namespace: logging
spec:
  ports:
    - name: grpc
      port: 11800
    - name: http
      port: 12800
  selector:
    app: skywalking-oap
```

#### 2. 部署 SkyWalking UI

**创建 UI Deployment**

创建文件 `ui-deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: skywalking-ui
  namespace: logging
spec:
  replicas: 1
  selector:
    matchLabels:
      app: skywalking-ui
  template:
    metadata:
      labels:
        app: skywalking-ui
    spec:
      containers:
        - name: ui
          image: apache/skywalking-ui:10.3.0
          ports:
            - containerPort: 8080
          env:
            - name: SW_OAP_ADDRESS
              value: http://skywalking-oap.logging.svc.cluster.local:12800
```

**创建 UI Service**

创建文件 `ui-service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: skywalking-ui
  namespace: logging
spec:
  type: ClusterIP
  selector:
    app: skywalking-ui
  ports:
    - port: 8080
      targetPort: 8080
      protocol: TCP
      name: http
```

**创建 Ingress（可选）**

如果需要通过域名访问 SkyWalking UI，创建文件 `ingress.yaml`:
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: skywalking-ingress
  namespace: logging
  annotations:
    traefik.ingress.kubernetes.io/router.entrypoints: web,websecure
    traefik.ingress.kubernetes.io/router.tls: "true"
spec:
  tls:
    - hosts:
        - skywalking-dev.jobcher.com        #改为你自己的域名
      secretName: jobcher-com-tls           #改为你自己的证书
  rules:
    - host: skywalking-dev.jobcher.com      #改为你自己的域名
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: skywalking-ui
                port:
                  number: 8080
```

> **注意**: 
> - 请将 `skywalking-dev.jobcher.com` 替换为你的实际域名
> - 请将 `jobcher-com-tls` 替换为你的 TLS 证书 Secret 名称
> - 如果使用其他 Ingress Controller（如 Nginx），请相应修改 annotations

#### 3. 部署 SkyWalking 组件

在包含所有 SkyWalking YAML 文件的目录下执行：

```bash
kubectl apply -f .
```

验证部署状态：

```bash
# 检查 OAP 状态
kubectl get pods -n logging | grep skywalking-oap

# 检查 UI 状态
kubectl get pods -n logging | grep skywalking-ui

# 查看 OAP 日志
kubectl logs -n logging -l app=skywalking-oap --tail=50
```

> **提示**: 如果使用 Ingress，部署完成后可以通过配置的域名访问 SkyWalking UI。也可以通过 `kubectl port-forward` 临时访问：
> ```bash
> kubectl port-forward -n logging svc/skywalking-ui 8080:8080
> ```
> 然后访问 `http://localhost:8080`

---

## 配置 SkyWalking Agent

要在 Java 应用中集成 SkyWalking，需要配置 SkyWalking Agent。本节以 Spring Boot 应用为例。

### 1. 添加 Maven 依赖

在项目的 `pom.xml` 文件中添加以下依赖：
```xml
    <!-- Logback 1.2.3 for SkyWalking compatibility -->
    <dependency>
      <groupId>ch.qos.logback</groupId>
      <artifactId>logback-classic</artifactId>
      <version>1.2.3</version>
    </dependency>
    <dependency>
      <groupId>ch.qos.logback</groupId>
      <artifactId>logback-core</artifactId>
      <version>1.2.3</version>
    </dependency>

    <!-- SkyWalking Logback Toolkit for trace ID -->
    <dependency>
      <groupId>org.apache.skywalking</groupId>
      <artifactId>apm-toolkit-logback-1.x</artifactId>
      <version>9.5.0</version>
    </dependency>
```

> **说明**: 
> - Logback 1.2.3 版本与 SkyWalking Agent 9.5.0 兼容
> - `apm-toolkit-logback-1.x` 用于在日志中自动添加 TraceID

### 2. 配置 Logback

修改 `logback-spring.xml` 文件，添加 SkyWalking 日志 Appender 和 TraceID 支持：
```xml
<!-- SkyWalking Log Appender -->
	<appender name="SKYWALKING" class="org.apache.skywalking.apm.toolkit.log.logback.v1.x.log.GRPCLogClientAppender">
		<encoder class="ch.qos.logback.core.encoder.LayoutWrappingEncoder">
			<layout class="org.apache.skywalking.apm.toolkit.log.logback.v1.x.mdc.TraceIdMDCPatternLogbackLayout">
				<Pattern>%d{yyyy-MM-dd HH:mm:ss.SSS} [%X{tid}] [%thread] %-5level %logger{36} -%msg%n</Pattern>
			</layout>
		</encoder>
	</appender>

    <appender name="console" class="ch.qos.logback.core.ConsoleAppender">
        <filter class="ch.qos.logback.classic.filter.ThresholdFilter">
        <level>Info</level>
        </filter>
        <encoder class="ch.qos.logback.core.encoder.LayoutWrappingEncoder">
        <layout class="org.apache.skywalking.apm.toolkit.log.logback.v1.x.TraceIdPatternLogbackLayout">
            <pattern>%d{yyyy-MM-dd HH:mm:ss.SSS} %contextName [%thread] %-5level %logger{36} :%line [%tid] - %msg%n</pattern>
        </layout>
        </encoder>
    </appender>

    <appender name="file" class="ch.qos.logback.core.rolling.RollingFileAppender">
        <file>${log.path}</file>
        <rollingPolicy class="ch.qos.logback.core.rolling.TimeBasedRollingPolicy">
        <fileNamePattern>${log.path}.%d{yyyy-MM-dd}.zip</fileNamePattern>
        </rollingPolicy>
        <encoder class="ch.qos.logback.core.encoder.LayoutWrappingEncoder">
        <layout class="org.apache.skywalking.apm.toolkit.log.logback.v1.x.TraceIdPatternLogbackLayout">
            <pattern>%date %level [%thread] %logger{36} [%file : %line] [%tid] %msg%n</pattern>
        </layout>
        </encoder>
    </appender>

<!-- 其他配置内容 ... -->

  <root level="info">
    <appender-ref ref="SKYWALKING"/>  <!-- 添加 SkyWalking 日志 Appender -->
    <appender-ref ref="console"/>
    <appender-ref ref="file"/>
  </root>
```

> **说明**: 
> - `SKYWALKING` Appender 会将日志发送到 SkyWalking OAP
> - `TraceIdPatternLogbackLayout` 和 `TraceIdMDCPatternLogbackLayout` 会在日志中自动添加 TraceID
> - 日志格式中的 `[%tid]` 或 `[%X{tid}]` 会显示追踪 ID

### 3. Kubernetes 部署配置

#### 创建应用 Deployment

创建文件 `deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: java-skywalking-test
  namespace: logging
spec:
  replicas: 1
  selector:
    matchLabels:
      app: java-skywalking-test
  template:
    metadata:
      labels:
        app: java-skywalking-test
    spec:
      initContainers:
        - name: skywalking-agent
          image: apache/skywalking-java-agent:9.5.0-java8
          command: ['sh', '-c']
          args:
            - |
              cp -r /skywalking/agent /shared/
          volumeMounts:
            - name: skywalking-agent
              mountPath: /shared
      containers:
        - name: java-skywalking-test
          image: dev/java-skywalking-test:latest  # 修改为你的应用镜像地址
          imagePullPolicy: Always
          ports:
            - containerPort: 9080
              name: http
          env:
            # SkyWalking Agent 配置（通过 JAVA_TOOL_OPTIONS）
            - name: JAVA_TOOL_OPTIONS
              value: >-
                -javaagent:/skywalking/agent/skywalking-agent.jar
                -DSW_AGENT_NAME=java-skywalking-test
                -DSW_AGENT_COLLECTOR_BACKEND_SERVICES=skywalking-oap.logging.svc.cluster.local:11800
                -DSW_LOGGING_LEVEL=INFO
                -DSW_AGENT_AUTHENTICATION=
                -DSW_LOGGING_ENABLED=true
                -DSW_LOGGING_FILE_NAME=skywalking-api.log
                -DSW_LOGGING_DIR=/tmp/logs
                -DSW_AGENT_SAMPLE_N_PER_3_SECS=-1
                -DSW_AGENT_FORCE_RECONNECTION_PERIOD=10
                -DSW_AGENT_IS_OPEN_DEBUGGING_CLASS=true
                -DSW_AGENT_SPAN_LIMIT_PER_SEGMENT=500
                -DSW_AGENT_IGNORE_SUFFIX=.jpg,.jpeg,.js,.css,.png,.bmp,.gif,.ico,.mp3,.mp4,.html,.svg
                -DSW_GRPC_LOG_SERVER_HOST=skywalking-oap.logging.svc.cluster.local
                -DSW_GRPC_LOG_SERVER_PORT=11800
            # SkyWalking Agent 环境变量配置（备用方式）
            - name: SW_AGENT_COLLECTOR_BACKEND_SERVICES
              value: skywalking-oap.logging.svc.cluster.local:11800
            - name: SW_AGENT_NAME
              value: java-skywalking-test
            - name: SW_LOGGING_LEVEL
              value: INFO
            - name: SW_LOGGING_ENABLED
              value: "true"
            - name: SW_AGENT_SAMPLE_N_PER_3_SECS
              value: "-1"
            # 增加每个 segment 的 span 限制，避免 "More than 300 spans" 警告
            - name: SW_AGENT_SPAN_LIMIT_PER_SEGMENT
              value: "500"
            # 忽略静态资源的追踪
            - name: SW_AGENT_IGNORE_SUFFIX
              value: ".jpg,.jpeg,.js,.css,.png,.bmp,.gif,.ico,.mp3,.mp4,.html,.svg"
            # GRPCLogClientAppender 配置 - 日志发送到 OAP
            - name: SW_GRPC_LOG_SERVER_HOST
              value: skywalking-oap.logging.svc.cluster.local
            - name: SW_GRPC_LOG_SERVER_PORT
              value: "11800"
          volumeMounts:
            - name: skywalking-agent
              mountPath: /skywalking
          # 根据实际需求配置资源限制
          # resources:
          #   requests:
          #     memory: "512Mi"
          #     cpu: "250m"
          #   limits:
          #     memory: "1Gi"
          #     cpu: "500m"
      volumes:
        - name: skywalking-agent
          emptyDir: {}
```

> **关键配置说明**:
> - `initContainers`: 使用 init 容器将 SkyWalking Agent 复制到共享卷
> - `JAVA_TOOL_OPTIONS`: 通过 JVM 参数启动 SkyWalking Agent
> - `SW_AGENT_NAME`: 服务名称，在 SkyWalking UI 中显示
> - `SW_AGENT_COLLECTOR_BACKEND_SERVICES`: OAP 服务地址
> - `SW_AGENT_SAMPLE_N_PER_3_SECS=-1`: 采样率，-1 表示采样所有请求
> - `SW_AGENT_SPAN_LIMIT_PER_SEGMENT=500`: 每个 segment 的 span 限制，避免警告
> - `SW_GRPC_LOG_SERVER_HOST/PORT`: 日志发送到 OAP 的地址

**创建应用 Service**

创建文件 `service.yaml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: java-skywalking-test
  namespace: logging
spec:
  type: NodePort
  selector:
    app: java-skywalking-test
  ports:
    - port: 9080
      targetPort: 9080
      protocol: TCP
      name: http
      nodePort: 30980
```

#### 部署应用

在包含应用 YAML 文件的目录下执行：

```bash
kubectl apply -f .
```

验证部署状态：

```bash
# 检查 Pod 状态
kubectl get pods -n logging | grep java-skywalking-test

# 查看应用日志
kubectl logs -n logging -l app=java-skywalking-test --tail=50

# 查看 SkyWalking Agent 日志
kubectl exec -n logging -it <pod-name> -- cat /tmp/logs/skywalking-api.log
```

---

## 验证和测试

### 1. 访问 SkyWalking UI

通过 Ingress 或 port-forward 访问 SkyWalking UI，你应该能看到：

- **服务列表**: 显示已注册的服务（如 `java-skywalking-test`）
- **服务拓扑**: 可视化服务之间的调用关系
- **追踪数据**: 查看详细的调用链信息
- **日志关联**: 通过 TraceID 关联日志和追踪数据

### 2. 测试应用追踪

访问你的应用并执行一些操作，然后在 SkyWalking UI 中：

1. 进入 **Topology** 页面查看服务拓扑
2. 进入 **Trace** 页面查看追踪数据
3. 进入 **Log** 页面查看关联的日志

### 3. 常见问题排查

**问题**: 应用未出现在 SkyWalking UI 中

- 检查 Pod 是否正常运行
- 检查 SkyWalking Agent 日志：`kubectl logs -n logging <pod-name>`
- 确认 `SW_AGENT_COLLECTOR_BACKEND_SERVICES` 配置正确
- 确认 OAP 服务可访问：`kubectl get svc -n logging skywalking-oap`

**问题**: 日志未发送到 SkyWalking

- 检查 Logback 配置是否正确
- 确认 `SW_GRPC_LOG_SERVER_HOST` 和 `SW_GRPC_LOG_SERVER_PORT` 配置正确
- 查看 OAP 日志确认日志接收器是否正常

**问题**: TraceID 未出现在日志中

- 确认已添加 `apm-toolkit-logback-1.x` 依赖
- 检查 Logback 配置中是否使用了 `TraceIdPatternLogbackLayout`

---

## 总结

本文介绍了如何在 Kubernetes 集群中部署 SkyWalking APM 监控系统，包括：

1. ✅ 部署 BanyanDB 和 etcd 作为存储后端
2. ✅ 部署 SkyWalking OAP Server 和 UI
3. ✅ 配置 Java 应用集成 SkyWalking Agent
4. ✅ 配置日志关联和 TraceID 追踪

通过以上配置，你可以获得完整的应用性能监控能力，包括分布式追踪、服务拓扑可视化和日志关联等功能。

> **提示**: 生产环境建议：
> - 根据实际负载调整资源限制
> - 配置持久化存储以确保数据安全
> - 设置合适的采样率以平衡性能和监控需求
> - 配置告警规则以便及时发现问题
