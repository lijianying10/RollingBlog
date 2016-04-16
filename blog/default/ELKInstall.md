title: ELK安装与调试
date: 2016-03-30 11:06:19
categories: 技术
tags: [elk,elasticsearch,kibina,logstash,docker]
---

![](http://7viiaq.com1.z0.glb.clouddn.com/download.png)

## VERSION:1.1

## 更新

### 2016年04月16日17:30:24
1. ubuntu镜像更新
2. 添加百度网盘下载
3. 修正logstash 检查config的bug添加一个注意事项


## 环境Version

```
java 8u77
elasticsearch 2.2.1
logstash 2.2.2
kibana 4.4.2
```
### 快速下载通道：

[http://pan.baidu.com/s/1qYCFvvY](http://pan.baidu.com/s/1qYCFvvY)

## 背景：

随着业务量的增长对于日志的处理显得越来越重要。主要需求来源于各种数据的统计，数据的分析。以及对于可以用户的查询，系统故障的排查与监控。

本文为入门比较浅显讲述ELK的安装，不涉及到生产方面的内容，在未来还会更新更佳深入的实践或者线上实战的内容。

## 我当时学习的参考文档：

[elastic](https://www.elastic.co/guide/index.html)

在文档上官方详细程度足够，ELK文档质量好评，结合本文进行实践相信您一定会非常顺利的使用ELK。

## 方法论

通过脚本自动化构建docker容器。通过固定的脚本减少重复工作带来的辛苦。专注于技术研究。

## 环境依赖包下载

```
wget https://www.reucon.com/cdn/java/jdk-8u77-linux-x64.tar.gz
wget https://download.elasticsearch.org/elasticsearch/release/org/elasticsearch/distribution/tar/elasticsearch/2.2.1/elasticsearch-2.2.1.tar.gz
wget https://download.elastic.co/logstash/logstash/packages/debian/logstash_2.2.2-1_all.deb
wget https://download.elastic.co/kibana/kibana/kibana-4.4.2-linux-x64.tar.gz
tar xf elasticsearch-2.2.1.tar.gz
tar xf jdk-8u77-linux-x64.tar.gz
tar xf kibana-4.4.2-linux-x64.tar.gz
mkdir -p el ki lo
mv kibana-4.4.2-linux-x64 kibana
mv elasticsearch-2.2.1 elasticsearch
mv jdk1.8.0_77 jdk
mv elasticsearch el
cp -rf jdk el
mv jdk lo
mv kibana ki
mv logstash_2.2.2-1_all.deb lo
```

执行内容：

1. 下载安装包
2. 解压安装包
3. 整理ELK到固定的文件夹，用名字开头两个字母来区分
4. 吧解压好的安装包改了名字放到固定的位置上去。

其中：

1. kibana是不需要java环境的。
2. 我的java环境下载是通过CDN下载的效果还不错。
3. 被墙的很厉害。建议到国外的VPS上实践，我用的[Vultr](http://www.vultr.com/?ref=6876444) (连接有小尾巴，谢谢支持。)
4. 注意内存一定要大于1gb ELK启动后空载的内存使用为900mb左右。

## 准备好之后检查所有包的位置是否正确

``` shell
root@elk:~# ls el
elasticsearch  jdk
root@elk:~# ls ki
kibana
root@elk:~# ls lo
jdk  logstash_2.2.2-1_all.deb
```

## Dockerfile

Docker 可以让环境的构建的问题可以复现，构建好的环境可以分发。用来研究问题非常方便。

`注意：` 在本环境中所有的配置文件没有做成环境变量与气动脚本，因此变量都是固定的，尤其要注意网卡的创建`elastic`避免不必要的时间浪费。
在未来的研究中笔者将会构建生产中使用的ELK容器，会考虑扩展性与灵活性。在本例中只考虑安装特性研究配合看查看官方文档学习。


### Elasticsearch

``` shell
cat >> /root/el/Dockerfile << EOF
FROM ubuntu:14.04.4
ADD jdk /usr/local/jdk
ADD elasticsearch /usr/local/elasticsearch
ENV JAVA_HOME /usr/local/jdk
ENV PATH $PATH:/usr/local/jdk/bin:/usr/local/elasticsearch/bin/
RUN useradd -d /home/elasticsearch -m elasticsearch
RUN mkdir -p /data && echo "path.data: /data" >> /usr/local/elasticsearch/config/elasticsearch.yml && echo "network.host: 0.0.0.0" >> /usr/local/elasticsearch/config/elasticsearch.yml &&chown -R elasticsearch /usr/local/elasticsearch && chown -R elasticsearch /data
USER elasticsearch
CMD elasticsearch
EOF
```

请注意，由于Elasticsearch不可以使用root用户运行。因此需要针对data文件夹的挂载创建系统所需运行账户。
在本例中，使用如下命令创建用户。`useradd -d /home/elasticsearch -G root -s /bin/bash elasticsearch`

### kibana

```
cat >> /root/ki/Dockerfile << EOF
FROM ubuntu:14.04.4
ADD kibana /usr/local/kibana
RUN echo "port: 5601" >> /usr/local/kibana/config/kibana.yml && echo "host: 0.0.0.0" >> /usr/local/kibana/config/kibana.yml && echo "elasticsearch_url: http://elasticsearch.elastic:9200" >> /usr/local/kibana/config/kibana.yml
ENV PATH $PATH:/usr/local/kibana/bin/
CMD kibana
EOF
```


### Logstash

```
cat >> /root/lo/Dockerfile << EOF
FROM ubuntu:14.04.4
ADD jdk /usr/local/jdk
ADD logstash_2.2.2-1_all.deb /logstash_2.2.2-1_all.deb
ENV JAVA_HOME /usr/local/jdk
RUN dpkg -i /logstash_2.2.2-1_all.deb
ENV PATH $PATH:/opt/logstash/bin:/usr/local/jdk/bin
CMD logstash agent -f /conf
EOF
```

## 构建

```
docker build -t el:0 /root/el/
docker build -t ki:0 /root/ki/
docker build -t lo:0 /root/lo/
```


## 构建之后查看容器结果

```
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
lo                  0                   6592f152333d        3 minutes ago       763.8 MB
ki                  0                   fd93c06a3d4e        4 minutes ago       312.3 MB
el                  0                   4395211af2d0        4 minutes ago       619.5 MB
ubuntu              14.04.3             3876b81b5a81        9 weeks ago         187.9 MB
```

## logstash 测试用配置文件

```

input {
    file {
        path => "/syslog"
        start_position => beginning
        ignore_older => 0
    }

}

output {
    elasticsearch {
        hosts => "elasticsearch.elastic"
    }
    file {
        path => "/data/%{+yyyy/MM/dd/HH}/%{host}.log.gz"
        gzip => true
    }
}
```

在本例配置文件中，入口方式为监控syslog输出为elasticsearch 以及按照天与容器到压缩包。 方便备份。
配合上面的dockerfile请将上面的配置文件放到`/root/conf/`

## logstash 配置文件的检查。

由于配置文件格式复杂，并且运行之后不能轻易修改（技术不到位）因此在运行之前一定要做好检查工作。

```
docker run -it --rm -v /root/conf/:/conf -v /root/logstash/:/data/ --net=elastic lo:0 logstash -t -f /conf
```

### `注意检查中会检查文件夹权限，如果错误是不会通过的。所以需要先运行下面的elastic 和kibana 运行脚本`

## 运行ELK

```
useradd -d /home/elasticsearch -s /bin/bash elasticsearch
mkdir -p /root/data
chown -R elasticsearch /root/data
docker network create elastic
docker run -it -d --name elasticsearch -v /root/data:/data --net=elastic el:0
docker run -it -d --name kibana -p 5601:5601 --net=elastic ki:0
docker run -it -d --name logstash -v /root/conf/:/conf -p 12201:12201/udp -v /root/logstash/:/data/ -v /var/log/syslog:/syslog --net=elastic lo:0
```

1. 创建用户
2. 创建Elasticsearch所需的文件夹
3. 修改文件夹权限
4. 创建docker网卡
5. 运行Elasticsearch
6. 运行 kibana
7. 运行logstash


## 总结

由于时间有限，有一些问题没有能够及时补充发现，很抱歉。在实践中的确遇到了很多坑，尤其是对java虚拟机的不熟悉造成的一些问题。
占用了我不少时间，虽然查看文档需要占用很多时间，但是从中等时间的角度考虑是占优的十分推荐官方文档的阅读。
