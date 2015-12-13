title: "RunC v1 尝鲜"
date: 2015-07-17 16:52:40
categories: 技术
tags: [RunC,docker,linux,golang,hexo]
---

## 背景
RunC今天凌晨终于发了第一个Release,git version：0.0.1 , Code version 0.2。不过这里官方应该是笔误了。在代码查询中版本为0.2 但是在Git tag中使用的版本为0.0.1。轻量级LibContainer 未来会支持Windows（你敢信？）。RunC是OCP的产物，在DockerCon 2015中亮相其官网为：[http://runc.io/](http://runc.io/)它与Docker的不同在于它不需要守护进程，只需要一个配置文件一个开放的文件夹。一个进程就可以跑。另外它的依赖很少
```
ldd /bin/runc
	linux-vdso.so.1 (0x00007fffab589000)
	libpthread.so.0 => /lib64/libpthread.so.0 (0x00007ff714366000)
	libc.so.6 => /lib64/libc.so.6 (0x00007ff713fbb000)
	/lib64/ld-linux-x86-64.so.2 (0x00007ff714583000)
```
可以直接复制文件到Busybox或者puppy 这种特别小的Linux上就可以运行。

## 意义
对于我来说：
1. 我需要运行那么多环境，我可以不用docker来管理了。我只需要管理我的进程即可。
2. 对于脚本开发更爽了。直接在文件夹里面替换脚本（php python等）重启容器，部署就完成了。
3. 对于我的blog hexo这种运行一下就ok的情况。更是方便，blog源码可以直接同步。

下面我们一步一步来构建hexo环境

## 构建
首先我的环境在昨天的Blog中有所描述[http://www.philo.top/2015/07/16/pc-docker/](http://www.philo.top/2015/07/16/pc-docker/)
```
docker pull golang:1.4.2 # 下载golang环境
docker run -it -v /mnt/:/aaa/ golang:1.4.2 /bin/bash
```
进入golang环境之后开始准备构建：
1. 本环境中 `/go` 为GOPATH
2. git已经有了
3. 基于Debian
4. 上面第二条命令中-v是为了把编译出来的文件导出用的
5. 打开之后直接进入bash
6. 进入后PWD为`/go`注意下面的命令初始化PWD也为`/go`(就是进去了不要乱动)

```
mkdir -p src/github.com/opencontainers #创建目录
cd src/github.com/opencontainers 
git clone https://github.com/opencontainers/runc #下载源码
cd runc 
make # 构建源码
cp ./runc /aaa/ # 导出目标
```

编译生成好了之后，直接`Ctrl+d`退出golang容器即可。
然后就可以在/mnt/runc 中找到编译好的RunC
我这里是根据我的情况来完成编译部署的，您可以更换其他target文件夹从/mnt中更换到其他地方。

## 试运行
```
 /mnt/runc --help
NAME:
   runc - Open Container Project runtime

runc is a command line client for running applications packaged according to the Open Container Format (OCF) and is
a compliant implementation of the Open Container Project specification.

runc integrates well with existing process supervisors to provide a production container runtime environment for
applications. It can be used with your existing process monitoring tools and the container will be spawned as a direct
child of the process supervisor.

After creating a spec for your root filesystem with runc, you can execute a simple container in your shell by running:

    cd /mycontainer
    runc


USAGE:
   runc [global options] command [command options] [arguments...]

VERSION:
   0.2

COMMANDS:
   checkpoint	checkpoint a running container
   events	display container events such as OOM notifications and cpu, memeory, IO, and network stats
   restore	restore a container from a previous checkpoint
   spec		create a new specification file
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --id "root"			specify the ID to be used for the container
   --debug			enable debug output for logging
   --root "/var/run/ocf"	root directory for storage of container state (this should be located in tmpfs)
   --criu "criu"		path to the criu binary used for checkpoint and restore
   --help, -h			show help
   --version, -v		print the version
```

## 尝试运行hexo

### 坑介绍：
1. 权限问题
2. RunC文档不完善

下面开始爬坑
```
docker pull lijianying10/hexo:3.0.0 # 下载我曾经做的hexo image
docker export $(docker create lijianying10/hexo:3.0.0) >hexo.tar #导出image到tar包中
mkdir rootfs # 容器的根目录
tar -C rootfs -xf nginx.tar # 解压打包文件到根目录
/mnt/runc spec > config.json # 给runc生成配置文件
```

修改配置文件：config.json
args sh->bash 为了方便用shell
readonly true->false 为了log之类的
cwd 设定为`/hexo` (初始化目录很方便)
namespace中删除整个
	```
	                        {
	                                "type": "network",
	                                "path": ""
	                        },
	```
	这样容器就与host共享一个网络堆栈了。（外网可以访问，因为说明书（几乎没有）没有找映射相关的所以只能先这样）


这里贴出完整Config.json
```
{
	"version": "pre-draft",
	"platform": {
		"os": "linux",
		"arch": "amd64"
	},
	"process": {
		"terminal": true,
		"user": {
			"uid": 0,
			"gid": 0,
			"additionalGids": null
		},
		"args": [
			"bash"
		],
		"env": [
			"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
			"TERM=xterm"
		],
		"cwd": "/hexo"
	},
	"root": {
		"path": "rootfs",
		"readonly": false
	},
	"hostname": "shell",
	"mounts": [
		{
			"type": "bind",
			"source": "/mnt/blogWork/blog",
			"destination": "/hexo",
			"options": "rbind,rw"
		},
		{
			"type": "proc",
			"source": "proc",
			"destination": "/proc",
			"options": ""
		},
		{
			"type": "tmpfs",
			"source": "tmpfs",
			"destination": "/dev",
			"options": "nosuid,strictatime,mode=755,size=65536k"
		},
		{
			"type": "devpts",
			"source": "devpts",
			"destination": "/dev/pts",
			"options": "nosuid,noexec,newinstance,ptmxmode=0666,mode=0620,gid=5"
		},
		{
			"type": "tmpfs",
			"source": "shm",
			"destination": "/dev/shm",
			"options": "nosuid,noexec,nodev,mode=1777,size=65536k"
		},
		{
			"type": "mqueue",
			"source": "mqueue",
			"destination": "/dev/mqueue",
			"options": "nosuid,noexec,nodev"
		},
		{
			"type": "sysfs",
			"source": "sysfs",
			"destination": "/sys",
			"options": "nosuid,noexec,nodev"
		},
		{
			"type": "cgroup",
			"source": "cgroup",
			"destination": "/sys/fs/cgroup",
			"options": "nosuid,noexec,nodev,relatime,ro"
		}
	],
	"linux": {
		"uidMapping": null,
		"gidMapping": null,
		"rlimits": null,
		"sysctl": null,
		"resources": {
			"disableOOMKiller": false,
			"memory": {
				"limit": 0,
				"reservation": 0,
				"swap": 0,
				"kernel": 0,
				"swappiness": -1
			},
			"cpu": {
				"shares": 0,
				"quota": 0,
				"period": 0,
				"realtimeRuntime": 0,
				"realtimePeriod": 0,
				"cpus": "",
				"mems": ""
			},
			"blockIO": {
				"blkioWeight": 0,
				"blkioWeightDevice": "",
				"blkioThrottleReadBpsDevice": "",
				"blkioThrottleWriteBpsDevice": "",
				"blkioThrottleReadIopsDevice": "",
				"blkioThrottleWriteIopsDevice": ""
			},
			"hugepageLimits": null,
			"network": {
				"classId": "",
				"priorities": null
			}
		},
		"namespaces": [
			{
				"type": "process",
				"path": ""
			},

			{
				"type": "ipc",
				"path": ""
			},
			{
				"type": "uts",
				"path": ""
			},
			{
				"type": "mount",
				"path": ""
			}
		],
		"capabilities": [
			"AUDIT_WRITE",
			"KILL",
			"NET_BIND_SERVICE"
		],
		"devices": [
			"null",
			"random",
			"full",
			"tty",
			"zero",
			"urandom"
		]
	}
}
```

## 目录挂载
虽然没有说明文档但是阅读了大量代码之后发现
```
		{
			"type": "bind",
			"source": "/root/blog",
			"destination": "/hexo",
			"options": "rbind,rw"
		},
```


虽然文档  [https://github.com/opencontainers/runc/blob/master/Godeps/_workspace/src/github.com/opencontainers/specs/config.md](https://github.com/opencontainers/runc/blob/master/Godeps/_workspace/src/github.com/opencontainers/specs/config.md)    中提到了挂载的方法。但是依然有好使又有不好使的情况。

注意权限正常情况的情况：
```
root@shell:/hexo# ls -all
total 24
drwxr-xr-x  8 root root  280 Jul 17 23:22 .
drwxr-xr-x 22 root root  480 Jul 17 23:31 ..
drwxr-xr-x 16 root root  360 Jul 17 23:22 .deploy_git
-rw-r--r--  1 root root   65 Jul 17 23:22 .gitignore
-rw-r--r--  1 root root 2059 Jul 17 23:22 _config.yml
-rw-r--r--  1 root root  174 Jul 17 23:22 db.json
-rw-r--r--  1 root root 2056 Jul 17 23:22 debug.log
drwxr-xr-x 18 root root  360 Jul 17 23:22 node_modules
-rw-r--r--  1 root root  518 Jul 17 23:22 package.json
drwxr-xr-x 15 root root  340 Jul 17 23:22 public
drwxr-xr-x  2 root root  100 Jul 17 23:22 scaffolds
drwxr-xr-x  5 root root  120 Jul 17 23:22 source
-rw-r--r--  1 root root   65 Jul 17 23:22 synck.sh
drwxr-xr-x  4 root root   80 Jul 17 23:22 themes
```

`我遇到的坑： 我从rsync同步的blog内容因为是带着权限的导致权限错误。`

```
runc # 运行runc 
hexo g # 开始生成
hexo s # 试运行
```
运行runc时它会自动检查config.json

自此结束该tag版本的测试

## 总结
目前还不能用作生产
1. 文档几乎没有。
2. 如果想实践请查找所有相关.md文档，包括它的依赖。
3. 这个项目引用了docker源代码，所以工程复杂度没有那么简单
4. 挂载还是有诡异的问题我甚至找到了源码文件`Godeps/_workspace/src/github.com/docker/docker/pkg/mount/flags.go`试过了所有的flag都不能解决这个问题。
5. 挂载这个flag有点类似mount -o
6. 但是调试好了。真心比docker方便，只是资源调度上，等一些细节问题还有待提升。
