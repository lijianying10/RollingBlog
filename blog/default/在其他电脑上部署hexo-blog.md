title: hexo blog 迁移
date: 2015-02-13 15:21:09
tags: [hexo,git,blog,node]
---

最近回家了终于能使用我的大电脑了。
但是迁移的时候发现了很多问题。
一一解决之后总结出来如下的标准步骤就可以搞定了

## 步骤Mapping
1. clone回来自己的blog代码
2. theme用了submod需要给自己的主题一个单独的git repo (如果有的话就直接克隆回来就好了)
3. 在跟目录下 npm install --registry=https://registry.npm.taobao.org
