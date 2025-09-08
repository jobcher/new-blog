---
title: "Argo cd 安装和部署"
date: 2023-12-04
draft: false
authors: "jobcher"
featuredImage: "/images/argo-logo.png"
featuredImagePreview: "/images/argo-logo.png"
images: ["/images/argo-logo.png"]
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---

# Argo cd 安装和部署

Argo CD 是一个为 Kubernetes 而生的，遵循声明式 GitOps 理念的持续部署（CD）工具。Argo CD 可在 Git 存储库更改时自动同步和部署应用程序
![argocd](/images/argocd-1.png)

## 安装

> 前提：你已经安装好了 k8s 环境，我们将在国内的k8s环境下部署argocd  
  
```sh
k3s kubectl create namespace argocd
kubectl apply -n argocd -f https://github.jobcher.com/gh/https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

### 检查是否正常部署
```sh
kubectl get po -n argocd
```
![argocd-1](/images/argocd-1.jpg)  
如果没有错误的情况下应该是全部都runnning，但是如果出现`argocd-repo-server CrashLoopBackOff`错误有以下解决途径：  
1. 使用以下补丁修补了部署。删除后，错误消失，repo 服务器可以启动。  
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: argocd-repo-server
spec:
  template:
    spec:
      securityContext:
        seccompProfile:
          type: RuntimeDefault
```
如果出现`argocd-dex-server imagepullbackoff`错误有以下解决方法：  
```sh
docker pull ghcr.io/dexidp/dex:v2.37.0
docker tag ghcr.io/dexidp/dex:v2.37.0 harbor/dexidp/dex:v2.37.0
docker push harbor/dexidp/adex:v2.37.0
```
**使用自定义镜像**


### 安装 Argo CD CLI

Argo CD CLI 是用于管理 Argo CD 的命令行工具,Mac 系统可以直接使用 brew install 进行安装

```sh
brew install argocd
```

### 发布 Argo CD 服务

默认情况下， Argo CD 服务不对外暴露服务，可以通过 LoadBalancer 或者 NodePort 类型的 Service、Ingress、Kubectl 端口转发等方式将 Argo CD 服务发布到 Kubernetes 集群外部。  
通过 NodePort 服务的方式暴露 Argo CD 到集群外部

```sh
kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "NodePort"}}'
```

### 获取端口号
```sh
kubectl get svc -n argocd
```

### argo nginx-ingress 配置
修改配置
```sh
kubectl patch configmap argocd-cmd-params-cm \
-n argocd \
--type merge \
-p '{"data":{"server.insecure":"true"}}'
```
修改argo cd 服务
```sh
kubectl patch deployment argocd-server \
-n argocd --type='json'  \
-p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--insecure"}]'
```
检查服务是否正常并获取密码
```sh
kubectl get svc -n ingress-nginx nginx-ingress-controller
kubectl get secret argocd-initial-admin-secret -n argocd -o jsonpath="{.data.password}" | base64 -d && echo
```
### 创建 ingress 路由
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: argocd-ingress
  namespace: argocd
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - host: argocd.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: argocd-server
            port:
              number: 80
```
## 使用

![argocd-2](/images/argocd_2.png)

### 获取 Argo CD 密码

默认情况下 `admin`  
 帐号的初始密码是自动生成的，会以明文的形式存储在 `Argo CD` 安装的命名空间中`argocd-initial-admin-secret` 的 `Secret `对象下的 `password`

```sh
kubectl -n argocd get secret \
argocd-initial-admin-secret \
-o jsonpath="{.data.password}" | base64 -d
```

### 命令行可以使用以下方式登录
> argocd login <节点 IP>:<端口>  
  
### 进入UI界面配置
输入地址： 例如：https://10.10.1.1:31751，输入用户名`admin`和密码  
![argocd-2](/images/argocd-2.jpg)  
#### 因为我是私有化仓库部署所以要配置私有化gitlab仓库  
![argocd-3](/images/argocd-3.jpg)  
![argocd-4](/images/argocd-4.jpg)  
![argocd-5](/images/argocd-5.jpg)  
## 部署服务
### 创建yaml文件
我提前在代码仓库写好了nginx的yaml文件  
`deployment.yaml`  
```yml
apiVersion: apps/v1	#与k8s集群版本有关，使用 kubectl api-versions 即可查看当前集群支持的版本
kind: Deployment	#该配置的类型，我们使用的是 Deployment
metadata:	        #译名为元数据，即 Deployment 的一些基本属性和信息
  name: nginx-deployment	#Deployment 的名称
  labels:	    #标签，可以灵活定位一个或多个资源，其中key和value均可自定义，可以定义多组，目前不需要理解
    app: nginx	#为该Deployment设置key为app，value为nginx的标签
spec:	        #这是关于该Deployment的描述，可以理解为你期待该Deployment在k8s中如何使用
  replicas: 1	#使用该Deployment创建一个应用程序实例
  selector:	    #标签选择器，与上面的标签共同作用，目前不需要理解
    matchLabels: #选择包含标签app:nginx的资源
      app: nginx
  template:	    #这是选择或创建的Pod的模板
    metadata:	#Pod的元数据
      labels:	#Pod的标签，上面的selector即选择包含标签app:nginx的Pod
        app: nginx
    spec:	    #期望Pod实现的功能（即在pod中部署）
      containers:	#生成container，与docker中的container是同一种
      - name: nginx	#container的名称
        image: nginx:1.7.9	#使用镜像nginx:1.7.9创建container，该container默认80访问
        livenessProbe: # 存活探针检测
          httpGet:
            path: /nginx_status
            port: 80
            scheme: HTTP
          timeoutSeconds: 1
          periodSeconds: 10
          successThreshold: 1
          failureThreshold: 3
        readinessProbe: # 就绪探针检测
          httpGet:
            path: /nginx_status
            port: 80
            scheme: HTTP
          timeoutSeconds: 1
          periodSeconds: 10
          successThreshold: 1
          failureThreshold: 3
        startupProbe: # 启动探针检测
          httpGet:
            path: /nginx_status
            port: 80
            scheme: HTTP
          timeoutSeconds: 1
          periodSeconds: 10
          successThreshold: 1
          failureThreshold: 3
```
`svc.yaml`
```yml
apiVersion: v1
kind: Service
metadata:
  name: nginx-service	#Service 的名称
  labels:     	#Service 自己的标签
    app: nginx	#为该 Service 设置 key 为 app，value 为 nginx 的标签
spec:	    #这是关于该 Service 的定义，描述了 Service 如何选择 Pod，如何被访问
  selector:	    #标签选择器
    app: nginx	#选择包含标签 app:nginx 的 Pod
  ports:
  - name: nginx-port	#端口的名字
    protocol: TCP	    #协议类型 TCP/UDP
    port: 80	        #集群内的其他容器组可通过 80 端口访问 Service
    nodePort: 32600   #通过任意节点的 32600 端口访问 Service
    targetPort: 80	#将请求转发到匹配 Pod 的 80 端口
  type: NodePort	#Serive的类型，ClusterIP/NodePort/LoaderBalancer
```
### 创建app
![argocd-6](/images/argocd-6.jpg)  
![argocd-7](/images/argocd-7.jpg)  
![argocd-8](/images/argocd-8.jpg)  
![argocd-9](/images/argocd-9.jpg)  
![argocd-10](/images/argocd-10.jpg)  
### 部署成功
![argocd-11](/images/argocd-11.jpg)  
