---
title: "Linux 系统收包流程以及内核参数优化"
date: 2023-08-30
draft: false
author: "jobcher"
featuredImage: '/images/networking-dc52a683.svg'
featuredImagePreview: '/images/networking-dc52a683.svg'
images: ['/images/nginx.png']
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---
## 简介
高并发的系统架构中，任何细微调整，稍有不注意便会引起连锁反应，只有系统地了解整个网络栈，在处理疑难杂症或者系统优化工作中，才能做到手中有粮心中不慌。在本节，我们概览一个 Linux 系统收包的流程，以便了解高并发系统所面临的性能瓶颈问题以及相关的优化策略。  
### 收包过程
![收包过程](/images/networking-dc52a683.svg)  
1. 网卡 eth0 收到数据包。
2. 网卡通过 DMA 将数据包拷贝到内存的环形缓冲区(Ring Buffer，在网卡中有 RX Ring 和 TX Ring 两种缓冲)。
3. 数据从网卡拷贝到内存后, 网卡产生 IRQ（Interupt ReQuest，硬件中断）告知内核有新的数据包达到。
4. 内核收到中断后, 调用相应中断处理函数，开始唤醒 ksoftirqd 内核线程处理软中断。
5. 内核进行软中断处理，调用 NAPI poll 接口来获取内存环形缓冲区(ring buffer)的数据包，送至更上层处理。
6. 内核中网络协议栈：L2 处理。
7. 内核中网络协议栈：L3 处理。
8. 内核中网络协议栈：L4 处理。
9. 网络协议栈处理数据后，并将其发送到对应应用的 socket 接收缓冲区。

### 高并发瓶颈
- 用户进程调用系统调用陷入内核态的开销。
- CPU 响应包的硬中断 CPU 开销
- ksoftirqd 内核线程的软中断上下文开销。

## RX/TX Ring 优化
处理一个数据包会有各类的中断、softirq 等处理，因为分配给 Ring Buffer 的空间是有限的，当收到的数据包速率大于单个 CPU 处理速度的时，Ring Buffer 可能被占满并导致新数据包被自动丢弃。一个 CPU 去处理 Ring Buffer 数据会很低效，这个时候就产生 RSS、RPS 等多核并发机制来提升内核网络包的处理能力。  
  
但是注意，`开启多核并发特性`，会挤压业务代码的执行时间，`如果业务属于 CPU 密集型`，会导致业务`性能下降`。是否开启多核处理，需要根据业务场景考虑，根据笔者的经验来看，例如此类`负载均衡服务器`、`网关`、`集群核心转发节点`等`网络I/O 密集型场景`可以尝试`优化 RSS、RPS 等配置`。
  
### 1. 判断是否需进行优化
当类似 LVS、集群核心交换节点、负载均衡服务器的场景 PPS（Packet Per Second，包每秒）指标存在一定的优化空间且这些核心节点影响了集群业务，那我们可以查看系统状态以决定是否进行内核优化。  
首先我们确定`是否存在丢包状况`，如果存在则进行相应的调整。`查询网卡收包情况` (RX 为收到的数据、TX 为发送的数据)。  
```sh
ifconfig eth0 | grep -E 'RX|TX'
```
> `eth0` 是你具体服务器的网卡地址  
  
出现以下结果说明RX dropped 表示数据包已经进入了 Ring Buffer，但是由于内存不够等系统原因，导致在拷贝到内存的过程中被丢弃，`RX overruns` 为 Ring Buffer 传输的 IO 大于 kernel 能够处理的 IO 导致，`overruns` 的错误意味着 CPU 无法即使的处理中断是造成 `Ring Buffer 溢出`。
```sh
RX packets 490423734  bytes 193802774970 (180.4 GiB)
RX errors 12732344  dropped 9008921  overruns 3723423  frame 0
TX packets 515280693  bytes 140609362555 (130.9 GiB)
TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
```
### 2. RSS 下的多队列调整
RSS（receive side steering）利用网卡多队列特性，将每个核分别跟网卡的一个首发队列绑定，以达到网卡硬中断和软中断均衡的负载在各个 CPU 中，RSS 要求网卡必须要支持多队列特性。（注意：对于大部分驱动，`修改以下配置会使网卡先 down 再 up，因此会造成丢包`！）  
查询 RX/TX 队列配置和使用情况。  
```sh
ethtool -l eth0
```
```sh
Channel parameters for eth0:
Pre-set maximums:
RX:		0
TX:		0
Other:		0
Combined:	8
Current hardware settings:
RX:		0
TX:		0
Other:		0
Combined:	4
```
可以看到硬件最多支持 6 个，当前使用了 4 个。将 RX 和 TX queue 数量都设为 8。
```sh
ethtool -L eth0 combined 8
```

