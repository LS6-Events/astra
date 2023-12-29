package types

type Comment struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	Author Author `json:"author"`
}

func (c Comment) MarshalJSON() ([]byte, error) {
	return []byte("\"This is a comment\""), nil
}
