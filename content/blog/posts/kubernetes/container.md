---
title: "Kubernetes — containerd 安装和部署"
date: 2023-12-13
draft: false
featuredImage: "/images/containerd-horizontal-color.png"
featuredImagePreview: "/images/containerd-horizontal-color.png"
images: ["/images/containerd-horizontal-color.png"]
authors: "jobcher"
tags: ["k8s"]
categories: ["k8s"]
series: ["k8s系列"]
---
# containerd
现在很多人说起容器都会说到docker，docker凭借镜像（images）快捷的部署，占领了极大的技术市场，docker公司将自己的核心依赖 Contanerd 捐给了 CNCF，这个就是contanerd的由来，containerd 在kubernetes在 v1.24之后的版本作为底层核心进行使用。  
## Containerd架构
![](/images/container-core.png)  
可以看到 Containerd 仍然采用标准的 `C/S` 架构，服务端通过 `GRPC` 协议提供稳定的 `API`，客户端通过调用服务端的 API 进行高级的操作。
为了解耦，Containerd 将不同的职责划分给不同的组件，每个组件就相当于一个`子系统（subsystem）`。连接不同子系统的组件被称为模块。
总体上 Containerd 被划分为两个子系统：
- **Bundle** : 在 Containerd 中，Bundle 包含了配置、元数据和根文件系统数据，你可以理解为容器的文件系统。而 Bundle 子系统允许用户从镜像中提取和打包 Bundles。
- **Runtime** : Runtime 子系统用来执行 Bundles，比如创建容器。

