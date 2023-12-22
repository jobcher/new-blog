---
title: "探索midjourney(二)：midjourney prompt初体验"
date: 2023-07-09
draft: false
featuredImage: "/images/ChatGPT-Plus/homeelegance_controllercolorfultranslucent_body_designed_bydiet_261aaa05-08ad-4fb5-b8a9-5961e85a50a1.png.webp"
featuredImagePreview: "/images/ChatGPT-Plus/homeelegance_controllercolorfultranslucent_body_designed_bydiet_261aaa05-08ad-4fb5-b8a9-5961e85a50a1.png.webp"
images: ["/images/ChatGPT-Plus.webp/homeelegance_controllercolorfultranslucent_body_designed_bydiet_261aaa05-08ad-4fb5-b8a9-5961e85a50a1.png.webp"]
authors: "jobcher"
tags: ["ChatGPT","midjourney"]
categories: ["人工智能"]
series: ["人工智能系列"]
---
紧接上文，我们继续探索midjourney prompt，这次我们来看看midjourney prompt的基础用法。

## midjourney 参数使用
参数是添加到提示中的选项，用于更改图像的生成方式。参数可以更改图像的长宽比、在中途模型版本之间切换、更改使用的 Upscaler 等等。参数始终添加到提示的末尾。您可以向每个提示添加多个参数  
![prompt](/images/ChatGPT-Plus/MJ_Parameters_example.png) 
> 许多 Apple 设备会自动将双连字符 (--) 更改为长破折号 (—)。midjourney接受两者！

### 基础参数列表
| 参数 | 描述 |
| --- | --- |
| --aspect, --ar | 更改生成图像的长宽比。 |
| --chaos <nummber 8-100> | 更改生成图像的混乱程度。 |
| --fast | 使用快速模式运行单个作业 |
| --iw <0-2> | 设置相对于文本粗细的图像提示粗细。默认值为 1。 |
| --no | 负面提示, `--no plant` 将生成没有植物的图像。 |
| --quality <.25, .5 ,1>, --q <.25, .5, 1>| 您想要花费多少渲染质量时间。默认值为 1。值越高，使用的 GPU 分钟数越多；较低的值使用较少 |
| --repeat <1-40>, --r <1-40> | 从单个提示创建多个作业。 `--repeat` 对于快速重新运行作业多次很有用。|
| --seed <integer between 0–4294967295> |Midjourney 机器人使用种子号来创建视觉噪声场（如电视静态），作为生成初始图像网格的起点。种子数是为每个图像随机生成的，但可以使用 `--seed` 或 `--sameseed` 参数指定。使用相同的种子编号和提示将产生相似的结局图像。|
| --stop <integer between 10–100> |  使用 --stop 参数在流程中途完成作业。以较早的百分比停止作业可能会产生更模糊、不太详细的结果。|
| --tile | 参数生成可用作重复图块以创建无缝图案的图像。|
| --Turbo | 使用 Turbo 模式运行单个作业。|
| --Weird <number 0-3000> | 使用实验性 --weird 参数探索不寻常的美学。|

### 参数示例
我了解那么多的参数变化，但是我不知道如何使用它们。让我们看看一些示例，以便您可以开始使用它们。

#### --aspect
`--aspect` 参数允许您更改生成图像的长宽比。默认值为 1，这意味着生成的图像将是正方形的。如果您想要更宽的图像，可以使用 `--aspect 1.5`。如果您想要更窄的图像，可以使用 `--aspect 0.5`。
```sh
apple --aspect 1.5
```

#### --chaos
`--chaos` 参数允许您更改生成图像的混乱程度。默认值为 50。值越高，图像越混乱。值越低，图像越清晰。如果您想要更清晰的图像，可以使用 `--chaos 25`。如果您想要更混乱的图像，可以使用 `--chaos 75`。
```sh
apple --chaos 25
```

#### --fast
`--fast` 参数允许您使用快速模式运行单个作业。快速模式将跳过 Upscaler 阶段，这意味着您将获得更快的结果，但图像的质量可能会降低。如果您想要更快的结果，可以使用 `--fast`。
```sh
apple --fast
```

