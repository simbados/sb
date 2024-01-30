package types

type Context struct {
	Config Config `json:"config"`
	Paths  Paths  `json:"path"`
}
