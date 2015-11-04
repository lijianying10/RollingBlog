title: "godoc技巧与注意事项"
date: 2015-07-10 10:28:56
categories: 技术
tags: [golang,godoc]
---

#意义
文档对于代码的意义不用多说。
在golang bolg中已经给出了详细的描述[http://blog.golang.org/godoc-documenting-go-code](http://blog.golang.org/godoc-documenting-go-code)
我在实战中踩到了不少坑，这里给出更详细的解释以及注意事项。

我们针对golang源码中的注释进行分析得到如下结果

## 针对Package的文档

### Synopsis

参考[http://golang.org/pkg/](http://golang.org/pkg/)中的Synopsis.
这句话主要出现在针对Package注释中的开头位置。

### OverView

参考[http://golang.org/pkg/archive/tar/](http://golang.org/pkg/archive/tar/)
是针对Package中的注释出现的。如果出现连接，无需标注，生成文档的时候自动会处理成连接

### 参考例子与注意事项
包： [$GOROOT/src/encoding/json]
文件：encode.go
```
// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package json implements encoding and decoding of JSON objects as defined in
// RFC 4627. The mapping between JSON objects and Go values is described
// in the documentation for the Marshal and Unmarshal functions.
//
// See "JSON and Go" for an introduction to this package:
// http://golang.org/doc/articles/json_and_go.html
package json
```
从注释中可以看出第四行是断开的，从第四行开始到package json都为针对包的注释。
目录中Synopsis出现内容为：Package json implements encoding and decoding of JSON objects as defined in RFC 4627.
参考注意事项：
1. 在代码的package上面
2. 在上面不能有空行
3. 注释不能断开(中间不能有空行)
4. 最前面一句话会模块的summary会出现在package index中
5. 第一句话以及之后的内容会出现在OverView中

对比文件：decode.go
```
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Represents JSON data structure using native Go types: booleans, floats,
// strings, arrays, and maps.

package json
```
在package上面有空行，因此只是针对文件的注释不显示在godoc中

## 针对Function
例子：
```
// Marshaler is the interface implemented by objects that
// can marshal themselves into valid JSON.
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}
```
我们可以看到：
1. 在函数上面进行注释
2. 中间不能有空行
3. 开始需要 [空格]FunctionName[空格] Summary
4. 然后继续说明
5. 想圈起来说明参数： 加缩进
进阶技巧：
例子同理于：Function Package
```
// Marshaler is the interface implemented by objects that
/*
can marshal themselves into valid JSON.
*/ 
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}
```
这样不算断开，写文档的时候就方便多了。

## 针对BUG

```
// BUG(src): Mapping between XML elements and data structures is inherently flawed:
// an XML element is an order-dependent collection of anonymous
// values, while a data structure is an order-independent collection
// of named values.
// See package json for a textual representation more suitable
// to data structures.
```
godoc会先查找:[空格]BUG
然后显示在Package说明文档最下面，例子：[http://golang.org/pkg/encoding/xml/](http://golang.org/pkg/encoding/xml/)

## 针对Example
1. 文件名惯用：example_test.go（其他也可以）
2. 包名： apckage_test
3. 方法名：
	1. OverView中： Example
	2. 方法中：      Example[FuncName]
	3. 方法中+一些模式：Example[FuncName]_[Mod]

例子查看：
[http://golang.org/pkg/errors/](http://golang.org/pkg/errors/)

Example文件(example_test.go)：
```
// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"fmt"
	"time"
)

// MyError is an error implementation that includes a time and message.
type MyError struct {
	When time.Time
	What string
}

func (e MyError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

func oops() error {
	return MyError{
		time.Date(1989, 3, 15, 22, 30, 0, 0, time.UTC),
		"the file system has gone away",
	}
}

func Example() {
	if err := oops(); err != nil {
		fmt.Println(err)
	}
	// Output: 1989-03-15 22:30:00 +0000 UTC: the file system has gone away
}
```
1. 注意文件名为：example_test.go
2. 注意package名为 errors_test
3. 针对Function的注释会出现在网页的Example中
4. 如果函数名直接叫Example会直接显示在OverView中

参考文件(errors_test.go)：
```
// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors_test

import (
	"errors"
	"fmt"
	"testing"
)

func TestNewEqual(t *testing.T) {
	// Different allocations should not be equal.
	if errors.New("abc") == errors.New("abc") {
		t.Errorf(`New("abc") == New("abc")`)
	}
	if errors.New("abc") == errors.New("xyz") {
		t.Errorf(`New("abc") == New("xyz")`)
	}

	// Same allocation should be equal to itself (not crash).
	err := errors.New("jkl")
	if err != err {
		t.Errorf(`err != err`)
	}
}

func TestErrorMethod(t *testing.T) {
	err := errors.New("abc")
	if err.Error() != "abc" {
		t.Errorf(`New("abc").Error() = %q, want %q`, err.Error(), "abc")
	}
}

func ExampleNew() {
	err := errors.New("emit macho dwarf: elf header corrupted")
	if err != nil {
		fmt.Print(err)
	}
	// Output: emit macho dwarf: elf header corrupted
}

// The fmt package's Errorf function lets us use the package's formatting
// features to create descriptive error messages.
func ExampleNew_errorf() {
	const name, id = "bimmler", 17
	err := fmt.Errorf("user %q (id %d) not found", name, id)
	if err != nil {
		fmt.Print(err)
	}
	// Output: user "bimmler" (id 17) not found
}

```
1. ExampleNew就是针对New的例子
2. ExampleNew_errorf 给例子加名字详细效果可以查看[这里](http://golang.org/pkg/errors/#example_New_errorf)

## 针对godoc命令
我常用两种方式：
1. `godoc -http=:6060` 直接运行网页上的版本，很方便
2. `godoc package [name ...]` 在开发的时候文档速查

## 总结
一般工程中搞定这些基本就够了。
详细的还是要动手做一做。
我没搞定的：怎么能显示成Main函数的，并且能跑Goplayground