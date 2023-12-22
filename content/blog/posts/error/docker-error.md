---
title: "docker 问题处理"
date: 2022-08-04
draft: false
author: "jobcher"
tags: ["error"]
categories: ["问题库"]
series: ["问题库系列"]
---

### docker 无法启动

打开服务器输入`docker ps`,输出错误  
Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?

怀疑是不是`docker.services` 部署没成功，`systemctl start docker` 启动 docker，结果服务器还是报错

Job for docker.service failed because the control process exited with error code.  
See "systemctl status docker.service" and "journalctl -xe" for details.

`systemctl status docker.service` 输出日志：

```sh
● docker.service - Docker Application Container Engine
     Loaded: loaded (/lib/systemd/system/docker.service; enabled; vendor preset: enabled)
     Active: failed (Result: exit-code) since Thu 2022-08-04 11:43:05 CST; 2min 57s ago
TriggeredBy: ● docker.socket
       Docs: https://docs.docker.com
    Process: 30432 ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock (code=exited, status=1/FAILURE)
   Main PID: 30432 (code=exited, status=1/FAILURE)

Aug 04 11:43:05 master01 systemd[1]: docker.service: Scheduled restart job, restart counter is at 3.
Aug 04 11:43:05 master01 systemd[1]: Stopped Docker Application Container Engine.
Aug 04 11:43:05 master01 systemd[1]: docker.service: Start request repeated too quickly.
Aug 04 11:43:05 master01 systemd[1]: docker.service: Failed with result 'exit-code'.
Aug 04 11:43:05 master01 systemd[1]: Failed to start Docker Application Container Engine.
```

`journalctl -xe` 输出日志：

```sh
Aug 04 11:46:49 master01 systemd[1]: Starting Docker Socket for the API.
-- Subject: A start job for unit docker.socket has begun execution
-- Defined-By: systemd
-- Support: http://www.ubuntu.com/support
--
-- A start job for unit docker.socket has begun execution.
--
-- The job identifier is 58900.
Aug 04 11:46:49 master01 systemd[1]: Listening on Docker Socket for the API.
-- Subject: A start job for unit docker.socket has finished successfully
-- Defined-By: systemd
-- Support: http://www.ubuntu.com/support
--
-- A start job for unit docker.socket has finished successfully.
--
-- The job identifier is 58900.
Aug 04 11:46:49 master01 systemd[1]: Starting Docker Application Container Engine...
-- Subject: A start job for unit docker.service has begun execution
-- Defined-By: systemd
-- Support: http://www.ubuntu.com/support
--
-- A start job for unit docker.service has begun execution.
--
-- The job identifier is 58830.
Aug 04 11:46:49 master01 dockerd[30544]: unable to configure the Docker daemon with file /etc/docker/daemon.json: EOF
Aug 04 11:46:49 master01 systemd[1]: docker.service: Main process exited, code=exited, status=1/FAILURE
-- Subject: Unit process exited
-- Defined-By: systemd
-- Support: http://www.ubuntu.com/support
--
-- An ExecStart= process belonging to unit docker.service has exited.
--
-- The process' exit code is 'exited' and its exit status is 1.
Aug 04 11:46:49 master01 systemd[1]: docker.service: Failed with result 'exit-code'.
-- Subject: Unit failed
-- Defined-By: systemd
-- Support: http://www.ubuntu.com/support
--
-- The unit docker.service has entered the 'failed' state with result 'exit-code'.
Aug 04 11:46:49 master01 systemd[1]: Failed to start Docker Application Container Engine.
-- Subject: A start job for unit docker.service has failed
-- Defined-By: systemd
-- Support: http://www.ubuntu.com/support
--
-- A start job for unit docker.service has finished with a failure.
--
-- The job identifier is 58830 and the job result is failed.
Aug 04 11:46:49 master01 sudo[30535]: pam_unix(sudo:session): session closed for user root
```

我运行了:

> sudo dockerd --debug

输出日志：  
unable to configure the Docker daemon with file /etc/docker/daemon.json: EOF

- 解决方法：

```sh
#如果 /etc/docker/daemon.json 为空
vim /etc/docker/daemon.json
# 添加
{
}
#保存退出
:wq
```

重启 docker

> systemctl restart docker
