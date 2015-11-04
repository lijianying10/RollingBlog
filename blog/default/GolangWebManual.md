title: 使用Golang快速构建WEB应用
date: 0000-00-00 00:00:00
categories: 技术
tags: [golang,web,oop,js,requirejs,mongodb]
---

## AUTH:[PHILO](http://philo.top/about) VERSION 2

我们从来都不开发代码，我们只是代码的搬运工。
**-- 阿飞**

希望大家都变卡卡西。 **--啊贱**

大家copy愉快，文档只做参考。如果发现问题或者有好的建议请回复我我回及时更正。



## 1.Abstract
在学习web开发的过程中会遇到很多困难，因此写洗一篇类似综述类的文章。作为路线图从web开发要素的index出发来介绍golang开发的学习流程以及Example代码。
在描述中多是使用代码来描述使用方法不会做过多的说明。最后可以方便的copy代码来实现自己的需求。

本文适应对象：
1. 对web开发有一定经验的人
2. 能够灵活使用ajax的人（至少懂得前后分离）
3. golang web 开发有一定了解，至少略读过一些golang web开发的书籍

看完本文之后您会收获：
1. golang web开发的一些技巧
2. golang web开发的一些实用API

本文在描述的时候为了解释尽量详细，已经把解释写到代码注释中。

## 2.golang web 开发check list

> 略过的部分：基本流程控制，OOP等基础语法知识。

### 2.1本章节提供golang web开发的知识面参考。

