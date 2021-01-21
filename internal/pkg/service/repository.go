package service

type Repository interface {
	ForumInfo() (uint64, error)
	PostInfo() (uint64, error)
	ThreadInfo() (uint64, error)
	UserInfo() (uint64, error)

	DropForum() error
	DropPost() error
	DropThread() error
	DropUser() error
	DropVotes() error
}
