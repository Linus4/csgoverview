package main

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

func main() {
	instDir := getInstallationDir()
	defaultFontPath := fontName + ".ttf"
	defaultOverviewDirectory := "./overviews/"
	if instDir != "" {
		defaultFontPath = fmt.Sprintf("%v\\%v.ttf", instDir, fontName)
		defaultOverviewDirectory = fmt.Sprintf("%v\\assets\\maps\\", instDir)
	}
	conf := parseArgs(&Config{
		FontPath:    defaultFontPath,
		OverviewDir: defaultOverviewDirectory,
	})

	err := run(conf)
	if err != nil {
		log.Fatalln(err)
	}
}

func getInstallationDir() string {
	instDirKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `Software\csgoverview`, registry.QUERY_VALUE)
	if err != nil {
		log.Println("Probably not an installation. Failed to open csgoverview registry key:", err)
		return ""
	}
	defer instDirKey.Close()

	instDir, _, err := instDirKey.GetStringValue("InstallLocation")
	if err != nil {
		log.Fatalln("Failed to get install directory from registry key:", err)
	}

	return instDir
}
