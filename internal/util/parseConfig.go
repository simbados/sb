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
	localConfig := extractLocalConfig(context)
	rootConfig := extractRootConfig(context)
	return mergeConfigs(localConfig, rootConfig)
}

func extractLocalConfig(context *types.Context) *types.SbConfig {
	localConfigPath, localConfigExists := LocalConfigPath(&context.Paths, context.Config.BinaryName)
	if localConfigExists {
		context.Paths.LocalConfigPath = localConfigPath
		localConfig := parseJsonConfig(&context.Paths, localConfigPath, context.Config.Commands, 1)
		log.LogDebug("Using local config file at path ", localConfigPath)
		return localConfig
	} else {
		log.LogDebug("No local config file found at: ", localConfigPath)
		log.LogDebug("Proceeding without local config")
	}
	return nil
}

func extractRootConfig(context *types.Context) *types.SbConfig {
	binaryGlobalConfigPath, exists := doesRootConfigExists(context)
	if exists {
		globalConfig := parseJsonConfig(&context.Paths, binaryGlobalConfigPath, context.Config.Commands, 1)
		return globalConfig
	}
	return nil
}

func doesRootConfigExists(context *types.Context) (string, bool) {
	if doesRootConfigDirExist(context.Paths.RootConfigPath) {
		binaryGlobalConfigPath := context.Paths.RootConfigPath + "/" + context.Config.BinaryName + ".json"
		binaryPathExists, _ := DoesPathExist(binaryGlobalConfigPath)
		if !binaryPathExists {
			log.LogDebug("No root config for binary found. You might want to create a config file at: ", context.Paths.RootConfigPath)
		} else {
			log.LogDebug("Using global config file")
			return binaryGlobalConfigPath, true
		}
	} else {
		log.LogDebug("No root config file found at: ", context.Paths.RootConfigPath)
		log.LogDebug("Proceeding without global config")
	}
	return "", false
}

func parseJsonConfig(paths *types.Paths, path string, commands []string, depth int) *types.SbConfig {
	if depth > 3 {
		log.LogErr("Nesting of the json config is only allowed for a depth of 2")
	}
	mapping := buildCommandMap(commands)
	mapping[types.RootConfigKey] = true
	configJson := ParseJson(path)
	configs := parseExtendedConfig(paths, path, commands, depth, configJson)
	for key, val := range configJson {
		var command string
		if strings.Contains(key, "*") {
			command = strings.TrimSuffix(key, "*")
		} else {
			command = key
		}
		if exists := mapping[command]; exists {
			configs = append(configs, parseOptionsForCommand(paths, path, val))
		} else {
			log.LogDev(fmt.Sprintf("No config found for key: %v in path %v", key, path))
		}
	}
	if len(configs) == 0 {
		log.LogErr(fmt.Sprintf("You have a config file at path %v, but no keys were found", path))
	}
	return mergeConfigs(configs...)
}

func parseExtendedConfig(paths *types.Paths, path string, commands []string, depth int, configJson map[string]interface{}) []*types.SbConfig {
	var configs []*types.SbConfig
	if value, exists := configJson[types.ExtendsConfigKey]; exists {
		if extendPath, isString := value.(string); isString {
			exists, err := DoesPathExist(extendPath)
			if err != nil {
				log.LogErr(err)
			}
			if exists {
				log.LogDebug("Extending config with config of path: ", extendPath)
				configs = append(configs, parseJsonConfig(paths, extendPath, commands, depth+1))
			} else {
				log.LogWarn("Path which was provided for extending the config does not exists: ", path)
			}
		} else {
			log.LogWarn(types.ExtendsConfigKey, " key is not a string at path", extendPath)
		}
	}
	return configs
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

func parseOptionsForCommand(paths *types.Paths, path string, val interface{}) *types.SbConfig {
	permissions, valid := parseNextJsonLevel(val)
	if !valid {
		log.LogErr("Malformed root config json, please check your config at path: ", path)
	}
	return parseConfigIntoStruct(paths, permissions, path)
}

func parseNextJsonLevel(config interface{}) (map[string]interface{}, bool) {
	rootConf, ok := config.(map[string]interface{})
	return rootConf, ok
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
		sbConfig.Read = parseIfExists(paths, read, path, "read")
	}
	if writeExists {
		sbConfig.Write = parseIfExists(paths, write, path, "write")
	}
	if readWriteExists {
		sbConfig.ReadWrite = parseIfExists(paths, readWrite, path, "read-write")
	}
	if processExists {
		sbConfig.Process = parseIfExists(paths, process, path, "process")
	}
	if netOutExists {
		if value, valid := netOut.(bool); valid {
			sbConfig.NetworkOutbound = &types.BoolOrNil{Value: value}
		} else {
			log.LogErr(fmt.Sprintf("Your net-out config at path %v, is not a boolean value", path))
		}
	}
	if netInExists {
		if value, valid := netIn.(bool); valid {
			sbConfig.NetworkInbound = &types.BoolOrNil{Value: value}
		} else {
			log.LogErr(fmt.Sprintf("Your net-in config at path %v, is not a boolean value", path))
		}
	}
	return sbConfig
}

func parseIfExists(paths *types.Paths, jsonKey interface{}, path string, configName string) []string {
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

func mergeConfigs(configToMerge ...*types.SbConfig) *types.SbConfig {
	newConfig := &types.SbConfig{}
	for _, config := range configToMerge {
		if config != nil {
			newConfig.Write = appendUniqueStrings(newConfig.Write, config.Write...)
			newConfig.Read = appendUniqueStrings(newConfig.Read, config.Read...)
			newConfig.Process = appendUniqueStrings(newConfig.Process, config.Process...)
			newConfig.ReadWrite = appendUniqueStrings(newConfig.ReadWrite, config.ReadWrite...)
			// Network in/out-bound can not really be merged if it is allowed once it should be allowed
			newConfig.NetworkOutbound = mergeBoolOrNil(newConfig.NetworkOutbound, config.NetworkOutbound)
			newConfig.NetworkInbound = mergeBoolOrNil(newConfig.NetworkInbound, config.NetworkInbound)
		}
	}
	return newConfig
}

func mergeBoolOrNil(a *types.BoolOrNil, b *types.BoolOrNil) *types.BoolOrNil {
	if a == nil && b == nil {
		return nil
	} else if a == nil {
		return b
	} else if b == nil {
		return a
	} else {
		return &types.BoolOrNil{Value: a.Value && b.Value}
	}
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
