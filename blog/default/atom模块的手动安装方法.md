title: atom模块的手动安装方法
date: 2015-02-12 16:47:59
categories: 技术
tags: atom
---

最近发现部分地区apm已经不能用了因此在这给大家介绍一下手动安装的过程。
主要分为如下三个步骤
1. 找到repo并且克隆到自己的电脑上。
2. 处理node依赖
3. 连接到atom中。

## 找到代码
1. 首先到atom.io中找到你需要的package
2. 在package首页上又一个repo的连接，非常醒目。
3. 跳转到github的页面之后，就可以克隆到自己的电脑中了

## 处理依赖
1. 使用npm处理依赖就可以了。
2. npm在国内是由mirror的，比如说淘宝。下面的例子就使用套高的mirror来处理依赖
3. 使用命令npm install

## 连接到atom中
3. 在插件代码的根目录下使用`apm link`完成连接。

## 总结
```shell
mkdir -p ~/.atom/git
cd ~/.atom/git/
git clone [REPO] ## 克隆代码回来
cd [REPO NAME]
npm install  --registry=https://registry.npm.taobao.org ## 使用淘宝的mirror
apm link
#以上命令请酌情考虑使用。
```
