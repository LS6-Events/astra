package types

type Comment struct {
	ID     int    `json:"id" binding:"required"`
	Body   string `json:"body" binding:"required"`
	Author Author `json:"author" binding:"required"`
}
