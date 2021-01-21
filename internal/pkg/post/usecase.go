package post

import "github.com/Scotfarel/db-tp-api/internal/pkg/models"

type Usecase interface {
	GetPost(uint64, []string) (*models.PostFull, error)
	UpdatePost(*models.Post) (*models.Post, error)
}
