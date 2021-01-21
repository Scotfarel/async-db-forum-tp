package user

import "github.com/Scotfarel/db-tp-api/internal/pkg/models"

type Repository interface {
	InsertUserInto(*models.User) error
	GetUserByNickname(string) (*models.User, error)
	CheckNicknames([]*models.Post) (bool, error)
	GetUserByEmail(string) (*models.User, error)
	GetUsersByForum(uint64, uint64, string, bool) ([]*models.User, error)
	UpdateUser(*models.User) error
}
