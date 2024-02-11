package types

import "time"

type Post struct {
	ID          string    `json:"id" binding:"required,uuid4"`
	Name        string    `json:"name" binding:"required"`
	Body        string    `json:"body" binding:"max=512"`
	PublishedAt time.Time `json:"published_at" binding:"required"`
	Author      Author    `json:"author" binding:"required"`
	Comments    []Comment `json:"comments" binding:"required,unique"`
}

type PostDTO struct {
	Name     string `json:"name" binding:"required"`
	Body     string `json:"body" binding:"max=512"`
	AuthorID string `json:"author_id" binding:"required,uuid4"`
}
