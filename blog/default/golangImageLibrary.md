title: golang 图片处理，剪切，base64数据转换，文件存储
date: 2015-03-05 20:11:59
categories: 技术
tags: golang
---

## AUTH:[PHILO](http://philo.top/about) VERSION:2

本文主要介绍：
1. 图片文件的读写。
2. 图片在go缓存中如何与base64互相转换
3. 图片裁剪

本文中，为了方便查看，去掉所有错误判断

## base64 -> file
```go
ddd, _ := base64.StdEncoding.DecodeString(datasource) //成图片文件并把文件写入到buffer
err2 := ioutil.WriteFile("./output.jpg", ddd, 0666)   //buffer输出到jpg文件中（不做处理，直接写到文件）
```
`datasource` base64 string

## base64 -> buffer
```go
ddd, _ := base64.StdEncoding.DecodeString(datasource) //成图片文件并把文件写入到buffer
bbb := bytes.NewBuffer(ddd)                           // 必须加一个buffer 不然没有read方法就会报错
```
转换成buffer之后里面就有Reader方法了。才能被图片API decode

## buffer-> ImageBuff（图片裁剪,代码接上面）
```go
m, _, _ := image.Decode(bbb)                                       // 图片文件解码
rgbImg := m.(*image.YCbCr)
subImg := rgbImg.SubImage(image.Rect(0, 0, 200, 200)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
```
## img -> file(代码接上面)

```go
f, _ := os.Create("test.jpg")     //创建文件
defer f.Close()                   //关闭文件
jpeg.Encode(f, subImg, nil)       //写入文件
```

## img -> base64(代码接上面)
```go
emptyBuff := bytes.NewBuffer(nil)                  //开辟一个新的空buff
jpeg.Encode(emptyBuff, subImg, nil)                //img写入到buff
dist := make([]byte, 50000)                        //开辟存储空间
base64.StdEncoding.Encode(dist, emptyBuff.Bytes()) //buff转成base64
fmt.Println(string(dist))                          //输出图片base64(type = []byte)
_ = ioutil.WriteFile("./base64pic.txt", dist, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
```

## imgFile -> base64
```go
ff, _ := ioutil.ReadFile("output2.jpg")               //我还是喜欢用这个快速读文件
bufstore := make([]byte, 5000000)                     //数据缓存
base64.StdEncoding.Encode(bufstore, ff)               // 文件转base64
_ = ioutil.WriteFile("./output2.jpg.txt", dist, 0666) //直接写入到文件就ok完活了。
```

大概就是这些代码基本上一些小网站都够用。
缩放什么的可以先靠前端。后端有个裁剪就够了。
