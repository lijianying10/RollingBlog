package main

import "RollingBlog/assets"

func main() {
	// Scan all article in path blog
	assets.Scan()

	// open every file and store to ArticleObject
	assets.OpenFile()
	assets.ParseTimeStamp()
	assets.RecentArticleGEN()
	assets.TagGEN()
}
