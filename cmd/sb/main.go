package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sb/internal/osHelper"
	"sb/internal/parse"
	"sb/internal/sandbox"
	"sb/internal/types"
	"sb/internal/util"
)

func main() {
	input := os.Args
	context := types.Context{}

	// set relevant paths
	configAllPath(&context)

	// set config parameter
	setConfigParams(&context, input)
	// parse config files
	if context.Config.CliConfig == nil {
		context.Config.SbConfig = parse.ConfigFileParsing(&context)
	} else {
		util.LogDebug("Using cli options")
		context.Config.SbConfig = context.Config.CliConfig
	}

	// build sandbox profile
	profile := sandbox.BuildSandboxProfile(&context)

	util.LogDebug(util.PrettyJson(&context))

	// Run the sandbox
	args := append(append(append(append(append(append([]string{}, "sandbox-exec"), "-p"), profile), context.Config.BinaryName), context.Config.Commands...))
	osHelper.Run(args)
}

func setConfigParams(context *types.Context, args []string) {
	cliOptions, cliConfig, commands := parse.OptionsParsing(&context.Paths, args[1:])
	context.Config.CliConfig = cliConfig
	if len(cliOptions) != 0 {
		context.Config.CliOptions = cliOptions
	}
	context.Paths.BinaryPath = getPathToExecutable(commands[0])
	context.Config.BinaryName = commands[0]
	context.Config.Commands = commands[1:]
}

func getPathToExecutable(executableName string) string {
	cmd, err := exec.LookPath(executableName)
	if err != nil {
		util.LogErr(executableName + " binary does not exists\n")
	}
	return cmd
}

func configAllPath(context *types.Context) {
	context.Paths.HomePath = getHomePath()
	context.Paths.RootConfigPath = context.Paths.HomePath + types.ConfigRepo
	context.Paths.WorkingDir = getWorkingDir()
	context.Paths.BinPath = "/usr/bin"
	context.Paths.LocalConfigPath = context.Paths.WorkingDir + types.LocalConfigPath
	sbBinPath, err := os.Executable()
	if err != nil {
		util.LogErr(err)
	}
	context.Paths.SbBinaryPath = sbBinPath
}

func getHomePath() string {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Can not get the current user")
		os.Exit(1)
	}
	return currentUser.HomeDir
}

func getWorkingDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Can not get the working directory")
		os.Exit(1)
	}
	return currentDir
}
