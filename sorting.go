package main

import "sort"

func sortByDate(articles []Article, ascending bool) []Article {
	sort.Slice(articles, func(i, j int) bool {
		if ascending {
			return articles[i].Created < articles[j].Created
		}
		return articles[i].Created > articles[j].Created
	})

	return articles
}
