package vote

import "github.com/Scotfarel/db-tp-api/internal/pkg/models"

type Repository interface {
	GetVotes(uint64) (int64, error)
	AddVote(*models.Vote) error
}
