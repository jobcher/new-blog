---
title: "3D Gaussian Splatting：3D模型渲染"
date: 2023-08-30
draft: false
featuredImage: "/images/graphdeco_logo.png"
featuredImagePreview: "/images/graphdeco_logo.png"
images: ["/images/graphdeco_logo.png"]
author: "jobcher"
tags: ["daliy"]
categories: ["福利"]
series: ["福利系列"]
---
## 简介
只需要一组照片或者一个视频，就能快速地生成一个3D模型。  
它使了用一种叫做“3D高斯函数”的数学工具来表示这个3D模型，并找到了一种更快的算法来渲染（即生成）3D模型。  

## 介绍
想象一下，你有一堆照片或视频，你想从一个全新的角度看这些场景。3D Gaussian Splatting 就是一个能让你做到这一点的高级工具。它用一种特别快和高质量的方式来“重建”这些场景，让你能够从任何角度观看它们，就像你实际站在那里一样。  
  
这个工具的“大脑”使用了一种叫做3D高斯的数学模型，这个模型能够非常精确地描述场景的每一个细节。更酷的是，这个工具还能实时地显示这些新角度的场景，这意味着你不必等待很长时间就能看到结果。  
  
简而言之，这是一个能让你以全新、快速和高质量的方式探索照片和视频场景的工具。  

## 使用
### 硬件要求
- 具有计算能力 7.0+ 的 CUDA 就绪 GPU
- 24 GB VRAM（用于训练论文评估质量）
### 软件要求
- Conda（推荐使用，以便于设置）
- 用于 PyTorch 扩展的 C++ 编译器（我们使用 Visual Studio 2019 for Windows）
- 用于 PyTorch 扩展的 CUDA SDK 11，在 Visual Studio 之后安装（我们使用 11.8，11.6 存在已知问题）
- C++编译器和CUDA SDK必须兼容
  
### 设置
```sh
SET DISTUTILS_USE_SDK=1 # Windows only
conda env create --file environment.yml
conda activate gaussian_splatting
```
请注意，此过程假设您安装了 CUDA SDK 11，而不是 12。有关修改，请参阅下文。
```sh
conda config --add pkgs_dirs <Drive>/<pkg_path>
conda env create --file environment.yml --prefix <Drive>/<env_path>/gaussian_splatting
conda activate <Drive>/<env_path>/gaussian_splatting
```
运行
```sh
python train.py -s <path to COLMAP or NeRF Synthetic dataset>
```
评估
```sh
python train.py -s <path to COLMAP or NeRF Synthetic dataset> --eval # Train with train/test split
python render.py -m <path to trained model> # Generate renderings
python metrics.py -m <path to trained model> # Compute error metrics on renderings
```

项目地址及演示：[链接](https://repo-sam.inria.fr/fungraph/3d-gaussian-splatting/)  
GitHub：[链接](https://github.com/graphdeco-inria/gaussian-splatting)  

## 演示视频
{{< youtube kShNYOuDnlI >}}  
<br/>
{{< youtube mD0oBE9LJTQ >}}  