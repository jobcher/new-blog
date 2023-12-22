---
title: "githubAction set-output弃用错误"
date: 2022-10-21
draft: false
authors: "jobcher"
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---

# githubAction set-output 弃用错误

The `set-output` command is deprecated and will be disabled soon. Please upgrade to using Environment Files. For more information see: https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/

## 原因

如果您有一个使用 设置输出的`GitHub Actionsecho ::set-output key=value`工作流程，您已经开始看到无用的弃用警告。这是修复它的方法。查看官方链接基本上得不到什么帮助！

## 修复方法

1. 更新其它人的 action 方法

```sh
将 @actions/core 提升到 1.10.0
```

2. 修改自己的 aciton 方法

```sh
run: echo "::set-output name=KEY::VALUE"
## 改为
run: echo "KEY=VALUE" >>$GITHUB_OUTPUT
```

> 建议：使用自己的方法

## 总结

平台经营者非常肆意妄为的修改自己的代码内容弃用功能，无限的权力滋生傲慢……我相信大部分开发这并没有注意到这个告警，知道流水线服务报错之后才会注意到，希望微软可以对能更加包容不同的开发者，尊重开发者社区。
