package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	fontName string = "DejaVuSans"
)

func parseArgs(conf *Config) *Config {
	flag.StringVar(&conf.FontPath, "fontpath", conf.FontPath, "Path to font file (.ttf)")
	flag.StringVar(&conf.OverviewDir, "overviewdir", conf.OverviewDir, "Path to overview directory")
	flag.BoolVar(&conf.PrintVersion, "version", conf.PrintVersion, "Print version number")
	flag.Parse()

	if conf.FontPath == "" {
		log.Fatalln("fontpath must be set.")
	}
	_, err := os.Stat(conf.FontPath)
	if os.IsNotExist(err) {
		log.Fatalf("Font file '%v' does not exist.", conf.FontPath)
	}

	if conf.OverviewDir == "" {
		log.Fatalln("overviewdir must be set")
	}
	_, err = os.Stat(conf.OverviewDir)
	if os.IsNotExist(err) {
		log.Fatalf("Directory with overviews '%v' does not exist.", conf.OverviewDir)
	}

	return conf
}

func unixDefaultFontPath() string {
	cmd := fmt.Sprintf("fc-list | grep %v.ttf", fontName)
	fontPathsB, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Printf("trying to find path to font: %v.ttf not installed on system\n", fontName)
	}
	return strings.Split(string(fontPathsB), ":")[0]
}

func unixDefaultOverviewDirectory() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("trying to get user home directory:", err)
	}
	return fmt.Sprintf("%v/.local/share/csgoverview/assets/maps", userHomeDir)
}
