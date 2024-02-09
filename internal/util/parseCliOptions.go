package util

import (
	"errors"
	"path/filepath"
	"regexp"
	"sb/internal/log"
	"sb/internal/types"
	"strings"
)

func OptionsParsing(paths *types.Paths, args []string) ([]string, *types.SbConfig, []string) {
	var options []string
	var cliConfigSb types.SbConfig
	var cliConfig *types.SbConfig = nil
	optionsUntilIndex := 0
	for index, value := range args {
		if strings.HasPrefix(value, "-") {
			split, splitValue := parseCliConfigParam(value)
			if _, exist := types.ValidCliOptions[split]; exist {
				options = append(options, value)
			} else if _, configExists := types.AllowedConfigKeys[split]; configExists && len(splitValue) > 0 {
				cliConfig = addToCliConfig(paths, cliConfig, cliConfigSb, splitValue, split)
			} else {
				log.LogErr("You passed a wrong cli option: ", value)
			}
		} else {
			optionsUntilIndex = index
			break
		}
	}
	var commands []string
	if len(args) != len(options) {
		commands = args[optionsUntilIndex:]
	}
	if cliConfig != nil {
		return options, cliConfig, commands
	} else {
		return options, nil, commands
	}
}

func addToCliConfig(paths *types.Paths, cliConfig *types.SbConfig, cliConfigSb types.SbConfig, splitValue string, split string) *types.SbConfig {
	cliConfig = &cliConfigSb
	if splitValue == "true" || splitValue == "false" {
		addToConfig(cliConfig, split, splitValue)
	} else if arr := strings.Split(splitValue, ","); len(arr) > 0 {
		for _, val := range arr {
			if path, err := expandPaths(paths, val); err == nil {
				addToConfig(cliConfig, split, path)
			} else {
				log.LogErr(err)
			}
		}
	}
	return cliConfig
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
			log.LogErr("You must provide true or false value for cli config: ", value)
		}
		config.NetworkInbound = &types.BoolOrNil{Value: boolVal}
		break
	case "--net-out":
		boolVal, exists := parseStringBoolean(value)
		if !exists {
			log.LogErr("You must provide true or false value for cli config: ", value)
		}
		config.NetworkOutbound = &types.BoolOrNil{Value: boolVal}
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

const globSingleReplace = `[^\/]*`

func buildSymbolToPathMatching(paths *types.Paths) map[string]string {
	return map[string]string{
		"~":        paths.HomePath,
		"[home]":   paths.HomePath,
		"[wd]":     paths.WorkingDir,
		"[bin]":    paths.BinPath,
		"[target]": paths.BinaryPath,
	}
}

func expandPaths(paths *types.Paths, value string) (string, error) {
	initialPath := value
	matching := buildSymbolToPathMatching(paths)
	for key, path := range matching {
		value = strings.ReplaceAll(value, key, path)
	}
	if strings.HasPrefix(value, "../") {
		value = paths.HomePath + "/" + value
	}
	if strings.Contains(value, "../") {
		splits := strings.Split(value, "/")[1:]
		for index := 0; index < len(splits); index++ {
			if splits[index] == ".." {
				if index+1 > len(splits) || index-1 < 0 {
					return "", errors.New("Can not resolve path of: " + initialPath + " in one of your config files or cli args")
				}
				splits = append(splits[0:index-1], splits[index+1:]...)
				index = -1
			}
		}
		value = "/" + filepath.Join(splits...)
	}
	if strings.HasPrefix(value, "./") {
		value = strings.ReplaceAll(value, "./", paths.WorkingDir+"/")
	}
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
	return value, nil
}
