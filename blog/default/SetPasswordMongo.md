title: 设定mongodb 密码
date: 0000-00-00 00:00:00
categories: 技术
tags: mongodb
---

#mongodb 密码设定
```bash
use admin
db.addUser('root','123456')
db.system.users.find()
#如果出现了你输入的账号就证明成功了
```
