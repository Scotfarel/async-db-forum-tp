package thread

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
)

type Usecase interface {
	InsertThreadInto(*models.Thread) (*models.Thread, error)
	CreateThreadPosts(string, []*models.Post) ([]*models.Post, error)
	GetThreadBySlugOrID(string) (*models.Thread, error)
	GetThreadPosts(string, uint64, uint64, string, bool) ([]*models.Post, error)
	UpdateThread(string, *models.Thread) (*models.Thread, error)
	InsertVote(string, *models.Vote) (*models.Thread, error)
}
