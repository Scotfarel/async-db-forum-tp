package forum

import "github.com/Scotfarel/db-tp-api/internal/pkg/models"

type Usecase interface {
	InsertIntoForum(*models.Forum) (*models.Forum, error)
	GetForumBySlug(string) (*models.Forum, error)
	GetForumUsers(string, uint64, string, bool) ([]*models.User, error)
	GetForumThreads(string, uint64, string, bool) ([]*models.Thread, error)
}
