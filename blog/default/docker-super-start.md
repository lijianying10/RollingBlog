title: "docker新手的超级开始"
date: 2015-04-01 10:17:28
tags: docker
categories: 技术
---
## 超级开始
在本文中简短的描述一下如何开始使用docker工作
理论研究以及介绍什么的左上角，google去一抓一大把。
这里只给出一别必要的操作，既然说他是好东西。
我们就拉出来616看看是什么情况。（能动手就别吵吵了）
命令遇到问题就见招拆招。都运行完自己就有自己的体会了。

## 这时候推荐在看一看docker`官方的`用户引导[https://docs.docker.com/userguide/](https://docs.docker.com/userguide/)

## 最后看命令的manual基本上就了解透彻了，再出什么问题就只能见招拆招了。

## 准备环境：
ubuntu 14.04 虚拟机云主机测试服务器都可以
准备账号： daocloud
以内国内网络经常出现错误因此我们需要换一个其他地方来下载docker镜像
在注册完进入到console里面之后，`很明显的位置上`会告诉你怎么用的。
我这里就不说我的账号了。(*^__^*) 嘻嘻……

## 安装docker

```
wget -qO- https://get.docker.com/ | sh
```

正常一路走下去之后啥都不用管了

然后根据daocloud来配置你的image源

比如说你想下载mongodb 以及mysql 这两个服务。
`一定要注意，docker需要root权限，这货是个服务`
那么输入如下命令
```
docker pull mysql:5.6
docker pull mongo:2.6
```
一定要加tags 不然就自动下载latest都是最新的 ，有可能不适用。

有了镜像之后，就可以新建instance来运行你的程序了。
这里直接给出运行模板。方便使用。配置管理。

```
docker run --name mysqlInstance -e MYSQL_ROOT_PASSWORD=1 -v [host绝对路径]:/var/lib/mysql -d -p 3307:3306 mysql:5.6
```

参数解释：
```
run是新运行一个instance
--name 非重要，不然就是一个hash不好操作。也不好分类别名之类的
-e 容器内部的环境参数 这里的参数是设置mysql root 密码的
-v 挂载容器内部的文件夹到外部，但是需要的是绝对路径（具有持久化方面的应用还是需要挂载本地磁盘的，不然删除之后数据就没了，肯定有其他方法导出来或者做管理的。但是我感觉这种方法是最方便的）
-p 端口映射 host:con
image:tag
```

## 查看所有运行的容器：
```
docker ps -a
```
显示了所有信息

## 运行或者关闭一个docker
```
docker start [name]
docker stop [name]
```
## 删除一个容器
```
docker rm [name]
```

访问一个正在运行的容器（调试的时候非常重要）
首先我们不要安装ssh了。没必要因为docker都给你准备好接口了
```
docker exec -it [name] /bin/bash
```
就可以运行一个容器内部的bash了。非常方便的调试你的服务。

## 异常处理
主要是针对内存比较小的情况做一个异常处理

### 首当其冲阿里云
默认路由与docker默认内网冲突
以及阿里云说自己面向daocker了但是image下载不下来。

```
sudo route del -net 172.16.0.0 netmask 255.240.0.0
```

删掉路由表不然服务启动不起来，删除之后不要忘记重启

### 512mb左右内存运行docker的问题

mysql报错 ： `InnoDB: Cannot allocate memory for the buffer pool`

最后innodb启动不起来
google之
```shell
dd if=/dev/zero of=/swapfile bs=1M count=512
mkswap /swapfile
swapon /swapfile
/swapfile swap swap defaults 0 0 to /etc/fstab
```
内存小就添加swap来搞定问题
