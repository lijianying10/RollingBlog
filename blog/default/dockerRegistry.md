title: Docker 私有镜像最基本调试
date: 2016-04-21 23:28:09
categories: 技术
tags: [nginx,docker,registry]
---


## VERSION:1

## 背景
Registry 作为Docker的重要基本组件，无论是在自己的内网搭建，或者是公网私有带宽搭建都是托管私有Image的好方法。


## Step1 Docker Image准备

```
docker pull index.alauda.cn/library/registry:2.3
docker pull index.alauda.cn/library/nginx
docker tag index.alauda.cn/library/registry:2.3 registry:2.3
docker tag index.alauda.cn/library/nginx:latest nginx:latest
```

## Step2 Nginx 配置文件准备

```
docker run -it --rm -v /root/:/data/ nginx cp -rf /etc/nginx /data/ngreg
cat > /root/ngreg/conf.d/reg.conf << EOF
server {
    listen       443;
    server_name XXX.XXX.XXX;

    ssl on;
    ssl_certificate /etc/nginx/XXX.XXX.XXX/fullchain1.pem;
    ssl_certificate_key /etc/nginx/XXX.XXX.XXX/privkey1.pem;
    add_header 'Docker-Distribution-Api-Version' 'registry/2.0' always;
    client_max_body_size 0;
    chunked_transfer_encoding on;

    location / {

        if ($http_user_agent ~ "^(docker\/1\.(3|4|5(?!\.[0-9]-dev))|Go ).*$" ) {
          return 404;
        }

        proxy_pass                          http://reg.prod:5000;
        proxy_set_header  Host              $http_host;   # required for docker client's sake
        proxy_set_header  X-Real-IP         $remote_addr; # pass on real client's IP
        proxy_set_header  X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header  X-Forwarded-Proto $scheme;
        proxy_read_timeout                  900;
  }
}
EOF
```

在上面的脚本重需要修改的是： 
1. XXX.XXX.XXX 需要修改成您自己的域名
2. 如果需要申请证书请参考：[http://www.philo.top/2016/04/06/letsencryptSSLApply/](http://www.philo.top/2016/04/06/letsencryptSSLApply/)

配置解释：
1. 默认使用HTTPS协议。
2. header添加指定registry版本。
3. 设置body长度为0无限制长度。
4. 拒绝低于docker1.6版本的访问。
5. 反向代理到docker registry 容器。
6. 其他Headers都是参考文档得到的。

## Step3 简单账号配置：

### 配置修改

```
        auth_basic "Restricted";
        auth_basic_user_file /etc/nginx/htpasswd;
```

上面这两段，放到location下面就可以了。

账号密码文件[生成工具](http://tool.oschina.net/htpasswd)

## Step4 开始服务

```
docker network create prod
docker run -it -d --name ngreg --net=prod -v /root/ngreg/:/etc/nginx -p 0.0.0.0:443:443 nginx
docker run -it -d --name reg --net=prod -v /data:/regdata registry:2.3
```

## 总结

1. 需要根据自己的需要来进行配置文件的修改。
2. 次文档为最低成本限度的registry部署。

之后要面临的挑战

1. 使用registry2.4进行GC
2. auth上微服务。
3. 存储管理。(OSS,UFile,等等这些产品。)
4. Public Private 控制。
