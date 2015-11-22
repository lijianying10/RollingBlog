package assets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

func ArticleGEN() {
	t, err := template.New("mainbody.tmpl").ParseFiles(`themeBase/mainbody.tmpl`)
	if err != nil {
		fmt.Println("文章模板出错,mainbody.tmpl")
		os.Exit(2)
	}

	for _, article := range ArticleObject {

		var wrt bytes.Buffer
		err = t.Execute(&wrt, ArticleObject[0])
		if err != nil {
			fmt.Println(err.Error())
		}

		os.MkdirAll("public/"+article.URI, 0777)
		err = ioutil.WriteFile("public/"+article.URI+"index.html", []byte(PageGEN(wrt.String())), 0644)
		if err != nil {
			fmt.Println("文件存储异常", err.Error())
		}
	}
}
