---
title: "gitlab与github同步项目"
date: 2021-12-27
draft: false
authors: "jobcher"
tags: ["gitlab"]
categories: ["gitlab"]
series: ["git入门系列"]
---

# gitlab 与 github 同步项目

1. 本地同步项目

```sh
git clone
```

2. 创建一个同名的项目,命令行终端中添加 remote 地址

```sh
git remote add githubOrigin git@github.com:sjtfreaks/blog.git
```

3. 项目同步到 Github 上

```sh
git push -u githubOrigin main
```

4. 分别同步 github 与 gitlab 即可

```sh
git push -u githubOrigin main
git push -u origin main
```
