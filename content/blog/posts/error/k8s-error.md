---
title: "K8S 问题排查：cgroup 内存泄露问题"
date: 2022-08-30
draft: false
authors: "jobcher"
featuredImage: '/images/error.jpeg'
featuredImagePreview: '/images/error.jpeg'
images: ['/images/error.jpeg']
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---

# K8S 问题排查：cgroup 内存泄露问题

```log
unable to ensure pod container exists: failed to create container for [kubepods besteffort pod5f26dae8-0421-4eab-a3f7-aa51c6848e2b] : mkdir /sys/fs/cgroup/memory/kubepods/besteffort/pod5f26dae8-0421-4eab-a3f7-aa51c6848e2b: cannot allocate memory
```

## 查看 linux 内核

```sh
cat /proc/version
uname -a
```

可以发现 linux 版本是 3.0 版本

## 原因

> `cgroup` 的 `kmem account` 特性在 `Linux 3.x` 内核上有内存泄露问题，然后`k8s`用了这个特性，导致后面创建不出新的`pod`来了

## 解决方法

```sh
# 修改/etc/default/grub 为
GRUB_CMDLINE_LINUX="crashkernel=auto rd.lvm.lv=centos/root rd.lvm.lv=centos/swap rhgb quiet cgroup.memory=nokmem"
#加上了 cgroup.memory=nokmem
# 生成配置
/usr/sbin/grub2-mkconfig -o /boot/grub2/grub.cfg

# 重启机器
reboot
```

## 验证

```sh
cat /sys/fs/cgroup/memory/kubepods/burstable/pod*/*/memory.kmem.slabinfo
```

输出信息

```sh
cat: /sys/fs/cgroup/memory/kubepods/burstable/pod0fe273ca-42e0-4223-9fe8-16d8dd1774e9/0fdd5d9c16929fd600dbdf313b5c3ebabad912dc0cb076ed6e7799e028b31481/memory.kmem.slabinfo: 输入/输出错误
cat: /sys/fs/cgroup/memory/kubepods/burstable/pod0fe273ca-42e0-4223-9fe8-16d8dd1774e9/aa30198d0c5413b70bf488c9daa350a85c7fc6b677235c5adaf2dde6caf95ec4/memory.kmem.slabinfo: 输入/输出错误
cat: /sys/fs/cgroup/memory/kubepods/burstable/pod5be86c5d-d012-4cc2-b693-4882a15eda90/059b26b00f4f286b0f52e759b83dad79c7676e1705ee0f3f175a277fd1e5ea5a/memory.kmem.slabinfo: 输入/输出错误
cat: /sys/fs/cgroup/memory/kubepods/burstable/pod5be86c5d-d012-4cc2-b693-4882a15eda90/bfa9db0c23fd056a0c05ee5b2b377dd551451cc0f18ddd5db82f9693674a4677/memory.kmem.slabinfo: 输入/输出错误
cat: /sys/fs/cgroup/memory/kubepods/burstable/pod8f588e3d-fe89-4716-ab36-3ef606c70367/6fab9f4f7a83bf4c79a68277b214807bd566a8f13212a0fdb5742e4eee4d75d5/memory.kmem.slabinfo: 输入/输出错误
cat: /sys/fs/cgroup/memory/kubepods/burstable/pod8f588e3d-fe89-4716-ab36-3ef606c70367/b04594732f7e38a47500ffe1150705110cfa683b585aa7eaf0965cc48ba2a46d/memory.kmem.slabinfo: 输入/输出错误
cat: /sys/fs/cgroup/memory/kubepods/burstable/pod9c449bdf-492b-4adc-a623-4ced323ac6d4/c9e73e56ddae0c6c301e43852a51165419eb293b05fa65407f8cb3fe449daf5d/memory.kmem.slabinfo: 输入/输出错误
cat: /sys/fs/cgroup/memory/kubepods/burstable/pod9c449bdf-492b-4adc-a623-4ced323ac6d4/d7392b1e4be1728fd739dc7117e6efc723a4727f143b91a1b386ad35dc1d3a2e/memory.kmem.slabinfo: 输入/输出错误
cat: /sys/fs/cgroup/memory/kubepods/burstable/podd316acd7-69fe-4ad5-963a-6e19174b7cb0/0ec18ac0509e6ab454ebe637bde002330afc9eb70eff6f23fe8caa12880e82f6/memory.kmem.slabinfo: 输入/输出错误
cat: /sys/fs/cgroup/memory/kubepods/burstable/podd316acd7-69fe-4ad5-963a-6e19174b7cb0/dc5be82c01802c349c9505375c37dc054898b9b84e57cb4e671044e8a6459aac/memory.kmem.slabinfo: 输入/输出错误

```
