package types

import "time"

type PostID struct {
	ID int `json:"id" yaml:"id" uri:"postID" binding:"required"`
}

type Post struct {
	PostID
	Name        string    `json:"name" yaml:"name"`
	Body        string    `json:"body" yaml:"body"`
	PublishedAt time.Time `json:"published_at" yaml:"publishedAt"`
	Author      Author    `json:"author" yaml:"author"`
	Comments    []Comment `json:"comments" yaml:"comments"`
}

type PostDTO struct {
	Name     string `json:"name" yaml:"name"`
	Body     string `json:"body" yaml:"body"`
	AuthorID int    `json:"author_id" yaml:"authorID"`
}
