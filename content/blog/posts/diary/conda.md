---
title: "优雅的使用Conda管理python环境"
date: 2022-12-14
draft: false
featuredImage: "/images/python-conda.png"
featuredImagePreview: "/images/python-conda.png"
author: "jobcher"
tags: ["daliy"]
categories: ["日常"]
series: ["日常系列"]
---

## 背景

很多时候,避免不了同时使用 python2 和 python3 的环境,也避免不了不同的工作所需要不同版本的库文件,比如在想用 TensorFlow 较早版本的同时;还想运行 Pytorch 最新版；还想顺便学习 Nao 机器人编程,学习 Django 后台,这个时候,一款非常好用的包管理工具就显得十分重要了,这就是我写这篇博客的原因,这篇博客将会讲解：

- [x] 如何安装 conda;
- [x] 如何更换 conda 的下载源;
- [x] 如何使用 canda;

## Linux 安装 conda
在 Ubuntu 上安装 Anaconda 的步骤如下：  
  
1. 首先，你需要下载 Anaconda 的安装包。你可以从 Anaconda 的官方网站上下载最新版本的 Anaconda for Linux。选择适合你的系统的版本（Python 3.x）。  
  
   访问下载链接：[https://www.anaconda.com/products/distribution#download-section](https://www.anaconda.com/products/distribution#download-section)  
  
2. 下载完成后，你可以在终端中导航到下载的文件所在的目录。你可以使用 `cd` 命令来改变目录。例如，如果你的下载文件在 Downloads 文件夹中，你可以输入以下命令：
  
```bash
cd ~/Downloads
```

3. 然后，你需要运行 bash 命令来安装 Anaconda。假设你下载的 Anaconda 文件名为 "Anaconda3-2020.02-Linux-x86_64.sh"，你可以输入以下命令：

```bash
bash Anaconda3-2020.02-Linux-x86_64.sh
```
请注意，你需要将上述命令中的 "Anaconda3-2020.02-Linux-x86_64.sh" 替换为你实际下载的文件名。  
4. 接下来，你会看到 Anaconda 的许可协议。按 `Enter` 键滚动到底部，然后输入 'yes' 来接受许可协议。  
5. 然后，你需要确认 Anaconda 的安装位置。你可以选择默认位置或输入新的位置。
6. 安装完成后，你会看到一个提示，询问你是否希望 Anaconda3 添加到你的 PATH。你应该输入 'yes'。
7. 最后，你需要激活安装。你可以通过关闭并重新打开终端或运行以下命令来完成此操作：
```bash
source ~/.bashrc
```
8. 验证安装。在终端中输入以下命令：
```bash
conda list
```
如果安装成功，这个命令会显示一个已安装的包的列表。  

## windows 安装 conda

在安装时这两个选项需要点上：  
![conda_install](/images/conda_install.png)

### 更换 conda 的下载源

`Conda官方`的下载源太慢了,而且经常会出现 HTTPERROR 之类的错误,如果想要用 Conda 愉快的创建不同工作环境,愉快的下载安装各种库,那么换下载源是必不可少的

```sh
conda config --add channels https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/free/
conda config --add channels https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud/conda-forge
conda config --add channels https://mirrors.tuna.tsinghua.edu.cn/anaconda/cloud/msys2/

conda config --set show_channel_urls yes
# 设置搜索时显示通道地址
```

具体操作同时按 Win+R 键打开运行窗口,输入 cmd,回车：  
![cmd-conda](/images/cmd-conda.png)

将上面的命令全部复制,到命令行里单击右键就会自动执行复制的命令,添加清华源

## 使用 conda

查看环境

```sh
conda info -e
conda info --envs
```

创建环境

```sh
conda create -n name python=3.6
# name参数指定虚拟环境的名字,python参数指定要安装python的版本,但注意至少需要指定python版本或者要安装的包,在后一种情况下,自动安装最新python版本
# 例如
conda create -n jobcher pillow numpy python=2.7.14
# 创建名字为naoqi,Python版本为2.7.14的虚拟环境,同时还会安装上pillow numpy这两个库
```

环境切换

```sh
conda activate jobcher
# 切换到jobcher环境下,在切换环境后,所执行的Pip命令,Python命令,都是更改当前环境下的,不会影响到其他的环境
conda deactivate
# 退出当前环境
```


>欢迎关注我的博客[test.jobcher.com](https://test.jobcher.com/)

