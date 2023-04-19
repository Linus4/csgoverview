package mapinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

type SaikoTechMapProvider struct{}

func NewSaikoTechMapProvider() *SaikoTechMapProvider {
	return &SaikoTechMapProvider{}
}

func (p SaikoTechMapProvider) RetrieveInfo(mapName string, crc string) (*MapInfo, error) {
	mapFilesDir, err := MapFilesDir(mapName, crc)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(mapFilesDir); err != nil {
		err = os.MkdirAll(mapFilesDir, os.ModePerm)
		if err != nil {
			return nil, err
		}

		if err = p.downloadMapFiles(mapName, crc, mapFilesDir); err != nil {
			return nil, err
		}
	}

	mapInfoFile, err := os.Open(fmt.Sprintf("%s/info.json", mapFilesDir))
	if err != nil {
		return nil, err
	}

	mapInfo, err := p.readMapInfo(mapName, mapInfoFile)
	if err != nil {
		return nil, err
	}
	mapInfo.RadarImagePath = fmt.Sprintf("%s/radar.png", mapFilesDir)

	return mapInfo, nil
}

func (p SaikoTechMapProvider) downloadMapFiles(mapName string, crc string, targetFolder string) error {
	const HOST = "https://radar-overviews.csgo.saiko.tech"
	basePath := fmt.Sprintf("%s/%s/%s", HOST, mapName, crc)

	files := []string{"radar.png", "info.json"}
	for _, file := range files {
		source := fmt.Sprintf("%s/%s", basePath, file)
		target := fmt.Sprintf("%s/%s", targetFolder, file)
		err := fetchFile(source, target)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p SaikoTechMapProvider) readMapInfo(mapName string, file *os.File) (*MapInfo, error) {
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data map[string]map[string]string
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	posX, err := strconv.ParseFloat(data[mapName]["pos_x"], 32)
	if err != nil {
		return nil, err
	}

	posY, err := strconv.ParseFloat(data[mapName]["pos_y"], 32)
	if err != nil {
		return nil, err
	}

	scale, err := strconv.ParseFloat(data[mapName]["scale"], 32)
	if err != nil {
		return nil, err
	}

	mapInfo := MapInfo{
		PosX:  float32(posX),
		PosY:  float32(posY),
		Scale: float32(scale),
	}

	return &mapInfo, nil
}
