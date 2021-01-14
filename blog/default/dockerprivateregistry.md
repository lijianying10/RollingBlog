title: 小型团队内网创建docker registry
date: 2021-01-14 12:07:26
categories: 技术
tags: [docker,registry,ssl,nginx]
---

docker registry 是作为容器相关自动化关键底层基础设施，作为一种重要的artifact存储形式，对接CI/CD Kubernetes 都非常方便。

结构设计：

```

                                       ---> registry (run in docker container)
http traffic -> host nginx with ssl ---|
                                       ---> registry UI (run in docker container)


```

## 部署过程


### 1. 部署docker容器

```
docker run -d -p 5000:5000 --restart=always --name registry -v /data/registry:/var/lib/registry registry:2
docker run -d -p 5001:80 --name registry-ui -e DELETE_IMAGES=true joxit/docker-registry-ui:static
```

注意：

1. regisry 存储放到了 host 的 `/data/registry` 位置需要可以修改
2. ui 允许删除数据，更多选项参考：[这里](https://hub.docker.com/r/joxit/docker-registry-ui) 的Run the static interface 段落
3. 我快速的读了一下这个 registry-ui 的代码其实它是个静态网页项目，dockerfile 的内容显示其实里面就是个nginx 所以说明书上面的 `跨域` 和 `registry_url` 和 `SSL` 相关配置可以直接如下的 Nginx reverse proxy配置中直接干净快速的解决。

### 2. 部署Nginx

```
$ cat /etc/nginx/sites-enabled/registry
server {
        listen 443 ssl;
        server_name hub.philo.top;
        ssl_certificate     /etc/ssl/1_hub.philo.top_bundle.crt;
        ssl_certificate_key /etc/ssl/2_hub.philo.top.key;
        ssl_protocols       TLSv1 TLSv1.1 TLSv1.2;
        ssl_ciphers         HIGH:!aNULL:!MD5;
        client_max_body_size 2048M;
        location / {
                proxy_pass http://127.0.0.1:5001;
        }
        location /v2 {
                proxy_pass http://127.0.0.1:5000;
        }
}
```

注意：

1. 第一行是命令是提示文件存放位置
2. 安装nginx的方法是 `apt-get update && apt-get install -y nginx`
3. `server_name` 命令需要改成你自己的域名
4. 证书淘宝买 5 块钱 店铺名字 `鼎森网络科技有限公司` 因为是内网使用的证书所以用Let's encrypt比较麻烦。
5. 命令 `client_max_body_size` 不要裁剪，因为docker pull 和 push 的 http body 很大。
6. `location /` 的作用是路由到 ui container
7. `location /v2` 的作用是路由到 docker registry
8. 别忘了设置域名A记录
9. `SSL` 证书最好是要配置上的，原因是docker pull过程默认是要证书的不然需要特别配置trust同理 kubernetes 也有类似的需求稍微衡量一下5块钱还是值得的
