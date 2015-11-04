title: RHEL添加本地源
date: 0000-00-00 00:00:00
categories: 技术
tags: [Linux,RHEL] 
---

# 之前经常用CentOS
最近看到swip的计算服务器里面有不少系统bug想着写脚本帮忙解决了。
现搞个HREL6.3再说。
结果遇到了没有C6Media的配置
经常用CentOS的都知道实用的C6 热泪盈眶。


RHEL添加本地源（光盘）
```bash
cat /etc/yum.repos.d/rhel-local.repo

[rhel6.3-local]
name=RHEL 6.3 local repository
baseurl=file:///opt/yum/rhel6.3/
gpgcheck=1
gpgkey=file:///etc/pki/rpm-gpg/RPM-GPG-KEY-redhat-release
enabled=1
```
