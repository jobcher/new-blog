---
title: "windows-exporter 监控"
date: 2022-10-08
draft: false
author: "jobcher"
tags: ["skywalking"]
categories: ["监控"]
series: ["运维监控系列"]
---

# windows-exporter 监控安装

### [windows_exporter](https://github.com/prometheus-community/windows_exporter)

### [下载安装](https://github.com/prometheus-community/windows_exporter/releases/download/v0.20.0/windows_exporter-0.20.0-amd64.msi)

## 启动

下载 msi 版本，输入一下命令启动

```sh
msiexec /i C:\Users\Administrator\Downloads\windows_exporter.msi ENABLED_COLLECTORS="ad,iis,logon,memory,process,tcp,scheduled_task" TEXTFILE_DIR="C:\custom_metrics\"
```

卸载

```sh
msiexec /uninstall C:\Users\Administrator\Downloads\windows_exporter.msi
```

## 添加 prometheus 监控

prometheus.yaml

```yml
# 新增 windows-exporter
- job_name: "windows-exporter"
  file_sd_configs:
    - files:
        - "./file_sd/windows-exporter.yaml"
```

./file_sd/windows-exporter.yaml

```yml
# 新增 windows-exporter
- targets: ["192.168.0.6:9182"]
  labels:
    instance: windows-task
```

## 添加 alertmanager 告警

```yml
# 告警信息
groups:
  - name: sanjiang windows 任务计划程序告警
    rules:
      - alert: windows实例任务告警
        expr: windows_scheduled_task_state{state="disabled",task=~"/ETL_kettle_tasks/.*"}==1
        for: 30s
        labels:
          severity: critical
          target: "{{$labels.job}}"
        annotations:
          summary: "sanjiang: {{$labels.job}} windows 任务异常"
          description: "{{$labels.task}} of job {{$labels.job}} 该任务断联已超过1分钟"
```
