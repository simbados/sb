package util

import (
	"fmt"
	"sb/internal/log"
	"sb/internal/types"
	"slices"
	"strings"
)

func doesRootConfigDirExist(path string) bool {
	isRootConfigExisting, err := DoesPathExist(path)
	if !isRootConfigExisting {
		log.LogWarn("Root config directory does not exist", err)
		// TODO: this could be on first run we might want to make sure it exists at this point
		// Built first run logic to create default files and directories if non existent
		return false
	}
	return true
}

func GetSubdirectories(path string, homePath string) []string {
	var allPaths []string
	currentPath := ""
	paths := strings.Split(path, "/")
	for _, singlePath := range paths[1:] {
		currentPath += "/" + singlePath
		// We just go as high as the home directory
		if !strings.Contains(homePath, currentPath) {
			allPaths = append(allPaths, currentPath)
		}
	}
	slices.Reverse(allPaths)
	return allPaths
}

func LocalConfigPath(paths *types.Paths, binaryName string) (string, bool) {
	allLocalPath := GetSubdirectories(paths.WorkingDir, paths.HomePath)
	var localConfigExists bool
	var localConfigPath string
	for _, localPath := range allLocalPath {
		localConfigPath = localPath + types.LocalConfigPath + "/" + binaryName + ".json"
		localConfigExists, _ = DoesPathExist(localConfigPath)
		if localConfigExists {
			break
		}
	}
	return localConfigPath, localConfigExists
}

// ConfigFileParsing Here we only parse the config files, cli configs have already been parsed
func ConfigFileParsing(context *types.Context) *types.SbConfig {
	globalConfig := &types.SbConfig{}
	sbConfig, localConfigExists := extractLocalConfig(context)
	if localConfigExists {
		return sbConfig
	}
	config, rootConfigExists := extractRootConfig(context, globalConfig)
	if rootConfigExists {
		return config
	}
	return nil
}

func extractLocalConfig(context *types.Context) (*types.SbConfig, bool) {
	localConfigPath, localConfigExists := LocalConfigPath(&context.Paths, context.Config.BinaryName)
	if localConfigExists {
		context.Paths.LocalConfigPath = localConfigPath
		localConfig := parseJsonConfig(&context.Paths, localConfigPath, context.Config.Commands)
		log.LogDebug("Using local config file at path ", localConfigPath)
		return localConfig, true
	} else {
		log.LogDebug("No local config file found at: ", localConfigPath)
		log.LogDebug("Proceeding without local config")
	}
	return nil, false
}

func extractRootConfig(context *types.Context, globalConfig *types.SbConfig) (*types.SbConfig, bool) {
	if doesRootConfigDirExist(context.Paths.RootConfigPath) {
		binaryGlobalConfigPath := context.Paths.RootConfigPath + "/" + context.Config.BinaryName + ".json"
		binaryPathExists, _ := DoesPathExist(binaryGlobalConfigPath)
		if !binaryPathExists {
			log.LogWarn("No config for binary found. You might want to create a config file at: ", context.Paths.RootConfigPath)
		} else {
			globalConfig = parseJsonConfig(&context.Paths, binaryGlobalConfigPath, context.Config.Commands)
			log.LogDebug("Using global config file")
			return globalConfig, true
		}
	} else {
		log.LogDebug("No root config file found at: ", context.Paths.RootConfigPath)
		log.LogDebug("Proceeding without global config")
	}
	return nil, false
}

func parseJsonConfig(paths *types.Paths, path string, commands []string) *types.SbConfig {
	mapping := buildCommandMap(commands)
	mapping["__root-config__"] = true
	var configs []*types.SbConfig
	configJson := ParseJson(path)
	for key, val := range configJson {
		var command string
		if strings.Contains(key, "*") {
			command = strings.TrimSuffix(key, "*")
		} else {
			command = key
		}
		if exists := mapping[command]; exists {
			configs = parseOptionsForCommand(paths, path, val, configs)
		} else {
			log.LogDev(fmt.Sprintf("No config found for key: %v in path %v", key, path))
		}
	}
	if len(configs) == 0 {
		log.LogErr(fmt.Sprintf("You have a config file at path %v, but no keys were found", path))
	}
	return mergeConfig(configs...)
}

func buildCommandMap(commands []string) map[string]bool {
	aggregate := ""
	mapping := map[string]bool{}
	for _, command := range commands {
		if aggregate != "" {
			aggregate = aggregate + " " + command
			mapping[aggregate] = true
		} else {
			aggregate = command
		}
		mapping[command] = true
	}
	return mapping
}

