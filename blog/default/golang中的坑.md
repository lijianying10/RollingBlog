title: golang 类型转换
date: 0000-00-00 00:00:00
categories: 技术
tags: golang
---

#golang中的坑

## 类型转换。
1. golang中的类型转换全都是通过类似C语言这种的atoi itoa这种的
1. session中的坑的解决办法
```go
strconv.Itoa(int(userSession.Get("loginTime").(int)))+"<br/>"
```
(import strconv)
