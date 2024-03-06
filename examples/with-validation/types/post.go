package types

import "time"

type Post struct {
	ID          int       `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at" binding:"required"`
	Author      Author    `json:"author" binding:"required"`
	Comments    []Comment `json:"comments" binding:"required,unique"`
}

type PostDTO struct {
	Name     string `json:"name" binding:"required"`
	Body     string `json:"body" binding:"required"`
	AuthorID int    `json:"author_id" binding:"required"`
}
