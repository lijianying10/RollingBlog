title: 在容器中部署博客
date: 2015-06-13 21:59:09
categories: 技术
tags: [golang,web,git,busybox]
---


## 意义
Pages服务更新速度慢，使用DaoCloud免费配可以快速更新Blog。

## 工作流
1. Hexo生成静态文件
2. 按照正常流程Deploy
3. WebHook到DaoCloud中的目标容器
4. 目标容器更新部署好的博客
5. 静态文件服务器同时也做好更新

## 准备工作
### 容器快速调试脚本
```
docker stop busygit
docker rm busygit
docker rmi busyboxgit
# 清理运行环境

# 重新构建环境
docker build -t busyboxgit /root/busyboxgit
docker run -it -d --name busygit -p 8080:8080 busyboxgit
docker exec -it busygit /bin/sh
```

### busybox构建容器使用的操作系统
```
rm -rf /rootfs
mkdir /rootfs #新建根
cd /rootfs
mkdir bin etc dev dev/pts lib proc sys tmp #根据根的格式创建目录
touch etc/resolv.conf # DNS服务器地址
cp /etc/nsswitch.conf etc/nsswitch.conf #名字解析配置
echo root:x:0:0:root:/:/bin/sh > etc/passwd # 新建用户
echo root:x:0: > etc/group #新建用户组
ln -s lib lib64 # 创建软连接
ln -s bin sbin 
cp /bin/busybox bin # 复制bin文件，不然没有命令用
cp $(which git) bin/
$(which busybox) --install -s bin #同上两步操作，用which是因为需要绝对路径
bash -c "cp /lib/x86_64-linux-gnu/lib{c,z,pcre,m,dl,crypt,util,rt,nsl,nss_*,pthread,resolv}.so.* lib" # 这一步很重要，如果使用了Go调用外部动态库这里需要复制进去，
bash -c "cp /lib/x86_64-linux-gnu/ld* lib" #同上
bash -c "mkdir -p usr/lib/git-core/" #添加Git依赖
bash -c "cp -rf /usr/lib/git-core/* usr/lib/git-core/" #添加Git依赖

cp /lib/x86_64-linux-gnu/ld-linux-x86-64.so.2 lib #同上
tar cf /rootfs.tar . #根目录打包
for X in console null ptmx random stdin stdout stderr tty urandom zero ; do tar uf /rootfs.tar -C/ ./dev/$X ; done 
mv /rootfs.tar /root/busyboxgit/rootfs.tar
```
此段Fork于[这里](http://philo.top/2015/06/06/BuildAndDeplyTheThinnestDockerImage/)

### Golang静态网页服务，以及自动更新源码
``` golang
package main

import (
	"net/http"
	"os"
	"os/exec"
	"time"
)
import "fmt"

var LastTime int64 // FREQ control

func main() {
	fmt.Println("Server start") // For log
	http.Handle("/", http.FileServer(http.Dir("/public/"))) // Static file service 
	http.HandleFunc("/manualupdate", func(w http.ResponseWriter, r *http.Request) {
		if time.Now().Unix() == LastTime {
			w.Write([]byte("Over Frequency"))
			return
		}
		LastTime = time.Now().Unix()
		go func() {
			cmd := "git"
			args := []string{"-C", "/public", "pull"}
			if err := exec.Command(cmd, args...).Run(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			fmt.Println("Update Finished")
		}()
		w.Write([]byte("Task is running"))
	})
	http.ListenAndServe(":8080", nil)
}
```

### DockerFile For DaoCloud
```
FROM scratch
MAINTAINER Jianying Li <lijianying12@gmail.com>

ADD ./rootfs.tar / 
RUN mkdir -p /app
WORKDIR /app
COPY ./static /app/
RUN git clone git://gitcafe.com/lijianying12/lijianying12.git /public
EXPOSE 8080

CMD ["/app/static"]
```

## 快速上手
通过本方式部署静态博客到DaoCloud请执行如下步骤:

 - Fork [https://github.com/lijianying10/HexoCan/](https://github.com/lijianying10/HexoCan/)  修改第八行DockerFile到自己部署的Git地址。`注意我裁剪的系统只支持Git协议，可以部署在GitCafe或者GitHub上`
 - 在DaoCloud中正常构建部署。
 - 添加WebHooks在配置中添加地址： http://XXX.daoapp.io/manualupdate

每次静态博客更新之后都会自动调用WebHook来自动在容器中更新到最新版本的博客。


##总结
0. 操作系统+golang运行文件+整个blog=37mb 构建部署速度非常快
1. 更换服务之后更新时间以及速度可以由自己把握，非常方便。
2. 更加灵活，可以开发自己的接口实现更多有意义的功能（搜索）。
3. 功能上更多的想象空间


