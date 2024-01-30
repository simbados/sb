package parse

import (
	"sb/internal/types"
	"testing"
)

func TestExpandPath(t *testing.T) {
	tests := []struct {
		a        string
		expected string
	}{
		{"~/whatsup", "(literal \"/Users/test/whatsup\")"},
		{"[home]/whatsup", "(literal \"/Users/test/whatsup\")"},
		{"[bin]/whatsup", "(literal \"/usr/bin/whatsup\")"},
		{"[wd]/whatsup", "(literal \"/Users/test/sb/whatsup\")"},
		{"./whatsup", "(literal \"/Users/test/sb/whatsup\")"},
		{"/Users/test/**", "(regex #\"/Users/test/(.+)\")"},
		{"/Users/**/test", "(regex #\"/Users/(.+\\/)?test\")"},
		{"/Users/test/*.js", "(regex #\"/Users/test/^[^\\/]*\\.js\")"},
	}

	paths := types.Paths{LocalConfigPath: "Users/test/sb", HomePath: "/Users/test", RootConfigPath: "Users/test", WorkingDir: "/Users/test/sb", BinPath: "/usr/bin", BinaryPath: "/usr/bin/ls"}
	for _, tt := range tests {
		res := expandPaths(&paths, tt.a)
		if res != tt.expected {
			t.Errorf("expandPaths(paths, %v) = %v; but got %v", tt.a, tt.expected, res)
		}
	}
}
