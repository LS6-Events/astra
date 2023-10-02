package types

type Comment struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	Author Author `json:"author"`
}
