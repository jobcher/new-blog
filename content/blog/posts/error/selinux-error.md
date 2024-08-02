---
title: "SELinux 问题：导致端口无法创建，无法访问"
date: 2024-08-02
draft: false
authors: "jobcher"
featuredImage: '/images/SELinux.jpg'
featuredImagePreview: '/images/SELinux.jpg'
images: ['/images/SELinux.jpg']
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---
## 背景
今天有同事在使用nginx部署一个服务，部署完成后发现无法访问，nginx创建端口无法创建，无法访问  
![nginx-error](/images/selinux-error-1.png)  
```sh
nginx: [emerg] bind() to 0.0.0.0:8081 failed (13: Permission denied)
```

## 解决方法
查看日志发现是SELinux导致的，SELinux是Linux系统的安全机制，它会限制进程访问文件和网络端口等资源。  
查看SELinux状态
```sh
sudo getenforce
```
当 SELinux 处于 enforcing 模式时，会阻止进程访问不允许的资源。有三种方法可以解决

### 1. 临时关闭SELinux
```sh
sudo setenforce 0
```
### 2. 永久关闭SELinux
```sh
sudo vim /etc/selinux/config
# 修改SELINUX=enforcing 为 SELINUX=disabled
```
重启服务器
```sh
sudo reboot
```
### 3. 设置为宽容模式
```sh
semanage permissive -a http_port_t
```
这个命令会将 http_port_t 类型的端口设置为宽容模式（permissive mode），使得 semanage 不再对该类型的端口进行访问控制。  
## 总结
 SELinux 是 Linux 系统的安全机制，它会限制进程访问文件和网络端口等资源。在使用 SELinux 时，需要根据实际情况选择合适的解决方案。

