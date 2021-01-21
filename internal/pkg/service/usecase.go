package service

import "github.com/Scotfarel/db-tp-api/internal/pkg/models"

type Usecase interface {
	GetInfo() (*models.Status, error)
	Drop() error
}
