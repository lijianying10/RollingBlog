package assets

import "time"

type Article struct {
	Title      string
	Date       string
	TimeStamp  time.Time
	URI        string
	Categories string
	Tags       []string
	Body       string
}

type TagOne struct {
	times    int
	Articles []Article
}
