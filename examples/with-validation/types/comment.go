package types

type Comment struct {
	ID     string `json:"id" binding:"required,uuid4"`
	Body   string `json:"body" binding:"required,max=256"`
	Author Author `json:"author" binding:"required"`
}
