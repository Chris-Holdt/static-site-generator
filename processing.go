package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type ContentTag string

// Looks for matching //-- >< tag
var (
	FEATURED_ARTICLE ContentTag = "//-- FEATURED_ARTICLE"
	LATEST_ARTICLES  ContentTag = "//-- LATEST_ARTICLES"
	LIST_ARTICLES    ContentTag = "//-- LIST_ARTICLES"
	META             ContentTag = "//-- META"
	END              ContentTag = "//-- ><"

	ContentTags = []ContentTag{
		FEATURED_ARTICLE,
		LATEST_ARTICLES,
		LIST_ARTICLES,
		META,
	}
)

/*
* 1. Create list of content to parse, don't care about dir structure
*    the files should contain metadata that tells us where to put them
*    We will need to recurse
*
* 2. Convert all the markdown files to HTML and link them correctly
*
* 3. Convert the SCSS to CSS (Can be done in parallel?) - Not sure if I will do this, looks like more work than it is worth
*
* 4. Generate the articles page with a list of the articles we found
 */

func startProcessing(config *Config) {
	if !hasDir(".", config.Output) {
		createDir(config.Output)
	}

	outputPath := fmt.Sprintf("./%s", config.Output)
	if !hasDir(outputPath, config.ArticleDir) {
		createDir(fmt.Sprintf("%s/%s", outputPath, config.ArticleDir))
	}

	mdPaths := findMarkdownDocs(config.Content)
	content := loadContent(mdPaths)

	content = getAllMetadata(content, config)

	content = processContent(content, config)

	content = convertToHtml(content, config)

	for _, c := range content {
		if len(c.Metadata.Title) <= 0 {
			continue
		}

		fileName := fmt.Sprintf("%s.html", c.Metadata.FileName)
		fileName = strings.ToLower(fileName)

		if c.Metadata.ContentType == "article" {
			fileName = fmt.Sprintf("articles/%s", fileName)
		}

		filePath := fmt.Sprintf("%s/%s", config.Output, fileName)

		writeFile(filePath, []byte(c.HTML))
	}

	handleCss(config)

}

func handleCss(config *Config) {
	out := fmt.Sprintf("./%s", config.Output)
	if !hasDir(out, "styles") {
		createDir(fmt.Sprintf("%s/styles", out))
	}

	css := readFile("./styles/main.css")
	writeFile(fmt.Sprintf("%s/styles/main.css", out), css)

}

func findMarkdownDocs(baseDir string) []string {
	paths := []string{}

	err := filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if !info.IsDir() {

			name := info.Name()
			fileType := name[len(name)-3:]
			if fileType == ".md" {
				paths = append(paths, path)
			}
		}

		return nil
	})

	if err != nil {
		handleError(err, fmt.Sprintf("Error reading content dir: %s\n", baseDir))
	}

	return paths
}

func loadContent(paths []string) []Content {
	contentStore := []Content{}

	for _, p := range paths {
		c := Content{}

		pathSplit := strings.Split(p, "/")
		c.Filename = pathSplit[len(pathSplit)-1]

		c.Raw = string(readFile(p))

		contentStore = append(contentStore, c)
	}

	return contentStore
}

func getAllMetadata(content []Content, config *Config) []Content {

	for i, c := range content {
		raw := c.Raw

		tagStart := strings.Index(raw, string(META))
		if tagStart < 0 {
			fmt.Printf("File: %s doesn't have Metadata, is this a mistake?\n", c.Filename)
			continue
		}

		tagEnd := strings.Index(raw[tagStart:], string(END))
		tagEnd += len(END) + tagStart

		content[i].Metadata, raw = getMeta(tagStart, tagEnd, raw)

		if content[i].Metadata.ContentType == "article" {
			addArticleDataToConfig(content[i].Metadata, config)
		}

		content[i].Raw = raw

	}

	return content
}

func processContent(content []Content, config *Config) []Content {

	for i, c := range content {

		raw := c.Raw

		for _, t := range ContentTags {

			tagStart := strings.Index(raw, string(t))
			if tagStart < 0 {
				continue
			}

			tagEnd := strings.Index(raw[tagStart:], string(END))
			tagEnd += len(END) + tagStart

			switch t {
			case FEATURED_ARTICLE:
				raw = injectFeatured(tagStart, tagEnd, raw, config)

			case LATEST_ARTICLES:
				raw = injectLatest(tagStart, tagEnd, raw, config)

			case LIST_ARTICLES:
				raw = injectList(tagStart, tagEnd, raw, config)

				// fmt.Println(raw)
			}

		}

		content[i].Clean = raw

		if c.Metadata.ContentType == "article" {
			content[i] = processArticle(content[i])
		}

	}

	return content
}

func getMeta(s int, e int, raw string) (Metadata, string) {
	m := Metadata{}

	metaTagLen := len(META)
	endTagLen := len(END)

	yamlStartIndex := s + metaTagLen
	yamlEndIndex := e - endTagLen

	metaStr := raw[yamlStartIndex:yamlEndIndex]

	if err := yaml.Unmarshal([]byte(metaStr), &m); err != nil {
		handleError(err, "Error unmarshalling metadata to yaml")
	}

	// Debug
	// fmt.Printf("\nMETA\n%+v\n", m)

	// Remove the processed data
	raw = strings.Replace(raw, raw[s:e], "", -1)

	return m, raw
}
