package models

type Status struct {
	ForumsCount  uint64 `json:"forum"`
	PostsCount   uint64 `json:"post"`
	ThreadsCount uint64 `json:"thread"`
	UsersCount   uint64 `json:"user"`
}
