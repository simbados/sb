package sandbox

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sb/internal/types"
	"sb/internal/util"
)

func Init(paths *types.Paths) {
	sbConfigPath := filepath.Join(paths.HomePath, ".sb-config")
	util.LogInfoLn("Init of sandbox called")
	util.LogInfoSl(fmt.Sprintf("Do you want to have a root config with default values at %v, (y/n) ", sbConfigPath))
	var answer string
	_, err := fmt.Scanln(&answer)
	handleError(err, "Sorry your input can not be read, please type y or n \nExiting")
	if answer == "y" {
		createSbConfig(sbConfigPath)
	} else {
		util.LogWarn("Skipping creating root config directory")
	}
	util.LogInfoSl(fmt.Sprintf("Do you want to move the binary to the root config in %v and add it to your PATH or keep it at the current location and add it to the path manually? (y/n) ", sbConfigPath))
	var moveLocation string
	_, errMove := fmt.Scanln(&moveLocation)
	handleError(errMove, "Sorry your input can not be read, please type y or n \nExiting")
	if moveLocation == "y" {
		moveBinaryAndAddToPath(paths, sbConfigPath)
		util.LogInfoLn("Please source your shell config, so that you can use sb")
		util.LogInfoLn("source ~/.zshrc")
	} else {
		util.LogWarn("Skipping moving the binary and adding it to path")
		util.LogWarn("To use sb you have to add the binary to your path manually")
	}
	os.Exit(0)
}

func moveBinaryAndAddToPath(paths *types.Paths, sbConfigPath string) {
	binPath := filepath.Join(sbConfigPath, "bin")
	createDir(binPath)
	err := os.Rename(paths.SbBinaryPath, filepath.Join(binPath, "sb"))
	handleError(err, "Could not move binary to new location")
	localShellConfigPath := filepath.Join(paths.HomePath, ".zshrc")
	file, err := os.OpenFile(localShellConfigPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	handleError(err, "Could not open shell configuration file")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			handleError(err, "")
		}
	}(file)
	if _, err := file.WriteString(fmt.Sprintf("export PATH=$PATH:%v\n", filepath.Join(sbConfigPath, "bin"))); err != nil {
		handleError(err, "Failed to write to file")
	}
}

//go:embed configs
var configDir embed.FS

func createSbConfig(sbConfigPath string) {
	util.LogInfoLn("Creating .sb-config...")
	if val, _ := util.DoesPathExist(sbConfigPath); val {
		util.LogWarn(fmt.Sprintf("There is already an existing config directory at %v \nPlease create a backup or remove it to have the default configuration", sbConfigPath))
	} else {
		util.LogInfoLn(fmt.Sprintf("Creating directory at destination %v", sbConfigPath))
		createDir(sbConfigPath)
		err := fs.WalkDir(configDir, "configDir", func(path string, d fs.DirEntry, err error) error {
			if err := copyEmbeddedFilesToDestination(sbConfigPath); err != nil {
				util.LogErr("Failed to extract embedded files")
				util.LogErr(err)
				util.LogErr("Please submit an issue in the github repository")
			}
			util.LogInfoLn("sb config successfully copied to destination")
			return nil
		})
		handleError(err, "Something went wrong wile copying files")
	}
}

func createDir(path string) {
	cmd := exec.Command("mkdir", path)
	if err := cmd.Run(); err != nil {
		handleError(err, fmt.Sprintf("Could not create directory at %v", path))
	}
}

func copyEmbeddedFilesToDestination(configPath string) error {
	return fs.WalkDir(configDir, "configs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "configs" {
			return nil
		}
		outPath := filepath.Join(configPath, d.Name())
		data, err := fs.ReadFile(configDir, path)
		if err != nil {
			return err
		}
		return os.WriteFile(outPath, data, os.FileMode(0o644))
	})
}

func handleError(err error, message string) {
	if err != nil {
		fmt.Println(message, "Error: ", err)
		os.Exit(1)
	}
}
