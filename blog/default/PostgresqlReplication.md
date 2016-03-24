title: PostgreSQL Replication (hot standby) 配置
date: 2016-03-24 14:28:15
categories: 技术
tags: [postgresql]
---

![](http://www.postgresql.org/media/img/layout/hdr_left.png)

## VERSION:1

## PostgreSQL Version 9.5.1

## 背景：

由于数据安全需要，我们对线上服务器的数据库进行主备配置。
在本文中使用的所有配置例子都是为了达到hot_standby 效果而做的。
里面涉及了很多权限放开的问题。请大家参考我的例子的同时要仔细思考一下。
具体的线上服务器应该如何配置。

在学习的时候一定要挑选质量比较高的资料，比如说DO厂的资料一般都比较好，
另外需要重点参考官方文档。虽然在短期内可能由于资料看的比较多学习成本比较高。
但是在后面的研究成本会明显降低。更加节约时间。

## 我当时学习的参考文档：

[DigitamOcean](https://www.digitalocean.com/community/tutorials/how-to-set-up-master-slave-replication-on-postgresql-on-an-ubuntu-12-04-vps)

[PostgreSQL Document Chapter 25](http://www.postgresql.org/docs/9.5/static/high-availability.html)

## 测试环境

在本文描述中使用的测试环境为Docker环境，不使用虚拟机可以让调试节奏更加快速。

### PostgreSQL Docker image 构建    

Dockerfile：

``` shell
FROM ubuntu:14.04.4
RUN sed -i 's/archive.ubuntu/mirrors.aliyun/g' /etc/apt/sources.list && apt-get update && apt-get install -y wget vim telnet
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main 9.5" > /etc/apt/sources.list.d/postgresql.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add - && apt-get update && apt-get install -y postgresql-9.5
```

镜像基于Ubuntu 14.04 ，先安装依赖包： wget vim telnet ，分别是用来下载key，编辑配置文件，测试服务端口是否监听。
下面两个layer是使用官方的方法安装数据库。

### 容器运行方法

``` shell
docker network create aaa
docker run -it -d --name pg2 --net=aaa -v /root/share:/share -h slave psql /bin/bash
docker run -it -d --name pg1 --net=aaa -v /root/share:/share -h master psql /bin/bash
```

### 本文中master节点的配置方法

``` shell
cat >>  /etc/postgresql/9.5/main/postgresql.conf << EOF

listen_addresses = '*'          # what IP address(es) to listen on;
wal_level = hot_standby
max_wal_senders = 5
wal_keep_segments = 64
synchronous_standby_names = '*'
EOF

cat > /etc/postgresql/9.5/main/pg_hba.conf << EOF
local   all             postgres                                trust
local   all             all                                     peer
host    all             all             0.0.0.0/0               md5
host    replication     all             0.0.0.0/0               trust
host    all             all             ::1/128                 md5
EOF
```

需要重点关注的参数：
`synchronous_standby_names` 这个参数里面的设置对应的是 slave recovery.conf 文件中的 `primary_conninfo` 中的 `application_name` 
写法比如说:

``` shell
  primary_conninfo = ' application_name=app1 host=192.168.0.2 user=repluser password=123 port=5432'
```

如果recovery 不不设置这个，这里就需要些成 * 这里使用数组的方式比如说:` 'app1,app2'`

在HBA文件中：

- 第一列 local代表本地unix连接，这里直接设置成trust(无需密码登陆)
- 第二列 数据库名字，但是 replication是需要独立设置的这里注意
- 第三列 用户名 
- 第四列 认证方式 peer获取操作系统用户名 trust 无密码登陆 md5 密码登陆 

这里的IP地址 `0.0.0.0/0` 是allow地址，当前的数据为各种开，不限制。 `/` 后面为掩码长度。


### 本文中slave节点的配置方法

``` shell
cat >>  /etc/postgresql/9.5/main/postgresql.conf << EOF
listen_addresses = '*'          # what IP address(es) to listen on;
wal_level = hot_standby
max_connections = 1000
hot_standby = on
hot_standby_feedback = on
EOF

cat > /etc/postgresql/9.5/main/pg_hba.conf << EOF
local   all             postgres                                trust
local   all             all                                     peer
host    all             all             0.0.0.0/0               md5
host    replication     all             0.0.0.0/0               trust
host    all             all             ::1/128                 md5
EOF

cat >  /var/lib/postgresql/9.5/main/recovery.conf << EOF
standby_mode = on
primary_conninfo = 'host=192.168.0.2 user=repluser password=123 port=5432'
EOF
```

需要注意的点：
这里的conninfo需要自己根据情况修改，host地址为主机的地址，用户建立的方法注意后面文档的写法。

## 具体操作方法

1. Master节点操作：

```
# 启动数据库
service postgresql start
# 修改账号密码
su postgres
psql -c "ALTER USER postgres WITH PASSWORD '123';"
psql -c " CREATE ROLE repluser REPLICATION LOGIN PASSWORD '123';"

# 修改配置
sudo su
sh /share/master.sh

# 重启数据库
service postgresql restart

```

2. 同步数据

```
# 设置数据库进入备份状态
psql -U postgres -c "select pg_start_backup('initial_backup');"

# 远程集群服务器配置方法
# rsync -cva --inplace /var/lib/postgresql/9.5/main/ slave_IP_address:/var/lib/postgresql/9.5/main/

# docker 双容器同步方法
cp -rf /var/lib/postgresql/9.5/main /share/

# 数据库退出备份状态
psql -U postgres -c "select pg_stop_backup();"
```

3. Slave节点操作方法

``` shell

# 配置服务器
sh /share/slave.sh

# 同步数据（如果是rsync 集群内同步就不需要了）
cp -rf /share/main /var/lib/postgresql/9.5

# 修改数据库文件owner
chown -R postgres.postgres main
```

## 总结

整个过程还是比较简单的，难点在于找到好用的材料比较困难。

说明一下同步的意义：当slave同步的内容有master wal 状态为minial产生的数据，就会出现wal设置问题。
很长时间都没明白是为啥，后来才懂。
具体报错内容参考这里[http://stackoverflow.com/questions/9123458/log-shipping-error-postgres](http://stackoverflow.com/questions/9123458/log-shipping-error-postgres)
重点错误特征内容：
```
WARNING:  WAL was generated with wal_level=minimal, data may be missing
HINT:  This happens if you temporarily set wal_level=minimal without taking a new base backup.
FATAL:  hot standby is not possible because wal_level was not set to "hot_standby" on the master server
```

其实最后一个FATAL的解释应该是数据在迁移前(或者您在执行的时候没有迁移)产生的未同步的数据。在master server上产生的时候wal_level 不是 hot_standby
我经过两天的摸索才明白是这个意思。

