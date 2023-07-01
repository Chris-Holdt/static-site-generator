package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type ContentTag string

// Looks for matching //-- >< tag
var (
	FEATURED_ARTICLE ContentTag = "//-- FEATURED_ARTICLE"
	LATEST_ARTICLES  ContentTag = "//-- LATEST_ARTICLES"
	META             ContentTag = "//-- META"
	END              ContentTag = "//-- ><"

	ContentTags = []ContentTag{
		FEATURED_ARTICLE,
		LATEST_ARTICLES,
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
* 3. Convert the SCSS to CSS (Can be done in parallel?)
*
* 4. Generate the articles page with a list of the articles we found
 */

func startProcessing(config *Config) {
	mdPaths := findMarkdownDocs(config.Content)
	content := loadContent(mdPaths)

	content = processContent(content)

	// for _, x := range content {
	// 	fmt.Println(x.Raw)
	// }
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

		c.Raw = string(readFile(p))

		contentStore = append(contentStore, c)
	}

	return contentStore
}

func processContent(content []Content) []Content {

	for _, c := range content {

		raw := c.Raw

		for _, t := range ContentTags {

			tagStart := strings.Index(raw, string(t))
			if tagStart < 0 {
				continue
			}

			tagEnd := strings.Index(raw[tagStart:], string(END))
			tagEnd += tagStart

			sub := raw[tagStart:tagEnd]

			c.Raw = strings.Replace(raw, sub, "", -1)

			sub = strings.Replace(sub, string(t), "", -1)
			sub = strings.Replace(sub, string(END), "", -1)

			// fmt.Printf("\nTAG CONTENT\n%s", sub)
			fmt.Printf("\nCONTENT\n%s", c.Raw)
		}

	}

	return content
}
