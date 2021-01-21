package forum

import "github.com/Scotfarel/db-tp-api/internal/pkg/models"

type Repository interface {
	InsertIntoForum(f *models.Forum) error
	GetForumBySlug(slug string) (*models.Forum, error)
}
