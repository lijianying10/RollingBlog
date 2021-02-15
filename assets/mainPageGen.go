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
		res = res + `<a href="/` + article.URI + `">` + article.Title + `</a></br>`
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
		PageTitle     string
	}{
		TagLeft:       left,
		TagRight:      right,
		RecentArticle: recnetArticleHTMLGEN(),
		MainBody:      `{{.MainBody}}`,
		PageTitle:     `{{.PageTitle}}`,
	}
	err = t.Execute(&wrt, data)
	if err != nil {
		fmt.Println(err.Error())
	}
	MotherPage = wrt.String()
}

func PageGEN(mainBody, pageTitle string) string {
	if pageTitle == "" {
		pageTitle = "philo.top"
	}
	t, err := template.New("page").Parse(MotherPage)
	if err != nil {
		fmt.Println("生成页面时模板错误", err.Error())
	}

	var w bytes.Buffer

	data := struct {
		MainBody  string
		PageTitle string
	}{
		MainBody:  mainBody,
		PageTitle: pageTitle,
	}
	err = t.Execute(&w, data)
	if err != nil {
		fmt.Println("运行模板出错")
	}
	return w.String()
}
