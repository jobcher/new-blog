---
title: "【福利】免费！本地部署！去除视频中移动的物体"
date: 2024-07-22
draft: false
featuredImage: "/images/ProPainter_pipeline.png"
featuredImagePreview: "/images/ProPainter_pipeline.png"
images: ["/images/ProPainter_pipeline.png"]
authors: "jobcher"
tags: ["daliy"]
categories: ["福利"]
series: ["福利系列"]
---
分享一款 去除视频中移动的物体。的本地部署软件，完全免费！  
## 效果
### 物体移除
![物体移除1](/images/object_removal1.gif)  
![物体移除2](/images/object_removal2.gif)  
### 水印去除
![水印去除1](/images/video_completion1.gif)  
![水印去除2](/images/video_completion2.gif)  
### 下载地址
[国内下载](https://pan.baidu.com/s/1XkQhzCzTtzVfgQg5heQQrA?pwd=jo38)  
### 代码仓库
[ProPainter](https://github.com/sczhou/ProPainter)  
## 安装
```sh
git clone https://github.com/sczhou/ProPainter.git
```
```sh
conda create -n propainter python=3.8 -y
conda activate propainter
cd ProPainter
pip3 install -r requirements.txt
```
### 版本要求
- CUDA >= 9.2
- PyTorch >= 1.7.1
- Torchvision >= 0.8.2 
## 开始使用
### 准备预训练模型
预训练模型从版本 V0.1.0 下载到 weights 文件夹。  
https://github.com/sczhou/ProPainter/releases/tag/v0.1.0  
```sh
weights
   |- ProPainter.pth
   |- recurrent_flow_completion.pth
   |- raft-things.pth
   |- i3d_rgb_imagenet.pt (for evaluating VFID metric)
   |- README.md
```
### 快速测试
```sh
# The first example (object removal)
python inference_propainter.py --video inputs/object_removal/bmx-trees --mask inputs/object_removal/bmx-trees_mask 
# The second example (video completion)
python inference_propainter.py --video inputs/video_completion/running_car.mp4 --mask inputs/video_completion/mask_square.png --height 240 --width 432
```
### 内存高效推理
视频修复通常需要大量 GPU 内存。在这里，我们提供了各种有助于内存高效推理的功能，有效避免内存不足（OOM）错误。您可以使用以下选项进一步减少内存使用量：  
- 通过减少 `--neighbor_length` （默认 10）来减少本地邻居的数量。
- 通过增加 `--ref_stride` （默认 10）来减少全局引用的数量。
- 设置 `--resize_ratio` （默认1.0）以调整处理视频的大小。
- 通过指定 `--width` 和 `--height` 设置较小的视频大小。
- 设置 `--fp16` 在推理过程中使用 fp16（半精度）。
- 减少子视频 `--subvideo_length` 的帧数（默认80），有效解耦GPU内存消耗和视频长度。  
  
### 训练
```sh
 # For training Recurrent Flow Completion Network
 python train.py -c configs/train_flowcomp.json
 # For training ProPainter
 python train.py -c configs/train_propainter.json
 ```

## web端
[Hugging Face 在线演示](https://openxlab.org.cn/apps/detail/ShangchenZhou/ProPainter)  
[OpenXLab 在线演示](https://huggingface.co/spaces/sczhou/ProPainter)  