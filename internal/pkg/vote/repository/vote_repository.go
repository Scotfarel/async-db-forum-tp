package repository

import (
	"github.com/jackc/pgx"

	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/vote"
)

type VoteRepository struct {
	db *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) vote.Repository {
	return &VoteRepository{
		db,
	}
}

func (repository *VoteRepository) GetVotes(id uint64) (int64, error) {
	var votesNum int64
	if err := repository.db.QueryRow(
		"SELECT votes FROM threads WHERE id = $1",
		id,
	).Scan(&votesNum); err != nil {
		return 0, err
	}

	return votesNum, nil
}

func (repository *VoteRepository) AddVote(vote *models.Vote) error {
	voteID := uint64(0)
	err := repository.db.QueryRow(
		"SELECT id FROM votes WHERE thread = $1 AND author = $2",
		vote.ThreadID,
		vote.UserID,
	).Scan(&voteID)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	if err == pgx.ErrNoRows {
		if _, err := repository.db.Exec(
			"INSERT INTO votes (author, thread, vote) VALUES ($1, $2, $3)",
			vote.UserID,
			vote.ThreadID,
			vote.Voice,
		); err != nil {
			return err
		}

		return nil
	}

	if _, err = repository.db.Exec(
		"UPDATE votes SET vote = $2 WHERE id = $1",
		voteID,
		vote.Voice,
	); err != nil {
		return err
	}

	return nil
}
