package parse

import (
	"sb/internal/types"
	"slices"
	"strings"
	"testing"
)

func TestRootConfig(t *testing.T) {
	var commands = []string{"run", "some"}
	paths := types.Paths{LocalConfigPath: "Users/test/sb", HomePath: "/Users/test", RootConfigPath: "Users/test", WorkingDir: "/Users/test/sb", BinPath: "/usr/bin", BinaryPath: "/usr/bin/ls"}
	config := parseRootBinaryConfig(&paths, "./test.json", commands)
	if len(config.Write) != 2 {
		t.Errorf("parseRootBinaryConfig should have 2 entries for write but was %v", config.Write)
	}
	if !slices.ContainsFunc(config.Write, func(s string) bool {
		return strings.Contains(s, "yes")
	}) {
		t.Errorf("parseRootBinaryConfig config entry should have 'yes' in path but was %v", config.Write)
	}
	if !slices.ContainsFunc(config.Write, func(s string) bool {
		return strings.Contains(s, "okay")
	}) {
		t.Errorf("parseRootBinaryConfig config entry should have 'okay' in path but was %v", config.Write)
	}
}
