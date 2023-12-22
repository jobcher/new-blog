---
title: "mysql数据库备份迁移"
date: 2021-12-24
draft: false
authors: "jobcher"
tags: ["mysql"]
categories: ["数据库"]
series: ["mysql系列"]
---

# mysql 数据库备份迁移

使用 mydumper 做数据备份迁移

## 备份数据库

1. 安装

```sh
# 安装 centos
yum install https://github.com/mydumper/mydumper/releases/download/v0.11.5/mydumper-0.11.5-1.el7.x86_64.rpm
yum install https://github.com/mydumper/mydumper/releases/download/v0.11.5/mydumper-0.11.5-1.el8.x86_64.rpm
```

```sh
# 安装 ubuntu
apt-get install libatomic1
wget https://github.com/mydumper/mydumper/releases/download/v0.11.5/mydumper_0.11.5-1.$(lsb_release -cs)_amd64.deb
dpkg -i mydumper_0.11.5-1.$(lsb_release -cs)_amd64.deb

```

2. 备份

```sh
nohup mydumper -h '备份数据库' \
-u '用户名' \
-p '密码' \
--threads=16 \
-B 备份数据库 \
-v 3 \
--outputdir=./backup --rows=100000 \
-L mydumper-logs.log &
```

## 迁移数据库

1. 还原数据

```sh
nohup myloader -h '迁移数据库' \
-u '用户名' \
-p '密码' \
--directory=./backup \
-s 来源数据库 \
-B 还原数据库 \
-t 16 \
-v 3 \
-e 2>myloader-logs.log &
```

## mydumper/myloader 参数

1. mydumper

```sh
Usage:
  mydumper [OPTION...] multi-threaded MySQL dumping

Help Options:
  -?, --help                  Show help options

Application Options:
  -B, --database              需要备份的数据库，一个数据库一条命令备份，要不就是备份所有数据库，包括mysql。
  -T, --tables-list           需要备份的表，用逗号分隔。
  -o, --outputdir             备份文件目录
  -s, --statement-size        生成插入语句的字节数，默认1000000，这个参数不能太小，不然会报 Row bigger than statement_size for tools.t_serverinfo
  -r, --rows                  试图用行块来分割表，该参数关闭--chunk-filesize
  -F, --chunk-filesize        行块分割表的文件大小，单位是MB
  -c, --compress              压缩输出文件
  -e, --build-empty-files     即使表没有数据，也产生一个空文件
  -x, --regex                 正则表达式匹配，如'db.table'
  -i, --ignore-engines        忽略的存储引擎，用逗号分隔
  -m, --no-schemas            不导出表结构
  -d, --no-data               不导出表数据
  -G, --triggers              导出触发器
  -E, --events                导出事件
  -R, --routines              导出存储过程
  -k, --no-locks              不执行共享读锁 警告：这将导致不一致的备份
  --less-locking              减到最小的锁在innodb表上.
  -l, --long-query-guard      设置长查询时间,默认60秒，超过该时间则会报错：There are queries in PROCESSLIST running longer than 60s, aborting dump
  -K, --kill-long-queries kill掉长时间执行的查询，备份报错：Lock wait timeout exceeded; try restarting transaction
  -D, --daemon                启用守护进程模式
  -I, --snapshot-interval     dump快照间隔时间，默认60s，需要在daemon模式下
  -L, --logfile               使用日志文件，默认标准输出到终端
  --tz-utc                    备份的时候允许备份Timestamp，这样会导致不同时区的备份还原会出问题，默认关闭，参数：--skip-tz-utc to disable.
  --skip-tz-utc
  --use-savepoints            使用savepoints来减少采集metadata所造成的锁时间，需要SUPER权限
  --success-on-1146           Not increment error count and Warning instead of Critical in case of table doesn't exist
  --lock-all-tables           锁全表，代替FLUSH TABLE WITH READ LOCK
  -U, --updated-since         Use Update_time to dump only tables updated in the last U days
  --trx-consistency-only      Transactional consistency only
  -h, --host                  The host to connect to
  -u, --user                  Username with privileges to run the dump
  -p, --password              User password
  -P, --port                  TCP/IP port to connect to
  -S, --socket                UNIX domain socket file to use for connection
  -t, --threads               备份执行的线程数,默认4个线程
  -C, --compress-protocol     在mysql连接上使用压缩协议
  -V, --version               Show the program version and exit
  -v, --verbose               更多输出, 0 = silent, 1 = errors, 2 = warnings, 3 = info, default 2
```

2. myloader

```sh
Usage:
  myloader [OPTION...] multi-threaded MySQL loader

Help Options:
  -?, --help                        Show help options

Application Options:
  -d, --directory                   备份文件所在的目录
  -q, --queries-per-transaction     每个事务的query数量, 默认1000
  -o, --overwrite-tables            如果表存在则先删除，使用该参数，需要备份时候要备份表结构，不然还原会找不到表
  -B, --database                    指定需要还原的数据库
  -s, --source-db                   还原的数据库
  -e, --enable-binlog               启用二进制日志恢复数据
  -h, --host                        The host to connect to
  -u, --user                        Username with privileges to run the dump
  -p, --password                    User password
  -P, --port                        TCP/IP port to connect to
  -S, --socket                      UNIX domain socket file to use for connection
  -t, --threads                     使用的线程数量，默认4
  -C, --compress-protocol           连接上使用压缩协议
  -V, --version                     Show the program version and exit
  -v, --verbose                     更多输出, 0 = silent, 1 = errors, 2 = warnings, 3 = info, default 2
```
