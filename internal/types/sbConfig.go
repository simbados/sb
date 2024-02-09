package types

type BoolOrNil struct {
	Value bool
}

type SbConfig struct {
	Read            []string   `json:"read"`
	Write           []string   `json:"write"`
	ReadWrite       []string   `json:"read-write"`
	Process         []string   `json:"process"`
	NetworkOutbound *BoolOrNil `json:"net-out"`
	NetworkInbound  *BoolOrNil `json:"net-in"`
}

var AllowedConfigKeys = map[string]string{
	"process":    "process",
	"write":      "write",
	"read":       "read",
	"read-write": "read-write",
	"net-in":     "net-in",
	"net-out":    "net-out",
	"--process":  "--process",
	"--write":    "--write",
	"--read":     "--read",
	"--net-in":   "--net-in",
	"--net-out":  "--net-out",
}
