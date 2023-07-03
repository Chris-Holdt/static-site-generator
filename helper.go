package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func isBlankStr(s string) bool {
	var b string

	return s == b
}

func hasDir(path string, dir string) bool {
	contents, err := os.ReadDir(path)
	if err != nil {
		handleError(err, fmt.Sprintf("Error reading directory contents of %s", path))
	}

	for _, c := range contents {
		if c.Name() == dir && c.IsDir() {
			return true
		}
	}

	return false

}

func handleError(e error, desc string) {
	log.Println(desc)
	log.Fatalln(e)
}

func titlecase(text string) string {
	first := text[:1]
	back := text[1:]

	first = strings.ToUpper(first)
	return fmt.Sprintf("%s%s", first, back)
}

func strToDate(sDate string) time.Time {
	split := strings.Split(sDate, "/")

	year, err := strconv.Atoi(split[2])
	if err != nil {
		handleError(err, fmt.Sprintf("Error converting date string to time.Time, failed to convert year: %s", sDate))
	}

	month, err := strconv.Atoi(split[1])
	if err != nil {
		handleError(err, fmt.Sprintf("Error converting date string to time.Time, failed to convert month: %s", sDate))
	}

	day, err := strconv.Atoi(split[0])
	if err != nil {
		handleError(err, fmt.Sprintf("Error converting date string to time.Time, failed to convert day: %s", sDate))
	}

	d := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)

	return d
}
