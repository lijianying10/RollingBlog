title: bazel macro 的开发
date: 2021-02-06 13:17:41
categories: 技术
tags: [bazel]
---

当你有了很多rules可以用在你的工程，当你做一些事务去构建你的目标你会发现BUILDfile开始出现大段的只有参数不同的重复代码降低了代码的复用率。

bazel 的不同于其他构建过程的配置并不是创造一个DSL出来解决这个问题。而通过 Starlark 语言来做到 target 分析。 其关键点在于bazel的运行分为三个阶段 ([three phase](https://docs.bazel.build/versions/master/skylark/concepts.html#evaluation-model)) 其中我们使用 starlark 语言去定义我们的构建目标的过程是在第二个阶段。

## 一个Macro的例子

``` starlark
load(":render.bzl", "blender_render")
load("//ffmpeg:ffmpeg.bzl", "ffmpeg_combine_video")

def blender_render_batch(targets, out):
    render_targets = []
    for t in targets:
        lb = Label(t)
        render_targets.append(t+"_render")
        blender_render(
            name = lb.name+"_render",
            blender_project = t,
            out = lb.name+".mp4"
        )
    ffmpeg_combine_video(
        name = out+"_render",
        input = render_targets,
        out = out+".mp4",
    )
```

CodeReview提示：进入函数之后我们通过 `blender_render` 这个 `rule` 把数组内标记的 `blender project` 渲染成 `mp4` 视频文件并且记录相关的渲染过程的 `target` 到数组 `render_targets` 中最后放到 `ffmpeg_combine_video` 这个 `rule` 把这一系列的mp4文件合并成为一个整体视频。这样我们就可以通过这个方法实现视频的工业化生产减少人工剪辑的介入。

让我们来观察上面的macro实现

1. 我们需要把上面的代码放到 `render_batch.bzl` 中。
   1. 然后我们可以通过 `load("@rules_3dmodule//blender:render_batch.bzl","blender_render_batch")` 如上的load代码就可以准备好调用上面的函数实现了。
2. 我们可以看到其实还是与写 `rule` 很像，只是我们并不需要定义 `rule` 函数中也不需要有 `ctx` 更不需要 `actions` 有所执行。
3. 其中核心的功能实现是通过 `load` rule 之后把 rule 当成函数调用即可。
4. 在写 `starlark` 代码中需要注意只有数组(`array`)和字典(`dict`)是可以修改(`mutable`)的变量其他的都是不可修改(`immutable`)变量这是为了并行, 这在代码开发中至关重要。按照普通开发语言的开发思路去写会因此掉到坑里。

## 深入解析

通过上面的例子我们可以通过分析了解到： `Analysis phase` 从原理上来讲是通过 `(ctx.action)[https://docs.bazel.build/versions/master/skylark/lib/actions.html]` 作为 `(DAG)[https://en.wikipedia.org/wiki/Directed_acyclic_graph]` 的 `Node` 我们可以从[这里](https://docs.bazel.build/versions/master/build-ref.html#dependencies) 来印证我们的观察。 因此我们写的代码或者准确来说 `Macro` 是通过运行之后帮助 `bazel` 来构建整个 `DAG` 的过程。 掌握了这个核心思想之后再开发 `Macro` 会轻松很多。

## 状态处理

当我们学会了把一些列固定的操作写到函数当中很快你会发现你需要对状态进行处理。 例如在我的工程中，我需要对每个 blender target 标记需要并且让构建运行时知道自己在整体工程中的位置，从视频剪辑的角度来讲叫做场序（或者说你是第几个视频片段）因此我们写的Macro需要处理运行状态问题。

一个例子：

以下为文件 `counter.bzl` 的内容

```
def video_scene_append(target_list,target):
    c = len(target_list)
    target_list.append("//%s:%s"%(native.package_name(),target))
    return c
```

文件说明： 我们把target变量append到target_list当中，并且返回target在数组中的index作为函数返回。 这样我们就可以定义出target在构建中的序号了。

以下为 BUILD 的文件内容

```
load("counter.bzl","video_scene_append")

# 1. storage video sequence list
# 2. as a counter for video move sequence
video_scene_list = []

video_scene_append(video_scene_list,"stag1")
```

我们把状态存储的变量放到 `BUILD` 文件当中， 并且调用函数实现功能。

## 配置处理建议

我们可以使用 Jsonnet 来作为配置生成的入口下面是一个例子：

在 `BUILD` file 当中

```
jsonnet_to_json(
    name = "config_gen_stag4", 
    src = "databargroup_config.jsonnet", 
    outs = ["config_gen_stag4.json"], 
    ext_code = {
        "config": """
        {
            default_shift_between_bar:2.6,
            animation_config+:{
                data_bar_keep_frames:24*2,
            },
            num_panel_cfg+:{
                data_division:1.0
            }
        }
        """,
        "data_bar_count":str(first_video_consts_get(stag4_dbg_count_key_name)),
    }, 
)
```

在 `databargroup_config.jsonnet` 当中

```
local tmpl = {
  "default_shift_between_bar": 1.3,
  "title_panel_size_x": 2,
  "title_panel_size_y": 1,
  "title_panel_scale": 0.9,
  "title_panel_distance_to_data_bar": -1.2,
  "title_panel_distance_to_camera": 0.4,
  "num_panel_exist": true,
  "num_panel_cfg": {
    "blah": "blah",
  },
  "animation_config": {
    "blah": "blah",,
  }
};

tmpl + std.extVar('config') + {data_bar_count:std.extVar('data_bar_count')} 
```

这样我们在真正执行构建 `stag4` 这个视频片段就可以使用 `config_gen_stag4` 渲染好的 `json` 来执行 `3d` 建模生产 `blender project`。

其中 Jsonnet 语法可以参考关于 [OOP](https://jsonnet.org/learning/tutorial.html#oo) 的语法解释。

## 配置硬编码问题

当你解决了配置生成问题之后会遇到配置硬编码的问题这里没什么好分享的可以直接[参考文档](https://docs.bazel.build/versions/master/skylark/config.html)。

具体使用中可以参考 [`.bazelrc`](https://docs.bazel.build/versions/master/guide.html#bazelrc-the-bazel-configuration-file) 的说明。

从上面的参考文档我们可以知道 `flag` 的名字是可以使用 bazel label 作为名字的因此我们的 bazelrc 还是一种代码生成很友好的解决方案。

## 最后

我经过了上面的学习和实践同时我也学习到了 [starlark-go](https://github.com/google/starlark-go/) 和 [go-jsonnet](https://github.com/google/go-jsonnet) 结合在一起会是一个非常好的组合， 无论是在配置，编排，还是自动化领域都是一个不错工具。 在未来的产品开发和构建当中我应该会使用这种组合来提升我的人效。

