---
title: "mysql 笔记（1）"
date: 2021-12-12
draft: false
authors: "jobcher"
tags: ["mysql"]
categories: ["数据库"]
series: ["mysql系列"]
---

# mysql 学习笔记（1）

本文章不涉及到关于 mysql 开放上的问题，主要记录关于 mysql 出现的问题，以及如何去维护 mysql 数据的日常。

## mysql 各类信息的收集

1. 收集变量信息

```sql
show global variables;
```

2. 收集进程信息

```sql
show PROCESSLIST;
```

3. 收集错误日志

```sql
show global variables like 'log_error';
```

4. 收集慢日志信息

```sql
show global variables like 'slow_querry_log_file';
```

5. 收集锁信息，高峰时期运行三次，每次间隔 10s

```sql
SELECT  locked_table,
        locked_index,
        locked_type,
        blocking_pid,
        T2.USER  blocking_user,
        T2.HOST blocking_host,
        blocking_lock_mode,
        blocking_trx_rows_modified,
        waiting_pid,
        T3.USER waiting_user,
        T3.HOST waiting_host,
        waiting_lock_mode,
        waiting_trx_row_modified,
        wait_age_secs,
        waiting_query
FROM sys.x$innodb_lock_waits T1
LEFT JOIN INFROMATION_SCHEMA.processlist T2 ON T1.blocking_pid=T2.ID
LEFT JOIN INFROMATION_SCHEMA.processlist T3 ON T3.ID=T1.waiting_pid;

```

6. 收集 mysql 状态信息

```sql
show global status;
show engine innodb status;
show engine innodb mutex;
```

## mysql 基础语法

1. 连接数据库

```sql
mysql -u <用户名> -p
```

2. 创建数据库

```sql
CREATE DATABASE <数据库名称>;
```

3. 删除数据库

```sql
drop database <数据库名称>;
```

4. 选择数据库

```sql
use <数据库名称>;
```

5. 创建表

```sql
CREATE table <数据表名> (
    <字段名1> <数据类型> [约束条件],
    <字段名2> <数据类型> [约束条件],
    <字段名3> <数据类型> [约束条件]
)
#例如
CREATE TABLE IF NOT EXISTS `nbtyfood_tbl`(
   `nbtyfood_id` INT UNSIGNED AUTO_INCREMENT,
   `nbtyfood_title` VARCHAR(100) NOT NULL,
   `nbtyfood_author` VARCHAR(40) NOT NULL,
   `submission_date` DATE,
   PRIMARY KEY ( `nbtyfood_id` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

6. 删除表

```sql
DROP TABLE <数据表名>;
```

7. 插入数据

```sql
INSERT INTO table_name ( field1, field2,...fieldN )
VALUES
value1, value2,...valueN );
```

8. 更新数据

```sql
UPDATE <数据表名> SET <字段名1>='更新' WHERE <字段名2>=3;
```

9. 删除数据

```sql
DELETE FROM <数据表名> WHERE <字段名2>=3;
```

10. 查询数据

```sql
select * from <数据表名>;
```
