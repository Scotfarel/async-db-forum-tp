package post

import "github.com/Scotfarel/db-tp-api/internal/pkg/models"

type Repository interface {
	InsertIntoPost([]*models.Post) error
	GetPostByThread(uint64, uint64, uint64, string, bool) ([]*models.Post, error)
	CheckPostParentPosts([]*models.Post, uint64) (bool, error)
	GetPostByID(uint64) (*models.Post, error)
	GetPostCountByForumID(uint64) (uint64, error)
	UpdatePost(*models.Post) error
}
