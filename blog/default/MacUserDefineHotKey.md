title: Mac自定义快捷键，以选择焦点屏幕为例
date: 2015-05-07 10:53:21
categories: 技术
tags: [minecraft,java]
---

## 看完本文之后你可以知道：

如何使用Automator自定义服务，并自定义快捷键调用。

## Case： Mac工作站多显示器切换

### 原理：cliclick 命令调用来点击不同的屏幕。 

### 安装cliclick ： brew install cliclick

## 步骤：

### 步骤1：

打开automator

### 步骤2：

![](http://7viiaq.com1.z0.glb.clouddn.com/automator.png)

```
打开之后选择服务，进入到上图中的界面
选择实用工具，运行shell脚本
输入命令 /usr/local/bin/cliclick c:1000,0 
	解释： 这个地方不会载入PATH配置因此使用绝对路径
		  坐标根据自己的情况来定，随便选个不常用的点即可。
		  此命令是模拟鼠标点击，来达到切换屏幕的目的。
运行一下试试效果，如果可以运行最下角会有一个绿色的对号。
保存服务即可。注意：文件名即为服务名
```

![](http://7viiaq.com1.z0.glb.clouddn.com/hotkey.png)

```
成功之后打开设置如上图所示
在服务中可以看到刚才建立的服务，设置快捷键即可运行脚本
点击右键选择automator可以调整命令或者加入其他工作流。
```
## 总结
本人也是刚刚接触automator没想到作用这么强大。
虽然mac提供了很多快捷键，但是对多屏幕的支持偶尔还是要自己动手。
就比如说上面的case，mac工作站多显示器之间的切换这样的问题解决。
