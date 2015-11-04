title: "golang 日志与配置文件的设计与实践"
date: 2015-07-06 11:21:10
categories: 技术
tags: [golang,yaml] 
---

## version : 2

# 意义
1. 在项目开始之前，先根据大概业务设计日志系统。
2. 一个好的日志系统在项目维护时非常重要。

我的日志系统中的内容：
1. 时间戳
2. 日志调用地点：文件+行号
3. 报错级别 致命 错误 警告 信息 调试
4. Summary 
5. 关键内存Dump

这些信息方便排查故障。

配置文件：
直接使用yaml

包： gopkg.in/yaml.v2

配置文件越简单越好。不然容易出错

# 代码
``` golang
package Logger

import (
	"encoding/json"
	"os"
	"runtime"
	"strconv"
	"time"
)
import "caidanhezi-go/utility"

var outFile *os.File

// Logone 单条Log的结构
type Logone struct {
	Timestamp int64
	Codeline  string
	Level     int // 出错级别： 0致命 1错误 2警告 3信息 4调试
	Info      string
	Detail    map[string]interface{}
}

// init 只负责打开文件
func init() {
	outFile, _ = os.OpenFile(utility.Conf.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
}

// Write 写入Log
func Write(Level int, Info string, Detail map[string]interface{}) {
	var newLog Logone
	newLog.Timestamp = time.Now().Unix()
	_, file, line, _ := runtime.Caller(1)
	newLog.Codeline = file + " " + string(strconv.Itoa(line))
	newLog.Level = Level
	newLog.Detail = Detail
	newLog.Info = Info
	outFile.Write(func() []byte {
		b, _ := json.Marshal(newLog)
		b = append(b, '\n')
		return b
	}())
}

// CloseLogFile function: close logger file
func CloseLogFile() {
	outFile.Close()
}

```

## runtime.Caller 可以返回调用日志的堆栈位置。其他没什么特别的了。

# 为啥不直接用系统中的log?
因为默认log缺少缓存。
缓存系统可以让系统具有日志激增的缓冲措施。

# 日志监控方案

包：	"github.com/ActiveState/tail"

```golang
var err error
t, _ := tail.TailFile(utility.Conf.LogFile, //日志文件位置
	tail.Config{Follow: true, Location: &tail.SeekInfo{0, 2}})
for line := range t.Lines { //每次写入日志都会进入这个循环
	for index := 0; index < len(connected); index++ {
		if err = websocket.Message.Send(&connected[index], line.Text); err != nil { //对每个连接监控单元都发送日志数据
			fmt.Println("Can't send")//掉线的监控T掉。
			// index2 := index + 1
			connected = append([]websocket.Conn{}, connected[:index]...)
			if len(connected) >= index+1 {
				connected = append(connected, connected[index+1:]...)
			}
			continue
		}
	}
}
```

## 注意：SeekInfo 这里的API说明在这里[https://golang.org/pkg/os/#File.Seek](https://golang.org/pkg/os/#File.Seek)
我这里是从最尾部开始。（之前文件的内容都忽略了）

## 8月16日更新
经过一个多月的使用。最后log的write参数改成interface，就不用构建map了。
然后半个月前在读代码的时候发现有个包叫debug。 调试的时候非常划算。

write具体代码：
```
// Write 写入Log
func Write(Level int, Info string, Detail interface{}) {
    var newLog Logone
    newLog.Timestamp = time.Now().Unix()
    _, file, line, _ := runtime.Caller(1)
    newLog.Codeline = file + " " + string(strconv.Itoa(line))
    newLog.Level = Level
    newLog.Detail = Detail
    newLog.Info = Info
    outFile.Write(func() []byte {
        b, _ := json.Marshal(newLog)
        b = append(b, '\n')
        return b
    }())
}
```

## 堆栈堆栈堆栈！
对于不喜欢gdb之类断点，单步调试的我来说对于堆栈的输出是至关重要的。
虽然在golang报错的时候会输出堆栈，但是对于在测试中输出目标执行位置的堆栈还是非常重要的。

`UtilLog.Write(1, "gob encode error"+err.Error(), string(debug.Stack()))`

如果想从控制台上直接输出堆栈也有个好办法。
`debug.PrintStack()`
输出这些信息对于异常处理来说是非常划算的。测试调试的过程都轻松了不少。

