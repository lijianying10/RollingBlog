title: "RancherOS 初期使用的感受"
date: 2015-11-29 12:14:16
tags: [rancher,docker,linux]
---

### VER: 0-prerelease

## 意义
随着我们团队线下业务的压力越来越高。并且我们团队是一个并没有运维的团队想部署与管理集群光安装都是比较麻烦的事情。
在实际情况中我们团队一共有8台服务器12颗CPU 约300G内存的样子。
之后我们是打算部署K8S 或者Mesos。Hadoop之类的做一些线下的数据处理工作。
因为我们团队依赖Docker非常多希望能够快速部署Docker然后直接就能生产了。

之前一直打算使用CoreOS，但是就国内的网络情况来说，实现他们的基本技术特性还是挺难的。
但是后来接触到了RancherOS之后情况就有所改变了。
因为安装ISO只有`20mb`怎么都下载回来了，同时PID1就是Docker，安装的时候可以通过国内的Image来加速安装这样的话，在公司内快速部署安装还是挺有戏的。
最后在3台服务器上部署RancherOS 包括开机时间，下载操作系统的时间每台机器跑秒安装使用了`2分20秒`的时间完成。极大的节省了人力。

## 安装过程

### ROS
学会这个命令是RancherOS开始最重要第一环。学会了它可以让你事半功倍。因为RancherOS的配置都是使用这个工具来完成的。
的确非常强大方便。

参考文档在这里[http://docs.rancher.com/rancher/](http://docs.rancher.com/rancher/)

### 安装RancherOS 到硬盘

注意启动之后默认的账号密码为rancher:rancher

[http://docs.rancher.com/os/running-rancheros/server/install-to-disk/](http://docs.rancher.com/os/running-rancheros/server/install-to-disk/)

参考上面连接的方法。但是对于国内网络来说我们嗨需要另外一个参数 `-i`

`sudo ros install -c cloud_config.yml -i index.tenxcloud.com/philo/rancheros:v0.4.1 -d /dev/sda`

我这里准别好了0.4.1版本的放在时速云上了`index.tenxcloud.com/philo/rancheros:v0.4.1`

`注意！`

一定要配置yml文件，不然安装完之后自己就登陆不上去了。
一定要注意安装硬盘的位置，别装错地方了。注意RancherOS版本，我用的是0.4.1

注意RancherOS默认NS服务器是google的，需要自己做调整，修改配置文件/etc/resolv.conf可以解决这个问题。

#### 国内安装加速点：

```
index.tenxcloud.com/philo/rancheros:v0.4.1
index.tenxcloud.com/philo/rancheros:v0.4.3


```

## RancherOS 的结构

系统启动非常快，里面只有两个关键部分，一个是System-docker另外一个是docker

系统的docker跑了所有系统中需要的进程

