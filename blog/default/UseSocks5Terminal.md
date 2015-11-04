title: 在终端中使用Socks5代理
date: 2015-05-25 21:53:23
categories: 技术
tags: [shell,linux,mac,delegate] 
---

## 目的
让终端下也能使用代理，加速GIT，Brew 等等更新或者安装

## Socks5
`ssh -D`的代理接口为Socks5，SS提供的代理接口也为Socks5。Terminial貌似只提供了http代理。不能用。
所以这里使用tsocks来使用socks5代理。

## Linux 安装
yum apt自便

## Mac 安装
安装命令
```shell
wget 'http://ftp1.sourceforge.net/tsocks/tsocks-1.8beta5.tar.gz' 
tar -xzf tsocks-1.8beta5.tar.gz && cd tsocks-1.8 
wget 'http://marc-abramowitz.com/download/tsocks-1.8_macosx.patch' 
patch < tsocks-1.8_macosx.patch && autoconf 
./configure && make && sudo make install
```

Tip：注意每一句是什么意思。根据自己情况安装

## 配置例子：

file: `/etc/tsocks.conf`
```
server = 192.168.0.1
server_type = 5
server_port = 1080
```

## 使用方法
有点类似于sudo，在命令前加上tsocks效果拔群。