package types

type CliOptionsStr struct {
	DebugEnabled        bool   `json:"debug-enabled"`
	DryRunEnabled       bool   `json:"dry-run-enabled"`
	CreateExeEnabled    bool   `json:"create-exe-enabled"`
	ShowConfigEnabled   bool   `json:"show-config-enabled"`
	HelpEnabled         bool   `json:"help"`
	VersionEnabled      bool   `json:"version-enabled"`
	EditEnabled         bool   `json:"edit-enabled"`
	VigilantModeEnabled bool   `json:"vigilant_mode_enabled"`
	ConfigModeEnabled   string `json:"config-mode-enabled"`
}

type EnvsStr struct {
	DevModeEnabled bool `json:"dev_mode_enabled"`
}

var CliOptions = CliOptionsStr{}

var Envs = EnvsStr{}

// Mapping from cli option name to how many arguments are expected
var ValidCliOptions = map[string]int{
	"-c":         1,
	"--config":   1,
	"--debug":    1,
	"-d":         1,
	"--dry-run":  1,
	"-dr":        1,
	"--help":     1,
	"-h":         1,
	"--version":  1,
	"-v":         1,
	"--init":     1,
	"-i":         1,
	"--edit":     1,
	"-e":         1,
	"-s":         1,
	"--vigilant": 1,
	"-vi":        1,
}
