package assets

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// using global variable to generage
func tagHTMLGEN() string {
	res := ``
	for tag, v := range Tag {
		res = res + `<a href="/tags/` + tag + `">` + tag + `(` + strconv.Itoa(v.times) + `)</a></br>` + "\n"
	}
	return res
}

func recnetArticleHTMLGEN() string {
	res := ``
	for _, article := range RecentArticle {
		res = res + `<a href="` + article.URI + `">` + article.Title + `</a></br>`
	}
	return res
}

func MotherPageGEN() {
	t, err := template.New("main.tmpl").ParseFiles(`themeBase/main.tmpl`)
	if err != nil {
		fmt.Println("主模板出错,main.tmpl")
		os.Exit(2)
	}

	taghtml := strings.Split(tagHTMLGEN(), "\n")
	left := ""
	right := ""
	for idx, tagh := range taghtml {
		if idx < len(taghtml)/2 {
			left = left + tagh
			continue
		}
		right = right + tagh
	}

	var wrt bytes.Buffer
	data := struct {
		TagLeft       string
		TagRight      string
		RecentArticle string
		MainBody      string
	}{
		TagLeft:       left,
		TagRight:      right,
		RecentArticle: recnetArticleHTMLGEN(),
		MainBody:      `{{.MainBody}}`,
	}
	err = t.Execute(&wrt, data)
	if err != nil {
		fmt.Println(err.Error())
	}
	MotherPage = wrt.String()
}
