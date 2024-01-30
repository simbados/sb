package parse

import (
	"fmt"
	"os"
	"sb/internal/types"
	"sb/internal/util"
	"strings"
)

func doesRootConfigDirExist(path string) bool {
	isRootConfigExisting, err := util.DoesPathExist(path)
	fmt.Println("is existing", isRootConfigExisting, path)
	if !isRootConfigExisting {
		util.LogWarn("Root config directory does not exist", err)
		// TODO: this could be on first run we might want to make sure it exists at this point
		// Built first run logic to create default files and directories if non existent
		return false
	}
	return true
}

// Here we only parse the config files, cli configs have already been parsed
func ParseConfigFiles(context *types.Context) *types.SbConfig {
	globalConfig := &types.SbConfig{}
	localConfigPath := context.Paths.LocalConfigPath + "/" + context.Config.BinaryName + ".json"
	localConfigExists, _ := util.DoesPathExist(localConfigPath)
	localConfig := &types.SbConfig{}
	if localConfigExists {
		localConfig = parseRootBinaryConfig(&context.Paths, localConfigPath, context.Config.Commands)
		util.LogInfo("Using local config file")
		return localConfig
	} else {
		util.LogInfo("No local config file found at: ", localConfigPath)
		util.LogInfo("Proceeding without local config")
	}
	if doesRootConfigDirExist(context.Paths.RootConfigPath) {
		binaryGlobalConfigPath := context.Paths.RootConfigPath + "/" + context.Config.BinaryName + ".json"
		binaryPathExists, _ := util.DoesPathExist(binaryGlobalConfigPath)
		if !binaryPathExists {
			util.LogWarn("No config for binary found. You might want to create a config file at: ", context.Paths.RootConfigPath)
		} else {
			globalConfig = parseRootBinaryConfig(&context.Paths, binaryGlobalConfigPath, context.Config.Commands)
			util.LogInfo("Using global config file")
			return globalConfig
		}
	} else {
		util.LogInfo("No root config file found at: ", context.Paths.RootConfigPath)
		util.LogInfo("Proceeding without global config")
	}
	return nil
}

func parseRootBinaryConfig(paths *types.Paths, path string, commands []string) *types.SbConfig {
	commands = append(commands, "__root-config__")
	var configs []*types.SbConfig
	configJson := util.ParseJson(path)
	for key, val := range configJson {
		for _, command := range commands {
			fmt.Println(strings.Contains(command, key))
			if strings.Contains(command, key) {
				if val == nil {
					util.LogInfo(fmt.Sprintf("No config found for key: %v in path %v", key, path))
				} else {
					permissions, err := parseNextJsonLevel(val)
					if err {
						util.LogErr("Malformed root config json, please check your config at path: ", path)
					}
					configs = append(configs, parseConfigIntoStruct(paths, permissions, path))
				}
			}
		}
	}
	if len(configs) == 0 {
		util.LogErr(fmt.Sprintf("You have a config file at path %v, but no keys were found", path))
	}
	fmt.Println(configs[0].Read)
	return mergeConfig(configs...)
}

//func parseRootCommandConfig(paths *Paths, path string, command string) *SbConfig {
//	configJson := parseJson(path)
//	conf := configJson[command]
//	if conf == nil {
//		logWarn("No root config for argument found: ", command)
//		return &SbConfig{}
//	}
//	permissions, err := parseNextJsonLevel(conf)
//	if err {
//		logErr("Malformed root config json, please check your config at path: ", path)
//	}
//	return parseConfigIntoStruct(paths, permissions, path)
//}

func parseNextJsonLevel(config interface{}) (map[string]interface{}, bool) {
	rootConf, ok := config.(map[string]interface{})
	return rootConf, !ok
}

func parseConfigIntoStruct(paths *types.Paths, binaryConfig map[string]interface{}, path string) *types.SbConfig {
	sbConfig := &types.SbConfig{}
	for key := range binaryConfig {
		_, exists := types.AllowedConfigKeys[key]
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
	fmt.Println(read, readExists, binaryConfig["read"])
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
			valueStrings = append(valueStrings, expandPaths(paths, strings.Trim(str, " ")))
		} else {
			util.LogErr("Malformed value in config for following array and value: ", jsonArray, value)
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
			newConfig.NetworkOutbound = newConfig.NetworkOutbound && config.NetworkOutbound
			newConfig.NetworkInbound = newConfig.NetworkInbound && config.NetworkInbound
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
