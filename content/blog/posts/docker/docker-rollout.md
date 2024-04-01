---
title: "docker的零停机部署"
date: 2024-04-01
draft: false
authors: "jobcher"
featuredImage: "/images/docker.png"
featuredImagePreview: "/images/docker.png"
images: ['/images/docker.png']
tags: ["docker"]
categories: ["docker"]
series: ["docker进阶系列"]
---
## 背景
使用 `docker compose up` 部署新版本的服务会导致停机，因为应用容器在创建新容器之前已停止。如果你的应用程序需要一段时间才能启动，用户可能会注意到这一点。为了保障服务用户无感，可以使用`docker rollout`  
  
>适合没必要用 K8S 轻量级小项目  
  
## 安装
[项目地址](https://github.com/Wowu/docker-rollout)  
  
```sh
# 为 Docker cli 插件创建目录
mkdir -p ~/.docker/cli-plugins

# 下载 docker-rollout 脚本到 Docker cli 插件目录
curl https://github.jobcher.com/gh/https://raw.githubusercontent.com/wowu/docker-rollout/master/docker-rollout -o ~/.docker/cli-plugins/docker-rollout

# 使脚本可执行
chmod +x ~/.docker/cli-plugins/docker-rollout
```
## 使用
### 注意事项！！！
- 服务不能在 `docker-compose.yml` 中定义 `container_name` 和 `ports` ，因为不可能运行具有相同名称或端口映射的多个容器。
- 需要像 `Traefik` 或 `nginx-proxy` 这样的代理来路由流量。
- 每次部署都会`增加`容器名称中的索引（例如 `project-web-1` -> `project-web-2` ）
### 使用示范
```sh
# 下载代码
git pull
# 构建新的应用程序映像
docker compose build web
# 运行数据库迁移
docker compose run web rake db:migrate
# 部署新版本
docker rollout web
```
或者使用`docker-compose.yaml`
```sh
docker rollout -f docker-compose.yml <service-name>
```
### 参数
- -f | --file FILE - （非必需）- 撰写文件的路径，可以多次指定，如 docker compose 中。
- -t | --timeout SECONDS -（非必需）- 如果容器在 Dockerfile 或 docker-compose.yml 中定义了运行状况检查，则等待新容器变得健康的超时时间（以秒为单位）。默认值：60
- -w | --wait SECONDS - （非必需）- 如果未定义 healthcheck，则等待新容器准备就绪的时间。默认值：10
- --env-file FILE - （非必需）- env 文件的路径，可以多次指定，如 docker compose 中。