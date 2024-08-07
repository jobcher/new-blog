---
title: "GNU/Linux 一键更换系统软件源脚本"
date: 2023-09-21
draft: false
featuredImage: "/images/gnulinux-logo-gnu-linux-logo-png.png"
featuredImagePreview: "/images/gnulinux-logo-gnu-linux-logo-png.png"
images: ["/images/gnulinux-logo-gnu-linux-logo-png.png"]
authors: "jobcher"
tags: ["运维"]
categories: ["基础"]
series: ["基础知识系列"]
---

## 背景
有很多小伙伴，在接受公司老项目后，想要更新服务器源时，发现镜像源失效了，手动添加也很不方便，我这里提供了一个脚本供大家使用，轻松切换镜像源。  
支持系统：  
<table align="center">
<tr>
    <td><a href="https://www.debian.org" target="_blank"><img src="/images/icon/debian.svg" width="16" height="16" style="vertical-align: -0.45em"></a>&nbsp;Debian</td>
    <td align="center">8.0 ~ 12</td>
</tr>
<tr>
    <td><a href="https://cn.ubuntu.com" target="_blank"><img src="/images/icon/ubuntu.svg" width="16" height="16" style="vertical-align: -0.15em"></a>&nbsp;Ubuntu</td>
    <td align="center">14.04 ~ 23</td>
</tr>
<tr>
    <td><a href="https://www.kali.org" target="_blank"><img src="/images/icon/kali-linux.svg" width="16" height="16" style="vertical-align: -0.15em"></a>&nbsp;Kali Linux</td>
    <td align="center">2.0 ~ 2023</td>
</tr>
<tr>
    <td><a href="https://access.redhat.com/products/red-hat-enterprise-linux" target="_blank"><img src="/images/icon/redhat.svg" width="16" height="16" style="vertical-align: -0.15em"></a>&nbsp;Red Hat Enterprise Linux</td>
    <td align="center">7.0 ~ 9</td>
</tr>
<tr>
    <td><a href="https://fedoraproject.org/zh-Hans" target="_blank"><img src="/images/icon/fedora.ico" width="16" height="16" style="vertical-align: -0.15em"></a>&nbsp;Fedora</td>
    <td align="center">30 ~ 38</td>
</tr>
<tr>
    <td><a href="https://www.centos.org" target="_blank"><img src="/images/icon/centos.svg" width="16" height="16" style="vertical-align: -0.15em"></a>&nbsp;CentOS</td>
    <td align="center">7.0 ~ 8.5 / Stream 8 ~ 9</td>
</tr>
<tr>
    <td><a href="https://rockylinux.org/zh_CN" target="_blank"><img src="/images/icon/rocky-linux.svg" width="16" height="16" style="vertical-align: -0.25em"></a>&nbsp;Rocky Linux</td>
    <td align="center">8 ~ 9</td>
</tr>
<tr>
    <td><a href="https://almalinux.org/zh-hans" target="_blank"><img src="/images/icon/almalinux.svg" width="16" height="16" style="vertical-align: -0.25em"></a>&nbsp;AlmaLinux</td>
    <td align="center">8 ~ 9</td>
</tr>
<tr>
    <td><a href="https://www.opencloudos.org" target="_blank"><img src="/images/icon/opencloudos.png" width="16" height="16" style="vertical-align: -0.25em"></a>&nbsp;OpenCloudOS</td>
    <td align="center">8.6 / 9.0</td>
</tr>
<tr>
    <td><a href="https://www.openeuler.org/zh" target="_blank"><img src="/images/icon/openeuler.ico" width="16" height="16" style="vertical-align: -0.15em"></a>&nbsp;openEuler</td>
    <td align="center">21.03 ~ 23</td>
</tr>
<tr>
    <td><a href="https://www.opensuse.org" target="_blank"><img src="/images/icon/opensuse.svg" width="16" height="16" style="vertical-align: -0.15em"></a>&nbsp;openSUSE</td>
    <td align="center">Leep 15 / Tumbleweed</td>
</tr>
<tr>
    <td><a href="https://archlinux.org" target="_blank"><img src="/images/icon/arch-linux.ico" width="16" height="16" style="vertical-align: -0.15em"></a>&nbsp;Arch Linux</td>
    <td align="center">all</td>
</tr>
</table>

## 执行命令
### 国内使用
```bash
bash <(curl -sSL https://www.jobcher.com/ChangeMirrors.sh)
```
### 教育网使用
```bash
bash <(curl -sSL https://www.jobcher.com/ChangeMirrors.sh) --edu
```
### 海外使用
```bash
bash <(curl -sSL https://www.jobcher.com/ChangeMirrors.sh) --abroad
```

## 注意事项
1. 需使用 `Root` 用户执行脚本
    切换命令为 `sudo -i` 或 `su root`，不同系统环境使用的命令不一样，因为有些系统没有在初始安装时为 `Root` 用户设置固定密码所以需要使用 `sudo` 指令来提权
2. 建议使用 `SSH` 远程工具
    如果你使用的系统终端界面不支持 UTF-8 编码那么将无法正常显示中文内容，导致无法正确选择交互内容。大部分系统都会自动开启该服务
3. 首次在新系统上执行脚本
    当前执行方式依赖 `curl` 指令来获取脚本内容并执行，所以需要先通过包管理工具来安装该软件包，否则会报错 `Command not found`，若无法安装就复制源码到本地新建.sh脚本，然后通过 bash 手动执行

## 其他问题
有一部分小伙伴会出现这些问题：  
>gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-EPEL-7  
  
这个配置告诉YUM，这个repository的GPG key存在于磁盘上。而当YUM在路径 /etc/pki/rpm-gpg 下找不到这个GPG key的时候，就会报如上的错误了  
### 解决
```bash
cd /etc/pki/rpm-gpg
yum install wget -y
wget https://archive.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-7
```