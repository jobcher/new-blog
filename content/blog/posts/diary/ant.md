---
title: "Ant中如何添加第三方jar包依赖"
date: 2022-04-14
draft: false
author: "jobcher"
tags: ["javascript"]
categories: ["日常"]
series: ["日常系列"]
---

# Ant 中如何添加第三方 jar 包依赖

如果使用 ant 进行 java 项目的编译部署，那怎么添加第三方 jar 包的依赖呢？方法如下：

1. 在项目的根目录下创建 lib 目录，并把所有需要的第三方 jar 包放到此目录下。
2. 在 build.xml 中依次添加：path、property，并在 javac 中添加 classpath，添加 unjar。完整配置如下：

```xml
<?xml version="1.0" encoding="UTF-8"?>
<project name="MyTool" default="build" basedir=".">
  <description>The ant project to build MyTool.</description>
  <property name="srcDir" location="src" description="源文件的存放目录" />
  <property name="libDir" location="lib" description="第三方jar包的存放目录" />
  <property name="antDir" location="ant" description="编译后所有文件存放的根目录" />
  <property name="binDir" location="${antDir}/bin" description="编译后class文件的存放目录" />
  <property name="jarDir" location="${antDir}/jar" description="打包后jar包的存放目录" />
  <property name="jarFile" location="${jarDir}/MyTool.jar" description="打包后jar包存放的完整路径" />
  <property name="package" value="com.xiboliya.mytool" description="包名" />
  <property name="mainClass" value="MyTool" description="主类名" />
  <property name="resFromDir" location="res" description="资源文件的源目录" />
  <property name="resToDir" location="${binDir}/res" description="资源文件的目标目录" />
  <path id="libPath" description="编译时依赖的第三方jar包的存放路径">
    <fileset dir="${libDir}" erroronmissingdir="false">
      <include name="*.jar" />
    </fileset>
  </path>
  <property name="classpath" refid="libPath" description="编译时依赖的第三方jar包的存放路径" />
  <target name="init" description="初始化">
    <delete dir="${binDir}" />
    <delete dir="${jarDir}" />
    <mkdir dir="${binDir}" />
    <mkdir dir="${jarDir}" />
    <copy todir="${resToDir}" description="复制资源文件">
      <fileset dir="${resFromDir}" />
    </copy>
  </target>
  <target name="compile" depends="init" description="编译代码">
    <javac srcdir="${srcDir}" destdir="${binDir}" classpath="${classpath}" includeAntRuntime="false">
      <compilerarg value="-Xlint:unchecked"/>
      <compilerarg value="-Xlint:deprecation"/>
    </javac>
  </target>
  <target name="unjarLib" depends="init" description="解压第三方jar包，以便于重新打包入程序jar包中">
    <unjar dest="${binDir}">
      <fileset dir="${libDir}">
        <include name="**/*.jar" />
      </fileset>
      <patternset>
        <exclude name="META-INF"/>
        <exclude name="META-INF/MANIFEST.MF"/>
      </patternset>
    </unjar>
  </target>
  <target name="makeJar" depends="init,compile,unjarLib" description="生成jar包">
    <jar destfile="${jarFile}" basedir="${binDir}" excludes="**/Thumbs.db" description="打包为jar文件，并排除Thumbs.db文件">
      <manifest>
        <attribute name="Main-Class" value="${package}.${mainClass}" />
      </manifest>
    </jar>
  </target>
  <target name="build" depends="init,compile,makeJar" description="编译并打包">
    <echo message="Ant is building the project." />
  </target>
</project>

```

欢迎关注我的博客[www.jobcher.com](https://www.jobcher.com/)
