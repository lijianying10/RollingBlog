title: docker下进行Android编译
date: 2016-01-03 23:04:14
categories: 技术
tags: [docker,android,ci]
---

## VERSION:1

![](http://7viiaq.com1.z0.glb.clouddn.com/dockerandroid.jpeg)

## 警告：此文可能过于严肃

## 意义

1. 极大的缩短安卓开发到`测试`到`产品`到`渠道`的距离。
2. 给安卓程序员减轻负担。
3. Google做的环境已经特别到位了，放到docker里面明显不会有多大的坑（误，逃）。

## Dockerfile

[https://github.com/lijianying10/FixLinux/blob/master/dockerfiles/androidautobuild/Dockerfile](https://github.com/lijianying10/FixLinux/blob/master/dockerfiles/androidautobuild/Dockerfile)

## 团队协作的故事

在敏捷开发的站立会议上，我作为一个后端程序员发现安卓程序员发布的时间大概是半个小时左右的时间我觉得时间太长了应该缩短一些，成为了研发这个东西的目的。
但是研究了一段时间之后发现了很多需要解决的问题：

1. 很多东西是被墙的。
2. 很多依赖不能够复用（各种pom包）。
3. gradlew怎么能快速安装不需要从网上下载。
4. 自动对齐。
5. 自动签名。
6. 自动混淆。

## 考虑范围

1. 系统底层依赖
2. JDK
3. Andorid-SDK
4. Gradlew
5. 项目依赖

## 构建解释

`建议：使用国外vps构建，不然要等很长时间`

### 构建变量

```
ENV JAVA_HOME /jdk1.8.0_65
ENV ANDROID_HOME /opt/android-sdk-linux/
ENV ANDROID_SDK_FILENAME android-sdk_r24.4.1-linux.tgz
ENV ANDROID_SDK_URL http://dl.google.com/android/${ANDROID_SDK_FILENAME}
ENV PATH ${PATH}:${ANDROID_HOME}/tools:${ANDROID_HOME}/platform-tools:${JAVA_HOME}/bin/
```


### 底层依赖

gcc一类的，注意我们需要安装32位编译环境，以及git wget。

```
RUN sudo apt-get update && sudo apt-get install -y gcc-multilib lib32z1 lib32stdc++6 git wget && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
```

### 安装JDK

```
RUN cd / && wget --no-check-certificate --no-cookies --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/8u65-b17/jdk-8u65-linux-x64.tar.gz &&\
    tar xf jdk-8u65-linux-x64.tar.gz \
    && rm -rf $JAVA_HOME/src.zip $JAVA_HOME/javafx-src.zip $JAVA_HOME/man /jdk-8u65-linux-x64.tar.gz
```

这里使用Cookie来确定同意协议。

### 安装SDK

```
RUN cd /opt && \
    wget -q ${ANDROID_SDK_URL} && \
    tar -xzf ${ANDROID_SDK_FILENAME} && \
    rm ${ANDROID_SDK_FILENAME} &&\
    echo y | android update sdk --no-ui --all --filter tools,platform-tools,extra-android-m2repository,android-21
RUN echo y | android update sdk --no-ui --all --filter android-22,build-tools-21.1.2,build-tools-22.0.1
```

1. 因为最好每一个layer控制在1G以内所以这里切割用了两个run。
2. 注意`SDK用您项目中需要的最高的版本安装到image里面然后向下安装，不然会出现tool这个文件夹无法运行工具的情况。`
3. 注意`上面的SDK plateform等都是根据我们的项目来的，详细的摸索一下项目代码就知道依赖什么了。`

## 准备项目

1. 使用git clone 同步项目目录。
2. 进行第一次手动构建编译。 命令为：`gradlew assembleDebug`

第二点中目的有三个：

1. 查看自己的依赖是否正确(android update sdk)这里，如果多了精简掉，如果少了加上。
2. 自动下载项目中所有的依赖。
3. 安装gradlew。

需要备份的点有两个

1. `/root/.gradle` 备份这个目录可以在以后自动化构建的时候不需要重复安装gradlew。
2. `$PROJDIR/.gradle`项目依赖的备份，备份了。($PROJDIR 为您的项目根目录位置)

都备份之后下次编译就不需要网络了(容器就不需要梯子了,这点对提升速度很重要）。

## 根据项目构建Image

因为每个项目的依赖不尽相同所以需要针对项目定制化。大概运行目标如下：

1. 创建容器。
1. 找个方法同步代码git，FTP，NFS等等方法。
2. 把上面两个备份点放到指定位置等待使用。
3. 执行构建输出。
4. 销毁容器。

`其实只要能做到上面这一点，加一个git hook 加上简单的发布就是一个简单的CI了。`

## 对齐，签名，混淆

根据下面参考文档可以对项目的build.gradlew进行调整

签名是在Android节点下面加入如下代码：

``` json
signingConfigs {
release{

            storeFile file("../xxxxxxx.keystore")
            storePassword "xxxxxx"
            keyAlias "xxxxx"
            keyPassword "xxxxx"
}}
```

在buildTypes 下面的release下面加入如下选项：

``` json
signingConfig signingConfigs.release
```

对齐方面根据安卓官方文档说明按照上面两步代码修改之后已经对齐。可以准备安装了。

混淆(proguard)

在buildTypes 下面的release 下面加入如下选项：

``` josn
proguardFiles getDefaultProguardFile('proguard-android.txt'), 'proguard-rules.pro'
```

## 总结

经过研究以及实战如果是用普通的笔记本电脑我们的应用30多个渠道大概需要使用31分钟的时间来进行构建。
如果使用RancherOS服务器Xeon X5675 两颗CPU 48G内存的刀片服务器构建的时间是`1分6秒`。

因为写文章的时间仓促，很多地方写的不明白希望大家能够指出来，方便我改进，另外本人安卓水平非常一般请大神们批评指正。十分感谢。

注:`下面的文献非常具有参考价值。`

## 主要参考文献：
[1] Building and Running from the Command Line [http://developer.android.com/intl/pt-br/tools/building/building-cmdline.html](http://developer.android.com/intl/pt-br/tools/building/building-cmdline.html)

[2] Configuring ProGuard [http://developer.android.com/intl/pt-br/tools/help/proguard.html](http://developer.android.com/intl/pt-br/tools/help/proguard.html)

[3] Signing Your Applications [http://developer.android.com/intl/pt-br/tools/publishing/app-signing.html](http://developer.android.com/intl/pt-br/tools/publishing/app-signing.html)
