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

func main() {
	conf := DefaultConfig
	flag.BoolVar(&conf.PrintVersion, "version", false, "Print version number")
	cmd := fmt.Sprintf("fc-list | grep %v.ttf", fontName)
	fontPathsB, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Printf("trying to find path to font: %v.ttf not installed on system\n", fontName)
	}
	defaultFontPath := strings.Split(string(fontPathsB), ":")[0]
	flag.StringVar(&conf.FontPath, "fontpath", defaultFontPath, "Path to font file (.ttf)")
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("trying to get user home directory:", err)
	}
	defaultOverviewDirectory := fmt.Sprintf("%v/.local/share/csgoverview/assets/maps", userHomeDir)
	flag.StringVar(&conf.OverviewDir, "overviewdir", defaultOverviewDirectory, "Path to overview directory")
	flag.Parse()

	err = run(&conf)
	if err != nil {
		log.Fatalln(err)
	}
}
