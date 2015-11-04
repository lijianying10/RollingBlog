title: "golang调优之GC"
date: 2015-05-29 12:35:27
categories: 技术
tags: [profiling,golang,GC] 
---

## Continue

在[上篇](/2015/05/29/golangProfilingAndGC/)中，主要针对个函数使用时间来对程序分析进行调优。
在本片中由于发现GC使用时间过长来对程序进行分析。
针对内存进行调优。

## 内存调优，查找什么位置引发GC时间过多

第二个版本的pprof结果
```
(pprof) top5
Total: 1652 samples
     197  11.9%  11.9%      382  23.1% scanblock
     189  11.4%  23.4%     1549  93.8% main.FindLoops
     130   7.9%  31.2%      152   9.2% sweepspan
     104   6.3%  37.5%      896  54.2% runtime.mallocgc
      98   5.9%  43.5%      100   6.1% flushptrbuf
```

可以看到GC的时间（runtime.mallocgc）占用到了（54.2%）。

内存profile输出方法
```
var memprofile = flag.String("memprofile", "", "write memory profile to this file")
...

    FindHavlakLoops(cfgraph, lsgraph)
    if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal(err)
        }
        pprof.WriteHeapProfile(f)
        f.Close()
        return
    }
```

调用方法
```
$ make havlak3.mprof
go build havlak3.go
./havlak3 -memprofile=havlak3.mprof
```

结果
```
$ go tool pprof havlak3 havlak3.mprof
Adjusting heap profiles for 1-in-524288 sampling rate
Welcome to pprof!  For help, type 'help'.
(pprof) top5
Total: 82.4 MB
    56.3  68.4%  68.4%     56.3  68.4% main.FindLoops
    17.6  21.3%  89.7%     17.6  21.3% main.(*CFG).CreateNode
     8.0   9.7%  99.4%     25.6  31.0% main.NewBasicBlockEdge
     0.5   0.6% 100.0%      0.5   0.6% itab
     0.0   0.0% 100.0%      0.5   0.6% fmt.init
```

分析过程：
首先说明一下 `1-in-524288` 的意思是内存每变化越0.5MB取样一次。
因为GC的垃圾回收，那么针对内存热点进行分析，看看为啥FindLoops 分配了那么多内存是关键！

```
(pprof) list FindLoops
Total: 82.4 MB
ROUTINE ====================== main.FindLoops in /home/rsc/g/benchgraffiti/havlak/havlak3.go
  56.3   56.3 Total MB (flat / cumulative)
...
   1.9    1.9  268:     nonBackPreds := make([]map[int]bool, size)
   5.8    5.8  269:     backPreds := make([][]int, size)
     .      .  270:
   1.9    1.9  271:     number := make([]int, size)
   1.9    1.9  272:     header := make([]int, size, size)
   1.9    1.9  273:     types := make([]int, size, size)
   1.9    1.9  274:     last := make([]int, size, size)
   1.9    1.9  275:     nodes := make([]*UnionFindNode, size, size)
     .      .  276:
     .      .  277:     for i := 0; i < size; i++ {
   9.5    9.5  278:             nodes[i] = new(UnionFindNode)
     .      .  279:     }
...
     .      .  286:     for i, bb := range cfgraph.Blocks {
     .      .  287:             number[bb.Name] = unvisited
  29.5   29.5  288:             nonBackPreds[i] = make(map[int]bool)
     .      .  289:     }
...
```
在29.5mb的位置可以看到与上一篇内容中出现了一样的错误。能用简单数据结构搞定的问题用了map

因为要看GC所以需要查看程序到底在这个位置上分配了多少次内存
```
$ go tool pprof --inuse_objects havlak3 havlak3.mprof
Adjusting heap profiles for 1-in-524288 sampling rate
Welcome to pprof!  For help, type 'help'.
(pprof) list FindLoops
Total: 1763108 objects
ROUTINE ====================== main.FindLoops in /home/rsc/g/benchgraffiti/havlak/havlak3.go
720903 720903 Total objects (flat / cumulative)
...
     .      .  277:     for i := 0; i < size; i++ {
311296 311296  278:             nodes[i] = new(UnionFindNode)
     .      .  279:     }
     .      .  280:
     .      .  281:     // Step a:
     .      .  282:     //   - initialize all nodes as unvisited.
     .      .  283:     //   - depth-first traversal and numbering.
     .      .  284:     //   - unreached BB's are marked as dead.
     .      .  285:     //
     .      .  286:     for i, bb := range cfgraph.Blocks {
     .      .  287:             number[bb.Name] = unvisited
409600 409600  288:             nonBackPreds[i] = make(map[int]bool)
     .      .  289:     }
...
(pprof)
```

