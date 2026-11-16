---
title: "How to learn MLOps, in practice"
date: 2026-11-16
draft: false
authors: "jobcher"
tags: ["MLOps", "MLflow", "Kubernetes"]
categories: ["AI"]
series: ["Artificial Intelligence"]
featuredImage: "/images/mlops.png"
featuredImagePreview: "/images/mlops.png"
images: ["/images/mlops.png"]
description: "A staged path into MLOps from Linux, Docker, Kubernetes, and CI/CD: a minimal machine-learning foundation, reproducible experiments, training pipelines, online serving, monitoring, and rollback."
---

MLOps is the work of making a model trainable by someone else, blocked by evaluation, deployed, observed, and reversible to the previous version when it goes wrong.

If you already know Linux, Docker, Kubernetes, and CI/CD, what is missing is wiring data, experiments, and models into that delivery path. The stages below start from there. If your background is algorithms or data analysis, replace stage 0 with the engineering gap (Git, tests, Docker, one scheduler). The later stages stay the same.

## How MLOps maps onto the ops you already do

A model path usually looks like this:

data in → validation and features → training → offline evaluation → model registry → deploy → online inference → watch drift and outcomes → decide whether to retrain

| What you already do | What you also have to manage for a model |
| --- | --- |
| Code versions | Code, data, and model versions that line up |
| Building an image | The image also pins the model artifact and the feature transformer |
| Unit tests | Add data checks and offline evaluation, and fail the pipeline when they miss the bar |
| Release and rollback | Move the model alias from `staging` to `production`; rollback points the alias at the previous version |
| CPU, error rate, latency | Also watch input distribution, prediction distribution, and business metrics that arrive late |

A complete platform is still a broken path if you cannot answer which data, which code, and which image produced the model that is serving traffic.

## Stage 0: the smallest machine-learning set (2–4 weeks)

The goal is to read training code, read metrics, and decide whether two experiments can be compared. You do not need to start with papers.

What to learn:

- How to split training, validation, and test. Information from the validation set must not leak into features or preprocessing.
- Common metrics for classification and regression: precision, recall, F1, AUC, MAE, RMSE. When classes are extremely imbalanced, accuracy is almost useless.
- When overfitting, underfitting, regularization, and early stopping show up.
- Missing values, categorical encoding, and numeric scaling. A transformer fitted at training time must be used unchanged at inference time.
- One tabular model with scikit-learn or XGBoost: read data, train, print metrics, save the model to a file.
- A command-line script as the training entry point. Fix the random seed. Put parameters in a config file.

Done when: you pick a public tabular dataset (UCI Adult is fine) and write `train.py`. The command line takes a data path and a config. When it finishes, you have metrics and a model file.

Stay with tabular models in this stage. Backprop derivations, paper reading, and large-model fine-tuning come after the path actually runs.

## Stage 1: one experiment that someone else can reproduce (3–4 weeks)

After someone clones the repo, they can rerun the same experiment. That is the first brick of MLOps.

Use Python, Git, MLflow, and DVC.

- Use `venv` or `uv`, and lock dependency versions.
- Do not commit data or model files. DVC commits only the pointer file. Large files live in remote storage.
- Each training run records parameters, metrics, and the model with MLflow Tracking.
- One experiment is one YAML config.
- Cover feature transforms with pytest. Leave a small-sample smoke test on the training script so the command can finish.

These commands should be stable when you check the stage off:

```bash
dvc pull
python train.py --config conf/baseline.yaml
mlflow ui
```

Run again with a different parameter set. The two MLflow runs can be compared side by side. On another machine, `git clone` followed by `dvc pull` reproduces the metrics.

## Stage 2: turn the script into a pipeline (3–4 weeks)

A pipeline is something that can be scheduled, retried, and rerun from an intermediate step. Split it in this order:

1. Validate raw data: field types, null ratios, value ranges.
2. Split the dataset, and version the split on its own.
3. Fit the feature transform, and save the transformer with the model.
4. Train.
5. Evaluate on the held-out set. If the score is below a threshold written down in advance, the process exits non-zero.
6. Only after evaluation passes, register the model and assign the `staging` alias.

