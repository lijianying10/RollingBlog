title: deploy docker registry under private network
date: 2021-02-15 15:30:35
categories: 技术
tags: [docker,registry,ssl,nginx]
---

Docker registry is an essential infrastructure of docker daemon or Kubernetes. We package project artifacts by docker image while storage and distribution by registry service. Today we will show you how we are setting up a straightforward and small registry implementation by docker official. It convenience a docking workflow for CI/CD.

Deploy structure:

![image](https://user-images.githubusercontent.com/3077762/107914469-76e1f300-6f9d-11eb-8b78-aaf7dad4c258.png)

## Deploy services

### Container

``` bash
docker run -d -p 5000:5000 --restart=always --name registry -v /data/registry:/var/lib/registry registry:2
docker run -d -p 5001:80 --name registry-ui -e DELETE_IMAGES=true joxit/docker-registry-ui:static
```

1. Replace path `/data/registry` to your own storage path
2. `-e DELETE_IMAGES=true` intends docker images can delete through UI operation
   1. Reference document by link [https://hub.docker.com/r/joxit/docker-registry-ui](https://hub.docker.com/r/joxit/docker-registry-ui)
   2. In paragraph `Run the static interface`
3. Code review image `joxit/docker-registry-ui:static` docker file we can know:
   1. The HTTP service is just an nginx process with a bunch of static HTTP static files.
   2. The `cross region` and `registry_url` and `SSL` config can be moved to our nginx deploy for more flexible and clean config management.

### Host nginx deploy

Generally, we install nginx by Linux package management such as apt.

We can install nginx under ubuntu by the command `sudo apt-get update && sudo apt-get install -y nginx`

Then we install the following config file under your config dir. The default path is `/etc/nginx/sites-enabled/`

Here we storage the config file in `/etc/nginx/sites-enabled/registry`

```
server {
    listen 443 ssl;
    server_name [[REPLACE: YOUR OWN DOMAIN NAME]];
    ssl_certificate     /etc/ssl/[[REPLACE: YOUR DOMAIN SSL CRT FILE]];
    ssl_certificate_key /etc/ssl/[[REPLACE: YOUR DOMAIN SSL KEY FILE]];
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

server{
    listen 80;
    server_name [[REPLACE: YOUR OWN DOMAIN NAME]];
    return 301 https://$host$request_uri;
}
```

`ATTENTION:` please replace the config with your environment situation.

1. The parameter `server_name` must replace with your domain name.
1. We recommend using `Let’s encrypt` DNS-01 challenge to verify your domain and get an SSL cert file.
1. The parameter `ssl_certificate` must replace with your domain crt file.
1. The parameter `ssl_certificate_key` must replace with your domain key file.
1. The parameter `client_max_body_size` at 2GB since we usually push a large docker image layer in practice.
2. `location /` route to registry UI container.
3. `location /v2` route to registry service.
4. Don't forget to set A record for your domain.
5. We highly recommend setting up nginx `HTTPS` for your service since the docker daemon or kubelet needs other configs to trust your registry.
6. The second `server` under the config file which helps us force switch from `HTTP` to `HTTPS`

`SAFETY WARNING:`

1. `Do not` deploy this solution in the public network.
2. Use it in a small team under a private network.


