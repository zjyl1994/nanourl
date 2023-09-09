package util

import (
	"encoding/json"
	"errors"
	"os"
)

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

func LoadJson[T any](filename string) (T, error) {
	var result T
	data, err := os.ReadFile(filename)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(data, &result)
	return result, err
}
