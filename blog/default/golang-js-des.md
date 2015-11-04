title: golang 与 js 的des加密
date: 2015-03-18 01:13:43
categories: 技术
tags: [golang,js,DES]
---

## AUTH:[PHILO](http://www.philo.top/about/) VERSION:1

![](http://7viiaq.com1.z0.glb.clouddn.com/images.jpeg)

## 1. Abstract

### 1.1 des好处
1. Des加密之后密文比原文长不了多少。
2. 加密参数少。没有iv不用同步比较方便。

### 1.2 des难度
1. 单独语言的加解密example还是比较多的，但是两种语言结果能做到一致就比较难了。
2. 猜出现乱码很难，有可能是字符集。

### 1.3 解决办法
直接复制我整理好的

## 2. js加解密及其测试（来源：github）
引用原文：[https://gist.github.com/ufologist/5581486](https://gist.github.com/ufologist/5581486)
但是做了一些修改。（CryptoJS.pad.ZeroPadding）
```html
//cryptoJS
<script src="./tripledes.js"></script>
<script src="./mode-ecb-min.js"></script>
<script src="./pad-zeropadding-min.js"></script>
```

```js
function encryptByDES(message, key) {
    var keyHex = CryptoJS.enc.Utf8.parse(key);
    var encrypted = CryptoJS.DES.encrypt(message, keyHex, {
        mode: CryptoJS.mode.ECB,
        padding: CryptoJS.pad.ZeroPadding
    });
    return encrypted.toString();
    }
    function decryptByDES(ciphertext, key) {
        var keyHex = CryptoJS.enc.Utf8.parse(key);
        var decrypted = CryptoJS.DES.decrypt({
            ciphertext: CryptoJS.enc.Base64.parse(ciphertext)
        }, keyHex, {
            mode: CryptoJS.mode.ECB,
            padding: CryptoJS.pad.ZeroPadding
        });
        return decrypted.toString(CryptoJS.enc.Utf8);
    }

    var message = 't';
var key = '5e8487e6';

var ciphertext = encryptByDES(message, key);
// ciphertext: 8dKft9vkZ4I=
console.info('ciphertext:', ciphertext);
var plaintext = decryptByDES(ciphertext, key);
// plaintext : Message
console.info('plaintext :', plaintext);
```

## 3. golang 加解密及其测试
引用原文：[https://gist.github.com/cuixin/10612934](https://gist.github.com/cuixin/10612934)
des.go 文件没做修改

与js一起使用的例子：
``` golang
key := []byte("5e8487e6")

//解密，reply是从js中收到的加密字符串
fmt.Println("Received back from client: " + reply)
ddd, _ := base64.StdEncoding.DecodeString(reply)//js加密后的结果是base64的，要转成byte的。
destext, _ := DesDecrypt(ddd, key)
fmt.Println("获取解密结果：", string(destext))//拿到结果

//贴心tip: string->[]byte看这里 []byte("XXX")

//加密，接上面
outs, _ := DesEncrypt(destext, key)
dist := make([]byte, 2048) //开辟存储空间
base64.StdEncoding.Encode(dist, outs)
fmt.Println("加密送出:", string(dist))
```