#### --iw
`--iw` 参数允许您设置相对于文本粗细的图像提示粗细。默认值为 1。值越高，图像越粗。值越低，图像越细。如果您想要更粗的图像，可以使用 `--iw 2`。如果您想要更细的图像，可以使用 `--iw 0.5`。
```sh
apple --iw 2
```

#### --no
`--no` 参数允许您生成没有某些物体的图像。例如，如果您想要生成没有植物的图像，可以使用 `--no plant`。如果您想要生成没有人的图像，可以使用 `--no person`。如果您想要生成没有动物的图像，可以使用 `--no animal`。
```sh
apple --no plant
```

#### --quality
`--quality` 参数允许您设置您想要花费多少渲染质量时间。默认值为 1。值越高，使用的 GPU 分钟数越多；较低的值使用较少。如果您想要更高质量的图像，可以使用 `--quality 2`。如果您想要更低质量的图像，可以使用 `--quality 0.5`。
```sh
apple --quality 2
```

#### --repeat
`--repeat` 参数允许您从单个提示创建多个作业。如果您想要从单个提示创建 10 个作业，可以使用 `--repeat 10`。
```sh
apple --repeat 10
```

#### --seed
`--seed` 参数允许您使用种子号来创建视觉噪声场（如电视静态），作为生成初始图像网格的起点。种子数是为每个图像随机生成的，但可以使用 `--seed` 或 `--sameseed` 参数指定。使用相同的种子编号和提示将产生相似的结局图像。如果您想要使用种子号 123456789 生成图像，可以使用 `--seed 123456789`。
```sh
apple --seed 123456789
```

#### --stop
`--stop` 参数允许您在流程中途完成作业。以较早的百分比停止作业可能会产生更模糊、不太详细的结果。如果您想要在 50% 完成作业，可以使用 `--stop 50`。
```sh
apple --stop 50
```

#### --tile
`--tile` 参数允许您生成可用作重复图块以创建无缝图案的图像。如果您想要生成可用作重复图块以创建无缝图案的图像，可以使用 `--tile`。
```sh
apple --tile
```

#### --Turbo
`--Turbo` 参数允许您使用 Turbo 模式运行单个作业。如果您想要使用 Turbo 模式运行单个作业，可以使用 `--Turbo`。
```sh
apple --Turbo
```

#### --Weird
`--Weird` 参数允许您使用实验性 `--weird` 参数探索不寻常的美学。如果您想要探索不寻常的美学，可以使用 `--Weird`。
```sh
apple --Weird
```

## 分享几个有趣的demo
#### 1.方罐芒果味口香糖包装设计
```sh
Square cans Mango flavored gum packaging design
```
![Square](/images/ChatGPT-Plus/homeelegance_Square_cans_Mango_flavored_gum_packaging_design_36ba5831-88f1-439f-b776-2aca6bc98b43.png.webp)

#### 2.一个未知的世界，在一个等轴测立方体内部，发光到灰色表面。鲜艳的色彩
```sh
an unknown world, inside a isometric cube emitting Light onto a grey surface. Vibrant color
```
![2](/images/ChatGPT-Plus/homeelegance_an_unknown_world_inside_a_isometric_cube_emitting__fda72beb-e92b-47e3-abc0-e9adac8c6e77.png.webp)

#### 3.控制器，彩色，半透明机身，由节食公羊设计，高细节，8k，工业设计，白色背景，工作室照明3D，C4D，搅拌器，OCrender，Pinterest，Dribble，高细节
```sh
controller,colorful,translucent body, designed bydieter rams,high detail,8k,industrial design,whitebackground,studio lighting 3d, c4d, blender, OCrenderer, pinterest, dribbble, high detail
```
![3](/images/ChatGPT-Plus/homeelegance_controllercolorfultranslucent_body_designed_bydiet_261aaa05-08ad-4fb5-b8a9-5961e85a50a1.png.webp)


## 感谢关注