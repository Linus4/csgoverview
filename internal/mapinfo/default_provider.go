package mapinfo

import (
	"encoding/json"
	"fmt"
	"os"

	_ "embed"
)

//go:embed assets/default_map_info.json
var defaultMapInfoJSON []byte

type record struct {
	Info         *MapInfo `json:"info"`
	RadarSources []string `json:"radar_sources"`
}

type DefaultMapProvider struct{}

func NewDefaultMapProvider() *DefaultMapProvider {
	return &DefaultMapProvider{}
}

func (p DefaultMapProvider) RetrieveInfo(mapName string, crc string) (*MapInfo, error) {
	var data map[string]record
	err := json.Unmarshal(defaultMapInfoJSON, &data)
	if err != nil {
		return nil, err
	}

	record, ok := data[mapName]
	if !ok {
		return nil, fmt.Errorf("local info for map %s not found", mapName)
	}

	mapFilesDir, err := MapFilesDir(mapName, "default")
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(mapFilesDir); err != nil {
		err = os.MkdirAll(mapFilesDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	// TODO: Check if the radar image is already locally available (should already be for windows instalation, for example)
	// TODO: Iterate over options of sources (if available) in case some fail
	radarTargetFile := fmt.Sprintf("%s/radar.png", mapFilesDir)
	err = fetchFile(record.RadarSources[0], radarTargetFile)
	if err != nil {
		return nil, err
	}

	mapInfo := record.Info
	mapInfo.RadarImagePath = radarTargetFile

	return mapInfo, nil
}
