package user

import "github.com/Scotfarel/db-tp-api/internal/pkg/models"

type Usecase interface {
	InsertUserInto(string, *models.User) ([]*models.User, error)
	GetUserByNickname(string) (*models.User, error)
	UpdateUser(string, *models.User) error
}
