---
title: "nginx exporter 安装配置"
date: 2022-06-08
draft: false
authors: "jobcher"
tags: ["prometheus"]
categories: ["监控"]
series: ["运维监控系列"]
---

# nginx exporter 安装配置

1. 二进制安装

```sh
wget https://github.com/nginxinc/nginx-prometheus-exporter/releases/download/v0.10.0/nginx-prometheus-exporter_0.10.0_linux_amd64.tar.gz
tar -zxvf nginx-prometheus-exporter_0.10.0_linux_amd64.tar.gz -C ./nginx-exporter
```

2. 在 nginx 上配置

```sh
./configure \
… \
--with-http_stub_status_module
make
sudo make install
```

在 nginx.config 上配置

```conf
server {
    # 新增
    location /nginx_status {
        stub_status on;
        access_log off;
    }
}
```

重启 nginx 服务

```sh
nginx -t
nginx -s reload
```

3. 启动 nginx exporter

```sh
nginx-prometheus-exporter -nginx.scrape-uri http://<nginx>:8080/nginx_status
```

4. 配置 prometheus
   添加 prometheus.yml

```yml
- job_name: "nginx-exporter"
  file_sd_configs:
    - files:
        - "./file_sd/nginx-exporter.yaml"
```

在 ./file_sd/新建 nginx-exporter.yaml

```yml
- targets: ["<IP>:9113"]
  labels:
    instance: <nginx名称>
```
