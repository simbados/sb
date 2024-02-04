package parse

import (
	"fmt"
	"sb/internal/log"
	"sb/internal/types"
	"sb/internal/util"
	"slices"
	"strings"
)

func doesRootConfigDirExist(path string) bool {
	isRootConfigExisting, err := util.DoesPathExist(path)
	if !isRootConfigExisting {
		log.LogWarn("Root config directory does not exist", err)
		// TODO: this could be on first run we might want to make sure it exists at this point
		// Built first run logic to create default files and directories if non existent
		return false
	}
	return true
}

func getSubdirectories(path string, homePath string) []string {
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

// ConfigFileParsing Here we only parse the config files, cli configs have already been parsed
func ConfigFileParsing(context *types.Context) *types.SbConfig {
	globalConfig := &types.SbConfig{}
	allLocalPath := getSubdirectories(context.Paths.WorkingDir, context.Paths.HomePath)
	var localConfigExists bool
	var localConfigPath string
	for _, localPath := range allLocalPath {
		localConfigPath = localPath + types.LocalConfigPath + "/" + context.Config.BinaryName + ".json"
		localConfigExists, _ = util.DoesPathExist(localConfigPath)
		if localConfigExists {
			break
		}
	}
	if localConfigExists {
		localConfig := parseRootBinaryConfig(&context.Paths, localConfigPath, context.Config.Commands)
		log.LogDebug("Using local config file")
		return localConfig
	} else {
		log.LogDebug("No local config file found at: ", localConfigPath)
		log.LogDebug("Proceeding without local config")
	}
	if doesRootConfigDirExist(context.Paths.RootConfigPath) {
		binaryGlobalConfigPath := context.Paths.RootConfigPath + "/" + context.Config.BinaryName + ".json"
		binaryPathExists, _ := util.DoesPathExist(binaryGlobalConfigPath)
		if !binaryPathExists {
			log.LogWarn("No config for binary found. You might want to create a config file at: ", context.Paths.RootConfigPath)
		} else {
			globalConfig = parseRootBinaryConfig(&context.Paths, binaryGlobalConfigPath, context.Config.Commands)
			log.LogDebug("Using global config file")
			return globalConfig
		}
	} else {
		log.LogDebug("No root config file found at: ", context.Paths.RootConfigPath)
		log.LogDebug("Proceeding without global config")
	}
	return nil
}

func parseRootBinaryConfig(paths *types.Paths, path string, commands []string) *types.SbConfig {
	commands = append(commands, "__root-config__")
	var configs []*types.SbConfig
	configJson := util.ParseJson(path)
	for key, val := range configJson {
		for _, command := range commands {
			if strings.Contains(command, key) {
				if val == nil {
					log.LogDebug(fmt.Sprintf("No config found for key: %v in path %v", key, path))
				} else {
					permissions, err := parseNextJsonLevel(val)
					if err {
						log.LogErr("Malformed root config json, please check your config at path: ", path)
					}
					configs = append(configs, parseConfigIntoStruct(paths, permissions, path))
					log.PrettyJson(configs)
				}
			}
		}
	}
	if len(configs) == 0 {
		log.LogErr(fmt.Sprintf("You have a config file at path %v, but no keys were found", path))
	}
	return mergeConfig(configs...)
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
	read, readExists := binaryConfig["read"].([]interface{})
	write, writeExists := binaryConfig["write"].([]interface{})
	readWrite, readWriteExists := binaryConfig["read-write"].([]interface{})
	process, processExists := binaryConfig["process"].([]interface{})
	netOut, netOutExists := binaryConfig["net-out"].(bool)
	netIn, netInExists := binaryConfig["net-in"].(bool)
	if readExists {
		sbConfig.Read = convertJsonArrayToStringArray(paths, read)
	}
	if writeExists {
		sbConfig.Write = convertJsonArrayToStringArray(paths, write)
	}
	if readWriteExists {
		sbConfig.ReadWrite = convertJsonArrayToStringArray(paths, readWrite)
	}
	if processExists {
		sbConfig.Process = convertJsonArrayToStringArray(paths, process)
	}
	if netOutExists {
		sbConfig.NetworkOutbound = netOut
	}
	if netInExists {
		sbConfig.NetworkInbound = netIn
	}
	return sbConfig
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
