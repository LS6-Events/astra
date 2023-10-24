package types

type Author struct {
	ID        int    `json:"id" yaml:"id"`
	FirstName string `json:"first_name" yaml:"firstName"`
	LastName  string `json:"last_name" yaml:"lastName"`
}
