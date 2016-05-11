package util

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// depth of start seed
// 种子起始深度
const (
	SEED_START_DEPTH = 0
)

type seedData []string

func GetSeedFromFile(path string, isDownload bool) (map[string]bool, error) {
	if !IsFileExist(path) {
		return nil, errors.New("file not exist")
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var decode seedData
	err = json.Unmarshal(data, &decode)
	if err != nil {
		return nil, err
	}

	result := map[string]bool{}
	for _, url := range decode {
		result[url] = isDownload
	}
	return result, nil
}
