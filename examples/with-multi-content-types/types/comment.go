package types

type Comment struct {
	ID     int    `json:"id" yaml:"id"`
	Body   string `json:"body" yaml:"body"`
	Author Author `json:"author" yaml:"author"`
}