其中，每一个子系统的行为都由一个或多个模块协作完成（架构图中的 Core 部分）。每一种类型的模块都以插件的形式集成到 Containerd 中，而且插件之间是相互依赖的。例如，上图中的每一个长虚线的方框都表示一种类型的插件，包括 `Service Plugin`、`Metadata Plugin`、`GC Plugin`、`Runtime Plugin` 等，其中 Service Plugin 又会依赖 `Metadata Plugin`、`GC Plugin` 和 `Runtime Plugin`。每一个小方框都表示一个细分的插件，例如 Metadata Plugin 依赖 Containers Plugin、Content Plugin 等。 总之，万物皆插件，插件就是模块，模块就是插件。
### 常用插件
- **Content Plugin** : 提供对镜像中可寻址内容的访问，所有不可变的内容都被存储在这里。
- **Snapshot Plugin** : 用来管理容器镜像的文件系统快照。镜像中的每一个 layer 都会被解压成文件系统快照，类似于 Docker 中的 graphdriver。
- **Metrics** : 暴露各个组件的监控指标。
## 安装
### 卸载docker
首先要保证环境干净整洁，如果你有安装docker服务，需要先卸载docker，如果没有安装可以跳过
```sh
sudo apt-get remove docker docker-engine docker.io containerd runc
```
删除数据
```sh
sudo rm -rf /var/lib/docker
sudo rm -rf /var/lib/containerd
```
### 准备包环境
更换软件源，可以参考这篇文章 [GNU/Linux 一键更换系统软件源脚本](https://test.jobcher.com/linux-mirror/)  
```sh
bash <(curl -sSL https://test.jobcher.com/ChangeMirrors.sh)
```
更新apt，允许使用https
```sh
sudo apt-get updates
```
```sh
 sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
```
添加docker官方GPG key
```sh
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
```
设置软件仓库源
```sh
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```
### 安装containerd
```sh
# 安装containerd
sudo apt-get update
sudo apt-get install -y containerd.io

# 如果是安装docker则执行：
sudo apt-get install docker-ce docker-ce-cli containerd.io

# 查看运行状态
systemctl enable containerd
systemctl status containerd
```
### 检查安装
```sh
ctr version
```
## 命令
containerd 相比于docker , 多了namespace概念, 每个image和container 都会在各自的namespace下可见, 目前k8s会使用k8s.io 作为命名空间
### 查看 images
```sh
ctr image list
# 缩写
ctr i list
ctr i ls
```
### images 打tag
```sh
ctr -n k8s.io i tag pause:3.2 k8s.gcr.io/pause:3.2
#注意: 若新镜像reference 已存在, 需要先删除新reference, 或者如下方式强制替换
ctr -n k8s.io i tag --force pause:3.2 k8s.gcr.io/pause:3.2
```
### 删除镜像
```sh
ctr -n k8s.io i rm k8s.gcr.io/pause:3.2
```
### 拉取镜像
```sh
ctr -n k8s.io i pull -k k8s.gcr.io/pause:3.2
```
### 推送镜像
```sh
ctr -n k8s.io i push -k k8s.gcr.io/pause:3.2
```
### 导出镜像
```sh
ctr -n k8s.io i export pause.tar k8s.gcr.io/pause:3.2
```
### 导入镜像
```sh
ctr -n k8s.io i import pause.tar
```
### 查看容器相关操作
```sh
ctr c
```
### 运行容器
- –null-io: 将容器内标准输出重定向到/dev/null
- –net-host: 主机网络
- -d: 当task执行后就进行下一步shell命令,如没有选项,则会等待用户输入,并定向到容器内
- –mount 挂载本地目录或文件到容器
- –env 环境变量
```sh
ctr -n k8s.io run --null-io --net-host -d \
–env PASSWORD="123456"
–mount type=bind,src=/etc,dst=/host-etc,options=rbind:rw
```
### 查看日志
```sh
ctr -n k8s.io run --log-uri file:///var/log/xx.log
```
### ctr和docker比较
|Containerd命令	|Docker命令	|描述|
|:----|:----|:----|
|ctr task ls	|docker ps	|查看运行容器|
|ctr image ls	|docker images	|获取image信息|
|ctr image pull pause	|docker pull pause	pull |应该pause镜像|
|ctr image push pause-test	|docker push pause-test	|改名|
|ctr image import pause.tar	|docker load 镜像	|导入本地镜像|
|ctr run -d pause-test pause	|docker run -d --name=pause pause-test	|运行容器|
|ctr image tag pause pause-test	|docker tag pause pause-test	|tag应该pause镜像|
### crictl 命令
```sh
# 通过在配置文件中设置端点 --config=/etc/crictl.yaml
root@k8s-node-0001:~$ cat /etc/crictl.yaml
runtime-endpoint: unix:///run/containerd/containerd.sock
```
列出业务容器状态
```sh
crictl inspect ee20ec2346fc5
```
查看运行中容器
```sh
root@k8s-node-0001:~$ crictl pods
POD ID              CREATED             STATE               NAME                                                     NAMESPACE           ATTEMPT             RUNTIME
b39a7883a433d       10 minutes ago      Ready               canal-server-quark-b477b5d79-ql5l5                       mbz-alpha           0                   (default)
```
打印某个固定pod
```sh
root@k8s-node-0001:~$ crictl pods --name canal-server-quark-b477b5d79-ql5l5
POD ID              CREATED             STATE               NAME                                 NAMESPACE           ATTEMPT             RUNTIME
b39a7883a433d       12 minutes ago      Ready               canal-server-quark-b477b5d79-ql5l5   mbz-alpha           0                   (default)
```
打印镜像
```sh
root@k8s-node-0001:~$ crictl images
IMAGE                                                          TAG                             IMAGE ID            SIZE
ccr.ccs.tencentyun.com/koderover-public/library-docker         stable-dind                     a6e51fd179fb8       74.6MB
ccr.ccs.tencentyun.com/koderover-public/library-nginx          stable                          588bb5d559c28       51MB
ccr.ccs.tencentyun.com/koderover-public/nsqio-nsq              v1.0.0-compat                   2714222e1b39d       22MB
```
只打印镜像 ID
```sh
root@k8s-node-0001:~$ crictl images -q
sha256:a6e51fd179fb849f4ec6faee318101d32830103f5615215716bd686c56afaea1
sha256:588bb5d559c2813834104ecfca000c9192e795ff3af473431497176b9cb5f2c3
sha256:2714222e1b39d8bd6300da72b0805061cabeca3b24def12ffddf47abd47e2263
```
打印容器清单
```sh
root@k8s-node-0001:~$ crictl ps -a
CONTAINER           IMAGE               CREATED             STATE               NAME                     ATTEMPT             POD ID
ee20ec2346fc5       c769a1937d035       13 minutes ago      Running             canal-server             0                   b39a7883a433d
76226ddb736be       cc0c524d64c18       34 minutes ago      Running             mbz-rescue-manager       0                   2f9d48c49e891
e2a19ff0591b4       eb40a52eb437d       About an hour ago   Running             export                   0                   9844b5ea5fdbc
```
打印正在运行的容器清单
```sh
root@k8s-node-0001:~$ crictl ps
CONTAINER           IMAGE               CREATED             STATE               NAME                   ATTEMPT             POD ID
ee20ec2346fc5       c769a1937d035       13 minutes ago      Running             canal-server           0                   b39a7883a433d
```
容器上执行命令
```sh
root@k8s-node-0001:~$ crictl exec -i -t ee20ec2346fc5 ls
app.sh  bin  canal-server  health.sh  node_exporter  node_exporter-0.18.1.linux-arm64
```
获取容器的所有日志
```sh
root@k8s-node-0001:~$ crictl logs ee20ec2346fc5
DOCKER_DEPLOY_TYPE=VM
==> INIT /alidata/init/02init-sshd.sh
==> EXIT CODE: 0
==> INIT /alidata/init/fix-hosts.py
```
获取最近的 N 行日志
```sh
root@k8s-node-0001:~$ crictl logs --tail=2 ee20ec2346fc5
start canal successful
==> START SUCCESSFUL ...
```
拉取镜像
```sh
crictl pull busybox
```

## 配置
### 生成配置文件
`Containerd` 的默认配置文件为 `/etc/containerd/config.toml`，我们可以通过命令来生成一个默认的配置：
```sh
mkdir /etc/containerd
containerd config default > /etc/containerd/config.toml
```
### 镜像加速
由于某些不可描述的因素，在国内拉取公共镜像仓库的速度是极慢的，为了节约拉取时间，需要为Containerd 配置镜像仓库的 mirror。Containerd 的镜像仓库 mirror 与 Docker 相比有两个区别：  
- `Containerd` 只支持通过 CRI 拉取镜像的 mirror，也就是说，只有通过 crictl 或者 Kubernetes 调用时 mirror 才会生效，通过 ctr 拉取是不会生效的。
- `Docker` 只支持为 Docker Hub 配置 mirror，而 Containerd 支持为任意镜像仓库配置 mirror。
#### Containerd 的配置结构
```conf
[plugins]
  [plugins."io.containerd.gc.v1.scheduler"]
    pause_threshold = 0.02
    deletion_threshold = 0
    mutation_threshold = 100
    schedule_delay = "0s"
    startup_delay = "100ms"
  [plugins."io.containerd.grpc.v1.cri"]
    disable_tcp_service = true
    stream_server_address = "127.0.0.1"
    stream_server_port = "0"
    stream_idle_timeout = "4h0m0s"
    enable_selinux = false
    sandbox_image = "k8s.gcr.io/pause:3.1"
    stats_collect_period = 10
    systemd_cgroup = false
    enable_tls_streaming = false
    max_container_log_line_size = 16384
    disable_cgroup = false
    disable_apparmor = false
    restrict_oom_score_adj = false
    max_concurrent_downloads = 3
    disable_proc_mount = false
    [plugins."io.containerd.grpc.v1.cri".containerd]
      snapshotter = "overlayfs"
      default_runtime_name = "runc"
      no_pivot = false
      [plugins."io.containerd.grpc.v1.cri".containerd.default_runtime]
        runtime_type = ""
        runtime_engine = ""
        runtime_root = ""
        privileged_without_host_devices = false
      [plugins."io.containerd.grpc.v1.cri".containerd.untrusted_workload_runtime]
        runtime_type = ""
        runtime_engine = ""
        runtime_root = ""
        privileged_without_host_devices = false
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes]
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
          runtime_type = "io.containerd.runc.v1"
          runtime_engine = ""
          runtime_root = ""
          privileged_without_host_devices = false
    [plugins."io.containerd.grpc.v1.cri".cni]
      bin_dir = "/opt/cni/bin"
      conf_dir = "/etc/cni/net.d"
      max_conf_num = 1
      conf_template = ""
    [plugins."io.containerd.grpc.v1.cri".registry]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
          endpoint = ["https://registry-1.docker.io"]
    [plugins."io.containerd.grpc.v1.cri".x509_key_pair_streaming]
      tls_cert_file = ""
      tls_key_file = ""
  [plugins."io.containerd.internal.v1.opt"]
    path = "/opt/containerd"
  [plugins."io.containerd.internal.v1.restart"]
    interval = "10s"
  [plugins."io.containerd.metadata.v1.bolt"]
    content_sharing_policy = "shared"
  [plugins."io.containerd.monitor.v1.cgroups"]
    no_prometheus = false
  [plugins."io.containerd.runtime.v1.linux"]
    shim = "containerd-shim"
    runtime = "runc"
    runtime_root = ""
    no_shim = false
    shim_debug = false
  [plugins."io.containerd.runtime.v2.task"]
    platforms = ["linux/amd64"]
  [plugins."io.containerd.service.v1.diff-service"]
    default = ["walking"]
  [plugins."io.containerd.snapshotter.v1.devmapper"]
    root_path = ""
    pool_name = ""
    base_image_size = ""
```
每一个顶级配置块的命名都是 `plugins."io.containerd.xxx.vx.xxx"` 这种形式，其实每一个顶级配置块都代表一个插件，其中 `io.containerd.xxx.vx` 表示插件的类型，`vx` 后面的 `xxx` 表示插件的 `ID`。可以通过 `ctr` 一览无余：
```sh
ctr plugin ls
```
镜像加速的配置就在 cri 插件配置块下面的 registry 配置块，所以需要修改的部分如下：
```conf
    [plugins."io.containerd.grpc.v1.cri".registry]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
          endpoint = ["https://dockerhub.mirrors.nwafu.edu.cn"]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."k8s.gcr.io"]
          endpoint = ["https://registry.aliyuncs.com/k8sxio"]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."gcr.io"]
          endpoint = ["xxx"]
```
- `registry.mirrors.“xxx”` : 表示需要配置 mirror 的镜像仓库。例如，registry.mirrors."docker.io" 表示配置 docker.io 的 mirror。
- `endpoint` : 表示提供 mirror 的镜像加速服务。例如，这里推荐使用西北农林科技大学提供的镜像加速服务作为 docker.io 的 mirror。
>gcr.io 暂时没有公共加速服务，各位可以自己搭建或者找个公共服务
### 存储配置
Containerd 有两个不同的存储路径，一个用来保存持久化数据，一个用来保存运行时状态。
```sh
root = "/var/lib/containerd"
state = "/run/containerd"
```
- root用来保存持久化数据，包括 Snapshots, Content, Metadata 以及各种插件的数据。每一个插件都有自己单独的目录，Containerd 本身不存储任何数据，它的所有功能都来自于已加载的插件，真是太机智了。
- state 用来保存临时数据，包括 sockets、pid、挂载点、运行时状态以及不需要持久化保存的插件数据。
### OOM配置
```sh
oom_score = 0
```
Containerd 是容器的守护者，一旦发生内存不足的情况，理想的情况应该是先杀死容器，而不是杀死 Containerd。所以需要调整 Containerd 的 `OOM 权重`，减少其被 `OOM Kill` 的几率。最好是将 oom_score 的值调整为比其他守护进程略低的值。这里的 oom_socre 其实对应的是 `/proc/<pid>/oom_socre_adj`，在早期的 Linux 内核版本里使用 `oom_adj` 来调整权重, 后来改用 `oom_socre_adj` 了  
在计算最终的 `badness score` 时，会在计算结果是中加上 `oom_score_adj` ,这样用户就可以通过该在值来保护某个进程不被杀死或者每次都杀某个进程。其取值范围为 `-1000` 到 `1000`。  
如果将该值设置为 `-1000`，则进程永远不会被杀死，因为此时 `badness score` 永远返回0。  
建议 Containerd 将该值设置为`-999` 到 `0`之间。如果作为 Kubernetes 的 Worker 节点，可以考虑设置为 `-999`
## nerdctl+buildkitd编译docker镜像
Containerd不存在编译镜像的命令，因此我们要找其他工具代替，我们使用nerdctl代替镜像编译  
### 下载
[代码地址](https://github.com/containerd/nerdctl/releases)  
这里我们下载完整版(完整版不仅有 netdctl 命令，还包含了 buildkitd、 buildctl、ctr、runc 等 containerd 相关的命令,如果嫌臃肿可以单独下载安装buildkitd)，amd版，各位可以根据需求下载   
```sh
wget https://github.jobcher.com/gh/https://github.com/containerd/nerdctl/releases/download/v1.7.2/nerdctl-full-1.7.2-linux-amd64.tar.gz
```
```sh
tar -zxvf nerdctl-full-1.7.2-linux-amd64.tar.gz
mv nerdctl /usr/bin/
```
### 配置
完整版的 lib 目录下有现成的 `buildkit.service` 文件，不过需要注意 `buildkitd` 命令的路径，文件内默认的路径是 `/usr/local/bin/buildkitd`，需要把二进制文件放到指定路径下，或者修改文件的默认路径  
```sh
cp lib/systemd/system/buildkit.service /lib/systemd/system/
```
>执行 nerdctl build 需要保证 buildctl 命令在系统 PATH 环境变量中可查  

### 启动 buildkit
```sh
systemctl enable buildkit.service --now
```
Dockerfile 测试
```Dockerfile
FROM alpine:3.16.3

ENV LANG=en_US.UTF-8
ENV TZ="Asia/Shanghai"

RUN echo '/bin/sleep 315360000' > start.sh
CMD ["sh","start.sh"]
```
构建
```sh
nerdctl build -t alpine:3.16.3-test .
```
### 登录到私有harbor
```sh
nerdctl login harbor.dev.com
```
### 推送镜像
```sh
nerdctl push harbor.dev.com/test/alpine:3.16.3-test
```