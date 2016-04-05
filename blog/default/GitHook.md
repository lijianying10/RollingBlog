title: GitHook自动检查版本冲突问题
date: 2016-04-05 14:14:07
categories: 技术
tags: [git]
---

## VERSION: 1

![](https://git-scm.com/images/logo@2x.png)

## 背景

为了解决团队代码同步的痛点，每次commit之前都需要检查一次当前主仓库的情况。及时了解其他人的PR动向。

## Fork主仓库的情况

从主仓库Fork到自己的github账号下的工作模式。

``` shell
echo BEFORE COMMIT
ex=$(git remote -v  | awk '{printf "%s\n",$1}' | grep wothing | wc -l)
if [[ $ex == 0 ]];then
    git remote add wothing https://github.com/XXXXX/XXXX.git
fi
git fetch wothing
echo -e "UR \033[32m ahead $(git rev-list  --left-only develop...wothing/develop | wc -l)  \033[0m commits before" 
echo -e "UR \033[31m behind $(git rev-list  --right-only develop...wothing/develop | wc -l) \033[0m commits after"
```

添加此文件到`.git/hooks/pre-commit` 即可
提示在您Commit之前落后主仓库多少版本。

上面的脚本注意修改主仓库的地址。

## 缺点

1. fetch远程仓库比较费时间尤其是github上

## Feature Branch 方法


``` shell
git fetch 
echo -e "UR \033[32m ahead $(git rev-list  --left-only master...develop | wc -l)  \033[0m commits before" 
echo -e "UR \033[31m behind $(git rev-list  --right-only master...develop | wc -l) \033[0m commits after"
```

文件同上一种方法一致。

## 总结

Github的速度真的是非常慢，所以还是建议在GOGS上面做比较好。 
