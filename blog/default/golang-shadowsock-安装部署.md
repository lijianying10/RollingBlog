title: golang shadowsock 安装部署
date: 2015-02-05 09:25:02
categories: 技术
tags: [golang,ss,Linux]
---

最近一直在用ss但是python的并发并不是很好所以我改换成golang 的ss了。
[代码下载](https://github.com/lijianying10/shadowsocks-go)

[编译好的直接下载](http://dl.chenyufei.info/shadowsocks/)

server配置（跟官方的不一样）：
1. 新建文件夹
2. 进入文件夹后新建文件：config.json  pid.cfg  start.sh  test.log
  1. config 是server运行的配置文件
  2. pid.cfg 是server运行的进程号记录的位置
  3. start.sh 是服务启动时候运行的脚本
  4. test.log 这个是排查故障的时候用的错误记录文件


### config.json
``` json
{
    "server":"XXX.com", //服务器ip地址或者绑定的域名
    "server_port":8088, //  运行的端口
    "local_port":1080, //本地运行端口
    "password":"XXXXXXX", // 端口的密码
    "timeout":600, // 不打算解释了
    "method":"aes-128-cfb" // 推荐的加密算法，128 强度其实足够了。
}
```

### pid.cfg 留空

### start.sh
```shell
ssserver-go -d start --log-file test.log --pid-file pid.cfg
```
ssseerver-go 是我从网站上下载的服务器端，然后改了个名字放到path里面了
-d是托管服务的意思（离开控制台了之后还能运行）
--log-file 不解释
--pid-file 不解释

### test.log 留空

## ss多用户配置

因为我们用的是文件夹管理的不同进程不同配置。
我们只需要复制文件夹然后修改里面的配置文件就可以了。
如果我们的系统扩容的时候需要多个用户独立管理。
### 复制文件夹之后
1. 需求修改端口跟密码。
2. 如果不喜欢后台运行可以使用screen命令随时切换回进程
3. 定时停止服务： at 21:00 tomorrow # 类似这种命令什么的网上有的是 at命令怎么定时运行。
4. 进入at之后，可以写入当时需要运行的shell脚本， 比如说kill 某个pid就ok了。
5. 如果任务比较多比较复杂那我比较推荐crontab 这里就不在多说了。

## 客户端的问题

1. 支持原来的python客户端。
2. golang命令登陆：
  1. 命令名自己搞定吧
  2. 命令的参数如下：-s [ServerIP] -p [ServerPort] -k [ServerPassword] -m aes-128-cfb -b 127.0.0.1 -l 1080
