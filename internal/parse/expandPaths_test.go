package parse

import (
	"sb/internal/types"
	"testing"
)

func TestExpandDotsError(t *testing.T) {
	tests := []struct {
		a        string
		expected string
	}{
		{"[wd]/../../../../hello", "(literal \"/Users/test/hello\")"},
		{"../../../../hello", "(literal \"/Users/test/hello\")"},
		{"what/is/this../../../../", "(literal \"/Users/test/hello\")"},
	}

	paths := types.Paths{LocalConfigPath: "Users/test/sb", HomePath: "/Users/test", RootConfigPath: "Users/test", WorkingDir: "/Users/test/sb", BinPath: "/usr/bin", BinaryPath: "/usr/bin/ls"}
	for _, tt := range tests {
		_, err := expandPaths(&paths, tt.a)
		if err == nil {
			t.Errorf("expandPaths(paths, %v) = %v; but got return error", tt.a, tt.expected)
		}
	}
}

func TestExpandPathSuccess(t *testing.T) {
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
		{"[wd]/../hello", "(literal \"/Users/test/hello\")"},
		{"[wd]/../hello/../what", "(literal \"/Users/test/what\")"},
		{"[wd]/../hello/../what*", "(regex #\"/Users/test/what^[^\\/]*\")"},
	}

	paths := types.Paths{LocalConfigPath: "Users/test/sb", HomePath: "/Users/test", RootConfigPath: "Users/test", WorkingDir: "/Users/test/sb", BinPath: "/usr/bin", BinaryPath: "/usr/bin/ls"}
	for _, tt := range tests {
		res, err := expandPaths(&paths, tt.a)
		if res != tt.expected {
			t.Errorf("expandPaths(paths, %v) = %v; but got %v", tt.a, tt.expected, res)
		}
		if err != nil {
			t.Errorf("expandPaths(paths, %v) = %v; but got error %v", tt.a, tt.expected, err)
		}
	}
}
