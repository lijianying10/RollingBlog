title: StackStorm介绍与入门
date: 2017-04-30 17:25:57
categories: 技术
tags: [docker,StackStorm2]
---

## 背景

在日常运维工作中经常使用的脚本为了提供给别人用,为了降低出错概率与操作复杂度，所以这里介绍一下StackStorm2的最基本使用与安装。

在日常开发工作中很多积累下来的知识需要`脚本化`->`自动化`在这一过程中需要做到`知识可自动执行`节省培训时间。

请先通读之后再做实践。

## 介绍

StackStorm2 是一种DevOps工具，包括顺序简单的工作流，包括Mistral工作流，包括触发器等。

在下面这个列表中列出一些笔者认为能够快速直观的学习ST2的官方文档列表。
按照顺序读可以快速掌握使用方法方便接下来继续学习。可以节约很多时间。

1. [功能介绍阅读](https://docs.stackstorm.com/overview.html)
2. [ST2全局关系图,在学习的时候适合来回做参考](https://docs.stackstorm.com/install/overview.html)
3. [st2命令掌握,cli是最强大的,图形界面适合给别人调用你写的功能](https://docs.stackstorm.com/start.html)
4. [开始尝试写action，写完了之后使用`st2 reload --register-all`注册你写的代码](https://docs.stackstorm.com/actions.html)
5. [开始尝试把action结合起来一起运行](https://docs.stackstorm.com/actionchain.html)

## 快速安装用的 Dockerfile

请查看链接[https://github.com/lijianying10/FixLinux/tree/master/st2](https://github.com/lijianying10/FixLinux/tree/master/st2)

构建流程描述：[https://github.com/lijianying10/FixLinux/blob/master/st2/Dockerfile](https://github.com/lijianying10/FixLinux/blob/master/st2/Dockerfile)

``` sh
FROM ubuntu:16.04
COPY st2ctl st2.conf supervisord.conf docker-entrypoint.sh 
RUN apt-get update && \
# 安装基本组件
apt-get install -y build-essential wget gnupg-curl curl sudo apache2-utils vim apt-utils supervisor && \

# 安装 ST2
os=ubuntu dist=xenial curl -s https://packagecloud.io/install/repositories/StackStorm/stable/script.deb.sh | sudo bash && \
apt-get update && \
apt-get install -y st2  && \

# 安装nginx 
apt-key adv --fetch-keys http://nginx.org/keys/nginx_signing.key && \
echo 'deb http://nginx.org/packages/ubuntu/ xenial nginx' >> /etc/apt/sources.list.d/nginx.list && \
apt-get update && \
apt-get install -y st2web nginx && \
rm /etc/nginx/conf.d/default.conf && \
cp /usr/share/doc/st2/conf/nginx/st2.conf /etc/nginx/conf.d/ && \

# 安装kubectl tini
curl -o /bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.6.1/bin/linux/amd64/kubectl && \
chmod +x /bin/kubectl && \
wget https://github.com/krallin/tini/releases/download/v0.14.0/tini-amd64 -O /bin/tini && \
chmod +x /bin/tini && \

# 复制配置文件
mv /st2.conf /etc/st2/ && \
mv supervisord.conf /etc/supervisor/supervisord.conf && \
mv /st2ctl /usr/bin/st2ctl && \
chmod +x /usr/bin/st2ctl && \

# 安装python pip
curl -SsL https://bootstrap.pypa.io/get-pip.py | python && \

# 安装docker client
pip install && \
 wget https://get.docker.com/builds/Linux/x86_64/docker-17.05.0-ce.tgz && \
 tar xf docker-17.05.0-ce.tgz && \
mv docker/docker /bin/docker && \
rm -rf docker docker-17.05.0-ce.tgz

ENTRYPOINT ["/bin/tini", "--"]
CMD bash /docker-entrypoint.sh
```

`注意：` 笔者裁剪掉了chatops以及Mistral原因是在前期刚学习的时候并不能用上。还能节省大量配置时间。

在构建中加入了entrypoint脚本，是启动的时候执行的脚本流程描述如下：

启动流程描述 [https://github.com/lijianying10/FixLinux/blob/master/st2/docker-entrypoint.sh](https://github.com/lijianying10/FixLinux/blob/master/st2/docker-entrypoint.sh)

``` sh
# 生成SSH key
echo generate ssh key
ssh-keygen -f /root/.ssh/id_rsa -P "" && cp /root/.ssh/id_rsa.pub /root/.ssh/authorized_keys

# 生成 ST2 账号密码
echo generate user
printf "%s\n" "${USER_NAME:?Need to set USER_NAME non-empty}"
printf "%s\n" "${USER_PASSWORD:?Need to set USER_PASSWORD non-empty}"
echo $USER_PASSWORD | sudo htpasswd -i /etc/st2/htpasswd $USER_NAME

# 生成 证书
echo generate cert
sudo mkdir -p /etc/ssl/st2
sudo openssl req -x509 -newkey rsa:2048 -keyout /etc/ssl/st2/st2.key -out /etc/ssl/st2/st2.crt \
-days 365 -nodes -subj "/C=US/ST=California/L=Palo Alto/O=StackStorm/OU=Information \
Technology/CN=$(hostname)"

# 检查环境变量是否完备
printf "%s\n" "${CONN_RMQ:?Need to set CONN_RMQ non-empty}"
printf "%s\n" "${MONGO_HOST:?Need to set MONGO_HOST non-empty}"
printf "%s\n" "${MONGO_DB:?Need to set MONGO_DB non-empty}"
printf "%s\n" "${MONGO_PORT:?Need to set MONGO_PORT non-empty}"

# 生成st2配置文件
cat >> /etc/st2/st2.conf <<EOF
[system_user]
user = root
ssh_key_file = /root/.ssh/id_rsa
[messaging]
url = $CONN_RMQ
[ssh_runner]
remote_dir = /tmp
[database]
host = $MONGO_HOST
port = $MONGO_PORT
db_name = $MONGO_DB
EOF

# 启动supervisor
/usr/bin/supervisord -c /etc/supervisor/supervisord.conf
```

## 简单上手

### 容器化一次性安装启动（无负担安装最快上手）

``` sh
docker network create stn
docker run -itd --hostname st2-mongo  --name st2-mongo  -v /var/lib/mongo:/data/db --net=stn daocloud.io/library/mongo:3.4.3
docker run -itd --hostname st2-etcd --name st2-etcd --net=stn index.tenxcloud.com/coreos/etcd:2.3.1 /etcd -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 -advertise-client-urls http://127.0.0.1:2379,http://127.0.0.1:4001
docker run -itd --hostname st2-rabbit --name st2-rabbit -e RABBITMQ_DEFAULT_USER=root -e RABBITMQ_DEFAULT_PASS=123456 --net=stn daocloud.io/library/rabbitmq:3.6.9
docker run -itd --hostname st2 --name st2 --net=stn -e USER_NAME=admin -e USER_PASSWORD=123456 -e CONN_RMQ=amqp://root:123456@st2-rabbit.stn:5672/ -e MONGO_HOST=st2-mongo.stn -e MONGO_DB=st2 -e MONGO_PORT=27017 -e ETCD_ENDPOINT=http://st2-etcd.stn:2379 -p 80:80 -p 443:443 index.tenxcloud.com/philo/stackstorm:2.2.1
```

`注意： ` 笔者使用了daocloud和tenxcloud提供的dockerhub为大家提供服务。

`注意： ` 笔者这持久化了mongodb的存储，如果不需要可删除挂载。

`注意： ` ETCD用来存储全局变量。用作分布式锁。在未来使用ST2调度`kubernetes`中会介绍如何使用。

### 第一个action

经过之前的文档基本概念的了解，相信读者已经对系统概念有基本了解。
Action是任务执行的最小单元。成功的学会Shell 和Python的Action和简单的WorkFlow就能够应付非常多的工作。

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

