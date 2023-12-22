---
title: "Gitlab批量导出用户"
date: 2022-01-14
draft: false
authors: "jobcher"
tags: ["自动化"]
categories: ["gitlab"]
series: ["git入门系列"]
---

# Gitlab 批量导出用户

登陆 Gitlab 服务器进行数据库登陆、数据查询及信息导出操作。

## 操作步骤

1. 根据配置文件，定位数据库相关信息

```sh
cat /var/opt/gitlab/gitlab-rails/etc/database.yml
```

![exportuser1](/images/exportuser1.png)

2. 查看 Gitlab 对应的系统用户

```sh
cat /etc/passwd | grep gitlab
```

![exportuser2](/images/exportuser2.png)

3. 切换用户 gitlab-psql

```sh
su - gitlab-psql
```

4. 登陆数据库（-h 指定 host，-d 指定数据库） 使用第 1 步获取的信息

```sh
psql -h /var/opt/gitlab/postgresql -d gitlabhq_production
```

(1) 查看帮助信息

```sh
gitlabhq_production=# \h
```

(2) 查看数据库

```sh
gitlabhq_production=# \l
```

(3) 查看库中的表（执行命令后，按回车键显示更多表信息）

```sh
gitlabhq_production=# \dt
```

(4) 通过筛查，可在库中找到 users 表，相关用户信息都记录在表中！

```sh
gitlabhq_production=# \d users
```

(5) 查看表信息

```sh
gitlabhq_production=# SELECT * FROM users;
```

(6) 查看 users 表中的 name 字段

```sh
gitlabhq_production=# SELECT name FROM users;
```

(7)登出数据库

```sh
gitlabhq_production=# \q
```

5. 根据需要提取的信息，确定表 users 中的字段，进行导出操作

```sh
echo 'select name,username,email,state from users;' |psql -h /var/opt/gitlab/postgresql -d gitlabhq_production > userinfo.txt
```

存储在/var/opt/gitlab/postgresql/userinfo.txt
