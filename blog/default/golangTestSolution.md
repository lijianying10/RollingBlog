title: golang 项目测试方案
date: 2015-8-17 15:33:57
categories: 技术
tags: [golang,python]
---

## 解决问题

形成本次测试方案的前提
1. 工程只有先有质量才能有性能。
2. 在不完全了解方案的时候匆忙的上框架赶进度的后果极有可能是重构。或者你坚持之后发现对于未来的问题很难解决，因为你没读过框架代码。
3. 所以我只给出方案不给出框架的原因就是用我方案的人对我公开的代码可控。
4. 东西越是小才越精美，你跟我扯一些有的没的，我感觉没啥用。
5. 开再发的时候就要不断的测试，有自己的测试方案，测试脚本。
6. 如果直接涉及web业务那么一定要用集成测试的case来进行，条件覆盖，逻辑覆盖。

首先为了更加灵活的`集成测试`的方案实现主要使用了Python代码,并不难外部的包调用也是比较少的。

结果：
1. 这套完成之后，一遍测试一遍做效果果然不错。
2. 基本解决golang 做ajax借口测试的难处。
3. 做完的测试脚本（yaml文件）可以给其他端（Android，IOS，Web）做参考，出现的case一定是通过的case。

## 主体测试文件
参考这里：[https://github.com/lijianying10/FixLinux/blob/master/prob/test.py](https://github.com/lijianying10/FixLinux/blob/master/prob/test.py)
各种注释全都有了。

1. 意义在于，根据后端系统的安全性要求以及通讯设计的不同可以在短时间之内写出一套适应自己想买的测试方案。
2. 常用的模拟登录之类的都好说，直接套用就行。
3. 针对不同的用户角色权限之类的都可以有。

testCase sample：
```
server: 'http://127.0.0.1:9090'
login: 'customer'
case:
    - src: '/permission/XXXn_time'
      data: 'null'

    - src: /XXX/XXX
      data:
        a: dfdfd
        b: dfdf

    - src: /XXX/XXX
      data:
        a: dfdfd
        b: dfdf
```

为啥用yaml:
1. 因为手写好写
2. 因为给其他不懂的人好解释
3. 给出case 参考性强，好合作。

## 业务黑盒白盒测试

Sample:
```golang
package math

import "testing"

func TestAverage(t *testing.T) {
  var v float64
  v = Average([]float64{1,2})
  if v != 1.5 {
    t.Error("Expected 1.5, got ", v)
  }
}
```

上面是一种最简单的黑盒测试情况。
在开发中简单的白盒测试可以通过以上这种手段进行。
业务内的条件覆盖已经逻辑覆盖,多关注内存使用。
尽量减少GC频率（虽然go1.5 的GC好多了）。

## CI持续集成中的在线测试：
1. 在不影响线上业务的情况下转发一份request到测试服务器上。不需要修改端口更不需要暂停服务等操作，可以replay之前的测试请求，可以转发至任意测试服务器。
2. 使用工具Gor
3. 简单Sample： `./gor --input-raw :3000 --output-http :50001`
4. 在上面的例子中，在线业务监听3000端口，gor转发同时监听3000端口，然后转发到50001端口中。 
5. 更加详细的文档：[https://github.com/buger/gor](https://github.com/buger/gor)

