package main

import (
	"log"
)

func main() {
	conf := parseArgs(&Config{
		FontPath:    unixDefaultFontPath(),
		OverviewDir: unixDefaultOverviewDirectory(),
	})
	err := run(conf)

	if err != nil {
		log.Fatalln(err)
	}
}
