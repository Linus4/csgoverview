package mapinfo

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type MapInfo struct {
	RadarImagePath string
	PosX           float32 `json:"pos_x"`
	PosY           float32 `json:"pos_y"`
	Scale          float32 `json:"scale"`
}

type MapProvider interface {
	RetrieveInfo(mapName string, crc string) (*MapInfo, error)
}

func ResolveMapProvider() MapProvider {
	// TODO: Maybe resolve to default or Saiko based on network?
	return NewDefaultMapProvider()
	// return NewSaikoTechMapProvider()
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

func fetchFile(sourcePath string, targetPath string) error {
	req, err := http.NewRequest("GET", sourcePath, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", "github.com/Linus4/csgoverview")

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
