title: "RancherOS iPXE 启动以及ISO重新打包的方法"
date: 2016-06-11 22:32:57
tags: [rancher,docker,linux,iPXE,PXE,syslinux,virtualbox]
---

### VER: 1

## 意义
最近一直在尝试自己写一服务器平台方便自己和别人开发的，能把流程做进去最好了。
但是万事开头难，在尝试了CoreOS iPXE 启动方法之后发现Image太大大概要200mb+的样子。
Boot2Docker始终没找到正确的看文档的方式，这种方法我个人预估也只能是Docker ToolBox使用比较靠谱了。
然后开始尝试RancherOS的IPXE启动操作系统但是但是很不稳定，有的时候Cloud-init下载不了。
大概2次就有一次无法启动不能忍，但是不用cloud-init还好，但是也是去了在我这里使用的意义。
这种只适合大规模的集群部署。不得不说他们的工作真的很赞。

那么要满足我的需求：

1. Image 要小
2. 开机迅速，可定制。各种系统配置。
3. 稳定不容易出错。
4. 自动持久化。 
5. 可以持续集成打包。

那么就只有一个路线syslinux使用`mkisofs`这个工具重新打包RancherOS最后我自己打包的结果30mb 比官方的还小1mb还可以对docker 进行自动持久化。方便在不同系统种分发Docker daemon

## iPXE 启动的方法：

VirtualBox参考启动：

```
VBoxManage createvm --name "service" --register
VBoxManage modifyvm "service" --memory 512 --acpi on --boot1 dvd
VBoxManage modifyvm "service" --nic1 hostonly --hostonlyadapter1 vboxnet0 --nicpromisc1 allow-all
VBoxManage modifyvm "service" --nic2 nat  --natnet2 "192.168/16" --natpf2 "guestssh,tcp,,2222,,22" --nicpromisc2 allow-all
VBoxManage modifyvm "service" --ostype Linux

VBoxManage createhd --filename /mnt/vdis/io.vdi --size 10000
VBoxManage storagectl "service" --name "IDE Controller" --add ide
VBoxManage storageattach "service" --storagectl "IDE Controller"  \
    --port 0 --device 0 --type hdd --medium /mnt/vdis/io.vdi
VBoxManage storageattach "service" --storagectl "IDE Controller" \
    --port 1 --device 0 --type dvddrive --medium /mnt/git/ipxe/src/bin/ipxe.iso
```

因为总是需要重新启动调试所以这里把所有的VBOX操作都自动化。提升调试速度。

### ipxe.iso 重新打包

因为要针对我们的using case使用，我们需要自动化boot一个固定的操作系统
所以这里要针对我们的脚本重新打包。


iso重新打包命令如下
```
git clone git://git.ipxe.org/ipxe.git
cd ipxe/src
make bin/ipxe.iso EMBED=bootscript
```

打包依赖：
```
gcc (version 3 or later)
binutils (version 2.18 or later)
make
perl
syslinux (for isolinux, only needed for building .iso images)
liblzma or xz header files
```

bootscript 的内容：

```
#!ipxe

dhcp
chain http://192.168.56.1:8089/boot.html
```
注意： boot的时候请求的整个地址返回的内容为启动脚本。

启动脚本内容：
```
#!ipxe
# Boot a persistent RancherOS to RAM

# Location of Kernel/Initrd images
set base-url http://192.168.56.1:8089

kernel ${base-url}/vmlinuz044 rancher.state.formatzero=true rancher.state.autoformat=[/dev/sda] rancher.cloud_init.datasources=['url:http://192.168.56.1:8089/rancheros.yml']
initrd ${base-url}/initrd044
boot
```

然后就可以启动RancherOS了。

启动流程说明： `虚拟机启动->启动ipxe DVD->Chain load 请求 boot.html -> iPXE 根据boot.html下载 vmlinuz initrd -> 带入参数启动Linux（其实PE之类的也都好使）`

