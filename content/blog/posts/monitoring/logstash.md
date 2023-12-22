---
title: "logstash 多管道部署"
date: 2022-07-19
draft: false
authors: "jobcher"
tags: ["logstash"]
categories: ["监控"]
series: ["运维监控系列"]
---

# logstash 多管道部署

找到 logstash 目录位置，一般来说在 `/etc/logstash` 路径下,修改 `logstash.yml`

```yaml
#增加 日志记录
path.logs: /var/log/logstash
```

## 增加管道

增加 `conf.d`目录下 test.conf

```conf
input {
    beats {
        host => "0.0.0.0"
        port => 23000 # 修改端口IP
    }
}
filter {
    mutate{
        add_field => {
            "cluster" => "test" # 修改标签
            "job" => "logstash"
        }
    }
}
output {
        file {
            path => "/data/路径名称" # 路径名称
            gzip => false  #匹配以空格开头的行
        }
}

```

修改 `pipelines.yml`

```yml
- pipeline.id: 名称
  path.config: "/etc/logstash/conf.d/配置文件.conf"
  queue.type: persisted
```

## 启动 logstash 文件

```sh
/usr/share/logstash/bin/logstash &
```
