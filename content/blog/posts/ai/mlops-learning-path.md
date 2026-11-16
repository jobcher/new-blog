---
title: "往 MLOps 发展，具体怎么学"
date: 2026-11-16
draft: false
authors: "jobcher"
tags: ["MLOps", "MLflow", "Kubernetes"]
categories: ["AI"]
series: ["人工智能系列"]
featuredImage: "/images/mlops.png"
featuredImagePreview: "/images/mlops.png"
images: ["/images/mlops.png"]
description: "从已有的 Linux、Docker、Kubernetes 和 CI/CD 出发，按阶段学习 MLOps：最小机器学习基础、可复现实验、训练流水线、在线服务、监控与回滚。"
---

MLOps 是一个典型的“交叉学科”岗位。它不像纯算法工程师那样只钻研模型结构，也不像纯后端开发那样只钻研业务逻辑，它要求你**既要懂模型原理（知道模型是怎么练出来的），又要懂工程架构（知道模型是怎么跑起来的）**。

如果你想往 MLOps 方向发展，你的学习路径应该是 **“由点及面，从工程向下渗透”**。我建议你分五个阶段进行，这不仅是技能的积累，更是思维方式的转变。

---

### 第一阶段：扎实的工程基础（MLOps 的“骨架”）
很多算法工程师转 MLOps 的最大痛点是“工程能力弱”。在 MLOps 领域，**代码的质量、规范性和部署能力**比你写出一个 SOTA（最先进）的模型更重要。

1.  **Python 高级编程：**
    *   不仅是写简单的脚本，要学习面向对象（OOP）、装饰器、生成器、异常处理、类型提示 (Type Hinting)。
    *   **核心目标：** 能够写出生产级别的、可维护的 Python 代码，而不是一堆乱糟糟的 Jupyter Notebook。
2.  **Linux 与 Shell 编程：**
    *   掌握常用的 Linux 命令、文件权限、进程管理、环境变量。
    *   **核心目标：** 能够熟练在远程服务器（通常是 Linux）上进行环境配置和任务调度。
3.  **Docker 容器化（极其重要）：**
    *   理解镜像 (Image)、容器 (Container)、Dockerfile、数据卷 (Volume) 的概念。
    *   **核心目标：** 学会如何将“环境+代码+依赖”封装进一个 Docker 镜像中，实现“一次构建，到处运行”。

### 第二阶段：数据工程与模型生命周期（MLOps 的“血液”）
MLOps 的核心是处理数据的流动。

1.  **数据处理与存储：**
    *   **SQL：** 必须精通，因为数据大多存储在关系型数据库中。
    *   **数据流处理：** 了解 Pandas (离线) 和简单的流处理思想。
    *   **特征工程：** 理解特征是如何从原始数据中提取出来的，以及如何保证**训练和推理时特征的一致性**。
2.  **模型训练管理：**
    *   **实验追踪 (Experiment Tracking)：** 学会使用 **MLflow** 或 **Weights & Biases (W&B)**。
    *   **核心目标：** 解决“我上次跑这个实验用的参数是什么？”以及“哪个版本的模型效果最好？”这两个问题。
3.  **模型版本控制：**
    *   理解模型不仅是权重文件，还包括对应的代码版本和数据版本（了解 **DVC - Data Version Control**）。

### 第三阶段：自动化与流水线（MLOps 的“灵魂”）
这是区分“调包侠”和“MLOps 工程师”的分水岭：**自动化 (Automation)**。

1.  **CI/CD (持续集成/持续部署)：**
    *   学习使用 **GitHub Actions** 或 **GitLab CI**。
    *   **核心目标：** 实现“当你提交代码后，系统自动运行测试、构建镜像、并部署到测试环境”。
2.  **工作流编排 (Orchestration)：**
    *   学习 **Apache Airflow** (通用任务调度) 或 **Kubeflow Pipelines** (专门为机器学习设计的编排)。
    *   **核心目标：** 把“拉取数据 → 数据清洗 → 模型训练 → 模型评估 → 模型部署”这个长长的链条，串成一个自动化的流水线。

