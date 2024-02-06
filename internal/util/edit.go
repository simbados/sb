package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sb/internal/log"
	"sb/internal/osHelper"
	"sb/internal/types"
)

func EditFile(commands []string, paths types.Paths) {
	if len(commands) != 2 {
		showError()
	}
	cmd := exec.Command("sh", "-c", "echo $EDITOR")
	output, err := cmd.Output()
	if err != nil {
		showErrorIfError(err)
	}
	var defaultEditor, path string
	if len(output) < 2 {
		defaultEditor = "vim"
	} else {
		defaultEditor = string(output[:len(output)-1])
	}
	if commands[0] == "root" {
		path = getPath(commands[1], paths.RootConfigPath, path)
	} else if commands[0] == "local" {
		path = getPath(commands[1], paths.LocalConfigPath, path)
	} else {
		showError()
	}
	osHelper.Run([]string{defaultEditor, path})
	log.LogHighlight("Successfully edit of file")
	os.Exit(0)
}

func getPath(binaryName string, configPath, path string) string {
	path = filepath.Join(configPath, binaryName+".json")
	if exists, _ := DoesPathExist(path); !exists {
		var answer string
		log.LogInfoSl(fmt.Sprintf("File does not exist, do you want to create it at path %v? (y)es/(n)o ", path))
		if _, scanErr := fmt.Scanln(&answer); scanErr != nil {
			showErrorIfError(scanErr)
		} else {
			if configPathExists, err := DoesPathExist(configPath); !configPathExists {
				showErrorIfError(err)
			}
		}
	}
	return path
}

func showErrorIfError(error error) {
	if error != nil {
		log.LogErr(error)
	}
}

func showError() {
	log.LogErr(fmt.Printf(`
You must provide two values to edit a configuration file
Valid options are:
sb -e local <binary name>
sb -e root <binary name>
`))
}
