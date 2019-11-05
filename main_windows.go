package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	fontName string = "DejaVuSans"
)

func main() {
	conf := DefaultConfig
	flag.Float64Var(&conf.FrameRate, "framerate", conf.FrameRate, "Fallback GOTV Framerate")
	flag.Float64Var(&conf.TickRate, "tickrate", conf.TickRate, "Fallback Gameserver Tickrate")
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("trying to get user home directory:", err)
	}
	defaultFontPath := fmt.Sprintf("%v\\csgoverview\\%v.ttf", userHomeDir, fontName)
	defaultOverviewDirectory := fmt.Sprintf("%v\\csgoverview\\", userHomeDir)
	flag.StringVar(&conf.FontPath, "fontpath", defaultFontPath, "Path to font file (.ttf)")
	flag.StringVar(&conf.OverviewDir, "overviewdir", defaultOverviewDirectory, "Path to overview directory")
	flag.Parse()

	err = run(&conf)
	if err != nil {
		log.Fatalln(err)
	}
}
