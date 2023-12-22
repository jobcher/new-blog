---
title: "安装配置 Terraform"
date: 2021-12-27
draft: false
author: "jobcher"
tags: ["Terraform"]
categories: ["日常"]
series: ["日常系列"]
---

# 安装配置 Terraform

## 安装

1. macOS 苹果系统安装

```sh
#安装
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
# 更新
brew update
brew upgrade hashicorp/tap/terraform
#验证安装
terraform -help
```

2. windows 系统安装

```sh
#安装
choco install terraform
#直接到这个url里下载64位系统
https://www.terraform.io/downloads
#验证安装
terraform -help
```

3. Linux 安装

```sh
curl -fsSL https://apt.releases.hashicorp.com/gpg | sudo apt-key add -
sudo apt-add-repository "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"
sudo apt-get update && sudo apt-get install terraform
#验证安装
terraform -help
```

```sh
wget -O- https://rpm.releases.hashicorp.com/RHEL/hashicorp.repo | sudo tee /etc/yum.repos.d/hashicorp.repo
sudo yum install terraform -y
```

## terrafrom 控制 proxmox 虚拟机

来源：https://github.com/Telmate/terraform-provider-proxmox

### 首先你要有一台 pve 主机

安装过程本篇文章就不想了，主要是要写一下关于他的配置  
https://pve.proxmox.com/pve-docs/

1. 下载

```sh
wget https://github.com/Telmate/terraform-provider-proxmox/releases/download/v2.9.3/terraform-provider-proxmox_2.9.3_linux_amd64.zip
unzip terraform-provider-proxmox_2.9.3_linux_amd64.zip
```

### 编写 terrafrom 程序

1. 虚拟机 main.tf

```sh
terraform {
  required_version = ">= 0.14"
  required_providers {
    proxmox = {
      source  = "telmate/proxmox"
    }
  }
}

provider "proxmox" {
  # 配置选项
    pm_tls_insecure = true
    pm_api_url = "https://localhost:8006/api2/json"
    pm_user = "root@pam"
    pm_password = "passwd"
    pm_otp = ""
}

# 创建VM
resource "proxmox_vm_qemu" "cloudinit-test" {
    name = "terraform-test-vm"
    desc = "A test for using terraform and cloudinit"

    #节点名称必须与集群内的名称相同
    #这可能不包括 FQDN
    target_node = "pve"

    #新虚拟机的目标资源池
    pool = "pool0"

    #从中克隆这个虚拟机的模板名称
    clone = "node0"

    #为这个虚拟机激活 QEMU 代理
    agent = 1

    os_type = "cloud-init"
    cores = 2
    sockets = 1
    vcpus = 0
    cpu = "host"
    memory = 2048
    scsihw = "lsi"

    #设置磁盘
    disk {
        size = 32
        type = "virtio"
        storage = "local-lvm"
        storage_type = "lvmthin"
        iothread = 1
        ssd = 1
        discard = "on"
    }

    #设置网络接口并分配一个 vlan 标签：256
    network {
        model = "virtio"
        bridge = "vmbr0"
        tag = 256
    }
}

```

2. 运行 terrafrom

```sh
# 初始化
terraform init
# 查看产生的变更
terraform plan
# 运行
terraform apply
```

## 配置

配置这个 terraform 我们这个需要持续更新，首先我们先配置 Azure 吧

### Azure 配置

1. 安装 azurecli

```sh
# linux
curl -L https://aka.ms/InstallAzureCli | bash
apt install azure-cli
# macOS
brew update && brew install azure-cli
```

2. 登录 azure

```sh
# 中国区azure
az cloud set --name AzureCloud
az login -u <账户> -p <密码>

#海外azure
az cloud set --name AzureChinaCloud
az login -u <账户> -p <密码>
```

3. 创建 terrafrom 代码
   创建 main.tf

```tf
# 正在使用的 Azure 提供程序源和版本
terraform {
    required_version = ">=0.12"

    required_providers {
        azurerm = {
            source = "hashicorp/azurerm"
            version = "~>2.0"
        }
    }
}
# 配置 Microsoft Azure 提供程序
provider "azurerm" {
  features {}
}

# 资源组前缀
resource "random_pet" "rg-name" {
  prefix    = var.resource_group_name_prefix
}

# 创建资源组
resource "azurerm_resource_group" "rg" {
  name      = random_pet.rg-name.id
  location  = var.resource_group_location
}
```

创建 variable.tf

```tf
variable "resource_group_name_prefix" {
  default       = "rg"
  description   = "Prefix of the resource group name that's combined with a random ID so name is unique in your Azure subscription."
}

variable "resource_group_location" {
  default = "eastus"
  description   = "Location of the resource group."
}
```
