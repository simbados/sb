package types

type Config struct {
	BinaryPath    string    `json:"binary-path"`
	BinaryDirPath string    `json:"binary-dir-path"`
	BinaryName    string    `json:"binary-name"`
	SbConfig      *SbConfig `json:"root-config"`
	Commands      []string  `json:"commands"`
	CliOptions    []string  `json:"cli-options"`
	CliConfig     *SbConfig `json:"cli-config"`
}