- [1.Abstract](#1-Abstract)
- [2.golang web 开发check list](#2-golang_web_开发check_list)
    - [2.1 本章节提供golang_web开发的知识面参考。](#2-1本章节提供golang_web开发的知识面参考。)
- [3.路由器](#3-路由器)
    - [3.1 手动路由](#3-1手动路由)
    - [3.2 手动路由的绑定](#3-2_手动路由的绑定)
        - [3.2.1 静态文件](#3-2-1_静态文件)
        - [3.2.2 固定函数与资源获取](#3-2-2_固定函数)
- [4.页面加载](#4-页面加载)
    - [4.1 纯静态页（HTML）](#4-1_纯静态页（HTML）)
    - [4.2 模板页面的加载](#4-2_模板页面的加载)
- [5.表示层脚本](#5-表示层脚本)
    - [5.1 require.js](#5-1_require-js)
        - [5.1.1 加载](#5-1-1_加载)
        - [5.1.2 页面Business](#5-1-2_页面Business)
    - [5.2 JQuery](#5-2_JQuery)
- [6.业务层](#6-业务层)
- [7.持久层](#7-持久层)
    - [7.1 Mysql](#7-1_Mysql)
    - [7.2 Mongodb](#7-2_Mongodb)
- [8.单元测试注意事项](#8-单元测试注意事项)
- [9.LOG](#9-LOG)
- [总结](#总结)

## 3.路由器
> 路由器是整个网站对外的灵魂，如果路由做的不好URL会非常恶心。
所以这部分设计成第一个要说的内容。

> 路由分两种一种是手动路由为了通过tul调度固定的功能，另外一点就是资源
的获取，通过url的分析来模仿静态页的方式来获取资源（类似get）

> 自动路由，主要使用OOP的COMMAND模式来实现。所有功能使用post，
统一入口，方便权限管理，安全管理，跨域管理。但是如此强大的功能还是
交给框架来做吧。这里就不给新手做参考了。

### 3.1手动路由

```go
package main

import (
  "log"
  "net/http"
  )

  func main() {
    RouterBinding() // 路由绑定函数
    err := http.ListenAndServe(":9090", nil) //设置监听的端口
    if err != nil {
      log.Fatal("ListenAndServe: ", err)
    }
  }

```
在httpserver运行之前先绑定路由

### 3.2 手动路由的绑定

#### 3.2.1 静态文件
```go
http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("./pages"))))
```
#### 3.2.2 固定函数与资源获取

  他们都是一样的
  ```go
  http.HandleFunc("/images/", fileUpload.DownloadPictureAction)
  ```


## 4.页面加载
### 4.1 纯静态页（HTML）
> 直接交给路由就行了。自动就访问那个文件夹了。不过生产环境果然还得是cdn，如果自己服务器比较多。可以nginx反向代理。
主要好处前后分离，能上CDN就是通讯次数多了。不过通过优化改善之类的都还ok啦。

### 4.2 模板页面的加载

```go
commonPage, err := template.ParseFiles("pages/common/head.gtpl", //加载模板
"pages/common/navbar.gtpl", "pages/common/tail.gtpl")
if err != nil {
  panic(err.Error())
}
navArgs := map[string]string{"Home": "home", "User": "yupengfei"}//复杂的参数开始往里塞

knowledgePage, err := template.ParseFiles("pages/knowledge/knowledge.gtpl")
knowledgeArgs := map[string]interface{}{"Head": "This is a test title",
"Author": "kun.wang", "PublishDatetime": "2014-09-14",
"Content": template.HTML("<p style=\"text-indent: 2em\">为什么要用语义呢？</p>")}//不是不好，只是做字符串分析会影响工程效率
commonPage.ExecuteTemplate(w, "header", nil)// render 开始
commonPage.ExecuteTemplate(w, "navbar", navArgs)
knowledgePage.ExecuteTemplate(w, "knowledge", knowledgeArgs)
commonPage.ExecuteTemplate(w, "tail", nil)
```

仅提供关键代码。
 > 1. 其他的都还挺好，就是页面渲染用服务器是不是有点太奢侈了。
 > 2. 字符串数组作为输入参数差错比较困难
 > 3. 总结：虽然减少的通讯次数，但是没办法上CDN蛋疼，另外，模板的mapping蛋疼。

## 5.表示层脚本

表示层脚本做的比较困难也不是很好学。
但是一旦搞定了，代码的复用性会有非常可观的提升。
>就普通情况而言JS开发效率是非常高的灵活度高，并且使用的是客户端的cpu
性能好，免费资源多，学习的人也多，好招聘。


### 5.1 require.js

#### 5.1.1 加载

  ```js
  <script data-main="/reqmod/login_main" language="JavaScript" defer async="true" src="js/r.js"></script>
  ```
  整个网页之留这么一个加载脚本的入口（每个页面最好只有一个js文件）

好处
```
  js是延迟加载。不会出现网页卡死的情况
  最大化使用缓存。（HTTP 304）
  一个网页只用一个js
  dom事件绑定，不用在html控件上写js绑定了
```

坏处
```
  学习比较难
  网站更新始终有缓存没更新的浏览器。造成错误（所以有些情况客户自己就知道多刷新几次了，已经成用户习惯了）
```

参数解释
```
  `data-main` 业务逻辑入口，载入当前字符串.js这个文件
  `language` 不解释
  `defer async` 字面意思
  `src` r.js就是require.js的意思。代码到处都能搞到。
```

#### 5.1.2 页面Business

加载依赖文件
  ```js
  require.baseUrl = "/"
  require.config({
    baseUrl: require.baseUrl,
    paths: {
      "jquery": "js/jquery-1.10.2.min",
      "domready" : "reqmod/domReady",
      "pm" : "reqmod/pmodal",
      "cookie":"reqmod/cookie",
      "user":"reqmod/user",
      "bootstrap": "reqmod/bootstrap.min",
      "nav":"reqmod/nav"
    },
    shim: {
      'bootstrap': {
        deps: ['jquery']
      }
    }
  });
  //直接copy全搞定。
  ```

执行页面business

> 执行里面做的最多的就是dom跟事件绑定而已。加载各种js库直接引用。
代码美观，开发效率，执行效率都是非常棒的。

    ```js
    require(['nav','domready', 'jquery', 'user','pm'], function (nav,doc, $, user,pm){
      //这个函数的第一个`数组`参数是选择的依赖的模块。1. 网站绝对路径。 2. 使用加载依赖模块的时候选择export的内容
      //数组的顺序要跟function顺序一致，如果有两个模块依赖比如说jquery插件，就写道最后不用变量，直接使用`$`
      doc(function () { // domReady
        pm.load();//加载各种插件HTML模板之类的都ok
        $('#btn_login')[0].onclick = function(){user.login();}//button 事件绑定
      });
    });
    ```
页面MODEL
    ```js
    define(['jquery','reqmod/cookie','user','bootstrap'],function ($,cookie,user){
        //define 函数的参数内容require是一样的。
        // 这里依赖的模块要在调用此模块中的模块中有path配置。不然会死的很惨，报错的时候不会说缺少什么什么地方错了。
      var nav_load = function () { // 注意函数定义的方式copy就行
        $.get('/nav.html', function(result){
          var newNode = document.createElement("div");
          newNode.innerHTML = result;
          $('body')[0].insertBefore(newNode,$('body')[0].firstChild);
          //document.body.innerHTML = result + document.body.innerHTML;
          $('#btn_login')[0].onclick = function(){user.login();}
          $('#btn_reg')[0].onclick = function(){window.location='/register.html'}
          $.post('/login_check',{},function(data){
            if(data==0){
              Form_login.style.display=""
            }
            else{
              form_userInfo.style.display=""
            }
          })
        });

      }

      return {//这里类似微型路由。非常灵活，非常方便
        load :nav_load
      };
    });
    ```

### 5.2 JQuery
 >  JQ的功能只要require.js引用了之后基本上都是一样的。
如果有需要可以到w3school上学习一下。

## 6.业务层

Post分析
    ```go
    func XXXAction(w http.ResponseWriter, r *http.Request) {
      r.parseForm() //有这个才能获取参数
      r.Form["Email"] // 获取Email 参数（String）
      // 写接下来的业务。
    }
    ```

资源入口函数资源require分析（url分析固定写法）
    ```go
    func Foo(w http.ResponseWriter, r *http.Request) {
      queryFile := strings.Split(r.URL.Path, "/")
      queryResource := queryFile[len(queryFile)-1] // 解析文件
    }
    //完成字符串分割之后，按照需求来获取资源就可以了。
    ```

直接输入object

```go
data, err := ioutil.ReadAll(r.Body) //直接读取form为 json 字符串
	if err != nil {
		utility.SimpleFeedBack(w, 10, "failed to read body")
		pillarsLog.PillarsLogger.Print("failed to read body")
		return
	}
	k := 【BUSINESS OBJECT】
	err = json.Unmarshal(data, &k)
	if err != nil {
		utility.SimpleFeedBack(w, 13, "Pramaters failed!")
		pillarsLog.PillarsLogger.Print("Pramaters failed!")
		return
	}
//方便快捷。再访问参数的时候，直接调用结构体参数就可以了。
//注意ajax调用函数的时候需要做出一些调整代码如下：
```

```js
$.ajax([dist],JSON.stringify([data]),function(){},'json');//注意JSON
```



## 7.持久层

### 7.1 Mysql
  > 其实不管什么语言的Mysql驱动都是从PRO\*C来的，所以会PRO\*\C之后，啥都好说

Insert Delete Update
    ```go
    stmt, err := mysqlUtility.DBConn.Prepare("INSERT INTO credit (credit_code, user_code, credit_rank) VALUES (?, ?, ?)")
    if err != nil {
      pillarsLog.PillarsLogger.Print(err.Error())
      return false, err
    }
    defer stmt.Close()
    _, err = stmt.Exec(credit.CreditCode, credit.UserCode, credit.CreditRank)
    if err 	!= nil {
      return false, err
      } else {
        return true, err
      }
      //还是比较方便的
    ```

Query
    ```go
    stmt, err := mysqlUtility.DBConn.Prepare(`SELECT commodity_code, commodity_name, description, picture,
      price, storage, count, status,
      insert_datetime, update_datetime FROM commodity WHERE commodity_code = ?`)
      if err != nil {
        return nil, err
      }
      defer stmt.Close()
      result, err := stmt.Query(commodityCode)
      if err != nil {
        return nil, err
      }
      defer result.Close()
      var commodity utility.Commodity
      if result.Next() {
        err = result.Scan(&(commodity.CommodityCode), &(commodity.CommodityName), &(commodity.Description),
        &(commodity.Picture), &(commodity.Price), &(commodity.Storage), &(commodity.Count), &(commodity.Status),
        &(commodity.InsertDatetime), &(commodity.UpdateDatetime))
        if err != nil {
          pillarsLog.PillarsLogger.Print(err.Error())
          return nil, err
        }
      }
      return &commodity, err
    ```

### 7.2 Mongodb
  ```go
  err := 	mongoUtility.PictureCollection.Find(bson.M{"picturecode":*pictureCode}).One(&picture)
  ```
  这里只给出最简单的例子。具体的看mgo的开发文档就ok。还是比较简单的。

## 8.单元测试注意事项
  1. 测试命令 go test -v （没有其他参数了！！！） `如果不带-v只显示结果，不显示调试过程，主要是调试开发的时候用`
  1. 文件格式 xxx_test.go 但是建议改成 xxx_test0.go 或者喜欢改成别的也可以。
    1. 由于测试先行的原则，在开发的时候一次测试也就一两个函数。
    1. 这样相当于把其他测试注释掉
  1. 测试的时候的配置文件要放到测试目录下面。别忘了。
  1. 心态，错误太多一个一个来，要有个好心态。

## 9.LOG
  1. 注意在调试中Log的不可缺失性。
  下面api如果不知道从何而来直接doc搜索就可以了。

  ```go
  package utility

  import "log"
  import "os"
  import "fmt"

  // Logger Model min variable.
  var Logger *log.Logger

  var outFile *os.File

  // init function if Logger if not inited will invoke this function
  func init() {
    if Logger == nil {
      propertyMap := ReadProperty("pic.properties")
      logFileName := propertyMap["LogFile"]
      fmt.Println("Initial and Open log file ", logFileName)
      var err error
      outFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)

      if err != nil {
        panic(err.Error())
      }

      Logger = log.New(outFile, "", log.Ldate|log.Ltime|log.Llongfile)
    }
  }

  // CloseLogFile function : close Logger invoke file.
  func CloseLogFile() {
    outFile.Close()
  }
  ```

  使用方法：
  ```go
  utility.Logger.Println("Log test")
  ```
## 总结
1. 看完这里copy代码日常工作还是能好应付一点。
2. 如果是新手看完这个之后，看那么厚的书就有一定的目标性了。能方便一点