### 第四阶段：大规模部署与监控（MLOps 的“外壳”）
当模型需要服务于成千上万的用户时，你需要处理并发和稳定性。

1.  **模型服务化 (Model Serving)：**
    *   学习如何用 **FastAPI** 写高性能的 API 接口。
    *   学习专业的推理引擎：**NVIDIA Triton** 或 **TF Serving**。
2.  **Kubernetes (K8s) —— 终极挑战：**
    *   K8s 是现代云原生架构的基石。你需要理解 Pod、Service、Deployment、ConfigMap 的概念。
    *   **核心目标：** 能够实现模型的**弹性伸缩**（用户多了自动增加模型实例，用户少了自动减少以省钱）。
3.  **监控与告警 (Monitoring)：**
    *   学习 **Prometheus** (指标收集) 和 **Grafana** (数据可视化)。
    *   **核心目标：** 建立一个仪表盘，实时观察模型的响应延迟、内存占用以及最重要的——**数据漂移 (Data Drift)**。

### 第五阶段：进阶——大规模分布式与云原生

*   **分布式训练：** 了解如何使用多卡/多机进行大规模训练（PyTorch Distributed, DeepSpeed）。
*   **云平台实践：** 深入学习 AWS (SageMaker), GCP (Vertex AI) 或 Azure ML 的一套完整方案。

---

### 🛠 总结：你的“学习清单”优先级

如果时间有限，请按以下优先级进阶：

1.  **入门期 (必须掌握):** Python 规范化 → Docker → MLflow → FastAPI。
2.  **进阶期 (职业分水岭):** Git/CI-CD → SQL/特征存储 → Airflow → K8s。
3.  **专家期 (高薪领域):** 漂移检测 → 分布式训练 → K8s 高级调度 → 完善的 MLOps 平台搭建。

**💡 一个练习建议：**
不要只是看书！尝试做一个项目：
*   找一个开源数据集（比如 MNIST 或房价预测）。
*   写一个 Python 脚本进行训练。
*   用 **MLflow** 记录实验。
*   用 **Docker** 把模型打包。
*   用 **FastAPI** 把模型做成 Web 服务。
*   最后，试着把整个过程写成一个 **GitHub Action** 脚本，实现“改一下代码，模型就自动重新训练并部署”的效果。

---
MLOps 做的是这件事：让一个模型可以被别人重复训练、被评估拦住、被部署、被观察，出了问题能回到上一版。

如果你已经会 Linux、Docker、Kubernetes 和 CI/CD，缺的是把数据、实验和模型接进这条交付链路。下面按这个起点排。算法或数据分析背景的人，把阶段 0 换成工程补课（Git、测试、Docker、一个调度器），后面的阶段相同。

## MLOps 和现有运维的对应关系

一条模型链路通常是：

数据进入 → 校验和特征 → 训练 → 离线评估 → 写入模型注册表 → 部署 → 在线推理 → 看漂移和效果 → 决定是否重训

| 你已经在做的 | 到了模型这里要多管的 |
| --- | --- |
| 代码版本 | 代码、数据、模型三个版本对得上 |
| 构建镜像 | 镜像里还要绑定模型制品和特征变换器 |
| 单元测试 | 再加数据校验和离线评估，不达标就让流水线失败 |
| 发版、回滚 | 模型别名从 `staging` 切到 `production`，回滚是把别名指回上一版 |
| CPU、错误率、延迟 | 再加输入分布、预测分布、延迟到达的业务指标 |

平台搭得再完整，对不上「哪份数据、哪次代码、哪个镜像产出了线上这个模型」，这条链路就是断的。

## 阶段 0：机器学习最小集（2–4 周）

目标是看懂训练代码、读懂指标、判断两次实验能不能比。不需要先读论文。

要掌握的内容：

