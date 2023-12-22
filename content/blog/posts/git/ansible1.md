---
title: "ansible 命令"
date: 2022-04-14
draft: false
author: "jobcher"
tags: ["ansible"]
categories: ["gitlab"]
series: ["git入门系列"]
---

# ansible 命令

- `Inventory`：Ansible 管理的主机信息，包括 IP 地址、SSH 端口、账号、密码等
- `Modules`：任务均有模块完成，也可以自定义模块，例如经常用的脚本。
- `Plugins`：使用插件增加 Ansible 核心功能，自身提供了很多插件，也可以自定义插件。例如 connection 插件，用于连接目标主机。
- `Playbooks`：“剧本”，模块化定义一系列任务，供外部统一调用。Ansible 核心功能。

## 编辑主机清单

```sh
[webservers]
192.168.0.20 ansible_ssh_user=root ansible_ssh_pass=’200271200’
192.168.0.21 ansible_ssh_user=root ansible_ssh_pass=’200271200’
192.168.0.22 ansible_ssh_user=root ansible_ssh_pass=’200271200’

[dbservers]
10.12.0.100
10.12.0.101
```

```sh
sed -i "s/#host_key_checking = .*/host_key_checking = False/g" /etc/ansible/ansible.cfg
```

## 命令行

```sh
ansible all -m ping
ansible all -m shell -a "ls /root" -u root -k
```

## 常用模块

在目标主机执行 shell 命令。

1. shell

```yml
- name: 将命令结果输出到指定文件
  shell: somescript.sh >> somelog.txt
- name: 切换目录执行命令
  shell:
    cmd: ls -l | grep log
    chdir: somedir/
- name: 编写脚本
  shell: |
    if [ 0 -eq 0 ]; then
       echo yes > /tmp/result
    else
       echo no > /tmp/result
    fi
  args:
    executable: /bin/bash
```

2. copy
   将文件复制到远程主机。

```yml
- name: 拷贝文件
  copy:
    src: /srv/myfiles/foo.conf
    dest: /etc/foo.conf
    owner: foo
    group: foo
    mode: u=rw,g=r,o=r
    # mode: u+rw,g-wx,o-rwx
    # mode: '0644'
    backup: yes
```

3. file
   管理文件和文件属性。

```yml
- name: 创建目录
  file:
    path: /etc/some_directory
    state: directory
    mode: "0755"
- name: 删除文件
  file:
    path: /etc/foo.txt
    state: absent
- name: 递归删除目录
  file:
    path: /etc/foo
    state: absent
```

- present，latest：表示安装
- absent：表示卸载

4. yum
   软件包管理。

```yml
- name: 安装最新版apache
  yum:
    name: httpd
    state: latest
- name: 安装列表中所有包
  yum:
    name:
      - nginx
      - postgresql
      - postgresql-server
    state: present
- name: 卸载apache包
  yum:
    name: httpd
    state: absent
- name: 更新所有包
  yum:
    name: "*"
    state: latest
- name: 安装nginx来自远程repo
  yum:
    name: http://nginx.org/packages/centos/6/noarch/RPMS/nginx-release-centos-6-0.el6.ngx.noarch.rpm
    # name: /usr/local/src/nginx-release-centos-6-0.el6.ngx.noarch.rpm
    state: present
```

5. service/systemd
   管理服务

```yml
- name: 服务管理
  service:
    name: httpd
    state: started
    #state: stopped
    #state: restarted
    #state: reloaded
- name: 设置开机启动
  service:
    name: httpd
    enabled: yes
```

6. unarchive
   解压

```yml
- name: 解压
  unarchive: src=test.tar.gz
    dest=/tmp
```

7. debug
   执行过程中打印语句。

```yml
- debug:
    msg: System {{ inventory_hostname }} has uuid {{ ansible_product_uuid }}

- name: 显示主机已知的所有变量
  debug:
    var: hostvars[inventory_hostname]
    verbosity: 4
```
