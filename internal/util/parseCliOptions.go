package util

import (
	"sb/internal/log"
	"sb/internal/types"
	"strings"
)

func OptionsParsing(paths *types.Paths, args []string) (map[string][]string, *types.SbConfig, []string) {
	var options = map[string][]string{}
	var cliConfigSb types.SbConfig
	var cliConfig *types.SbConfig = nil
	index := 0
	for index < len(args) {
		value := args[index]
		if strings.HasPrefix(value, "-") {
			split, splitValue := splitOptionIfPossible(value)
			if multiOptVal, multiOptionExists := types.ValidCliOptions[split]; multiOptionExists {
				index, options[split] = parseCliOptions(split, splitValue, args, multiOptVal, index)
			} else if _, configExists := types.AllowedConfigKeys[split]; configExists && len(splitValue) > 0 {
				cliConfig = addToCliConfig(paths, cliConfig, cliConfigSb, splitValue, split)
				index++
			} else {
				log.LogErr("You passed a wrong cli option: ", value)
			}
		} else {
			break
		}
	}
	var commands []string
	if len(args) != len(options) {
		commands = args[index:]
	}
	return options, cliConfig, commands
}

func parseCliOptions(split string, splitValue string, args []string, multiOptVal int, index int) (int, []string) {
	if splitValue == "" {
		if len(args) < multiOptVal+index {
			log.LogErr("You passed an option which requires additional arguments", split)
		}
		return index + multiOptVal, args[index : multiOptVal+index]
	} else {
		return index + 1, []string{split, splitValue}
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

func splitOptionIfPossible(s string) (string, string) {
	split := strings.SplitN(s, "=", 2)
	if len(split) == 2 {
		return split[0], split[1]
	}
	return split[0], ""
}