- 训练、验证、测试怎么切。验证集的信息不能漏进特征或预处理。
- 分类和回归的常用指标：Precision、Recall、F1、AUC、MAE、RMSE。类别极不平衡时，Accuracy 几乎没有参考价值。
- 过拟合、欠拟合、正则、早停分别在什么现象下出现。
- 缺失值、类别编码、数值缩放。训练时拟合出来的变换器，推理时必须原样使用。
- 用 scikit-learn 或 XGBoost 走通一个表格模型：读数据、训练、打印指标、把模型存成文件。
- 训练入口是命令行脚本。随机种子固定，参数放在配置文件里。

验收标准：选一个公开表格数据（例如 UCI Adult），写出 `train.py`。命令行传入数据路径和配置，运行结束后得到指标和模型文件。

这一阶段先把表格模型做扎实。反向传播推导、论文阅读、大模型微调都放到链路跑通之后。

## 阶段 1：一次实验可以复现（3–4 周）

别人克隆仓库之后，能跑出同一次实验。这是 MLOps 的第一块砖。

推荐组合：Python、Git、MLflow、DVC。

- 环境用 `venv` 或 `uv`，依赖锁版本。
- 数据和模型文件不进 Git。DVC 只把指针文件交进去，大文件放远端存储。
- 每次训练用 MLflow Tracking 记下参数、指标和模型。
- 一次实验对应一份 YAML 配置。
- 用 pytest 覆盖特征变换。训练脚本留一个小样本冒烟测试，保证命令能跑完。

验收时这几条命令的结果是稳定的：

```bash
dvc pull
python train.py --config conf/baseline.yaml
mlflow ui
```

换一组参数再跑一次，MLflow 里两次 run 可以并排比较。换一台机器，`git clone` 之后执行 `dvc pull`，指标能够复现。

## 阶段 2：把脚本收成流水线（3–4 周）

能被调度、重试、在中间步骤重跑的，才是流水线。按顺序拆开：

1. 校验原始数据：字段类型、空值比例、取值范围。
2. 切分数据集，切分结果单独版本化。
3. 拟合特征变换，把变换器和模型一起保存。
4. 训练。
5. 在留出集上评估。低于事先写死的阈值，进程以非 0 退出。
6. 通过评估之后，才把模型注册，并打上 `staging` 别名。

编排器选一个。想快，用 Makefile 或一段 Python 把步骤串起来。要和现有调度对齐，用 Airflow 或 Prefect。Kubeflow Pipelines 适合已经有 Kubernetes、而且团队准备统一训练平台的时候，不适合当作第一门编排课。

CI 用现有的 GitLab CI 或 GitHub Actions 接上。每次 push 跑 lint、pytest 和小样本冒烟；打 tag 或手动触发才跑完整训练。

```yaml
stages: [test, train]

test:
  stage: test
  script:
    - pip install -r requirements.txt
    - pytest

train:
  stage: train
  when: manual
  script:
    - dvc pull
    - python src/train.py --config conf/baseline.yaml
    - python src/evaluate_gate.py   # 低于阈值则非 0 退出
```

模型注册用 MLflow Model Registry 的别名，加载时写 `models:/churn@production`。晋升就是把 `production` 别名指到新版本，回滚就是指回旧版本。别名比把版本号写进部署配置更不容易对错。

验收：改一行特征代码，CI 能跑完；评估低于阈值时，流水线失败，`production` 别名不动。

## 阶段 3：做成可以调用的服务（3–4 周）

先做在线推理，反馈最快。离线批推是同一套模型的另一种入口，等接口稳定后再加。

服务只做三件事：启动时加载模型和变换器，按请求做特征变换，返回预测。加载放在进程启动，不要放在每次请求里。

