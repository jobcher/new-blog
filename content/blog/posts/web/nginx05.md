---
title: "nginx ssh-key connection exception"
date: 2022-06-15
draft: false
authors: "jobcher"
tags: ["nginx", "linux"]
categories: ["web 服务器"]
series: ["web服务入门系列"]
---

# nginx ssh-key connection exception

Not long ago, I wanted to restart the company's gitlab server.I couldn't coonect to ssh when it restarted.emm……I try copy the ssh rsa.pub,but it didn't work.

error log:

```log
identity_sign: private key ~/.ssh/id_rsa contents do not match public
```

what is happen？

# solution

reconfigure gitlab ssh key!

1. create new ssh key

```sh
ssh-keygen -t rsa -C 'git@gitlab.com' -f ~/.ssh/gitlab-rsa
```

2. update config file,enter ~./ssh,open config

```sh
# add host
Host gitlab.com
    HostName gitlab.com
    IdentityFile ~/.ssh/gitlab_id-rsa

```

3. enter http://gitlab.com ,Profile Settings-->SSH Keys-->Add SSH Key

`You are done`
