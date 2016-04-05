title: 使用aglio自动化管理API文档
date: 2016-04-05 22:25:28
categories: 技术
tags: [npm,nodejs,aglio,docker,shell,git,nginx]
---

## VERSION: 1

## 背景

前后端需要良好的沟通交流环境无法离开实时更新的文档，我没无法离开实时更新的API接口文档。使用Github hooks在服务器上实现文档自动更新是一个不错的选择。

`注意本文中hook代码中删掉了所有安全验证的内容`

## 托管容器

文档的托管容器只管所有的关键其他的所有东西都没有。
既： 环境与业务分离，以任务执行的不同来分割容器。

``` sh
FROM node
RUN apt-get update && apt-get install -y git build-essential curl vim && apt-get clean
RUN npm install aglio -g
CMD bash /e.sh
```

`注意：`容器构建时由于包本身需要构建的原因，使用内存会达到2gb。

## 自动运行脚本 `e.sh`

```
echo "build time" $(date)
cd /root/
rm -rf XXX-doc
git clone --depth=1 git@github.com:XXX/XXX-doc.git
cd XXX-doc
nohup aglio -i android.md -h 0.0.0.0 -p 3001 -s >> /tmp/log 2>&1 &
nohup aglio -i admin.md -h 0.0.0.0 -p 3002 -s >> /tmp/log 2>&1 &
tail -F /tmp/log
```

首先给出自动化build时间,紧接着在容器中下载代码，进入后直接使用命令运行。最后不要忘记输出日志。

## 钩子容器

``` go
package main

import (
    "bytes"
    "io"
    "net/http"
    "os"
    "os/exec"

    "github.com/wothing/log"
)

// hello 文档管理的首页，提供强制更新容器的方法。
func hello(w http.ResponseWriter, r *http.Request) {
    resp := `<h1>CURRENT STATUS <a href="/afejid">force update</a></h1><pre>`
    o, e := CMD("docker logs --tail=100 docs_running_container")
    resp = resp + o
    resp = resp + e
    resp = resp + "</pre>"
    io.WriteString(w, resp)
}

// afejid 重启容器触发,容器重启之后会从新clone代码以及重新生成文档。
func afejid(w http.ResponseWriter, r *http.Request) {
    // !!! 这里不要忘了安全验证。
    io.WriteString(w, "updateing already started")
    go CMD("docker restart docs_running_container")
}

func main() {
    http.HandleFunc("/", hello)
    http.HandleFunc("/afejid", afejid)
    http.ListenAndServe(":8000", nil)
}

// CMD 调用系统命令。
func CMD(order string) (string, string) {
    log.Tinfof("", "RUN: %s", order)
    cmd := exec.Command("/usr/bin/script", "-e", "-q", "-c", order)
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        log.Terrorf("", "%s --> %s \n order return none zero code STDERR: \n %s", order, err.Error(), stderr.String())
        os.Exit(1)
    }
    return out.String(), stderr.String()
}
```


## ngxin 配合

这里主要涉及的问题是一个端口给两个端口的应用使用。
我们这里使用路由分开。在本文中给出了两个例子，一个是安卓的API另外一个是管理端的API

使用技术内容： nginx `rewrite + proxy_pass`

这里我们给出访问安卓的方法 `doc.XXX.XXX/android`

nginx 配置参考
```
    location ^~/android {
        auth_basic "Restricted";
        auth_basic_user_file /etc/nginx/htpasswd;
        rewrite /android(.*) /$1  break;
        proxy_pass http://docs_running_container.aaa:3001/;
    }
```


`注意： `location 内第一行第二行的内容为账号密码验证。


## 总结

在此次实践中想了非常多的东西对依赖的分析执行的计划思考，以及更新hook调度的范围，以及真个项目的复用都进行了成本考量。还是花了不少时间。在这了做了一些简单的小终结。
