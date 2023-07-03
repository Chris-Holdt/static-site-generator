package main

import (
	"fmt"
	"strings"
)

func injectFeatured(s int, e int, raw string, config *Config) string {
	featured := "Looks like I haven't set this!"
	if config.HomeConfig.ShowFeatured && config.HomeConfig.Featured != "" {
		title := titlecase(config.HomeConfig.Featured)
		featured = fmt.Sprintf("[%s](%s/%s.html)", title, config.ArticleDir, config.HomeConfig.Featured)
	}

	raw = strings.Replace(raw, raw[s:e], featured, -1)

	return raw
}

func injectLatest(s int, e int, raw string, config *Config) string {

	sorted := sortByDate(config.ArticleList, false)

	latest := "Oh dear, I haven't set this up properly!"
	if config.HomeConfig.ShowArticles {
		latest = ""

		for i := 0; i <= config.HomeConfig.ArticleCount && i < len(sorted); i++ {
			a := sorted[i].Filename
			title := titlecase(a)

			latest += fmt.Sprintf("[%s](%s/%s.html) - %s  |  ", title, config.ArticleDir, a, sorted[i].Created)
		}

		latest = strings.TrimSuffix(latest, "  |  ")
	}

	raw = strings.Replace(raw, raw[s:e], latest, -1)

	return raw
}

func injectList(s int, e int, raw string, config *Config) string {
	list := ""

	for _, cat := range config.CategoryList {
		if list != "" {
			list += "\n"
		}

		list += fmt.Sprintf("**%s**\n", cat.Title)

		for _, art := range cat.Articles {
			title := titlecase(art)
			link := fmt.Sprintf("%s/%s.html", config.ArticleDir, art)
			list += fmt.Sprintf("[%s](%s)\n", title, link)
		}
	}

	list += "\n"

	raw = strings.Replace(raw, raw[s:e], list, -1)

	return raw

}
