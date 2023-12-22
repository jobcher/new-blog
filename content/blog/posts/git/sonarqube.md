---
title: "sonarqube docker安装和配置"
date: 2023-11-02
draft: false
featuredImage: "/images/sonarqube.png"
featuredImagePreview: "/images/sonarqube.png"
images: ["/images/sonarqube.png"]
author: "jobcher"
tags: ["gitlab"]
categories: ["gitlab"]
series: ["git入门系列"]
---
## 背景
我是在ubuntu服务器安装docker服务，我已经安装好了docker和docker-compose服务，这里我就不写这些服务的安装过程，直接开始安装sonarqube服务
## 安装 sonarqube服务器
### 1.执行脚本文件 `config.sh`
```bash
#!/bin/bash
sysctl -w vm.max_map_count=262144
sysctl -w fs.file-max=65536
ulimit -n 65536
ulimit -u 4096
# # 永久改变
# echo "vm.max_map_count=262144" >> /etc/sysctl.conf
# sysctl -p
```
```sh
sh config.sh
```

### 2.执行docker-compose文件
```yaml
version: '3'
services:
  postgres:
    image: postgres:15
    container_name: postgres
    ports:
      - 5432:5432
    volumes:
      - ./sonar/postgres/postgresql:/var/lib/postgresql
      - ./sonar/postgres/data:/var/lib/postgresql/data
    environment:
      TZ: Asia/Shanghai
      POSTGRES_USER: user #数据库用户
      POSTGRES_PASSWORD: password #数据库密码
      POSTGRES_DB: sonar

  sonarqube:
    depends_on:
      - postgres
    image: sonarqube:9.9.0-community
    container_name: sonarqube
    ports:
      - 9000:9000
    volumes:
      - ./sonar/sonarqube/extensions:/opt/sonarqube/extensions
      - ./sonar/sonarqube/logs:/opt/sonarqube/logs
      - ./sonar/sonarqube/data:/opt/sonarqube/data
      - ./sonar/sonarqube/conf:/opt/sonarqube/conf
    environment:
      SONARQUBE_JDBC_USERNAME: user #数据库用户
      SONARQUBE_JDBC_PASSWORD: password #数据库密码
      SONARQUBE_JDBC_URL: jdbc:postgresql://postgres:5432/sonar
```
```sh
docker-compose up -d
```
### 3. 生成用户令牌
![1698896646032.jpg](/images/1698896646032.jpg)  

## 安装 sonar scanner扫描器
```sh
wget https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-5.0.1.3006-linux.zip
unzip sonar-scanner-cli-5.0.1.3006-linux.zip
mv sonar-scanner-cli-5.0.1.3006-linux.zip /opt/
```
连接sonarqube服务器
```sh
cd /opt/sonar-scanner-5.0.1.3006-linux/conf/
vim sonar-scanner.properties
```
修改配置文件
```sh
#No information about specific project should appear here

#----- 修改服务器地址
sonar.host.url=http://sonar-server:9000

#----- Default source code encoding
sonar.sourceEncoding=UTF-8
```
安装scanner
```sh
cd sonar-scanner-5.0.1.3006-linux/
sudo ln -s /opt/sonar-scanner-5.0.1.3006-linux/bin/sonar-scanner /usr/local/bin/sonar-scanner
sonar-scanner -v
```

## 执行扫描
```sh
cd /code/github/project/path
sonar-scanner -D sonar.login=<sonar-key> -Dsonar.projectKey=<my-project> -Dsonar.projectName=<my-project-name> -Dsonar.projectVersion=1.0 -Dsonar.sourceEncoding=UTF-8 -Dsonar.java.binaries=target/classes -Dsonar.java.libraries=target/*.jar
```

## 在sonarsever上查看
![1698896917631.jpg](/images/1698896917631.jpg)