package assets

import (
	"container/list"
	"fmt"
	"time"
)

/*
1. recent article
2. Tags
*/

func ParseTimeStamp() {
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")
	var err error
	for i := 0; i < len(ArticleObject); i++ {
		ArticleObject[i].TimeStamp, err = time.ParseInLocation(timeLayout, ArticleObject[i].Date, loc)
		if err != nil {
			fmt.Println("文章： ", ArticleObject[i].Title, "时间格式有问题！", ArticleObject[i].Date)
		}
	}
}

func RecentArticleGEN() {
	l := list.New()
	l.PushBack(ArticleObject[0])
	for i := 1; i < len(ArticleObject); i++ {
		for e := l.Front(); e != nil; e = e.Next() {
			// front judge
			if ArticleObject[i].TimeStamp.Unix() >= l.Front().Value.(Article).TimeStamp.Unix() {
				l.PushFront(ArticleObject[i])
				break
			}
			// tail judge
			if ArticleObject[i].TimeStamp.Unix() <= l.Back().Value.(Article).TimeStamp.Unix() {
				l.PushBack(ArticleObject[i])
				break
			}

			if ArticleObject[i].TimeStamp.Unix() <= e.Value.(Article).TimeStamp.Unix() && ArticleObject[i].TimeStamp.Unix() >= e.Next().Value.(Article).TimeStamp.Unix() {
				l.InsertAfter(ArticleObject[i], e)
				break
			}
		}
	}
	for e := l.Front(); e != nil; e = e.Next() {
		RecentArticle = append(RecentArticle, e.Value.(Article))
		if len(RecentArticle) == 8 {
			break
		}
	}
}

func TagGEN() {
	Tag = make(map[string]TagOne)
	for _, article := range ArticleObject {
		for _, tag := range article.Tags {
			var tagone TagOne
			if Tag[tag].times != 0 {
				tagone = Tag[tag]
			}
			tagone.times += 1
			tagone.Articles = append(Tag[tag].Articles, article)
			Tag[tag] = tagone
		}
	}

	for k, v := range Tag {
		fmt.Println(k, v.times)
	}
}
