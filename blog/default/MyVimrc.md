title: "我的vim golang 开发环境"
date: 2015-07-27 09:25:57
categories: 技术
tags: [vim,linux,macos]
---

## 意义
最近一直蛋疼图形界面实在不想因为编辑器的问题使用图形界面。太不简洁看着难受。
最后决定在Docker中配置我自己的Golnag开发环境为例子做个教程。

## 开始
1. 在Mac上使用brew安装vim`brew install vim` 因为10.10.4 默认的vim竟然是7.3的不能忍。
2. 修改path能够直接打开7.4的不要如果需要lua支持不要忘了`--with-lua`
3. 安装vundle `mkdir -p ~/.vim/bundle && git clone https://github.com/gmarik/Vundle.vim.git ~/.vim/bundle/Vundle.vim`
4. 安装vimrc `curl -Ssl https://raw.githubusercontent.com/lijianying10/FixLinux/master/vimrc -o ~/.vimrc`
5. 安装配色`mkdir ~/.vim/colors/ && curl -Ssl https://raw.githubusercontent.com/tomasr/molokai/master/colors/molokai.vim -o ~/.vim/colors/molokai.vim`
6. 安装插件： vim下：`:PluginInstall` 可以自动补全。

经过如上配置基本上vim环境就已经成型。

## 快捷键
1. F7: NERDtree 内部快捷键：
```
和编辑文件一样，通过h j k l移动光标定位
o 打开关闭文件或者目录，如果是文件的话，光标出现在打开的文件中
go 效果同上，不过光标保持在文件目录里，类似预览文件内容的功能
i和s可以水平分割或纵向分割窗口打开文件，前面加g类似go的功能
t 在标签页中打开
T 在后台标签页中打开
p 到上层目录
P 到根目录
K 到同目录第一个节点
J 到同目录最后一个节点
m 显示文件系统菜单（添加、删除、移动操作）
? 帮助
q 关闭
```

2. F8: Outline

3. Leader+t 快速打开文件自动匹配哦。 一般都是工程的时候用就ok了。
安装依赖ruby还有gcc
```
cd ~/.vim/bundle/command-t/ruby/command-t
ruby extconf.rb
make
```
使用
```
ctrl+j/k 上下选择文件，选中后回车打开文件
ctrl+t 以tab方式打开文件
ctrl+s/v 可以水平或垂直分割窗口打开文件
ctrl+c 退出该模式
```

4. 多tab情况翻页： Ctrl+0 向后 Ctrl+9 向前 
```
:tabnew [++opt选项] ［＋cmd］ 文件            建立对指定文件新的tab
  :tabc       关闭当前的tab
  :tabo       关闭所有其他的tab
  :tabs       查看所有打开的tab
  :tabp      前一个
  :tabn      后一个
```

5. 跳转的向前向后(比如说GoDef跳转到其他地方想跳转回来)：Ctrl-o 向前：Ctrl-i

## 部署在Docker中的坑：
1. Command-T 依赖GCC以及Ruby 和Ruby-dev
2. 调整Docker从8色到156色
3. 调整Docker支持UTF-8
4. PowerLine字体支持。
其他的从Dockerfile中一点一点解决。因为文章中环境的构建大量依赖Github以及golang官方网站，因此这里我们推荐使用境外服务器构建，Dockerfile 完整参考：[https://github.com/lijianying10/FixLinux/blob/master/golangdev/Dockerfile](https://github.com/lijianying10/FixLinux/blob/master/golangdev/Dockerfile)在连接中。
详细解释： 
1. `apt-get update` 准备安装各种依赖
2. build-essential 各种编译器，curl 下载配置文件用，git 平时写代码的时候用，vim-nox 此vim有+ryby +lua 所以选他安装，ctags  是tagbar的依赖。其他的不是重点，但都不可少。
3. 第6-9行都是从github上下载各种包还有配置文件放到容器中。
4. 本句命令`vim "+PluginInstall" "+GoInstallBinaries" "+qall"` 是打开vim执行两个命令，然后退出，第一个命令为下载vimrc中描述的插件，第二个命令为自动下载golnag开发工具包。
5. 所有go get命令都为常用golnag开发时用的包。
6. `echo "en_US.UTF-8 UTF-8" > /etc/locale.gen && locale-gen "en_US.UTF-8"`  编译本地定义文件的一个列表,如此才能支持utf-8
7. 字体安装`mkdir ~/.font/ && cd ~/.font/ && git clone https://github.com/eugeii/consolas-powerline-vim.git && cd consolas-powerline-vim/ && cp *.ttf .. && cd .. && rm -rf consolas-powerline-vim/ && mkfontscale && mkfontdir && fc-cache -vf`
8. Command-T 依赖安装`cd ~/.vim/bundle/command-t/ruby/command-t && ruby extconf.rb && make ` 这个插件十分有用类似sublime中Ctrl+p的功能。

bashrc 及其解释:
```
export TERM='xterm-256color' # 使用256色，
export LANG=en_US.UTF-8      # 解决乱码
export LC_CTYPE="en_US.UTF-8"# 解决乱码
export LC_ALL=en_US.UTF-8    # 解决乱码

stty stop ''   # vim使用Ctrl+s 与shell快捷键冲突解决
stty start ''  # vim使用Ctrl+s 与shell快捷键冲突解决
stty -ixon     # vim使用Ctrl+s 与shell快捷键冲突解决
stty -ixoff    # vim使用Ctrl+s 与shell快捷键冲突解决
```

vimrc 因为文件太长不适宜放到blog中说明，直接注释了[https://github.com/lijianying10/FixLinux/blob/master/golangdev/vimrc](https://github.com/lijianying10/FixLinux/blob/master/golangdev/vimrc)

直接下载请参考这里：https://registry.hub.docker.com/u/lijianying10/golangdev 注意不要用latest 是调试版本。
