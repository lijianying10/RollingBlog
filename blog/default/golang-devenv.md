title: golang 环境配置建议
date: 2016-04-05 13:39:00
categories: 技术
tags: [golang,dev,env,vim,docker]
---

## VERSION:1.3

![](http://golang.org/doc/gopher/project.png)

## 摘要
[在之前的实践中](/2015/02/06/golang-环境配置建议/)满足开发环境所有特征的情况下进行了大量方式上的升级。

经过`8`次的版本升级，调整，爬坑，终于觉得这次版本升级足够有意义替代之前版本的建议。

我总觉得花一少部分的时间调整开发工具让自己的工作效率更高是非常值得的。

本次祭出大招`Docker`来作为项目开发的主要脚手架。

并且我现在认为一个好的开发工具应该满足：

1. 项目管理
1. 快速文件跳转
1. 自动语法检查
1. 自动补全
1. 查找定义
1. 启动速度快
1. 如果自己有需求的话插件可以随便写
1. 灵活的运行配置

但是根据10个月的开发实践，除了上面这些之外还需要兼顾：

1. 国内的网络环境
2. 升级方便
3. 随时随地快速部署自己的开发环境
4. 能在线上服务器进行开发，随时随地升级版本修改代码。

## Docker Registry

`docker pull index.tenxcloud.com/philo/golangdev:1.3`

`docker pull lijianying10/golangdev:1.3`

如果你还不会docker请参考[这里](/2015/04/01/docker-super-start/)

## Dockerfile

[点击这里查看Dockerfile](https://github.com/lijianying10/FixLinux/blob/master/golangdev/Dockerfile)

## 升级日志

### 1.2.1

```
修复vim下escape有时间延迟（解决方法见.vimrc最后一行）
Ctrl+s保存的时候添加代码格式检查（其实就是追加命令 :GoMetaLinter）
更新golang到1.5.2 based on Debian 8
```

### 1.3 (2016年04月04日)

1. 环境升级：`debian 8`, `golang 1.6`
2. VIM插件升级 : powerline-> airline 修复字体问题
3. 修复VIM+TMUX背景颜色不一致的问题。
4. BASHRC添加了手动下载最新环境变量的alias

注意：如果您想解决乱码问题需要下载[PowerLine字体](https://github.com/powerline/fonts)设置Term软件到这里面的字体就可以了。常用的编程字体里面都有
如果您不想用PowerLine字体请注释掉：`let g:airline_powerline_fonts = 1` 此行代码位置在`~/.vimrc` 


## 特征解释

### 兼顾国内网络情况

1. 使用Dockerfile从国外VPS构建，然后推送到时速云备用。这种构建方式适合调试。
2. `推荐！`如果您在国外没有VPS推荐使用时速云TCE来构建，从香港节点自动化构建随时能看到日志。[参考文档](http://doc.tenxcloud.com/doc/v1/ci/client-download.html)

### 升级方便

1. 直接修改Dockerfile完成升级，调整From就可以调整底层系统使用。
2. 可根据您的需要随时定制自己的版本。非常方便。

### 随时能够快速部署

1. Docker启动速度非常快。
2. 如果您没有Image在内网该Image也只有1GB大小可非常快速的传输到您的电脑。
3. 国内准备好了加速源，详细查看Docker Registry部分。时速云确实挺快的。
4. 如果您的工作站安装的是CoreOS or Rancher这种的Docker Linux 不但安装快，部署开发环境也是一瞬间完成。

### 能够在线上服务器进行线上代码调整

1. 只要部署到线上服务器直接就可以使用。
2. 老板再也不用担心我的集成新功能速度太慢了。


## 使用方法

### 文件跳转([Command-T](https://wincent.com/products/command-t))

快捷键： `<leader>t`

注意：`<leader>`在我的vim配置里面是反斜杠,插件快捷键参考官方文档。

![](http://7viiaq.com1.z0.glb.clouddn.com/QQ20151213-0@2x.png)

### 文件管理(NERD_tree)

快捷键： `M-u`。

注意： 插件快捷键参考官方文档。

![](http://7viiaq.com1.z0.glb.clouddn.com/QQ20151213-1@2x.png)

### 自动语法检查

触发： 每次保存文件。

命令： `:GoMetaLinter, which invokes all possible linters (golint, vet, errcheck, deadcode, etc..) and shows the warnings/errors`

![](http://7viiaq.com1.z0.glb.clouddn.com/QQ20151213-2@2x.png)

例子中：Struct默认要求有注释，不然就会报警。对于常用的拼写比如说ID有严格的检查需要符合大众的拼写习惯。

其他正确性检查这里不再赘述。

### 自动补全

![](http://7viiaq.com1.z0.glb.clouddn.com/QQ20151213-3@2x.png)

注意：在最上面会显示API文档,想关闭文档快捷键：`M-c` 。

### 查找定义位置

快捷键：`M-n`。

### Outline 快速跳转(Tagbar)

快捷键：`M-p`。

![](http://7viiaq.com1.z0.glb.clouddn.com/QQ20151213-4@2x.png)

### 快速关闭文件

快捷键：`C-c`。


### 文件标签切换

快捷键： `M-i` 切换到上一个。
快捷键： `M-o` 切换到下一个。

### 保存文件

快捷键 `C-s`。

注意： Stop tty已经被我关闭，不必担心tty被锁。

## 总结

在上面我总结的全部都是我加的快捷键用起来比较舒服的。如果您需要自己修改快捷键请fork[我的github REPO](https://github.com/lijianying10/FixLinux/blob/master/dotfile/.vimrc)。

如果您有任何改进意见请回复留言给我发Email。先谢过。

vim还有很多默认的快捷键这里就不再多说，最好的学习方法是想到自己有什么习惯或者需要快捷键支持去google找找。