```python
from contextlib import asynccontextmanager

import joblib
import mlflow.pyfunc
import pandas as pd
from fastapi import FastAPI
from pydantic import BaseModel

MODEL_URI = "models:/churn@production"


class PredictRequest(BaseModel):
    tenure: int
    monthly_charges: float
    contract: str

    def to_frame(self) -> pd.DataFrame:
        return pd.DataFrame([self.model_dump()])


class State:
    model = None
    transformer = None


state = State()


@asynccontextmanager
async def lifespan(app: FastAPI):
    state.model = mlflow.pyfunc.load_model(MODEL_URI)
    state.transformer = joblib.load("artifacts/transformer.joblib")
    yield


app = FastAPI(lifespan=lifespan)


@app.post("/predict")
def predict(req: PredictRequest):
    features = state.transformer.transform(req.to_frame())
    proba = state.model.predict(features)
    return {"proba": float(proba[0]), "model": MODEL_URI}
```

接着补上这几件线上一定会遇到的事：

- 请求和响应用 Pydantic 固定字段。多出来的字段、缺字段要直接拒绝。
- 日志带上 `model_version` 和请求 id，否则出了问题对不上是哪一版在答。
- `/health` 表示模型已经加载好。`/metrics` 交给 Prometheus。
- 打包成镜像，本地 `docker run` 就能预测。
- 表格模型用 FastAPI 足够。GPU 上的大模型推理再考虑 NVIDIA Triton。

部署沿用现有的 Kubernetes：一个 Deployment 加一个 Service。KServe 留给模型数量上来、需要统一推理平面的时候。

验收：同一条样本，离线脚本的预测和接口返回一致。镜像标签、模型别名、变换器文件三者能对上。

## 阶段 4：监控和回滚（3–4 周）

模型上线之后变差，经常是预测悄悄漂了，进程还是健康的。分四层看：

- 服务：QPS、延迟、错误率。Prometheus + Grafana，和别的服务同一套。
- 输入：特征缺失率、均值、分位数。和训练集或上一周比。
- 输出：预测分数的分布、正类占比。分布突然变形，先查上游数据。
- 效果：标签会延迟到达时，等到齐了再算真实指标。没有标签时，用一个业务代理指标，例如拒绝率、人工复核率。

漂移可以先写一个每日任务，用 PSI 或 Evidently 对比「今天的请求样本」和「训练样本」，把结果打进日报或指标。先做到每天有一份对比。

发布按风险递进：

1. 全量替换，保留上一版 Deployment，出问题就把 `production` 别名和镜像一起退回。
2. 影子流量：新模型只记录预测，不返回给用户。
3. 按比例切流：1% → 10% → 100%，每一步看延迟和输出分布。

自动重训放在门禁、数据窗口和回滚都稳定之后。触发条件没写清楚时，坏数据会直接进入下一版模型。

验收：人为把某个特征的分布改掉，日报或面板能在一天内指出来。执行回滚后，线上日志里的模型版本回到上一版。

## 阶段 5：模型变多之后再做平台

一个人或一个小组，阶段 1 到 4 的一条链路就够用。下面这些是模型数量和人数上来之后的事：

- 训练和在线共用同一份特征计算，避免同一条样本两边结果不同。第一步是变换器跟模型一起发布。特征商店（例如 Feast）等这套约定守不住了再引入。
- 集群上的 GPU 队列和配额。已有 Kubernetes 的话，用 Kueue 或 Volcano 接上现有集群。
- 血缘：Production 模型能追溯到数据版本、Git 提交和镜像 digest。
- 权限：谁可以把别名改到 `production`，操作要有审计。
- 如果方向是大模型，再加 LLMOps：提示词和评测集版本、RAG 索引更新、单次请求的成本和延迟。底层仍然是版本、门禁、部署、监控。

## 固定一条技术栈

每个环节留一个工具，用到能排障。

| 环节 | 就学这个 | 什么时候再加 |
| --- | --- | --- |
| 实验和注册 | MLflow | 团队已经统一用别的平台时再换 |
| 数据版本 | DVC | 数据进湖、要分支和权限时再看 lakeFS |
| 流水线 | Airflow 或 Prefect | 训练任务都已经在 Kubernetes 上，再评估 Kubeflow |
| 服务 | FastAPI + Docker | GPU 大模型推理再上 Triton |
| 运行 | Kubernetes + 现有 CI | 模型数量多了再上 KServe |
| 监控 | Prometheus + Grafana，漂移用 PSI 或 Evidently | 日报看不过来再谈专门的模型监控产品 |
| 云 | SageMaker、Vertex AI、Azure Machine Learning 三选一 | 开源链路跑通之后，对照它替你省掉了哪一段 |