解决问题
[https://code.google.com/p/benchgraffiti/source/diff?name=796913012f93&r=d856c2f698c1&format=unidiff&path=/havlak/havlak.go](https://code.google.com/p/benchgraffiti/source/diff?name=796913012f93&r=d856c2f698c1&format=unidiff&path=/havlak/havlak.go)

尊重原作者，直接贴diff连接。
scan一次代码之后可以很直接的看到map通过算法调整全都不用了。
这也是前文中提到的能用简单的数据结构解决的问题使用了map来解决。

# 但是本次调优更精华的内容在于
1. 直观的看到make的次数减少了。
2. 以变量nonBackPreds为例：
	1. 在原来的执行过程当中通过不断的新建临时map来达到集合的效果（不重复）
	2. 后来调优者使用了append的方式不断加在内存中。
3. 算法中数据存储大部分都集中在了这个变量中（这是避免GC）
4. append slices 中有两个属性一个是len一个是cap使用append从内存角度来看，是操作很长时间才由于cap不够再次申请一次内存，而原来是达到一定周期数
	之后马上alloc内存从频率上就有很多不同，说白了在用slices避免GC
5. 不过我没验证过，在make slices的过程中如果给出cap parameter虽然内存使用可能会变多。但是alloc的次数会变少，会不会有性能优化？
6. 看来算法执行上，越是避免内存操作性能就会越好。

本次调优结果显示性能已经超过C++：
```
$time ./havlak4
# of loops: 76000 (including 1 artificial root node)

real	0m8.965s
user	0m8.875s
sys	0m0.083s
```

接下来发现GC的时间还是很长
```
(pprof) list FindLoops
...
     .      .  270: func FindLoops(cfgraph *CFG, lsgraph *LSG) {
     .      .  271:     if cfgraph.Start == nil {
     .      .  272:             return
     .      .  273:     }
     .      .  274:
     .      .  275:     size := cfgraph.NumNodes()
     .      .  276:
     .    145  277:     nonBackPreds := make([][]int, size)
     .      9  278:     backPreds := make([][]int, size)
     .      .  279:
     .      1  280:     number := make([]int, size)
     .     17  281:     header := make([]int, size, size)
     .      .  282:     types := make([]int, size, size)
     .      .  283:     last := make([]int, size, size)
     .      .  284:     nodes := make([]*UnionFindNode, size, size)
     .      .  285:
     .      .  286:     for i := 0; i < size; i++ {
     2     79  287:             nodes[i] = new(UnionFindNode)
     .      .  288:     }
...
(pprof)
```

在L277位置上发现不断的make导致GC
因此下一步精简内存alloc（只有精简alloc才能减少GC，道理很直接）
做成Global Cache之后测试结果为：

```
$time ./havlak5
# of loops: 76000 (including 1 artificial root node)

real	0m5.742s
user	0m5.664s
sys	0m0.063s
```

之后经过解开多层循环，等等一些列的优化。（profile已经看不出来）
成绩：

```
$time ./havlak6
# of loops: 76000 (including 1 artificial root node)

real	0m2.307s
user	0m2.276s
sys	0m0.031s
```

同等情况下C++的结果(原作中C++结果为1.99s)：
```
$time ./a.out
# of loops: 76000 (including 1 artificial root node)

real	0m4.741s
user	0m4.650s
sys	0m0.081s

O3 优化结果
$time ./a.out
# of loops: 76000 (including 1 artificial root node)

real	0m2.279s
user	0m2.178s
sys	0m0.095s

```

# 总结
首先golang其实性能是可以接近C++的。
当然我们最大的收获是避免GC的具体方法。以及一次珍贵的。
近距离观察大神的算法调优过程。再一次学习了如何减少内存alloc的次数方法。
