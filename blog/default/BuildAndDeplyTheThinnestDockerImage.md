title: 构建最小Docker Image运行网站程序并部署到DaoCloud中
date: 2015-06-06 17:00:22
categories: 技术
tags: [DaoCloud,Docker,Linux,Ubuntu,BusyBox,Golang]
---

![](http://7viiaq.com1.z0.glb.clouddn.com/docker.jpg)

## 意义：
在前人的工作中：[创建尽可能小的Docker容器](http://blog.xebia.com/2014/07/04/create-the-smallest-possible-docker-container/)，[中文翻译版本](http://www.tuicool.com/articles/3yiume)。经过对比我们可以发现 Adriaan de Jonge的的工作方式是通过CGO实现golang的静态编译以达到目的。我们认为这种方式虽然很棒，但是操作起来并不容易，而且在很多项目中编译起来颇为麻烦。

为了精益求精，我们在极致精简与正常使用操作系统（比如说基于ubuntu14.04）做工程角度的平衡，我们还是基于Scratch空Docker image进行构建，使用busybox重新构建了一个大小为5MB左右裁剪过的Linux作为golang运行环境，该系统可以直接使用普通golang编译器输出amd64架构的程序，整个image构建完成之后大小为11.75mb。在DaoCloud中可以达到秒级构建，经过4次测试3次在一分钟以下最短45秒。

Dockerfile以及环境参考：[https://github.com/lijianying10/DaoCloudStaticBlog](https://github.com/lijianying10/DaoCloudStaticBlog)

## golang HTTP 主程序

```golang
package main

import "net/http"
import "fmt"

func main() {
	fmt.Println("Server start") //最简单的日志提示已经开始工作
	http.Handle("/",http.FileServer(http.Dir("/public/"))) // 设置静态文件路由
	http.ListenAndServe(":8080", nil)// 开始http server
}
```

Golang给出了足够多的Web开发工具包，可以编写简短的代码写出高性能web静态服务器。

## Docker image中底层系统的构建
如果方便使用ubuntu 14.04 amd64的情况下：

```shell
apt-get install busybox-static #安装busybox
mkdir /rootfs #新建根
mkdir bin etc dev dev/pts lib proc sys tmp #根据根的格式创建目录
touch etc/resolv.conf # DNS服务器地址
cp /etc/nsswitch.conf etc/nsswitch.conf #名字解析配置
echo root:x:0:0:root:/:/bin/sh > etc/passwd # 新建用户
echo root:x:0: > etc/group #新建用户组
ln -s lib lib64 # 创建软连接
ln -s bin sbin 
cp /bin/busybox bin # 复制bin文件，不然没有命令用
$(which busybox) --install -s bin #同上两步操作，用which是因为需要绝对路径
bash -c "cp /lib/x86_64-linux-gnu/lib{c,m,dl,crypt,util,rt,nsl,nss_*,pthread,resolv}.so.* lib" # 这一步很重要，如果使用了Go调用外部动态库这里需要复制进去，
bash -c "cp /lib/x86_64-linux-gnu/ld* lib" #同上
cp /lib/x86_64-linux-gnu/ld-linux-x86-64.so.2 lib #同上
tar cf /rootfs.tar . #根目录打包
for X in console null ptmx random stdin stdout stderr tty urandom zero ; do tar uf /rootfs.tar -C/ ./dev/$X ; done #复制设备，注意自己的程序调用的设备，比如说random，顺便吐槽一下等待cpu周期太长
```

构建调试过程中的注意事项：
1. 在复制lib文件夹的底层动态库的时候需要不断调试动态库的依赖。
2. 因为容器中要运行不同的项目，有可能对设备与权限上有不同的需求，需要注意定制脚本及时调整。
3. 整个打包过程每次项目几乎只需要一次，调整项目大部分情况都是增量打包因此精简系统文件尽量仔细做。

boot2Docker 跨平台解决方案的Dockerfile(仿造上面来)：
```
FROM ubuntu:14.04
RUN apt-get update -q
RUN apt-get install -qy busybox-static
RUN mkdir /rootfs
WORKDIR /rootfs
RUN mkdir bin etc dev dev/pts lib proc sys tmp
RUN touch etc/resolv.conf
RUN cp /etc/nsswitch.conf etc/nsswitch.conf
RUN echo root:x:0:0:root:/:/bin/sh > etc/passwd
RUN echo root:x:0: > etc/group
RUN ln -s lib lib64
RUN ln -s bin sbin
RUN cp /bin/busybox bin
RUN $(which busybox) --install -s bin
RUN bash -c "cp /lib/x86_64-linux-gnu/lib{c,m,dl,crypt,util,rt,nsl,nss_*,pthread,resolv}.so.* lib"
RUN bash -c "cp /lib/x86_64-linux-gnu/ld* lib"
RUN cp /lib/x86_64-linux-gnu/ld-linux-x86-64.so.2 lib
RUN tar cf /rootfs.tar .
RUN for X in console null ptmx random stdin stdout stderr tty urandom zero ; do tar uf /rootfs.tar -C/ ./dev/$X ; done
```

## 非安全快速解决方案
我们这里的安全性是指，我提供的tar根目录中也许包含了老版本的代码漏洞，或者各种技术债务。如果不是特殊场景下不需要自己手动构建Linux根。

具体解决方案如下：
1. 直接在Dockerhub上使用progrium/busybox
2. 直接使用我们在github上共享的rootfs.tar

方式1不是很推荐，因为构建的适合需要从网上下载依赖，可能会降低构建速度。
非安全情况下推荐方式2构建所需要的docker image

方式1构建的Dockerfile:
```
FROM progrium/busybox
MAINTAINER Jianying Li <lijianying12@gmail.com>

RUN mkdir -p /app
WORKDIR /app
COPY ./static /app/
COPY ./public /public
EXPOSE 80

CMD ["/app/static"]
```

方式2构建使用的Dockerfile：
```
FROM scratch
MAINTAINER Jianying Li <lijianying12@gmail.com>

ADD ./rootfs.tar /
RUN mkdir -p /app
WORKDIR /app
COPY ./static /app/
COPY ./public /public
EXPOSE 80

CMD ["/app/static"]
```

Tip: static为第一节编译好的golang静态程序
public为静态网站根目录

Tip：在构建时请不要使用`sudo docker build -< somefile`
要使用`sudo docker build .`

Tip： 方式2是我在Github中提供的方式。

## 跨平台支持
在本方案中跨平台的支持主要有两个方面：
1. golang源程序的跨平台编译：GOX
2. busybox 系统根(方案二中的Dockerfile中的文件rootfs.tar)的编译可以通过boot2docker来实现。

## DaoCloud部署
DaoCloud给我们提供了一个非常好的测试平台，验证性能工作方式等等，不需要电脑中有Docker环境也可非常方便测试。
在这里我们介绍一下上面提供的[github仓库](https://github.com/lijianying10/DaoCloudStaticBlog)的使用方法。
1. Fork我们提供的代码仓库到您的github账号中。
2. 登陆到[DaoCloud.io](https://dashboard.daocloud.io/)使用代码构建功能
3. 给项目取名字
4. 登陆到github中系统会列出您在Github中的所有项目。并选中DaoCloudStaticBlog项目点击开始创建

在这之后DaoCloud会帮您创建目标项目的DockerImage。
在镜像仓库中找到项目点击部署给容器起名选择环境即可运行。

**注意事项**
1. Golang程序中的端口一定要与Dockerfile中的EXPOSE的号码对应，不然外网无法从DaoCloud访问您的容器。
2. CMD中注意要直接运行网站。
3. 如果条件允许请在本地调试运行之后再部署到DaoCloud中。

## 总结
在本次工作中我们使用了BusyBox来替代原来CGO的实施方案，有效的解决了编译时的麻烦。建立了底层依赖后直接复制应用程序进入Docker image完成在DaoCloud上的部署。在Dockerfile的依赖上我们直接使用了scratch没有任何依赖，加快了构建的速度达到秒级构建。