云平台建议跟着工作环境选。先在自己机器上把链路跑通，再看托管服务。直接在控制台点出来的模型，排障时会不知道自己部署了什么。

## 用一个项目串起来

题目不用大。用户流失、工单分类、信用违约都可以。仓库里最终有这些东西：

```text
conf/baseline.yaml          # 一次实验一份配置
src/train.py
src/features.py
src/serve.py
src/evaluate_gate.py        # 低于阈值就失败
tests/
Dockerfile
.gitlab-ci.yml
monitoring/drift.py         # 每日输入对比
artifacts/*.dvc             # 数据和模型的指针，不是文件本身
README.md
```

README 写三件事：怎么复现某一次实验，怎么把模型晋升到 `production`，怎么回滚。

简历上可以写成：搭建表格模型的训练、注册、服务和回滚链路；用 MLflow 别名管理线上版本；评估不达标时 CI 阻断发布；对输入漂移做每日检测。

## 时间怎么排

按每周 8–10 小时估算：

- 第 1 个月：阶段 0 和阶段 1，得到一个能复现的训练仓库。
- 第 2 个月：阶段 2 和阶段 3，流水线、注册表、一个可调用的服务。
- 第 3 个月：阶段 4，监控、漂移、回滚，把 README 补完。
- 之后：在工作和这个仓库里维护这条链路，再按岗位补一个云平台或 GPU 推理。

Kubernetes 和 CI 已经熟的人，阶段 3 会快很多。时间主要花在阶段 0 和阶段 1。这两段跳过去，后面的平台没有东西可以托管。

## 面试时能讲清楚的四件事

工具名单帮助不大。围绕「模型坏了怎么发现、怎么处理」准备例子：

- 一次训练对不上历史指标，怎么查到是数据切分或依赖版本变了。
- 线上效果下降时，怎么区分服务故障、输入漂移和标签本身变了。
- 评估失败时，模型为什么到不了 `production`。
- 回滚时，镜像、模型别名、特征变换器、配置是否一起回去。

## 资料顺序

1. Chip Huyen《Designing Machine Learning Systems》（中文版《设计机器学习系统》）。先读数据、训练、部署、监控几章，把它当地图。
2. [MLOps Zoomcamp](https://github.com/DataTalksClub/mlops-zoomcamp)。免费，作业就是一条能跑的链路。
3. Coursera 上 DeepLearning.AI 的 [Machine Learning Engineering for Production](https://www.coursera.org/specializations/machine-learning-engineering-for-production-mlops)。用来对齐术语。
4. [MLflow 文档](https://mlflow.org/docs/latest/)里的 Tracking 和 Model Registry，加上 [DVC 入门](https://dvc.org/doc/start)。对着自己的仓库做。
5. Google 的 [Practitioners Guide to MLOps](https://cloud.google.com/architecture/practitioners-guide-to-mlops)。看成熟团队怎么给流程分阶段。

论文和框架发布笔记，放到这条链路已经能跑之后再看。

## 路径变长的几种走法

- 把时间花在 Kubernetes 认证上，仓库里却没有一个能复现的训练入口。集群是运行模型的地方。
- 实验跟踪、注册、部署各学三套，结果没有一套能独立完成一次发布。
- 只有 notebook。CI 调不到的训练，就不能进流水线。
- 版本、评估门禁和回滚还没有时，先做自动重训和特征平台。自动化会把错误放大。
- 只报告训练集指标，不去核对同一条样本在线上接口的预测是否一致。

先让一个模型能够被重复训练、被拒绝发布、被观察、被回滚。这四件事做到了，再谈平台和团队规范，才有东西可以落地。
