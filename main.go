package main

import (
	"RollingBlog/assets"
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
}
