title: 修改mysql默认密码
date: 0000-00-00 00:00:00
categories: 技术
tags: [mysql]
---

#修改mysql 密码

## 方法1 （我常用的方法）
比较简单
```
set password=password('1');
flush privileges;
```
登陆的时候别忘了。 mysql -u root
带上用户名
此方法常用于新装不带密码的数据库。
