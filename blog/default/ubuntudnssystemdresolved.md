title: 网关服务器部署DNS的冲突问题解决
date: 2021-01-14 10:40:11
categories: 技术
tags: [ubuntu,dns,systemd-resolved]
---

因为5.4内核对网络的管理修改很多，在自己部署DNS服务器时，发现systemd-resolved.service占用端口53

所以这里Ubuntu 20.04 的运维SOP如下:

1. netplan 设置nameserver 127.0.0.1 apply
2. systemd-resolved.service stop && disable
3. 部署你自己的dns server
4. delete softlink /etc/resolv.conf 写入 `nameserver 127.0.0.1`

注意，Docker daemon会默认指定 8.8.8.8 8.8.4.4 作为dns 不会再取系统配置 (Docker 20.10.1 行为)
所以如果做了特别的DNS配置需要对Daemon配置方法参考这里 https://forums.docker.com/t/local-dns-or-public-dns-why-not-both-etc-docker-daemon-json/54544

