package types

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	gorm.Model
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	AuthorID    uint      `json:"-"`
	Author      Author    `json:"author"`
	Comments    []Comment `json:"comments"`
}

type PostDTO struct {
	Name     string `json:"name"`
	Body     string `json:"body"`
	AuthorID uint   `json:"author_id"`
}
