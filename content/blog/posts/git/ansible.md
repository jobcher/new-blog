---
title: "ansible 安装和部署"
date: 2021-12-30
draft: false
authors: "jobcher"
tags: ["ansible"]
categories: ["gitlab"]
series: ["git入门系列"]
---

# ansible 安装和部署

Ansible 默认通过 SSH 协议管理机器.

## 安装 ansible

1. 下载安装

```sh
# ubuntu 安装
apt-get install software-properties-common
apt-add-repository ppa:ansible/ansible
apt-get update
apt-get install ansible
# centos 安装
yum install ansible
```

2. 检查文件

```sh
#检查
ansible --version
```

## ansible 配置

1. 添加主机

```sh
vim /etc/ansible/hosts
#添加你需要添加的被控主机地址和IP
```

2. 配置 SSH key 授权访问

```sh
# 控制主机生成ssh 密钥对（一路回车）
ssh-keygen -t rsa
# 复制公钥IP到被控主机
ssh-copy-id -i ~/.ssh/id_rsa.pub root@192.168.0.2
ssh-copy-id -i ~/.ssh/id_rsa.pub root@192.168.0.3
ssh-copy-id -i ~/.ssh/id_rsa.pub root@192.168.0.4
# ssh-copy-id命令会自动将id_rsa.pub文件的内容追加到远程主机root用户下.ssh/authorized_keys文件中。
```

3. 更改 ansible 配置

```sh
vim /etc/ansible/ansible.cfg

#禁用每次执行ansbile命令检查ssh key host
host_key_checking = False
# 开启日志记录
log_path = /var/log/ansible.log
```

4. 测试

```sh
# 控制主机
ansible all -m ping
```

## Inventory

1. 配置组别

```ini
vim /etc/ansible/hosts
# 添加组别

[pve]
192.168.0.2
192.168.0.3
192.168.0.4

#测试
ansible pve -m ping
```

2. Inventory 参数
   把你的 inventory 文件 和 变量 放入 git repo 中,以便跟踪他们的更新,这是一种非常推荐的方式.

| 参数                                  | 作用                                                                                                                                                                            |
| :------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| ansible_ssh_host                      | 将要连接的远程主机名.与你想要设定的主机的别名不同的话,可通过此变量设置.                                                                                                         |
| ansible_ssh_port                      | ssh 端口号.如果不是默认的端口号,通过此变量设置.                                                                                                                                 |
| ansible_ssh_user                      | 默认的 ssh 用户名                                                                                                                                                               |
| ansible_ssh_pass                      | ssh 密码(这种方式并不安全,我们强烈建议使用 --ask-pass 或 SSH 密钥)                                                                                                              |
| ansible_sudo_pass                     | sudo 密码(这种方式并不安全,我们强烈建议使用 --ask-sudo-pass)                                                                                                                    |
| ansible_sudo_exe (new in version 1.8) | sudo 命令路径(适用于 1.8 及以上版本)                                                                                                                                            |
| ansible_connection                    | 与主机的连接类型.比如:local, ssh 或者 paramiko. Ansible 1.2 以前默认使用 paramiko.1.2 以后默认使用 'smart','smart' 方式会根据是否支持 ControlPersist, 来判断'ssh' 方式是否可行. |
| ansible_ssh_private_key_file          | ssh 使用的私钥文件.适用于有多个密钥,而你不想使用 SSH 代理的情况.                                                                                                                |
| ansible_shell_type                    | 目标系统的 shell 类型.默认情况下,命令的执行使用 'sh' 语法,可设置为 'csh' 或 'fish'.                                                                                             |
| ansible_python_interpreter            | 目标主机的 python 路径.适用于的情况: 系统中有多个 Python, 或者命令路径不是"/usr/bin/python",比如 \*BSD, 或者 /usr/bin/python                                                    |

## Roles

Roles 是基于已知文件结构自动加载某些变量文件，任务和处理程序的方法。按角色对内容进行分组，适合构建复杂的部署环境
Roles 目录结构：

```yaml
site.yml
webservers.yml
fooservers.yml
roles/
common/
tasks/
handlers/
files/
templates/
vars/
defaults/
meta/
webservers/
tasks/
defaults/
meta/
```

- tasks 包含角色要执行的任务的主要列表
- handlers 包含处理程序，此角色甚至在此角色之外的任何地方都可以使用这些处理程序
- defaults 角色的默认变量
- vars 角色的其他变量
- files 包含可以通过此角色部署的文件
- templates 包含可以通过此角色部署的模版
- meta 为此角色定义一些元数据
