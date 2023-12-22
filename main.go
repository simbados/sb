package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

type Config struct {
	BinaryPath    string    `json:"binary-path"`
	BinaryDirPath string    `json:"binary-path"`
	BinaryName    string    `json:"binary-name"`
	SbConfig      *SbConfig `json:"root-config"`
	Commands      []string  `json:"commands"`
	CliOptions    []string  `json:"cli-options"`
	CliConfig     *SbConfig `json:"cli-config"`
}

type Paths struct {
	HomePath        string `json:"home-path"`
	RootConfigPath  string `json:"root-config-path"`
	WorkingDir      string `json:"working-dir"`
	LocalConfigPath string `json:"local-config-path"`
	BinPath         string `json:"bin-path"`
	TargetBinPath   string `json:"target-bin-path"`
}

type CliOptionsStr struct {
	DebugEnabled     bool `json:"debug-enabled"`
	PrintEnabled     bool `json:"print-enabled"`
	DryRunEnabled    bool `json:"dry-run-enabled"`
	CreateExeEnabled bool `json:"create-exe-enabled"`
	HelpEnabled      bool `json:"help"`
	VersionEnabled   bool `json:"version-enabled"`
}

// Globals, might want to consider putting them in some state management, but meh
var CliOptions = CliOptionsStr{
	DebugEnabled:     false,
	PrintEnabled:     false,
	DryRunEnabled:    false,
	CreateExeEnabled: false,
}

const configRepo = "/.sb-conf"
const localConfigPath = "/.sb.conf"

var validCliOptions = map[string]string{
	"--debug":      "--debug",
	"-d":           "-d",
	"--print":      "--print",
	"-p":           "-p",
	"--dry-run":    "--dry-rin",
	"-dr":          "-dr",
	"--create-exe": "--create-exe",
	"-ce":          "-ce",
	"--help":       "--help",
	"-h":           "-h",
	"--version":    "--version",
	"-v":           "-v",
}

type Context struct {
	Config Config `json:"config"`
	Paths  Paths  `json:"path"`
}

var AllowedConfigKeys = map[string]string{
	"process":   "process",
	"write":     "write",
	"read":      "read",
	"net-in":    "net-in",
	"net-out":   "net-out",
	"--process": "--process",
	"--write":   "--write",
	"--read":    "--read",
	"--net-in":  "--net-in",
	"--net-out": "--net-out",
}

func main() {
	input := os.Args
	context := Context{}

	// set relevant paths
	configAllPath(&context)

	// set config parameter
	setConfigParams(&context, input)

	// parse config files
	parseConfigFiles(&context)

	// build sandbox profile
	buildSandboxProfile(&context)

	logInfo(prettyJson(context))
}

func setConfigParams(context *Context, args []string) {
	cliOptions, cliConfig, commands := parseOptions(&context.Paths, args[1:])
	context.Config.CliConfig = cliConfig
	if len(cliOptions) != 0 {
		context.Config.CliOptions = cliOptions
	}
	context.Config.BinaryPath = getPathToExecutable(commands[0])
	context.Config.BinaryDirPath = getDirForPath(context.Config.BinaryPath)
	context.Config.BinaryName = commands[0]
	context.Config.Commands = commands[1:]
}

func getDirForPath(path string) string {
	return filepath.Dir(path)
}

func getPathToExecutable(executableName string) string {
	cmd, err := exec.LookPath(executableName)
	if err != nil {
		logErr("%s binary does not exists\n", executableName)
	}
	return cmd
}

func configAllPath(context *Context) {
	context.Paths.HomePath = getHomePath()
	context.Paths.RootConfigPath = context.Paths.HomePath + configRepo
	context.Paths.WorkingDir = getWorkingDir()
	context.Paths.BinPath = "/bin"
	context.Paths.LocalConfigPath = context.Paths.WorkingDir + localConfigPath
}

func setOption(option string) {
	currentOption := validCliOptions[option]
	if currentOption == "--debug" || currentOption == "-d" {
		CliOptions.DebugEnabled = true
	} else if currentOption == "--print" || currentOption == "-p" {
		CliOptions.PrintEnabled = true
	} else if currentOption == "--dry-run" || currentOption == "-dr" {
		CliOptions.DryRunEnabled = true
	} else if currentOption == "--create-exe" || currentOption == "-ce" {
		CliOptions.CreateExeEnabled = true
	} else if currentOption == "--help" || currentOption == "-h" {
		printHelp()
	} else if currentOption == "--version" || currentOption == "-v" {
		showVersion()
	}
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

type SbConfig struct {
	Read            []string `json:"read"`
	Write           []string `json:"write"`
	ReadWrite       []string `json:"read-write"`
	Process         []string `json:"process"`
	NetworkOutbound bool     `json:"net-out"`
	NetworkInbound  bool     `json:"net-in"`
}
