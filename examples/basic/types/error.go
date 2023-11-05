package types

type Error struct {
	Error   string `json:"error" yaml:"error"`
	Details string `json:"details" yaml:"details"`
}
