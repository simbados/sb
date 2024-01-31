package parse

import (
	"regexp"
	"sb/internal/sandbox"
	"sb/internal/types"
	"sb/internal/util"
	"strings"
)

func setOption(paths *types.Paths, option string) {
	currentOption := types.ValidCliOptions[option]
	if currentOption == "--debug" || currentOption == "-d" {
		types.CliOptions.DebugEnabled = true
	} else if currentOption == "--print" || currentOption == "-p" {
		types.CliOptions.PrintEnabled = true
	} else if currentOption == "--dry-run" || currentOption == "-dr" {
		types.CliOptions.DryRunEnabled = true
	} else if currentOption == "--create-exe" || currentOption == "-ce" {
		types.CliOptions.CreateExeEnabled = true
	} else if currentOption == "--help" || currentOption == "-h" {
		util.PrintHelp()
	} else if currentOption == "--version" || currentOption == "-v" {
		util.ShowVersion()
	} else if currentOption == "--init" || currentOption == "-i" {
		sandbox.Init(paths)
	}
}

func OptionsParsing(paths *types.Paths, args []string) ([]string, *types.SbConfig, []string) {
	var options []string
	var cliConfigSb types.SbConfig
	var cliConfig *types.SbConfig = nil
	re := regexp.MustCompile("^-.*")
	optionsUntilIndex := 0
	for index, value := range args {
		if re.MatchString(value) {
			split, splitValue := parseCliConfigParam(value)
			if _, exist := types.ValidCliOptions[split]; exist {
				options = append(options, value)
				setOption(paths, value)
			} else if _, configExists := types.AllowedConfigKeys[split]; configExists && len(splitValue) > 0 {
				cliConfig = &cliConfigSb
				if splitValue == "true" || splitValue == "false" {
					addToConfig(cliConfig, split, splitValue)
				} else if arr := strings.Split(splitValue, ","); len(arr) > 0 {
					for _, val := range arr {
						addToConfig(cliConfig, split, expandPaths(paths, val))
					}
				}
			} else {
				util.LogErr("You passed a wrong cli option: ", value)
			}
		} else {
			optionsUntilIndex = index
			break
		}
	}
	if len(options) == len(args) {
		util.LogErr("Please specify the program that you want to sandbox")
	}
	if cliConfig != nil {
		return options, cliConfig, args[optionsUntilIndex:]
	} else {
		return options, nil, args[optionsUntilIndex:]
	}
}

func parseStringBoolean(s string) (bool, bool) {
	if s == "true" {
		return true, true
	} else if s == "false" {
		return false, true
	}
	return false, false
}

func addToConfig(config *types.SbConfig, key string, value string) *types.SbConfig {
	switch key {
	case "--read":
		config.Read = append(config.Read, value)
		break
	case "--write":
		config.Write = append(config.Write, value)
		break
	case "--process":
		config.Process = append(config.Process, value)
		break
	case "--net-in":
		boolVal, exists := parseStringBoolean(value)
		if !exists {
			util.LogErr("You must provide true or false value for cli config: ", value)
		}
		config.NetworkInbound = boolVal
		break
	case "--net-out":
		boolVal, exists := parseStringBoolean(value)
		if !exists {
			util.LogErr("You must provide true or false value for cli config: ", value)
		}
		config.NetworkOutbound = boolVal
		break
	}
	return config
}

func parseCliConfigParam(s string) (string, string) {
	split := strings.SplitN(s, "=", 2)
	if len(split) == 2 {
		return split[0], split[1]
	}
	return split[0], ""
}

const dotReplace = "\\."
const dotRegex = `\.`
const globEndRegex = `\*\*$`
const globEndReplace = `(.+)`
const globMiddleRegex = `\*\*\/`
const globMiddleReplace = `(.+\/)?`
const globSingleRegex = `\*`

const globSingleReplace = `^[^\/]*`

func expandPaths(paths *types.Paths, value string) string {
	value = strings.ReplaceAll(value, "~", paths.HomePath)
	value = strings.ReplaceAll(value, "[home]", paths.HomePath)
	if strings.HasPrefix(value, "./") {
		value = strings.ReplaceAll(value, "./", paths.WorkingDir+"/")
	}
	value = strings.ReplaceAll(value, "[wd]", paths.WorkingDir)
	value = strings.ReplaceAll(value, "[bin]", paths.BinPath)
	if strings.Contains(value, "*") || strings.Contains(value, ".") {
		value = regexp.MustCompile(dotRegex).ReplaceAllString(value, dotReplace)
		if val, err := regexp.Compile(globEndRegex); err == nil && val.MatchString(value) {
			value = val.ReplaceAllString(value, globEndReplace)
		} else if val, err := regexp.Compile(globMiddleRegex); err == nil && val.MatchString(value) {
			value = val.ReplaceAllString(value, globMiddleReplace)
		} else if val, err := regexp.Compile(globSingleRegex); err == nil && val.MatchString(value) {
			value = val.ReplaceAllString(value, globSingleReplace)
		}
		value = "(regex #\"" + value + "\")"
	} else {
		value = "(literal \"" + value + "\")"
	}
	return value
}
