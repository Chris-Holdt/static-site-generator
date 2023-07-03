package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func convertToHtml(content []Content, config *Config) []Content {

	if isBlankStr(config.HTMLConfig.Top) {
		handleError(errors.New("HTML topper file not set in config"), "")
	}

	if isBlankStr(config.HTMLConfig.Bottom) {
		handleError(errors.New("HTML bottomer file not set in config"), "")
	}

	htmlTop := readFile(config.HTMLConfig.Top)
	htmlBottom := readFile(config.HTMLConfig.Bottom)

	for i, c := range content {
		pageTitle := titlecase(c.Metadata.Title)
		pageTitle = fmt.Sprintf("<title>%s | Chris' Thoughts and Stuff</title>", pageTitle)
		titledTop := strings.Replace(string(htmlTop), "<title></title>", pageTitle, -1)

		body := mdToHtml(c.Clean, c.Metadata.Title)
		full := fmt.Sprintf("%s%s%s", titledTop, body, htmlBottom)
		content[i].HTML = full
	}

	return content
}

func mdToHtml(md string, title string) string {
	var buf bytes.Buffer

	mdParser := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	if err := mdParser.Convert([]byte(md), &buf); err != nil {
		handleError(err, fmt.Sprintf("Error converting markdown to HTML, %s", title))
	}

	return buf.String()
}
