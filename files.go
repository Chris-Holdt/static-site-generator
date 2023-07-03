package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Index       string      `json:"index"`
	Content     string      `json:"contentDir"`
	Output      string      `json:"outputDir"`
	ArticleDir  string      `json:"articleDir"`
	HomeConfig  HomeConfig  `json:"homepageConfig"`
	StyleConfig StyleConfig `json:"styling"`
	HTMLConfig  HTMLConfig  `json:"htmlConfig"`

	ArticleList  []Article           `json:"-"`
	CategoryList map[string]Category `json:"-"`
}

type Article struct {
	Created  string
	Filename string
}

type Category struct {
	Title    string
	Articles []string
}

type HomeConfig struct {
	ShowArticles bool   `json:"showArticles"`
	ShowFeatured bool   `json:"showFeatured"`
	ArticleCount int    `json:"articlesToShow"`
	Featured     string `json:"featuredArticle"`
}

type StyleConfig struct {
	ScssDir string `json:"scssDir"`
}

type HTMLConfig struct {
	Top    string `json:"top"`
	Bottom string `json:"bottom"`
}

type Content struct {
	Filename string
	Raw      string
	Clean    string
	Metadata Metadata
	HTML     string
}

type Metadata struct {
	Title       string   `yaml:"title"`
	FileName    string   `yaml:"filename"`
	Date        string   `yaml:"date"`
	LastUpdate  string   `yaml:"last-update"`
	Tags        []string `yaml:"tags"`
	Categories  []string `yaml:"categories"`
	WordCount   int      `yaml:"wordcount"`
	ContentType string   `yaml:"content-type"`
}

func loadConfig() *Config {
	var c Config

	bConfig := readFile("config.json")

	if err := json.Unmarshal(bConfig, &c); err != nil {
		handleError(err, "Error processing config file")
	}

	c.ArticleList = []Article{}
	c.CategoryList = make(map[string]Category)

	return &c
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		handleError(err, fmt.Sprintf("There was an error reading the file: %s\n", path))
	}

	return data
}

func writeFile(path string, data []byte) {
	if err := os.WriteFile(path, data, 0644); err != nil {
		handleError(err, fmt.Sprintf("Error writing to file: %s", path))
	}
}

func createDir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		handleError(err, fmt.Sprintf("Error creating directory: %s", path))
	}
}
