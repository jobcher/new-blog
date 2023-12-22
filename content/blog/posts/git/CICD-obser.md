---
title: "CI/CD 可观察性-基于grafana"
date: 2023-12-12
draft: false
featuredImage: "/images/grafana-logo.png"
featuredImagePreview: "/images/grafana-logo.png"
images: ["/images/grafana-logo.png"]
author: "jobcher"
tags: ["gitlab"]
categories: ["gitlab"]
series: ["git入门系列"]
---
## 背景
什么是 CI/CD 可观察性，我们如何为更多可观察的管道铺平道路？可观察性不仅仅是观察错误或监控基本健康信号。相反，它会更深入，以便您可以了解系统内行为背后的“原因”。  
CI/CD 可观察性在其中发挥着关键作用。这是关于深入了解持续集成和部署系统的整个管道——查看每个代码签入、每个测试、每个构建和每个部署。当您组合所有这些数据时，您可以全面了解整个软件交付过程，揭示效率领域、瓶颈和潜在故障点。  
CI/CD 可观察性是可观察性的一个子集，专注于软件开发生命周期。它有助于以多种方式确保流程可靠、相关且易于理解：  
- `积极主动解决问题`。没有可观察性，我们只能对问题做出反应。有了它，我们可以在问题升级之前预见并解决问题，从而节省时间和资源。
- `更好的决策`。通过了解 CI/CD 流程的细节，团队可以在资源分配、流程变更和工具采用方面做出更明智的决策。
- `建立信心`。通过对 CI/CD 管道的清晰洞察，开发人员、测试人员和运营团队可以对他们发布的软件更有信心。它减少了“对部署的恐惧”并培育了持续改进的文化。
- `问责制和透明度`。可观察性确保 CI/CD 流程的每一步都是可追溯的。这意味着，如果出现问题，可以追溯到其源头，促进问责并帮助解决根本原因，而不仅仅是解决症状。
## 问题
CI/CD 系统并非没有自身的挑战。破坏 CI/CD 管道平稳运行的常见问题是`不稳定`、`性能下降`和`配置错误`。
### Flakiness 片状
片状测试是 CI/CD 方程中不可预测的变量。当测试在代码没有任何更改的情况下产生不同的结果（通过或失败）时，该测试被认为是“不稳定的”。出现片状现象通常有以下几个原因：  
- 外部依赖和环境问题。如果这些依赖项不能始终可用，则依赖于外部服务、数据库或特定环境设置的测试可能会产生不可预测的结果。如果环境设置不正确或意外拆除，也可能会发生这种情况。从本质上讲，先前测试的残留或外部服务的不可用可能会扭曲结果，使其不可靠。
- 比赛条件。当系统的行为依赖于不可控事件的顺序或时间时，就会出现这种情况。特别是在异步操作中，如果管理不当，事件序列的不可预测性可能会导致偶发故障。
### Performance regression 性能回归
随着 CI/CD 流程的发展并变得更加复杂，系统性能可能会开始下降。这种回归可能不会立即显现出来，但长期的累积效应可能会阻碍 CI/CD 管道的效率。以下是常见的原因：  
- 测试执行效率低下。某些测试可能会运行比必要的时间更长的时间，这可能是因为冗余操作、设置的等待时间太长或查询效率低下。这在集成和端到端测试中尤其明显。
- 代码和测试膨胀。当我们添加更多功能和测试而不解决技术债务或进行优化时，我们的构建时间可能会增加。有些测试从添加之日起可能会很慢。如果不解决这些问题，整个构建和测试过程可能会比需要的时间更长。

### Misconfigurations 配置错误
即使是最深思熟虑的管道也可能因配置错误而失败。这可能导致：
- `次优测试计划`。 CI/CD 管道遵循一条关键路径，其中每个步骤都依赖于前一个步骤。如果步骤未设置为按正确的顺序执行或正在等待非依赖项，则可能会导致效率低下。
- `次优容量规划`。未配置足够的资源或对所需工作负载规划不当可能会导致管道出现瓶颈。如果 CI/CD 流程在关键阶段没有必要的能力，则可能会减慢整个工作流程或导致中断和故障。
## DORA 指标
- `Deployment frequency` 部署频率 (DF)：组织成功发布到生产环境的频率
- `Mean Lead time for changes` 变更平均前置时间 (MLT)：从代码提交到代码在生产中运行所需的时间
- `Mean time to recover` 平均恢复时间 (MTTR)：发生服务事件或缺陷后恢复服务需要多长时间
- `Change failure rate` 变更失败率 (CFR)：导致失败的变更百分比
## 优化 CI/CD 可观察性
目前GraCIe 是 Grafana 正在使用的应用程序插件，旨在为用户提供一种简单的方法来了解他们的 CI/CD 系统。它非常适合评估构建性能、识别测试结果中的不一致以及分析构建输出。该应用程序简化了这些流程，旨在轻松提供有关管道的见解。  
通过利用 [Grafana Tempo](https://grafana.com/oss/tempo/?pg=blog&plcmt=body-txt&src=tw&mdm=social&cnt=youre_probably_familiar_w&camp=blog)、[Grafana Loki](https://grafana.com/oss/loki/?pg=blog&plcmt=body-txt&src=tw&mdm=social&cnt=youre_probably_familiar_w&camp=blog) 和 Prometheus 的强大功能，我们为一个全新的领域（即 CI/CD 可观测性）构建了一种固执己见的体验，因为它最终依赖于更成熟的可观测性用例中使用的相同遥测信号。此外，通过依赖 OpenTelemetry，GraCIe 几乎可以与任何 CI/CD 平台无缝协作，为用户提供相同的无与伦比的见解，而无需自定义设置或配置。  
![cicd-observability-2.png](/images/cicd-observability-2.png)  
![cicd-observability-3.png](/images/cicd-observability-3.png)  
  
  
[转载自grafana](https://grafana.com/blog/2023/11/20/ci-cd-observability-via-opentelemetry-at-grafana-labs/?utm_campaign=blog&utm_content=youre_probably_familiar_w&utm_medium=social&utm_source=tw)  