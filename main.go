package main

import "log"

func main() {

	config := loadConfig()
	startProcessing(config)

}

func handleError(e error, desc string) {
	log.Println(desc)
	log.Fatalln(e)
}
