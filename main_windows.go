package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

const (
	fontName string = "DejaVuSans"
)

func main() {
	conf := DefaultConfig
	flag.BoolVar(&conf.PrintVersion, "version", false, "Print version number")
	flag.Float64Var(&conf.FrameRate, "framerate", conf.FrameRate, "Fallback GOTV Framerate")
	flag.Float64Var(&conf.TickRate, "tickrate", conf.TickRate, "Fallback Gameserver Tickrate")
	instDirKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\csgoverview`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalln("trying to open csgoverview registry key:", err)
	}
	defer instDirKey.Close()
	instDir, _, err := instDirKey.GetStringValue("InstallLocation")
	if err != nil {
		log.Fatalln("trying to get install directory from registry key:", err)
	}
	defaultFontPath := fmt.Sprintf("%v\\%v.ttf", instDir, fontName)
	defaultOverviewDirectory := fmt.Sprintf("%v\\assets\\maps\\", instDir)
	flag.StringVar(&conf.FontPath, "fontpath", defaultFontPath, "Path to font file (.ttf)")
	flag.StringVar(&conf.OverviewDir, "overviewdir", defaultOverviewDirectory, "Path to overview directory")
	flag.Parse()

	err = run(&conf)
	if err != nil {
		log.Fatalln(err)
	}
}
