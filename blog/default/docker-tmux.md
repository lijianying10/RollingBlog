title: docker+tmux 加强容器调度
date: 2015-12-13 13:07:32
categories: 技术
tags: [tmux,docker]
---

## VERSION:1

![](http://7viiaq.com1.z0.glb.clouddn.com/u=1306363718,2217319814&fm=15&gp=0.jpg)

## 摘要

为了让自己做事更加自动化，把重复的工作尽可能降到最低，平时不但需要写很多固定操作的脚本来加快工作效率。
搞搞调度环境也是需要的。

本篇通过Docker+Tmux在RancherOS上做开发平台来实现最快速的Docker调度方便自己开发。

1. 可以最快速度进入到调度容器中。
2. 该容器有docker deamon 的所有控制权限。
3. 可以在容器内的Tmux中跳转到其他容器中。方便调度开发。

经过`2`个版本的迭代终于搞定。到达1.0版本。

## Docker Registry

`docker pull index.tenxcloud.com/philo/dmonit:1.0`

## Dockerfile

[https://github.com/lijianying10/FixLinux/blob/master/dockerfiles/dmonit/Dockerfile](https://github.com/lijianying10/FixLinux/blob/master/dockerfiles/dmonit/Dockerfile)

## 主要功能

### 启动方法

```
docker run -it --name kkk -d -p 445:22 -v /usr/local/bin/docker:/usr/local/bin/docker -v /var/run/docker.sock:/var/run/docker.sock -e 'PUBKEY=ssh-rsa XXXX' index.tenxcloud.com/philo/dmonit:1.0
```

参数解释：

1. 映射22端口到其他位置，防止冲突。
2. 挂载docker命令到容器中。
3. 挂载Docker API的Named PIPE控制docker。
4. 环境变量：PUBKEY 写入控制机的ssh 的 publickey。

### 进入控制方法

方便登陆Docker容器的配置文件。

```
# cat ~/.ssh/config
Host dmmm
hostname 192.168.99.100
user root
port 445
```

输入命令：`ssh dmmm` 可进入调度容器。

### 解释为啥使用ssh

主要是看了这个[Docker ISSUE](https://github.com/docker/docker/issues/8755)
然而他们并没有解决`docker exec -it` 和`docker run -it`不能使用`tmux`的问题。

为了能获得一个好用的tty所以，也为了节省时间所以就用了OpenSSH。

### xdev

此命令用来开一个开发tmux还可以进入之前开过的tmux window。

![](http://7viiaq.com1.z0.glb.clouddn.com/QQ20151213-5@2x.png)

上面会标记项目名，预设：编辑器，运行窗口，测试窗口，日志窗口，数据库查看窗口。

后面有当前内存使用，当前时间，当前Unix时间戳。

1. xdev 有只有一个参数是给session命名的。
2. 在不同的终端输入一样的xdev命令会进入到同一个session中。
3. 非常方便的窗口恢复切换。

### e 

如果你跟我一样无法忍受`docker exec -it [container] /bin/bash`。
打太多次打到烦。
所以这个脚本是这样的：

```
[#2#root@75477389dbdf ~]$cat $(which e)
docker exec -it $1 /bin/bash
```

因为挂载了docker程序以及named pipe 所以在这里面是可以管理docker的。

### tmux

切换开发Tab：

快捷键： `M-h`切换到上一个Tab。

快捷键： `M-l`切换到下一个Tab。

可以和vim很好的结合。包括其他容器内的vim都可以。

![](http://7viiaq.com1.z0.glb.clouddn.com/QQ20151213-6@2x.png)

如图所示： 上面为vim的tab，下面为Tmux的tab。



## 总结

有了这个容器之后，可以非常方便的调度其他容器。可以提升开发效率。减少操作次数。频率。如果有好的意见一定要提醒我哦。先谢过。

Tmux 的配置在这里：[https://github.com/lijianying10/FixLinux/blob/master/dotfile/.tmux.conf](https://github.com/lijianying10/FixLinux/blob/master/dotfile/.tmux.conf)如果需要定制请FORK我的REPO。


