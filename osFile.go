package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func doesPathExist(filePath string) (bool, string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, "File or directory does not exist"
	} else if err != nil {
		return false, fmt.Sprintf("File could not be checked + %v", err)
	} else {
		return true, ""
	}
}

func parseJson(path string) map[string]interface{} {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		logErr("Error while reading file at path: ", path)
		os.Exit(1)
	}
	var jsonFile map[string]interface{}

	if err := json.Unmarshal(fileContent, &jsonFile); err != nil {
		logErr("Error while parsing json file at path: ", path, err)
		os.Exit(1)
	}
	return jsonFile
}
