package assets

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/russross/blackfriday"

	"gopkg.in/yaml.v2"
)

func Scan() {
	err := filepath.Walk("./blog", Fliter)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

// handle only markdown file
func Fliter(path string, f os.FileInfo, err error) error {
	if len(path) <= 3 {
		return nil
	}
	if path[len(path)-3:] == ".md" || path[len(path)-3:] == ".MD" {
		ArticlePath = append(ArticlePath, path)
	}
	return nil
}

// OpenFile Open Article file and store to ArticleObject
func OpenFile() {
	for _, path := range ArticlePath {
		mainBody, _ := ioutil.ReadFile(path)

		First := 0
		//find ---
		for idx := 0; idx < len(mainBody)-3; idx++ {
			if mainBody[idx] == byte(45) && mainBody[idx+1] == byte(45) && mainBody[idx+2] == byte(45) {
				First = idx
				break
			}
		}

		// yml part
		var thisArticle Article
		err := yaml.Unmarshal(mainBody[:First-1], &thisArticle)
		if err != nil {
			fmt.Println("文章配置出错:  Path: "+path, err.Error())
		}
		html := blackfriday.MarkdownBasic(mainBody[First+4:])
		thisArticle.Body = string(html)
		ArticleObject = append(ArticleObject, thisArticle)
	}
}
