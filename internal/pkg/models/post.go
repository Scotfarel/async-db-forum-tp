package models

import (
	"time"
)

type Post struct {
	Author       string    `json:"author"`
	CreationDate time.Time `json:"created"`
	Forum        string    `json:"forum"`
	ID           uint64    `json:"id"`
	IsEdited     bool      `json:"isEdited"`
	Message      string    `json:"message"`
	ParentID     uint64    `json:"parent"`
	ThreadID     uint64    `json:"thread"`
	ForumID      uint64    `json:"-"`
	AuthorID     uint64    `json:"-"`
}

type PostFull struct {
	Author   *User   `json:"author"`
	Forum    *Forum  `json:"forum"`
	PostData *Post   `json:"post"`
	Thread   *Thread `json:"thread"`
}
