title: StackStorm介绍与入门
date: 2017-04-30 17:25:57
categories: 技术
tags: [docker,StackStorm2]
---

## 背景

在日常运维工作中经常使用的脚本为了提供给别人用,为了降低出错概率与操作复杂度，所以这里介绍一下StackStorm2的最基本使用与安装。

在日常开发工作中很多积累下来的知识需要`脚本化`->`自动化`在这一过程中需要做到`知识可自动执行`节省培训时间。

请先通读之后再做实践。

## Dockerfile

请查看链接[https://github.com/lijianying10/FixLinux/tree/master/st2](https://github.com/lijianying10/FixLinux/tree/master/st2)

构建流程描述：[https://github.com/lijianying10/FixLinux/blob/master/st2/Dockerfile](https://github.com/lijianying10/FixLinux/blob/master/st2/Dockerfile)

1. 安装最基本组件（如果你的依赖比较多最好安装build-essential）
2. 安装ST2的官方deb包
3. 安装与配置nginx
4. 安装kubectl命令以及安装[tini](https://github.com/krallin/tini)避免在kubernetes运行中产生僵尸进程
5. 安装配置supervisor 与st2ctl `这些修改了st2默认使用systemctl为supervisor`

`注意：` 笔者裁剪掉了chatops以及Mistral原因是在前期刚学习的时候并不能用上。还能节省大量配置时间。


启动流程描述 [https://github.com/lijianying10/FixLinux/blob/master/st2/docker-entrypoint.sh](https://github.com/lijianying10/FixLinux/blob/master/st2/docker-entrypoint.sh)

1. 生成sshkey
2. 生成 账号密码
3. 自签证书
4. 生成st2配置
5. 启动st2服务堆栈

## 简单上手

### 容器化一次性安装启动（无负担安装最快上手）

``` sh
docker network create stn
docker run -itd --hostname st2-mongo  --name st2-mongo  -v /var/lib/mongo:/data/db --net=stn daocloud.io/library/mongo:3.4.3
docker run -itd --hostname st2-etcd --name st2-etcd --net=stn index.tenxcloud.com/coreos/etcd:2.3.1 /usr/local/bin/etcd -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 -advertise-client-urls http://127.0.0.1:2379,http://127.0.0.1:4001
docker run -itd --hostname st2-rabbit --name st2-rabbit -e RABBITMQ_DEFAULT_USER=root -e RABBITMQ_DEFAULT_PASS=123456 --net=stn daocloud.io/library/rabbitmq:3.6.9
docker run -itd --hostname st2 --name st2 --net=stn -e USER_NAME=admin -e USER_PASSWORD=123456 -e CONN_RMQ=amqp://root:123456@st2-rabbit.stn:5672/ -e MONGO_HOST=st2-mongo.stn -e MONGO_DB=st2 -e MONGO_PORT=27017 -e ETCD_ENDPOINT=http://st2-etcd.stn:2379 -p 80:80 -p 443:443 index.tenxcloud.com/philo/stackstorm:2.2.1
```

`注意： ` 笔者这持久化了mongodb的存储，如果不需要可删除挂载。

`注意： ` ETCD用来存储全局变量。用作分布式锁。在未来使用ST2调度`kubernetes`中会介绍如何使用。


### 最值得参考的文档列表：

在下面这个列表中列出一些笔者认为使用最少时间掌握StackStorm2入门的文档。
按照顺序读可以快速掌握使用方法方便接下来继续学习。可以节约很多时间。

1. [功能介绍阅读](https://docs.stackstorm.com/overview.html)
2. [ST2全局关系图,在学习的时候适合来回做参考](https://docs.stackstorm.com/install/overview.html)
3. [st2命令掌握,cli是最强大的,图形界面适合给别人调用你写的功能](https://docs.stackstorm.com/start.html)
4. [开始尝试写action，写完了之后使用`st2 reload --register-all`注册你写的代码](https://docs.stackstorm.com/actions.html)
5. [开始尝试把action结合起来一起运行](https://docs.stackstorm.com/actionchain.html)

### 第一个action

第一个shell action：

``` sh
#!/usr/bin/env bash

SERVER=$1
MESSAGE=$2
echo ${SERVER} ${MESSAGE}
```

第一个action yaml申明

``` yaml
---
name: "my_first_action"
runner_type: "local-shell-cmd"
description: "first automate"
enabled: true
entry_point: "first.sh"
parameters:
    server:
        type: "string"
        description: "server address"
        required: true
        position: 0
    message:
        type: "string"
        description: "the information"
        required: true
        position: 1
```

`注意： ` 在开始复制文档之前看清楚PACK是如何组织架构开发的。运行后在容器中查看文件夹`/opt/stackstorm/packs/packs`来学习如何开发一个pack是最直观最有效的。


第一个workflow：

`注意： ` workflow 在simple action中只能串行执行，并且只有成功失败两个路径可以走。

有两个文件组成`workflow.yaml` 在这里开发workflow chain。

第二个文件是 `xx.meta.yaml` 申明workflow。

case请看参考文档5。

# 总结

经过查看Dockerimage是如何构建与运行的，快速在你的服务器上运行StackStorm2，到简单编写Action和workflow之后您就已经入门了StackStorm2。
在未来的文章中我会继续更新更高级的应用方法。

