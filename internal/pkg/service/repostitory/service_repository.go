package repostitory

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/service"

	"github.com/jackc/pgx"
)

type Repository struct {
	db *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) service.Repository {
	return &Repository{
		db,
	}
}

func (repository *Repository) ForumInfo() (uint64, error) {
	var info uint64
	if err := repository.db.QueryRow("SELECT count(*) from forums").Scan(&info); err != nil {
		return 0, err
	}
	return info, nil
}

func (repository *Repository) PostInfo() (uint64, error) {
	var info uint64
	if err := repository.db.QueryRow("SELECT count(*) from posts").Scan(&info); err != nil {
		return 0, err
	}
	return info, nil
}

func (repository *Repository) ThreadInfo() (uint64, error) {
	var info uint64
	if err := repository.db.QueryRow("SELECT count(*) from threads").Scan(&info); err != nil {
		return 0, err
	}
	return info, nil
}

func (repository *Repository) UserInfo() (uint64, error) {
	var info uint64
	if err := repository.db.QueryRow("SELECT count(*) from users").Scan(&info); err != nil {
		return 0, err
	}
	return info, nil
}

func (repository *Repository) DropForum() error {
	if _, err := repository.db.Exec("TRUNCATE TABLE forums CASCADE"); err != nil {
		return err
	}
	return nil
}

func (repository *Repository) DropPost() error {
	if _, err := repository.db.Exec("TRUNCATE TABLE posts CASCADE"); err != nil {
		return err
	}
	return nil
}

func (repository *Repository) DropThread() error {
	if _, err := repository.db.Exec("TRUNCATE TABLE threads CASCADE"); err != nil {
		return err
	}
	return nil
}

func (repository *Repository) DropUser() error {
	if _, err := repository.db.Exec("TRUNCATE TABLE users CASCADE"); err != nil {
		return err
	}
	return nil
}

func (repository *Repository) DropVotes() error {
	if _, err := repository.db.Exec("TRUNCATE TABLE votes CASCADE"); err != nil {
		return err
	}
	return nil
}
