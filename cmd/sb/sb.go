package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sb/internal/log"
	"sb/internal/osHelper"
	"sb/internal/parse"
	"sb/internal/sandbox"
	"sb/internal/types"
	"sb/internal/util"
	"strconv"
)

func main() {
	input := os.Args
	context := types.Context{}

	// set relevant paths
	configAllPath(&context)

	parseEnvs()

	// set config parameter
	setConfigParams(&context, input)
	// parse config files
	if context.Config.CliConfig == nil {
		context.Config.SbConfig = parse.ConfigFileParsing(&context)
	} else {
		log.LogDebug("Using cli options")
		context.Config.SbConfig = context.Config.CliConfig
	}

	// build sandbox profile
	profile := sandbox.BuildSandboxProfile(&context)

	log.LogDebug(log.PrettyJson(&context))

	if !types.CliOptions.DryRunEnabled {
		// Run the sandbox
		args := append(append(append(append(append(append([]string{}, "sandbox-exec"), "-p"), profile), context.Config.BinaryName), context.Config.Commands...))
		osHelper.Run(args)
	}
}

func parseEnvs() {
	for _, env := range types.AllEnvs {
		val := os.Getenv(env)
		if env == types.DEV_MODE && val != "" {
			if val, err := strconv.ParseBool(val); err == nil {
				types.Envs.DevModeEnabled = val
			} else {
				log.LogErr(fmt.Sprintf("Can not parse %v env variable please provide true or false", types.DEV_MODE))
			}
		}
	}
}

func setConfigParams(context *types.Context, args []string) {
	cliOptions, cliConfig, commands := parse.OptionsParsing(&context.Paths, args[1:])
	if types.CliOptions.EditEnabled {
		util.EditFile(commands, context.Paths)
	}
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
		log.LogErr(executableName + " binary does not exists\n")
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
		log.LogErr(err)
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
