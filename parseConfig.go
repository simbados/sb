package main

import (
	"fmt"
	"os"
	"strings"
)

func doesRootConfigExist(path string) bool {
	isRootConfigExisting, err := doesPathExist(path)
	if !isRootConfigExisting {
		logWarn("Root config directory does not exist", err)
		// TODO: this could be on first run we might want to make sure it exists at this point
		// Built first run logic to create default files and directories if non existent
		return false
	}
	return true
}

func parseConfigFiles(context *Context) {
	globalConfig := &SbConfig{}
	commandConfig := &SbConfig{}
	if doesRootConfigExist(context.Paths.RootConfigPath) {
		binaryGlobalConfigPath := context.Paths.RootConfigPath + "/" + context.Config.BinaryName + ".json"
		binaryPathExists, _ := doesPathExist(binaryGlobalConfigPath)
		if !binaryPathExists {
			logWarn("No config for binary found. You might want to create a config file at: ", context.Paths.RootConfigPath)
		} else {
			globalConfig = parseRootBinaryConfig(&context.Paths, binaryGlobalConfigPath)
			if len(context.Config.Commands) > 0 {
				commandConfig = parseRootCommandConfig(&context.Paths, binaryGlobalConfigPath, context.Config.Commands[0])
			}
		}
	}
	localConfigPath := context.Paths.LocalConfigPath
	localConfigExists, _ := doesPathExist(localConfigPath)
	localConfig := &SbConfig{}
	if localConfigExists {
		localConfig = parseLocalBinaryConfig(&context.Paths, localConfigPath, context.Config.BinaryName)
	}
	// Merge all configs together
	context.Config.SbConfig = mergeConfig(globalConfig, commandConfig, localConfig, context.Config.CliConfig)
}

func parseLocalBinaryConfig(paths *Paths, path string, binaryName string) *SbConfig {
	configJson := parseJson(path)
	emptyConf := &SbConfig{}

	binaryConfig := configJson[binaryName]
	if binaryConfig == nil {
		logWarn("Local config found but no local config for binary")
		return emptyConf
	}
	rootConfLocal, err := parseNextJsonLevel(binaryConfig)
	if err {
		logWarn("Error while parsing local config for binary: ", binaryName)
		return emptyConf
	}

	rootConf := rootConfLocal["__root-conf__"]
	if rootConf == nil {
		logWarn("Local config found for binary but no root config for binary:", binaryName)
		return emptyConf
	}

	rootConfMap, errRootConfMap := parseNextJsonLevel(rootConf)
	if errRootConfMap {
		logWarn("Error while parsing local config for binary: ", binaryName)
		return emptyConf
	}
	return parseConfigIntoStruct(paths, rootConfMap, path)
}

func parseRootBinaryConfig(paths *Paths, path string) *SbConfig {
	configJson := parseJson(path)
	rootConf := configJson["__root-conf__"]
	if rootConf == nil {
		fmt.Printf("Malformed config json, please check your config at path: %s\n", path)
		os.Exit(1)
	}
	permissions, err := parseNextJsonLevel(rootConf)
	if err {
		fmt.Printf("Malformed root config json, please check your config at path: %s\n", path)
		os.Exit(1)
	}
	return parseConfigIntoStruct(paths, permissions, path)
}

func parseRootCommandConfig(paths *Paths, path string, command string) *SbConfig {
	configJson := parseJson(path)
	conf := configJson[command]
	if conf == nil {
		logWarn("No root config for argument found: ", command)
		return &SbConfig{}
	}
	permissions, err := parseNextJsonLevel(conf)
	if err {
		logErr("Malformed root config json, please check your config at path: ", path)
		os.Exit(1)
	}
	return parseConfigIntoStruct(paths, permissions, path)
}

func parseNextJsonLevel(config interface{}) (map[string]interface{}, bool) {
	rootConf, ok := config.(map[string]interface{})
	return rootConf, !ok
}

func parseConfigIntoStruct(paths *Paths, binaryConfig map[string]interface{}, path string) *SbConfig {
	sbConfig := &SbConfig{}
	for key := range binaryConfig {
		_, exists := AllowedConfigKeys[key]
		if !exists {
			fmt.Printf("Found unsupported key in your binary config, please remove this key: %s at path %s\n", key, path)
			os.Exit(1)
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

func convertJsonArrayToStringArray(paths *Paths, jsonArray []interface{}) []string {
	var valueStrings []string
	for _, value := range jsonArray {
		if str, ok := value.(string); ok {
			valueStrings = append(valueStrings, expandPaths(paths, strings.Trim(str, " ")))
		} else {
			logErr("Malformed value in config for following array and value: ", jsonArray, value)
		}
	}
	return valueStrings
}

func mergeConfig(configToMerge ...*SbConfig) *SbConfig {
	newConfig := &SbConfig{}
	for _, config := range configToMerge {
		if config != nil {
			newConfig.Write = appendUniqueStrings(newConfig.Write, config.Write...)
			newConfig.Read = appendUniqueStrings(newConfig.Read, config.Read...)
			newConfig.Process = appendUniqueStrings(newConfig.Process, config.Process...)
			newConfig.ReadWrite = appendUniqueStrings(newConfig.ReadWrite, config.ReadWrite...)
			// Network in/out-bound can not really be merged if it is prohibited once it should be enforced
			fmt.Println(newConfig.NetworkInbound, config.NetworkInbound)
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
