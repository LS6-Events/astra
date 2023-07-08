package types

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID   uint   `json:"-"`
	Body     string `json:"body"`
	AuthorID uint   `json:"-"`
	Author   Author `json:"author"`
}
