---
title: "Maven 安装编译"
date: 2022-03-03
draft: false
author: "jobcher"
tags: ["work"]
categories: ["work"]
series: ["问题库系列"]
---

# Maven 安装编译

Maven 就是专门为 Java 项目打造的管理和构建工具，它的主要功能有：

- 提供了一套标准化的项目结构；
- 提供了一套标准化的构建流程（编译，测试，打包，发布……）；
- 提供了一套依赖管理机制。

默认结构：

```sh
a-maven-project
├── pom.xml
├── src
│   ├── main
│   │   ├── java
│   │   └── resources
│   └── test
│       ├── java
│       └── resources
└── target

```

项目的根目录`a-maven-project`是项目名，  
它有一个`项目描述`文件`pom.xml`，  
存放`Java源码`的目录是`src/main/java`，  
存放`资源文件`的目录是`src/main/resources`，  
存放`测试源码`的目录是`src/test/java`，  
存放`测试资源`的目录是`src/test/resources`，  
最后，所有编译、`打包生成的文件`都放在`target`目录里。  
这些就是一个 Maven 项目的标准目录结构。

pom.xml 文件:

```xml
<project ...>
	<modelVersion>4.0.0</modelVersion>
	<groupId>com.itranswarp.learnjava</groupId>
	<artifactId>hello</artifactId>
	<version>1.0</version>
	<packaging>jar</packaging>
	<properties>
        ...
	</properties>
	<dependencies>
        <dependency>
            <groupId>commons-logging</groupId>
            <artifactId>commons-logging</artifactId>
            <version>1.2</version>
        </dependency>
	</dependencies>
</project>
```

`groupId`类似于 Java 的包名，通常是公司或组织名称，  
`artifactId`类似于 Java 的类名，通常是项目名称，  
`version`，一个 Maven 工程就是由`groupId，artifactId和version`作为唯一标识。  
我们在引用其他第三方库的时候，也是通过这 3 个变量确定。

依赖`commons-logging`：

```xml
<dependency>
    <groupId>commons-logging</groupId>
    <artifactId>commons-logging</artifactId>
    <version>1.2</version>
</dependency>
```

使用`<dependency>`声明一个依赖后，Maven 就会自动下载这个依赖包并把它放到`classpath`中。

## 安装 Maven
