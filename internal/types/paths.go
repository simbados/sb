package types

type Paths struct {
	HomePath        string `json:"home-path"`
	RootConfigPath  string `json:"root-config-path"`
	WorkingDir      string `json:"working-dir"`
	LocalConfigPath string `json:"local-config-path"`
	BinPath         string `json:"bin-path"`
	BinaryPath      string `json:"binary-path"`
	SbBinaryPath    string `json:"sb-binary-path"`
}
