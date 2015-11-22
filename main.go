package main

import (
	"RollingBlog/assets"
	"fmt"
	"net/http"
	"os"
)

func main() {
	// Scan all article in path blog
	assets.Scan()

	// open every file and store to ArticleObject
	assets.OpenFile()
	assets.ParseTimeStamp()
	assets.RecentArticleGEN()
	assets.TagGEN()

	// Clean output path
	os.RemoveAll("public/")

	assets.MotherPageGEN()
	assets.ArticleGEN()
	assets.IndexGEN()

	// copy static files
	assets.CopyDir("themeBase/_static", "public/")
	assets.ArchiveGEN()

	fmt.Println("Finish Gen lesson 5k and serve")
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./public"))))
	http.ListenAndServe(":5000", nil)
}
