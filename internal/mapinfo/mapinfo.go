package mapinfo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type MapInfo struct {
	RadarImagePath string
	PosX           float32 `json:"pos_x"`
	PosY           float32 `json:"pos_y"`
	Scale          float32 `json:"scale"`
}

func MapFilesDir(mapName string, crc string) (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"%s/.local/share/csgoverview/assets/maps/%s/%s",
		userHomeDir,
		mapName,
		crc,
	), nil
}

func ResolveMapInfo(mapName string, crc string) (*MapInfo, error) {
	mapFilesDir, err := MapFilesDir(mapName, crc)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(mapFilesDir); err != nil {
		err = os.MkdirAll(mapFilesDir, os.ModePerm)
		if err != nil {
			return nil, err
		}

		if err = downloadMapFiles(mapName, crc, mapFilesDir); err != nil {
			return nil, err
		}
	}

	mapInfoFile, err := os.Open(fmt.Sprintf("%s/info.json", mapFilesDir))
	if err != nil {
		return nil, err
	}

	mapInfo, err := readMapInfo(mapName, mapInfoFile)
	if err != nil {
		return nil, err
	}
	mapInfo.RadarImagePath = fmt.Sprintf("%s/radar.png", mapFilesDir)

	return mapInfo, nil
}

func readMapInfo(mapName string, file *os.File) (*MapInfo, error) {
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

func downloadMapFiles(mapName string, crc string, targetFolder string) error {
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

func fetchFile(sourcePath string, targetPath string) error {
	req, err := http.NewRequest("GET", sourcePath, nil)
	if err != nil {
		return err
	}
	// TODO(augustoccesar)[2022-12-18]: Change to original repo path (or another user agent for better identification)
	req.Header.Add("User-Agent", "github.com/augustoccesar/csgoverview")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"failed to fetch remote file from %s. Status: %d",
			sourcePath,
			resp.StatusCode,
		)
	}

	file, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
