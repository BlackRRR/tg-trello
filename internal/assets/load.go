package assets

import (
	"encoding/json"
	"os"
)

func LoadTexts(path string) (map[string]string, error) {
	texts := make(map[string]string)
	ru, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(ru, &texts)
	if err != nil {
		return nil, err
	}

	return texts, nil
}
