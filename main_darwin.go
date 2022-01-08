package main

import (
	"log"
)

func main() {
	conf := parseArgs(unixDefaultFontPath(), unixDefaultOverviewDirectory())
	err = run(&conf)
	if err != nil {
		log.Fatalln(err)
	}
}
