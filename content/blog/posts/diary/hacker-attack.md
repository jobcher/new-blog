---
title: "异常流量分析：图片库服务黑客入侵"
date: 2025-03-25
draft: false
featuredImage: "/images/hacker-attack.jpg"
featuredImagePreview: "/images/hacker-attack.jpg"
images: ["/images/hacker-attack.jpg"]
authors: "jobcher"
tags: ["daliy"]
categories: ["日常"]
series: ["日常系列"]
---
## 背景
最近这两天prometheus一直再对服务器的下行带宽告警，检查告警发现是一台图片服务器，这台图片服务器是很早部署的一台nginx反向代理本地静态文件的服务器，当时没有做任何的安全防护，只是简单做了一个上传接口暴露给外网，导致这台服务器被黑客入侵，上传了大量文件，同时对这些文件有大量的下载请求，导致服务器的带宽被打满，同时也导致了prometheus的告警。

## 发现
1. 发现这台服务器的带宽被打满了，所以只能先从这台服务器开始排查。先检查nginx所有大量请求的日志  
```sh
cd /var/log/nginx
awk '{print $7}' access.log | sort | uniq -c | sort -nr | head -n 10
```
![server-log](/images/hacker-attack01.png)  

2. 这个是一个图片服务器，看起来都是图片的请求，好像没有发现异常的请求。但是这个图片为什么有那么高的流量呢？检查这个图片，发现图片无法打开，所以我换二进制的方式打开图片  
![hex-open](/images/hacker-attack02.png)  
> 发现图片的二进制头部是 `FFmpeg`, 所以这是一个伪装成图片的视频文件  
  
## 解决
1. 检查nginx的配置文件，上传接口拒绝所有外部请求
```
server{
    location /upload {
        proxy_pass http://192.168.1.1:8080;

        # 拒绝所有外部请求
        allow 192.168.1.0/24;
        deny all;
    }
}
```
  
2. 查找最早上传的文件,创建python文件，并在图片服务器执行
```
import os

# 查找最早包含 FFmpeg 的 JPG 文件
def find_earliest_jpg_with_ffmpeg(directory):
    earliest_file = None
    earliest_timestamp = float('inf')

    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.lower().endswith('.jpg'):
                file_path = os.path.join(root, file)
                try:
                    with open(file_path, 'rb') as f:
                        file_header = f.read(55)
                        print(file_header)
                        f.close()
                        if b'FFmpeg\tService01' in file_header:
                            timestamp = os.path.getmtime(file_path)
                            if timestamp < earliest_timestamp:
                                earliest_timestamp = timestamp
                                earliest_file = file_path
                except Exception as e:
                    print(f"Error processing file {file_path}: {e}")

    if earliest_file:
        print(f"最早的包含 FFmpeg 的 JPG 文件是: {earliest_file}")

# 调用函数并传入目录路径
find_earliest_jpg_with_ffmpeg('/var/www/html/images/')
```
3. 删除这些`FFmpeg`文件，创建python文件，并在图片服务器执行
```
import os

# 检查文件头是否包含 FFmpeg
def delete_file_if_header_contains_ffmpeg(file_path):
    try:
        with open(file_path, 'rb') as file:
            file_header = file.read(55)
            print(f"检查文件: {file_path}，文件头: {file_header}")
            file.close()
            if b'FFmpeg\tService01' in file_header:
                os.remove(file_path)
                print(f"文件 {file_path} 已被删除")
            else:
                print(f"文件 {file_path} 不符合删除条件")
    except FileNotFoundError:
        print(f"文件 {file_path} 找不到")
    except Exception as e:
        print(f"发生错误: {e}")

# 遍历目录并删除 JPG 文件
def traverse_and_delete_jpg_files(directory):
    try:
        for root, dirs, files in os.walk(directory):
            for file in files:
                if file.lower().endswith('.jpg'):
                    file_path = os.path.join(root, file)
                    delete_file_if_header_contains_ffmpeg(file_path)
    except Exception as e:
        print(f"遍历目录时发生错误: {e}")

# 调用函数并传入目录路径
traverse_and_delete_jpg_files('/var/www/html/images/')
```

## 总结
大家一定要增加服务器权限的管控，开放到外部的接口，一定要做好访问权限控制，避免被有心人利用。