title: bootstrap pmodal with requre.js
date: 0000-00-00 00:00:00
categories: 技术
tags: [bootstrap,require.js] 
---

# 很多时候需要数据弹窗。
[bootstrap弹窗插件地址](http://getbootstrap.com/javascript/#modals)

# deploy时候的情况。

## require.js AMD mod
```js
define(['jquery','reqmod/cookie','bootstrap'],function ($,cookie){
    var modal_load = function () {
        $.get('/modal.html', function(result){
            $('body').append(result);
        });
    }

    var modal_show = function(title,contain){
        if(title=='e'){modal_title.innerHTML="<a style='color: red'>错误</a>";}
        if(title=='w'){modal_title.innerHTML="<a style='color: orange'>警告</a>";}
        if(title=='i'){modal_title.innerHTML="提示";}
        modal_contain.innerHTML=contain;
        $('#myModal').modal('show');
    //btn_modal.click();

    }

    return {
        load :modal_load,
        show :modal_show
    };
});
```

## auto load html
```html
<div class="modal fade" id="myModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">&times;</span><span class="sr-only">Close</span></button>
                <h4 id="modal_title" class="modal-title" id="myModalLabel">NULL</h4>
            </div>
            <div id="modal_contain" class="modal-body">NULL
            </div>
            <a style="color: red" id="info"></a>
            <div class="modal-footer">
                <button id="btn_close" type="button" class="btn btn-success" data-dismiss="modal">Close</button>
            </div>
        </div>
    </div>
</div>
```


## dependance
1. jquery
2. bootstrap

## feature
1. load once in a require page
2. invoke every where webiste , just invoke load
3. have fun
