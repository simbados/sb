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
		showErrorWith(err)
	}
	var defaultEditor, path string
	if len(output) < 2 {
		defaultEditor = "vim"
	} else {
		defaultEditor = string(output[:len(output)-1])
	}
	if commands[0] == "root" {
		path = filepath.Join(paths.RootConfigPath, commands[1]+".json")
		if exists, err := DoesPathExist(path); !exists {
			showErrorWith(err)
		}
	} else if commands[0] == "local" {
		path = filepath.Join(paths.LocalConfigPath, commands[1]+".json")
		if exists, err := DoesPathExist(path); !exists {
			showErrorWith(err)
		}
	} else {
		showError()
	}
	osHelper.Run([]string{defaultEditor, path})
	log.LogHighlight("Successfully edit of file")
	os.Exit(0)
}

func showErrorWith(error error) {
	log.LogErr(error)
}

func showError() {
	log.LogErr(fmt.Printf(`
You must provide two values to edit a configuration file
Valid options are:
sb -e local <binary name>
sb -e root <binary name>
`))
}
