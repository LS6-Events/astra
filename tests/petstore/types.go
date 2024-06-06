package petstore

// Tag the tag model.
type Tag struct {
	ID   string `json:"id" binding:"required,uuid4"`
	Name string `json:"name" binding:"required"`
}

// Pet the pet model.
type Pet struct {
	ID        string   `json:"id" binding:"required,uuid4"`
	Name      string   `json:"name" binding:"required"`
	PhotoURLs []string `json:"photoUrls,omitempty" binding:"dive,url"`
	Status    string   `json:"status,omitempty"`
	Tags      []Tag    `json:"tags,omitempty"`
}

// PetDTO the pet dto.
type PetDTO struct {
	Name      string   `json:"name"  binding:"required"`
	PhotoURLs []string `json:"photoUrls,omitempty" binding:"dive,url"`
	Status    string   `json:"status,omitempty"`
	Tags      []Tag    `json:"tags,omitempty"`
}
