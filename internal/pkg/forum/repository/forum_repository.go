package repository

import (
	"github.com/Scotfarel/db-tp-api/internal/pkg/forum"
	"github.com/Scotfarel/db-tp-api/internal/pkg/models"
	"github.com/Scotfarel/db-tp-api/internal/pkg/utils"

	"github.com/jackc/pgx"
)

type ForumRepository struct {
	db *pgx.ConnPool
}

func CreateRepository(db *pgx.ConnPool) forum.Repository {
	return &ForumRepository{
		db: db,
	}
}

func (repository *ForumRepository) InsertIntoForum(forum *models.Forum) error {
	if _, err := repository.db.Exec(
		"INSERT INTO forums (slug, admin, title) VALUES ($1, $2, $3)",
		forum.Slug,
		forum.AdminID,
		forum.Title,
	); err != nil {
		return err
	}

	return nil
}

func (repository *ForumRepository) GetForumBySlug(slug string) (*models.Forum, error) {
	oldForum := &models.Forum{}
	if err := repository.db.QueryRow(
		"SELECT f.id, f.slug, u.nickname, f.title, f.threads, f.posts FROM forums as f JOIN users as u ON (u.id = f.admin) WHERE lower(slug) = lower($1)",
		slug,
	).Scan(
		&oldForum.ID,
		&oldForum.Slug,
		&oldForum.AdminNickname,
		&oldForum.Title,
		&oldForum.ThreadsCount,
		&oldForum.PostsCount,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.ErrDoesntExists
		}
		return nil, err
	}

	return oldForum, nil
}
