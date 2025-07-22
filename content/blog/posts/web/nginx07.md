---
title: "使用TLSv1.3 升级nginx和openssl"
date: 2025-07-22
draft: false
featuredImage: "/images/nginx.png"
featuredImagePreview: "/images/nginx.png"
images: ["/images/nginx.png"]
authors: "jobcher"
tags: ["nginx", "linux"]
categories: ["web 服务器"]
series: ["web服务"]
---
## 介绍
TLS v1.3（Transport Layer Security version 1.3）是传输层安全协议的最新正式版本，用于在计算机网络中提供加密通信。它由 [IETF（Internet Engineering Task Force）](https://www.ietf.org/) 于 **2018 年 8 月正式发布**，是对 TLS v1.2 的重大改进。  
![tls_ssl_development](/images/tls_ssl_development_timeline.png)  
由于原有老的nginx版本不支持新的TLSv1.3,需要升级nginx和openssl。  

---

### 🔐 TLS 的用途

TLS 常用于以下场景：

* HTTPS（浏览器访问网站）
* 邮件客户端与服务器通信（如 IMAP/SMTP over TLS）
* VPN、聊天工具等需要安全传输的场合

---

## ✨ 相比 TLS 1.2，TLS 1.3 有哪些主要改进？

| 特性   | TLS 1.2               | TLS 1.3           |
| ---- | --------------------- | ----------------- |
| 握手轮数 | 至少 2 次往返              | 最多 1 次往返，支持 0-RTT |
| 加密套件 | 多且复杂，包含弱算法            | 简化，仅支持强加密算法       |
| 前向保密 | 可选                    | 强制启用              |
| 加密内容 | 一部分未加密                | 握手后所有内容都加密，包括证书   |
| 安全性  | 存在旧漏洞（如 BEAST、POODLE） | 移除已知不安全特性         |
| 性能   | 较慢                    | 更快（特别是在移动网络）      |

---

## 🔧 移除的内容（相比 TLS 1.2）

* RSA 密钥交换（只保留 ECDHE/DHE）
* 不安全的加密算法（如 RC4、3DES、MD5）
* 静态密钥协商、不再支持非前向保密
* 会话恢复机制被简化为基于票据（session tickets）

---

## ✅ TLS 1.3 的优势总结

* **更安全**：移除所有已知不安全或弱加密机制
* **更快**：减少握手延迟，适合移动/高延迟网络
* **更私密**：握手阶段信息也加密，防监听分析

---

## 💡 哪些应用已经支持 TLS 1.3？

* 现代浏览器（Chrome、Firefox、Safari、Edge）
* 常用 Web 服务器（Nginx、Apache、LiteSpeed）
* 后端库和操作系统（OpenSSL 1.1.1+、BoringSSL、Windows 10+）
---

## 下载编译
### 首先编译openssl
首先，下载 OpenSSL 的源代码，建议下载 1.1.1 或更高版本，因为这些版本支持 TLSv1.3。我这里是放到`/usr/local/src`  
```sh
cd /usr/local/src
sudo wget https://www.openssl.org/source/openssl-1.1.1l.tar.gz
sudo tar -xvzf openssl-1.1.1l.tar.gz
cd openssl-1.1.1l
```  
在编译 OpenSSL 时，你需要指定安装目录。通常，OpenSSL 会安装到 `/usr/local/ssl`，这样不会干扰系统的默认 OpenSSL 版本。当然你也可以选择你自己的喜欢的位置，我这里就安装在`/usr/local/ssl`  
```sh
sudo ./config --prefix=/usr/local/ssl --openssldir=/usr/local/ssl
```
### 安装和配置环境
将编译并安装 OpenSSL 到指定目录 `/usr/local/ssl`
```sh
sudo make
sudo make install
```
编辑 /etc/profile 或用户的 .bashrc 文件，添加以下内容：  
```bash
export PATH=/usr/local/ssl/bin:$PATH
export LD_LIBRARY_PATH=/usr/local/ssl/lib:$LD_LIBRARY_PATH
export C_INCLUDE_PATH=/usr/local/ssl/include:$C_INCLUDE_PATH
export CPLUS_INCLUDE_PATH=/usr/local/ssl/include:$CPLUS_INCLUDE_PATH
```
加载新的环境变量  
```sh
source /etc/profile  # 或者 source ~/.bashrc
```
检查 OpenSSL 安装
```sh
openssl version -a
```
### 编译nginx
下载和解压 Nginx 源代码  
```sh
cd /usr/local/src
sudo wget https://nginx.org/download/nginx-1.29.0.tar.gz
sudo tar -zxvf nginx-1.29.0.tar.gz
cd nginx-1.29.0
```
在编译 Nginx 时，指定 OpenSSL 的安装路径。确保在配置时使用 --with-openssl 参数，指向 OpenSSL 源代码的路径。`/usr/local/src/openssl-1.1.1l`
```sh
sudo ./configure --prefix=/usr/share/nginx --conf-path=/etc/nginx/nginx.conf \
--http-log-path=/var/log/nginx/access.log --error-log-path=/var/log/nginx/error.log \
--with-http_ssl_module --with-http_realip_module --with-http_gzip_static_module \
--with-http_image_filter_module --with-http_stub_status_module --with-pcre-jit \
--with-openssl=/usr/local/src/openssl-1.1.1l --with-debug
```
>注意完成这一步，如果要实现不停机升级nginx千万不要执行`make install`  
  
### 备份老的nginx
```sh
which nginx
cd /usr/sbin
sudo cp nginx nginx.bak
```
### 替换新的nginx
```sh
cd /usr/local/src/nginx-1.29.0
make
cp objs/nginx /usr/sbin/nginx
```
### 测试新的nginx
```sh
nginx -t
nginx -s reload
nginx -V
```
## 配置TLSv1.3
打开 Nginx 配置文件 `/etc/nginx/nginx.conf` 或你相关的站点配置文件，确保启用了 TLSv1.3。  
```conf
server {
    listen 443 ssl;
    server_name example.com;

    ssl_certificate /etc/nginx/ssl/example.com.crt;
    ssl_certificate_key /etc/nginx/ssl/example.com.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:TLS_AES_128_CCM_SHA256:TLS_AES_256_CCM_SHA384';
    ssl_prefer_server_ciphers off;

    # Other server configurations...
}
```
### 测试新配置
```sh
nginx -t
```
### 重启nginx
```sh
nginx -s reload
```
使用 openssl s_client 来验证是否启用了 TLSv1.3
```sh
openssl s_client -connect example.com:443 -tls1_3
```
或者浏览器访问  
![chrome](/images/nginx07.png)  
  
## 总结
- 编译 OpenSSL 1.1.1 并将其安装到 /usr/local/ssl。
- 重新编译 Nginx，确保它链接到新安装的 OpenSSL。
- 在 Nginx 配置文件 中启用 TLSv1.3 支持。
- 验证配置 并重新加载 Nginx。
