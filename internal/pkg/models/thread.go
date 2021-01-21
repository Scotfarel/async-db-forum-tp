package models

import (
	"time"
)

type Thread struct {
	Author       string    `json:"author"`
	AuthorID     uint64    `json:"-"`
	CreationDate time.Time `json:"created"`
	Forum        string    `json:"forum"`
	ForumID      uint64    `json:"-"`
	ID           uint64    `json:"id"`
	About        string    `json:"message"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	Votes        int64     `json:"votes"`
}
