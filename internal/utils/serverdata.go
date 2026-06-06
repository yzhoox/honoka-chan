package honokautils

import (
	"encoding/json"
	"honoka-chan/pkg/utils"
	"path/filepath"
)

const serverDataPath = "assets/serverdata"

func LoadServerData[T any](fileName string) (T, error) {
	var data T
	filePath := filepath.Join(serverDataPath, fileName)
	err := json.Unmarshal([]byte(utils.ReadAllText(filePath)), &data)
	return data, err
}
