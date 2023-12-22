---
title: "ant build.xml 编写"
date: 2022-04-21
draft: false
authors: "jobcher"
tags: ["java"]
categories: ["日常"]
series: ["日常系列"]
---

# ant build.xml 编写

## 生成 build.xml

`Eclipse` 自动生成 Ant 的`Build.xml` 配置文件,生成的方法很隐蔽  
选择你要生成`Build.xml`文件的项目,右键. `Export-> General -> Ant Buildfiles .`  
点 Next,选择项目，再点`Finish`.

## 编写 build.xml

```xml
<?xml version="1.0" encoding="UTF-8" standalone="no"?>

<!-- 每个构建文件对应一个项目。<project>标签时构建文件的根标签。它可以有多个内在属性，就如代码中所示，其各个属性的含义分别如下。
(1) default表示默认的运行目标，这个属性是必须的。
(2) basedir表示项目的基准目录。
(3) name表示项目名。
(4) description表示项目的描述。
 -->
<project default="build" name="Sort">
	<!-- 设置属性或文件路径，读取属性使用${property}，value路径默认项目根目录 -->
	<property file="ant/builds.properties" />

	<property name="src.dir" value="src/statics" />

	<property name="classes.dir" value="ant/classes" />

	<property name="lib.dir" value="lib" />

	<property name="dist.dir" value="ant/dist" />

	<!-- 定义classpath -->
	<path id="master-classpath">
		<fileset file="${lib.dir}/*.jar" />
		<pathelement path="${classes.dir}" />
	</path>

	<!--一个项目标签Project包含多个target标签，一个target标签可以依赖其他的target标签
		在生成可执行文件之前必须先编译该文件，因策可执行文件的target依赖于编译程序的 target。

		(1).name表示标明，这个属性是必须的。
		(2).depends表示依赖的目标。
		(3)if表示仅当属性设置时才执行。
		(4)unless表示当属性没有设置时才执行。
		(5)description表示项目的描述。
		Ant的depends属性指定了target的执行顺序。Ant会依照depends属性中target出现顺序依次执行每个target。在执行之前，首先需要执行它所依赖的target。
	 -->
	<!-- 初始化任务 -->
	<target name="init">
		<!-- 输出标签 ，${init}是builds.properties中的属性 -->
		<echo message="  Available Targets:"/>
        <echo message="-------------------------------------------------------"/>
        <echo message="  init ${init}   ..."/>
        <echo message="-------------------------------------------------------"/>
	</target>

	<!-- 编译 -->
	<target name="compile" depends="init" description="compile the source files">
		<!-- 删除文件夹 -->
		<delete dir="${classes.dir}" />
		<!-- 创建文件夹 -->
		<mkdir dir="${classes.dir}" />
		<!-- 编译java生成class文件 ，其属性如下
			(1).srcdir表示源程序的目录。
			(2).destdir表示class文件的输出目录。
			(3).include表示被编译的文件的模式。
			(4).excludes表示被排除的文件的模式。
			(5).classpath表示所使用的类路径。
			(6).debug表示包含的调试信息。
			(7).optimize表示是否使用优化。
			(8).verbose 表示提供详细的输出信息。
			(9).fileonerror表示当碰到错误就自动停止。
		 -->
		<javac srcdir="${src.dir}" destdir="${classes.dir}">
			<!-- 编译需要的jar包 引用前面设置的class-path -->
			<classpath refid="master-classpath" />
		</javac>
	</target>


	<!-- 打包成jar -->
	<target name="pack" description="make .jar file">

		<delete dir="${dist.dir}" />

		<mkdir dir="${dist.dir}" />
		<!-- 该标签用来生成一个JAR文件，其属性如下
			(1) destfile表示JAR文件名。
			(2) basedir表示被归档的文件名。要操作的文件路径
			(3) includes表示别归档的文件模式。
			(4) exchudes表示被排除的文件模式。
		 -->
		<jar destfile="${dist.dir}/hello.jar" basedir="${classes.dir}">
			<!-- 不包含的类或内容 -->
			<exclude name="**/*Test.*" />
		</jar>

	</target>

	<!-- 生成zip压缩包 -->
	<target name="zip">
	    <delete dir="${release-dir}" />
		<mkdir dir="${release-dir}" />
		<!-- 该标签用来生成一个zip文件，其属性如下
			(1) destfile表示zip文件名。
			(2) basedir表示被归档的文件名。 要操作的文件路径
			(3) includes表示别归档的文件模式。
			(4) exchudes表示被排除的文件模式。
		 -->
		<zip destfile="${release-dir}/antTest.zip" update="true"
			       basedir="ant" />
	</target>
</project>
```

欢迎关注我的博客[www.jobcher.com](https://www.jobcher.com/)
