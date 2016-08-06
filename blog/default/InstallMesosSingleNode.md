title: 在单个节点上最快速度安装Mesos
date: 2016-07-28 23:16:25
categories: 技术
tags: [mesos,centos,marathon,etcd,flannel]
---

## 背景

在写Mesos应用的时候会用到这个Case

## [网络安装脚本](https://github.com/lijianying10/FixLinux/blob/master/mesos/mesos-systemd/singlenode/netinstall.sh)

## [本地安装脚本](https://github.com/lijianying10/FixLinux/blob/master/mesos/mesos-systemd/singlenode/localinstall.sh)

注意安装文件本地化的位置。在脚本里面了请提前下载


## 配置

```
HostIP='192.168.56.112'
IFACE='enp0s3'
```

配置本机IP地址以及对外访问的Interface名字，可能的都有`ethX`,`emX`,`enpXsX`等。


## 安装组件：

```
1. Mesos
2. Marathon
3. Flannel
4. ETCD
```

注意Docker需要自行[安装](http://get.daocloud.io/#install-docker)

## 总结

速度很快，成功率很高。适合自己在单点的时候开发应用。

