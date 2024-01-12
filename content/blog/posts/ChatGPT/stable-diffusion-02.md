---
title: ubuntu 安装 ComfyUI
date: 2024-01-12
draft: false
author: 'jobcher'
featuredImage: '/images/ComfyUI-logo.png'
featuredImagePreview: '/images/ComfyUI-logo.png'
images: ['/images/ComfyUI-logo.png']
tags: ['stable diffusion']
categories: ['stable diffusion']
series: ['stable diffusion']
---
## 背景
ComfyUI 是用于稳定扩散的基于节点的用户界面。ComfyUI 由 Comfyanonymous 于 2023 年 1 月创建，他创建了该工具来学习稳定扩散的工作原理。
## 效果
webui和ComfyUI之间的区别,相比较webUI，ComfyUI更工业化，更符合高级使用者的配置  
![](/images/b5i0v5krtcdb1-1024x435.png)  
## 安装
### 安装本体
下载软件
```sh
mkdir ~/sd-web
cd ~/sd-web
git clone https://github.jobcher.com/gh/https://github.com/comfyanonymous/ComfyUI.git
```
环境依赖
```sh
cd ~/sd-web/ComfyUI
conda create -n ComfyUI python=3.10
pip install -r requirements.txt -i https://pypi.douban.com/simple --trusted-host=pypi.douban.com
```
下载sd_xl_turbo模型
```sh
aria2c --console-log-level=error -c -x 16 -s 16 -k 1M https://huggingface.jobcher.com/https://huggingface.co/stabilityai/sdxl-turbo/resolve/main/sd_xl_turbo_1.0_fp16.safetensors -d ~/sd-web/ComfyUI/models/checkpoints -o sd_xl_turbo_1.0_fp16.safetensors
```
启动服务
```sh
cd ~/sd-web/ComfyUI
python main.py --listen --port 6006 --cuda-device 1
```
## webUI共享模型
```sh
cd 
mv extra_model_paths.yaml..example extra_model_paths.yaml
```
编辑参数
```sh
vim extra_model_paths.yaml
```
修改 `base_path: path/to/stable-diffusion-webui/` 改为你的webui实际地址，例如：
`base_path: ~/sd-web/stable-diffusion-webui/ `
```yaml
#config for a1111 ui
#all you have to do is change the base_path to where yours is installed
a111:
    base_path: path/to/stable-diffusion-webui/ ## 这里改为你实际的webui地址

    checkpoints: models/Stable-diffusion
    configs: models/Stable-diffusion
    vae: models/VAE
    loras: |
         models/Lora
         models/LyCORIS
    upscale_models: |
                  models/ESRGAN
                  models/RealESRGAN
                  models/SwinIR
    embeddings: embeddings
    hypernetworks: models/hypernetworks
    controlnet: models/ControlNet
```
重启服务
```sh
python main.py --listen --port 6006 --cuda-device 1
```