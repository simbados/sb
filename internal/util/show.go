package util

import (
	"os"
	"sb/internal/log"
	"sb/internal/types"
)

func ShowConfig(config *types.Config, profile string) {
	log.LogHighlight("\nSandbox profile which will be applied\n")
	log.LogInfoLn(profile)
	log.LogHighlight("\nConfig parsed from config files: \n")
	log.LogInfoLn(log.PrettyJson(config.SbConfig))
	os.Exit(0)
}
