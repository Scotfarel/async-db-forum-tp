package models

type Forum struct {
	ID            uint64 `json:"-"`
	PostsCount    uint64 `json:"posts"`
	Slug          string `json:"slug"`
	ThreadsCount  uint64 `json:"threads"`
	Title         string `json:"title"`
	AdminNickname string `json:"user"`
	AdminID       uint64 `json:"-"`
}
