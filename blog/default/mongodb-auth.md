title: mongodb 身份验证问题解决
date: 2015-03-05 13:22:07
tags: [mongodb,Linux]
categories: 技术
---

MongoDB默认正在就是可信的内网情况使用的因此刚开始的时候是没有身份验证的。
但是我们在做项目的时候是需要使用身份验证的。
因此本篇解决如下问题
```
1. 添加账号密码
2. 在shell状态下验证账号密码是否设置成功
```

## 添加账号密码

```shell
$ mongo
> use SthDB # 切换数据库
> db.addUser('sa','sa') //添加账号密码
> db.system.users.find()   //查看所有账号密码
```

## 在shell下验证账号密码

```shell
mongo localhost:27017/SthDB -u sa -p // [ip:port/DB] -u [user] -p
```

错误情况：
> 2015-03-05T13:20:55.344+0800 Error: 18 { ok: 0.0, errmsg: "auth failed", code: 18 } at src/mongo/shell/db.js:1210
exception: login failed

正确情况：
> connecting to: localhost/SthDB

## 通过IP地址访问局域网内的mongodb的坑。
1. 配置里面要重新绑定你的ip地址。（/etc/mongodb.conf 中 bind_ip =）
2. 在本地登陆的时候 也要用你bind的ip才能登陆不然就报错。
