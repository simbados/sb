package main

import (
	"fmt"
	"regexp"
	"strings"
)

func parseOptions(paths *Paths, args []string) ([]string, *SbConfig, []string) {
	var options []string
	var cliConfig SbConfig
	re := regexp.MustCompile("^-.*")
	optionsUntilIndex := 0
	hasCliConfig := false
	for index, value := range args {
		fmt.Println("value", value)
		if re.MatchString(value) {
			split, splitValue := parseCliConfigParam(value)
			if _, exist := validCliOptions[split]; exist {
				options = append(options, value)
				setOption(value)
			} else if _, configExists := AllowedConfigKeys[split]; configExists && len(splitValue) > 0 {
				addToConfig(&cliConfig, split, expandPaths(paths, splitValue))
				hasCliConfig = true
			} else {
				logErr("You passed a wrong cli option: ", value)
			}
		} else {
			optionsUntilIndex = index
			break
		}
	}
	if len(options) == len(args) {
		logErr("Please specify the program that you want to sandbox")
	}
	if hasCliConfig {
		return options, &cliConfig, args[optionsUntilIndex:]
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

func addToConfig(config *SbConfig, key string, value string) *SbConfig {
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
			logErr("You must provide true or false value for cli config: ", value)
		}
		config.NetworkInbound = boolVal
		break
	case "--net-out":
		boolVal, exists := parseStringBoolean(value)
		if !exists {
			logErr("You must provide true or false value for cli config: ", value)
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
const globEndReplace = `(*.)`
const globMiddleRegex = `\*\*\/`
const globMiddleReplace = `(.+\/)?`
const globSingleRegex = `\*`
const globSingleReplace = `([^\/]+)?`

func expandPaths(paths *Paths, value string) string {
	fmt.Println(value)
	value = strings.ReplaceAll(value, "~", paths.HomePath)
	value = strings.ReplaceAll(value, "[home]", paths.HomePath)
	value = strings.ReplaceAll(value, "./", paths.WorkingDir+"/")
	value = strings.ReplaceAll(value, "[wd]", paths.WorkingDir)
	value = strings.ReplaceAll(value, "[bin]", paths.BinPath)
	if strings.Contains(value, "*") {
		value = regexp.MustCompile(dotRegex).ReplaceAllString(value, dotReplace)
		value = regexp.MustCompile(globEndRegex).ReplaceAllString(value, globEndReplace)
		value = regexp.MustCompile(globMiddleRegex).ReplaceAllString(value, globMiddleReplace)
		value = regexp.MustCompile(globSingleRegex).ReplaceAllString(value, globSingleReplace)
		fmt.Println(value)
		value = "(regex #\"" + value + "\")"
	} else {
		value = "(literal \"" + value + "\")"
	}
	return value
}
