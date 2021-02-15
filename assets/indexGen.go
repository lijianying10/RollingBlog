package assets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

func IndexGEN() {
	t, err := template.New("mainbody.tmpl").ParseFiles(`themeBase/mainbody.tmpl`)
	if err != nil {
		fmt.Println("文章模板出错,mainbody.tmpl")
		os.Exit(2)
	}

	RecentArticles := ``

	for _, article := range RecentArticle {

		var wrt bytes.Buffer
		err = t.Execute(&wrt, article)
		if err != nil {
			fmt.Println(err.Error())
		}
		RecentArticles += wrt.String()
	}

	err = ioutil.WriteFile("public/index.html", []byte(PageGEN(RecentArticles, "")), 0644)
	if err != nil {
		fmt.Println("文件存储异常", err.Error())
	}
}
