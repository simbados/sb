package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sb/internal/log"
)

func DoesPathExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, errors.New("file or directory does not exist")
	} else if err != nil {
		return false, errors.New(fmt.Sprintf("File could not be checked + %v", err))
	} else {
		return true, nil
	}
}

func ParseJson(path string) map[string]interface{} {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		log.LogErr("Error while reading file at path: ", path)
		os.Exit(1)
	}
	var jsonFile map[string]interface{}

	if err := json.Unmarshal(fileContent, &jsonFile); err != nil {
		log.LogErr("Error while parsing json file at path: ", path, err)
		os.Exit(1)
	}
	return jsonFile
}
