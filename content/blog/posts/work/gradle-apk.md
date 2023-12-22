---
title: "Jenkins 编译Android apk 流水线"
date: 2022-11-25
draft: false
featuredImage: "/images/logo-jenkins.png"
featuredImagePreview: "/images/logo-jenkins.png"
author: "jobcher"
tags: ["work"]
categories: ["work"]
series: ["问题库系列"]
---

## 背景

Jenkins 编译 Android apk，上传 apk 包，生成下载二维码，并推送钉钉

## 安装 Android 环境

### 安装 JDK

```sh
# 这里使用的是openjdk 1.8.0版本，有需要的话需要到java官网上进行下载对应的JDK版本。
$ yum install java -y

# 其他版本JDK的安装方式
$ mv jdk1.8.0_161 /usr/local/
$ ln -s /usr/local/jdk1.8.0_161 /usr/local/jdk
$ vim /etc/profile     #配置JDK的环境变量
export JAVA_HOME=/usr/local/jdk
export PATH=$JAVA_HOME/bin:$JAVA_HOME/jre/bin:$PATH
export CLASSPATH=.$CLASSPATH:$JAVA_HOME/lib:$JAVA_HOME/jre/lib:$JAVA_HOME/lib/tools.jar
$ source /etc/profile    #重新加载系统环境变量
$ java -version    #查看java版本

```

### Android SDK 安装

```sh
# 下载sdk工具包
$ wget https://dl.google.com/android/repository/sdk-tools-linux-3859397.zip

# 创建sdk工具文件夹和解压工具包
$ mkdir  -p /opt/android/sdk
$ unzip sdk-tools-linux-3859397.zip -d /opt/android/sdk

# 使用sdkmanager工具配置构建工具和平台版本
$ cd /opt/android/sdk/tools/bin/
$ ./sdkmanager "build-tools;29.0.6" "platforms;android-29" "platform-tools"
$ ./sdkmanager --list    #可以查看有哪些版本，自行选择对应的版本

# 增加系统环境变量
$ vim /etc/profile
export ANDROID_HOME=/opt/android/sdk
PATH=$PATH:$ANDROID_HOME:$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools:$ANDROID_HOME/emulator:$ANDROID_HOME/tools/bin

$ adb version
Android Debug Bridge version 1.0.41
Version 29.0.6-6198805
Installed as /opt/android/sdk/platform-tools/adb

```

### 安装 gradle

```sh
$ wget https://downloads.gradle-dn.com/distributions/gradle-6.3-all.zip
$ mkdir /opt/gradle
$ unzip gradle-6.3-all.zip -d /opt/gradle/
$ vim /etc/profile
export PATH=$PATH:/opt/gradle/gradle-6.3/bin
$ source /etc/profile
$ gradle -v
------------------------------------------------------------
Gradle 6.3
------------------------------------------------------------

Build time:   2020-03-24 19:52:07 UTC
Revision:     bacd40b727b0130eeac8855ae3f9fd9a0b207c60

Kotlin:       1.3.70
Groovy:       2.5.10
Ant:          Apache Ant(TM) version 1.10.7 compiled on September 1 2019
JVM:          1.8.0_201 (Oracle Corporation 25.201-b09)
OS:           Linux 3.10.0-693.el7.x86_64 amd64
```

## jenkins 流水线配置

### gradle-jdk.sh 打包脚本

```sh
#!/bin/sh
source /etc/profile
# /home/编译目录
cd /home/编译目录 && gradle clean
cd /home/编译目录 && gradle assembleRelease
# 打包编译完成，在项目的app/build/outputs/apk中可以找到debug版本或者是release版本。
```

### JenkinsFile 脚本

```json
pipeline {
    //使用标签 'master' 的节点
    agent {label 'master'}
    //环境变量
    environment {
        //成功运行特征
        JEN_FEATURE = ''
        //日志路径
        JEN_LOG = ''
    }

    stages {
        stage ('编译打包'){
            steps {
                sh 'cd /home/编译目录 && sudo git pull'
                sh 'sh gradle-apk.sh'

            }
        }

        stage ('上传'){
            steps {
                sh 'sshpass -p "password" scp /home/编译目录/app/build/outputs/apk/release/app-release.apk user@www.nginx.com:/home/上传目录/app-release.apk'
            }
        }

        stage ('生成qr'){
            steps {
                echo "生成test qr"
                sh "pwd && myqr 'https://www.nginx.com/app-release.apk' -n qrcode-`date +'%Y-%m-%d-%H%M%S'`.png -d /home/code/image"
                sh "cd /home/code/image && git add ."
                sh "cd /home/code/image && git commit -m 'new images'"
                sh "cd /home/code/image && git push origin main"
                echo "结束 end"
            }
        }
    }


    post {
        success {
            script {
                env.DATETIME = sh(script:"date", returnStdout: true).trim()
                env.QRCODE = sh(script:"ls -lhst /home/code/image | awk 'NR==2' | awk '{print \$10}'",returnStdout: true).trim()
                def job_name = "# ${JOB_NAME} 流水线 执行成功"
                def jenkinsid = """构建:  第 ${BUILD_ID} 次执行"""
                def JEN_production = "> 部署节点： jenkins"
                def build_url = "> 部署详情： [详情](${BUILD_URL})"
                def jen_date = "> 执行时间： ${env.DATETIME}"
                def jen_qrcode = "![qr](https://gitlab.com/jobcher/image/-/raw/main/${env.QRCODE})"

                dingtalk (
                    robot: 'e3999649-d3f-钉钉key-4c57333a327b',
                    type: 'MARKDOWN',
                    title: job_name,
                    text: [
                        job_name,
                        jenkinsid,
                        '',
                        '---',
                        jen_qrcode,
                        JEN_production,
                        '',
                        build_url,
                        '',
                        jen_date
                        ],
                    at: [
                        '手机号'
                        ]
                )
            }
        }

        failure {
            script {
                env.DATETIME = sh(script:"date", returnStdout: true).trim()
                def job_name = "# ${JOB_NAME} 流水线 执行失败"
                def jenkinsid = """构建:  第 ${BUILD_ID} 次执行"""
                def JEN_production = "> 部署节点： jenkins"
                def build_url = "> 部署详情： [详情](${BUILD_URL})"
                def jen_date = "> 执行时间： ${env.DATETIME}"

                dingtalk (
                    robot: 'e3999649-d3f-钉钉key-4c57333a327b',
                    type: 'MARKDOWN',
                    title: job_name,
                    text: [
                        job_name,
                        jenkinsid,
                        '',
                        '---',
                        JEN_production,
                        '',
                        build_url,
                        '',
                        jen_date
                        ],
                    at: [
                        '手机号'
                        ]
                )
            }
        }
    }
}

```
