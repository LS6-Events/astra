package types

type Author struct {
	ID        string `json:"id" binding:"required,uuid4"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}
