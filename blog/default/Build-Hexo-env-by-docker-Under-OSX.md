title: "在OSX下使用docker构建hexo环境"
date: 2015-05-28 16:31:45
categories: 技术
tags: [Linux,osx,mac,docker,virtualbox,ubuntu,hexo] 
---

### Version 1.1

![](http://cdn.liginc.co.jp/wp-content/uploads/2014/08/233.png)

## 意义
1. 截止目前为止，在docker hub上还看不到hexo 3+版本的镜像构建。
2. 由于Node环境的安装时间比较长，对于电脑比较多的人相对来说还是比较麻烦的。
3. 基于DockerFile构建自己的HEXO环境非常灵活，可以根据自己的情况进行定制。
4. 可以导入导出复制环境部署到其他电脑中。

# 实践开始

## docker安装

大部分Linux，比如说Ubuntu，Debian都可以通过如下命令安装：
```
curl -sSL https://get.daocloud.io/docker | sh
```

TIP：别忘了配置好[DaoCloud](https://www.daocloud.io/)加速，不然构建速度很慢。

## 构建 
`docker build -t hexo3 - < hexo3.dockerfile`

hexo3.dockerfile
```
FROM node:slim

MAINTAINER Jianying Li <lijianying12@gmail.com>

# instal basic tool 
RUN apt-get update && apt-get install -y git ssh-client ca-certificates --no-install-recommends && rm -r /var/lib/apt/lists/*
# set time zone
RUN echo "Asia/Shanghai" > /etc/timezone && dpkg-reconfigure -f noninteractive tzdata
# install hexo
RUN npm install hexo@3.0.0 -g
# set base dir
RUN mkdir /hexo
# set home dir
WORKDIR /hexo
 
EXPOSE 4000

CMD ["/bin/bash"]
```

定制自己的image请注意，现在最简洁的三个包内容为： 
	1. git，部署的时候用（如果不用git部署请去掉）。
	2. ssh-client（ssh方式的git部署依赖）。
	3. ca-certificates（https方式的git部署依赖）。

TIP: 在shell中或者lib中调用https方式通讯的时候如果报错`Problem with the SSL CA cert (path? access rights?)`可以通过安装包：ca-certificates来解决问题，yum apt中都是如此。

构建时间大概十几分钟完成。

## 准备把实体机(host)上的文件挂载到docker中

1. 安装Guest Additions, 因为要使用Shared Floader。
2. 使用命令 `sudo mount -t vboxsf [sharename] [dist]`来挂载共享目录。

## 运行
`docker run -it -d -p 4000:4000 -v /root/blog:/hexo/ --name hexo hexo3 `

注意路径 `/root/blog/` 是我VirtualBox 虚拟机中blog存储的位置。

注意参数`/root/blog/`需要使用绝对路径

其他的参数可以很容易的在[manual](https://docs.docker.com/userguide/)中找到意义。

## 备份与还原

```
	#docker save hexo3 > /root/hexo3.tar

	#docker load < /root/hexo3.tar
```

`注意这里使用save而不是export 因为需要保存历史层`

参考导出大小：
```
du -h /root/hexo3.tar
261M	/root/hexo3.tar
```
从以上所有的工作中，对比虚拟机进行环境的构建打包，docker具有构建环境时间更短，打包文件更小的特点。

## 使用容器操作blog
`docker exec -it hexo /bin/bash`
Tip: 虽然做到了用docker构建一个非常方便移植的hexo环境，但是运行命令hexo的时候有点慢，但不是那种忍受不了的慢。

## 小技巧
在调试的时候可以使用 docker rm $(docker ps -q -a) 一次性删除所有的容器，docker rmi $(docker images -q) 一次性删除所有的镜像。
