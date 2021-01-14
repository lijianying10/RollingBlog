title: iptables 使用 Redsocks 时候 Docker 没有网络问题的分析
date: 2021-01-13 20:33:10
categories: 技术
tags: [iptables,docker,redsocks,network]
---

今天我自己遇到了一个故障，Docker里面没有网络。

最开始以为是我iptables设置有错误。通过仔细阅读文档发现估计并不是iptables的配置有问题。

后来发现一个特征，凡是走Redsocks RETURN target 全都可以联通，发现这个主要是我无意间发现内网是通的。

所以定位核心故障是Redsocks不通的，这时候发现Redsocks的Listen address 是有问题的.
因为Docker是通过网桥来到local的所以IP地址不是本地所以redsocks需要监听0.0.0.0


## 所以正确的配置如下

### chnroute 生成中国的ip地址

```
curl 'http://ftp.apnic.net/apnic/stats/apnic/delegated-apnic-latest' | grep ipv4 | grep CN | awk -F\| '{ printf("%s/%d\n", $4, 32-log($5)/log(2)) }    ' > chnroute.txt
```

我们默认绕过中国的ip地址

### redsocks 配置文件

```
base {
        // debug: connection progress & client list on SIGUSR1
        log_debug = off;

        // info: start and end of client session
        log_info = on;

        /* possible `log' values are:
         *   stderr
         *   "file:/path/to/file"
         *   syslog:FACILITY  facility is any of "daemon", "local0"..."local7"
         */
        log = "syslog:daemon";

        // detach from console
        daemon = on;

        /* Change uid, gid and root directory, these options require root
         * privilegies on startup.
         * Note, your chroot may requre /etc/localtime if you write log to syslog.
         * Log is opened before chroot & uid changing.
         */
        user = redsocks;
        group = redsocks;
        // chroot = "/var/chroot";

        /* possible `redirector' values are:
         *   iptables   - for Linux
         *   ipf        - for FreeBSD
         *   pf         - for OpenBSD
         *   generic    - some generic redirector that MAY work
         */
        redirector = iptables;
}

redsocks {
        /* `local_ip' defaults to 127.0.0.1 for security reasons,
         * use 0.0.0.0 if you want to listen on every interface.
         * `local_*' are used as port to redirect to.
         */
        local_ip = 0.0.0.0;
        local_port = 12345;

        // `ip' and `port' are IP and tcp-port of proxy-server
        // You can also use hostname instead of IP, only one (random)
        // address of multihomed host will be used.
        ip = 1.2.3.4;
        port = 8089;


        // known types: socks4, socks5, http-connect, http-relay
        type = socks5;

        // login = "foobar";
        // password = "baz";
}
```

### iptables 配置 


```
set -x
set -e

# 下面两句是添加中国的ip表
sudo ipset create chnroute hash:net
cat chnroute.txt | sudo xargs -I ip ipset add chnroute ip

# 在 nat 表中创建新链
iptables -t nat -N REDSOCKS
## 首先这个代理要屏蔽网络代理出口服务（直接return嘛）

# 下面是添加例外端口的例子
#iptables -t nat -A REDSOCKS -p tcp --dport 59237 -j RETURN
# 下面是添加例外ip地址的例子
#iptables -t nat -A REDSOCKS -d 45.76.241.57 -j RETURN
#iptables -t nat -A REDSOCKS -d 45.32.79.211 -j RETURN
#iptables -t nat -A REDSOCKS -d 144.202.84.244 -j RETURN

# 过滤Docker container 网段 （默认的）
iptables -t nat -A REDSOCKS -d 172.17.0.0/16 -j RETURN
# 过滤局域网段
iptables -t nat -A REDSOCKS -d 10.0.0.0/8 -j RETURN
iptables -t nat -A REDSOCKS -d 127.0.0.0/8 -j RETURN
iptables -t nat -A REDSOCKS -d 169.254.0.0/16 -j RETURN
iptables -t nat -A REDSOCKS -d 172.16.0.0/12 -j RETURN
iptables -t nat -A REDSOCKS -d 192.168.0.0/16 -j RETURN
iptables -t nat -A REDSOCKS -d 224.0.0.0/4 -j RETURN
iptables -t nat -A REDSOCKS -d 240.0.0.0/4 -j RETURN
# enable 中国的ipset
iptables -t nat -A REDSOCKS -p tcp -m set --match-set chnroute dst -j RETURN
# FWD 需要走代理的情况，直接打到redsock的tcp port
iptables -t nat -A REDSOCKS -p tcp -j REDIRECT --to-ports 12345

# 让上面这个过滤用途的Chain 在PREROUTING中使用（本Linux系统作为 Gateway时候使用）
iptables -t nat -I PREROUTING -p tcp -j REDSOCKS
# 为本机的代理append chain
iptables -t nat -I OUTPUT -p tcp -j REDSOCKS
# 使用docker的时候FWD的DROP的据说是为了安全，但是这里，得给FWD了。不然转发机器无法上网
iptables -P FORWARD ACCEPT
```

## 注意事项

1. 如果你要用本机作为gateway别忘了做这个操作： https://linuxconfig.org/how-to-turn-on-off-ip-forwarding-in-linux
2. Redsocks 的 upstream socks5 要改成你自己的

