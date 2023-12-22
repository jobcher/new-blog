---
title: "在windows上安装appium"
date: 2023-05-15
draft: false
authors: "jobcher"
tags: ["运维"]
featuredImage: "/images/appium-logo-horiz.png"
featuredImagePreview: "/images/appium-logo-horiz.png"
categories: ["基础"]
series: ["基础知识系列"]
---
## 在windows上安装appium
`Appium` 主要用于软件测试自动化领域，以帮助确定给定应用程序的功能是否按预期工作。与其他类型的软件测试相比，UI 自动化允许测试人员编写代码，在应用程序的实际 UI 中演练用户方案，尽可能模拟现实世界中发生的情况，同时实现自动化的各种好处，包括速度、规模和一致性。

### 安装nodejs
Appium是基于Node.js构建的,所以首先需要安装Node.js  
[下载地址：https://nodejs.org/en/download/](https://nodejs.org/en/download/)  
下载并安装，验证是否安装成功  
```sh
node -v
```

### 安装appium
安装Appium。在命令提示符下运行:
```sh
npm install -g appium
```

### 安装appium-doctor
安装appium-doctor。在命令提示符下运行:
```sh
npm install -g appium-doctor
```

### 安装appium-desktop
安装appium-desktop。在命令提示符下运行:
```sh
npm install -g appium-desktop
```

### 安装jdk
下载地址：[https://www.oracle.com/java/technologies/javase/javase-jdk8-downloads.html](https://www.oracle.com/java/technologies/javase/javase-jdk8-downloads.html)
下载并安装，验证是否安装成功
```sh
java -version
```

### 安装android-sdk
下载地址：[https://developer.android.com/studio](https://developer.android.com/studio)
下载并安装，验证是否安装成功
```sh
adb version
```

### 配置环境变量
在系统环境变量中添加以下变量
```sh
ANDROID_HOME = D:\Android\sdk
JAVA_HOME = C:\Program Files\Java\jdk1.8.0_311
Path = %ANDROID_HOME%\platform-tools;%ANDROID_HOME%\tools;%JAVA_HOME%\bin
```
验证是否配置成功
```sh
appium-doctor
```

### 启动appium
```sh
appium
```

### 启动appium-desktop
```sh
appium-desktop
```

### 测试脚本
创建测试脚本test.py
```python
from appium import webdriver
import time

desired_caps = {}

desired_caps['platformName'] = 'Android'
desired_caps['platformVersion'] = '10'
desired_caps['deviceName'] = 'Android Emulator'
desired_caps['appPackage'] = 'com.android.calculator2'
desired_caps['appActivity'] = '.Calculator'

driver = webdriver.Remote('http://localhost:4723/wd/hub', desired_caps)

driver.find_element_by_id("com.android.calculator2:id/digit_2").click()
driver.find_element_by_id("com.android.calculator2:id/op_add").click()
driver.find_element_by_id("com.android.calculator2:id/digit_3").click()
driver.find_element_by_id("com.android.calculator2:id/eq").click()

time.sleep(5)

driver.quit()
```

### 运行测试脚本
```sh
python test.py
```
