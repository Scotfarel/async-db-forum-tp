package thread

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
)

type Repository interface {
	InsertThreadInto(*models.Thread) error
	GetThreadByID(uint64) (*models.Thread, error)
	GetThreadBySlug(string) (*models.Thread, error)
	GetThreadByForumSlug(string, uint64, string, bool) ([]*models.Thread, error)
	GetThreadCountByForumID(uint64) (uint64, error)
	UpdateThread(*models.Thread) error
}
