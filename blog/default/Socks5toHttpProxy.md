title: "Socks5代理转换成HTTP代理"
date: 2015-05-28 16:22:52
categories: 技术
tags: [Linux,ShadowSocks,Cow,golang] 
---

## 意义
之前有朋友跟我说很多SDK的更新经常失败，可能是国际入口带宽减少的原因。因此我们需要使用http搭理寻找更好的线路帮忙更新开发SDK。

另外SHELL只要设定http_proxy变量就可以使用http代理非常方便。

## 解决方案：Cow

安装：`curl -L git.io/cow | bash`

配置：
```
#开头的行是注释，会被忽略
# 本地 HTTP 代理地址
# 配置 HTTP 和 HTTPS 代理时请填入该地址
# 或者在自动代理配置中填入 http://127.0.0.1:7777/pac
listen = http://127.0.0.1:7777

# SS 二级代理
proxy = ss://aes-128-cfb:[Password]@[server]:8388
```

好处： 能帮你直接export出HTTP代理。方便shell使用。
`更多说明详见：[https://github.com/cyfdecyf/cow](https://github.com/cyfdecyf/cow)`

## 测试


撞墙
```
$wget http://pbs.twimg.com/profile_images/603341077702049792/hzGDhXNe_bigger.png
--2015-05-28 15:11:15--  http://pbs.twimg.com/profile_images/603341077702049792/hzGDhXNe_bigger.png
Resolving pbs.twimg.com... 199.96.57.7
Connecting to pbs.twimg.com|199.96.57.7|:80... ^C
```

翻墙
```
http_proxy="127.0.0.1:7777" wget http://pbs.twimg.com/profile_images/603341077702049792/hzGDhXNe_bigger.png
--2015-05-28 15:11:44--  http://pbs.twimg.com/profile_images/603341077702049792/hzGDhXNe_bigger.png
Connecting to 127.0.0.1:7777... connected.
Proxy request sent, awaiting response... 200 OK
Length: 3151 (3.1K) [image/png]
Saving to: 'hzGDhXNe_bigger.png'

hzGDhXNe_bigger.png                         100%[==========================================================================================>]   3.08K  --.-KB/s   in 0s

2015-05-28 15:11:45 (429 MB/s) - 'hzGDhXNe_bigger.png' saved [3151/3151]
```
