title: "golang调优之clock ticks"
date: 2015-05-29 12:04:33
categories: 技术
tags: [profiling,golang] 
---

## 本blog的来源

昨天在找工作面试的时候我与面试官聊到了golang的问题。当然讨论的热点就是调优与GC。
结果面试变成了技术讨论与研究，聊了接近一个小时，真的很开心。
下面的研究内容来自goblog [https://blog.golang.org/profiling-go-programs](https://blog.golang.org/profiling-go-programs)
我也只是想浓缩一遍上面的内容方便大家研习。当然文章可能比较老了。
因此我在这里重新走一遍大神之路：

## 问题来源：

来自论文：[http://research.google.com/pubs/pub37122.html](http://research.google.com/pubs/pub37122.html)
提出的挑战，在这篇文章中golang的性能是最低的。因此在本blog中就针对这篇文章中的算法进行调优。

先说明我的各种版本号：

```
[#11#ljy@ljydeiMac ~/cpp_pro]$g++ --version
Configured with: --prefix=/Library/Developer/CommandLineTools/usr --with-gxx-include-dir=/usr/include/c++/4.2.1
Apple LLVM version 6.0 (clang-600.0.56) (based on LLVM 3.5svn)
Target: x86_64-apple-darwin14.0.0
Thread model: posix
[#12#ljy@ljydeiMac ~/cpp_pro]$go version
go version go1.4.2 darwin/amd64

C++ 成绩：
real	0m16.791s
user	0m16.093s
sys	0m0.687s

golang 成绩：
real	0m26.582s
user	0m26.393s
sys	0m0.161s
```

## pprof运行原理and解释and调优

```
(pprof) top10
Total: 2525 samples
     298  11.8%  11.8%      345  13.7% runtime.mapaccess1_fast64
     268  10.6%  22.4%     2124  84.1% main.FindLoops
     251   9.9%  32.4%      451  17.9% scanblock
     178   7.0%  39.4%      351  13.9% hash_insert
     131   5.2%  44.6%      158   6.3% sweepspan
     119   4.7%  49.3%      350  13.9% main.DFS
      96   3.8%  53.1%       98   3.9% flushptrbuf
      95   3.8%  56.9%       95   3.8% runtime.aeshash64
      95   3.8%  60.6%      101   4.0% runtime.settype_flush
      88   3.5%  64.1%      988  39.1% runtime.mallocgc
```
pprof模块通过每秒大概100次的对runtime 中的 `stack` 进行取样来进行统计的。下面来解释一下报表为啥是上面这个样子。
首先 Total 2525 程序大概运行了25s+
-------这一部分是针对单个函数的统计
col1： 在取样中作为栈顶的次数
col2： 作为堆顶的百分比，以第一行为例统计关系：298/2525 约等于 11.8% 就好理解了
col3： 排名结果的累加，都是这个位置的数的上面加左面获取的结果，有了这个就可以大概看出来几个热点占用的总比例，非常方便
-------这一部分是对整个堆栈的统计。与上面的区别是不考虑是否在堆栈顶部。
col4： 在sample堆栈中出现的次数，不管是waiting还是return只要出现就计入统计。
col5： 出现次数百分比，与左边报表左边类似。
col6： 略

这种统计方法不但不会影响太多程序性能，而且可以很好的把握程序热点在何位置。
在Intel Vtune中它会帮你完全统计出函数所用的时间。虽然非常爽但是其实没有什么大作用。
有个大概百分比就基本够用了。
不失为一种定性与定量的中间选择。其实我在做log系统的时候也可以仿照他的来做。

```
(pprof) list DFS
Total: 2525 samples
ROUTINE ====================== main.DFS in /home/rsc/g/benchgraffiti/havlak/havlak1.go
   119    697 Total samples (flat / cumulative)
     3      3  240: func DFS(currentNode *BasicBlock, nodes []*UnionFindNode, number map[*BasicBlock]int, last []int, current int) int {
     1      1  241:     nodes[current].Init(currentNode, current)
     1     37  242:     number[currentNode] = current
     .      .  243:
     1      1  244:     lastid := current
    89     89  245:     for _, target := range currentNode.OutEdges {
     9    152  246:             if number[target] == unvisited {
     7    354  247:                     lastid = DFS(target, nodes, number, last, lastid+1)
     .      .  248:             }
     .      .  249:     }
     7     59  250:     last[number[currentNode]] = lastid
     1      1  251:     return lastid
(pprof)
```

虽然不需要解释，但是很容易看出来那句话执行时间是最长的(L:247)
主要热点问题在于使用了map进行搜索。
在这个blog中提出了使用[]int的方式给map增加类似索引的东西。效果不错。（(*^__^*) 嘻嘻……其实我在自己做cache搜索的时候也这么做。）

Tip：作者的compiler是6g很老版本的。这里补充一下go1.4.2的成绩：
```

调优前
$time ./havlak1
# of loops: 76000 (including 1 artificial root node)

real	0m21.686s
user	0m21.578s
sys	0m0.111s

按照上面方法调优后
time ./havlak2
# of loops: 76000 (including 1 artificial root node)

real	0m12.588s
user	0m12.486s
sys	0m0.103s

对比我们上一次测试的源代码修改了之后的成绩，可以看到与C++的成绩非常接近了。(GO:26s,CPP:16s)
time ./go_pro
Welcome to LoopTesterApp, Go edition
Constructing Simple CFG...
15000 dummy loops
Constructing CFG...
Performing Loop Recognition
1 Iteration
Another 50 iterations...
..................................................
# of loops: 76000 (including 1 artificial root node)

real	0m18.447s
user	0m18.297s
sys	0m0.134s
```

调优过程源码：
`hg clone https://code.google.com/p/benchgraffiti`

## 总结

在本篇中，我们可以看到简单的使用pprof模块就可以针对程序的热点进行大幅度的性能改进。
当然我依然坚持认为，在项目prototype开发以及alpha版本中不适合任何角度的调优。
只考虑架构性能已经是最多了（或者说是技术方向性）
但是需要注意的是，pprof本身不会帮你调优，还是要看对golang的熟悉程度。
在这里虽然用了Index来进行调优，但是我们在实战的过程当中可能会更加复杂。
也许路还很远，下一篇GC内存调优。Continue。