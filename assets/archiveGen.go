package assets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

func ArchiveBodyGEN(articles []Article, archivetitle string) {
	t, err := template.New("archive.tmpl").ParseFiles(`themeBase/archive.tmpl`)
	if err != nil {
		fmt.Println("文章模板出错,archive.tmpl", err.Error())
		os.Exit(2)
	}

	wrts := ""
	for _, article := range articles {

		var wrt bytes.Buffer
		err = t.Execute(&wrt, article)
		if err != nil {
			fmt.Println(err.Error())
		}

		wrts += wrt.String()

	}
	os.MkdirAll("public/"+archivetitle+"/", 0777)
	err = ioutil.WriteFile("public/"+archivetitle+"/"+"index.html", []byte(PageGEN(wrts, "philo.top: Archive")), 0644)
	if err != nil {
		fmt.Println("文件存储异常", err.Error())
	}
}

func ArchiveGEN() {
	ArchiveBodyGEN(ArticleObject, "archive")
	for key, tag := range Tag {
		ArchiveBodyGEN(tag.Articles, "tags/"+key)
	}
}
