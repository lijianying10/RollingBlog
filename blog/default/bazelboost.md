title: bazel 使用经验和记录
date: 2021-01-15 19:01:18
categories: 技术
tags: [bazel]
---

近期我自己做的项目已经用上 Bazel 这款产品了。 它是我目前学过的最难学的一门技术（工具）了。 下面总结一下我的学习和使用过程和方法。
一方面我想提升自己的总结能力，另外一方面我自己学习的时候很寂寞身边没有人会这门技术，大多都是用一用编译就好了， 并不想深入了解。

## 我为什么要用bazel

我最近在做3d建模视频我的主要业务逻辑使用的go语言，blender只能使用python语言，不少操作的驱动需要使用python实现。
不少SVG从矢量到标量的计算用到了NodeJS，并且我的项目并不能统一语言去做。

依赖复杂，我有很多种action而且未来还会有更多。例如目前我有的操作：

1. SVG矢量转标量
2. 依赖上面的结果让两个svg path 正交 绘制 mesh 并且输出 Wavefront OBJ 3d 模型的格式 （未来还会有更多的 Mesh 建模算法来应对各种不同的建模场景）
3. 运行blender background 模式调用我写的 python library 和 python RPC server 对接 golang 建模用业务层。
4. blender建模和动画的结果之间可以相互引用和复用。

所以这是一个相当复杂的 Dag （有向无环图） 因此每当我动一个比较底层的算法我需要 快速的 正确的 构建并且执行Dag依赖树。

当我编译好我的blender工程之后我就可以把blender文件上传到GPU集群来渲染视频。

一个更简单的为什么要使用bazel的例子：

当你自己写了一个 golang 程序帮助你输入一个 SQL 语句输出一个对应的驱动包来避免使用 ORM 当你修改了 code gen 的模板时你怎么知道该更新那些和测试那些依赖于它的目标 ( Target ) 呢?

并且你并不想写一大堆重复的类似的代码 干体力活 希望机器生成代码， 并且你知道 Golang 没有泛型造成很多重复代码时， 我们唯一的选择就只能是设计一个Code Generator。

## 学习路线

这里记录一下我的学习路线方便当我自己忘记知识点可以看这篇文字帮助自己复习， 我的学习目标并不是 使用 bazel  而是熟练的掌握开发 rules 应对各类情况。

### Concepts 了解基本概念全部读完。

举个例子：

```
Transitive dependencies
Bazel only reads dependencies listed in your WORKSPACE file. If your project (A) depends on another project (B) which lists a dependency on a third project (C) in its WORKSPACE file, youll have to add both B and C to your projects WORKSPACE file. This requirement can balloon the WORKSPACE file size, but limits the chances of having one library include C at version 1.0 and another include C at 2.0.
```

当你的 Workspace 越来越大时如果你能想起这里来。 你会发现并不是忘记读了什么文档使用什么方法而是就这样设计的。

Concepts 里面有非常多设计者的意图，是一个系统设计非常好的学习资料。

### Extending Bazel 中的 Extension Overview 和 Concepts

Extending Bazel 其实就是写rules来处理自己的工作情况。
这里一定都要读完，意义很大比如一个重要的概念 Depsets 你读完文档之后并不知道它是 immutable variable。
如果你并不理解文档的情况下就会一脸懵逼的碰到 rules 的设计问题，很多变量是只读的文档里面没有提到但是你要反应过来。
所以前期文档阅读量不足会耽误很多时间。因为遇到问题你反应不过来。

### 换个角度继续学

当我看完上面的文档时还不知道怎么下手， 云里雾里， 这时候有一个重要的 blog 出现：
https://www.jayconrod.com/posts/106/writing-bazel-rules--simple-binary-rule

这些文章需要拿出来反复琢磨：

