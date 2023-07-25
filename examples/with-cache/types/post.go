package types

import "time"

type Post struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	PublishedAt time.Time `json:"published_at"`
	Author      Author    `json:"author"`
	Comments    []Comment `json:"comments"`
}

type PostDTO struct {
	Name     string `json:"name"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}
