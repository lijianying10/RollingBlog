title: "时速云上部署Owncloud"
date: 2015-12-16 21:08:33
tags: [mysql,owncloud,docker,tenxcloud]
---

### VER:1

## 介绍
在很多办公场景下，我们需要私密安全的网盘备份空间。在这里推荐Owncloud为作为此场景的解决方案。
我们不推荐使用VPS来部署owncloud一方面是部署比较麻烦，另一方面花费也比较高。我推荐Tenxcloud（时速云）为首选部署平台。
因为在这上面部署省时，省力，省钱。

1. 镜像都是预先准备好的。
2. 不需要使用资源的时候（晚上睡觉，早上还没开始工作）可以释放资源。
3. 时速云价格便宜带宽计算下来比较便宜。

阅读建议：请您都看完有所了解之后再对照动手做。

## 过程

1. 创建云端磁盘。包括数据库与Owncloud存储卷。
2. 创建Owncloud并运行容器。
3. 创建Mysql并运行容器。
4. 配置Owncloud。
5. 开始使用。

## 磁盘准备

![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy1.png)
### 图1

在登录后，进入控制台按照图片的步骤执行。
执行完成第四步之后会弹出一个创建磁盘的对话框。
分别输入磁盘`名称`和`大小`。

建议配置如下图：

![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy2.png)
### 图2

注意： 右侧黄色按钮为格式化磁盘，如果您启动mysql错误请您格式化磁盘之后再试。

## owncloud容器配置与运行

如下图所示:在界面上按照下图顺序创建容器。

![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy3.png)
### 图3

按照下图操作顺序`选择镜像来源`。
![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy4.png)
### 图4

按照下图操作顺序`进行容器配置`。
![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy6.png)
### 图5

最后点击创建完成owncloud部署。（owncloud需要2分钟左右的启动时间）

在容器管理界面如图所示。
![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy7.png)
### 图6
您会看到owncloud已经在运行状态。注意owncloud的服务地址为您未来使用的服务器地址。

## MySQL容器配置与运行

入下图所示进入创建容器界面。
![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy3.png)
### 图3

在容器管理界面`选择镜像来源`按照图中操作。
![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy9.png)
### 图7

在`容器配置`界面按照如下操作顺序操作。
![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy10.png)
### 图8

在`高级配置`中按照如下操作顺序操作。
![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy11.png)
### 图9

选择启动即可启动mysql容器。

## 配置Owncloud

先进行配置信息的收集：

`控制台->服务->服务->你的MySQL容器(教程中的ownsql)点开左边小箭头->查看内网服务地址`
如下图所示。
![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeployQQ20151216-0@2x.png)
### 图10

记下地址别名，不需要记录端口号。

打开`owncloud服务地址`

看到如下图所示的配置界面：
![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy8.png)
### 图11

1. 用户名密码为您的owncloud管理员账号密码请您牢记
2. 按照步骤1 打开`存储与数据库`。
3. 选择`MySQL/MariaDB`
4. 数据库用户名：root
5. 密码为：图9中您输入的密码（第二个框里面的）
6. 数据库名为：mysql
7. localhost 这里（数据库服务器地址）替换成您刚才记下的mysql内网别名。

点击完成安装即可！

![](http://7viiaq.com1.z0.glb.clouddn.com/ownclouddeploy12.png)
### 图12

## 总结

总体步骤上来讲不是很难。安装完成之后客户端可以在各大应用市场中找到。或者登陆owncloud官网来下载最新版本客户端。




