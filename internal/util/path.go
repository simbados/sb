package util

import (
	"errors"
	"path/filepath"
	"regexp"
	"sb/internal/types"
	"strings"
)

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
		"[local]":  paths.LocalConfigPath,
		"[root]":   paths.RootConfigPath,
	}
}

func expandPaths(paths *types.Paths, value string) (string, error) {
	value, err := commonPathExpansion(paths, value)
	if err != nil {
		return "", err
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

func commonPathExpansion(paths *types.Paths, value string) (string, error) {
	initialPath := value
	matching := buildSymbolToPathMatching(paths)
	for key, path := range matching {
		value = strings.ReplaceAll(value, key, path)
	}
	if strings.HasPrefix(value, "../") || strings.HasPrefix(value, "./") {
		return "", errors.New("no relative path allowed in config, use an identifier like [wd] or [home], for more information consult the documentation")
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
	return value, nil
}