参考文档：

 - [代码下载以及编译参考](http://ipxe.org/download)
 - [iso重新打包,嵌入脚本参考](http://ipxe.org/embed)
 - [所有的BuildCFG参考](http://ipxe.org/buildcfg)

## 重新打包RancherOS 系统image

虽然官方提供了[Release下载](https://github.com/rancher/os/releases)但是这里面的iso对于我的使用场景来说还远远不够。
首先默认密码不能是固定的rancher，另外Docker 自动持久化的部分并没有做，另外就是Docker启动参数我不想用默认的我需要有一定的调整。这些原因通过iPXE.iso重新打包的方法学习之后促成了这一想法虽然简单，但是kernel参数的调整的确花费了我至少2-3小时的时间。为了避免自己再次踩坑来写一些说明。

### VBOX 调试脚本

```
VBoxManage createvm --name "service" --register
VBoxManage modifyvm "service" --memory 512 --acpi on --boot1 dvd
VBoxManage modifyvm "service" --nic1 hostonly --hostonlyadapter1 vboxnet0 --nicpromisc1 allow-all
VBoxManage modifyvm "service" --nic2 nat  --natnet2 "192.168/16" --natpf2 "guestssh,tcp,,2222,,22" --nicpromisc2 allow-all
VBoxManage modifyvm "service" --ostype Linux

VBoxManage createhd --filename /mnt/vdis/io.vdi --size 10000
VBoxManage storagectl "service" --name "IDE Controller" --add ide
VBoxManage storageattach "service" --storagectl "IDE Controller"  \
    --port 0 --device 0 --type hdd --medium /mnt/vdis/io.vdi
VBoxManage storageattach "service" --storagectl "IDE Controller" \
    --port 1 --device 0 --type dvddrive --medium /mnt/git/rancheros.iso
```

### 打包命令(build.sh)：

```
mkisofs -o ros.iso -b isolinux/isolinux.bin -c isolinux/boot.cat -no-emul-boot -boot-load-size 4 -boot-info-table rosiso/
```

### 打包后结果
```
├── build.sh
├── rosiso
│   ├── initrd
│   ├── isolinux
│   │   └── isolinux.bin
│   ├── isolinux.cfg
│   └── vmlinuz
└── ros.iso
```

文件来源： 

 - initrd 以及vmlinuz 从rancheros release 页面获取
 - isolinux.bin 如果是ubuntu文件的位置在： `/usr/lib/syslinux/isolinux.bin`
 - boot.cat 是自动生成的
 - ros.iso 生成的系统镜像

### isolinux.cfg 启动参数定制

```
default rancheros
label rancheros
    kernel /vmlinuz
    initrd /initrd
    append quiet rancher.password=rancher rancher.state.autoformat=[/dev/sda] rancher.state.formatzero=true 
```

参数解释：（从quiet之后开始）

 - rancher.password 设置默认账户密码。
 - rancher.state.autoformat 自动格式化磁盘。
 - Starts with 1 megabyte of zeros and `rancher.state.formatzero` is true.
 - 参数顺序会影响系统稳定性。
 - 所有参数查看需要等启动之后sudo ros c export -f > cfg.yml 几乎所有参数都可以调整，另外单独一个参数里面不能有空格kernel的参数检查比cloud-init模块本身严格多了。
 - 上面这种方式为cmdline[操作系统源码](https://github.com/rancher/os/blob/master/cmd/cloudinit/cloudinit.go) 在GetDatasources 函数这里。

### 参考文档:

 - [ISOLINUX](http://www.syslinux.org/wiki/index.php?title=ISOLINUX)
 - [对我就是从RancherOS官方启动参考](https://github.com/rancher/os/blob/master/scripts/isolinux.cfg)

## 总结

 - 新打包好的系统大小30MB
 - 启动时间10s左右（机械硬盘）
 - 强制关机，重新启动10次以上无错误。每次测试都有不同的小操作。

