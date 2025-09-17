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
1. 切换镜像源,选择你喜欢的镜像源，我这里选择腾讯云
```sh
bash <(curl -sSL https://linuxmirrors.cn/main.sh)
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

> **注意：** 接下来的步骤有所不同，如果你是初始化集群，继续跟着文档操作，如果你在已有集群的情况下，请跳到[加入集群](#加入集群)

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
## 加入集群
在主节点上执行`kubeadm token create --print-join-command` 命令，获取加入集群的命令，然后在其他节点上执行该命令，即可加入集群。
```sh
kubeadm token create --print-join-command
```
在你需要加入集群的节点上执行该命令，即可加入集群。
### 配置网络组件
创建 `/etc/cni/net.d` 目录
```sh
mkdir -p /etc/cni/net.d
```
下载 flannel 配置文件
```sh
cd /tmp
curl -L -o flannel.tgz   https://github.jobcher.com/gh/https://github.com/flannel-io/cni-plugin/releases/download/v1.1.2/cni-plugin-flannel-windows-amd64-v1.1.2.tgz
tar -xzvf flannel.tgz -C /opt/cni/bin
sudo ln -s /opt/cni/bin/flannel-amd64 /opt/cni/bin/flannel
ls -l /opt/cni/bin/flannel
```
重启 kubelet 服务
```sh
sudo systemctl daemon-reload
sudo systemctl restart kubelet
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

## 其他问题
### 节点状态为 NotReady
执行`systemctl status kubelet` 查看 kubelet 服务状态  
发现`[failed to find plugin "flannel" in path [/opt/cni/bin]]`
```sh
● kubelet.service - kubelet: The Kubernetes Node Agent
   Loaded: loaded (/usr/lib/systemd/system/kubelet.service; enabled; vendor preset: disabled)
  Drop-In: /usr/lib/systemd/system/kubelet.service.d
           └─10-kubeadm.conf
   Active: active (running) since Fri 2025-08-29 10:26:22 CST; 12min ago
     Docs: https://kubernetes.io/docs/
 Main PID: 31569 (kubelet)
    Tasks: 20
   Memory: 38.6M
   CGroup: /system.slice/kubelet.service
           └─31569 /usr/bin/kubelet --bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf --config=/var/lib/kubelet/config.yaml --cgroup-driver=cgroupfs --network-plugin=cni --pod-in...

Aug 29 10:38:24 dev-nbsj-node2-174 kubelet[31569]: "type": "portmap",
Aug 29 10:38:24 dev-nbsj-node2-174 kubelet[31569]: "capabilities": {
Aug 29 10:38:24 dev-nbsj-node2-174 kubelet[31569]: "portMappings": true
Aug 29 10:38:24 dev-nbsj-node2-174 kubelet[31569]: }
Aug 29 10:38:24 dev-nbsj-node2-174 kubelet[31569]: }
Aug 29 10:38:24 dev-nbsj-node2-174 kubelet[31569]: ]
Aug 29 10:38:24 dev-nbsj-node2-174 kubelet[31569]: }
Aug 29 10:38:24 dev-nbsj-node2-174 kubelet[31569]: : [failed to find plugin "flannel" in path [/opt/cni/bin]]
Aug 29 10:38:24 dev-nbsj-node2-174 kubelet[31569]: W0829 10:38:24.012395   31569 cni.go:237] Unable to update cni config: no valid networks found in /etc/cni/net.d
Aug 29 10:38:25 dev-nbsj-node2-174 kubelet[31569]: E0829 10:38:25.817847   31569 kubelet.go:2184] Container runtime network not ready: NetworkReady=false reason:NetworkPluginNotReady message:docker: network plugin is...ig uninitialized
```
### 解决方法
在新节点安装 CNI 插件时，要确保包含 flannel：
```sh
cd /tmp
curl -L https://github.jobcher.com/gh/https://github.com/flannel-io/cni-plugin/releases/download/v1.1.0/flannel-amd64 -o /opt/cni/bin/flannel
chmod +x /opt/cni/bin/flannel
```
确认文件存在
```sh
ls /opt/cni/bin/flannel
```
重启 kubelet 服务
```sh
systemctl restart kubelet
```
