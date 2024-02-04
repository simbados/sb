package types

type CliOptionsStr struct {
	DebugEnabled     bool `json:"debug-enabled"`
	DryRunEnabled    bool `json:"dry-run-enabled"`
	CreateExeEnabled bool `json:"create-exe-enabled"`
	HelpEnabled      bool `json:"help"`
	VersionEnabled   bool `json:"version-enabled"`
	EditEnabled      bool `json:"edit-enabled"`
}

var CliOptions = CliOptionsStr{
	DebugEnabled:     false,
	DryRunEnabled:    false,
	CreateExeEnabled: false,
}

var ValidCliOptions = map[string]string{
	"--debug":      "--debug",
	"-d":           "-d",
	"--dry-run":    "--dry-run",
	"-dr":          "-dr",
	"--create-exe": "--create-exe",
	"-ce":          "-ce",
	"--help":       "--help",
	"-h":           "-h",
	"--version":    "--version",
	"-v":           "-v",
	"--init":       "--init",
	"-i":           "-i",
	"--edit":       "--edit",
	"-e":           "-e",
}
