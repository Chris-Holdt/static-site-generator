package main

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

func processArticle(content Content) Content {
	titleIndex := strings.Index(content.Clean, "# ")
	newLineIndex := strings.Index(content.Clean[titleIndex:], "\n")

	title := content.Clean[:newLineIndex+1]

	body := content.Clean[newLineIndex:]

	count := countWords(body)
	content.Metadata.WordCount = count

	extra := fmt.Sprintf("Words: %d\n", count)

	content.Clean = fmt.Sprintf("%s%s%s", title, extra, body)

	return content
}

func countWords(body string) int {

	onlyAlphaNumerics := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	words := strings.FieldsFunc(body, onlyAlphaNumerics)
	count := len(words)

	return count

}

func addArticleDataToConfig(meta Metadata, config *Config) {
	a := Article{
		Created:  meta.Date,
		Filename: meta.FileName,
	}

	config.ArticleList = append(config.ArticleList, a)

	for _, c := range meta.Categories {

		if reflect.DeepEqual(config.CategoryList[c], (Category{})) {
			title := titlecase(c)
			config.CategoryList[c] = Category{
				Title:    title,
				Articles: []string{},
			}
		}

		if entry, ok := config.CategoryList[c]; ok {
			entry.Articles = append(entry.Articles, meta.FileName)
			config.CategoryList[c] = entry
		}

	}
}
