package service

import (
	"encoding/json"
	"os"
)

func ReadFileTo(path string, value interface{}) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}
