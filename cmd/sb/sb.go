package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"sb/internal/log"
	"sb/internal/osHelper"
	"sb/internal/sandbox"
	"sb/internal/types"
	"sb/internal/util"
	"strconv"
	"strings"
)

func main() {
	input := os.Args
	context := types.Context{}

	// set relevant paths
	configAllPath(&context)

	parseEnvs()

	// set config parameter
	setConfigParams(&context, input)

	if context.Config.CliConfig == nil {
		context.Config.SbConfig = util.ConfigFileParsing(&context)
	} else {
		log.LogDebug("Using cli options")
		context.Config.SbConfig = context.Config.CliConfig
	}

	minifiedProfile, profile := sandbox.BuildSandboxProfile(&context)

	if types.CliOptions.ShowConfigEnabled {
		util.ShowConfig(&context.Config, profile)
	}

	log.LogDebug(log.PrettyJson(&context))

	if !types.CliOptions.DryRunEnabled {
		// Run the sandbox
		args := append(append(append(append(append(append([]string{}, "sandbox-exec"), "-p"), minifiedProfile), context.Config.BinaryName), context.Config.Commands...))
		log.LogHighlight("Running sandbox exec with following command")
		if types.CliOptions.VigilantModeEnabled {
			log.LogInfoLn("The sandbox profile:")
			log.LogInfoLn(profile)
			log.LogHighlight("Do you want to run the command")
			log.LogHighlight(strings.Join(args[3:], " "))
			log.LogHighlightSl("(Y)es/(n)o ")
			var answer string
			_, err := fmt.Scanln(&answer)
			if err != nil {
				log.LogErr("Can not read your input")
			}
			if answer != "y" {
				log.LogInfoLn("Exiting program without running binary")
				os.Exit(0)
			}
		} else {
			log.LogHighlight(strings.Join(args[3:], " "))
		}
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

func checkCliOptions(context *types.Context, commands []string) {
	for currentOption, value := range context.Config.CliOptions {
		if currentOption == "--debug" || currentOption == "-d" {
			types.CliOptions.DebugEnabled = true
		} else if currentOption == "--dry-run" || currentOption == "-dr" {
			types.CliOptions.DryRunEnabled = true
		} else if currentOption == "--help" || currentOption == "-h" {
			onlyThisOptionAllowed(len(context.Config.CliOptions), currentOption)
			util.PrintHelp()
		} else if currentOption == "--version" || currentOption == "-v" {
			onlyThisOptionAllowed(len(context.Config.CliOptions), currentOption)
			util.ShowVersion()
		} else if currentOption == "--init" || currentOption == "-i" {
			onlyThisOptionAllowed(len(context.Config.CliOptions), currentOption)
			util.Init(&context.Paths)
		} else if currentOption == "--edit" || currentOption == "-e" {
			onlyThisOptionAllowed(len(context.Config.CliOptions), currentOption)
			util.EditFile(commands, &context.Paths)
		} else if currentOption == "--show" || currentOption == "-s" {
			onlyThisOptionAllowed(len(context.Config.CliOptions), currentOption)
			types.CliOptions.ShowConfigEnabled = true
		} else if currentOption == "--vigilant" || currentOption == "-vi" {
			types.CliOptions.VigilantModeEnabled = true
		} else if currentOption == "--config" || currentOption == "-c" {
			fmt.Println(currentOption, value)
			if len(value) != 2 {
				log.LogErr("You can not call config mode without providing local, root or path as config source")
			}
			types.CliOptions.ConfigModeEnabled = value[1]
		}
	}
}

func onlyThisOptionAllowed(length int, option string) {
	if length != 1 {
		log.LogErr(fmt.Sprintf("%v option can only be used in isolation and not in combination with other options", option))
	}
}

func setConfigParams(context *types.Context, args []string) {
	cliOptions, cliConfig, commands := util.OptionsParsing(&context.Paths, args[1:])
	context.Config.CliConfig = cliConfig
	if len(cliOptions) != 0 {
		context.Config.CliOptions = cliOptions
	}
	// Check which options are enabled from the user
	checkCliOptions(context, commands)
	if len(commands) < 1 {
		log.LogErr("Please specify the binary, no command provided")
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