func parseOptionsForCommand(paths *types.Paths, path string, val interface{}, configs []*types.SbConfig) []*types.SbConfig {
	permissions, err := parseNextJsonLevel(val)
	if err {
		log.LogErr("Malformed root config json, please check your config at path: ", path)
	}
	configs = append(configs, parseConfigIntoStruct(paths, permissions, path))
	return configs
}

func parseNextJsonLevel(config interface{}) (map[string]interface{}, bool) {
	rootConf, ok := config.(map[string]interface{})
	return rootConf, !ok
}

func parseConfigIntoStruct(paths *types.Paths, binaryConfig map[string]interface{}, path string) *types.SbConfig {
	sbConfig := &types.SbConfig{}
	for key := range binaryConfig {
		_, exists := types.AllowedConfigKeys[key]
		if !exists {
			log.LogErr(fmt.Sprintf("Found unsupported key in your binary config, please remove this key: %s at path %s\n", key, path))
		}
	}
	read, readExists := binaryConfig["read"]
	write, writeExists := binaryConfig["write"]
	readWrite, readWriteExists := binaryConfig["read-write"]
	process, processExists := binaryConfig["process"]
	netOut, netOutExists := binaryConfig["net-out"]
	netIn, netInExists := binaryConfig["net-in"]
	if readExists {
		sbConfig.Read = parseIfExists(paths, read, sbConfig, path, "read")
	}
	if writeExists {
		sbConfig.Write = parseIfExists(paths, write, sbConfig, path, "write")
	}
	if readWriteExists {
		sbConfig.ReadWrite = parseIfExists(paths, readWrite, sbConfig, path, "read-write")
	}
	if processExists {
		sbConfig.Process = parseIfExists(paths, process, sbConfig, path, "process")
	}
	if netOutExists {
		if value, exists := netOut.(bool); exists {
			sbConfig.NetworkOutbound = value
		} else {
			log.LogErr(fmt.Sprintf("Your net-out config at path %v, is not a boolean value", path))
		}
	}
	if netInExists {
		if value, exists := netIn.(bool); exists {
			sbConfig.NetworkInbound = value
		} else {
			log.LogErr(fmt.Sprintf("Your net-in config at path %v, is not a boolean value", path))
		}
	}
	return sbConfig
}

func parseIfExists(paths *types.Paths, jsonKey interface{}, sbConfig *types.SbConfig, path string, configName string) []string {
	if arr, exists := jsonKey.([]interface{}); exists {
		return convertJsonArrayToStringArray(paths, arr)
	} else {
		log.LogErr(fmt.Sprintf("Your %v config at path %v, contains a value which should be an array but is not.", configName, path))
	}
	return []string{}
}

func convertJsonArrayToStringArray(paths *types.Paths, jsonArray []interface{}) []string {
	var valueStrings []string
	for _, value := range jsonArray {
		if str, ok := value.(string); ok {
			if path, err := expandPaths(paths, strings.Trim(str, " ")); err == nil {
				valueStrings = append(valueStrings, path)
			} else {
				log.LogErr(err)
			}
		} else {
			log.LogErr("Malformed value in config for following array and value: ", jsonArray, value)
		}
	}
	return valueStrings
}

func mergeConfig(configToMerge ...*types.SbConfig) *types.SbConfig {
	newConfig := &types.SbConfig{}
	for _, config := range configToMerge {
		if config != nil {
			newConfig.Write = appendUniqueStrings(newConfig.Write, config.Write...)
			newConfig.Read = appendUniqueStrings(newConfig.Read, config.Read...)
			newConfig.Process = appendUniqueStrings(newConfig.Process, config.Process...)
			newConfig.ReadWrite = appendUniqueStrings(newConfig.ReadWrite, config.ReadWrite...)
			// Network in/out-bound can not really be merged if it is prohibited once it should be enforced
			newConfig.NetworkOutbound = newConfig.NetworkOutbound || config.NetworkOutbound
			newConfig.NetworkInbound = newConfig.NetworkInbound || config.NetworkInbound
		}
	}
	return newConfig
}

func appendUniqueStrings(array []string, stringsToMerge ...string) []string {
	unique := make(map[string]struct{})
	for _, s := range array {
		unique[s] = struct{}{}
	}
	for _, s := range stringsToMerge {
		unique[s] = struct{}{}
	}
	result := make([]string, 0, len(unique))
	for s := range unique {
		result = append(result, s)
	}
	return result
}
