title: 一个非常漂亮安全的Linux web shell
date: 0000-00-00 00:00:00
categories: 技术
tags: [shell,Linux,butterfly] 
---

http://paradoxxxzero.github.io/2014/02/28/butterfly.html
http://paradoxxxzero.github.io/2014/03/21/butterfly-with-ssl-auth.html

安装的时候注意ssl的devel libffi-dev就行了。
瞅瞅错误基本上就能编译过去了。

```sh
$ sudo butterfly.server.py --generate-certs --host="192.168.0.1" # Generate the root certificate for running on local network
$ sudo butterfly.server.py --generate-user-pkcs=foo              # Generate PKCS#12 auth file for user foo
```
cd /etc/butterfly/ssl/
butterfly_ca.crt这个是要 import到浏览器里面的
foo.p12这个哈。也是要import到浏览器里面的
之后访问host的那个ip地址就能用https 安全使用控制台了。

```sh
$ sudo butterfly.server.py --host="192.168.0.1" # Run the server
```
这是运行啦~~~


## 总结
1. 虽然blog看起来不起眼，但是这个配合zsh真心是非常舒服的
2. install的时候还是比较麻烦的，有的时候编译东西不是很舒服，经常出错。
3. 如果是远程服务器的话一定要配合ssl登陆，会安全一点。
4. 吐槽一下，，ssl现在也不怎么安全了。经常爆漏洞哈。
