---
title: "Kubernetes — RKE2 + kube-vip + cilium 部署"
date: 2025-09-18
draft: false
featuredImage: "/images/kubernetes0install.png"
featuredImagePreview: "/images/kubernetes0install.png"
images: ["/images/kubernetes0install.png"]
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s入门系列"]
---
## 准备工作
|节点名称|节点IP|
|:----|:----|
|k8s-master-1|10.10.10.151|
|k8s-master-2|10.10.10.152|
|k8s-master-3|10.10.10.153|
|kube-vip(虚拟IP)|10.10.10.150|

### RKE 安装 rancher
#### 在第一个 master 安装 RKE2 server
```sh
# 安装 RKE2
curl -sfL https://get.rke2.io | sh -
```
创建配置文件
```sh
mkdir -p /etc/rancher/rke2/
```
```sh
# 配置 server
cat <<EOF >/etc/rancher/rke2/config.yaml
write-kubeconfig-mode: "0644"
tls-san:
  - 10.10.10.150
  - rancher.jobcher.com
cni: cilium
disable-kube-proxy: true
EOF
```
```sh
# 启动 server
systemctl enable rke2-server --now
systemctl status rke2-server
```
```sh
ln -s /var/lib/rancher/rke2/bin/kubectl /usr/local/bin/kubectl
echo 'export KUBECONFIG=/etc/rancher/rke2/rke2.yaml' >> ~/.bashrc
source ~/.bashrc
```
#### kubectl 补全
```sh
apt-get install -y bash-completion
```
```sh
echo 'source <(kubectl completion bash)' >> ~/.bashrc
source ~/.bashrc
kubectl completion bash >/etc/bash_completion.d/kubectl
source /etc/bash_completion.d/kubectl
```
#### cilium 安装
下载 Cilium CLI

```sh
curl -L --remote-name https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz
tar xzvf cilium-linux-amd64.tar.gz
mv cilium /usr/local/bin/
cilium version

```
```sh
cilium status
kubectl get pods -n kube-system -l k8s-app=cilium
kubectl get nodes -o wide
```
```sh
vim rke2-cilium-config.yml
```
```yaml
# /var/lib/rancher/rke2/server/manifests/rke2-cilium-config.yml
apiVersion: helm.cattle.io/v1
kind: HelmChartConfig
metadata:
  name: rke2-cilium
  namespace: kube-system
spec:
  valuesContent: |-
    tunnelProtocol: geneve
    kubeProxyReplacement: true
    k8sServiceHost: localhost
    k8sServicePort: 6443
    hubble:
      enabled: true
      relay:
        enabled: true
      ui:
        enabled: true
```
配置 hubble-ui ，查看网络结构
```sh
kubectl apply -f rke2-cilium-config.yml

kubectl -n kube-system patch svc hubble-ui \
  -p '{"spec": {"type": "NodePort"}}'

```

#### 其他master 加入集群
获取token值
```sh
cat /var/lib/rancher/rke2/server/node-token
```
```sh
curl -sfL https://get.rke2.io | sh -
```
```sh
mkdir -p /etc/rancher/rke2/
```
```sh
cat <<EOF >/etc/rancher/rke2/config.yaml
server: https://10.10.10.150:9345
token: <token> # 输入token值
tls-san:
  - 10.10.10.150
  - rancher.jobcher.com
cni: cilium
disable-kube-proxy: true
EOF
```
```sh
# 启动 server
systemctl enable rke2-server --now
systemctl status rke2-server
```


#### 配置kubectl
```sh
ln -s /var/lib/rancher/rke2/bin/kubectl /usr/local/bin/kubectl
echo 'export KUBECONFIG=/etc/rancher/rke2/rke2.yaml' >> ~/.bashrc
source ~/.bashrc
```

```sh
kubectl get pod -n kube-system | grep kube-vip
```

#### kube-vip 安装
```sh
KVVERSION=$(curl -sL https://api.github.com/repos/kube-vip/kube-vip/releases | jq -r ".[0].name")
```
```sh
alias ctr="/var/lib/rancher/rke2/bin/ctr --address /run/k3s/containerd/containerd.sock"
```
```sh
alias kube-vip="ctr image pull ghcr.io/kube-vip/kube-vip:$KVVERSION; ctr run --rm --net-host ghcr.io/kube-vip/kube-vip:$KVVERSION vip /kube-vip"
```
```sh
wget https://kube-vip.io/manifests/rbac.yaml
mv rbac.yaml kube-vip-rbac.yaml
chmod +x kube-vip-rbac.yaml && kubectl apply -f kube-vip-rbac.yaml
```
```sh
# 运行
kube-vip manifest daemonset \
  --arp \
  --controlplane \
  --address 10.10.10.150\
  --interface eth0 \
  --leaderElection \
  --enableLoadBalancer \
  --inCluster \
  --taint > kube-vip.yaml
```
```sh
kubectl apply -f kube-vip.yaml
```
检测vip
```sh
ping 10.10.10.150
```
```sh
curl -k https://10.10.10.150:9345/v1-rke2/connect
```
#### rancher 安装
```sh
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

helm repo add rancher-stable https://releases.rancher.com/server-charts/stable
helm repo update

kubectl create namespace cattle-system
kubectl -n cattle-system create secret tls tls-rancher-ingress --cert=fullchain.pem --key=privkey.pem
```
```sh
helm upgrade --install rancher rancher-stable/rancher \
  --namespace cattle-system \
  --set hostname="rancher.jobcher.com" \
  --set ingress.tls.source="secret" \
  --set bootstrapPassword="输入你的密码"
```