---
title: "Kubernetes — k8s 手动安装 1.17.9"
date: 2024-08-08
draft: false
featuredImage: "/images/kubernetes0install.png"
featuredImagePreview: "/images/kubernetes0install.png"
images: ["/images/kubernetes0install.png"]
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---
## 背景
已经2024年了， k8s已经更新到 1.30.x的版本了，但是还有很多公司还在使用1.17.9版本，那么我们今天就来手动安装一下1.17.9版本的k8s。

## 安装
我们在测试`centos`服务器`192.168.40.1`安装单节点 Kubernetes 集群（Master 节点）使用 kubeadm 是一个相对直接的过程。
### 前提条件
- 确保主机满足以下要求：
  - 操作系统：CentOS 7.x 或更高版本
  - 内存：至少 2 GB 内存
  - 磁盘空间：至少 20 GB 磁盘空间
  - 网络：至少 2 个网络接口
  
### 配置主机名和 IP
```sh
sudo hostnamectl set-hostname k8s
echo "192.168.40.1 k8s" | sudo tee -a /etc/hosts
```

### 更新系统
1. 切换镜像源
```sh
bash <(curl -sSL https://www.jobcher.com/ChangeMirrors.sh)
```
2. 更新系统
```sh
sudo yum update -y
```
### 禁用 SELinux
```sh
sudo setenforce 0
sudo sed -i --follow-symlinks 's/^SELINUX=enforcing/SELINUX=permissive/' /etc/selinux/config
```
### 禁用 Swap
```sh
sudo swapoff -a
sudo sed -i '/swap/d' /etc/fstab
```
### 修改 /etc/sysctl.conf
```sh
# 如果有配置，则修改
sed -i "s#^net.ipv4.ip_forward.*#net.ipv4.ip_forward=1#g"  /etc/sysctl.conf
sed -i "s#^net.bridge.bridge-nf-call-ip6tables.*#net.bridge.bridge-nf-call-ip6tables=1#g"  /etc/sysctl.conf
sed -i "s#^net.bridge.bridge-nf-call-iptables.*#net.bridge.bridge-nf-call-iptables=1#g"  /etc/sysctl.conf
sed -i "s#^net.ipv6.conf.all.disable_ipv6.*#net.ipv6.conf.all.disable_ipv6=1#g"  /etc/sysctl.conf
sed -i "s#^net.ipv6.conf.default.disable_ipv6.*#net.ipv6.conf.default.disable_ipv6=1#g"  /etc/sysctl.conf
sed -i "s#^net.ipv6.conf.lo.disable_ipv6.*#net.ipv6.conf.lo.disable_ipv6=1#g"  /etc/sysctl.conf
sed -i "s#^net.ipv6.conf.all.forwarding.*#net.ipv6.conf.all.forwarding=1#g"  /etc/sysctl.conf
# 可能没有，追加
echo "net.ipv4.ip_forward = 1" >> /etc/sysctl.conf
echo "net.bridge.bridge-nf-call-ip6tables = 1" >> /etc/sysctl.conf
echo "net.bridge.bridge-nf-call-iptables = 1" >> /etc/sysctl.conf
echo "net.ipv6.conf.all.disable_ipv6 = 1" >> /etc/sysctl.conf
echo "net.ipv6.conf.default.disable_ipv6 = 1" >> /etc/sysctl.conf
echo "net.ipv6.conf.lo.disable_ipv6 = 1" >> /etc/sysctl.conf
echo "net.ipv6.conf.all.forwarding = 1"  >> /etc/sysctl.conf
# 执行命令以应用
sysctl -p
```
### 安装docker
因为k8s 1.17.9 不适用于高版本docker，所以我们要下载指定版本的docker
```sh
sudo yum remove docker*
sudo yum install -y yum-utils
#配置docker yum 源
sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
#安装docker 19.03.9
yum install -y docker-ce-3:19.03.9-3.el7.x86_64  docker-ce-cli-3:19.03.9-3.el7.x86_64 containerd.io

#启动服务
systemctl start docker
systemctl enable docker

sudo systemctl daemon-reload
sudo systemctl restart docker
```

### 拉取镜像
由于国内访问 k8s.gcr.io 仓库速度较慢或被阻止，您可以使用国内的镜像源来替代这些镜像地址。以下是一个解决方案，使用阿里云的镜像源来替代默认的 k8s.gcr.io 镜像源。
```sh
# 拉取阿里云的镜像
sudo docker pull registry.aliyuncs.com/google_containers/kube-apiserver:v1.17.9
sudo docker pull registry.aliyuncs.com/google_containers/kube-controller-manager:v1.17.9
sudo docker pull registry.aliyuncs.com/google_containers/kube-scheduler:v1.17.9
sudo docker pull registry.aliyuncs.com/google_containers/kube-proxy:v1.17.9
sudo docker pull registry.aliyuncs.com/google_containers/pause:3.1
sudo docker pull registry.aliyuncs.com/google_containers/etcd:3.4.3-0
sudo docker pull registry.aliyuncs.com/google_containers/coredns:1.6.5

# 重新标记镜像
sudo docker tag registry.aliyuncs.com/google_containers/kube-apiserver:v1.17.9 k8s.gcr.io/kube-apiserver:v1.17.9
sudo docker tag registry.aliyuncs.com/google_containers/kube-controller-manager:v1.17.9 k8s.gcr.io/kube-controller-manager:v1.17.9
sudo docker tag registry.aliyuncs.com/google_containers/kube-scheduler:v1.17.9 k8s.gcr.io/kube-scheduler:v1.17.9
sudo docker tag registry.aliyuncs.com/google_containers/kube-proxy:v1.17.9 k8s.gcr.io/kube-proxy:v1.17.9
sudo docker tag registry.aliyuncs.com/google_containers/pause:3.1 k8s.gcr.io/pause:3.1
sudo docker tag registry.aliyuncs.com/google_containers/etcd:3.4.3-0 k8s.gcr.io/etcd:3.4.3-0
sudo docker tag registry.aliyuncs.com/google_containers/coredns:1.6.5 k8s.gcr.io/coredns:1.6.5
```

### 安装 kubeadm、kubelet 和 kubectl
1. 配置阿里云的 Kubernetes 仓库:
```sh
cat <<EOF | sudo tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
```
2. 安装特定版本的 kubeadm、kubelet 和 kubectl：
```sh
sudo yum install -y kubelet-1.17.9 kubeadm-1.17.9 kubectl-1.17.9 --disableexcludes=kubernetes
sudo systemctl enable --now kubelet
```

### 初始化集群
1. 初始化集群：
```sh
sudo kubeadm init --kubernetes-version=v1.17.9 --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=192.168.40.1
```
2. 配置 kubectl：
```sh
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

### 安装 Pod 网络插件（Flannel）
这里由于版本关系，不能下载最新的flannel，所以我们下载低版本的flannel  
```sh
# 拉取 Flannel 镜像
sudo docker pull quay.io/coreos/flannel:v0.12.0-amd64

# 重新标记 Flannel 镜像
sudo docker tag quay.io/coreos/flannel:v0.12.0-amd64 quay.io/coreos/flannel:v0.12.0-amd64
```
创建 Flannel 配置文件
```sh
wget https://raw.githubusercontent.com/coreos/flannel/v0.12.0/Documentation/kube-flannel.yml
kubectl apply -f kube-flannel.yml
```

## 验证安装
查看节点状态：
```sh
kubectl get nodes
```
查看所有 Pod 状态：
```sh
kubectl get pods -A
```