Pick one orchestrator. To move fast, chain the steps with a Makefile or a short Python script. To match an existing scheduler, use Airflow or Prefect. Kubeflow Pipelines fits a team that already runs Kubernetes and is ready to standardize a training platform. It is a poor first orchestration course.

Hook CI to the GitLab CI or GitHub Actions you already have. Every push runs lint, pytest, and the small-sample smoke test. A full training run happens on a tag or a manual trigger.

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
    - python src/evaluate_gate.py   # exit non-zero below the threshold
```

Register models with MLflow Model Registry aliases, and load them as `models:/churn@production`. Promotion points the `production` alias at the new version. Rollback points it back at the old one. An alias is harder to get wrong than a version number written into the deploy config.

Done when: changing one line of feature code still lets CI finish; when evaluation is below the threshold, the pipeline fails and the `production` alias does not move.

## Stage 3: a service something else can call (3–4 weeks)

Start with online inference. The feedback is fastest. Offline batch scoring is another entry point to the same model. Add it after the interface is stable.

The service does three things: load the model and the transformer at startup, transform features per request, and return a prediction. Load at process start, not on every request.

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

Then add the things production always runs into:

- Fix request and response fields with Pydantic. Reject extra fields and missing fields immediately.
- Put `model_version` and a request id on the logs. Otherwise you cannot tell which version answered.
- `/health` means the model has finished loading. Hand `/metrics` to Prometheus.
- Pack it as an image. `docker run` locally should be enough to get a prediction.
- FastAPI is enough for a tabular model. For large-model inference on GPUs, look at NVIDIA Triton.

Deploy the way you already deploy on Kubernetes: one Deployment and one Service. Leave KServe until the number of models justifies a shared inference plane.

Done when: the same sample gets the same prediction from the offline script and from the API. The image tag, the model alias, and the transformer file all refer to the same release.

## Stage 4: monitoring and rollback (3–4 weeks)

After launch, a model often gets worse while the process stays healthy. The predictions just drift. Watch four layers:

- Service: QPS, latency, error rate. Prometheus and Grafana, the same stack as the other services.
- Input: missing-feature rate, mean, quantiles. Compare with the training set or with last week.
- Output: the distribution of prediction scores, and the share of the positive class. If the shape jumps, check upstream data first.
- Outcome: when labels arrive late, wait until they are in before computing the real metric. If there are no labels, use a business proxy such as rejection rate or manual-review rate.

Drift can start as a daily job. Compare today's request sample with the training sample using PSI or Evidently, and send the result to a daily report or a metric. The first bar is one comparison every day.

Release in increasing order of risk:

1. Replace all traffic, and keep the previous Deployment. If it goes wrong, roll the `production` alias and the image back together.
2. Shadow traffic: the new model only records predictions. It does not answer the user.
3. Split traffic by percentage: 1% → 10% → 100%. At each step, watch latency and the output distribution.

Automatic retraining comes after the gate, the data window, and rollback are all stable. If the trigger is vague, bad data walks straight into the next model.

Done when: you deliberately change the distribution of one feature, and a report or dashboard points at it within a day. After you roll back, the model version in the live logs is the previous one.

## Stage 5: a platform, after there are many models

One person or a small group only needs one path through stages 1 to 4. The items below show up when the number of models and the number of people grow:

- Training and online serving share the same feature computation, so one sample does not produce two different results. The first step is shipping the transformer with the model. A feature store such as Feast comes in when that convention no longer holds.
- GPU queues and quotas on the cluster. If you already have Kubernetes, attach Kueue or Volcano to the cluster you have.
- Lineage: a production model traces back to a data version, a Git commit, and an image digest.
- Permission: who can move an alias to `production`, and an audit record of that action.
- If the direction is large language models, add LLMOps: versioned prompts and eval sets, RAG index updates, and cost and latency per request. Underneath, it is still versions, gates, deployment, and monitoring.

## Stick to one stack

Keep one tool per step, and learn it well enough to debug it.

| Step | Learn this | Add something else when |
| --- | --- | --- |
| Experiments and registry | MLflow | The team has already standardized on another platform |
| Data versions | DVC | Data is in a lake and you need branches and permissions; then look at lakeFS |
| Pipelines | Airflow or Prefect | Training jobs already all run on Kubernetes; then evaluate Kubeflow |
| Serving | FastAPI + Docker | Large-model inference on GPUs; then Triton |
| Runtime | Kubernetes + the CI you have | Many models; then KServe |
| Monitoring | Prometheus + Grafana; drift with PSI or Evidently | A daily report is no longer enough; then look at a dedicated model-monitoring product |
| Cloud | Pick one of SageMaker, Vertex AI, or Azure Machine Learning | After the open-source path runs, see which piece the managed service removes |

Pick the cloud platform that matches the job. Run the path on your own machine first, then look at the managed service. A model clicked together in a console is hard to debug, because you no longer know what you deployed.

## Tie it together with one project

The topic can be small. Churn, ticket classification, or credit default all work. The repo ends up with:

```text
conf/baseline.yaml          # one config per experiment
src/train.py
src/features.py
src/serve.py
src/evaluate_gate.py        # fail below the threshold
tests/
Dockerfile
.gitlab-ci.yml
monitoring/drift.py         # daily input comparison
artifacts/*.dvc             # pointers to data and models, not the files
README.md
```

The README covers three things: how to reproduce a given experiment, how to promote a model to `production`, and how to roll it back.

On a resume this can read: built the train, register, serve, and rollback path for a tabular model; managed the live version with MLflow aliases; blocked release in CI when evaluation missed the bar; ran a daily check for input drift.

## How to schedule the time

At 8–10 hours a week:

- Month 1: stages 0 and 1. You have a training repo someone else can reproduce.
- Month 2: stages 2 and 3. A pipeline, a registry, and a service that can be called.
- Month 3: stage 4. Monitoring, drift, rollback, and a finished README.
- After that: keep the path alive at work and in this repo, then add one cloud platform or GPU inference for the role you want.

If Kubernetes and CI are already familiar, stage 3 goes much faster. Most of the time goes to stages 0 and 1. Skip those and the later platform has nothing to host.

## Four things you should be able to explain in an interview

A tool list does not help much. Prepare examples around "the model broke; how did you notice, and what did you do":

- A training run no longer matches a historical metric. How you found that the data split or a dependency version had changed.
- Live quality drops. How you tell apart a service failure, input drift, and a change in the labels themselves.
- Evaluation fails. Why that model never reaches `production`.
- On rollback, whether the image, the model alias, the feature transformer, and the config all go back together.

## What to read, in order

1. Chip Huyen, *Designing Machine Learning Systems*. Read the data, training, deployment, and monitoring chapters first, and use the book as a map.
2. [MLOps Zoomcamp](https://github.com/DataTalksClub/mlops-zoomcamp). Free. The assignments are a path you can actually run.
3. DeepLearning.AI on Coursera, [Machine Learning Engineering for Production](https://www.coursera.org/specializations/machine-learning-engineering-for-production-mlops). Use it to line up the vocabulary.
4. Tracking and Model Registry in the [MLflow docs](https://mlflow.org/docs/latest/), plus [Get Started with DVC](https://dvc.org/doc/start). Do the steps against your own repo.
5. Google's [Practitioners Guide to MLOps](https://cloud.google.com/architecture/practitioners-guide-to-mlops). See how a mature team stages the process.

Papers and framework release notes can wait until this path already runs.

## Ways the path gets longer than it needs to be

- Spending the time on a Kubernetes certification while the repo still has no reproducible training entry point. The cluster is where the model runs.
- Learning three tools each for experiment tracking, registry, and deployment, and ending with none that can finish a release alone.
- Notebooks only. Training that CI cannot invoke cannot enter a pipeline.
- Building automatic retraining and a feature platform before versions, an evaluation gate, and rollback exist. Automation amplifies the mistake.
- Reporting training-set metrics only, and never checking that the same sample gets the same prediction from the live API.

First make one model retrainable, rejectable, observable, and rollback-able. Once those four exist, a platform and a team standard have something to land on.
