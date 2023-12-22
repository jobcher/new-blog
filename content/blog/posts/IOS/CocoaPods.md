---
title: "CocoaPods 安装及碰到问题"
date: 2023-07-31
draft: false
featuredImage: "/images/cocoapods-orange-on-grey-1.jpg"
featuredImagePreview: "/images/cocoapods-orange-on-grey-1.jpg"
images: ['/images/cocoapods-orange-on-grey-1.jpg']
author: "jobcher"
tags: ["IOS 开发"]
categories: ["IOS 开发"]
series: ["IOS 开发"]
---
## 背景
CocoaPods 是OS X和IOS 下的第三类库管理工具，通过CocoaPods工具我们可以为项目添加被称为`Pods`的依赖库  

## 检查环境
```
ruby -v
gem -v
```
### 出现异常问题
> /System/Library/Frameworks/Ruby.framework/Versions/2.6/usr/lib/ruby/2.6.0/universal-darwin22/rbconfig.rb:21: warning: Insecure world writable dir /opt/homebrew/bin in PATH, mode 040777  
  
该警告信息表明在你的PATH环境变量中包含了一个“不安全可写”（Insecure world writable）的目录`/opt/homebrew/bin`。这可能会导致潜在的安全问题。  
为了解决这个警告，你需要修复`/opt/homebrew/bin`目录的权限，以使其不再被标记为“不安全可写”。  
#### 解决问题
```
chmod 755 /opt/homebrew/bin
chmod 755 /opt/homebrew
chmod 755 /opt/homebrew/sbin
```

## 安装cocoapods
输入安装命令
```sh
sudo gem install cocoapods
```
### 出现异常问题
>ERROR:  While executing gem ... (Gem::FilePermissionError)
    You don't have write permissions for the /System/Library/Frameworks/Ruby.framework/Versions/2.6/usr/lib/ruby/gems/2.6.0 directory.
        /Library/Ruby/Site/2.6.0/rubygems/installer.rb:714:in `verify_gem_home'
        /Library/Ruby/Site/2.6.0/rubygems/installer.rb:904:in `pre_install_checks'
        ……  
  
在 `macOS` 系统中，系统的Ruby目录通常是受保护的，并且普通用户没有对这些目录进行写操作的权限。为了解决这个问题，你应该避免在系统级别的Ruby目录中进行Gem的安装。相反，你应该使用用户级别的Gem安装目录。
#### 解决方案
```sh
mkdir -p ~/.gem/ruby/2.6.0
export PATH="$HOME/.gem/ruby/2.6.0/bin:$PATH"
source ~/.bash_profile
```
```sh
gem cocoapods install --user-install
```
### 第三方库
在iOS开发中，有许多常用的框架和库，用于加速应用程序的开发和提供各种功能。以下是一些常用的iOS开发框架和库：  
1. `UIKit`：UIKit是iOS开发的核心框架，它提供了构建用户界面和处理用户交互的基本组件，如视图（UIView）、控制器（UIViewController）、按钮（UIButton）、标签栏（UITabBarController）等。
2. `SwiftUI`：SwiftUI是一种较新的声明性用户界面框架，它使您可以使用Swift代码轻松构建用户界面，并自动处理界面状态的变化。它为多平台开发（iOS、iPadOS、macOS、watchOS）提供了便捷的方式。
3. `Core Data`：Core Data是Apple提供的数据持久化框架，用于在iOS应用程序中管理数据模型、数据存储和数据查询。
4. `Core Animation`：Core Animation是处理iOS界面动画的底层框架，可以实现各种动画效果，包括平移、旋转、缩放和淡入淡出等。
5. `Alamofire`：Alamofire是一个优秀的网络请求库，它简化了在iOS应用程序中进行HTTP请求的过程，支持多种请求方式和参数编码。
6. `Kingfisher`：Kingfisher是一个流行的图像加载和缓存库，它能够异步加载网络上的图片，并且可以自动进行内存和磁盘缓存。
7. `Firebase`：Firebase是谷歌提供的一组云服务，包括实时数据库、身份验证、云存储、推送通知等功能，用于简化iOS应用程序的后端开发。
8. `SwiftyJSON`：SwiftyJSON是一个简单易用的JSON解析库，用于在iOS应用程序中解析和处理JSON数据。
9. `MapKit`：MapKit是iOS提供的地图显示框架，可以在应用程序中集成地图、标记、路线等地理位置功能。
10. `CoreDataKit`：CoreDataKit是一个便捷的Core Data封装库，简化了Core Data的使用和管理。  
  
这只是iOS开发中的一小部分常用框架和库，还有很多其他优秀的开源库可供选择，根据您的应用程序需求，选择适合的框架和库可以加快开发速度并提供更好的用户体验。

## 初始化项目
在xcode中创建项目，进入xcode根目录  
```sh
pod init
```
编写Podfile
```yaml
# Uncomment the next line to define a global platform for your project
# platform :ios, '9.0'

target 'demo' do ### app名称
  # Comment the next line if you don't want to use dynamic frameworks
  use_frameworks!

  # Pods for demo
  pod 'SWXMLHash' ### 添加的依赖
  pod "ViewAnimator" ### 添加的依赖

  target 'demoTests' do
    inherit! :search_paths
    # Pods for testing
  end

  target 'demoUITests' do
    # Pods for testing
  end

end
```
安装依赖
```sh
pod install
```