1. Writing Bazel rules: simple binary rule 帮助你学习如何构建好一个target
2. Writing Bazel rules: library rule, depsets, providers 帮助你学习如何处理好依赖，很重要，因为你知道bazel的核心依赖算法是DAG
3. Writing Bazel rules: repository rules 帮助你把外部的软件转换成一个 bazel 系统可以接受的依赖很重要，比如说你构建时依赖了 blender 的二进制文件去运行， 它就会帮助你做到类似于 `rules_go` 一样下载 golang 的 toolchain

以上是写的非常好的非官方的tutorial。

我看到这时还是不会写。接下来开始读源码推荐阅读：

1. 一个复杂的： [GITHUB REPO](https://github.com/bazelbuild/rules_go) 去看首页文档写的重要API的实现。
2. 一个简单的： [GITHUB REPO](https://github.com/zaucy/rules_blender) 去看比较简单的case的关键rules实现。

当我看到这里时我基本上已经写了个大概。因为已经大概熟悉了 Skylark 这门语言。进入爬坑阶段。

### 最重要的我常常需要参考的关键数据结构

Actions https://docs.bazel.build/versions/master/skylark/lib/actions.html 这是核心中的核心
Files https://docs.bazel.build/versions/master/skylark/lib/File.html 写 rules 其实就是灵活的处理你的代码文件与编译器之间做沟通。
DefaultInfo https://docs.bazel.build/versions/master/skylark/lib/DefaultInfo.html 理解好这个概念可以帮助你与其他开发者写的 Rules 很好的做对接和依赖。

## 爬坑经验

### 代理方案选择

1. linux iptables ip filter 省事
2. windows Proxifier

windows 注意：

因为 Proxifier 是按照进程为最基本单位的， 当你看不到具体那个进程网络卡住了可以使用工具，可以看到进程树：

process monitor： https://docs.microsoft.com/en-us/sysinternals/downloads/procmon

目前我的进程rules: `bazel.exe;java.exe;fetch_repo.exe;Conhost.exe;go.exe;git.exe;git-remote-https.exe;` 供参考

### 管理你系统外的依赖

系统外依赖的定义是： 你并不关心实现的代码。仅仅是下载， 比如说 GJSON 这类的

Cons:

你在A项目中引用了B项目 B项目引用了GJSON 那么A的workspace 也需要管理 GJSON 依赖。 这时你的 WORKSPACE 文件会像气球一样膨胀
我遇到这个坑是因为我对 Concept 文档理解的不够深入。

Pros:

可以实现 Shadowing dependencies 在文档： https://docs.bazel.build/versions/master/external.html 提到的

### 关于 actions run 这个函数的坑

这些坑都是希望假如我身边有人会这技术我希望他能提前告诉我的

1. `declare_file` 一定要放到output里面。 这样当你的 actions.run 结束时文件没创建， 会显示编译错误。
2. `inputs` 一定要把你依赖的所有文件都放到这个数组里面。 就算 bazel query 时的确显示出了正确的依赖关系， 但是 inputs 没有声明这些文件在actions中依赖了你会你修改了文件不会重新编译。 会破坏编译的正确性。
3. `mnemonic` 只能是一个单词，放入要给动词可以帮助你理解正在并行编译时在做什么动作。
4. `executable` 不要依赖与你的path或者一个绝对路径那就出现了环境依赖，要使用 target 来保证整个系统时封闭的来保证你的 `正确性`
5. `use_default_shell_env` `False` = 有节操 `True` = 没节操 如果你选择省事至少放到Docker容器里面定义好Path运行。不然会出现环境依赖以及脱离 bazel 环境管理的软件依赖。
6. `tools` 假如你在构建中依赖了其他的二进制文件， 例如我在运行 blender 我的入口时 python 的 RPC Server， 我用 Golang 写 RPC client 你很熟练的通过 arguments 或者 env 把你的 Golang client 二进制文件放到了 执行的地方，但是你会新奇的发现 bazel 系统的构建 DAG 树的 Graphviz 图纸是有依赖的但是它并不构建， 然而你读了很多遍文档也不知道为什么你依赖的二进制 bazel target 不编译因为文档的描述是这样的： `tools: List or depset of any tools needed by the action. Tools are inputs with additional runfiles that are automatically made available to the action.` 现在我可以高兴的告诉你这时候要放到tools里面。出现这个问题并不是文档有问题，是 bazel 的设计者经验丰富抽象层次高造成的。但是我给跪了。

#### input 数组的开发技巧

例如说你依赖了 `py_library` 这个 bazel native rule

第一你需要debug看这个rules返回给你什么info

下面时rules label attr

```
"util":attr.label(),
```

输入参数label例子： 下面这个例子的util其实就是一个 `py_library rule`

```
util = "//blenderutil:util",
```

下面时skylark 调试代码, 可以看到有什么Info

```
print(ctx.attr.util)
```

一点点的调试最后可以得到下面的代码：

```
for f in ctx.attr.util[PyInfo].transitive_sources.to_list():
    input_list.append(f)
```

你会得到你所有依赖的内容如果你的依赖很复杂 比如说  `A->B->C`  `A->D->C` 你会发现C重复了。可以直接用depset解决这个问题 `new_var = depset([这里时你的数组，或者直接把变量放进去])`

#### executable 的例子

定义label

```
"_blenderRunner":attr.label(
    default = Label("//blenderRunner:blenderRunner"),
    executable = True,
    cfg = "exec"
),
```

action.run 参数例子

```
executable = ctx.executable._blenderRunner, 
```


## 一些感悟

### 为什么 golang 要用 bazel

目前 gomod 已经做的够好了，就算在我的这个工程里面也会使用 golang gomod 来开发和测试我的代码。
但是与 bazel 结合其实成本并不高有重要的工具叫 `gazelle`
他可以帮助你转换 gomod 项目成 bazel 项目。
也可以帮助你同步 gomod 到 bazel。

通过如下命令：

```
bazel run //:gazelle # 初始化 
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies # 更新
```

需要注意：

1. 你还是要显式从 WORKSPACE 或者你自己的 bazel function 中声明你的 `go_repository` rule 来定义依赖。
2. 关于 golang bazel 项目吐槽最多的 protobuf 生成之后 IDE 很难结合的问题可以结合这个思路 https://github.com/nikunjy/golink 另外我在自己读 bazel rule 开发文档的时候看到其实我们是可以通过创建软连接来帮助IDE找到我们已经生成的代码。
3. 上面的问题还有另外一个解决思路， 修改 Gocode 的实现，目前 Gocode 这个项目已经有 daemon 了。 但是它扫描代码并不积极。

### 一个恰到好处的设计

从actions.run这里可以看到，
每个参数设计的都很到位， 尤其是限制的很到位， 并且每个每个参数都不能删掉， 并且很好的兼顾了各种情况。
新手只可能对概念的理解和知识的掌握不到位， 很难产生一些错误， 这大概就是牛逼的架构师设计的接口。
把安全做到语言级别， 给你报错而不是你在写代码的时候给你各种需要主动遵守的规范。

如果你并不享受这些对你的限制，还不如用 python 或者 make 之类作为构建工具。

### 一些 bazel 相比其他构建系统的好处

1. 封装的好，会让构建过程产出结果稳定，受环境的影响可以做到尽量小，也可以让依赖多个版本并存
2. Starlark 语言设计的好， 类似 python 写起来贼舒服被裁剪的很好， 一方面是影响话编译环境的API都被裁掉了（你并不能很容易的写文件到磁盘到任意位置）内部帮助你构建的变量都是只读的。当然它的好是相对于 GnuMake 和 CMake 相比。
3. 因为拓展性好，基本上你能想到的语言都有了现成的Rules支持。
4. 对 Code Generation 非常友好。
5. `正确` `快速` 实至名归

## 最后

很开心我看到了 Bazel 的门， 希望我能早日入门。

