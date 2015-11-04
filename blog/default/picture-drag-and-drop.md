title: picture drag and drop
date: 2015-01-29 21:00:11
categories: 技术
tags: [html5,js,css]
---

#html5 图片拖动到div 生成base64
```css
<style type="text/css" media="screen">
        #drop_zone {border:2px #bbb dashed; padding:25px; text-align:center;
            -moz-border-radius:5px; -webkit-border-radius:5px; border-radius:5px;}
</style>
```

```html
<div id="drop_zone"><noscript>必须有js支持</noscript></div>

```

```JS
var box=document.getElementById('drop_zone'); //拖拽区域
       box.innerHTML = "拖动缩略图到这里";
       box.addEventListener("drop",function(e){
           e.preventDefault(); //取消默认浏览器拖拽效果
           var fileList = e.dataTransfer.files; //获取文件对象
           //检测是否是拖拽文件到页面的操作
           if(fileList.length == 0){
               return false;
           }
           //检测文件是不是图片
           if(fileList[0].type.indexOf('image') === -1){
               alert("您拖的不是图片！");
               return false;
           }

           //拖拉图片到浏览器，可以实现预览功能
           var img = window.webkitURL.createObjectURL(fileList[0]);
           var filename = fileList[0].name; //图片名称
           var filesize = Math.floor((fileList[0].size)/1024);
           if(filesize>500){
               pm.show('e','图片不能超过500kb');
               return false;
           }
           console.log(fileList[0]);
           var reader = new FileReader();
           reader.onload = (function() {
               return function(e) {
                   var dataUri = e.target.result,
                       base64 = dataUri.substr(dataUri.indexOf(',') + 1);
                   console.log(base64)
                        };
           })();
           // read file as data URI
           reader.readAsDataURL(fileList[0]);
       },false);

```
