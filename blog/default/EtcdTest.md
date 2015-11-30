title: "Etcd 线上服务器试用"
date: 2015-11-29 11:12:01
tags: [golang,etcd]
---

### VER:0-preRelease

## 介绍
Etcd 作为golnag的杀手级应用，还是非常值得研究的。下面几个链接是从我的认知角度上认为新手应该看的内容。

[http://www.oschina.net/p/etcd](http://www.oschina.net/p/etcd)

[https://coreos.com/etcd/](https://coreos.com/etcd/)

最后在一个不起眼的地方，etcd release 的目录里面有一个Documentation初学者的无尽宝藏
当然也可以在这里看[https://github.com/coreos/etcd/tree/master/Documentation](https://github.com/coreos/etcd/tree/master/Documentation)

## 对于我的意义

1. 它为我提供了更加可靠的线上服务的配置存储，这是最基本的。
2. 能够结合监控系统。针对服务进行监控。

## 驱动上

新手不推荐直接使用Etcd驱动还要读更多的文档，都不如直接使用RESTful了。
所以我在线上服务器也是这么搞的。

## 学习过程

### 基本信息获取

1. 下载[https://github.com/coreos/etcd/releases](https://github.com/coreos/etcd/releases)
2. 文档[https://coreos.com/etcd/docs/latest/](https://coreos.com/etcd/docs/latest/)

查看Release页面的时候一定要注意，查看更新内容，有没有什么新的操作方式。

### 基本操作

#### 单节点服务启动参考：

```
./etcd -name infra0 -initial-advertise-peer-urls http://0.0.0.0:2380 \
  -listen-peer-urls http://0.0.0.0:2380 \
  -listen-client-urls http://0.0.0.0:2379 \
  -advertise-client-urls http://0.0.0.0:2379 \
  -initial-cluster-token etcd-cluster-1 \
  -initial-cluster infra0=http://0.0.0.0:2380 \
  -initial-cluster-state new
```

这也是我日常开发调试的时候用的最常用的命令了。

#### 针对数据进行基本操作的参考文档

[https://github.com/coreos/etcd/tree/master/etcdctl](https://github.com/coreos/etcd/tree/master/etcdctl)


#### 基本针对数据操作学会了WebAPI参考文档自己写驱动用

[https://coreos.com/etcd/docs/latest/api.html](https://coreos.com/etcd/docs/latest/api.html)

根据CURL结合自己使用的语言来进行操作哦。相信自己根据需求来写是最简洁最不容易出错的。

#### 服务写好之后开始做集群

集群的服务启动参考文档[https://coreos.com/etcd/docs/latest/clustering.html](https://coreos.com/etcd/docs/latest/clustering.html)

集群服务器动态添加节点(etcd里面管叫Member)的方法[https://coreos.com/etcd/docs/latest/admin_guide.html](https://coreos.com/etcd/docs/latest/admin_guide.html)

### 做完这些之后要研究的东西方向就不会很一样了。

1. hack代码
2. 自己写更加全面的驱动或者取看别人的驱动是怎么写的，然后再用
3. 认证，安全问题
4. 结合业务看性能问题，看看能不能让更多的业务依赖ETCD

### 总结

我是走过了好多坑，看了好多冤枉的文档才总结出来这么一篇比较小的文档。
如果能按照我的顺序看完之后相信您可以很容易的在ETCD的道路上继续往下走。
当然我也不是一帆风顺的，目前跨集群部署的时候遇到了 floating ip 的问题。
这些都是没有文档的问题，只能自己慢慢想办法。
