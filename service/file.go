package service

import (
	"io/fs"
	"path/filepath"
)

// GetFileMapping
//
//	@param path
//	@return map[string]string
//	@return error
func GetFileMapping(path string) (map[string]string, error) {
	retData := map[string]string{}
	err := filepath.Walk(path, func(tmpPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Clean(path) == filepath.Clean(tmpPath) {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		name, _ := filepath.Rel(path, tmpPath)
		retData[name] = tmpPath
		return nil
	})

	return retData, err
}
