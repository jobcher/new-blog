---
title: "NVIDIA GPU 的 docker 安装和部署"
date: 2023-06-01
draft: false
featuredImage: "/images/nvidia-docker.png"
featuredImagePreview: "/images/nvidia-docker.png"
images: ['/images/nvidia-docker.png']
author: "jobcher"
tags: ["ChatGPT"]
categories: ["人工智能"]
series: ["人工智能系列"]
---
## 背景
在人工智能领域，GPU 是必不可少的。在本文中，我们将介绍如何在服务器上安装和部署 NVIDIA GPU 的 docker

### 升级你的CUDA
[官网链接](https://developer.nvidia.com/cuda-toolkit-archive)  
  
选择你的系统对应版本进行安装  
![cuda](/images/CUDA-download.png)  


## 安装
1. 确保你的系统已经安装了NVIDIA驱动和Docker引擎。确保驱动版本与Docker引擎兼容。你还需要安装nvidia-docker2软件包，它是NVIDIA Docker的一个插件。可以在[https://github.com/NVIDIA/nvidia-docker](https://github.com/NVIDIA/nvidia-docker)上找到安装说明  
1.1 卸载旧版本的nvidia-docker(如果已安装)
```bash
sudo apt-get remove nvidia-docker
```
1.2 添加nvidia-docker2仓库
```bash
distribution=$(. /etc/os-release;echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list
```
1.3 更新软件包列表和安装nvidia-docker2
```bash
sudo apt-get update
sudo apt-get install nvidia-docker2
```
1.4 重启Docker守护进程
```bash
sudo systemctl restart docker
```
2. 通过以下命令测试是否安装成功
```bash
docker run --gpus all nvidia/cuda:10.0-base nvidia-smi
```
3. 如果安装成功，你将看到类似下面的输出
```bash
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 440.33.01    Driver Version: 440.33.01    CUDA Version: 10.2     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  Tesla P100-PCIE...  On   | 00000000:00:1E.0 Off |                    0 |
| N/A   32C    P0    32W / 250W |      0MiB / 16280MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   1  Tesla P100-PCIE...  On   | 00000000:00:1F.0 Off |                    0 |
| N/A   32C    P0    32W / 250W |      0MiB / 16280MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   2  Tesla P100-PCIE...  On   | 00000000:00:20.0 Off |                    0 |
| N/A   32C    P0    32W / 250W |      0MiB / 16280MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   3  Tesla P100-PCIE...  On   | 00000000:00:21.0 Off |                    0 |
| N/A   32C    P0    32W / 250W |      0MiB / 16280MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   4  Tesla P100-PCIE...  On   | 00000000:00:22.0 Off |                    0 |
| N/A   32C    P0    32W / 250W |      0MiB / 16280MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   5  Tesla P100-PCIE...  On   | 00000000:00:23.0 Off |                    0 |
| N/A   32C    P0    32W / 250W |      0MiB / 16280MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   6  Tesla P100-PCIE...  On   | 00000000:00:24.0 Off |                    0 |
| N/A   32C    P0    32W / 250W |      0MiB / 16280MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
|   7  Tesla P100-PCIE...  On   | 00000000:00:25.0 Off |                    0 |
| N/A   32C    P0    32W / 250W |      0MiB / 16280MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+
```
## 部署
1. 创建一个新的dockerfile
```Dockerfile
FROM nvidia/cuda:11.7.0-cudnn8-runtime-ubuntu18.04

# 安装软件包
RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    wget \
    vim \
    openssh-server 

RUN mkdir -p /run/sshd
RUN echo "root:password" | chpasswd
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]

```
2. 创建一个新的docker-compose.yml
```yml
version: "3"
services:
  cuda-container:
    build: .
    runtime: nvidia
    container_name: cuda-container
    ports:
      - 2201:22
    volumes:
      - ./app:/app
```
3. 启动docker
```bash
docker-compose up -d
```
4. 验证
```bash
docker exec -ti cuda-container nvidia-smi
```