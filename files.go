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
	ArticleList string      `json:"articleList"`
	ArticleDir  string      `json:"articleDir"`
	HomeConfig  HomeConfig  `json:"homepageConfig"`
	StyleConfig StyleConfig `json:"styling"`
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

type Content struct {
	Raw      string
	Clean    string
	Metadata Metadata
	HTML     string
}

type Metadata struct {
	Title       string `yaml:"title"`
	Tags        string `yaml:"tags"`
	Category    string `yaml:"category"`
	WordCount   int    `yaml:"wordcount"`
	ContentType string `yaml:"content-type"`
}

func loadConfig() *Config {
	var c Config

	bConfig := readFile("config.json")

	if err := json.Unmarshal(bConfig, &c); err != nil {
		handleError(err, "Error processing config file")
	}

	return &c
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		handleError(err, fmt.Sprintf("There was an error reading the file: %s\n", path))
	}

	return data
}
