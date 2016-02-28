title: 微服务实践 - golang Thrift 开发一个月的感受 
date: 2016-02-28 23:31:14
categories: 技术
tags: [golang,thrift,microservice]
---

![](http://7viiaq.com1.z0.glb.clouddn.com/thrift.jpg)

## 背景
这是我对微服务向往已久的第一次实践虽然屡屡碰壁工期紧张，但最后还是按照正确的时间完成了任务。在这里记录一下使用Thrift的感受。整个工程希望能够有清晰的文档，清晰的通讯逻辑，干净整洁的代码。其中还有很多团队的约定。

注意以下代码都只谈技术不谈业务。

## RPC选型
我们看中了Thrift的性能。

## 实践

### IDL SAMPLE
```
struct UserInfo{
    1:required string ID (go.tag = 'json:"user_id" bson:"user_id"'),
    2:required string Name,
    4:optional i32 XXXXX,
}
```

优点

1. 编号使得Message具有弱化的版本功能使得不同版本之间的通讯成为可能。
2. 数据类型丰富。
3. 可以使用gotag来辅助定义golang中的结构体。

缺点

1. 与Mongodb的不兼容，Thrift没有提供自定义类型，比如说定义一个`ObjectId`类型，由于此造成时间浪费在结构体之间的数据互相转换。
    - 为啥不用UUID？ 因为我们需要时间搜索以及分页搜索，所以只能用ObjectId。
2. required 使用的是值类型，optional使用的是引用类型（一个带指针一个不带指针）所以使用此方法来`复用结构体`开发任务还是比较辛苦的。
3. 带编号还是挺烦的，编辑文件的时候出现错误时有发生。

``` golang
service UserService {
    // Check if the service is health
    string Ping(),

	// CreateUser: create a new user
	string CreateUser(1:string traceID, 2:string userID, 3:string password) throws (1:woerr.WoError we),
}
```

优点

1. 比较通用的注释风格。
2. 异常处理很方便。

缺点

1. 没有提供微服务重连的问题的解决方案（Ping）。
2. 没有提供服务踪(tracerID)。
3. 返回值只能有一个。

``` shell
thrift --gen go:package_prefix="github.com/wothing/thrift/" -out . user.thrift
```

最后生成的结果是一个包，并且里面的代码很长很难看。可以把类似上面的命令写到Makefile备用。

### 服务器端开发

``` golang
package main

import (
	"flag"
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/wothing/log"
	"github.com/wothing/thrift/user"
)

func main() {
	log.SetOutputLevel(log.Linfo)

	port := flag.String("p", "3001", "listening port")
	flag.Parse()
	listenAddr := fmt.Sprintf(":%s", *port)

    // Init dependance modules
    // ....................................................

	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	serverTransport, err := thrift.NewTServerSocket(listenAddr)
	if err != nil {
		log.Fatalf("error on creating server socket : %s", err.Error())
		return
	}

	handler := &UserServiceImpl{}
	processor := user.NewUserServiceProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	log.Infof("User Service servering in %s", listenAddr)
	if err = server.Serve(); err != nil {
		log.Errorf("User Service startup error: %s", err.Error())
	}
}
```

然后接下来就针对结构体`UserServiceImpl`进行Interface的实现就可以了。最快速最准确的方法是去找自动生成的代码复制过来，在类型中不要忘记包的引用就可以。因为他们处于不同的包的位置。

可以看到Thrift Server 初始化是依赖自动生成的包来完成服务器的初始化的，更底层的依赖于Thrift的golang源代码。如此的实现方式耦合性的确有点强。

### 客户端的实现（给Gateway用）

由于们某些原因往后写的时候发现微服务呈现网状模式，客户端调用不只是Gateway所以客户端的开发独立一个包出来。

通用服务连接

``` golang
func prepareConn(svcName string) (*thrift.TBinaryProtocolFactory, *thrift.TTransport, error) {
    // Sevice discoary -> address
    // ...........................
    
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(address.String())
	if err != nil {
		log.Errorf("new socket to '%s' service fail: %s", svcName, err)
		return nil, nil, errors.New(woerr.ErrMicroService)
	}

	useTransport := transportFactory.GetTransport(transport)
	registerConnection(svcName, useTransport)

	if err := useTransport.Open(); err != nil {
		log.Errorf("connect to '%s' service fail: %s", svcName, err)
		return nil, nil, errors.New(woerr.ErrMicroService)
	}

	return protocolFactory, &useTransport, nil
}
```

具体服务的连接

``` golang
func CheckUserConn(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	svcName := userSvcName
	if userService != nil {
		if _, err := userService.Ping(); err == nil {
			next(rw, r)
			return
		}
		Services[svcName].Close()
	}

	protocolFactory, useTransport, err := prepareConn(svcName)
	if err != nil {
		misc.RespondString(rw, fmt.Sprintf(`{"code":"%s", "message":"connect to %s service error"}`, woerr.ErrMicroService, svcName))
		return
	}

	userService = user.NewUserServiceClientFactory(*useTransport, protocolFactory)
	log.Infof("connected to '%s' service", svcName)
	next(rw, r)
}
```

此代码为negroni中间件，缺陷在于，到路由的时候都要对服务健康进行检查。

微服务的调用

``` golang
func XXXX(rw http.ResponseWriter, r *http.Request) {
	// Generate TracerID

	// Check data format

	userID, err := userService.CreateUser(traceID, userID, password)
	if err != nil {
		// handle error
	}

	// HTTP response
}
```


## 其他需要注意的问题

1. 控制Gateway 的代码体积。
2. 参数检查如果依赖于其他微服务的情况应该谨慎处理。
3. 很多check处理都可有通用函数，但是团队成员不一定能反应过来。

## 总结

Thrift 只是提供了通讯方案，其他都需要自己解决。

### 微服务的挑战
1. 性能不是首要问题。
2. 开发速度。
3. 横向拓展。
4. 负载均衡。
5. 在做技术选型时候的蝴蝶效应，一点点的决策失误都有可能带来大量的人工浪费。
