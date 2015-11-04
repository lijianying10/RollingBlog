title: datatables 重新刷新数据的方法
date: 0000-00-00 00:00:00
categories: 技术
tags: [datatables,js] 
---
# 报表报表！
报表系统真的是够让我头疼的。
自从有了dt这个控件之后我再也不怀念微软给我带来的datatables.datasource了。
233

# datatables (datatables.net)
dt 是非常常见的数据控件但是好用的开源的可以商业化的真心不多
但是这个就是一个非常灵活的demo。

参考：
[javascript as data source](http://www.datatables.net/examples/data_sources/js_array.html)

## 数据刷新：

### initial情况参考：
```javascript
$('#example').dataTable( {
    "destroy": true,//如果需要重新加载的时候请加上这个
    "data": dataSet,
    "columns": [
        { "title": "Engine" },
        { "title": "Browser" },
        { "title": "Platform" },
        { "title": "Version", "class": "center" },
        { "title": "Grade", "class": "center" }
    ]
} );
```
加上上面注释的数据主要是因为DT不可以reinitial

# 2014年12月1日补充

还总是有那么几列不听话什么的。
宽度上很不好看
```javascript
$('#example').dataTable( {
"autoWidth": true,
"data": dataset,
"columns": [
{ "title": "index","width": "5%"  },
{ "title": "guid" },
{ "title": "file" },
{ "title": "Version"},
{ "title": "author" }
]
} );
```
