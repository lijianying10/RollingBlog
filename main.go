package main

import "RollingBlog/assets"

func main() {
	// Scan all artile in path blog
	assets.Scan()

	// open every file and store to ArticleObject
	assets.OpenFile()
}
