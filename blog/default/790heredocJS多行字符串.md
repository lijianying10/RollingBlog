title: here doc JS 函数转字符串
date: 0000-00-00 00:00:00
categories: 技术
tags: js
---

直接上代码就可以了

```js
function heredoc(fn) {
    return fn.toString().split('\n').slice(1,-1).join('\n') + '\n'
    }

var tmpl = heredoc(function(){/*
!!! 5
html
include header
body
//if IE 6
.alert.alert-error
center 对不起，我们不支持IE6，请升级你的浏览器
a(href="http://windows.microsoft.com/zh-CN/internet-explorer/download-ie") | IE8官方下载
a(href="https://www.google.com/intl/en/chrome/browser/") | Chrome下载
include head
.container
.row-fluid
.span8
block main
include pagerbar
.span4
include sidebar
include footer
include script
*/});
```
