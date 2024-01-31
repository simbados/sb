package types

type CliOptionsStr struct {
	DebugEnabled     bool `json:"debug-enabled"`
	PrintEnabled     bool `json:"print-enabled"`
	DryRunEnabled    bool `json:"dry-run-enabled"`
	CreateExeEnabled bool `json:"create-exe-enabled"`
	HelpEnabled      bool `json:"help"`
	VersionEnabled   bool `json:"version-enabled"`
}

var CliOptions = CliOptionsStr{
	DebugEnabled:     false,
	PrintEnabled:     false,
	DryRunEnabled:    false,
	CreateExeEnabled: false,
}

var ValidCliOptions = map[string]string{
	"--debug":      "--debug",
	"-d":           "-d",
	"--print":      "--print",
	"-p":           "-p",
	"--dry-run":    "--dry-rin",
	"-dr":          "-dr",
	"--create-exe": "--create-exe",
	"-ce":          "-ce",
	"--help":       "--help",
	"-h":           "-h",
	"--version":    "--version",
	"-v":           "-v",
	"--init":       "--init",
	"-i":           "-i",
}
