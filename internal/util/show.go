package util

import (
	"fmt"
	"os"
	"path/filepath"
	"sb/internal/log"
	"sb/internal/types"
)

func ShowConfig(context *types.Context, binaryName string) {
	localConfigPath, exists := LocalConfigPath(&context.Paths, binaryName)
	if !exists {
		log.LogInfoLn("Local Config not found, trying root config")
		path := filepath.Join(context.Paths.RootConfigPath, binaryName+".json")
		configExists, err := DoesPathExist(path)
		if err != nil {
			log.LogErr(err)
		}
		if configExists {
			log.LogHighlight(fmt.Sprintf("Root config found at path %v", path))
			configJson := ParseJson(path)
			log.LogInfoLn(log.PrettyJson(configJson))
		}
	} else {
		log.LogHighlight(fmt.Sprintf("Local config found at path %v", localConfigPath))
		configJson := ParseJson(localConfigPath)
		log.LogInfoLn(log.PrettyJson(configJson))
	}
	os.Exit(0)
}
