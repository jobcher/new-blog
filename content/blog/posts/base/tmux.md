---
title: "Tmux 安装和使用教程"
date: 2024-06-13
draft: false
authors: "jobcher"
tags: ["运维"]
featuredImage: "/images/tmux-3.jpg"
featuredImagePreview: "/images/tmux-3.jpg"
images: ['/images/tmux-3.jpg','/images/tmux-1.png','/images/tmux-2.png',]
categories: ["基础"]
series: ["基础知识系列"]
---
**tmux** 是一个终端 multiplexer，它可以让你在一个终端中开启多个会话，并且可以在一个终端中切换多个会话。
## 安装
tmux 安装很简单，直接在终端中输入以下命令即可：
```sh
# Ubuntu 或 Debian
sudo apt-get install tmux

# CentOS 或 Fedora
sudo yum install tmux

# Mac
brew install tmux
```
## 使用
安装完成后，键入`tmux`命令，就进入了 `Tmux` 窗口。
```sh
tmux
```
![tmux-1](/images/tmux-3.jpg)  
按下`Ctrl+d`或者显式输入`exit`命令，就可以退出 Tmux 窗口。  

```sh
exit
```
### 前缀键
Tmux 窗口有大量的快捷键。所有快捷键都要通过前缀键唤起。默认的前缀键是`Ctrl+b`，即先按下`Ctrl+b`，快捷键才会生效。  
举例来说，帮助命令的快捷键是`Ctrl+b ?`。它的用法是，在 Tmux 窗口中，先按下`Ctrl+b`，再按下`?`，就会显示帮助信息。  
然后，按下 `ESC` 键或`q`键，就可以退出帮助。
### 快捷键
#### 面板（pane）指令
|前缀	|指令	|描述|
|:---:|:---:|:---|
|Ctrl+b|	"	|当前面板上下一分为二，下侧新建面板|
|Ctrl+b|	%	|当前面板左右一分为二，右侧新建面板|
|Ctrl+b|	x	|关闭当前面板（关闭前需输入y or n确认）|
|Ctrl+b|	z	|最大化当前面板，再重复一次按键后恢复正常（v1.8版本新增）|
|Ctrl+b|	!	|将当前面板移动到新的窗口打开（原窗口中存在两个及以上面板有效）|
|Ctrl+b|	;	|切换到最后一次使用的面板|
|Ctrl+b|	q	|显示面板编号，在编号消失前输入对应的数字可切换到相应的面板|
|Ctrl+b|	{	|向前置换当前面板|
|Ctrl+b|	}	|向后置换当前面板|
|Ctrl+b|	Ctrl+o	|顺时针旋转当前窗口中的所有面板|
|Ctrl+b|	方向键	|移动光标切换面板|
|Ctrl+b|	o	|选择下一面板|
|Ctrl+b|	空格键	|在自带的面板布局中循环切换|
|Ctrl+b|	Alt+方向键	|以5个单元格为单位调整当前面板边缘|
|Ctrl+b|	Ctrl+方向键	|以1个单元格为单位调整当前面板边缘（Mac下被系统快捷键覆盖）|
|Ctrl+b|	t	|显示时钟|

#### 系统指令
|前缀	|指令	|描述|
|:---:|:---:|:---|
|Ctrl+b|	?	|显示快捷键帮助文档|
|Ctrl+b|	d	|断开当前会话|
|Ctrl+b|	D	|选择要断开的会话|
|Ctrl+b|	Ctrl+z	|挂起当前会话|
|Ctrl+b|	r	|强制重载当前会话|
|Ctrl+b|	s	|显示会话列表用于选择并切换|
|Ctrl+b|	:	|进入命令行模式，此时可直接输入ls等命令|
|Ctrl+b|	[	|进入复制模式，按q退出|
|Ctrl+b|	]	|粘贴复制模式中复制的文本|
|Ctrl+b|	~	|列出提示信息缓存|

#### 窗口（window）指令
|前缀	|指令	|描述|
|:---:|:---:|:---|
|Ctrl+b|	?	|显示快捷键帮助文档|
|Ctrl+b|	d	|断开当前会话|
|Ctrl+b|	D	|选择要断开的会话|
|Ctrl+b|	Ctrl+z	|挂起当前会话|
|Ctrl+b|	r	|强制重载当前会话|
|Ctrl+b|	s	|显示会话列表用于选择并切换|
|Ctrl+b|	:	|进入命令行模式，此时可直接输入ls等命令|
|Ctrl+b|	[	|进入复制模式，按q退出|
|Ctrl+b|	]	|粘贴复制模式中复制的文本|
|Ctrl+b|	~	|列出提示信息缓存|