### 3. 队列大小调整
增大 ring buffer 可以在 PPS（packets per second）很大时`缓解丢包`问题  
查看队列大小。  
```sh
ethtool -g eth0
```
```sh
Ring parameters for eth0:
Pre-set maximums:
RX:		1024
RX Mini:	0
RX Jumbo:	0
TX:		1024
Current hardware settings:
RX:		512
RX Mini:	0
RX Jumbo:	0
TX:		512
```
以上输出显示网卡最多支持 1024 个 RX/TX 数据包大小，但是现在只用到了 512 个。 ethtool -G 修改 queue 大小。  
```sh
ethtool -G eth0 rx 1024
```

## 内核协议栈优化
一个传输少量数据的 TCP 短连接的生命周期中，握手、挥手阶段将近占用了 `70%` 的资源消耗，在一个高并发的服务场景，如负载均衡、数据库等，针对性的去优化较为保守内核参数提提升服务性能的必要手段。在本文中，内核协议栈参数优化方向为 TCP 握手流程中队列大小、挥手 TIME_WAITE、keepalive 保活机制以及拥塞控制。  
![TCP 握手](/images/TCP-ad017126.svg)  
### TCP 握手流程参数优化
握手流程中有两个队列较为关键，当队列满时，多余的连接将会被丢弃。  
- `SYN Queue` 也被称为`半连接队列`，是内核保持的未被 ACK 的 SYN 包最大队列长度，通过内核参数 `net.ipv4.tcp_max_syn_backlog` 配置，高并发的环境下建议设置为 `1024 或更高`。
- `Accept Queue` 也被称为`全连接队列`， 是一个 socket 上等待应用程序 accept 的最大队列长度。取值为 min(backlog，net.core.somaxconn)。
### TCP 连接保活参数优化
当 TCP 建立连接后，会有个发送一个空的 ACK 的探测行为来保持连接（keepalive）。keepalive 受以下参数影响：  
- `net.ipv4.tcp_keepalive_time` 最大闲置时间
- `net.ipv4.tcp_keepalive_intvl` 发送探测包的时间间隔
- `net.ipv4.tcp_keepalive_probes` 最大失败次数，超过此值后将通知应用层连接失效
在大规模的集群内部，如果 `keepalive_time` 设置`较短且发送较为频繁`，会产生大量的`空 ACK 报文`，存在塞满 RingBuffer 造成 TCP 丢包甚至连接断开问题。  
### TCP 连接断开参数优化
由于 `TCP 双全工`的特性，安全关闭一个连接需要四次挥手。但在一个复杂的网络环境下，会存在很多异常情况，异常断开连接会导致产生孤儿连(半连接)。 这种连接既不能发送数据，也无法接收数据，累计过多，会消耗大量系统资源。在高并发的场景下，孤儿连过多会引起资源不足，产生 `Address already in use: connect `类似的错误。  
![TCP 挥手](/images/tcp_disconnect-4051b998.svg)  
挥手流程中的主要优化为 `TIME_WAIT` 的参数调整。`TIME_WAIT` 是 TCP 挥手的最后一个状态。当收到被动方发来的 FIN 报文后，主动方回复 ACK，表示确认对方的发送通道已经关闭，继而进入TIME_WAIT 状态 ，等待 2MSL 时间后关闭连接。  
如果发起连接一方的 `TIME_WAIT` 状态过多，会占满了所有端口资源，则会导致无法创建新连接。`TIME_WAIT` 的问题在`反向代理中比较明显`，例如 nginx 默认行为下会对于 client 传来的每一个 request 都向 upstream server 打开一个新连接，`高 QPS 的反向代理`将会快速积累 TIME_WAIT 状态的 socket，直到没有可用的本地端口，无法继续向 upstream 打开连接，此时服务将不可用。  
- `net.ipv4.tcp_max_tw_buckets`，此数值定义系统在同一时间最多能有多少 TIME_WAIT 状态，当超过这个值时，系统会`直接删掉这个 socket 而不会留下 TIME_WAIT 的状态`
- `net.ipv4.ip_local_port_range`，TCP 建立连接时 client 会随机从该参数中定义的端口范围中选择一个作为源端口。可以`调整该参数范围增大可选择的端口范围`。

### 相关配置参考
```sh
net.ipv4.tcp_tw_recycle = 0
net.ipv4.tcp_tw_reuse = 1
net.ipv4.ip_local_port_range = 1024 65535
net.ipv4.tcp_rmem = 16384 262144 8388608
net.ipv4.tcp_wmem = 32768 524288 16777216
net.core.somaxconn = 8192
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.core.wmem_default = 2097152
net.ipv4.tcp_max_tw_buckets = 5000
net.ipv4.tcp_max_syn_backlog = 10240
net.core.netdev_max_backlog = 10240
net.netfilter.nf_conntrack_max = 1000000
net.ipv4.netfilter.ip_conntrack_tcp_timeout_established = 7200
net.core.default_qdisc = fq_codel
net.ipv4.tcp_congestion_control = bbr
net.ipv4.tcp_slow_start_after_idle = 0
```








