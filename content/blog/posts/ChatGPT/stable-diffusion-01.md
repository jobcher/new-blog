---
title: ubuntu 安装 stable-diffusion-webui
date: 2024-01-10
draft: false
author: 'jobcher'
featuredImage: '/images/sd-logo.jpeg'
featuredImagePreview: '/images/sd-logo.jpeg'
images: ['/images/wallpaper/sd-logo.jpeg']
tags: ['stable diffusion']
categories: ['stable diffusion']
series: ['stable diffusion']
---
## 背景
`Stable Diffusion` (稳定扩散) 是一个扩散模型，2022年8月由德国CompVis协同Stability AI和Runway发表论文，并推出相关程序。  
`AUTOMATIC1111`开发了`图形化界面`：「Stable Diffusion WebUI」，这是能用AI技术生成图片的开源软件，只要给定一组描述文本，AI就会开始绘图(准确的说是「算图」或「生图」)；亦能模仿现有的图片，生成另一张图片。甚至给它一部分涂黑的图片，AI也能按照你的意愿将图片填上适当的内容。除此之外还支持自行训练模型加强生图效果。  
本篇文章就是介绍如何安装 stable-diffusion-webui  
## 安装conda
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
## 搭建环境
安装软件
```sh
apt -y update -qq
apt -y install -qq aria2
```
创建python环境
```sh
conda create -n sd-web python=3.9
```
下载github代码
```sh
mkdir ~/sd-web
cd ~/sd-web/
git clone https://github.jobcher.com/gh/https://github.com/AUTOMATIC1111/stable-diffusion-webui.git
```
下载模型
```sh
aria2c --console-log-level=error -c -x 16 -s 16 -k 1M https://huggingface.jobcher.com/https://huggingface.co/runwayml/stable-diffusion-v1-5/resolve/main/v1-5-pruned-emaonly.safetensors -d ~/sd-web/stable-diffusion-webui/extensions/sd-webui-controlnet/models -o v1-5-pruned-emaonly.safetensors
```
安装依赖
```sh
cd ~/cd-web/stable-diffusion-webui
# 先安装requirements_versions.txt
pip install -r requirements_versions.txt -i https://pypi.douban.com/simple --trusted-host=pypi.douban.com --verbose basicsr --use-pep517
pip install -r requirements.txt -i https://pypi.douban.com/simple --trusted-host=pypi.douban.com --verbose basicsr --use-pep517
```
开放外部访问
```sh
export COMMANDLINE_ARGS="--listen"
```
运行软件
```sh
cd ~/cd-web/stable-diffusion-webui
bash ./webui.sh
```
出现一下画面说明已安装成功  
![success](/images/1704866846903.jpg)  
## 正常使用软件
进入你的地址例如 127.0.0.1:7860你会看到一下界面，开始正常训练模型吧    
![demo](/images/1704867402210.jpg)  
