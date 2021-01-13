title: ubuntu下使用netplan修改网络信息
date: 2021-01-13 12:23:15
categories: 技术
tags: [ubuntu,netplan,network]
---

从Ubuntu 17.10之后默认使用netplan来修改网络信息

一个例子：

```
l@l:~$ cat /etc/netplan/00-installer-config.yaml
# This is the network config written by 'subiquity'
network:
  ethernets:
    ens33:
      dhcp4: no
      addresses:
        - 192.168.123.44/24
      gateway4: "192.168.123.1"
      nameservers:
        addresses:
          - "192.168.123.102"
  version: 2
```

修改之后使用命令 enable config：

```
sudo netplan apply
```
