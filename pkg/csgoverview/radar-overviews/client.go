package radar_overviews

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const HOST = "https://radar-overviews.csgo.saiko.tech"

type MapInfo struct {
	PosX  float32 `json:"pos_x"`
	PosY  float32 `json:"pos_y"`
	Scale float32 `json:"scale"`
}

type Client struct {
	httpClient http.Client
	mapsDir    string
}

func NewClient() (*Client, error) {
	httpClient := http.Client{}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	mapsDir := fmt.Sprintf("%s/.local/share/csgoverview/assets/maps", userHomeDir)

	return &Client{httpClient, mapsDir}, err
}

func (c *Client) DownloadMapPNG(mapName string, crc string) error {
	path := fmt.Sprintf("%s/%s/%s/radar.png", HOST, mapName, crc)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	addHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch map image")
	}

	// Create target file
	// TODO(augustoccesar)[2022-12-18]:
	//  Currently is expected the images to be .png, but maybe should match the .png from the overviews
	mapFolder := fmt.Sprintf("%s/%s/%s", c.mapsDir, mapName, crc)
	err = os.MkdirAll(mapFolder, os.ModePerm)
	if err != nil {
		return err
	}

	targetPath := fmt.Sprintf("%s/radar.jpg", mapFolder)
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

func (c *Client) DownloadMapInfo(mapName string, crc string) error {
	path := fmt.Sprintf("%s/%s/%s/info.json", HOST, mapName, crc)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}
	addHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	mapFolder := fmt.Sprintf("%s/%s/%s", c.mapsDir, mapName, crc)
	err = os.MkdirAll(mapFolder, os.ModePerm)
	if err != nil {
		return err
	}

	targetPath := fmt.Sprintf("%s/info.json", mapFolder)
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

func (c *Client) ReadMapInfo(mapName string, crc string) (*MapInfo, error) {
	mapFolder := fmt.Sprintf("%s/%s/%s", c.mapsDir, mapName, crc)
	file, err := os.Open(fmt.Sprintf("%s/info.json", mapFolder))
	if err != nil {
		return nil, err
	}

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

func addHeaders(req *http.Request) {
	// TODO(augustoccesar)[2022-12-18]: Change to original repo path (or another user agent for better identification)
	req.Header.Add("User-Agent", "github.com/augustoccesar/csgoverview")
}
